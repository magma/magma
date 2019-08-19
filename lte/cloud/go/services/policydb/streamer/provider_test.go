/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streamer_test

import (
	"encoding/json"
	"os"
	"testing"

	"magma/lte/cloud/go/lte"
	plugin2 "magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	pdbstreamer "magma/lte/cloud/go/services/policydb/streamer"
	policydb_test_init "magma/lte/cloud/go/services/policydb/test_init"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/services/configurator"
	cfg_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/streamer"
	streamer_test_init "magma/orc8r/cloud/go/services/streamer/test_init"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
	"golang.org/x/net/context"
)

const testAgHwId = "Test-AGW-Hw-Id"

func TestPolicyStreamers(t *testing.T) {
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "1")
	cfg_test_init.StartTestService(t)
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
			Config: &models.PolicyRule{
				ID:            "r1",
				MonitoringKey: strPtr("foo"),
			},
		},
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r2",
			Config: &models.PolicyRule{
				ID:       "r2",
				Priority: uintPtr(42),
			},
		},
		{
			Type: lte.PolicyRuleEntityType,
			Key:  "r3",
			Config: &models.PolicyRule{
				ID:            "r3",
				MonitoringKey: strPtr("bar"),
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

func TestPolicydbStreamer_Legacy(t *testing.T) {
	// Setup - start services, register provider
	_ = os.Setenv(orc8r.UseConfiguratorEnv, "0")
	magmad_test_init.StartTestService(t)
	policydb_test_init.StartTestService(t)
	streamer_test_init.StartTestService(t)
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	testNetworkId, err := magmad.RegisterNetwork(&magmad_protos.MagmadNetworkRecord{Name: "Test Network 1"}, "policydb_streamer_test_network")
	assert.NoError(t, err)

	hwId1 := orcprotos.AccessGatewayID{Id: testAgHwId}
	_, err = magmad.RegisterGateway(testNetworkId, &magmad_protos.AccessGatewayRecord{HwId: &hwId1, Name: "bla"})
	assert.NoError(t, err)

	rule1 := &protos.PolicyRule{
		Id: "1",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{TcpSrc: 0},
			},
		},
		Priority: 10,
	}
	rule2 := &protos.PolicyRule{
		Id: "2",
		FlowList: []*protos.FlowDescription{
			{
				Match: &protos.FlowMatch{IpProto: 7},
			},
		},
		Priority: 15,
	}

	// Add policies
	err = policydb.AddRule(testNetworkId, rule1)
	assert.NoError(t, err)
	err = policydb.AddRule(testNetworkId, rule2)
	assert.NoError(t, err)

	policies, err := policydb.ListRuleIds(testNetworkId)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(policies))

	conn, err := registry.GetConnection(streamer.ServiceName)
	assert.NoError(t, err)

	grpcClient := orcprotos.NewStreamerClient(conn)
	streamerClient, err := grpcClient.GetUpdates(
		context.Background(),
		&orcprotos.StreamRequest{GatewayId: testAgHwId, StreamName: "policydb"},
	)
	assert.NoError(t, err)

	updateBatch, err := streamerClient.Recv()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(updateBatch.GetUpdates()))
	var p1, p2 protos.PolicyRule
	err = proto.Unmarshal(updateBatch.Updates[0].Value, &p1)
	assert.NoError(t, err)
	err = proto.Unmarshal(updateBatch.Updates[1].Value, &p2)
	assert.NoError(t, err)
	p1j, _ := json.Marshal(p1)
	p2j, _ := json.Marshal(p2)
	t.Logf("\nReceived Policies:\n\t%s\n\t%s", string(p1j), string(p2j))
}

func strPtr(s string) *string {
	return &s
}

func uintPtr(i uint32) *uint32 {
	return &i
}
