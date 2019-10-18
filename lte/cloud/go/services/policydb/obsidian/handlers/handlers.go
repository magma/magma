/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"magma/lte/cloud/go/plugin/handlers"
	"magma/orc8r/cloud/go/obsidian"
)

const (
	policiesRootPath         = obsidian.RestRoot + "/networks/:network_id/policies"
	policyRuleRootPath       = policiesRootPath + "/rules"
	policyRuleManagePath     = policyRuleRootPath + "/:rule_id"
	policyBaseNameRootPath   = policiesRootPath + "/base_names"
	policyBaseNameManagePath = policyBaseNameRootPath + "/:base_name"
)

// GetObsidianHandlers returns all obsidian handlers for policydb
func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		// base names
		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: handlers.ListBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: handlers.CreateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: handlers.GetBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: handlers.UpdateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: handlers.DeleteBaseName},

		// rules
		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: handlers.ListRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: handlers.CreateRule},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: handlers.GetRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: handlers.UpdateRule},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: handlers.DeleteRule},
	}
}
