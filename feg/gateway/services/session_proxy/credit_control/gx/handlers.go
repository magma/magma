/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gx

import (
	"context"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/glog"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/feg/gateway/services/session_proxy/relay"
	"magma/gateway/service_registry"
)

// ccaHandler parses a CCADiameterMessage received over Gx and returns the
// `KeyAndAnswer` packed inside the CCA message.
func ccaHandler(message *diam.Message) diameter.KeyAndAnswer {
	var cca CCADiameterMessage
	glog.V(2).Infof("Received Gx CCA message:\n%s\n", message)
	if err := message.Unmarshal(&cca); err != nil {
		metrics.GxUnparseableMsg.Inc()
		glog.Errorf("Received unparseable CCA over Gx: %s", err)
		return diameter.KeyAndAnswer{}
	}
	sid := diameter.DecodeSessionID(cca.SessionID)
	return diameter.KeyAndAnswer{
		Key: credit_control.GetRequestKey(credit_control.Gx, sid, cca.RequestNumber),
		Answer: &CreditControlAnswer{
			ResultCode:             cca.ResultCode,
			ExperimentalResultCode: cca.ExperimentalResult.ExperimentalResultCode,
			SessionID:              sid,
			OriginHost:             cca.OriginHost,
			RequestNumber:          cca.RequestNumber,
			RuleInstallAVP:         cca.RuleInstalls,
			RuleRemoveAVP:          cca.RuleRemovals,
			UsageMonitors:          cca.UsageMonitors[:],
			EventTriggers:          cca.EventTriggers,
			RevalidationTime:       cca.RevalidationTime,
		},
	}
}

type ReAuthHandler func(request *ReAuthRequest) *ReAuthAnswer

// Factory function for a RAR message handler which relays to the corresponding
// gateway.
func GetGxReAuthHandler(cloudRegistry service_registry.GatewayRegistry, policyDBClient policydb.PolicyDBClient) ReAuthHandler {
	return func(request *ReAuthRequest) *ReAuthAnswer {
		sid := diameter.DecodeSessionID(request.SessionID)
		imsi, err := relay.GetIMSIFromSessionID(sid)
		if err != nil {
			glog.Errorf("Error retrieving IMSI from session ID %s: %s", request.SessionID, err)
			return &ReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnknownSessionID,
			}
		}

		client, err := relay.GetSessionProxyResponderClient(cloudRegistry)
		if err != nil {
			glog.Error(err)
			return &ReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnableToDeliver,
			}
		}
		defer client.Close()

		gwReq := request.ToProto(imsi, sid, policyDBClient)
		ans, err := client.PolicyReAuth(context.Background(), gwReq)
		if err != nil {
			glog.Errorf("Error relaying Gx reauth request to gateway: %s", err)
			return &ReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnableToDeliver,
			}
		}
		return (&ReAuthAnswer{}).FromProto(request.SessionID, ans)
	}
}
