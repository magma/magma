/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin_test

import (
	"errors"
	"magma/orc8r/cloud/go/services/state/indexer"
	"testing"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/plugin/mocks"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type errorLoader struct{}

func (errorLoader) LoadPlugins() ([]plugin.OrchestratorPlugin, error) {
	return nil, errors.New("foobar")
}

type mockLoader struct {
	ret plugin.OrchestratorPlugin
}

func (m mockLoader) LoadPlugins() ([]plugin.OrchestratorPlugin, error) {
	return []plugin.OrchestratorPlugin{m.ret}, nil
}

func TestLoadAllPlugins(t *testing.T) {
	// Happy path - just make sure all functions on the plugin are called
	mockPlugin := &mocks.OrchestratorPlugin{}
	mockPlugin.On("GetMconfigBuilders").Return([]configurator.MconfigBuilder{})
	mockPlugin.On("GetMetricsProfiles", mock.Anything).Return([]metricsd.MetricsProfile{}).Once()
	mockPlugin.On("GetObsidianHandlers", mock.Anything).Return([]obsidian.Handler{})
	mockPlugin.On("GetSerdes").Return([]serde.Serde{})
	mockPlugin.On("GetServices").Return([]registry.ServiceLocation{})
	mockPlugin.On("GetStateIndexers").Return([]indexer.Indexer{})
	mockPlugin.On("GetStreamerProviders").Return([]providers.StreamProvider{})
	err := plugin.LoadAllPlugins(mockLoader{ret: mockPlugin})
	assert.NoError(t, err)
	mockPlugin.AssertNumberOfCalls(t, "GetMconfigBuilders", 1)
	mockPlugin.AssertNumberOfCalls(t, "GetMetricsProfiles", 1)
	mockPlugin.AssertNumberOfCalls(t, "GetObsidianHandlers", 1)
	mockPlugin.AssertNumberOfCalls(t, "GetSerdes", 1)
	mockPlugin.AssertNumberOfCalls(t, "GetServices", 1)
	mockPlugin.AssertNumberOfCalls(t, "GetStateIndexers", 1)
	mockPlugin.AssertNumberOfCalls(t, "GetStreamerProviders", 1)
	mockPlugin.AssertExpectations(t)

	// Error in the middle of registration - duplicate metrics profile
	mockPlugin.On("GetMetricsProfiles", mock.Anything).Times(1).Return(
		[]metricsd.MetricsProfile{
			{Name: "foo"},
			{Name: "foo"},
		},
	)
	err = plugin.LoadAllPlugins(mockLoader{ret: mockPlugin})
	assert.EqualError(t, err, "A metrics profile with the name foo already exists")
	mockPlugin.AssertExpectations(t)

	// Error from loader
	err = plugin.LoadAllPlugins(errorLoader{})
	assert.EqualError(t, err, "foobar")
}
