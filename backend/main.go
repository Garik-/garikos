package main

import (
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
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
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

func systemStatusHandler(logger *slog.Logger) http.HandlerFunc {
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

			_, err = io.WriteString(w, "data: ")
			if err != nil {
				logger.ErrorContext(ctx, "WriteString", slog.String("error", err.Error()))

				return
			}

			err = encoder.Encode(res)
			if err != nil {
				logger.ErrorContext(ctx, "Encode", slog.String("error", err.Error()))

				return
			}

			_, err = io.WriteString(w, "\n\n")
			if err != nil {
				logger.ErrorContext(ctx, "WriteString", slog.String("error", err.Error()))

				return
			}

			flusher.Flush()
		}
	}
}

func initServer(ctx context.Context, logger *slog.Logger, addr string) *http.Server {
	mux := http.NewServeMux()

	srv := newServer(ctx, addr)
	srv.Handler = mux

	mux.Handle("/system-status", systemStatusHandler(logger))

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
