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

func getStaticPassAll(ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32) *protos.PolicyRule {
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
		Priority:      swag.Uint32(priority),
		TrackingType:  trackingType,
		RatingGroup:   ratingGroup,
	}

	return rule.ToProto(ruleID)
}

func getStaticDenyAll(ruleID string, monitoringKey string, ratingGroup uint32, trackingType string, priority uint32) *protos.PolicyRule {
	rule := &models.PolicyRuleConfig{
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("DENY"),
				Match: &models.FlowMatch{
					Direction: swag.String("UPLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPV4Dst:   "0.0.0.0/0",
					IPV4Src:   "0.0.0.0/0",
				},
			},
			{
				Action: swag.String("DENY"),
				Match: &models.FlowMatch{
					Direction: swag.String("DOWNLINK"),
					IPProto:   swag.String("IPPROTO_IP"),
					IPV4Dst:   "0.0.0.0/0",
					IPV4Src:   "0.0.0.0/0",
				},
			},
		},
		MonitoringKey: monitoringKey,
		Priority:      swag.Uint32(priority),
		TrackingType:  trackingType,
		RatingGroup:   ratingGroup,
	}

	return rule.ToProto(ruleID)
}
