/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package plugin exposes the OrchestratorPlugin implementation for the module.
// This is so unit tests can register the plugin without building and loading
// it from disk.
package plugin

import (
	"magma/feg/cloud/go/feg"
	fegconfig "magma/feg/cloud/go/services/controller/config"
	fegh "magma/feg/cloud/go/services/controller/obsidian/handlers"
	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/service/serviceregistry"
	configregistry "magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/metricsd"
	stateregistry "magma/orc8r/cloud/go/services/state/registry"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/providers"
)

// FegOrchestratorPlugin is an implementation of OrchestratorPlugin for the
// feg module
type FegOrchestratorPlugin struct{}

func (*FegOrchestratorPlugin) GetName() string {
	return feg.ModuleName
}

func (*FegOrchestratorPlugin) GetServices() []registry.ServiceLocation {
	serviceLocations, err := serviceregistry.LoadServiceRegistryConfig(feg.ModuleName)
	if err != nil {
		return []registry.ServiceLocation{}
	}
	return serviceLocations

	// TODO: do we need FEG gateway services in the plugin or is that
	// the wild west?
	//{Name: session_proxy.ServiceName, Host: "localhost", Port: 9097},
	//{Name: s6a_proxy.ServiceName, Host: "localhost", Port: 9098},
	//{Name: csfb.ServiceName, Host: "localhost", Port: 9101},
	// feg hello went on 9093
}

func (*FegOrchestratorPlugin) GetConfigManagers() []configregistry.ConfigManager {
	return []configregistry.ConfigManager{
		&fegconfig.FegNetworkConfigManager{},
		&fegconfig.FegGatewayConfigManager{},
	}
}

func (*FegOrchestratorPlugin) GetStateSerdes() []stateregistry.StateSerde {
	return []stateregistry.StateSerde{}
}

func (*FegOrchestratorPlugin) GetMconfigBuilders() []factory.MconfigBuilder {
	return []factory.MconfigBuilder{
		&fegconfig.Builder{},
	}
}

func (*FegOrchestratorPlugin) GetMetricsProfiles() []metricsd.MetricsProfile {
	return []metricsd.MetricsProfile{}
}

func (*FegOrchestratorPlugin) GetObsidianHandlers() []handlers.Handler {
	return plugin.FlattenHandlerLists(
		fegh.GetObsidianHandlers(),
	)
}

func (*FegOrchestratorPlugin) GetStreamerProviders() []providers.StreamProvider {
	return []providers.StreamProvider{}
}
