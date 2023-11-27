package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/aayushjn/auto-scaler/config"
	"github.com/aayushjn/auto-scaler/scaler"
	"github.com/aayushjn/auto-scaler/util"
	"golang.org/x/sync/errgroup"
)

var configFile string

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	flag.StringVar(&configFile, "config-file", filepath.Join(cwd, "config.toml"), "path to config file")
	flag.Parse()

	var conf config.Config
	conf, err = config.LoadConfig(configFile)
	if err != nil {
		slog.Error("failed to load config: ", err)
		os.Exit(util.ConfigFailure)
	}

	autoScaler, err := scaler.NewAutoScaler(util.DefaultMonitoringPeriod, conf)
	if err != nil {
		slog.Error("failed to create auto scaler: ", err)
		os.Exit(util.AutoScalerFailure)
	}
	err = autoScaler.DiscoverServices()
	if err != nil {
		slog.Error("failed to discover services: ", err)
		os.Exit(util.AutoScalerFailure)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		autoScaler.Monitor()
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		_, cancel := context.WithTimeout(context.Background(), util.DefaultShutdownTimeout)
		defer cancel()
		autoScaler.Stop()
		return nil
	})

	err = g.Wait()

	if err != nil && err != context.Canceled {
		slog.Warn(err.Error())
	}
}
