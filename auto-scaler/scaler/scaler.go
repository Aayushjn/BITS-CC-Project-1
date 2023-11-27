package scaler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/aayushjn/auto-scaler/config"
	"github.com/aayushjn/auto-scaler/docker"
	scalerErrors "github.com/aayushjn/auto-scaler/errors"
	"github.com/aayushjn/auto-scaler/scaler/strategy"
	"github.com/aayushjn/auto-scaler/util"
	"github.com/aayushjn/auto-scaler/util/net"
	"github.com/docker/docker/client"
	"github.com/rs/xid"
)

type AutoScaler struct {
	dockerCli       *client.Client
	serviceConf     config.ServiceConfig
	ticker          *time.Ticker
	monitorTicker   *time.Ticker
	strategy        strategy.Strategy
	logger          slog.Logger
	loadBalancerUrl string
	backendMapping  map[string]string
	networkName     string
	pauseMonitoring bool
}

func (a *AutoScaler) DiscoverServices() error {
	for i := 0; i < util.MaxBackends; i++ {
		cData, err := docker.GetDockerContainer(
			a.dockerCli,
			fmt.Sprintf("%s-%s", a.serviceConf.Name, xid.New().String()),
			a.networkName,
		)
		if err != nil {
			continue
		}
		url := "http://" + cData.IPAddress
		var resp *http.Response
		var jsonBody []byte
		body := map[string]any{
			"url": url,
		}

		jsonBody, err = json.Marshal(body)
		if err != nil {
			continue
		}
		resp, err = net.MakeRequest("POST", a.loadBalancerUrl+"/_api/backend/", jsonBody)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusNoContent {
			continue
		}
		a.backendMapping[url] = cData.ID
	}

	if len(a.backendMapping) < 2 {
		err := a.ScaleUp(2 - len(a.backendMapping))
		if err != nil {
			return err
		}
	}
	return nil
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
		err := a.ScaleUp(delta)
		if err != nil {
			a.logger.Error(err.Error())
		}
	} else {
		a.logger.Info("Scaling down by ", fmt.Sprint(delta), " containers")
		err := a.ScaleDown(delta)
		if err != nil {
			a.logger.Error(err.Error())
		}
	}
}

func (a *AutoScaler) Stop() {
	a.monitorTicker.Stop()
	a.ticker.Stop()
}

func (a *AutoScaler) ScaleUp(count int) error {
	if len(a.backendMapping) == 10 {
		return &scalerErrors.ErrBackendLimitOverflow{}
	}

	var joinedErr error = nil
	var err error
	var container docker.Container

	for i := 0; i < count; i++ {
		container, err = docker.CreateAndStartDockerContainer(
			a.dockerCli,
			a.serviceConf.Image,
			fmt.Sprintf("%s-%s", a.serviceConf.Name, xid.New().String()),
			a.networkName,
		)
		if err != nil {
			joinedErr = errors.Join(joinedErr, err)
			continue
		}
		url := "http://" + container.IPAddress
		var resp *http.Response
		var jsonBody []byte
		body := map[string]any{
			"url": url,
		}

		jsonBody, err = json.Marshal(body)
		if err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to serialize body: %w", err))
			continue
		}
		resp, err = net.MakeRequest("POST", a.loadBalancerUrl+"/_api/backend/", jsonBody)
		if err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to register instance: %w", err))
			continue
		}
		if resp.StatusCode != http.StatusNoContent {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to register instance: %w", err))
			continue
		}
		a.backendMapping[url] = container.ID
	}

	return joinedErr
}
func (a *AutoScaler) ScaleDown(count int) error {
	if len(a.backendMapping) == 2 || count > 8 {
		return &scalerErrors.ErrBackendLimitUnderflow{}
	}

	var joinedErr error = nil

	stats := make(map[util.Entry[string, string]]docker.ContainerStats, len(a.backendMapping))
	keys := make([]util.Entry[string, string], 0, len(stats))
	for url, containerId := range a.backendMapping {
		stat, err := docker.GetDockerContainerStats(a.dockerCli, containerId)
		if err != nil {
			joinedErr = errors.Join(joinedErr, err)
			continue
		}
		entry := util.Entry[string, string]{Key: url, Value: containerId}
		stats[entry] = stat
		keys = append(keys, entry)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return stats[keys[i]].Norm() > stats[keys[j]].Norm()
	})
	keys = keys[:count]

	var err error
	for _, entry := range keys {
		url := "http://" + entry.Key
		var resp *http.Response
		var jsonBody []byte
		body := map[string]any{
			"url": url,
		}

		jsonBody, err = json.Marshal(body)
		if err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to serialize body: %w", err))
			continue
		}
		resp, err = net.MakeRequest("DELETE", a.loadBalancerUrl+"/_api/backend/", jsonBody)
		if err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to unregister instance: %w", err))
			continue
		}
		if resp.StatusCode != http.StatusNoContent {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to unregister instance: %w", err))
			continue
		}

		err := docker.StopAndRemoveDockerContainer(a.dockerCli, entry.Value)
		if err != nil {
			joinedErr = errors.Join(joinedErr, err)
			continue
		}
		delete(a.backendMapping, entry.Key)
	}

	return joinedErr
}

func NewAutoScaler(frequency time.Duration, conf config.Config) (*AutoScaler, error) {
	cli, dockerErr := client.NewClientWithOpts(client.FromEnv)
	if dockerErr != nil {
		return nil, fmt.Errorf("failed to initialize Docker client: %w", dockerErr)
	}

	cData, err := docker.GetDockerContainer(cli, "lb-"+conf.Service.Name, conf.Network)
	if err != nil {
		cData, err = docker.CreateAndStartDockerContainer(cli, conf.LoadBalancer.Image, fmt.Sprintf("lb-%s", conf.Service.Name), conf.Network)
		if err != nil {
			return nil, fmt.Errorf("failed to create and start docker container: %w", err)
		}
	}

	return &AutoScaler{
		dockerCli:       cli,
		ticker:          time.NewTicker(frequency),
		monitorTicker:   time.NewTicker(2 * time.Second),
		loadBalancerUrl: "http://" + cData.IPAddress + ":8080",
		backendMapping:  make(map[string]string, util.MaxBackends),
		networkName:     conf.Network,
		serviceConf:     conf.Service,
		logger:          *slog.New(slog.NewTextHandler(os.Stdout, nil)),
		strategy:        strategy.NewThresholdStrategy(30.5, 40.5),
	}, nil
}
