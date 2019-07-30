/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin_test

import (
	"testing"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/plugin"
	"magma/lte/cloud/go/protos/mconfig"
	cellular_models "magma/lte/cloud/go/services/cellular/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	builder := &plugin.Builder{}

	trueValue := true
	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkType: newDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: &trueValue,
			},
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayType, Key: "gw1",
		Config: newDefaultGatewayConfig(),
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularEnodebType, Key: "enb1"},
		},
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	enb := configurator.NetworkEntity{
		Type: lte.CellularEnodebType, Key: "enb1",
		Config:             newDefaultEnodebConfig(),
		ParentAssociations: []storage.TypeAndKey{lteGW.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
			{From: lteGW.GetTypeAndKey(), To: enb.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{
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
			EnbConfigsBySerial: map[string]*mconfig.EnodebD_EnodebConfig{
				"enb1": {
					Earfcndl:               39150,
					SubframeAssignment:     2,
					SpecialSubframePattern: 7,
					Pci:                    260,
					TransmitEnabled:        true,
					DeviceClass:            "Baicells ID TDD/FDD",
					BandwidthMhz:           20,
					Tac:                    15000,
					CellId:                 138777000,
				},
			},
		},
		"mobilityd": &mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "",
			CsfbMnc:                  "",
			Lac:                      1,
			RelayEnabled:             false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         true,
			AttachedEnodebTacs:       []int32{15000},
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
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_Build_BaseCase(t *testing.T) {
	builder := &plugin.Builder{}

	// no dnsd config, no enodebs
	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkType: newDefaultTDDNetworkConfig(),
		},
	}
	gw := configurator.NetworkEntity{
		Type: orc8r.MagmadGatewayType, Key: "gw1",
		Associations: []storage.TypeAndKey{
			{Type: lte.CellularGatewayType, Key: "gw1"},
		},
	}
	lteGW := configurator.NetworkEntity{
		Type: lte.CellularGatewayType, Key: "gw1",
		Config:             newDefaultGatewayConfig(),
		ParentAssociations: []storage.TypeAndKey{gw.GetTypeAndKey()},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{lteGW, gw},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{
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
			EnbConfigsBySerial:  map[string]*mconfig.EnodebD_EnodebConfig{},
		},
		"mobilityd": &mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  "192.168.128.0/24",
		},
		"mme": &mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      "001",
			Mnc:                      "01",
			Tac:                      1,
			MmeCode:                  1,
			MmeGid:                   1,
			NonEpsServiceControl:     mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			CsfbMcc:                  "",
			CsfbMnc:                  "",
			Lac:                      1,
			RelayEnabled:             false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       []int32{},
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
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func newDefaultTDDNetworkConfig() *cellular_models.NetworkCellularConfigs {
	return &cellular_models.NetworkCellularConfigs{
		Ran: &cellular_models.NetworkRanConfigs{
			BandwidthMhz:           20,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			TddConfig: &cellular_models.NetworkRanConfigsTddConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
		},
		Epc: &cellular_models.NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),
		},
	}
}

func newDefaultGatewayConfig() *cellular_models.GatewayCellularConfigs {
	return &cellular_models.GatewayCellularConfigs{
		AttachedEnodebSerials: []string{"enb1"},
		Ran: &cellular_models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: true,
		},
		Epc: &cellular_models.GatewayEpcConfigs{
			NatEnabled: true,
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &cellular_models.GatewayNonEpsServiceConfigs{
			CsfbMcc:              "",
			CsfbMnc:              "",
			Lac:                  1,
			CsfbRat:              0,
			Arfcn2g:              []uint32{},
			NonEpsServiceControl: 0,
		},
	}
}

func newDefaultEnodebConfig() *cellular_models.NetworkEnodebConfigs {
	return &cellular_models.NetworkEnodebConfigs{
		Earfcndl:               39150,
		SubframeAssignment:     2,
		SpecialSubframePattern: 7,
		Pci:                    260,
		CellID:                 138777000,
		Tac:                    15000,
		BandwidthMhz:           20,
		TransmitEnabled:        true,
		DeviceClass:            "Baicells ID TDD/FDD",
	}
}
