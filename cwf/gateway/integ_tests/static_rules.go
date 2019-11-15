/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package integ_tests

import (
	lteProtos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb/obsidian/models"

	"github.com/go-openapi/swag"
)

func getStaticPassAll(ruleID string, monitoringKey string) (*lteProtos.PolicyRule, error) {
	rule := &models.PolicyRule{
		FlowList: []*models.FlowDescription{
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: "UPLINK",
					IPProto:   swag.String("IPPROTO_IP"),
					IPV4Dst:   "0.0.0.0/0",
					IPV4Src:   "0.0.0.0/0",
				},
			},
			{
				Action: swag.String("PERMIT"),
				Match: &models.FlowMatch{
					Direction: "DOWNLINK",
					IPProto:   swag.String("IPPROTO_IP"),
					IPV4Dst:   "0.0.0.0/0",
					IPV4Src:   "0.0.0.0/0",
				},
			},
		},
		ID:            ruleID,
		MonitoringKey: swag.String(monitoringKey),
		Priority:      swag.Uint32(3),
		TrackingType:  models.PolicyRuleTrackingTypeONLYPCRF,
	}
	protoRule := &lteProtos.PolicyRule{}
	err := rule.ToProto(protoRule)
	if err != nil {
		return nil, err
	}
	return protoRule, nil
}
