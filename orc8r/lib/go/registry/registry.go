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
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"magma/orc8r/lib/go/protos"
	registry_client "magma/orc8r/lib/go/registry/client"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	ServiceRegistryServiceName = "service_registry"
	ServiceRegistryModeEnvVar  = "SERVICE_REGISTRY_MODE"
	DockerRegistryMode         = "docker"
	K8sRegistryMode            = "k8s"
	YamlRegistryMode           = "yaml"

	HttpServerPort  = 8080
	GrpcServicePort = 9180

	annotationFieldSeparator = ","
)

type ServiceRegistry struct {
	sync.RWMutex
	ServiceConnections map[string]*grpc.ClientConn
	ServiceLocations   map[string]ServiceLocation

	cloudConnMu      sync.RWMutex
	cloudConnections map[string]cloudConnection

	serviceRegistryMode string

	additionalOpts []grpc.DialOption
}

type cloudConnection struct {
	*grpc.ClientConn
	expiration time.Time
}

var localKeepaliveParams = keepalive.ClientParameters{
	Time:                31 * time.Second,
	Timeout:             10 * time.Second,
	PermitWithoutStream: true,
}

// New creates and returns a new registry
func New() *ServiceRegistry {
	registryMode := os.Getenv(ServiceRegistryModeEnvVar)
	return &ServiceRegistry{
		ServiceConnections:  map[string]*grpc.ClientConn{},
		ServiceLocations:    map[string]ServiceLocation{},
		cloudConnections:    map[string]cloudConnection{},
		serviceRegistryMode: registryMode,
	}
}

func NewWithDialOpts(opts ...grpc.DialOption) *ServiceRegistry {
	r := New()
	r.additionalOpts = opts
	return r
}

func NewWithMode(mode string) *ServiceRegistry {
	return &ServiceRegistry{
		ServiceConnections:  map[string]*grpc.ClientConn{},
		ServiceLocations:    map[string]ServiceLocation{},
		cloudConnections:    map[string]cloudConnection{},
		serviceRegistryMode: mode,
	}
}

// AddService add a new service.
// If the service already exists, overwrites the service config.
func (r *ServiceRegistry) AddService(location ServiceLocation) {
	r.Lock()
	defer r.Unlock()
	location.Name = strings.ToLower(location.Name)

	r.addUnsafe(location)
}

// AddServices adds new services to the registry.
// If any services already exist, their locations will be overwritten
func (r *ServiceRegistry) AddServices(locations ...ServiceLocation) {
	r.Lock()
	defer r.Unlock()

	for _, location := range locations {
		location.Name = strings.ToLower(location.Name)
		r.addUnsafe(location)
	}
}

// RemoveService removes a service from the registry.
// Has no effect if the service does not exist.
func (r *ServiceRegistry) RemoveService(service string) {
	r.Lock()
	defer r.Unlock()
	service = strings.ToLower(service)

	delete(r.ServiceLocations, service)
	delete(r.ServiceConnections, service)
}

// RemoveServicesWithLabel removes all services from the registry which have
// the passed label.
func (r *ServiceRegistry) RemoveServicesWithLabel(label string) {
	r.Lock()
	defer r.Unlock()

	for service, location := range r.ServiceLocations {
		if location.HasLabel(label) {
			delete(r.ServiceLocations, service)
			delete(r.ServiceConnections, service)
		}
	}
}

// ListAllServices lists the names of all registered services.
func (r *ServiceRegistry) ListAllServices() ([]string, error) {
	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode:
		client, err := r.getServiceRegistryServiceClient()
		if err != nil {
			return []string{}, err
		}
		return registry_client.ListAllServices(client)
	case YamlRegistryMode:
		fallthrough
	default:
		r.RLock()
		defer r.RUnlock()
		services := make([]string, 0, len(r.ServiceLocations))
		for service := range r.ServiceLocations {
			services = append(services, service)
		}
		return services, nil
	}
}

// FindServices returns the names of all registered services that have
// the passed label.
func (r *ServiceRegistry) FindServices(label string) ([]string, error) {
	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode:
		client, err := r.getServiceRegistryServiceClient()
		if err != nil {
			return []string{}, err
		}
		return registry_client.FindServices(client, label)
	case YamlRegistryMode:
		fallthrough
	default:
		r.RLock()
		defer r.RUnlock()
		var ret []string
		for service, location := range r.ServiceLocations {
			if location.HasLabel(label) {
				ret = append(ret, service)
			}
		}
		return ret, nil
	}
}

