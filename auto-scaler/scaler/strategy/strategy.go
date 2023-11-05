package strategy

import "github.com/aayushjn/auto-scaler/docker"

type Strategy interface {
	AddMeasurement(docker.ContainerStats)
	ResetMeasurements()
	AnalyzeAndPlan(int) int
}
