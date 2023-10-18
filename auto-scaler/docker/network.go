package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func CreateDockerNetwork(cli *client.Client, name string) error {
	ctx := context.Background()

	_, err := cli.NetworkCreate(ctx, name, types.NetworkCreate{})
	if err != nil {
		return err
	}
	return nil
}
