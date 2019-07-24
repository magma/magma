/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/orc8r/cloud/go/obsidian/handlers"
)

// GetObsidianHandlers returns all obsidian handlers for configurator
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		// Network
		{Path: ListNetworks, Methods: handlers.GET, HandlerFunc: listNetworks},
		{Path: RegisterNetwork, Methods: handlers.POST, HandlerFunc: registerNetwork},
		{Path: ManageNetwork, Methods: handlers.GET, HandlerFunc: getNetwork},
		{Path: ManageNetwork, Methods: handlers.PUT, HandlerFunc: updateNetwork},
		{Path: ManageNetwork, Methods: handlers.DELETE, HandlerFunc: deleteNetwork},
	}
}
