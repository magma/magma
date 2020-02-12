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
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/protos"
	sdbstreamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	cfg_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func TestSubscriberdbStreamer(t *testing.T) {
	cfg_test_init.StartTestService(t)
	_ = plugin.RegisterPluginForTests(t, &plugin2.LteOrchestratorPlugin{})

	err := configurator.CreateNetwork(configurator.Network{ID: "n1"})
	assert.NoError(t, err)
	_, err = configurator.CreateEntity("n1", configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "g1", PhysicalID: "hw1"})
	assert.NoError(t, err)

	// 1 sub without a profile on the backend (should fill as "default"), the
	// other inactive with a sub profile
	_, err = configurator.CreateEntities("n1", []configurator.NetworkEntity{
		{
			Type: lte.SubscriberEntityType, Key: "IMSI12345",
			Config: &models2.LteSubscription{
				State:   "ACTIVE",
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
		},
		{Type: lte.SubscriberEntityType, Key: "IMSI67890", Config: &models2.LteSubscription{State: "INACTIVE", SubProfile: "foo"}},
	})
	assert.NoError(t, err)

	pro := &sdbstreamer.SubscribersProvider{}
	expectedProtos := []*protos.SubscriberData{
		{
			Sid: &protos.SubscriberID{Id: "12345", Type: protos.SubscriberID_IMSI},
			Lte: &protos.LTESubscription{
				State:   protos.LTESubscription_ACTIVE,
				AuthKey: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
				AuthOpc: []byte("\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22\x22"),
			},
			NetworkId:  &orcprotos.NetworkID{Id: "n1"},
			SubProfile: "default",
		},
		{
			Sid:        &protos.SubscriberID{Id: "67890", Type: protos.SubscriberID_IMSI},
			Lte:        &protos.LTESubscription{State: protos.LTESubscription_INACTIVE},
			NetworkId:  &orcprotos.NetworkID{Id: "n1"},
			SubProfile: "foo",
		},
	}
	expected := funk.Map(
		expectedProtos,
		func(sub *protos.SubscriberData) *orcprotos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &orcprotos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err := pro.GetUpdates("hw1", nil)
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
		func(sub *protos.SubscriberData) *orcprotos.DataUpdate {
			data, err := proto.Marshal(sub)
			assert.NoError(t, err)
			return &orcprotos.DataUpdate{Key: "IMSI" + sub.Sid.Id, Value: data}
		},
	)
	actual, err = pro.GetUpdates("hw1", nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
