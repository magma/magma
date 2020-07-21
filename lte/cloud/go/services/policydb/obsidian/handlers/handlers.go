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
	orc8rhandlers "magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
)

const (
	policiesRootPath         = orc8rhandlers.ManageNetworkPath + obsidian.UrlSep + "policies"
	policyRuleRootPath       = policiesRootPath + obsidian.UrlSep + "rules"
	policyRuleManagePath     = policyRuleRootPath + obsidian.UrlSep + ":rule_id"
	policyBaseNameRootPath   = policiesRootPath + obsidian.UrlSep + "base_names"
	policyBaseNameManagePath = policyBaseNameRootPath + obsidian.UrlSep + ":base_name"

	ratingGroupsRootPath   = orc8rhandlers.ManageNetworkPath + obsidian.UrlSep + "rating_groups"
	ratingGroupsManagePath = ratingGroupsRootPath + obsidian.UrlSep + ":rating_group_id"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: ListBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: CreateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: GetBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteBaseName},

		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: ListRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: CreateRule},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: GetRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRule},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRule},

		{Path: ratingGroupsRootPath, Methods: obsidian.GET, HandlerFunc: ListRatingGroups},
		{Path: ratingGroupsRootPath, Methods: obsidian.POST, HandlerFunc: CreateRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.GET, HandlerFunc: GetRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRatingGroup},
	}
	return ret
}
