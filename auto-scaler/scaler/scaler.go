package scaler

import (
	"fmt"
	"net/url"
	"time"

	"github.com/aayushjn/auto-scaler/docker"
	"github.com/aayushjn/auto-scaler/errors"
	"github.com/docker/docker/client"
)

type AutoScaler struct {
	dockerCli       *client.Client
	ticker          *time.Ticker
	loadBalancerUrl *url.URL
	backendMapping  map[string]string
	networkId       string
}

func (a *AutoScaler) Start() {
	for range a.ticker.C {
		// TODO: Implement monitoring
	}
}

func (a *AutoScaler) Stop() {
	a.ticker.Stop()
}

func (a *AutoScaler) Register(url string) error {
	if len(a.backendMapping) == 10 {
		return &errors.ErrBackendLimitOverflow{}
	}

	if _, found := a.backendMapping[url]; found {
		return &errors.ErrBackendAlreadyExists{BackendUrl: url}
	}
	// TODO: Implement registration
	return nil
}

func NewAutoScaler(frequency time.Duration, networkName, loadBalancerUrl string) (*AutoScaler, error) {
	parsed, err := url.Parse(loadBalancerUrl)
	if err != nil {
		return nil, err
	}

	cli, dockerErr := client.NewClientWithOpts(client.FromEnv)
	if dockerErr != nil {
		return nil, err
	}

	var networkId string
	networkId, err = docker.GetDockerNetworkId(cli, networkName)
	if err != nil {
		networkId, err = docker.CreateDockerNetwork(cli, networkName)
		if err != nil {
			return nil, fmt.Errorf("failed to create docker network: %w", err)
		}
	}

	return &AutoScaler{
		dockerCli:       cli,
		ticker:          time.NewTicker(frequency),
		loadBalancerUrl: parsed,
		backendMapping:  make(map[string]string, 10),
		networkId:       networkId,
	}, nil
}
