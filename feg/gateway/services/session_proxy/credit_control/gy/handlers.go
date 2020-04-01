/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gy

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/glog"
	"golang.org/x/net/context"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/relay"
	"magma/gateway/service_registry"
	"magma/lte/cloud/go/protos"
)

// GetGyReAuthHandler returns the default handler for RAR messages by relaying
// them to the gateway, where session proxy will initiate a credit update and respond
// with an RAA
func GetGyReAuthHandler(cloudRegistry service_registry.GatewayRegistry) ChargingReAuthHandler {
	return ChargingReAuthHandler(func(request *ChargingReAuthRequest) *ChargingReAuthAnswer {
		sid := diameter.DecodeSessionID(request.SessionID)
		imsi, err := protos.ParseIMSIfromSessionIdWithPrefix(sid)
		if err != nil {
			glog.Errorf("Error retreiving IMSI from Session ID %s: %s", request.SessionID, err)
			return &ChargingReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnknownSessionID,
			}
		}

		client, err := relay.GetSessionProxyResponderClient(cloudRegistry)
		if err != nil {
			glog.Error(err)
			return &ChargingReAuthAnswer{SessionID: request.SessionID, ResultCode: diam.UnableToDeliver}
		}
		defer client.Close()

		ans, err := client.ChargingReAuth(context.Background(), getGyReAuthRequestProto(request, imsi, sid))
		if err != nil {
			glog.Errorf("Error relaying reauth request to gateway: %s", err)
		}
		return getGyReAuthAnswerDiamMsg(request.SessionID, ans)
	})
}

func getGyReAuthRequestProto(diamReq *ChargingReAuthRequest, imsi, sid string) *protos.ChargingReAuthRequest {
	protoReq := &protos.ChargingReAuthRequest{
		SessionId: sid,
		Sid:       imsi,
	}
	if diamReq.RatingGroup != nil {
		protoReq.ChargingKey = *diamReq.RatingGroup
		protoReq.Type = protos.ChargingReAuthRequest_SINGLE_SERVICE
		if diamReq.ServiceIdentifier != nil {
			protoReq.ServiceIdentifier = &protos.ServiceIdentifier{Value: *diamReq.ServiceIdentifier}
		}
	} else {
		protoReq.ChargingKey = 0
		protoReq.Type = protos.ChargingReAuthRequest_ENTIRE_SESSION
	}
	return protoReq
}

func getGyReAuthAnswerDiamMsg(
	sessionID string,
	protoAns *protos.ChargingReAuthAnswer,
) *ChargingReAuthAnswer {
	var resultCode uint32
	if protoAns.Result == protos.ChargingReAuthAnswer_UPDATE_INITIATED {
		resultCode = diam.LimitedSuccess
	} else if protoAns.Result == protos.ChargingReAuthAnswer_UPDATE_NOT_NEEDED {
		resultCode = diam.Success
	} else if protoAns.Result == protos.ChargingReAuthAnswer_SESSION_NOT_FOUND {
		resultCode = diam.UnknownSessionID
	} else {
		resultCode = diam.UnableToComply
	}
	return &ChargingReAuthAnswer{
		SessionID:  sessionID,
		ResultCode: resultCode,
	}
}
