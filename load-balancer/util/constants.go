package util

import "time"

type key int

const (
	Attempts key = iota
	Retries
)

const (
	MinBackends               = 2
	MaxBackends               = 10
	DefaultPort               = 8080
	DefaultShutdownTimeout    = 10 * time.Second
	DefaultHealthCheckTimeout = 10 * time.Second
	DefaultBalancingStrategy  = "round_robin"
)
