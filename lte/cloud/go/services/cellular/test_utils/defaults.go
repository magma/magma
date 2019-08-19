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
