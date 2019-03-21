/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers provided AKA Response handlers for supported AKA subtypes
package handlers

import (
	"fmt"
	"io"
	"log"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	swx_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
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

	servicers.AddHandler(aka.SubtypeIdentity, identityResponse)
}

// identityResponse implements handler for AKA Challenge, see https://tools.ietf.org/html/rfc4187#page-49 for reference
func identityResponse(s *servicers.EapAkaSrv, ctx *protos.EapContext, req eap.Packet) (eap.Packet, error) {
	identifier := req.Identifier()
	if ctx == nil {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		ctx.SessionId = eap.CreateSessionId()
		log.Printf("Missing Session ID for EAP: %x; Generated new SID: %s", req, ctx.SessionId)
	}
	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}
	var a eap.Attribute

	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		// Find first valid AT_IDENTITY attribute to get UE IMSI
		if a.Type() == aka.AT_IDENTITY {
			identity, imsi, err := getIMSIIdentity(a)
			if err == nil {
				if imsi[0] != '0' {
					log.Printf("AKA AT_IDENTITY '%s' (IMSI: %s) is non-permanent type", identity, imsi)
				} else {
					imsi = imsi[1:]
				}
				uc := s.GetLockedUserCtx(imsi)
				defer uc.Unlock()

				state, t := uc.State()
				if state > aka.StateCreated {
					log.Printf(
						"EAP AKA IdentityResponse: Overwriting unexpected user state: %d,%s for IMSI: %s",
						state, t, imsi)
				}
				uc.SetState(aka.StateIdentity)
				ans, err := swx_proxy.Authenticate(
					&swx_protos.AuthenticationRequest{
						UserName:             string(imsi),
						SipNumAuthVectors:    1,
						AuthenticationScheme: swx_protos.AuthenticationScheme_EAP_AKA})

				if err != nil {
					errCode := codes.Internal
					if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
						errCode = se.GRPCStatus().Code()
					}
					return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, errCode, err.Error())
				}
				if ans == nil || len(ans.SipAuthVectors) == 0 {
					return aka.EapErrorResPacket(
						identifier, aka.NOTIFICATION_FAILURE, codes.Internal, "Nil SWx Response")
				}
				if len(ans.SipAuthVectors) == 0 {
					return aka.EapErrorResPacket(
						identifier, aka.NOTIFICATION_FAILURE, codes.Internal, "Missing SWx Auth Vector: %+v", *ans)
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
				identifier += 1

				uc.Identifier = identifier
				uc.Rand = ra[:aka.RAND_LEN]
				autn := ra[aka.RAND_LEN:aka.RandAutnLen]
				uc.Xres = av.GetXres()
				// Clone EAP Challenge packet
				p := eap.Packet(make([]byte, challengeReqTemplateLen))
				copy(p, challengeReqTemplate)

				// Set current identifier
				p[eap.EapMsgIdentifier] = identifier

				// Set AT_RAND
				copy(p[atRandOffset:], uc.Rand)

				// Set AT_AUTN
				copy(p[atAutnOffset:], autn)

				// Calculate AT_MAC
				IK := av.GetIntegrityKey()
				CK := av.GetConfidentialityKey()
				_, uc.K_aut, uc.MSK, _ = aka.MakeAKAKeys([]byte(identity), IK, CK)
				mac := aka.GenMac(p, uc.K_aut)
				// Set AT_MAC
				copy(p[atMacOffset:], mac)
				// Update state
				uc.SetState(aka.StateChallenge)
				s.UpdateSession(uc, ctx.SessionId, aka.ChallengeTimeout())
				return p, nil // success - return EAP packet
			}
		}
	}
	if err != nil && err != io.EOF {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}
	return aka.EapErrorResPacket(
		identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition, "Missing AT_IDENTITY Attribute")
}

// see https://tools.ietf.org/html/rfc4187#section-4.1.1.4
func getIMSIIdentity(a eap.Attribute) (string, aka.IMSI, error) {
	if a.Type() != aka.AT_IDENTITY {
		return "", "", fmt.Errorf("Unexpected Attr Type: %d, AT_IDENTITY expected", a.Type())
	}
	if a.Len() <= 4 {
		return "", "", fmt.Errorf("AT_IDENTITY is too short: %d", a.Len())
	}
	val := a.Value()
	actualLen2 := int(val[0])<<8 + int(val[1]) + 2
	if actualLen2 > len(val) {
		return "", "", fmt.Errorf("Corrupt AT_IDENTITY Attribute: actual len %d > data len %d", actualLen2-2, len(val))
	}
	fullIdentity := string(val[2:actualLen2])
	atIdx := strings.Index(fullIdentity, "@")
	var imsi aka.IMSI
	if atIdx > 0 {
		imsi = aka.IMSI(fullIdentity[:atIdx])
	} else {
		imsi = aka.IMSI(fullIdentity)
	}
	return fullIdentity, imsi, imsi.Validate()
}
