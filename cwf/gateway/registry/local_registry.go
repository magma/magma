/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package registry

import (
	platform_registry "magma/orc8r/lib/go/registry"

	"google.golang.org/grpc"
)

const (
	ModuleName = "cwf"

	GatewayHealth = "HEALTH"
	UeSim         = "UESIM"
	Radiusd       = "RADIUSD"
)

// Add a new service.
// If the service already exists, overwrites the service config.
func AddService(serviceType, host string, port int) {
	platform_registry.AddService(platform_registry.ServiceLocation{Name: serviceType, Host: host, Port: port})
}

// Returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return platform_registry.GetServiceAddress(service)
}

// Provides a gRPC connection to a service in the registry.
func GetConnection(service string) (*grpc.ClientConn, error) {
	return platform_registry.GetConnection(service)
}

func addLocalService(serviceType string, port int) {
	AddService(serviceType, "localhost", port)
}

func init() {
	addLocalService(UeSim, 10101)
	addLocalService(GatewayHealth, 9107)
	addLocalService(Radiusd, 10102)
}
