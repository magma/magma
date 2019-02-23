/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package xservice_test

import (
	"encoding/json"
	"testing"

	"magma/orc8r/cloud/go/protos"
	checkinti "magma/orc8r/cloud/go/services/checkind/test_init"
	"magma/orc8r/cloud/go/services/checkind/test_utils"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/registry"
	configti "magma/orc8r/cloud/go/services/config/test_init"
	"magma/orc8r/cloud/go/services/magmad"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	magmadti "magma/orc8r/cloud/go/services/magmad/test_init"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	storagetu "magma/orc8r/cloud/go/services/materializer/gateways/storage/test_utils"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/xservice"

	"github.com/stretchr/testify/assert"
)

var cfg1 = &storagetu.Conf1{Value1: 1, Value2: "foo", Value3: []byte("bar")}
var cfg2 = &storagetu.Conf2{Value1: []string{"foo", "bar"}, Value2: 1}

func TestXserviceStorage_GetGatewayViewsForNetwork(t *testing.T) {
	// Test setup
	magmadti.StartTestService(t)
	configti.StartTestService(t)
	checkinti.StartTestService(t)

	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(storagetu.NewConfig1Manager())
	registry.RegisterConfigManager(storagetu.NewConfig2Manager())

	// Setup fixture data
	networkID, err := magmad.RegisterNetwork(&magmadprotos.MagmadNetworkRecord{Name: "foobar"}, "xservice1")
	assert.NoError(t, err)

	// Register gateways
	record1 := &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{Id: "hw1"},
	}
	record2 := &magmadprotos.AccessGatewayRecord{
		HwId: &protos.AccessGatewayID{Id: "hw2"},
	}
	_, err = magmad.RegisterGatewayWithId(networkID, record1, "gw1")
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(networkID, record2, "gw2")
	assert.NoError(t, err)

	// gw1 has cfg1 and cfg2, gw2 only has cfg2
	config.CreateConfig(
		networkID,
		storagetu.NewConfig1Manager().GetConfigType(),
		"gw1",
		cfg1,
	)
	config.CreateConfig(
		networkID,
		storagetu.NewConfig2Manager().GetConfigType(),
		"gw1",
		cfg2,
	)
	config.CreateConfig(
		networkID,
		storagetu.NewConfig2Manager().GetConfigType(),
		"gw2",
		cfg2,
	)

	// gw1 has status, gw2 does not
	checkinReq := &protos.CheckinRequest{
		GatewayId:       "hw1",
		MagmaPkgVersion: "1.2.3",
		Status: &protos.ServiceStatus{
			Meta: map[string]string{
				"hello": "world",
			},
		},
		SystemStatus: &protos.SystemStatus{
			CpuUser:   31498,
			CpuSystem: 8361,
			CpuIdle:   1869111,
			MemTotal:  1016084,
			MemUsed:   54416,
			MemFree:   412772,
		},
	}
	test_utils.Checkin(t, checkinReq)

	store := xservice.NewCrossServiceGatewayViewsStorage()
	actual, err := store.GetGatewayViewsForNetwork(networkID)
	assert.NoError(t, err)
	// Wipe out timestamps from status so we can compare the structs
	for _, state := range actual {
		if state.Status != nil {
			state.Status.CertExpirationTime = 0
			state.Status.Time = 0
		}
	}

	expected := map[string]*storage.GatewayState{
		"gw1": {
			GatewayID: "gw1",
			Config: map[string]interface{}{
				storagetu.NewConfig1Manager().GetConfigType(): cfg1,
				storagetu.NewConfig2Manager().GetConfigType(): cfg2,
			},
			Record: record1,
			Status: &protos.GatewayStatus{
				Checkin:            checkinReq,
				Time:               0,
				CertExpirationTime: 0,
			},
		},
		"gw2": {
			GatewayID: "gw2",
			Config: map[string]interface{}{
				storagetu.NewConfig2Manager().GetConfigType(): cfg2,
			},
			Record: record2,
		},
	}

	marshaledExpected, err := json.Marshal(expected)
	assert.NoError(t, err)
	marshaledActual, err := json.Marshal(actual)
	assert.NoError(t, err)
	assert.Equal(t, marshaledExpected, marshaledActual)
}
