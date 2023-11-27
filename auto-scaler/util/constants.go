package util

import "time"

type key int

const (
	MinBackends             = 2
	MaxBackends             = 10
	DefaultPort             = 8081
	DefaultShutdownTimeout  = 10 * time.Second
	DefaultMonitoringPeriod = 10 * time.Second

	ConfigFailure     = 10
	AutoScalerFailure = 11
)
