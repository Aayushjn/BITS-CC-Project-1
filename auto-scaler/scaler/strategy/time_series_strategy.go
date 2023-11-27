package strategy

import (
	"math"

	"github.com/aayushjn/auto-scaler/docker"
)

type TimeSeriesStrategy struct {
	smoothed        docker.ContainerStats
	trend           docker.ContainerStats
	current         docker.ContainerStats
	alpha           float64
	gamma           float64
	cpuThreshold    float64
	memoryThreshold float64
}

func expSmoothing(alpha, current, smooth, trend float64) float64 {
	return (alpha * current) + ((1 - alpha) * (smooth - trend))
}

func (s *TimeSeriesStrategy) AddMeasurement(stat docker.ContainerStats) {
	prevSmooth := s.smoothed
	s.current = stat
	s.smoothed.CPU = expSmoothing(s.alpha, s.current.CPU, s.smoothed.CPU, s.trend.CPU)
	s.smoothed.Memory = expSmoothing(s.alpha, s.current.Memory, s.smoothed.Memory, s.trend.Memory)

	s.trend.CPU = (s.gamma * (s.smoothed.CPU - prevSmooth.CPU)) + ((1 - s.gamma) * s.trend.CPU)
	s.trend.Memory = (s.gamma * (s.smoothed.Memory - prevSmooth.Memory)) + ((1 - s.gamma) * s.trend.Memory)
}

func (s *TimeSeriesStrategy) ResetMeasurements() {}

func (s *TimeSeriesStrategy) AnalyzeAndPlan(numBackends int) int {
	predCpu := s.smoothed.CPU + s.trend.CPU
	predMem := s.smoothed.Memory + s.trend.Memory

	deltaCpu := int(math.Ceil(float64(numBackends) * (1.0 - (predCpu / s.cpuThreshold))))
	deltaMem := int(math.Ceil(float64(numBackends) * (1.0 - (predMem / s.memoryThreshold))))

	if deltaCpu >= deltaMem {
		return deltaCpu
	} else {
		return deltaMem
	}
}

func NewTimeSeriesStrategy(cpuThreshold, memoryThreshold float64) *TimeSeriesStrategy {
	return &TimeSeriesStrategy{
		alpha:           0.6,
		gamma:           0.2,
		cpuThreshold:    cpuThreshold,
		memoryThreshold: memoryThreshold,
	}
}
