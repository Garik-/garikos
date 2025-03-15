package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
	"golang.org/x/sync/errgroup"
)

// set from ldflags.
var (
	Version = "" //nolint:gochecknoglobals
)

const (
	defaultHTTPAddr   = ":8002"
	readHeaderTimeout = 2 * time.Second
	defaultInterval   = 5 * time.Second // Значение по умолчанию
)

func newServer(ctx context.Context, addr string) *http.Server {
	return &http.Server{
		ReadHeaderTimeout: readHeaderTimeout,
		Addr:              addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
}

type Response struct {
	CPU     []float64                 `json:"cpu"`
	MEM     *mem.VirtualMemoryStat    `json:"mem"`
	Sensors []sensors.TemperatureStat `json:"sensors"`
}

func newResponse(ctx context.Context, interval time.Duration) (*Response, error) {
	var err error

	res := &Response{}
	res.CPU, err = cpu.PercentWithContext(ctx, interval, false)

	if err != nil {
		return nil, fmt.Errorf("PercentWithContext: %w", err)
	}

	res.MEM, err = mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("VirtualMemoryWithContext: %w", err)
	}

	res.Sensors, err = sensors.TemperaturesWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("TemperaturesWithContext: %w", err)
	}

	return res, nil
}

func parseInterval(value string) time.Duration {
	if value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}

	return defaultInterval
}

func newJSONEncoder(w http.ResponseWriter, r *http.Request) (*json.Encoder, func() error) {
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gw := gzip.NewWriter(w)

		return json.NewEncoder(gw), func() error {
			return gw.Close()
		}
	}

	return json.NewEncoder(w), func() error {
		return nil
	}
}

func diskHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		usage, err := disk.UsageWithContext(ctx, "/")

		if err != nil {
			logger.ErrorContext(ctx, "UsageWithContext", slog.String("error", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		encoder, encoderClose := newJSONEncoder(w, r)
		defer func() {
			if err := encoderClose(); err != nil {
				logger.ErrorContext(ctx, "encoderClose", slog.String("error", err.Error()))
			}
		}()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=60")

		w.WriteHeader(http.StatusOK)

		err = encoder.Encode(usage)
		if err != nil {
			logger.ErrorContext(ctx, "Encode", slog.String("error", err.Error()))

			return
		}
	}
}

func sendEvent(w io.Writer, encoder *json.Encoder, flusher http.Flusher, res *Response) error {
	_, err := io.WriteString(w, "data: ")
	if err != nil {
		return fmt.Errorf("writeString: %w", err)
	}

	err = encoder.Encode(res)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	_, err = io.WriteString(w, "\n\n")
	if err != nil {
		return fmt.Errorf("writeString: %w", err)
	}

	flusher.Flush()

	return nil
}

func systemHandler(logger *slog.Logger) http.HandlerFunc {
	var lastResponse atomic.Value
	lastResponse.Store(nil)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		flusher, ok := w.(http.Flusher)
		if !ok {
			logger.ErrorContext(ctx, "streaming unsupported")
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)

			return
		}

		encoder := json.NewEncoder(w)
		interval := parseInterval(r.URL.Query().Get("interval"))
		logger.DebugContext(ctx, "cpu handler", slog.Duration("interval", interval))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		if res := lastResponse.Load().(*Response); res != nil { // TODO: check TTL
			err := sendEvent(w, encoder, flusher, res)
			if err != nil {
				logger.ErrorContext(ctx, "sendEvent", slog.String("error", err.Error()))

				return
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			res, err := newResponse(ctx, interval)
			if err != nil {
				logger.ErrorContext(ctx, "newResponse", slog.String("error", err.Error()))

				return
			}

			err = sendEvent(w, encoder, flusher, res)
			if err != nil {
				logger.ErrorContext(ctx, "sendEvent", slog.String("error", err.Error()))

				return
			}

			lastResponse.Store(res)
		}
	}
}

func initServer(ctx context.Context, logger *slog.Logger, addr string) *http.Server {
	mux := http.NewServeMux()

	srv := newServer(ctx, addr)
	srv.Handler = mux

	mux.Handle("/system", systemHandler(logger))
	mux.Handle("/disk", diskHandler(logger))

	return srv
}

func main() {
	addr := flag.String("addr", defaultHTTPAddr, "HTTP server address")
	showVersion := flag.Bool("v", false, "Show version information")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if *showVersion {
		logger.InfoContext(ctx, Version)

		return
	}

	ctxSignal, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.InfoContext(ctxSignal, "starting server...", slog.String("addr", *addr))

	server := initServer(ctxSignal, logger, *addr)

	g, gCtx := errgroup.WithContext(ctxSignal)

	g.Go(func() error {
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		return server.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.ErrorContext(ctxSignal, "errgroup.Wait", slog.Any("error", err))
	}

	logger.InfoContext(ctxSignal, "done")
}
