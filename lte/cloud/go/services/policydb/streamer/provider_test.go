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
	plugin2 "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/protos"
	pdbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorTestInit "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func TestPolicyStreamers(t *testing.T) {
	configuratorTestInit.StartTestService(t)
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
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

	policyPro := &pdbstreamer.PoliciesProvider{}
	expectedProtos := []*protos.PolicyRule{
		{
			Id:            "r1",
			MonitoringKey: "foo",
			FlowList: []*protos.FlowDescription{
				{
					Match: &protos.FlowMatch{
						Direction: protos.FlowMatch_UPLINK,
						IpProto:   protos.FlowMatch_IPPROTO_IP,
						Ipv4Dst:   "192.168.160.0/24",
						Ipv4Src:   "192.168.128.0/24",
					},
					Action: protos.FlowDescription_PERMIT,
				},
			},
		},
		{
			Id:       "r2",
			Priority: 42,
			Redirect: &protos.RedirectInformation{
				Support:       protos.RedirectInformation_ENABLED,
				AddressType:   protos.RedirectInformation_IPv4,
				ServerAddress: "https://www.google.com",
			},
		},
		{Id: "r3", MonitoringKey: "bar"},
	}
	expected := funk.Map(
		expectedProtos,
		func(r *protos.PolicyRule) *orcprotos.DataUpdate {
			data, err := proto.Marshal(r)
			assert.NoError(t, err)
			return &orcprotos.DataUpdate{Key: r.Id, Value: data}
		},
	)
	actual, err := policyPro.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	bnPro := &pdbstreamer.BaseNamesProvider{}
	expectedBNProtos := []*protos.ChargingRuleBaseNameRecord{
		{Name: "b1", RuleNamesSet: &protos.ChargingRuleNameSet{RuleNames: []string{"r1", "r2"}}},
		{Name: "b2", RuleNamesSet: &protos.ChargingRuleNameSet{RuleNames: []string{"r3"}}},
	}
	expected = funk.Map(
		expectedBNProtos,
		func(bn *protos.ChargingRuleBaseNameRecord) *orcprotos.DataUpdate {
			data, err := proto.Marshal(bn)
			assert.NoError(t, err)
			return &orcprotos.DataUpdate{Key: bn.Name, Value: data}
		},
	)
	actual, err = bnPro.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
