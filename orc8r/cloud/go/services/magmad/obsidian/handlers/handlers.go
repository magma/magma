/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"magma/orc8r/cloud/go/obsidian"
)

// GetObsidianHandlers returns all obsidian handlers for magmad
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// V1
		{Path: RebootGatewayV1, Methods: obsidian.POST, HandlerFunc: rebootGateway},
		{Path: RestartServicesV1, Methods: obsidian.POST, HandlerFunc: restartServices},
		{Path: GatewayPingV1, Methods: obsidian.POST, HandlerFunc: gatewayPing},
		{Path: GatewayGenericCommandV1, Methods: obsidian.POST, HandlerFunc: gatewayGenericCommand},
		{Path: TailGatewayLogsV1, Methods: obsidian.POST, HandlerFunc: tailGatewayLogs},
	}
}
