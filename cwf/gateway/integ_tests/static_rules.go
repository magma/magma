/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	"magma/lte/cloud/go/plugin/models"
	"magma/lte/cloud/go/protos"

	"github.com/go-openapi/swag"
)

func getStaticPassAll(ruleID string, monitoringKey string, trackingType string) *protos.PolicyRule {
	rule := &models.PolicyRuleConfig{
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPV4Dst:   "0.0.0.0/0",
					IPV4Src:   "0.0.0.0/0",
				},
			},
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: swag.String("DOWNLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPV4Dst:   "0.0.0.0/0",
					IPV4Src:   "0.0.0.0/0",
				},
			},
		},
		MonitoringKey: monitoringKey,
		Priority:      swag.Uint32(3),
		TrackingType:  trackingType,
	}

	return rule.ToProto(ruleID)
}
