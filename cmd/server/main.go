package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

var version string

func main() {
	if err := run(context.Background()); err != nil {
		slog.Error("failed to start main", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
func run(runCtx context.Context) error {
	ctx, cancel := context.WithCancel(runCtx)

	go stopSignalNotify(ctx, cancel)

	return nil
}

func stopSignalNotify(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	reason := <-c
	slog.InfoContext(ctx, "program stopped", slog.String("reason", reason.String()))
	cancel()
}
