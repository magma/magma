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
	id1 := "r1"
	monitoringKey1 := swag.String("foo")
	id2 := "r2"
	priority2 := swag.Uint32(42)
	id3 := "r3"
	monitoringKey3 := swag.String("bar")

	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r1",
			Config: &models.PolicyRule{
				ID:            &id1,
				MonitoringKey: *monitoringKey1,
			},
		},
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r2",
			Config: &models.PolicyRule{
				ID:       &id2,
				Priority: priority2,
			},
		},
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r3",
			Config: &models.PolicyRule{
				ID:            &id3,
				MonitoringKey: *monitoringKey3,
			},
		},
	})
	assert.NoError(t, err)
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type:   lte.BaseNameEntityType,
			Key:    "b1",
			Config: &models.BaseNameRecord{Name: models.BaseName("b1")},
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "r1"},
				{Type: lte.PolicyRuleEntityType, Key: "r2"},
			},
		},
		{
			Type:   lte.BaseNameEntityType,
			Key:    "b2",
			Config: &models.BaseNameRecord{Name: models.BaseName("b2")},
			Associations: []storage.TypeAndKey{
				{Type: lte.PolicyRuleEntityType, Key: "r3"},
			},
		},
	})
	assert.NoError(t, err)

	policyPro := &pdbstreamer.PoliciesProvider{}
	expectedProtos := []*protos.PolicyRule{
		{Id: "r1", MonitoringKey: "foo"},
		{Id: "r2", Priority: 42},
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
