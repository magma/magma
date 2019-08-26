/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import "magma/orc8r/cloud/go/obsidian"

// GetObsidianHandlers returns all obsidian handlers for policydb
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// base names
		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: listBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: createBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: getBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: updateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteBaseName},

		// rules
		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: listRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: createRule},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: getRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: updateRule},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteRule},
	}
}
