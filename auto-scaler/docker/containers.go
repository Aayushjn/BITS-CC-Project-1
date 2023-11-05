package docker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type Container struct {
	ID        string
	IPAddress string
}

type ContainerStats struct {
	CPU    float64
	Memory float64
}

func CreateAndStartDockerContainer(cli *client.Client, containerImage, containerName, networkName string, startArgs ...string) (Container, error) {
	ctx := context.Background()
	containerConf := &container.Config{
		Hostname:   containerName,
		Image:      containerImage,
		Entrypoint: strslice.StrSlice(startArgs),
	}
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	_, err := cli.ContainerCreate(ctx, containerConf, hostConfig, nil, nil, containerName)
	if err != nil {
		return Container{}, fmt.Errorf("failed to create container: %w", err)
	}
	err = cli.NetworkConnect(ctx, networkName, containerName, nil)
	if err != nil {
		return Container{}, fmt.Errorf("failed to connect container to network: %w", err)
	}

	err = StartDockerContainer(cli, containerName)
	if err != nil {
		return Container{}, err
	}

	return GetDockerContainer(cli, containerName, networkName)
}

func StopAndRemoveDockerContainer(cli *client.Client, containerId string) error {
	err := StopDockerContainer(cli, containerId)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = cli.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}
	return nil
}

func StartDockerContainer(cli *client.Client, containerId string) error {
	ctx := context.Background()
	err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil
}

func StopDockerContainer(cli *client.Client, containerId string) error {
	ctx := context.Background()
	err := cli.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

func GetDockerContainer(cli *client.Client, containerName, networkName string) (Container, error) {
	ctx := context.Background()
	cData := Container{}
	inspectResp, err := cli.ContainerInspect(ctx, containerName)
	if err != nil {
		return cData, fmt.Errorf("failed to inspect container: %w", err)
	}

	cData.ID = inspectResp.ID
	cData.IPAddress = inspectResp.NetworkSettings.Networks[networkName].IPAddress
	return cData, nil
}

func GetDockerContainerStats(cli *client.Client, containerName string) (ContainerStats, error) {
	ctx := context.Background()
	cStats := ContainerStats{}

	stats, err := cli.ContainerStats(ctx, containerName, false)
	if err != nil {
		return cStats, err
	}

	defer stats.Body.Close()

	var dockerStats types.Stats
	err = json.NewDecoder(stats.Body).Decode(&dockerStats)
	if err != nil {
		return cStats, err
	}

	cacheStats, ok := dockerStats.MemoryStats.Stats["cache"]
	if !ok {
		cacheStats = 0
	}
	usedMemory := dockerStats.MemoryStats.Usage - cacheStats
	cpuDelta := dockerStats.CPUStats.CPUUsage.TotalUsage - dockerStats.PreCPUStats.CPUUsage.TotalUsage
	sysCpuDelta := dockerStats.CPUStats.SystemUsage - dockerStats.PreCPUStats.SystemUsage
	numCpus := dockerStats.CPUStats.OnlineCPUs
	if numCpus == 0 {
		numCpus = dockerStats.PreCPUStats.OnlineCPUs
	}
	if numCpus == 0 {
		numCpus = uint32(len(dockerStats.CPUStats.CPUUsage.PercpuUsage))
	}

	cStats.CPU = (float64(cpuDelta) / float64(sysCpuDelta)) * float64(numCpus) * 100.0
	cStats.Memory = (float64(usedMemory) / float64(dockerStats.MemoryStats.Limit)) * 100.0

	return cStats, nil
}
