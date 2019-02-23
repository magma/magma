/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package providers_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config/streaming/providers"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/config/streaming/storage/mocks"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestMconfigStreamProvider_GetUpdates(t *testing.T) {
	// Register 2 gateways
	magmad_test_init.StartTestService(t)
	_, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: "nw"},
		"mconfig_provider_network")
	assert.NoError(t, err)

	_, err = magmad.RegisterGatewayWithId(
		"mconfig_provider_network",
		&magmad_protos.AccessGatewayRecord{Name: "gw1", HwId: &protos.AccessGatewayID{Id: "id1"}},
		"gw1",
	)
	assert.NoError(t, err)
	_, err = magmad.RegisterGatewayWithId(
		"mconfig_provider_network",
		&magmad_protos.AccessGatewayRecord{Name: "gw2", HwId: &protos.AccessGatewayID{Id: "id2"}},
		"gw2",
	)
	assert.NoError(t, err)

	// Only 1 will have an mconfig
	mconfigFixture := map[string]proto.Message{"magmad": &mconfig.MagmaD{TierId: "test_tier"}}
	mockStorage := &mocks.MconfigStorage{}
	mockStorage.On("GetMconfig", "mconfig_provider_network", "gw1").
		Return(
			&storage.StoredMconfig{
				NetworkId: "mconfig_provider_network",
				GatewayId: "gw1",
				Mconfig:   getMconfig(t, mconfigFixture),
				Offset:    42,
			},
			nil,
		)
	mockStorage.On("GetMconfig", "mconfig_provider_network", "gw2").Return(nil, nil)

	prov := providers.NewMconfigStreamProvider(mockStorage)

	// Expecting nothing because offset is equal
	actual, err := prov.GetUpdates("id1", getStreamAny(t, 42))
	assert.NoError(t, err)
	assert.Equal(t, []*protos.DataUpdate{}, actual)

	// Should get the fixture here
	actual, err = prov.GetUpdates("id1", getStreamAny(t, 1))
	expected := &protos.OffsetGatewayConfigs{
		Configs: getMconfig(t, mconfigFixture),
		Offset:  42,
	}
	expectedMarshaled, err := protos.MarshalIntern(expected)
	assert.NoError(t, err)
	assert.Equal(t, []*protos.DataUpdate{{Key: "gw1", Value: expectedMarshaled}}, actual)

	// Expecting nothing because no mconfig has been computed
	actual, err = prov.GetUpdates("id2", getStreamAny(t, 42))
	assert.NoError(t, err)
	assert.Equal(t, []*protos.DataUpdate{}, actual)

	mockStorage.AssertExpectations(t)
	mockStorage.AssertNumberOfCalls(t, "GetMconfig", 3)
}

func getStreamAny(t *testing.T, offset int64) *any.Any {
	val := &protos.MconfigStreamRequest{Offset: offset}
	ret, err := ptypes.MarshalAny(val)
	assert.NoError(t, err)
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
