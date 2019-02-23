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
		{Path: policyBaseNameRootPath, Methods: handlers.GET, HandlerFunc: listBaseNameHandler},
		{Path: policyBaseNameRootPath, Methods: handlers.POST, HandlerFunc: createBaseNameHandler},
		{Path: policyBaseNameManagePath, Methods: handlers.GET, HandlerFunc: getBaseNameHandler},
		{Path: policyBaseNameManagePath, Methods: handlers.PUT, HandlerFunc: updateBaseNameHandler},
		{Path: policyBaseNameManagePath, Methods: handlers.DELETE, HandlerFunc: deleteBaseNameHandler},

		// rules
		{Path: policyRuleRootPath, Methods: handlers.GET, HandlerFunc: listRulesHandler},
		{Path: policyRuleRootPath, Methods: handlers.POST, HandlerFunc: createRuleHandler},
		{Path: policyRuleManagePath, Methods: handlers.GET, HandlerFunc: getRuleHandler},
		{Path: policyRuleManagePath, Methods: handlers.PUT, HandlerFunc: updateRuleHandler},
		{Path: policyRuleManagePath, Methods: handlers.DELETE, HandlerFunc: deleteRuleHandler},
	}
}
