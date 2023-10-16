package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aayushjn/load-balancer/balancer"
	"github.com/aayushjn/load-balancer/balancer/strategy"
	"github.com/aayushjn/load-balancer/server"
	"github.com/aayushjn/load-balancer/util"
	"golang.org/x/sync/errgroup"
)

var port int
var shutdownTimeout time.Duration
var healthCheckTimeout time.Duration
var balancingStrategy string

func main() {
	flag.IntVar(&port, "port", util.DefaultPort, "port to run server on")
	flag.DurationVar(&shutdownTimeout, "shutdown-timeout", util.DefaultShutdownTimeout, "time to wait before forcefully shutting down server")
	flag.DurationVar(&healthCheckTimeout, "health-check-timeout", util.DefaultHealthCheckTimeout, "frequency of health checks")
	flag.StringVar(&balancingStrategy, "balancing-strategy", util.DefaultBalancingStrategy, "balancing strategy to use")
	flag.Parse()

	if port <= 0 || port > 65535 {
		fmt.Println("invalid port")
		os.Exit(1)
	}
	ok := false
	for _, strat := range strategy.AllowedStrategies {
		if balancingStrategy == strat {
			ok = true
			break
		}
	}
	if !ok {
		fmt.Printf("invalid balancing strategy, must be one of %v\n", strategy.AllowedStrategies)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	g, gCtx := errgroup.WithContext(ctx)

	lb := balancer.NewLoadBalancer(strategy.NewLoadBalancingStrategy(balancingStrategy))

	srv := server.NewServer(lb, port)

	go func() {
		ticker := time.NewTicker(healthCheckTimeout)
		for {
			select {
			case <-ticker.C:
				fmt.Println("Starting health check...")
				lb.HealthCheck()
				fmt.Println("Health check completed")
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	})

	err := g.Wait()

	if err != nil && err != context.Canceled {
		fmt.Println(err)
	}
}
