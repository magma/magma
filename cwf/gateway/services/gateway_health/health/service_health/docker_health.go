/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package service_health

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const (
	dockerRequestTimeout = 3 * time.Second
)

// DockerServiceHealthProvider provides service health for
// docker containers through docker's API.
type DockerServiceHealthProvider struct {
	dockerClient *client.Client
}

// NewDockerServiceHealthProvider creates a new DockerServiceHealthProvider
// with an initialized docker client.
func NewDockerServiceHealthProvider() (*DockerServiceHealthProvider, error) {
	dockercli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &DockerServiceHealthProvider{
		dockerClient: dockercli,
	}, nil
}

// GetUnhealthyServices returns all docker services failing their health checks.
func (d *DockerServiceHealthProvider) GetUnhealthyServices() ([]string, error) {
	filter := filters.NewArgs()
	filter.Add("health", types.Unhealthy)
	unhealthyFilter := types.ContainerListOptions{
		Filters: filter,
	}
	unhealthyContainers, err := d.dockerClient.ContainerList(context.Background(), unhealthyFilter)
	if err != nil {
		return []string{}, err
	}
	var unhealthyServices []string
	for _, container := range unhealthyContainers {
		if len(container.Names) == 0 {
			continue
		}
		unhealthyServices = append(unhealthyServices, container.Names[0])
	}
	return unhealthyServices, nil
}

// Enable restarts the service provided
func (d *DockerServiceHealthProvider) Enable(service string) error {
	sessiondID, err := d.getContainerID(service)
	if err != nil {
		return err
	}
	timeout := dockerRequestTimeout
	return d.dockerClient.ContainerRestart(context.Background(), sessiondID, &timeout)
}

// Disable stops the service provided
func (d *DockerServiceHealthProvider) Disable(service string) error {
	sessiondID, err := d.getContainerID(service)
	if err != nil {
		return err
	}
	timeout := dockerRequestTimeout
	return d.dockerClient.ContainerStop(context.Background(), sessiondID, &timeout)
}

func (d *DockerServiceHealthProvider) getContainerID(serviceName string) (string, error) {
	filter := filters.NewArgs()
	filter.Add("name", serviceName)
	sessiondContainerFilter := types.ContainerListOptions{
		Filters: filter,
	}
	containers, err := d.dockerClient.ContainerList(context.Background(), sessiondContainerFilter)
	if err != nil || len(containers) == 0 {
		return "", err
	}
	return containers[0].ID, nil
}
