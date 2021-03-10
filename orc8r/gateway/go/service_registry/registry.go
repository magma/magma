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
// package cloud_registry provides CloudRegistry interface for Go based gateways
package service_registry

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

// ServiceRegistry defines interface for gateway's registry of local services
type GatewayRegistry interface {
	// platform's ServiceRegistry methods
	AddServices(locations ...platform_registry.ServiceLocation)
	AddService(location platform_registry.ServiceLocation)
	GetServiceAddress(service string) (string, error)
	GetServicePort(service string) (int, error)
	GetServiceProxyAliases(service string) (map[string]int, error)
	ListAllServices() ([]string, error)

	// Gateway specific methods
	//
	// GetCloudConnection Creates and returns a new GRPC service connection to the service.
	// GetCloudConnection always creates a new connection & it's responsibility of the caller to close it.
	GetCloudConnection(service string) (*grpc.ClientConn, error)
	// GetCloudConnectionFromServiceConfig returns a connection to the cloud
	// using a specific control_proxy service config map. This map must contain the cloud_address
	// and local_port params
	// Input: serviceConfig - ConfigMap containing cloud_address and local_port and optional proxy_cloud_connections,
	// 		  cloud_port, rootca_cert, gateway_cert/key fields if direct cloud connection is needed
	//        service - name of cloud service to connect to
	// Output: *grpc.ClientConn with connection to cloud service error if it exists
	GetCloudConnectionFromServiceConfig(controlProxyConfig *config.ConfigMap, service string) (*grpc.ClientConn, error)
	// GetSharedCloudConnection returns a new GRPC service connection to the service in the cloud for a gateway
	// either directly or via control proxy
	// GetSharedCloudConnection will return an existing cached cloud connection if it's available and healthy,
	// if not - it'll try to create, cache and return a new cloud connection
	GetSharedCloudConnection(service string) (*grpc.ClientConn, error)
	// GetSharedCloudConnectionFromServiceConfig returns a connection to the cloud
	// using a specific control_proxy service config map. This map must contain the cloud_address
	// and local_port params
	// GetSharedCloudConnectionFromServiceConfig will return an existing cached cloud connection if it's available and
	// healthy, if not - it'll try to create, cache and return a new cloud connection
	GetSharedCloudConnectionFromServiceConfig(controlProxyConfig *config.ConfigMap, service string) (*grpc.ClientConn, error)
	// GetClientConnection provides a gRPC connection to a service on the address addr.
	CleanupSharedCloudConnection(service string) bool

	// GetConnection provides a gRPC connection to a service in the registry.
	GetConnection(service string) (*grpc.ClientConn, error)
	GetConnectionImpl(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

// ProxiedRegistry a GW service registry which supports direct and proxied cloud connections
type ProxiedRegistry platform_registry.ServiceRegistry

// New returns a new ProxiedRegistry initialized with empty service & connection maps
func New() GatewayRegistry {
	return platform_registry.New()
}
