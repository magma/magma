/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service_health

import (
	"context"
	"fmt"
	"strings"
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
		All:     true,
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
		serviceName := container.Names[0]
		if strings.HasPrefix(serviceName, "/") {
			serviceName = strings.ReplaceAll(serviceName, "/", "")
		}
		unhealthyServices = append(unhealthyServices, serviceName)
	}
	return unhealthyServices, nil
}

// Restart restarts the service provided
func (d *DockerServiceHealthProvider) Restart(service string) error {
	sessiondID, err := d.getContainerID(service)
	if err != nil {
		return err
	}
	timeout := dockerRequestTimeout
	return d.dockerClient.ContainerRestart(context.Background(), sessiondID, &timeout)
}

// Stop stops the service provided
func (d *DockerServiceHealthProvider) Stop(service string) error {
	serviceID, err := d.getContainerID(service)
	if err != nil {
		return err
	}
	timeout := dockerRequestTimeout
	return d.dockerClient.ContainerStop(context.Background(), serviceID, &timeout)
}

func (d *DockerServiceHealthProvider) getContainerID(serviceName string) (string, error) {
	filter := filters.NewArgs()
	fullName := fmt.Sprintf("/%s", serviceName)
	filter.Add("name", fullName)
	sessiondContainerFilter := types.ContainerListOptions{
		Filters: filter,
		All:     true,
	}
	containers, err := d.dockerClient.ContainerList(context.Background(), sessiondContainerFilter)
	if err != nil || len(containers) == 0 {
		return "", err
	}
	// There's a chance that search may returns multiple containers where
	// one service's name is a prefix of the other service.
	for _, svc := range containers {
		for _, name := range svc.Names {
			if name == fullName {
				return svc.ID, nil
			}
		}
	}
	return "", fmt.Errorf("Could not find containerID for service: %s", serviceName)
}
