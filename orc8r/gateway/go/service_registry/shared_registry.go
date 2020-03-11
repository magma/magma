/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// package cloud_registry provides CloudRegistry interface for Go based gateways
package service_registry

import (
	"log"
	"sync/atomic"

	platform_registry "magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/serviceregistry"
)

// default service registry shared by all GW process services
var initializedRegistry atomic.Value

// Get returns default service registry which can be shared by all GW process services
func Get() GatewayRegistry {
	reg, ok := initializedRegistry.Load().(GatewayRegistry)
	if (!ok) || reg == nil {
		reg = NewDefaultRegistry()
		// Overwrite/Add from /etc/magma/service_registry.yml if it exists
		// moduleName is "" since all feg configs lie in /etc/magma without a module name
		locations, err := serviceregistry.LoadServiceRegistryConfig("")
		if err != nil {
			log.Printf("Error loading Gateway service_registry.yml: %v", err)
			// return registry, but don't store/cache it
			return reg
		}
		if len(locations) > 0 {
			reg.AddServices(locations...)
		}
		initializedRegistry.Store(reg)
	}
	return reg
}

// NewDefaultRegistry returns a new service registry populated by default services
func NewDefaultRegistry() GatewayRegistry {
	reg := platform_registry.Get()
	if _, err := reg.GetServiceAddress(platform_registry.ControlProxyServiceName); err != nil {
		addLocal(reg, platform_registry.ControlProxyServiceName, 5053)
	}
	return reg
}

// Returns the RPC address of the service.
// The service needs to be added to the registry before this.
func GetServiceAddress(service string) (string, error) {
	return Get().GetServiceAddress(service)
}

func addLocal(reg GatewayRegistry, serviceType string, port int) {
	reg.AddService(platform_registry.ServiceLocation{Name: serviceType, Host: "localhost", Port: port})
}
