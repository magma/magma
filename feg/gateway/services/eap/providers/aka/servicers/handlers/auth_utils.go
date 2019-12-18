/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package handlers

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	swx_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	"magma/feg/gateway/services/swx_proxy"
)

var (
	challengeReqTemplate eap.Packet
	challengeReqTemplateLen,
	// Offsets in the challenge template of corresponding attribute values
	atRandOffset,
	atAutnOffset,
	atMacOffset int
)

func init() {
	var err error
	p := eap.NewPacket(eap.RequestCode, 0, []byte{aka.TYPE, byte(aka.SubtypeChallenge), 0, 0})
	atRandOffset = len(p) + aka.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		aka.AT_RAND, append(
			[]byte{0, 0}, // reserved
			make([]byte, aka.RAND_LEN)...)))
	if err != nil {
		panic(err)
	}
	atAutnOffset = len(p) + aka.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		aka.AT_AUTN, append(
			[]byte{0, 0}, // reserved
			make([]byte, aka.AUTN_LEN)...)))
	if err != nil {
		panic(err)
	}
	atMacOffset = len(p) + aka.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		aka.AT_MAC, append(
			[]byte{0, 0}, // reserved
			make([]byte, aka.MAC_LEN)...)))
	if err != nil {
		panic(err)
	}
	challengeReqTemplateLen = len(p)
	challengeReqTemplate = p
}

func createChallengeRequest(
	s *servicers.EapAkaSrv,
	lockedCtx *servicers.UserCtx,
	identifier uint8,
	resyncInfo []byte) (eap.Packet, error) {

	metrics.SwxRequests.Inc()
	swxStartTime := time.Now()

	ans, err := swx_proxy.Authenticate(
		&swx_protos.AuthenticationRequest{
			UserName:             string(lockedCtx.Imsi),
			SipNumAuthVectors:    1,
			AuthenticationScheme: swx_protos.AuthenticationScheme_EAP_AKA,
			ResyncInfo:           resyncInfo,
			RetrieveUserProfile:  true,
		})

	metrics.SWxLatency.Observe(time.Since(swxStartTime).Seconds())

	if err != nil {
		metrics.SwxFailures.Inc()
		errCode := codes.Internal
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			errCode = se.GRPCStatus().Code()
		}
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, errCode, err.Error())
	}
	if ans == nil {
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.Internal, "Error: Nil SWx Response")
	}
	if len(ans.SipAuthVectors) == 0 {
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.Internal, "Error: Missing/empty SWx Auth Vector: %+v", *ans)
	}
	av := ans.SipAuthVectors[0] // Use first vector for now
	ra := av.GetRandAutn()
	if len(ra) < aka.RandAutnLen {
		return aka.EapErrorResPacket(
			identifier,
			aka.NOTIFICATION_FAILURE,
			codes.Internal,
			"Invalid SWx RandAutn len (%d, expected: %d) in Response: %+v",
			len(ra), aka.RandAutnLen, *ans)
	}

	identifier++

	lockedCtx.Identifier = identifier
	lockedCtx.Rand = ra[:aka.RAND_LEN]
	autn := ra[aka.RAND_LEN:aka.RandAutnLen]
	lockedCtx.Xres = av.GetXres()
	lockedCtx.Profile = ans.GetUserProfile()
	lockedCtx.AuthSessionId = ans.GetSessionId()

	// Clone EAP Challenge packet
	p := eap.Packet(make([]byte, challengeReqTemplateLen))
	copy(p, challengeReqTemplate)

	// Set current identifier
	p[eap.EapMsgIdentifier] = identifier

	// Set AT_RAND
	copy(p[atRandOffset:], lockedCtx.Rand)

	// Set AT_AUTN
	copy(p[atAutnOffset:], autn)

	// Calculate AT_MAC
	IK := av.GetIntegrityKey()
	CK := av.GetConfidentialityKey()
	_, lockedCtx.K_aut, lockedCtx.MSK, _ = aka.MakeAKAKeys([]byte(lockedCtx.Identity), IK, CK)
	mac := aka.GenMac(p, lockedCtx.K_aut)
	// Set AT_MAC
	copy(p[atMacOffset:], mac)
	return p, nil
}
