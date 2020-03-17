/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import "magma/feg/cloud/go/protos"

func getDynamicPassAll(ruleID, monitoringKey string, precedence uint32) *protos.RuleDefinition {
	return &protos.RuleDefinition{
		RuleName:         ruleID,
		Precedence:       precedence,
		FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
		MonitoringKey:    monitoringKey,
	}
}