// GetServiceAddress returns the RPC address of the service.
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetServiceAddress(service string) (string, error) {
	service = strings.ToLower(service)

	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode:
		// Fetching the service registry service address is a special case
		// given that we cannot use the service registry service for discovering
		// it's own address
		if service == ServiceRegistryServiceName {
			addr := r.getServiceRegistryServiceAddress()
			if len(addr) == 0 {
				return "", fmt.Errorf("Service registry address is empty")
			}
			return addr, nil
		}
		client, err := r.getServiceRegistryServiceClient()
		if err != nil {
			return "", err
		}
		return registry_client.GetServiceAddress(client, service)
	case YamlRegistryMode:
		fallthrough
	default:
		r.RLock()
		defer r.RUnlock()
		location, ok := r.ServiceLocations[service]
		if !ok {
			return "", fmt.Errorf("service %s not registered", service)
		}
		if location.Port == 0 {
			return location.Host, nil
		}
		return fmt.Sprintf("%s:%d", location.Host, location.Port), nil
	}
}

// GetHttpServerAddress returns the HTTP address of the service from global registry
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetHttpServerAddress(service string) (string, error) {
	service = strings.ToLower(service)

	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode:
		client, err := r.getServiceRegistryServiceClient()
		if err != nil {
			return "", err
		}
		return registry_client.GetHttpServerAddress(client, service)
	case YamlRegistryMode:
		fallthrough
	default:
		r.RLock()
		defer r.RUnlock()
		location, ok := r.ServiceLocations[service]
		if !ok {
			return "", fmt.Errorf("service %s not registered", service)
		}
		if location.EchoPort == 0 {
			return "", fmt.Errorf("service %s is not available", service)
		}
		return fmt.Sprintf("%s:%d", location.Host, location.EchoPort), nil
	}
}

// GetServiceProxyAliases returns the proxy_aliases, if any, of the service.
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetServiceProxyAliases(service string) (map[string]int, error) {
	r.RLock()
	defer r.RUnlock()
	service = strings.ToLower(service)

	location, ok := r.ServiceLocations[service]
	if !ok {
		return nil, fmt.Errorf("failed to retrieve proxy alias: service '%s' not registered", service)
	}
	return location.ProxyAliases, nil
}

// GetServicePort returns the listening port for the RPC service.
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetServicePort(service string) (int, error) {
	r.RLock()
	defer r.RUnlock()
	service = strings.ToLower(service)

	// TODO: Update Docker/K8s service registry to return standard gRPC port
	// once controller image is split
	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode, YamlRegistryMode:
		fallthrough
	default:
		location, ok := r.ServiceLocations[service]
		if !ok {
			return 0, fmt.Errorf("service %s not registered", service)
		}
		if location.Port == 0 {
			return 0, fmt.Errorf("service %s not available", service)
		}
		return location.Port, nil
	}
}

// GetEchoServerPort returns the listening port for the service's echo server.
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetEchoServerPort(service string) (int, error) {
	r.RLock()
	defer r.RUnlock()
	service = strings.ToLower(service)

	// TODO: Update Docker/K8s service registry to return standard HTTP port
	// once controller image is split
	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode, YamlRegistryMode:
		fallthrough
	default:
		location, ok := r.ServiceLocations[service]
		if !ok {
			return 0, fmt.Errorf("failed to get echo server port: service '%s' not registered", service)
		}
		if location.EchoPort == 0 {
			return 0, fmt.Errorf("echo server port for service %s is not available", service)
		}
		return location.EchoPort, nil
	}
}

// GetAnnotation returns the annotation value for the passed annotation name.
func (r *ServiceRegistry) GetAnnotation(service, annotationName string) (string, error) {
	service = strings.ToLower(service)

	switch r.serviceRegistryMode {
	case DockerRegistryMode, K8sRegistryMode:
		client, err := r.getServiceRegistryServiceClient()
		if err != nil {
			return "", err
		}
		return registry_client.GetAnnotation(client, service, annotationName)
	case YamlRegistryMode:
		fallthrough
	default:
		r.RLock()
		defer r.RUnlock()
		location, ok := r.ServiceLocations[strings.ToLower(service)]
		if !ok {
			return "", fmt.Errorf("service %s not registered", service)
		}
		annotationValue, ok := location.Annotations[annotationName]
		if !ok {
			return "", fmt.Errorf("service %s doesn't have annotation values for %s", service, annotationName)
		}
		return annotationValue, nil
	}
}

// GetAnnotationList returns the comma-split fields of the value for the passed
// annotation name.
// First splits by field separator, then strips all whitespace
// (including newlines). Empty fields are removed.
func (r *ServiceRegistry) GetAnnotationList(service, annotationName string) ([]string, error) {
	annotationValue, err := r.GetAnnotation(service, annotationName)
	if err != nil {
		return nil, err
	}

	var values []string
	for _, s := range strings.Split(annotationValue, annotationFieldSeparator) {
		s = strings.Join(strings.Fields(s), "")
		if s != "" {
			values = append(values, s)
		}
	}

	return values, nil
}

// GetConnection provides a gRPC connection to a service in the registry.
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetConnection(service string) (*grpc.ClientConn, error) {
	return r.GetConnectionWithTimeout(service, GrpcMaxTimeoutSec*time.Second)
}

// GetConnectionWithTimeout is same as GetConnection, but caller can provide
// their own timeout.
// The service needs to be added to the registry before this.
func (r *ServiceRegistry) GetConnectionWithTimeout(service string, timeout time.Duration) (*grpc.ClientConn, error) {
	service = strings.ToLower(service)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return r.GetConnectionImpl(ctx, service, r.getGRPCDialOptions()...)
}

