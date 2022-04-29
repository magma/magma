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

package registry

import (
	"magma/orc8r/lib/go/protos"
	"time"

	platform_registry "magma/orc8r/lib/go/registry"

	"google.golang.org/grpc"
)

const (
	ModuleName = "cwf"

	Eventd        = "EVENTD"
	GatewayHealth = "HEALTH"
	UeSim         = "UESIM"
	Radiusd       = "RADIUSD"
)

// AddService Add a new service.
// If the service already exists, overwrites the service config.
func AddService(serviceType, host string, port int) {
	platform_registry.AddService(platform_registry.ServiceLocation{Name: serviceType, Host: host, Port: port})
}

// GetServiceAddress Returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return platform_registry.GetServiceAddress(service, protos.ServiceType_SOUTHBOUND)
}

// GetConnection Provides a gRPC connection to a service in the registry.
func GetConnection(service string) (*grpc.ClientConn, error) {
	return platform_registry.GetConnection(service, protos.ServiceType_SOUTHBOUND)
}

// GetConnectionWithTimeout Provides a gRPC connection with timeout to a service in the registry.
func GetConnectionWithTimeout(service string, timeout time.Duration) (*grpc.ClientConn, error) {
	return platform_registry.GetConnectionWithTimeout(service, timeout, protos.ServiceType_SOUTHBOUND)
}

func addLocalService(serviceType string, port int) {
	AddService(serviceType, "localhost", port)
}

func init() {
	addLocalService(Eventd, 50075)
	addLocalService(UeSim, 10101)
	addLocalService(GatewayHealth, 9107)
	addLocalService(Radiusd, 10102)
}
