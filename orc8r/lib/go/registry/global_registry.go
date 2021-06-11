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

// Package registry for Magma microservices
package registry

import (
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	GrpcMaxDelaySec        = 10
	GrpcMaxLocalTimeoutSec = 30
	GrpcMaxTimeoutSec      = 60
)

// globalRegistry is the global service registry instance
var globalRegistry = New()

// Get returns a reference to the instance of global platform registry
func Get() *ServiceRegistry {
	return globalRegistry
}

func SetDialOpts(opts ...grpc.DialOption) {
	globalRegistry.additionalOpts = opts
}

// PopulateServices populates the service registry based on the per-module
// config files at /etc/magma/configs/MODULE_NAME/service_registry.yml.
func PopulateServices() error {
	serviceConfigs, err := LoadServiceRegistryConfigs()
	if err != nil {
		return err
	}
	AddServices(serviceConfigs...)
	return nil
}

// MustPopulateServices is same as PopulateServices but fails on errors.
func MustPopulateServices() {
	if err := PopulateServices(); err != nil {
		glog.Fatalf("Error populating services: %v", err)
	}
}

// AddService add a new service to global registry.
// If the service already exists, overwrites the service config.
func AddService(location ServiceLocation) {
	globalRegistry.AddService(location)
}

// AddServices adds new services to the global registry.
// If any services already exist, their locations will be overwritten
func AddServices(locations ...ServiceLocation) {
	globalRegistry.AddServices(locations...)
}

// RemoveService removes a service from the registry.
// Has no effect if the service does not exist.
func RemoveService(service string) {
	globalRegistry.RemoveService(service)
}

// RemoveServicesWithLabel removes all services from the registry which have
// the passed label.
func RemoveServicesWithLabel(label string) {
	globalRegistry.RemoveServicesWithLabel(label)
}

// ListAllServices lists all services' names from global registry
func ListAllServices() ([]string, error) {
	return globalRegistry.ListAllServices()
}

// FindServices returns the names of all registered services that have
// the passed label.
func FindServices(label string) ([]string, error) {
	return globalRegistry.FindServices(label)
}

// GetServiceAddress returns the RPC address of the service from global registry
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return globalRegistry.GetServiceAddress(service)
}

// GetHttpServerAddress returns the HTTP address of the service from global registry
// The service needs to be added to the registry before this.
func GetHttpServerAddress(service string) (string, error) {
	return globalRegistry.GetHttpServerAddress(service)
}

// GetServiceProxyAliases returns the proxy_aliases, if any, of the service from global registry
// The service needs to be added to the registry before this.
func GetServiceProxyAliases(service string) (map[string]int, error) {
	return globalRegistry.GetServiceProxyAliases(service)
}

// GetServicePort returns the listening port for the RPC service.
// The service needs to be added to the registry before this.
func GetServicePort(service string) (int, error) {
	return globalRegistry.GetServicePort(service)
}

// GetEchoServerPort returns the listening port for the service's echo server.
// The service needs to be added to the registry before this.
func GetEchoServerPort(service string) (int, error) {
	return globalRegistry.GetEchoServerPort(service)
}

// GetAnnotation returns the annotation value for the passed annotation name.
// The service needs to be added to the registry before this.
func GetAnnotation(service, annotationName string) (string, error) {
	return globalRegistry.GetAnnotation(service, annotationName)
}

// GetAnnotationList returns the comma-split fields of the value for the passed
// annotation name.
// The service needs to be added to the registry before this.
func GetAnnotationList(service, annotationName string) ([]string, error) {
	return globalRegistry.GetAnnotationList(service, annotationName)
}

// GetConnection provides a gRPC connection to a service in the registry.
func GetConnection(service string) (*grpc.ClientConn, error) {
	return globalRegistry.GetConnection(service)
}

// GetConnectionWithTimeout is same as GetConnection, but caller can provide
// their own timeout.
func GetConnectionWithTimeout(service string, timeout time.Duration) (*grpc.ClientConn, error) {
	return globalRegistry.GetConnectionWithTimeout(service, timeout)
}

func GetConnectionImpl(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return globalRegistry.GetConnectionImpl(ctx, service, opts...)
}
