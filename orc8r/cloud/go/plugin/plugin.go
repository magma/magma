/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"plugin"
	"reflect"
	"strings"

	"magma/orc8r/cloud/go/obsidian/handlers"
	"magma/orc8r/cloud/go/registry"
	config_registry "magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/metricsd"
	state_registry "magma/orc8r/cloud/go/services/state/registry"
	"magma/orc8r/cloud/go/services/streamer/mconfig/factory"
	"magma/orc8r/cloud/go/services/streamer/providers"

	"github.com/golang/glog"
)

const (
	modulePluginDir       = "/var/opt/magma/plugins/"
	moduleFactoryFunction = "GetOrchestratorPlugin"
)

// OrchestratorPlugin defines the functionality that a plugin on the magma
// cloud side is expected to implement and provide. This interface is the
// formal surface area for integrating into and extending the magma
// orchestrator.
type OrchestratorPlugin interface {
	// GetName returns a unique name for the plugin.
	GetName() string

	// GetServices returns a list of services that this plugin runs to register
	// with the orc8r service registry.
	GetServices() []registry.ServiceLocation

	// GetConfigManagers returns a list of ConfigManagers to register with the
	// config service. This is how a plugin adds support for custom
	// configurations to the config service.
	GetConfigManagers() []config_registry.ConfigManager

	// GetStateSerdes returns a list of StateSerdeManagers to register
	// with the state service. This is how a plugin adds support for custom
	// marshall/unmarshall methods to the state service.
	GetStateSerdes() []state_registry.StateSerde

	// GetMconfigBuilders returns a list of MconfigBuilders to register with
	// the config streamer application. These builders are responsible for
	// constructing gateway mconfigs from cloud-stored configs.
	GetMconfigBuilders() []factory.MconfigBuilder

	// GetMetricsProfile returns the metricsd profiles that this module
	// supplies. This will make specific configurations available for metricsd
	// to load on startup. See MetricsProfile for additional documentation.
	GetMetricsProfiles() []metricsd.MetricsProfile

	// GetObsidianHandlers returns all the custom obsidian handlers for the
	// plugin to add functionality to the REST API.
	GetObsidianHandlers() []handlers.Handler

	// GetStreamerProviders returns streamer streams to expose to gateways.
	// These stream providers are the primary mechanism by which gateways
	// receive data from the orchestrator (e.g. configuration).
	GetStreamerProviders() []providers.StreamProvider
}

// LoadAllPluginsFatalOnError loads and registers all orchestrator plugins
// and calls os.Exit() on error. See LoadAllPlugins for additional
// documentation.
func LoadAllPluginsFatalOnError(loader OrchestratorPluginLoader) {
	if err := LoadAllPlugins(loader); err != nil {
		glog.Fatal(err)
	}
}

// LoadAlPlugins loads and registers all orchestrator plugins, returning the
// first error encountered during the process. Standard use-cases should pass
// DefaultOrchestratorPluginLoader.
//
// This function will NOT roll back registered plugins if it fails in the
// middle of execution. For this reason, you will likely prefer to use
// LoadAllPluginsFatalOnError which wraps this function with a glog.Fatal.
func LoadAllPlugins(loader OrchestratorPluginLoader) error {
	plugins, err := loader.LoadPlugins()
	if err != nil {
		return err
	}

	for _, p := range plugins {
		if err := registerPlugin(p); err != nil {
			return err
		}
	}
	return nil
}

// OrchestratorPluginLoader wraps the loading of OrchestratorPlugin impls.
// Standard use case is to use the provided DefaultOrchestratorPluginLoader
// in this package - only create a new impl if you need to customize the
// loading process in some way (e.g. loading from a different directory).
type OrchestratorPluginLoader interface {
	LoadPlugins() ([]OrchestratorPlugin, error)
}

// DefaultOrchestratorPluginLoader looks for all .so files in
// /var/opt/magma/plugins and tries to load each .so as an OrchestratorPlugin.
type DefaultOrchestratorPluginLoader struct{}

func (DefaultOrchestratorPluginLoader) LoadPlugins() ([]OrchestratorPlugin, error) {
	var ret []OrchestratorPlugin

	_, err := os.Stat(modulePluginDir)
	if err != nil {
		// No plugins to load
		if os.IsNotExist(err) {
			return ret, nil
		}
		return ret, fmt.Errorf("Failed to stat plugin directory: %s", err)
	}

	files, err := ioutil.ReadDir(modulePluginDir)
	if err != nil {
		return ret, fmt.Errorf("Failed to read plugin directory contents: %s", err)
	}

	for _, file := range files {
		isPlugin := strings.HasSuffix(file.Name(), ".so") && !file.IsDir()
		if !isPlugin {
			glog.Infof("Not loading file %s in plugin directory because it does not appear to be a valid plugin", file.Name())
			continue
		}

		p, err := plugin.Open(modulePluginDir + file.Name())
		if err != nil {
			return []OrchestratorPlugin{}, fmt.Errorf("Could not open plugin %s: %s", file.Name(), err)
		}
		pluginFactory, err := p.Lookup(moduleFactoryFunction)
		if err != nil {
			return []OrchestratorPlugin{}, fmt.Errorf(
				"Failed lookup for plugin factory function %s for plugin %s: %s",
				moduleFactoryFunction, file.Name(), err,
			)
		}
		castedPluginFactory, ok := pluginFactory.(func() OrchestratorPlugin)
		if !ok {
			return []OrchestratorPlugin{}, fmt.Errorf(
				"Failed to cast plugin factory function from plugin %s. Expected func() OrchestratorPlugin, got %s",
				file.Name(), reflect.TypeOf(pluginFactory),
			)
		}
		ret = append(ret, castedPluginFactory())
	}
	return ret, nil
}

func registerPlugin(orc8rPlugin OrchestratorPlugin) error {
	registry.AddServices(orc8rPlugin.GetServices()...)
	if err := config_registry.RegisterConfigManagers(orc8rPlugin.GetConfigManagers()...); err != nil {
		return err
	}
	if err := state_registry.RegisterStateSerdes(orc8rPlugin.GetStateSerdes()...); err != nil {
		return err
	}
	factory.RegisterMconfigBuilders(orc8rPlugin.GetMconfigBuilders()...)
	if err := metricsd.RegisterMetricsProfiles(orc8rPlugin.GetMetricsProfiles()...); err != nil {
		return err
	}
	if err := handlers.RegisterAll(orc8rPlugin.GetObsidianHandlers()); err != nil {
		return err
	}
	if err := providers.RegisterStreamProviders(orc8rPlugin.GetStreamerProviders()...); err != nil {
		return err
	}

	return nil
}
