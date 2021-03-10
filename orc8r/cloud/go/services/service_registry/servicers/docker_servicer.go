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

package servicers

import (
	"fmt"
	"strings"

	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"golang.org/x/net/context"
)

const (
	// Use custom label to select only orc8r applications
	orc8rServicePartOfLabel         = "part-of=orc8r"
	dockerComposeServiceLabelPrefix = "com.docker.compose.service"
)

type DockerClient interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
}

type DockerServiceRegistryServicer struct {
	dockerClient DockerClient
}

// NewDockerServiceRegistryServicer creates a new service registry servicer
// that is backed by Docker.
func NewDockerServiceRegistryServicer(dockerCli DockerClient) *DockerServiceRegistryServicer {
	return &DockerServiceRegistryServicer{
		dockerClient: dockerCli,
	}
}

// ListAllServices returns the service name of all services in the registry.
func (s *DockerServiceRegistryServicer) ListAllServices(ctx context.Context, req *protos.Void) (*protos.ListAllServicesResponse, error) {
	labelFilter := filters.NewArgs()
	labelFilter.Add("label", orc8rServicePartOfLabel)
	listAllFilter := types.ContainerListOptions{
		All:     true,
		Filters: labelFilter,
	}
	allContainers, err := s.dockerClient.ContainerList(context.Background(), listAllFilter)
	if err != nil {
		return &protos.ListAllServicesResponse{}, err
	}
	services := []string{}
	for _, container := range allContainers {
		if len(container.Names) == 0 {
			continue
		}
		services = append(services, s.parseContainerNames(container.Names))
	}
	return &protos.ListAllServicesResponse{
		Services: services,
	}, nil
}

// FindServices returns all services in that have the provided label.
func (s *DockerServiceRegistryServicer) FindServices(ctx context.Context, req *protos.FindServicesRequest) (*protos.FindServicesResponse, error) {
	if req == nil {
		return &protos.FindServicesResponse{}, fmt.Errorf("FindServicesRequest is nil")
	}
	labelFilter := filters.NewArgs()
	formattedLabel := fmt.Sprintf("%s=1", req.GetLabel())
	labelFilter.Add("label", formattedLabel)
	labelFilter.Add("label", orc8rServicePartOfLabel)
	labelListOptions := types.ContainerListOptions{
		Filters: labelFilter,
		All:     true,
	}
	allContainers, err := s.dockerClient.ContainerList(context.Background(), labelListOptions)
	if err != nil {
		return &protos.FindServicesResponse{}, err
	}
	services := []string{}
	for _, container := range allContainers {
		if len(container.Names) == 0 {
			continue
		}
		services = append(services, s.parseContainerNames(container.Names))
	}
	return &protos.FindServicesResponse{
		Services: services,
	}, nil
}

// GetServiceAddress return the address of the gRPC server for the provided
// service.
func (s *DockerServiceRegistryServicer) GetServiceAddress(ctx context.Context, req *protos.GetServiceAddressRequest) (*protos.GetServiceAddressResponse, error) {
	if req == nil {
		return &protos.GetServiceAddressResponse{}, fmt.Errorf("ServiceRequest is nil")
	}
	exists, err := s.doesComposeServiceExist(req.GetService())
	if err != nil {
		return &protos.GetServiceAddressResponse{}, err
	}
	if !exists {
		return &protos.GetServiceAddressResponse{}, fmt.Errorf("Docker compose service '%s' was not found", req.GetService())
	}
	// Since orc8r mesh services are attached to the same docker network,
	// the services don't need to expose any ports. We can use a standardized
	// port for the gRPC service address.
	return &protos.GetServiceAddressResponse{
		Address: fmt.Sprintf("%s:%d", req.GetService(), registry.GrpcServicePort),
	}, nil
}

// GetHttpServerAddress returns the address of the HTTP server for the provided
// service.
func (s *DockerServiceRegistryServicer) GetHttpServerAddress(ctx context.Context, req *protos.GetHttpServerAddressRequest) (*protos.GetHttpServerAddressResponse, error) {
	if req == nil {
		return &protos.GetHttpServerAddressResponse{}, fmt.Errorf("ServiceRequest is nil")
	}
	exists, err := s.doesComposeServiceExist(req.GetService())
	if err != nil {
		return &protos.GetHttpServerAddressResponse{}, err
	}
	if !exists {
		return &protos.GetHttpServerAddressResponse{}, fmt.Errorf("Docker compose service '%s' was not found", req.GetService())
	}
	// Since orc8r mesh services are attached to the same docker network,
	// the services don't need to expose any ports. We can use a standardized
	// port for the HTTP server address.
	return &protos.GetHttpServerAddressResponse{
		Address: fmt.Sprintf("%s:%d", req.GetService(), registry.HttpServerPort),
	}, nil
}

// GetAnnotation returns the annotation value for the provided service and
// annotation.
func (s *DockerServiceRegistryServicer) GetAnnotation(ctx context.Context, req *protos.GetAnnotationRequest) (*protos.GetAnnotationResponse, error) {
	serviceFilter := filters.NewArgs()
	formattedLabel := fmt.Sprintf("%s=%s", dockerComposeServiceLabelPrefix, req.GetService())
	serviceFilter.Add("label", formattedLabel)
	serviceFilter.Add("label", orc8rServicePartOfLabel)
	containerFilter := types.ContainerListOptions{
		Filters: serviceFilter,
		All:     true,
	}
	containers, err := s.dockerClient.ContainerList(context.Background(), containerFilter)
	if err != nil {
		return &protos.GetAnnotationResponse{}, err
	} else if len(containers) == 0 {
		return &protos.GetAnnotationResponse{}, fmt.Errorf("No containers were found for service '%s'", req.GetService())
	}
	// Docker Compose does not differentiate between annotations and labels
	annotationValue, ok := containers[0].Labels[req.GetAnnotation()]
	if !ok {
		return &protos.GetAnnotationResponse{}, fmt.Errorf("Annotation '%s' does not exist for service '%s'", req.GetAnnotation(), req.GetService())
	}
	if len(annotationValue) == 0 {
		return &protos.GetAnnotationResponse{}, fmt.Errorf("Annotation value for annotation '%s' is empty", req.GetAnnotation())
	}
	return &protos.GetAnnotationResponse{
		AnnotationValue: annotationValue,
	}, nil
}

func (s *DockerServiceRegistryServicer) doesComposeServiceExist(serviceName string) (bool, error) {
	composeServiceFilter := filters.NewArgs()
	formattedLabel := fmt.Sprintf("%s=%s", dockerComposeServiceLabelPrefix, serviceName)
	composeServiceFilter.Add("label", formattedLabel)
	composeServiceFilter.Add("label", orc8rServicePartOfLabel)
	composeServiceContainerFilter := types.ContainerListOptions{
		Filters: composeServiceFilter,
		All:     true,
	}
	containers, err := s.dockerClient.ContainerList(context.Background(), composeServiceContainerFilter)
	if err != nil || len(containers) == 0 {
		return false, err
	}
	return true, nil
}

func (s *DockerServiceRegistryServicer) parseContainerNames(names []string) string {
	// Docker allows multiple names to be specified for a container, returning
	// the containers names in format ['/service1, /service_alias, ...]
	// Return the non-aliased container name, removing the preceding slash
	serviceName := names[0]
	return strings.TrimPrefix(serviceName, "/")
}
