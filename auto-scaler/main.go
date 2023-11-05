package main

import (
	"time"

	"github.com/aayushjn/auto-scaler/scaler"
)

func main() {
	autoScaler, err := scaler.NewAutoScaler(10*time.Second, "service-1", "python:latest", "load-balancer:latest")
	if err != nil {
		panic(err)
	}
	autoScaler.Monitor()
	autoScaler.Stop()
}
