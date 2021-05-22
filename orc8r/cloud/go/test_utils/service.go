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
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"magma/gateway/config"
	cloud_service "magma/orc8r/cloud/go/service"
	"magma/orc8r/lib/go/registry"
	platform_service "magma/orc8r/lib/go/service"
	platform_cfg "magma/orc8r/lib/go/service/config"
)

// control_proxy.yml
// only local_port and cloud_port are needed
var testCPConfigYaml = `
nghttpx_config_location: /var/tmp/nghttpx.conf
rootca_cert: /var/opt/magma/certs/rootCA.pem
gateway_cert: /var/opt/magma/secrets/certs/gateway.crt
gateway_key: /var/opt/magma/secrets/certs/gateway.key

local_port: %s
cloud_address: controller.magma.foobar.com
cloud_port: 9999

bootstrap_address: bootstrapper-controller.magma.foobar.com
bootstrap_port: 4444
proxy_cloud_connections: True
allow_http_proxy: True`

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

// NewTestOrchestratorServiceWithControlProxy create a Orchestrator Service and a
// control_proxy.yml file. This service can be used by gateway services to test
// a cloud service. This service will configure a custom control_proxy.yml file
// matching local_port on control proxy with the listener port of the orc8r service.
// Remember to delete the temporary file once the test it is done os.RemoveAll(dir)
func NewTestOrchestratorServiceWithControlProxy(
	t *testing.T,
	moduleName string,
	serviceType string,
	labels map[string]string,
	annotations map[string]string,
) (*cloud_service.OrchestratorService, net.Listener, string) {
	srv, lis := NewTestOrchestratorService(
		t, moduleName, serviceType, labels, annotations)
	tempDir := setControlProxyConfig(t, lis.Addr())
	return srv, lis, tempDir
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

// setControlProxyConfig creates a temporal control_proxy.yml and returns its
// location. Remember to delete the temporary file once the test it is done
// os.RemoveAll(dir)
func setControlProxyConfig(t *testing.T, addrs net.Addr) string {
	if addrs == nil {
		t.Fatalf("listener address is nil. Can't create control_proxy.yml")
	}
	splitAddrs := strings.Split(addrs.String(), ":")
	if len(splitAddrs) == 0 {
		t.Fatalf("listener address is empty  %s. Can't create control_proxy.yml", splitAddrs)
	}

	dir, err := ioutil.TempDir("", "magma_cfg_test")
	if err != nil {
		t.Fatalf("can't create temp directory for control_proxy.yml test config: %s", err)
	}
	cfgFilePath := filepath.Join(dir, "control_proxy.yml")

	port := splitAddrs[len(splitAddrs)-1]
	testCPConfigYamlWithValues := fmt.Sprintf(testCPConfigYaml, port)

	err = ioutil.WriteFile(cfgFilePath, []byte(testCPConfigYamlWithValues), os.ModePerm)
	if err != nil {
		t.Fatalf("can't write control_proxy.yml test config: %s", err)
	}

	platform_cfg.SetConfigDirectories(dir, dir+"/foo", dir+"/bar")
	cfg := config.GetControlProxyConfigs()
	if port != fmt.Sprint(cfg.LocalPort) {
		t.Fatalf("control_proxy.yml doesnt have the right port (%s)\n%s", port, testCPConfigYamlWithValues)
	}
	return dir
}