// GetConnectionWithOptions is same as GetConnection, but allows caller to
// provide their own gRPC dial options.
func (r *ServiceRegistry) GetConnectionWithOptions(service string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	service = strings.ToLower(service)
	ctx, cancel := context.WithTimeout(context.Background(), GrpcMaxTimeoutSec*time.Second)
	defer cancel()
	return r.GetConnectionImpl(ctx, service, options...)
}

func (r *ServiceRegistry) GetConnectionImpl(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	service = strings.ToLower(service)

	// First try to get an existing connection with reader lock
	r.RLock()
	conn, ok := r.ServiceConnections[service]
	r.RUnlock()
	if ok && conn != nil {
		return conn, nil
	}

	// Attempt to connect outside of the lock
	// Each attempt to get client connection has a long timeout. Connecting
	// without the lock prevents callers from timing out waiting for the
	// lock to a bad connection.
	addr, err := r.GetServiceAddress(service)
	if err != nil {
		return nil, err
	}
	newConn, err := GetClientConnection(ctx, addr, opts...)
	if err != nil || newConn == nil {
		return newConn, fmt.Errorf("service %v connection error: %s", service, err)
	}

	r.Lock()
	defer r.Unlock()

	// Re-check after taking the lock
	conn, ok = r.ServiceConnections[service]
	if ok && conn != nil {
		// Another routine already added the connection for the service, clean up ours & return existing
		err := newConn.Close()
		if err != nil {
			glog.Errorf("Error closing unneeded gRPC connection: %v", err)
		}
		return conn, nil
	}

	r.ServiceConnections[service] = newConn
	return newConn, nil
}

func (r *ServiceRegistry) getServiceRegistryServiceClient() (protos.ServiceRegistryClient, error) {
	conn, err := r.GetConnection(ServiceRegistryServiceName)
	if err != nil {
		return nil, err
	}
	return protos.NewServiceRegistryClient(conn), nil
}

func (r *ServiceRegistry) getGRPCDialOptions() []grpc.DialOption {
	opts := []grpc.DialOption{grpc.WithBackoffMaxDelay(GrpcMaxDelaySec * time.Second), grpc.WithBlock()}
	if *grpcKeepAlive {
		opts = append(opts, grpc.WithKeepaliveParams(localKeepaliveParams))
	}
	var timeoutInterceptor = TimeoutInterceptor
	if r.serviceRegistryMode == K8sRegistryMode || r.serviceRegistryMode == DockerRegistryMode {
		timeoutInterceptor = CloudClientTimeoutInterceptor
	}
	opts = append(opts, grpc.WithUnaryInterceptor(timeoutInterceptor))
	opts = append(opts, r.additionalOpts...)
	return opts
}

func (r *ServiceRegistry) getServiceRegistryServiceAddress() string {
	// Use hardcoded address for service_registry service as we can't
	// dynamically discover the service registry service itself
	return fmt.Sprintf("orc8r-service-registry:%d", GrpcServicePort)
}

func (r *ServiceRegistry) addUnsafe(location ServiceLocation) {
	if r.ServiceLocations == nil {
		r.ServiceLocations = map[string]ServiceLocation{}
	}
	r.ServiceLocations[location.Name] = location
	delete(r.ServiceConnections, location.Name)
}

// ServiceLocation is an entry for the service registry which identifies a
// service by name and the host:port that it is running on.
type ServiceLocation struct {
	// Name of the service.
	Name string
	// Host name of the service.
	Host string
	// Port is the service's gRPC endpoint.
	Port int
	// EchoPort is the service's HTTP endpoint for providing obsidian handlers.
	EchoPort int
	// ProxyAliases provides the list of host:port aliases for the service.
	ProxyAliases map[string]int

	// Labels provide a way to identify the service.
	// Use cases include listing service mesh servicers the service implements.
	Labels map[string]string
	// Annotations provides a string-to-string map of per-service metadata.
	Annotations map[string]string
}

func (s ServiceLocation) HasLabel(label string) bool {
	_, ok := s.Labels[label]
	return ok
}

// String implements ServiceLocation stringer interface
// Returns string in the form: <service name> @ host:port (also known as: host:port, ...)
func (s ServiceLocation) String() string {
	alsoKnown := ""
	if len(s.ProxyAliases) > 0 {
		aliases := ""
		for host, port := range s.ProxyAliases {
			aliases += fmt.Sprintf(" %s:%d,", host, port)
		}
		alsoKnown = " (also known as:" + aliases[:len(aliases)-1] + ")"
	}
	return fmt.Sprintf("%s @ %s:%d%s", s.Name, s.Host, s.Port, alsoKnown)
}

// GetClientConnection provides a gRPC connection to a service on the address addr.
func GetClientConnection(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("address: %s gRPC Dial error: %s", addr, err)
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return conn, nil
}
