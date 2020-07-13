/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	lte_plugin "magma/lte/cloud/go/plugin"
	lte_protos "magma/lte/cloud/go/protos"
	lte_test_init "magma/lte/cloud/go/services/lte/test_init"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/lte/cloud/go/services/policydb/streamer"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func TestRatingGroupStreamers(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{})) // load remote providers
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	provider, err := providers.GetStreamProvider(lte.RatingGroupStreamName)
	assert.NoError(t, err)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	// create the rating groups
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.RatingGroupEntityType,
			Key:  "111",
			Config: &models.RatingGroup{
				ID:        111,
				LimitType: swag.String("FINITE"),
			},
		},
		{
			Type: lte.RatingGroupEntityType,
			Key:  "222",
			Config: &models.RatingGroup{
				ID:        222,
				LimitType: swag.String("INFINITE_METERED"),
			},
		},
		{
			Type: lte.RatingGroupEntityType,
			Key:  "333",
			Config: &models.RatingGroup{
				ID:        333,
				LimitType: swag.String("INFINITE_UNMETERED"),
			},
		},
	})
	assert.NoError(t, err)

	expectedProtos := []*lte_protos.RatingGroup{
		{
			Id:        111,
			LimitType: lte_protos.RatingGroup_FINITE,
		},
		{
			Id:        222,
			LimitType: lte_protos.RatingGroup_INFINITE_METERED,
		},
		{
			Id:        333,
			LimitType: lte_protos.RatingGroup_INFINITE_UNMETERED,
		},
	}
	expected := funk.Map(
		expectedProtos,
		func(r *lte_protos.RatingGroup) *protos.DataUpdate {
			data, err := proto.Marshal(r)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: swag.FormatUint32(r.Id), Value: data}
		},
	)
	actual, err := provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestPolicyStreamers(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{})) // load remote providers
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	provider, err := providers.GetStreamProvider(lte.PolicyStreamName)
	assert.NoError(t, err)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	// create the rules first otherwise base names can't associate to them
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r1",
			Config: &models.PolicyRuleConfig{
				FlowList: []*models.FlowDescription{
					{
						Action: swag.String("PERMIT"),
						Match: &models.FlowMatch{
							Direction: swag.String("UPLINK"),
							IPProto:   swag.String("IPPROTO_IP "),
							IPV4Dst:   "192.168.160.0/24",
							IPV4Src:   "192.168.128.0/24",
						},
					},
				},
				MonitoringKey: "foo",
			},
		},
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r2",
			Config: &models.PolicyRuleConfig{
				Priority: swag.Uint32(42),
				Redirect: &models.RedirectInformation{
					AddressType:   swag.String("IPv4"),
					ServerAddress: swag.String("https://www.google.com"),
					Support:       swag.String("ENABLED"),
				},
			},
		},
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r3",
			Config: &models.PolicyRuleConfig{
				MonitoringKey: "bar",
			},
		},
	})
	assert.NoError(t, err)
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type:   lte.BaseNameEntityType,
			Key:    "b1",
			Config: &models.BaseNameRecord{Name: "b1"},
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "r1"},
				{Type: lte.PolicyRuleEntityType, Key: "r2"},
			},
		},
		{
			Type:   lte.BaseNameEntityType,
			Key:    "b2",
			Config: &models.BaseNameRecord{Name: "b2"},
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "r3"},
			},
		},
	})
	assert.NoError(t, err)

	expectedProtos := []*lte_protos.PolicyRule{
		{
			Id:            "r1",
			MonitoringKey: []byte("foo"),
			FlowList: []*lte_protos.FlowDescription{
				{
					Match: &lte_protos.FlowMatch{
						Direction: lte_protos.FlowMatch_UPLINK,
						IpProto:   lte_protos.FlowMatch_IPPROTO_IP,
						Ipv4Dst:   "192.168.160.0/24",
						Ipv4Src:   "192.168.128.0/24",
					},
					Action: lte_protos.FlowDescription_PERMIT,
				},
			},
		},
		{
			Id:       "r2",
			Priority: 42,
			Redirect: &lte_protos.RedirectInformation{
				Support:       lte_protos.RedirectInformation_ENABLED,
				AddressType:   lte_protos.RedirectInformation_IPv4,
				ServerAddress: "https://www.google.com",
			},
		},
		{Id: "r3", MonitoringKey: []byte("bar")},
	}
	expected := funk.Map(
		expectedProtos,
		func(r *lte_protos.PolicyRule) *protos.DataUpdate {
			data, err := proto.Marshal(r)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: r.Id, Value: data}
		},
	).([]*protos.DataUpdate)

	actual, err := provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	bnPro := &streamer.BaseNamesProvider{}
	expectedBNProtos := []*lte_protos.ChargingRuleBaseNameRecord{
		{Name: "b1", RuleNamesSet: &lte_protos.ChargingRuleNameSet{RuleNames: []string{"r1", "r2"}}},
		{Name: "b2", RuleNamesSet: &lte_protos.ChargingRuleNameSet{RuleNames: []string{"r3"}}},
	}
	expected = funk.Map(
		expectedBNProtos,
		func(bn *lte_protos.ChargingRuleBaseNameRecord) *protos.DataUpdate {
			data, err := proto.Marshal(bn.RuleNamesSet)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: bn.Name, Value: data}
		},
	).([]*protos.DataUpdate)

	actual, err = bnPro.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestRuleMappingsProvider(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{})) // load remote providers
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	provider, err := providers.GetStreamProvider(lte.MappingsStreamName)
	assert.NoError(t, err)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.SubscriberEntityType, Key: "s1"},
			{Type: lte.SubscriberEntityType, Key: "s2"},
			{Type: lte.SubscriberEntityType, Key: "s3"},

			// r1 -> s1, r2 -> s2, r3 -> s1,s2
			{Type: lte.PolicyRuleEntityType, Key: "r1", Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "s1"}}},
			{Type: lte.PolicyRuleEntityType, Key: "r2", Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "s2"}}},
			{Type: lte.PolicyRuleEntityType, Key: "r3", Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "s1"}, {Type: lte.SubscriberEntityType, Key: "s2"}}},

			// b1 -> s1, b2 -> s2, b3 -> s1,s2
			{Type: lte.BaseNameEntityType, Key: "b1", Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "s1"}}},
			{Type: lte.BaseNameEntityType, Key: "b2", Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "s2"}}},
			{Type: lte.BaseNameEntityType, Key: "b3", Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "s1"}, {Type: lte.SubscriberEntityType, Key: "s2"}}},
		},
	)
	assert.NoError(t, err)

	expectedProtos := []*lte_protos.AssignedPolicies{
		{
			AssignedBaseNames: []string{"b1", "b3"},
			AssignedPolicies:  []string{"r1", "r3"},
		},
		{
			AssignedBaseNames: []string{"b2", "b3"},
			AssignedPolicies:  []string{"r2", "r3"},
		},
	}
	expected := funk.Map(
		expectedProtos,
		func(ap *lte_protos.AssignedPolicies) *protos.DataUpdate {
			data, err := proto.Marshal(ap)
			assert.NoError(t, err)
			return &protos.DataUpdate{Value: data}
		},
	).([]*protos.DataUpdate)
	expected[0].Key, expected[1].Key = "s1", "s2"

	actual, err := provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNetworkWideRulesProvider(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{})) // load remote providers
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	provider, err := providers.GetStreamProvider(lte.NetworkWideRulesStreamName)
	assert.NoError(t, err)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	_, err = configurator.CreateEntities(
		"n1",
		[]configurator.NetworkEntity{
			{Type: lte.PolicyRuleEntityType, Key: "r1"},
			{Type: lte.PolicyRuleEntityType, Key: "r2"},
			{Type: lte.PolicyRuleEntityType, Key: "r3"},

			{Type: lte.BaseNameEntityType, Key: "b1"},
			{Type: lte.BaseNameEntityType, Key: "b2"},
			{Type: lte.BaseNameEntityType, Key: "b3"},
		},
	)
	assert.NoError(t, err)
	config := &models.NetworkSubscriberConfig{
		NetworkWideBaseNames: []models.BaseName{"b1", "b2"},
		NetworkWideRuleNames: []string{"r1", "r2"},
	}
	assert.NoError(t, configurator.UpdateNetworkConfig("n1", lte.NetworkSubscriberConfigType, config))

	expectedProtos := []*lte_protos.AssignedPolicies{
		{
			AssignedBaseNames: []string{"b1", "b2"},
			AssignedPolicies:  []string{"r1", "r2"},
		},
	}
	expected := funk.Map(
		expectedProtos,
		func(ap *lte_protos.AssignedPolicies) *protos.DataUpdate {
			data, err := proto.Marshal(ap)
			assert.NoError(t, err)
			return &protos.DataUpdate{Value: data}
		},
	).([]*protos.DataUpdate)

	actual, err := provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
