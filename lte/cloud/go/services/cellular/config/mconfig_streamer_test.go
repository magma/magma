/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	"magma/lte/cloud/go/protos/mconfig"
	cellular_config "magma/lte/cloud/go/services/cellular/config"
	"magma/lte/cloud/go/services/cellular/test_utils"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/streaming"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"
	dnsd_config "magma/orc8r/cloud/go/services/dnsd/config"
	dnsd_protos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

// Long test case just to walk through a typical call flow
func TestCellularStreamer_ApplyMconfigUpdate(t *testing.T) {
	cs := &cellular_config.CellularStreamer{}

	// Create NW config, verify fields we're expecting
	inputMconfigs := map[string]*protos.GatewayConfigs{
		"gw1": {ConfigsByKey: map[string]*any.Any{}},
		"gw2": {ConfigsByKey: map[string]*any.Any{}},
	}
	nwUpdate := &streaming.ConfigUpdate{
		ConfigType: cellular_config.CellularNetworkType,
		ConfigKey:  "nw",
		NewValue:   test_utils.NewDefaultTDDNetworkConfig(),
		Operation:  streaming.ReadOperation,
	}

	_, err := cs.ApplyMconfigUpdate(nwUpdate, inputMconfigs)
	assert.NoError(t, err)

	expected := map[string]proto.Message{
		"enodebd": &mconfig.EnodebD{
			LogLevel:               protos.LogLevel_INFO,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			TddConfig: &mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz: 20,
			Tac:          1,
			PlmnidList:   "00101",
		},
		"mobilityd": &mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
		},
		"mme": &mconfig.MME{
			LogLevel:     protos.LogLevel_INFO,
			Mcc:          "001",
			Mnc:          "01",
			Tac:          1,
			MmeCode:      1,
			MmeGid:       1,
			RelayEnabled: false,
		},
		"pipelined": &mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			DefaultRuleId: "",
			Services: []mconfig.PipelineD_NetworkServices{
				mconfig.PipelineD_METERING,
				mconfig.PipelineD_DPI,
				mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			SubProfiles:  map[string]*mconfig.SubscriberDB_SubscriptionProfile{},
			RelayEnabled: false,
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &mconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: false,
		},
	}
	expectedMconfig := getExpectedMconfig(t, expected)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig},
		inputMconfigs,
	)

	// Create GW config, apply only to gw2, validate fields we're expecting
	gwUpdate := &streaming.ConfigUpdate{
		ConfigType: cellular_config.CellularGatewayType,
		ConfigKey:  "gw2",
		NewValue:   test_utils.NewDefaultGatewayConfig(),
		Operation:  streaming.UpdateOperation,
	}
	onlyGw2Input := map[string]*protos.GatewayConfigs{"gw2": inputMconfigs["gw2"]}
	_, err = cs.ApplyMconfigUpdate(gwUpdate, onlyGw2Input)
	assert.NoError(t, err)

	expectedGw2 := map[string]proto.Message{
		"enodebd": &mconfig.EnodebD{
			LogLevel:               protos.LogLevel_INFO,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			Pci:                    260,
			TddConfig: &mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz:        20,
			AllowEnodebTransmit: true,
			Tac:                 1,
			PlmnidList:          "00101",
			CsfbRat:             mconfig.EnodebD_CSFBRAT_2G,
			Arfcn_2G:            []int32{},
		},
		"mobilityd": &mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &mconfig.MME{
			LogLevel:             protos.LogLevel_INFO,
			Mcc:                  "001",
			Mnc:                  "01",
			Tac:                  1,
			MmeCode:              1,
			MmeGid:               1,
			NonEpsServiceControl: mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:              "",
			CsfbMnc:              "",
			Lac:                  1,
			RelayEnabled:         false,
		},
		"pipelined": &mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []mconfig.PipelineD_NetworkServices{
				mconfig.PipelineD_METERING,
				mconfig.PipelineD_DPI,
				mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			SubProfiles:  map[string]*mconfig.SubscriberDB_SubscriptionProfile{},
			RelayEnabled: false,
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &mconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: false,
		},
	}
	expectedGw2Mconfig := getExpectedMconfig(t, expectedGw2)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedGw2Mconfig},
		inputMconfigs,
	)

	// Create dnsd config, validate fields
	dnsdUpdate := &streaming.ConfigUpdate{
		ConfigType: dnsd_config.DnsdNetworkType,
		ConfigKey:  "nw",
		NewValue:   &dnsd_protos.NetworkDNSConfig{EnableCaching: true},
		Operation:  streaming.CreateOperation,
	}
	_, err = cs.ApplyMconfigUpdate(dnsdUpdate, inputMconfigs)
	assert.NoError(t, err)

	expected["mme"] = &mconfig.MME{
		LogLevel:         protos.LogLevel_INFO,
		Mcc:              "001",
		Mnc:              "01",
		Tac:              1,
		MmeCode:          1,
		MmeGid:           1,
		EnableDnsCaching: true,
		RelayEnabled:     false,
	}
	expectedGw2["mme"] = &mconfig.MME{
		LogLevel:             protos.LogLevel_INFO,
		Mcc:                  "001",
		Mnc:                  "01",
		Tac:                  1,
		MmeCode:              1,
		MmeGid:               1,
		NonEpsServiceControl: mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
		CsfbMcc:              "",
		CsfbMnc:              "",
		Lac:                  1,
		EnableDnsCaching:     true,
		RelayEnabled:         false,
	}

	expectedMconfig = getExpectedMconfig(t, expected)
	expectedGw2Mconfig = getExpectedMconfig(t, expectedGw2)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedGw2Mconfig},
		inputMconfigs,
	)

	// Delete dnsd config, verify that we only falsify EnableDnsCaching
	dnsdUpdate = &streaming.ConfigUpdate{
		ConfigType: dnsd_config.DnsdNetworkType,
		ConfigKey:  "nw",
		NewValue:   nil,
		Operation:  streaming.DeleteOperation,
	}
	_, err = cs.ApplyMconfigUpdate(dnsdUpdate, inputMconfigs)
	assert.NoError(t, err)

	expected["mme"] = &mconfig.MME{
		LogLevel:         protos.LogLevel_INFO,
		Mcc:              "001",
		Mnc:              "01",
		Tac:              1,
		MmeCode:          1,
		MmeGid:           1,
		EnableDnsCaching: false,
		RelayEnabled:     false,
	}
	expectedGw2["mme"] = &mconfig.MME{
		LogLevel:             protos.LogLevel_INFO,
		Mcc:                  "001",
		Mnc:                  "01",
		Tac:                  1,
		MmeCode:              1,
		MmeGid:               1,
		NonEpsServiceControl: mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
		CsfbMcc:              "",
		CsfbMnc:              "",
		Lac:                  1,
		EnableDnsCaching:     false,
		RelayEnabled:         false,
	}
	expectedMconfig = getExpectedMconfig(t, expected)
	expectedGw2Mconfig = getExpectedMconfig(t, expectedGw2)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedGw2Mconfig},
		inputMconfigs,
	)

	// Delete gw1 gateway config, we should clear everything only for gw1
	gwUpdate = &streaming.ConfigUpdate{
		ConfigType: cellular_config.CellularGatewayType,
		ConfigKey:  "gw1",
		NewValue:   nil,
		Operation:  streaming.DeleteOperation,
	}
	onlyGw1Input := map[string]*protos.GatewayConfigs{"gw1": inputMconfigs["gw1"]}
	_, err = cs.ApplyMconfigUpdate(gwUpdate, onlyGw1Input)
	assert.NoError(t, err)

	expected = map[string]proto.Message{}
	expectedMconfig = getExpectedMconfig(t, expected)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedGw2Mconfig},
		inputMconfigs,
	)

	// Delete network config, we should clear everything
	nwUpdate = &streaming.ConfigUpdate{
		ConfigType: cellular_config.CellularNetworkType,
		ConfigKey:  "nw",
		NewValue:   nil,
		Operation:  streaming.DeleteOperation,
	}
	_, err = cs.ApplyMconfigUpdate(nwUpdate, inputMconfigs)
	assert.NoError(t, err)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig},
		inputMconfigs,
	)
}

