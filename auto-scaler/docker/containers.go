package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Container struct {
	ID        string
	IPAddress string
}

func CreateNewDockerContainer(cli *client.Client, containerImage, containerName, networkName string) (Container, error) {
	cData := Container{}
	ctx := context.Background()
	containerConf := &container.Config{
		Hostname: containerName,
		Image:    containerImage,
	}
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name:              "unless-stopped",
			MaximumRetryCount: 3,
		},
	}

	createResp, err := cli.ContainerCreate(ctx, containerConf, hostConfig, nil, nil, containerName)
	if err != nil {
		return cData, err
	}
	err = cli.NetworkConnect(ctx, networkName, createResp.ID, nil)
	if err != nil {
		return cData, err
	}

	inspectResp, err := cli.ContainerInspect(ctx, createResp.ID)
	if err != nil {
		return cData, err
	}

	cData.ID = inspectResp.ID
	cData.IPAddress = inspectResp.NetworkSettings.Networks[networkName].IPAddress
	return cData, nil
}

func RemoveDockerContainer(cli *client.Client, containerId string) error {
	ctx := context.Background()
	err := cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func StartDockerContainer(cli *client.Client, containerId string) error {
	ctx := context.Background()
	err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}

func StopDockerContainer(cli *client.Client, containerId string) error {
	ctx := context.Background()
	err := cli.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		return err
	}
	return nil
}
