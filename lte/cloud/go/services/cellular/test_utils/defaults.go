/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test_utils

import (
	"magma/lte/cloud/go/services/cellular/protos"
)

func NewDefaultTDDNetworkConfig() *protos.CellularNetworkConfig {
	return &protos.CellularNetworkConfig{
		Ran: &protos.NetworkRANConfig{
			BandwidthMhz:           20,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			TddConfig: &protos.NetworkRANConfig_TDDConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
		},
		Epc: &protos.NetworkEPCConfig{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			RelayEnabled: false,
		},
	}
}

func NewDefaultFDDNetworkConfig() *protos.CellularNetworkConfig {
	return &protos.CellularNetworkConfig{
		Ran: &protos.NetworkRANConfig{
			BandwidthMhz: 20,
			Earfcndl:     1,
			FddConfig: &protos.NetworkRANConfig_FDDConfig{
				Earfcndl: 1,
				Earfcnul: 18001,
			},
		},
		Epc: &protos.NetworkEPCConfig{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:    []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf:   []byte("\x80\x00"),
			RelayEnabled: false,
		},
	}
}

func OldTDDNetworkConfig() *protos.CellularNetworkConfig {
	return &protos.CellularNetworkConfig{
		Ran: &protos.NetworkRANConfig{
			BandwidthMhz:           20,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
		},
		Epc: &protos.NetworkEPCConfig{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),
		},
	}
}

func OldFDDNetworkConfig() *protos.CellularNetworkConfig {
	return &protos.CellularNetworkConfig{
		Ran: &protos.NetworkRANConfig{
			BandwidthMhz: 20,
			Earfcndl:     1,
		},
		Epc: &protos.NetworkEPCConfig{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),
		},
	}
}

func NewDefaultGatewayConfig() *protos.CellularGatewayConfig {
	return &protos.CellularGatewayConfig{
		AttachedEnodebSerials: []string{"enb1"},
		Ran: &protos.GatewayRANConfig{
			Pci:             260,
			TransmitEnabled: true,
		},
		Epc: &protos.GatewayEPCConfig{
			NatEnabled: true,
			IpBlock:    "192.168.128.0/24",
		},
		NonEpsService: &protos.GatewayNonEPSConfig{
			CsfbMcc:              "",
			CsfbMnc:              "",
			Lac:                  1,
			CsfbRat:              protos.GatewayNonEPSConfig_CSFBRAT_2G,
			Arfcn_2G:             []int32(""),
			NonEpsServiceControl: protos.GatewayNonEPSConfig_NON_EPS_SERVICE_CONTROL_OFF,
		},
	}
}

func NewDefaultEnodebConfig() *protos.CellularEnodebConfig {
	return &protos.CellularEnodebConfig{
		Earfcndl:               39150,
		SubframeAssignment:     2,
		SpecialSubframePattern: 7,
		Pci:                    260,
		CellId:                 138777000,
		Tac:                    15000,
		BandwidthMhz:           20,
		TransmitEnabled:        true,
		DeviceClass:            "Baicells ID TDD/FDD",
	}
}
