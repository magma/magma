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
	"magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	assert "github.com/stretchr/testify/require"
	"github.com/thoas/go-funk"
)

func TestSubscriberdbStreamer(t *testing.T) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &lte_plugin.LteOrchestratorPlugin{})) // load remote providers
	lte_test_init.StartTestService(t)
	configurator_test_init.StartTestService(t)

	provider, err := providers.GetStreamProvider(lte.SubscriberStreamName)
	assert.NoError(t, err)

	err = configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	// 1 sub without a profile on the backend (should fill as "default"), the
	// other inactive with a sub profile
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.SubscriberEntityType, Key: "IMSI12345",
			Config: &models.LteSubscription{
				State:   "ACTIVE",
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
		},
		{Type: lte.SubscriberEntityType, Key: "IMSI67890", Config: &models.LteSubscription{State: "INACTIVE", SubProfile: "foo"}},
	})
	assert.NoError(t, err)

	expectedProtos := []*lte_protos.SubscriberData{
		{
			Sid: &lte_protos.SubscriberID{Id: "12345", Type: lte_protos.SubscriberID_IMSI},
			Lte: &lte_protos.LTESubscription{
				State:   lte_protos.LTESubscription_ACTIVE,
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "default",
		},
		{
			Sid:        &lte_protos.SubscriberID{Id: "67890", Type: lte_protos.SubscriberID_IMSI},
			Lte:        &lte_protos.LTESubscription{State: lte_protos.LTESubscription_INACTIVE},
			NetworkId:  &protos.NetworkID{Id: "n1"},
			SubProfile: "foo",
		},
	}
	expected := funk.Map(
		expectedProtos,
		func(sub *lte_protos.SubscriberData) *protos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err := provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Create policies and base name associated to sub
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.BaseNameEntityType, Key: "bn1",
			Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI12345"}},
		},
		{
			Type: lte.PolicyRuleEntityType, Key: "r1",
			Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI12345"}},
		},
		{
			Type: lte.PolicyRuleEntityType, Key: "r2",
			Associations: []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: "IMSI12345"}},
		},
	})
	assert.NoError(t, err)

	expectedProtos[0].Lte.AssignedPolicies = []string{"r1", "r2"}
	expectedProtos[0].Lte.AssignedBaseNames = []string{"bn1"}
	expected = funk.Map(
		expectedProtos,
		func(sub *lte_protos.SubscriberData) *protos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &protos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err = provider.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
