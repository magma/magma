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
	"os"
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

const (
	GrpcMaxDelaySec   = 10
	GrpxMaxTimeoutSec = 60
)

type serviceRegistry struct {
	sync.RWMutex
	serviceConnections map[string]*grpc.ClientConn
	serviceLocations   map[string]ServiceLocation
}

var registry = &serviceRegistry{
	serviceConnections: map[string]*grpc.ClientConn{},
	serviceLocations:   map[string]ServiceLocation{},
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

// AddServices adds new services to the registry.
// If any services already exist, their locations will be overwritten
func AddServices(locations ...ServiceLocation) {
	registry.Lock()
	defer registry.Unlock()
	for _, location := range locations {
		addUnsafe(location)
	}
}

// AddService add a new service.
// If the service already exists, overwrites the service config.
func AddService(location ServiceLocation) {
	registry.Lock()
	defer registry.Unlock()
	addUnsafe(location)
}

func addUnsafe(location ServiceLocation) {
	registry.serviceLocations[location.Name] = location
	delete(registry.serviceConnections, location.Name)
}

// GetServiceAddress returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	registry.RLock()
	defer registry.RUnlock()

	location, ok := registry.serviceLocations[service]
	if !ok {
		return "", fmt.Errorf("Service %s not registered", service)
	}
	if location.Port == 0 {
		return "", fmt.Errorf("Service %s is not available", service)
	}
	return fmt.Sprintf("%s:%d", location.Host, location.Port), nil
}

// GetServiceProxyAliases returns the proxy_aliases, if any, of the service.
// The service needs to be added to the registry before this.
func GetServiceProxyAliases(service string) (map[string]int, error) {
	registry.RLock()
	defer registry.RUnlock()
	location, ok := registry.serviceLocations[service]
	if !ok {
		return nil, fmt.Errorf("Service %s not registered", service)
	}
	return location.ProxyAliases, nil
}

// GetServicePort returns the listening port for the RPC service.
// The service needs to be added to the registry before this.
func GetServicePort(service string) (int, error) {
	registry.RLock()
	defer registry.RUnlock()
	location, ok := registry.serviceLocations[strings.ToUpper(string(service))]
	if !ok {
		return 0, fmt.Errorf("Service %s not registered", service)
	}

	if location.Port == 0 {
		return 0, fmt.Errorf("Service %s is not available", service)
	}
	return location.Port, nil
}

// GetConnection provides a gRPC connection to a service in the registry.
func GetConnection(service string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GrpxMaxTimeoutSec*time.Second)
	defer cancel()
	return GetConnectionImpl(
		ctx,
		service,
		grpc.WithBackoffMaxDelay(GrpcMaxDelaySec*time.Second),
		grpc.WithBlock(),
	)
}

func GetConnectionImpl(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	addr, err := GetServiceAddress(service)
	if err != nil {
		return nil, err
	}
	// First try to get an existing connection with reader lock
	registry.RLock()
	conn, ok := registry.serviceConnections[service]
	registry.RUnlock()
	if ok && conn != nil {
		return conn, nil
	}

	registry.Lock()
	defer registry.Unlock()
	// Re-check after taking the lock
	conn, ok = registry.serviceConnections[service]
	if ok && conn != nil {
		return conn, nil
	}
	conn, err = GetClientConnection(ctx, addr, opts...)
	if err != nil {
		err = fmt.Errorf("Service %v connection error: %s", service, err)
	} else {
		registry.serviceConnections[service] = conn
	}
	return conn, err
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

// ListAllServices lists all services' name.
func ListAllServices() []string {
	registry.RLock()
	defer registry.RUnlock()
	services := make([]string, 0, len(registry.serviceLocations))
	for service := range registry.serviceLocations {
		services = append(services, service)
	}
	return services
}

// ListControllerServices list all services that should run on a controller instances
// This is a comma separated list in an env var named CONTROLLER_SERVICES. This
// will be used for metricsd on controller to determine
// what services to pull metrics from.
func ListControllerServices() []string {
	ret := make([]string, 0)
	controllerServices, ok := os.LookupEnv("CONTROLLER_SERVICES")
	if !ok {
		return ret
	}
	return strings.Split(controllerServices, ",")
}
