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
	"context"
	"fmt"
	"testing"

	"magma/orc8r/lib/go/protos"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDockerListAllServices(t *testing.T) {
	dockerClient := &MockDockerClient{}
	servicer := NewDockerServiceRegistryServicer(dockerClient)
	req := &protos.Void{}
	ctx := context.Background()

	// Test happy path
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(getMockContainers(), nil).Once()
	response, err := servicer.ListAllServices(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.GetServices()))
	assert.Equal(t, []string{"service1", "service2"}, response.GetServices())

	// Test docker daemon unavailable
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{}, fmt.Errorf("docker daemon unavailable")).Once()
	_, err = servicer.ListAllServices(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.Error(t, err)
}

func TestDockerFindServices(t *testing.T) {
	dockerClient := &MockDockerClient{}
	servicer := NewDockerServiceRegistryServicer(dockerClient)
	req := &protos.FindServicesRequest{Label: "label1"}
	ctx := context.Background()

	// Test happy path
	containers := getMockContainers()
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(containers, nil).Once()
	response, err := servicer.FindServices(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service1", "service2"}, response.Services)

	req.Label = "label4"
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{containers[1]}, nil).Once()
	response, err = servicer.FindServices(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, []string{"service2"}, response.Services)

	// Test label not found
	req.Label = "label5"
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{}, nil).Once()
	response, err = servicer.FindServices(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response.Services))
}

func TestDockerGetServiceAddress(t *testing.T) {
	dockerClient := &MockDockerClient{}
	servicer := NewDockerServiceRegistryServicer(dockerClient)
	req := &protos.GetServiceAddressRequest{Service: "service1"}
	ctx := context.Background()

	// Test happy path
	containers := getMockContainers()
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(containers, nil).Once()
	response, err := servicer.GetServiceAddress(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, response.Address, "service1:9180")

	// Test non-existent service
	req.Service = "service3"
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{}, nil).Once()
	_, err = servicer.GetServiceAddress(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.Error(t, err)
}

func TestDockerGetHttpServerAddress(t *testing.T) {
	dockerClient := &MockDockerClient{}
	servicer := NewDockerServiceRegistryServicer(dockerClient)
	req := &protos.GetHttpServerAddressRequest{Service: "service1"}
	ctx := context.Background()

	// Test happy path
	containers := getMockContainers()
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(containers, nil).Once()
	response, err := servicer.GetHttpServerAddress(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, response.Address, "service1:8080")

	// Test non-existent service
	req.Service = "service3"
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{}, nil).Once()
	_, err = servicer.GetHttpServerAddress(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.Error(t, err)

}

func TestDockerGetAnnotation(t *testing.T) {
	dockerClient := &MockDockerClient{}
	servicer := NewDockerServiceRegistryServicer(dockerClient)
	req := &protos.GetAnnotationRequest{Service: "service1", Annotation: "label1"}
	ctx := context.Background()

	// Test happy path
	containers := getMockContainers()
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return(containers, nil).Once()
	response, err := servicer.GetAnnotation(ctx, req)
	dockerClient.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, response.AnnotationValue, "foo")

	// Test non-existent label
	req.Annotation = "label4"
	dockerClient.On("ContainerList", mock.Anything, mock.Anything).Return([]types.Container{containers[0]}, nil).Once()
	_, err = servicer.GetAnnotation(ctx, req)
	assert.Error(t, err)
	dockerClient.AssertExpectations(t)
}

func getMockContainers() []types.Container {
	return []types.Container{
		{
			ID:     "111",
			Names:  []string{"service1"},
			Image:  "service_image",
			Status: "running",
			Labels: map[string]string{
				"label1": "foo",
				"label2": "bar",
			},
		},
		{
			ID:     "222",
			Names:  []string{"/service2"},
			Image:  "service_image",
			Status: "running",
			Labels: map[string]string{
				"label1": "roo",
				"label3": "baz",
				"label4": "sop",
			},
		},
	}
}

type MockDockerClient struct {
	mock.Mock
}

func (m *MockDockerClient) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	args := m.Called(ctx, options)
	return args.Get(0).([]types.Container), args.Error(1)
}