func TestCellularStreamer_SeedNewGatewayMconfig(t *testing.T) {
	config_test_init.StartTestService(t)

	err := config.CreateConfig("network", cellular_config.CellularNetworkType, "network", test_utils.NewDefaultTDDNetworkConfig())
	assert.NoError(t, err)

	s := &cellular_config.CellularStreamer{}
	mconfigOut := &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	s.SeedNewGatewayMconfig("network", "gw", mconfigOut)

	expected := map[string]proto.Message{
		"enodebd": &mconfig.EnodebD{
			LogLevel:               protos.LogLevel_INFO,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			TddConfig: &mconfig.EnodebD_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
			BandwidthMhz: 20,
			Tac:          1,
			PlmnidList:   "00101",
		},
		"mobilityd": &mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
		},
		"mme": &mconfig.MME{
			LogLevel:     protos.LogLevel_INFO,
			Mcc:          "001",
			Mnc:          "01",
			Tac:          1,
			MmeCode:      1,
			MmeGid:       1,
			RelayEnabled: false,
		},
		"pipelined": &mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			DefaultRuleId: "",
			Services: []mconfig.PipelineD_NetworkServices{
				mconfig.PipelineD_METERING,
				mconfig.PipelineD_DPI,
				mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			SubProfiles:  map[string]*mconfig.SubscriberDB_SubscriptionProfile{},
			RelayEnabled: false,
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &mconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: false,
		},
	}
	expectedMconfig := getExpectedMconfig(t, expected)
	assert.Equal(t, expectedMconfig, mconfigOut)
}

func getExpectedMconfig(t *testing.T, expected map[string]proto.Message) *protos.GatewayConfigs {
	ret := &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	for k, v := range expected {
		vAny, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		ret.ConfigsByKey[k] = vAny
	}
	return ret
}
