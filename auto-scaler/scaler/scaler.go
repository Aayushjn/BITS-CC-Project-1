package scaler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/aayushjn/auto-scaler/docker"
	"github.com/aayushjn/auto-scaler/errors"
	"github.com/aayushjn/auto-scaler/scaler/strategy"
	"github.com/aayushjn/auto-scaler/util/net"
	"github.com/docker/docker/client"
	"github.com/rs/xid"
)

type AutoScaler struct {
	dockerCli       *client.Client
	serviceImage    string
	serviceName     string
	ticker          *time.Ticker
	monitorTicker   *time.Ticker
	strategy        strategy.Strategy
	logger          slog.Logger
	loadBalancerUrl string
	backendMapping  map[string]string
	networkName     string
	pauseMonitoring bool
}

func (a *AutoScaler) Monitor() {
	for range a.monitorTicker.C {
		if a.pauseMonitoring {
			continue
		}
		stats := docker.ContainerStats{CPU: 0.0, Memory: 0.0}
		for _, containerId := range a.backendMapping {
			a.logger.Debug("Monitoring container ", containerId, " ...")
			currentStats, err := docker.GetDockerContainerStats(a.dockerCli, containerId)
			if err != nil {
				a.logger.Error(err.Error())
				continue
			}
			stats.CPU += currentStats.CPU
			stats.Memory += currentStats.Memory
		}
		numBackends := len(a.backendMapping)
		stats.CPU /= float64(numBackends)
		stats.Memory /= float64(numBackends)
		a.strategy.AddMeasurement(stats)
	}
	a.pauseMonitoring = true
	delta := a.strategy.AnalyzeAndPlan(len(a.backendMapping))
	if delta == 0 {
		a.logger.Info("No need to scale")
	} else if delta > 0 {
		a.logger.Info("Scaling up by ", fmt.Sprint(delta), " containers")
		for i := 1; i <= delta; i++ {
			err := a.ScaleUp()
			a.logger.Error(err.Error())
		}
	}
}

func (a *AutoScaler) Stop() {
	a.monitorTicker.Stop()
	a.ticker.Stop()
}

func (a *AutoScaler) ScaleUp() error {
	if len(a.backendMapping) == 10 {
		return &errors.ErrBackendLimitOverflow{}
	}

	var err error
	var container docker.Container

	container, err = docker.CreateAndStartDockerContainer(
		a.dockerCli,
		a.serviceImage,
		fmt.Sprintf("%s-%s", a.serviceName, xid.New().String()),
		a.networkName,
	)
	if err != nil {
		return err
	}

	url := "http://" + container.IPAddress
	var resp *http.Response
	var jsonBody []byte
	body := map[string]any{
		"url": url,
	}

	jsonBody, err = json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to serialize body: %w", err)
	}
	resp, err = net.MakeRequest("POST", a.loadBalancerUrl+"/_api/backend/", jsonBody)
	if err != nil {
		return fmt.Errorf("failed to register instance: %w", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to register instance: %w", err)
	}
	a.backendMapping[url] = container.ID

	return nil
}

func NewAutoScaler(frequency time.Duration, serviceName, serviceImage, loadBalancerImage string) (*AutoScaler, error) {
	cli, dockerErr := client.NewClientWithOpts(client.FromEnv)
	if dockerErr != nil {
		return nil, fmt.Errorf("failed to initialize Docker client: %w", dockerErr)
	}
	networkName := "net-" + serviceName
	err := docker.CreateDockerNetwork(cli, networkName)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker network: %w", err)
	}

	var cData docker.Container
	cData, err = docker.GetDockerContainer(cli, "lb-"+serviceName, networkName)
	if err != nil {
		cData, err = docker.CreateAndStartDockerContainer(cli, loadBalancerImage, fmt.Sprintf("lb-%s", serviceName), networkName)
		if err != nil {
			return nil, fmt.Errorf("failed to create and start docker container: %w", err)
		}
	}

	return &AutoScaler{
		dockerCli:       cli,
		ticker:          time.NewTicker(frequency),
		monitorTicker:   time.NewTicker(2 * time.Second),
		loadBalancerUrl: "http://" + cData.IPAddress + ":8080",
		backendMapping:  make(map[string]string, 10),
		networkName:     networkName,
		serviceImage:    serviceImage,
		serviceName:     serviceName,
		logger:          *slog.New(slog.NewTextHandler(os.Stdout, nil)),
		strategy:        strategy.NewThresholdStrategy(30.5, 40.5),
	}, nil
}
