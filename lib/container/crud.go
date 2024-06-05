package container

import (
	"context"
	dockerContainerTypes "github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"	
	"github.com/hahahannes/e2e-go-utils/lib"
)

func ContainerExists(containerID string, ctx context.Context) bool {
	checkFunc := func() bool {
		cli, err := docker.NewClientWithOpts(
			docker.FromEnv, docker.WithAPIVersionNegotiation(),
		)
		if err != nil {
			panic(err)
		}
	
		containers, err := cli.ContainerList(ctx, dockerContainerTypes.ListOptions{
			All: true,
		})

		if err != nil {
			return false 
		}

		for _, container := range containers {
			if container.ID == containerID {
				return true
			}
		}
		return false
	}

	containerExists := lib.CheckConditionWithRetry(checkFunc, 3, 10)
	return containerExists
}

func ContainerIsRemoved(containerID string, ctx context.Context) bool {
	checkFunc := func() bool {
		cli, err := docker.NewClientWithOpts(
			docker.FromEnv, docker.WithAPIVersionNegotiation(),
		)
		if err != nil {
			panic(err)
		}
	
		containers, err := cli.ContainerList(ctx, dockerContainerTypes.ListOptions{
			All: true,
		})

		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			if container.ID == containerID {
				return false
			}
		}
		return true
	}

	containerExists := lib.CheckConditionWithRetry(checkFunc, 3, 10)
	return containerExists
}