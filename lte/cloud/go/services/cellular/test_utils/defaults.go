/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_utils

import (
	"magma/lte/cloud/go/services/cellular/obsidian/models"
)

func NewDefaultTDDNetworkConfig() *models.NetworkCellularConfigs {
	return &models.NetworkCellularConfigs{
		Ran: &models.NetworkRanConfigs{
			BandwidthMhz:           20,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
			TddConfig: &models.NetworkRanConfigsTddConfig{
				Earfcndl:               44590,
				SubframeAssignment:     2,
				SpecialSubframePattern: 7,
			},
		},
		Epc: &models.NetworkEpcConfigs{
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

func NewDefaultFDDNetworkConfig() *models.NetworkCellularConfigs {
	return &models.NetworkCellularConfigs{
		Ran: &models.NetworkRanConfigs{
			BandwidthMhz: 20,
			Earfcndl:     1,
			FddConfig: &models.NetworkRanConfigsFddConfig{
				Earfcndl: 1,
				Earfcnul: 18001,
			},
		},
		Epc: &models.NetworkEpcConfigs{
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

func OldTDDNetworkConfig() *models.NetworkCellularConfigs {
	return &models.NetworkCellularConfigs{
		Ran: &models.NetworkRanConfigs{
			BandwidthMhz:           20,
			Earfcndl:               44590,
			SubframeAssignment:     2,
			SpecialSubframePattern: 7,
		},
		Epc: &models.NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),
		},
	}
}

func OldFDDNetworkConfig() *models.NetworkCellularConfigs {
	return &models.NetworkCellularConfigs{
		Ran: &models.NetworkRanConfigs{
			BandwidthMhz: 20,
			Earfcndl:     1,
		},
		Epc: &models.NetworkEpcConfigs{
			Mcc: "001",
			Mnc: "01",
			Tac: 1,
			// 16 bytes of \x11
			LteAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
			LteAuthAmf: []byte("\x80\x00"),
		},
	}
}

func NewDefaultGatewayConfig() *models.GatewayCellularConfigs {
	return &models.GatewayCellularConfigs{
		AttachedEnodebSerials: []string{"enb1"},
		Ran: &models.GatewayRanConfigs{
			Pci:             260,
			TransmitEnabled: true,
		},
		Epc: &models.GatewayEpcConfigs{
			NatEnabled: true,
			IPBlock:    "192.168.128.0/24",
		},
		NonEpsService: &models.GatewayNonEpsServiceConfigs{
			CsfbMcc:              "",
			CsfbMnc:              "",
			Lac:                  1,
			CsfbRat:              0, //2G
			Arfcn2g:              []uint32{},
			NonEpsServiceControl: 0, //CONTROL_OFF
		},
	}
}

func NewDefaultEnodebConfig() *models.NetworkEnodebConfigs {
	return &models.NetworkEnodebConfigs{
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
