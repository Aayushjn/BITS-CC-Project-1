package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CreateNewDockerContainer(cli *client.Client, containerImage, containerName, networkId string) (string, error) {
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

	resp, err := cli.ContainerCreate(ctx, containerConf, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	err = cli.NetworkConnect(ctx, networkId, resp.ID, nil)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
