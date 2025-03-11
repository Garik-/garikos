package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/shirou/gopsutil/v4/sensors"
)

// set from ldflags.
var (
	Version = "" //nolint:gochecknoglobals
)

const (
	defaultHTTPAddr = ":8002"
	shutdownTimeout = 2 * time.Second
)

func main() {
	addr := flag.String("addr", defaultHTTPAddr, "HTTP server address")
	showVersion := flag.Bool("v", false, "Show version information")
	flag.Parse()

	if *showVersion {
		log.Println(Version)
		os.Exit(0)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Starting server on", *addr)

	v, _ := mem.VirtualMemoryWithContext(ctx)
	log.Println(v)

	c, _ := cpu.InfoWithContext(ctx)
	log.Println(c)

	p, _ := cpu.PercentWithContext(ctx, time.Second*2, false)
	log.Println(p)

	s, _ := sensors.TemperaturesWithContext(ctx)
	log.Println(s)

	procs, _ := process.ProcessesWithContext(ctx)
	for _, proc := range procs {
		m, _ := proc.MemoryInfoWithContext(ctx)
		log.Println(m)
	}
}
