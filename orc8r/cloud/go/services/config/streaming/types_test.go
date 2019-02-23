/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config/streaming"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/config/streaming/storage/mocks"
	"magma/orc8r/cloud/go/services/magmad"
	protos2 "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"
	upgradeprotos "magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTierUpdate_Apply_Create(t *testing.T) {
	// 3 gateways - 1 with matching tier, 1 with non-matching tier, 1 without
	// an mconfig

	// Update should only be applied to gw1
	magmad_test_init.StartTestService(t)
	_, err := magmad.RegisterNetwork(
		&protos2.MagmadNetworkRecord{Name: "nw"},
		"tierupdate_network")
	assert.NoError(t, err)

	_, err = magmad.RegisterGatewayWithId(
		"tierupdate_network",
		&protos2.AccessGatewayRecord{Name: "gw1", HwId: &protos.AccessGatewayID{Id: "id1"}},
		"gw1",
	)
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(
		"tierupdate_network",
		&protos2.AccessGatewayRecord{Name: "gw2", HwId: &protos.AccessGatewayID{Id: "id2"}},
		"gw2",
	)
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(
		"tierupdate_network",
		&protos2.AccessGatewayRecord{Name: "gw3", HwId: &protos.AccessGatewayID{Id: "id3"}},
		"gw3",
	)
	assert.NoError(t, err)

	mconfigFixtures := map[string]map[string]proto.Message{
		"gw1": {
			"magmad": &mconfig.MagmaD{
				TierId:         "thistier",
				PackageVersion: "0.0.0-0",
				Images:         []*mconfig.ImageSpec{},
			},
		},
		"gw2": {
			"magmad": &mconfig.MagmaD{
				TierId:         "not_thistier",
				PackageVersion: "1.0.0-0",
				Images:         []*mconfig.ImageSpec{{Name: "img", Order: 1}},
			},
		},
	}

	mockStorage := &mocks.MconfigStorage{}
	mockStorage.On("GetMconfigs", "tierupdate_network", mock.Anything).
		Return(getStorageFixtures(t, "tierupdate_network", mconfigFixtures), nil)
	expectedGw1Mconfig := map[string]proto.Message{
		"magmad": &mconfig.MagmaD{
			TierId:         "thistier",
			PackageVersion: "1.2.3-4",
			Images:         []*mconfig.ImageSpec{{Name: "foo", Order: 2}},
		},
	}
	expectedUpdates := []*storage.MconfigUpdateCriteria{
		{
			GatewayId:  "gw1",
			Offset:     3,
			NewMconfig: getMconfig(t, expectedGw1Mconfig),
		},
	}
	mockStorage.On("CreateOrUpdateMconfigs", "tierupdate_network", expectedUpdates).Return(nil)

	tierUpdate := &streaming.TierUpdate{
		NetworkId:   "tierupdate_network",
		Operation:   streaming.CreateOperation,
		TierId:      "thistier",
		TierVersion: "1.2.3-4",
		TierImages:  []*upgradeprotos.ImageSpec{{Name: "foo", Order: 2}},
	}
	err = tierUpdate.Apply(mockStorage, 100)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
	mockStorage.AssertNumberOfCalls(t, "GetMconfigs", 1)
	mockStorage.AssertNumberOfCalls(t, "CreateOrUpdateMconfigs", 1)
	mockStorage.AssertCalled(t, "CreateOrUpdateMconfigs", "tierupdate_network", expectedUpdates)
}

func TestTierUpdate_Apply_Delete(t *testing.T) {
	// 3 gateways - 1 with matching tier, 1 with non-matching tier, 1 without
	// an mconfig

	// Update should only be applied to gw1
	magmad_test_init.StartTestService(t)
	_, err := magmad.RegisterNetwork(
		&protos2.MagmadNetworkRecord{Name: "nw"},
		"tierupdate_network")
	assert.NoError(t, err)

	_, err = magmad.RegisterGatewayWithId(
		"tierupdate_network",
		&protos2.AccessGatewayRecord{Name: "gw1", HwId: &protos.AccessGatewayID{Id: "id1"}},
		"gw1",
	)
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(
		"tierupdate_network",
		&protos2.AccessGatewayRecord{Name: "gw2", HwId: &protos.AccessGatewayID{Id: "id2"}},
		"gw2",
	)
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(
		"tierupdate_network",
		&protos2.AccessGatewayRecord{Name: "gw3", HwId: &protos.AccessGatewayID{Id: "id3"}},
		"gw3",
	)
	assert.NoError(t, err)

	mconfigFixtures := map[string]map[string]proto.Message{
		"gw1": {
			"magmad": &mconfig.MagmaD{
				TierId:         "thistier",
				PackageVersion: "1.2.3-4",
				Images:         []*mconfig.ImageSpec{{Name: "foo", Order: 2}},
			},
		},
		"gw2": {
			"magmad": &mconfig.MagmaD{
				TierId:         "not_thistier",
				PackageVersion: "1.0.0-0",
				Images:         []*mconfig.ImageSpec{{Name: "img", Order: 1}},
			},
		},
	}

	mockStorage := &mocks.MconfigStorage{}
	mockStorage.On("GetMconfigs", "tierupdate_network", mock.Anything).
		Return(getStorageFixtures(t, "tierupdate_network", mconfigFixtures), nil)
	expectedGw1Mconfig := map[string]proto.Message{
		"magmad": &mconfig.MagmaD{
			TierId:         "thistier",
			PackageVersion: "0.0.0-0",
			Images:         []*mconfig.ImageSpec{},
		},
	}
	expectedUpdates := []*storage.MconfigUpdateCriteria{
		{
			GatewayId:  "gw1",
			Offset:     3,
			NewMconfig: getMconfig(t, expectedGw1Mconfig),
		},
	}
	mockStorage.On("CreateOrUpdateMconfigs", "tierupdate_network", expectedUpdates).Return(nil)

	tierUpdate := &streaming.TierUpdate{
		NetworkId: "tierupdate_network",
		Operation: streaming.DeleteOperation,
		TierId:    "thistier",
	}
	err = tierUpdate.Apply(mockStorage, 100)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
	mockStorage.AssertNumberOfCalls(t, "GetMconfigs", 1)
	mockStorage.AssertNumberOfCalls(t, "CreateOrUpdateMconfigs", 1)
	mockStorage.AssertCalled(t, "CreateOrUpdateMconfigs", "tierupdate_network", expectedUpdates)
}

func getStorageFixtures(t *testing.T, networkId string, mconfigs map[string]map[string]proto.Message) map[string]*storage.StoredMconfig {
	ret := map[string]*storage.StoredMconfig{}
	for k, v := range mconfigs {
		retMconfig := getMconfig(t, v)
		ret[k] = &storage.StoredMconfig{
			NetworkId: networkId,
			GatewayId: k,
			Mconfig:   retMconfig,
			Offset:    int64(len(k)),
		}
	}
	return ret
}

func getMconfig(t *testing.T, mc map[string]proto.Message) *protos.GatewayConfigs {
	retMconfig := &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	for cfgKey, cfgVal := range mc {
		cfgAny, err := ptypes.MarshalAny(cfgVal)
		assert.NoError(t, err)
		retMconfig.ConfigsByKey[cfgKey] = cfgAny
	}
	return retMconfig
}
