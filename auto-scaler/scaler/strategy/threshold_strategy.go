package strategy

import (
	"math"

	"github.com/aayushjn/auto-scaler/docker"
)

type ThresholdStrategy struct {
	measurements    []docker.ContainerStats
	cpuThreshold    float64
	memoryThreshold float64

	counter int
}

func (s *ThresholdStrategy) AddMeasurement(stat docker.ContainerStats) {
	if s.counter == 10 {
		s.counter = 0
	}
	s.measurements[s.counter] = stat
	s.counter++
}

func (s *ThresholdStrategy) ResetMeasurements() {
	s.counter = 0
}

func (s *ThresholdStrategy) AnalyzeAndPlan(numBackends int) int {
	cpuAvg := 0.0
	memAvg := 0.0

	for i := 0; i < s.counter; i++ {
		cpuAvg += s.measurements[i].CPU
		memAvg += s.measurements[i].Memory
	}

	numStats := len(s.measurements)
	cpuAvg /= float64(numStats)
	memAvg /= float64(numStats)

	deltaCpu := int(math.Ceil(float64(numBackends) * (1.0 - (cpuAvg / s.cpuThreshold))))
	deltaMem := int(math.Ceil(float64(numBackends) * (1.0 - (memAvg / s.memoryThreshold))))

	if deltaCpu >= deltaMem {
		return deltaCpu
	} else {
		return deltaMem
	}
}

func NewThresholdStrategy(cpuThreshold, memoryThreshold float64) *ThresholdStrategy {
	return &ThresholdStrategy{
		measurements:    make([]docker.ContainerStats, 10),
		cpuThreshold:    cpuThreshold,
		memoryThreshold: memoryThreshold,
		counter:         0,
	}
}
