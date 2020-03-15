/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package registry for Magma microservices
package registry

import (
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	GrpcMaxDelaySec        = 10
	GrpcMaxLocalTimeoutSec = 30
	GrpcMaxTimeoutSec      = 60
)

var globalRegistry = New() // global registry instance

// Get returns a reference to the instance of global platform registry
func Get() *ServiceRegistry {
	return globalRegistry
}

// AddServices adds new services to the global registry.
// If any services already exist, their locations will be overwritten
func AddServices(locations ...ServiceLocation) {
	globalRegistry.AddServices(locations...)
}

// AddService add a new service to global registry.
// If the service already exists, overwrites the service config.
func AddService(location ServiceLocation) {
	globalRegistry.AddService(location)
}

// GetServiceAddress returns the RPC address of the service from global registry
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return globalRegistry.GetServiceAddress(service)
}

// ListAllServices lists all services' names from global registry
func ListAllServices() []string {
	return globalRegistry.ListAllServices()
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

// GetConnection provides a gRPC connection to a service in the registry.
func GetConnection(service string) (*grpc.ClientConn, error) {
	return globalRegistry.GetConnection(service)
}

func GetConnectionImpl(ctx context.Context, service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return globalRegistry.GetConnectionImpl(ctx, service, opts...)
}

// ListControllerServices list all services that should run on a controller instances
// This is a comma separated list in an env var named CONTROLLER_SERVICES. This
// will be used for metricsd on controller to determine
// what services to pull metrics from.
func ListControllerServices() []string {
	controllerServices, ok := os.LookupEnv("CONTROLLER_SERVICES")
	if !ok {
		return make([]string, 0)
	}
	return strings.Split(controllerServices, ",")
}
