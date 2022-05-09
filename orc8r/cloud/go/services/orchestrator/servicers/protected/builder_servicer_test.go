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

package servicers_test

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/orchestrator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	servicers "magma/orc8r/cloud/go/services/orchestrator/servicers/protected"
	orchestrator_test_init "magma/orc8r/cloud/go/services/orchestrator/test_init"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
	mconfig_protos "magma/orc8r/lib/go/protos/mconfig"
)

func TestBaseOrchestratorMconfigBuilder_Build(t *testing.T) {
	orchestrator_test_init.StartTestService(t)
	syncInterval := uint32(500)
	expectedDefaultJitteredSyncIntervalGW1 := uint32(71)
	expectedJitteredSyncIntervalGW1 := uint32(592)
	expectedJitteredSyncIntervalGW2 := uint32(568)
	version := models.TierVersion("1.0.0-0")

	t.Run("test shared config", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			orc8r.NetworkSentryConfig: &models.NetworkSentryConfig{
				SampleRate:   swag.Float32(0.75),
				UploadMmeLog: true,
				URLPython:    "https://www.example.com/v1/api",
				URLNative:    "https://www.example.com/v1/api",
			},
		}}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  "gw1",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				DynamicServices:         []string{},
				FeatureFlags:            map[string]bool{},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw},
		}

		expected := map[string]proto.Message{
			"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
			"magmad": &mconfig_protos.MagmaD{
				LogLevel:                protos.LogLevel_INFO,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				AutoupgradeEnabled:      true,
				AutoupgradePollInterval: 300,
				PackageVersion:          "0.0.0-0",
				Images:                  nil,
				DynamicServices:         nil,
				FeatureFlags:            nil,
			},
			"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
			"td-agent-bit": &mconfig_protos.FluentBit{
				ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
				ThrottleRate:     1000,
				ThrottleWindow:   5,
				ThrottleInterval: "1m",
			},
			"eventd": &mconfig_protos.EventD{
				LogLevel:       protos.LogLevel_INFO,
				EventVerbosity: -1,
			},
			"state": &mconfig_protos.State{
				SyncInterval: expectedDefaultJitteredSyncIntervalGW1,
				LogLevel:     protos.LogLevel_INFO,
			},
			"shared_mconfig": &mconfig_protos.SharedMconfig{
				SentryConfig: &mconfig_protos.SharedSentryConfig{
					SampleRate:   0.75,
					UploadMmeLog: true,
					DsnPython:    "https://www.example.com/v1/api",
					DsnNative:    "https://www.example.com/v1/api",
				},
			},
		}
		actual, err := buildBaseOrchestrator(&nw, &graph, "gw1")
		assert.NoError(t, err)
		test_utils.AssertMapsEqual(t, expected, actual)
	})

	t.Run("no tier", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": &models.StateConfig{
				SyncInterval: syncInterval,
			},
		}}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  "gw1",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				DynamicServices:         []string{},
				FeatureFlags:            map[string]bool{},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw},
		}

		expected := map[string]proto.Message{
			"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
			"magmad": &mconfig_protos.MagmaD{
				LogLevel:                protos.LogLevel_INFO,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				AutoupgradeEnabled:      true,
				AutoupgradePollInterval: 300,
				PackageVersion:          "0.0.0-0",
				Images:                  nil,
				DynamicServices:         nil,
				FeatureFlags:            nil,
			},
			"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
			"td-agent-bit": &mconfig_protos.FluentBit{
				ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
				ThrottleRate:     1000,
				ThrottleWindow:   5,
				ThrottleInterval: "1m",
			},
			"eventd": &mconfig_protos.EventD{
				LogLevel:       protos.LogLevel_INFO,
				EventVerbosity: -1,
			},
			"state": &mconfig_protos.State{
				SyncInterval: expectedJitteredSyncIntervalGW1,
				LogLevel:     protos.LogLevel_INFO,
			},
			"shared_mconfig": &mconfig_protos.SharedMconfig{
				SentryConfig: nil,
			},
		}

		actual, err := buildBaseOrchestrator(&nw, &graph, "gw1")
		assert.NoError(t, err)
		test_utils.AssertMapsEqual(t, expected, actual)
	})

	// Put a tier in the graph
	t.Run("tiers work correctly", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": &models.StateConfig{
				SyncInterval: syncInterval,
			},
		}}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  "gw1",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				DynamicServices:         []string{},
				FeatureFlags:            map[string]bool{},
			},
		}

		tier := configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType,
			Key:  "default",
			Config: &models.Tier{
				Name:    "default",
				Version: &version,
				Images: []*models.TierImage{
					{Name: swag.String("Image1"), Order: swag.Int64(42)},
					{Name: swag.String("Image2"), Order: swag.Int64(1)},
				},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw, tier},
			Edges: []configurator.GraphEdge{
				{From: tier.GetTK(), To: gw.GetTK()},
			},
		}

		expected := map[string]proto.Message{
			"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
			"magmad": &mconfig_protos.MagmaD{
				LogLevel:                protos.LogLevel_INFO,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				AutoupgradeEnabled:      true,
				AutoupgradePollInterval: 300,
				PackageVersion:          "1.0.0-0",
				Images: []*mconfig_protos.ImageSpec{
					{Name: "Image1", Order: 42},
					{Name: "Image2", Order: 1},
				},
				DynamicServices: nil,
				FeatureFlags:    nil,
			},
			"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
			"td-agent-bit": &mconfig_protos.FluentBit{
				ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
				ThrottleRate:     1000,
				ThrottleWindow:   5,
				ThrottleInterval: "1m",
			},
			"eventd": &mconfig_protos.EventD{
				LogLevel:       protos.LogLevel_INFO,
				EventVerbosity: -1,
			},
			"state": &mconfig_protos.State{
				SyncInterval: expectedJitteredSyncIntervalGW1,
				LogLevel:     protos.LogLevel_INFO,
			},
			"shared_mconfig": &mconfig_protos.SharedMconfig{
				SentryConfig: nil,
			},
		}

		actual, err := buildBaseOrchestrator(&nw, &graph, "gw1")
		assert.NoError(t, err)
		test_utils.AssertMapsEqual(t, expected, actual)
	})

	t.Run("set list of files for log aggregation", func(t *testing.T) {
		testThrottleInterval := "30h"
		testThrottleWindow := uint32(808)
		testThrottleRate := uint32(305)

		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": &models.StateConfig{
				SyncInterval: syncInterval,
			},
		}}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  "gw1",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				DynamicServices:         nil,
				FeatureFlags:            nil,
				Logging: &models.GatewayLoggingConfigs{
					Aggregation: &models.AggregationLoggingConfigs{
						TargetFilesByTag: map[string]string{
							"thing": "/var/log/thing.log",
							"blah":  "/some/directory/blah.log",
						},
						ThrottleRate:     &testThrottleRate,
						ThrottleWindow:   &testThrottleWindow,
						ThrottleInterval: &testThrottleInterval,
					},
					EventVerbosity: swag.Int32(0),
				},
			},
		}
		tier := configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType,
			Key:  "default",
			Config: &models.Tier{
				Name:    "default",
				Version: &version,
				Images: []*models.TierImage{
					{Name: swag.String("Image1"), Order: swag.Int64(42)},
					{Name: swag.String("Image2"), Order: swag.Int64(1)},
				},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw, tier},
			Edges: []configurator.GraphEdge{
				{From: tier.GetTK(), To: gw.GetTK()},
			},
		}

		expected := map[string]proto.Message{
			"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
			"magmad": &mconfig_protos.MagmaD{
				LogLevel:                protos.LogLevel_INFO,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				AutoupgradeEnabled:      true,
				AutoupgradePollInterval: 300,
				PackageVersion:          "1.0.0-0",
				Images: []*mconfig_protos.ImageSpec{
					{Name: "Image1", Order: 42},
					{Name: "Image2", Order: 1},
				},
				DynamicServices: nil,
				FeatureFlags:    nil,
			},
			"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
			"td-agent-bit": &mconfig_protos.FluentBit{
				ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
				ThrottleRate:     305,
				ThrottleWindow:   808,
				ThrottleInterval: "30h",
				FilesByTag: map[string]string{
					"thing": "/var/log/thing.log",
					"blah":  "/some/directory/blah.log",
				},
			},
			"eventd": &mconfig_protos.EventD{
				LogLevel:       protos.LogLevel_INFO,
				EventVerbosity: 0,
			},
			"state": &mconfig_protos.State{
				SyncInterval: expectedJitteredSyncIntervalGW1,
				LogLevel:     protos.LogLevel_INFO,
			},
			"shared_mconfig": &mconfig_protos.SharedMconfig{
				SentryConfig: nil,
			},
		}

		actual, err := buildBaseOrchestrator(&nw, &graph, "gw1")
		assert.NoError(t, err)
		test_utils.AssertMapsEqual(t, expected, actual)
	})

	t.Run("check default values for log throttling", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{}}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  "gw1",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				DynamicServices:         nil,
				FeatureFlags:            nil,
				Logging: &models.GatewayLoggingConfigs{
					Aggregation: &models.AggregationLoggingConfigs{
						TargetFilesByTag: map[string]string{
							"thing": "/var/log/thing.log",
							"blah":  "/some/directory/blah.log",
						},
						// No throttle values
					},
				},
			},
		}
		tier := configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType,
			Key:  "default",
			Config: &models.Tier{
				Name:    "default",
				Version: &version,
				Images: []*models.TierImage{
					{Name: swag.String("Image1"), Order: swag.Int64(42)},
					{Name: swag.String("Image2"), Order: swag.Int64(1)},
				},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw, tier},
			Edges: []configurator.GraphEdge{
				{From: tier.GetTK(), To: gw.GetTK()},
			},
		}

		expected := map[string]proto.Message{
			"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
			"magmad": &mconfig_protos.MagmaD{
				LogLevel:                protos.LogLevel_INFO,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				AutoupgradeEnabled:      true,
				AutoupgradePollInterval: 300,
				PackageVersion:          "1.0.0-0",
				Images: []*mconfig_protos.ImageSpec{
					{Name: "Image1", Order: 42},
					{Name: "Image2", Order: 1},
				},
				DynamicServices: nil,
				FeatureFlags:    nil,
			},
			"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
			"td-agent-bit": &mconfig_protos.FluentBit{
				ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw1"},
				ThrottleRate:     1000,
				ThrottleWindow:   5,
				ThrottleInterval: "1m",
				FilesByTag: map[string]string{
					"thing": "/var/log/thing.log",
					"blah":  "/some/directory/blah.log",
				},
			},
			"eventd": &mconfig_protos.EventD{
				LogLevel:       protos.LogLevel_INFO,
				EventVerbosity: -1,
			},
			"state": &mconfig_protos.State{
				SyncInterval: expectedDefaultJitteredSyncIntervalGW1,
				LogLevel:     protos.LogLevel_INFO,
			},
			"shared_mconfig": &mconfig_protos.SharedMconfig{
				SentryConfig: nil,
			},
		}

		actual, err := buildBaseOrchestrator(&nw, &graph, "gw1")
		assert.NoError(t, err)
		test_utils.AssertMapsEqual(t, expected, actual)
	})

	// Test sync interval jitter
	t.Run("sync interval jitter works correctly", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": &models.StateConfig{
				SyncInterval: syncInterval,
			},
		}}
		gw := configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  "gw2",
			Config: &models.MagmadGatewayConfigs{
				AutoupgradeEnabled:      swag.Bool(true),
				AutoupgradePollInterval: 300,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				DynamicServices:         []string{},
				FeatureFlags:            map[string]bool{},
			},
		}
		tier := configurator.NetworkEntity{
			Type: orc8r.UpgradeTierEntityType,
			Key:  "default",
			Config: &models.Tier{
				Name:    "default",
				Version: &version,
				Images: []*models.TierImage{
					{Name: swag.String("Image1"), Order: swag.Int64(42)},
					{Name: swag.String("Image2"), Order: swag.Int64(1)},
				},
			},
		}
		graph := configurator.EntityGraph{
			Entities: []configurator.NetworkEntity{gw, tier},
			Edges: []configurator.GraphEdge{
				{From: tier.GetTK(), To: gw.GetTK()},
			},
		}

		expected := map[string]proto.Message{
			"control_proxy": &mconfig_protos.ControlProxy{LogLevel: protos.LogLevel_INFO},
			"magmad": &mconfig_protos.MagmaD{
				LogLevel:                protos.LogLevel_INFO,
				CheckinInterval:         60,
				CheckinTimeout:          10,
				AutoupgradeEnabled:      true,
				AutoupgradePollInterval: 300,
				PackageVersion:          "1.0.0-0",
				Images: []*mconfig_protos.ImageSpec{
					{Name: "Image1", Order: 42},
					{Name: "Image2", Order: 1},
				},
				DynamicServices: nil,
				FeatureFlags:    nil,
			},
			"metricsd": &mconfig_protos.MetricsD{LogLevel: protos.LogLevel_INFO},
			"td-agent-bit": &mconfig_protos.FluentBit{
				ExtraTags:        map[string]string{"network_id": "n1", "gateway_id": "gw2"},
				ThrottleRate:     1000,
				ThrottleWindow:   5,
				ThrottleInterval: "1m",
			},
			"eventd": &mconfig_protos.EventD{
				LogLevel:       protos.LogLevel_INFO,
				EventVerbosity: -1,
			},
			"state": &mconfig_protos.State{
				SyncInterval: expectedJitteredSyncIntervalGW2,
				LogLevel:     protos.LogLevel_INFO,
			},
			"shared_mconfig": &mconfig_protos.SharedMconfig{
				SentryConfig: nil,
			},
		}

		actual, err := buildBaseOrchestrator(&nw, &graph, "gw2")
		assert.NoError(t, err)
		test_utils.AssertMapsEqual(t, expected, actual)
	})
}

