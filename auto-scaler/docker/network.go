package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetDockerNetworkId(cli *client.Client, name string) (string, error) {
	ctx := context.Background()

	resp, err := cli.NetworkInspect(ctx, name, types.NetworkInspectOptions{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func CreateDockerNetwork(cli *client.Client, name string) (string, error) {
	ctx := context.Background()

	resp, err := cli.NetworkCreate(ctx, name, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}
