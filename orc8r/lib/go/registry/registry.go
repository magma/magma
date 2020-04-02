/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package registry for Magma microservices
package registry

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ServiceLocation is an entry for the service registry which identifies a
// service by name and the host:port that it is running on.
type ServiceLocation struct {
	Name         string
	Host         string
	Port         int
	ProxyAliases map[string]int
}

type ServiceRegistry struct {
	sync.RWMutex
	ServiceConnections map[string]*grpc.ClientConn
	ServiceLocations   map[string]ServiceLocation
}

// New creates and returns a new registry
func New() *ServiceRegistry {
	return &ServiceRegistry{
		ServiceConnections: map[string]*grpc.ClientConn{},
		ServiceLocations:   map[string]ServiceLocation{}}
}

// String implements ServiceLocation stringer interface
// Returns string in the form: <Service name> @ host:port (also known as: host:port, ...)
func (sl ServiceLocation) String() string {
	alsoKnown := ""
	if len(sl.ProxyAliases) > 0 {
		aliases := ""
		for host, port := range sl.ProxyAliases {
			aliases += fmt.Sprintf(" %s:%d,", host, port)
		}
		alsoKnown = " (also known as:" + aliases[:len(aliases)-1] + ")"
	}
	return fmt.Sprintf("%s @ %s:%d%s", sl.Name, sl.Host, sl.Port, alsoKnown)
}

// Initialize initializes the registry maps
func (registry *ServiceRegistry) Initialize() *ServiceRegistry {
	if registry == nil {
		return New()
	}
	registry.ServiceConnections = map[string]*grpc.ClientConn{}
	registry.ServiceLocations = map[string]ServiceLocation{}
	return registry
}

// AddServices adds new services to the registry.
// If any services already exist, their locations will be overwritten
func (registry *ServiceRegistry) AddServices(locations ...ServiceLocation) {
	registry.Lock()
	defer registry.Unlock()
	for _, location := range locations {
		registry.addUnsafe(location)
	}
}

// AddService add a new service.
// If the service already exists, overwrites the service config.
func (registry *ServiceRegistry) AddService(location ServiceLocation) {
	registry.Lock()
	defer registry.Unlock()
	registry.addUnsafe(location)
}

// GetServiceAddress returns the RPC address of the service.
// The service needs to be added to the registry before this.
func (registry *ServiceRegistry) GetServiceAddress(service string) (string, error) {
	registry.RLock()
	defer registry.RUnlock()

	location, ok := registry.ServiceLocations[service]
	if !ok {
		return "", fmt.Errorf("Service '%s' not registered", service)
	}
	if location.Port == 0 {
		return "", fmt.Errorf("Service %s is not available", service)
	}
	return fmt.Sprintf("%s:%d", location.Host, location.Port), nil
}

// ListAllServices lists all services' name.
func (registry *ServiceRegistry) ListAllServices() []string {
	registry.RLock()
	defer registry.RUnlock()
	services := make([]string, 0, len(registry.ServiceLocations))
	for service := range registry.ServiceLocations {
		services = append(services, service)
	}
	return services
}

// GetServiceProxyAliases returns the proxy_aliases, if any, of the service.
// The service needs to be added to the registry before this.
func (registry *ServiceRegistry) GetServiceProxyAliases(service string) (map[string]int, error) {
	registry.RLock()
	defer registry.RUnlock()
	location, ok := registry.ServiceLocations[service]
	if !ok {
		return nil, fmt.Errorf("failed to retrieve proxy alias: Service '%s' not registered", service)
	}
	return location.ProxyAliases, nil
}

// GetServicePort returns the listening port for the RPC service.
// The service needs to be added to the registry before this.
func (registry *ServiceRegistry) GetServicePort(service string) (int, error) {
	registry.RLock()
	defer registry.RUnlock()
	location, ok := registry.ServiceLocations[strings.ToUpper(string(service))]
	if !ok {
		return 0, fmt.Errorf("failed to get service port: Service '%s' not registered", service)
	}

	if location.Port == 0 {
		return 0, fmt.Errorf("Service %s is not available", service)
	}
	return location.Port, nil
}

func (registry *ServiceRegistry) addUnsafe(location ServiceLocation) {
	if registry.ServiceLocations == nil {
		registry.ServiceLocations = map[string]ServiceLocation{}
	}
	registry.ServiceLocations[location.Name] = location
	delete(registry.ServiceConnections, location.Name)
}

// GetConnection provides a gRPC connection to a service in the registry.
func (registry *ServiceRegistry) GetConnection(service string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GrpcMaxTimeoutSec*time.Second)
	defer cancel()
	return registry.GetConnectionImpl(
		ctx,
		service,
		grpc.WithBackoffMaxDelay(GrpcMaxDelaySec*time.Second),
		grpc.WithBlock(),
	)
}

func (registry *ServiceRegistry) GetConnectionImpl(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	addr, err := registry.GetServiceAddress(service)
	if err != nil {
		return nil, err
	}
	// First try to get an existing connection with reader lock
	registry.RLock()
	conn, ok := registry.ServiceConnections[service]
	registry.RUnlock()
	if ok && conn != nil {
		return conn, nil
	}

	// Attempt to connect outside of the lock
	newConn, err := GetClientConnection(ctx, addr, opts...)
	if err != nil || newConn == nil {
		return newConn, fmt.Errorf("Service %v connection error: %s", service, err)
	}

	registry.Lock()
	defer registry.Unlock()

	// Re-check after taking the lock
	conn, ok = registry.ServiceConnections[service]
	if ok && conn != nil {
		// another routine already added the connection for the service, clean up ours & return existing
		newConn.Close()
		return conn, nil
	}
	registry.ServiceConnections[service] = newConn
	return newConn, nil
}

// GetClientConnection provides a gRPC connection to a service on the address addr.
func GetClientConnection(ctx context.Context, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("Address: %s GRPC Dial error: %s", addr, err)
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return conn, nil
}
