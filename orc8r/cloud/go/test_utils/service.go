/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test_utils

import (
	"fmt"
	"net"
	"testing"

	cloud_service "magma/orc8r/cloud/go/service"
	"magma/orc8r/lib/go/registry"
	platform_service "magma/orc8r/lib/go/service"
)

// NewTestService creates and registers a basic test Magma service on a
// dynamically selected available local port.
// Returns the newly created service and listener it was registered with.
func NewTestService(t *testing.T, moduleName string, serviceType string) (*platform_service.Service, net.Listener) {
	srvPort, lis, err := getOpenPort()
	if err != nil {
		t.Fatal(err)
	}

	registry.AddService(registry.ServiceLocation{Name: serviceType, Host: "localhost", Port: srvPort})

	srv, err := cloud_service.NewTestService(t, moduleName, serviceType)
	if err != nil {
		t.Fatalf("Error creating service: %s", err)
	}
	return srv, lis
}

// NewTestOrchestratorService creates and registers a test orchestrator service
// on a dynamically selected available local port for the gRPC server and HTTP
// echo server. Returns the newly created service and listener it was
// registered with.
func NewTestOrchestratorService(
	t *testing.T,
	moduleName string,
	serviceType string,
	labels map[string]string,
	annotations map[string]string,
) (*cloud_service.OrchestratorService, net.Listener) {
	if labels == nil {
		labels = map[string]string{}
	}
	if annotations == nil {
		annotations = map[string]string{}
	}

	srvPort, lis, err := getOpenPort()
	if err != nil {
		t.Fatal(err)
	}
	echoPort, echoLis, err := getOpenPort()
	if err != nil {
		t.Fatal(err)
	}
	err = echoLis.Close()
	if err != nil {
		t.Fatal(err)
	}

	location := registry.ServiceLocation{
		Name:        serviceType,
		Host:        "localhost",
		EchoPort:    echoPort,
		Port:        srvPort,
		Labels:      labels,
		Annotations: annotations,
	}
	registry.AddService(location)

	srv, err := cloud_service.NewTestOrchestratorService(t, moduleName, serviceType)
	if err != nil {
		t.Fatalf("Error creating service: %s", err)
	}
	return srv, lis
}

func getOpenPort() (int, net.Listener, error) {
	lis, err := net.Listen("tcp", "")
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create listener: %s", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", lis.Addr().String())
	if err != nil {
		return 0, nil, fmt.Errorf("failed to resolve TCP address: %s", err)
	}
	return addr.Port, lis, err
}
