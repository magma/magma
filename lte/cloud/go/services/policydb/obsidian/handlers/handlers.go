/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import "magma/orc8r/cloud/go/obsidian/handlers"

// GetObsidianHandlers returns all obsidian handlers for policydb
func GetObsidianHandlers() []handlers.Handler {
	return []handlers.Handler{
		// base names
		{Path: policyBaseNameRootPath, Methods: handlers.GET, HandlerFunc: listBaseNameHandler, MigratedHandlerFunc: listBaseNames},
		{Path: policyBaseNameRootPath, Methods: handlers.POST, HandlerFunc: createBaseNameHandler, MigratedHandlerFunc: createBaseName},
		{Path: policyBaseNameManagePath, Methods: handlers.GET, HandlerFunc: getBaseNameHandler, MigratedHandlerFunc: getBaseName},
		{Path: policyBaseNameManagePath, Methods: handlers.PUT, HandlerFunc: updateBaseNameHandler, MigratedHandlerFunc: updateBaseName},
		{Path: policyBaseNameManagePath, Methods: handlers.DELETE, HandlerFunc: deleteBaseNameHandler, MigratedHandlerFunc: deleteBaseName},

		// rules
		{Path: policyRuleRootPath, Methods: handlers.GET, HandlerFunc: listRulesHandler, MigratedHandlerFunc: listRules},
		{Path: policyRuleRootPath, Methods: handlers.POST, HandlerFunc: createRuleHandler, MigratedHandlerFunc: createRule},
		{Path: policyRuleManagePath, Methods: handlers.GET, HandlerFunc: getRuleHandler, MigratedHandlerFunc: getRule},
		{Path: policyRuleManagePath, Methods: handlers.PUT, HandlerFunc: updateRuleHandler, MigratedHandlerFunc: updateRule},
		{Path: policyRuleManagePath, Methods: handlers.DELETE, HandlerFunc: deleteRuleHandler, MigratedHandlerFunc: deleteRule},
	}
}
