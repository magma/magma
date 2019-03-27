/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers provided AKA Response handlers for supported AKA subtypes
package handlers

import (
	"io"
	"log"
	"reflect"

	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
)

func init() {
	servicers.AddHandler(aka.SubtypeChallenge, challengeResponse)
}

// challengeResponse implements handler for AKA Challenge Response,
// see https://tools.ietf.org/html/rfc4187#page-49 for details
func challengeResponse(s *servicers.EapAkaSrv, ctx *protos.EapContext, req eap.Packet) (eap.Packet, error) {
	identifier := req.Identifier()
	if ctx == nil {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing Session ID")
	}
	imsi, uc, ok := s.FindSession(ctx.SessionId)
	if !ok {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No Session found for ID: %s", ctx.SessionId)
	}
	if uc == nil {
		s.UpdateSessionTimeout(ctx.SessionId, aka.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No IMSI '%s' found for SessionID: %s", imsi, ctx.SessionId)
	}

	state, _ := uc.State()
	if state != aka.StateChallenge {
		log.Printf(
			"AKA Challenge Response: Unexpected user state: %d for IMSI: %s, Session: %s",
			state, imsi, ctx.SessionId)
	}

	p := make([]byte, len(req))
	copy(p, req)
	scanner, err := eap.NewAttributeScanner(p)
	if err != nil {
		s.UpdateSessionUnlockCtx(uc, aka.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}

	var a, atMac, atRes eap.Attribute

attrLoop:
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		case aka.AT_MAC:
			atMac = a
			if atRes != nil {
				break attrLoop
			}
		case aka.AT_RES:
			atRes = a
			if atMac != nil {
				break attrLoop
			}
		case aka.AT_CHECKCODE: // Ignore CHECKCODE for now
		default:
			log.Printf("INFO: Unexpected EAP-AKA Challenge Response Attribute type %d", a.Type())
		}
	}

	if err != nil {
		s.UpdateSessionUnlockCtx(uc, aka.NotificationTimeout())
		if err == io.EOF {
			return aka.EapErrorResPacket(
				identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing AT_MAC | AT_RES")
		}
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}

	// Verify MAC
	macBytes := atMac.Marshaled()
	if len(macBytes) < aka.ATT_HDR_LEN+aka.MAC_LEN {
		s.UpdateSessionUnlockCtx(uc, aka.NotificationTimeout())
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Malformed AT_MAC")
	}
	ueMac := make([]byte, len(macBytes)-aka.ATT_HDR_LEN)
	copy(ueMac, macBytes[aka.ATT_HDR_LEN:])

	for i := aka.ATT_HDR_LEN; i < len(macBytes); i++ {
		macBytes[i] = 0
	}
	mac := aka.GenMac(p, uc.K_aut)
	if !reflect.DeepEqual(ueMac, mac) {
		s.UpdateSessionUnlockCtx(uc, aka.NotificationTimeout())
		log.Printf(
			"Invalid MAC for Session ID: %s; IMSI: %s; UE MAC: %x; Expected MAC: %x; EAP: %x",
			ctx.SessionId, imsi, ueMac, mac, req)
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.Unauthenticated,
			"Invalid MAC for Session ID: %s; IMSI: %s", ctx.SessionId, imsi)
	}

	// Verify AT_RES
	ueRes := atRes.Marshaled()[aka.ATT_HDR_LEN:]
	if !reflect.DeepEqual(ueRes, uc.Xres) {
		log.Printf("Invalid AT_RES for Session ID: %s; IMSI: %s\n\t%.3v !=\n\t%.3v",
			ctx.SessionId, imsi, ueRes, uc.Xres)
		s.UpdateSessionUnlockCtx(uc, aka.NotificationTimeout())
		return aka.EapErrorResPacketWithMac(
			identifier, aka.NOTIFICATION_FAILURE_AUTH, uc.K_aut, codes.Unauthenticated,
			"Invalid AT_RES for Session ID: %s; IMSI: %s", ctx.SessionId, imsi)
	}

	// All good, set IMSI & MSK & return SuccessCode
	ctx.Imsi = string(imsi)
	ctx.Msk = uc.MSK
	// For now - cleanup after successful auth, TBD: keep CTX long term...
	uc.Unlock()
	s.RemoveSession(ctx.SessionId)

	return []byte{eap.SuccessCode, identifier, 0, 4}, nil
}
