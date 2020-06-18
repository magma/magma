/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integration

import (
	"magma/feg/cloud/go/protos"

	"github.com/go-openapi/swag"
)

func getPassAllRuleDefinition(ruleID, monitoringKey string, ratingGroup *uint32, precedence uint32) *protos.RuleDefinition {
	rule := &protos.RuleDefinition{
		RuleName:         ruleID,
		Precedence:       precedence,
		FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
		MonitoringKey:    monitoringKey,
	}
	if ratingGroup != nil {
		rule.RatingGroup = swag.Uint32Value(ratingGroup)
	}
	return rule
}
