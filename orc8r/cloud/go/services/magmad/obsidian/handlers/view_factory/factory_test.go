/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package view_factory_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorti "magma/orc8r/cloud/go/services/configurator/test_init"
	configuratortu "magma/orc8r/cloud/go/services/configurator/test_utils"
	"magma/orc8r/cloud/go/services/device"
	deviceti "magma/orc8r/cloud/go/services/device/test_init"
	storagetu "magma/orc8r/cloud/go/services/magmad/obsidian/handlers/test_utils"
	"magma/orc8r/cloud/go/services/magmad/obsidian/handlers/view_factory"
	"magma/orc8r/cloud/go/services/state"
	stateti "magma/orc8r/cloud/go/services/state/test_init"
	statetu "magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/storage"

	"github.com/stretchr/testify/assert"
)

var cfg1 = &storagetu.Conf1{Value1: 1, Value2: "foo", Value3: []byte("bar")}
var cfg2 = &storagetu.Conf2{Value1: []string{"foo", "bar"}, Value2: 1}

func TestFullGatewayViewFactoryImpl_GetGatewayViewsForNetwork(t *testing.T) {
	// Test setup
	configuratorti.StartTestService(t)
	deviceti.StartTestService(t)
	stateti.StartTestService(t)

	err := serde.RegisterSerdes(
		storagetu.NewConfig1ConfiguratorManager(),
		storagetu.NewConfig2ConfiguratorManager(),
		state.NewStateSerde(orc8r.GatewayStateType, &models.GatewayStatus{}),
		serde.NewBinarySerde(device.SerdeDomain, orc8r.AccessGatewayRecordType, &models.GatewayDevice{}),
	)
	assert.NoError(t, err)

	// Setup fixture data
	networkID := "test_network"
	gatewayID1 := "gw1"
	gatewayID2 := "gw2"
	hwID1 := "hw1"
	hwID2 := "hw2"
	record1 := &models.GatewayDevice{HardwareID: hwID1}
	record2 := &models.GatewayDevice{HardwareID: hwID2}
	configuratortu.RegisterNetwork(t, networkID, "xservice1")
	configuratortu.RegisterGatewayWithName(t, networkID, gatewayID1, "111", record1)
	configuratortu.RegisterGatewayWithName(t, networkID, gatewayID2, "222", record2)

	// configs for gw1
	gw1config1 := configurator.NetworkEntity{
		Type:   storagetu.NewConfig1Manager().GetType(),
		Key:    gatewayID1,
		Config: cfg1,
	}
	gw1config2 := configurator.NetworkEntity{
		Type:   storagetu.NewConfig2Manager().GetType(),
		Key:    gatewayID1,
		Config: cfg2,
	}
	// configs for gw2
	gw2config1 := configurator.NetworkEntity{
		Type:   storagetu.NewConfig1Manager().GetType(),
		Key:    gatewayID2,
		Config: cfg1,
	}

	_, err = configurator.CreateEntities(networkID, []configurator.NetworkEntity{gw1config1, gw1config2, gw2config1})
	assert.NoError(t, err)

	// add config associations to gateways
	// gw1 has cfg1 and cfg2, gw2 only has cfg1
	updateGW1 := configurator.EntityUpdateCriteria{
		Type: orc8r.MagmadGatewayType,
		Key:  gatewayID1,
		AssociationsToSet: []storage.TypeAndKey{
			{Type: gw1config1.Type, Key: gatewayID1},
			{Type: gw1config2.Type, Key: gatewayID1}},
	}
	updateGW2 := configurator.EntityUpdateCriteria{
		Type: orc8r.MagmadGatewayType,
		Key:  gatewayID2,
		AssociationsToSet: []storage.TypeAndKey{
			{Type: gw2config1.Type, Key: gatewayID2}},
	}

	_, err = configurator.UpdateEntities(networkID, []configurator.EntityUpdateCriteria{updateGW1, updateGW2})
	assert.NoError(t, err)

	// put status into gw1
	ctx := statetu.GetContextWithCertificate(t, hwID1)
	gwStatus := models.NewDefaultGatewayStatus(hwID1)
	statetu.ReportGatewayStatus(t, ctx, gwStatus)

	fact := &view_factory.FullGatewayViewFactoryImpl{}
	actual, err := fact.GetGatewayViewsForNetwork(networkID)
	assert.NoError(t, err)
	// Wipe out timestamps from status so we can compare the structs
	for _, state := range actual {
		if state.Status != nil {
			state.Status.CertExpirationTime = 0
			state.Status.CheckinTime = 0
		}
	}

	expected := map[string]*view_factory.GatewayState{
		gatewayID1: {
			GatewayID: gatewayID1,
			Config: map[string]interface{}{
				storagetu.NewConfig1ConfiguratorManager().GetType(): cfg1,
				storagetu.NewConfig2ConfiguratorManager().GetType(): cfg2,
				orc8r.MagmadGatewayType:                             nil,
			},
			Name:   "111",
			Status: models.NewDefaultGatewayStatus(hwID1),
			Record: record1,
		},
		gatewayID2: {
			GatewayID: gatewayID2,
			Config: map[string]interface{}{
				storagetu.NewConfig1ConfiguratorManager().GetType(): cfg1,
				orc8r.MagmadGatewayType:                             nil,
			},
			Name:   "222",
			Record: record2,
		},
	}

	assert.Equal(t, expected, actual)

	// add an unrelated entity to gw1 and make sure only the config entities are loaded
	nonConfigEntity, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Key:    "random_entity",
		Type:   storagetu.NewConfig1Manager().GetType(),
		Config: cfg1,
	})
	assert.NoError(t, err)

	// add association from gw1 -> nonConfigEntity
	updateGW1.AssociationsToAdd = []storage.TypeAndKey{{Type: nonConfigEntity.Type, Key: nonConfigEntity.Key}}
	updateGW1.AssociationsToSet = nil
	_, err = configurator.UpdateEntities(networkID, []configurator.EntityUpdateCriteria{updateGW1})
	assert.NoError(t, err)

	actual, err = fact.GetGatewayViewsForNetwork(networkID)
	assert.NoError(t, err)
	// Wipe out timestamps from status so we can compare the structs
	for _, state := range actual {
		if state.Status != nil {
			state.Status.CertExpirationTime = 0
			state.Status.CheckinTime = 0
		}
	}
	// result should be the same as before, ignoring the non config ents
	assert.Equal(t, expected, actual)
}