func TestGetStateMconfig(t *testing.T) {
	syncInterval := uint32(30)
	expectedJitteredSyncIntervalGW1 := uint32(35)
	expectedJitteredSyncIntervalGW3 := uint32(32)
	expectedDefaultJitteredSyncIntervalGW1 := uint32(71)

	t.Run("interval set in net config", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": &models.StateConfig{
				SyncInterval: syncInterval,
			},
		}}
		gwKey := "gw1"
		expected := &mconfig_protos.State{
			LogLevel:     protos.LogLevel_INFO,
			SyncInterval: expectedJitteredSyncIntervalGW1,
		}

		actual := servicers.GetStateMconfig(nw, gwKey)
		assert.Equal(t, expected, actual)
	})
	t.Run("test jitter", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": &models.StateConfig{
				SyncInterval: syncInterval,
			},
		}}
		gwKey := "gw3"
		expected := &mconfig_protos.State{
			LogLevel:     protos.LogLevel_INFO,
			SyncInterval: expectedJitteredSyncIntervalGW3,
		}

		actual := servicers.GetStateMconfig(nw, gwKey)
		assert.Equal(t, expected, actual)
	})
	t.Run("interval not set in state config", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{}}
		gwKey := "gw1"
		expected := &mconfig_protos.State{
			LogLevel:     protos.LogLevel_INFO,
			SyncInterval: expectedDefaultJitteredSyncIntervalGW1,
		}

		actual := servicers.GetStateMconfig(nw, gwKey)
		assert.Equal(t, expected, actual)
	})
	t.Run("failed to cast state config", func(t *testing.T) {
		nw := configurator.Network{ID: "n1", Configs: map[string]interface{}{
			"state_config": nil,
		}}
		gwKey := "gw1"
		expected := &mconfig_protos.State{
			LogLevel:     protos.LogLevel_INFO,
			SyncInterval: expectedDefaultJitteredSyncIntervalGW1,
		}

		actual := servicers.GetStateMconfig(nw, gwKey)
		assert.Equal(t, expected, actual)
	})
}

func buildBaseOrchestrator(network *configurator.Network, graph *configurator.EntityGraph, gatewayID string) (map[string]proto.Message, error) {
	networkProto, err := network.ToProto(serdes.Network)
	if err != nil {
		return nil, err
	}
	graphProto, err := graph.ToProto(serdes.Entity)
	if err != nil {
		return nil, err
	}
	builder := mconfig.NewRemoteBuilder(orchestrator.ServiceName)
	res, err := builder.Build(networkProto, graphProto, gatewayID)
	if err != nil {
		return nil, err
	}

	configs, err := mconfig.UnmarshalConfigs(res)
	if err != nil {
		return nil, err
	}

	// Only return configs relevant to base orc8r
	ret := map[string]proto.Message{
		"control_proxy":  configs["control_proxy"],
		"metricsd":       configs["metricsd"],
		"state":          configs["state"],
		"shared_mconfig": configs["shared_mconfig"],
	}
	_, ok := configs["magmad"]
	if ok {
		ret["magmad"] = configs["magmad"]
		ret["td-agent-bit"] = configs["td-agent-bit"]
		ret["eventd"] = configs["eventd"]
	}

	return ret, nil
}
