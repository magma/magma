/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian"
)

// GetObsidianHandlers returns all obsidian handlers for configurator
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// Network
		{Path: ListNetworks, Methods: obsidian.GET, HandlerFunc: listNetworks},
		{Path: RegisterNetwork, Methods: obsidian.POST, HandlerFunc: registerNetwork},
		{Path: ManageNetwork, Methods: obsidian.GET, HandlerFunc: getNetwork},
		{Path: ManageNetwork, Methods: obsidian.PUT, HandlerFunc: updateNetwork},
		{Path: ManageNetwork, Methods: obsidian.DELETE, HandlerFunc: deleteNetwork},
	}
}
