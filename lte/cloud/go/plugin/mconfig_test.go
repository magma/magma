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
	models2 "magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/orc8r"
	orc8rplugin "magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build(t *testing.T) {
	_ = orc8rplugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	_ = orc8rplugin.RegisterPluginForTests(t, &plugin.LteOrchestratorPlugin{})
	builder := &plugin.Builder{}

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkType: models2.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
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
	rating1 := configurator.NetworkEntity{
		Type: lte.RatingGroupEntityType,
		Key:  "1",
		Config: &models2.RatingGroup{
			ID:        models2.RatingGroupID(uint32(1)),
			LimitType: swag.String("INFINITE_UNMETERED"),
		},
	}
	rating2 := configurator.NetworkEntity{
		Type: lte.RatingGroupEntityType,
		Key:  "2",
		Config: &models2.RatingGroup{
			ID:        models2.RatingGroupID(uint32(2)),
			LimitType: swag.String("INFINITE_METERED"),
		},
	}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{enb, lteGW, gw, rating1, rating2},
		Edges: []configurator.GraphEdge{
			{From: gw.GetTypeAndKey(), To: lteGW.GetTypeAndKey()},
			{From: lteGW.GetTypeAndKey(), To: enb.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{
		"enodebd": &mconfig.EnodebD{
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
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
			Arfcn_2G:            nil,
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
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
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
				mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			SubProfiles:  nil,
			RelayEnabled: false,
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel:                      protos.LogLevel_INFO,
			InfiniteMeteredChargingKeys:   []uint32{uint32(2)},
			InfiniteUnmeteredChargingKeys: []uint32{uint32(1)},
		},
		"sessiond": &mconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: false,
		},
	}

	// Happy path
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// Do break with non-allowed network service
	setEPCNetworkServices([]string{"0xdeadbeef"}, &nw)
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.EqualError(t, err, "unknown network service name 0xdeadbeef")

	// Don't break with deprecated network services
	setEPCNetworkServices([]string{"metering"}, &nw)
	expected["pipelined"] = &mconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     "192.168.128.0/24",
		NatEnabled:    true,
		DefaultRuleId: "",
		Services: []mconfig.PipelineD_NetworkServices{
			mconfig.PipelineD_METERING,
		},
	}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_Build_BaseCase(t *testing.T) {
	builder := &plugin.Builder{}

	// no dnsd config, no enodebs
	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkType: models2.NewDefaultTDDNetworkConfig(),
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
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
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
			Arfcn_2G:            nil,
			EnbConfigsBySerial:  nil,
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
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			RelayEnabled:             false,
			CloudSubscriberdbEnabled: false,
			AttachedEnodebTacs:       nil,
		},
		"pipelined": &mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []mconfig.PipelineD_NetworkServices{
				mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			SubProfiles:  nil,
			RelayEnabled: false,
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel:                      protos.LogLevel_INFO,
			InfiniteMeteredChargingKeys:   nil,
			InfiniteUnmeteredChargingKeys: nil,
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

// minimal configuration of enodeB, inherit rest of props from nw/gw configs
func TestBuilder_BuildInheritedProperties(t *testing.T) {
	builder := &plugin.Builder{}

	nw := configurator.Network{
		ID: "n1",
		Configs: map[string]interface{}{
			lte.CellularNetworkType: models2.NewDefaultTDDNetworkConfig(),
			orc8r.DnsdNetworkType: &models.NetworkDNSConfig{
				EnableCaching: swag.Bool(true),
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
		Config: &models2.EnodebConfiguration{
			CellID:          swag.Uint32(42),
			DeviceClass:     "Baicells ID TDD/FDD",
			TransmitEnabled: swag.Bool(true),
		},
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
			LogLevel: protos.LogLevel_INFO,
			Pci:      260,
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
			Arfcn_2G:            nil,
			EnbConfigsBySerial: map[string]*mconfig.EnodebD_EnodebConfig{
				"enb1": {
					Earfcndl:               44590,
					SubframeAssignment:     2,
					SpecialSubframePattern: 7,
					Pci:                    260,
					TransmitEnabled:        true,
					DeviceClass:            "Baicells ID TDD/FDD",
					BandwidthMhz:           20,
					Tac:                    1,
					CellId:                 42,
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
			CsfbMcc:                  "001",
			CsfbMnc:                  "01",
			Lac:                      1,
			RelayEnabled:             false,
			CloudSubscriberdbEnabled: false,
			EnableDnsCaching:         true,
			AttachedEnodebTacs:       []int32{1},
		},
		"pipelined": &mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     "192.168.128.0/24",
			NatEnabled:    true,
			DefaultRuleId: "",
			Services: []mconfig.PipelineD_NetworkServices{
				mconfig.PipelineD_ENFORCEMENT,
			},
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			SubProfiles:  nil,
			RelayEnabled: false,
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel:                      protos.LogLevel_INFO,
			InfiniteMeteredChargingKeys:   nil,
			InfiniteUnmeteredChargingKeys: nil,
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

func newDefaultGatewayConfig() *models2.GatewayCellularConfigs {
	return &models2.GatewayCellularConfigs{
		Ran: &models2.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: swag.Bool(true),
		},
		Epc: &models2.GatewayEpcConfigs{
			NatEnabled: swag.Bool(true),
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &models2.GatewayNonEpsConfigs{
			CsfbMcc:              "001",
			CsfbMnc:              "01",
			Lac:                  swag.Uint32(1),
			CsfbRat:              swag.Uint32(0),
			Arfcn2g:              nil,
			NonEpsServiceControl: swag.Uint32(0),
		},
	}
}

func newDefaultEnodebConfig() *models2.EnodebConfiguration {
	return &models2.EnodebConfiguration{
		Earfcndl:               39150,
		SubframeAssignment:     2,
		SpecialSubframePattern: 7,
		Pci:                    260,
		CellID:                 swag.Uint32(138777000),
		Tac:                    15000,
		BandwidthMhz:           20,
		TransmitEnabled:        swag.Bool(true),
		DeviceClass:            "Baicells ID TDD/FDD",
	}
}

func setEPCNetworkServices(services []string, nw *configurator.Network) {
	inwConfig := nw.Configs[lte.CellularNetworkType]
	cellularNwConfig := inwConfig.(*models2.NetworkCellularConfigs)
	cellularNwConfig.Epc.NetworkServices = services

	nw.Configs[lte.CellularNetworkType] = cellularNwConfig
}
