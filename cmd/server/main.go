package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/volodya-nrg/mm/internal/adapters/logger"
	"github.com/volodya-nrg/mm/internal/configs"
	"github.com/volodya-nrg/mm/pkg/mylib"
)

var version string

func main() {
	configFilepath := flag.String("config", "", "config filepath")
	flag.Parse()

	if err := run(context.Background(), *configFilepath); err != nil {
		slog.Error("failed to start main", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
func run(runCtx context.Context, configFilepath string) error {
	ctx, cancel := context.WithCancel(runCtx)

	go stopSignalNotify(ctx, cancel)

	cfg := configs.NewConfig(configFilepath)

	if err := cfg.AttachData(); err != nil {
		return fmt.Errorf("failed to attach config-data: %s", err)
	}
	if version != "" {
		cfg.SetVersion(version)
	}

	log, err := logger.InitSlog(cfg.ServiceName, cfg.Version, cfg.Log.Level, cfg.Log.Path)
	if err != nil {
		return fmt.Errorf("failed to init slog: %w", err)
	}

	defer func() {
		if err = log.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close logger", slog.String("error", err.Error()))
		} else {
			slog.InfoContext(ctx, "logger closed successfully")
		}
	}()

	myLib, err := mylib.NewMyLib(
		cfg.SearchFolder,
		cfg.DegreeOfParallelism,
		mylib.NewNodeGetMD5(),
	)
	if err != nil {
		return fmt.Errorf("failed to init mylib: %w", err)
	}

	for resp := range myLib.Run(ctx) {
		slog.InfoContext(
			ctx,
			"got response",
			slog.String("name", resp.Name),
			slog.String("result", resp.Result),
			slog.Any("error", resp.Err),
		)
	}

	return nil
}

func stopSignalNotify(ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	reason := <-c
	slog.InfoContext(ctx, "program stopped", slog.String("reason", reason.String()))
	cancel()
}
