package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func CreateDockerNetwork(cli *client.Client, name string) error {
	ctx := context.Background()

	_, err := cli.NetworkInspect(ctx, name, types.NetworkInspectOptions{})
	if err == nil {
		// since Docker network already exists, there is no need to create a new one
		return nil
	}

	_, err = cli.NetworkCreate(ctx, name, types.NetworkCreate{})
	if err != nil {
		return fmt.Errorf("failed to create network: %w", err)
	}
	return nil
}
