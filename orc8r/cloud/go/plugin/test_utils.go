/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import (
	"sync"
	"testing"

	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/lib/go/service/config"
)

type testPluginRegistry struct {
	sync.Mutex
	plugins map[string]OrchestratorPlugin
}

var testPlugins = &testPluginRegistry{plugins: map[string]OrchestratorPlugin{}}

var testMetricsConfigMap = &config.ConfigMap{
	RawMap: map[interface{}]interface{}{
		metricsd.Profile:                      "",
		metricsd.PrometheusPushAddresses:      []string{""},
		metricsd.PrometheusQueryAddress:       "",
		metricsd.PrometheusConfigServiceURL:   "",
		metricsd.AlertmanagerConfigServiceURL: "",
		metricsd.AlertmanagerApiURL:           "",
	},
}

// RegisterPluginForTests registers all components of a given plugin with the
// corresponding component registries exposed by the orchestrator. This should
// only be used in test code to avoid the cost of building and loading plugins
// from disk for unit tests, thus the required but unused *testing.T
// parameter. This function will not register a plugin which has already been
// registered as identified by its GetName().
func RegisterPluginForTests(_ *testing.T, plugin OrchestratorPlugin) error {
	testPlugins.Lock()
	defer testPlugins.Unlock()
	if _, ok := testPlugins.plugins[plugin.GetName()]; !ok {
		testPlugins.plugins[plugin.GetName()] = plugin
		return registerPlugin(plugin, testMetricsConfigMap)
	}
	// plugin has already been registered, no-op
	return nil
}
