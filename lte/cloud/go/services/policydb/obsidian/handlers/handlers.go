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
		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: listBaseNameHandler, MigratedHandlerFunc: listBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: createBaseNameHandler, MigratedHandlerFunc: createBaseName, MultiplexAfterMigration: true},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: getBaseNameHandler, MigratedHandlerFunc: getBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: updateBaseNameHandler, MigratedHandlerFunc: updateBaseName, MultiplexAfterMigration: true},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteBaseNameHandler, MigratedHandlerFunc: deleteBaseName, MultiplexAfterMigration: true},

		// rules
		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: listRulesHandler, MigratedHandlerFunc: listRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: createRuleHandler, MigratedHandlerFunc: createRule, MultiplexAfterMigration: true},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: getRuleHandler, MigratedHandlerFunc: getRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: updateRuleHandler, MigratedHandlerFunc: updateRule, MultiplexAfterMigration: true},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteRuleHandler, MigratedHandlerFunc: deleteRule, MultiplexAfterMigration: true},
	}
}
