/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package handlers provided AKA Response handlers for supported AKA subtypes
package handlers

import (
	"log"

	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
)

func init() {
	servicers.AddHandler(aka.SubtypeSynchronizationFailure, resyncResponse)
}

// resyncResponse implements handler for EAP-Response/AKA-Synchronization-Failure,
// see https://tools.ietf.org/html/rfc4187#section-9.6 for details
func resyncResponse(s *servicers.EapAkaSrv, ctx *protos.EapContext, req eap.Packet) (eap.Packet, error) {
	identifier := req.Identifier()
	if ctx == nil {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing Session ID")
	}
	imsi, ok := s.FindSession(ctx.SessionId)
	if !ok {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No Session found for ID: %s", ctx.SessionId)
	}

	p := make([]byte, len(req))
	copy(p, req)
	scanner, err := eap.NewAttributeScanner(p)
	if err != nil {
		s.RemoveSession(ctx.SessionId)
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}

	uc := s.FindLockedUserCtx(imsi)
	if uc == nil {
		s.RemoveSession(ctx.SessionId)
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No IMSI '%s' found for SessionID: %s", imsi, ctx.SessionId)
	}
	defer uc.Unlock()

	state, t := uc.State()
	if state != aka.StateChallenge {
		log.Printf(
			"AKA-Synchronization-Failure: Overwriting unexpected user state: %d,%s for IMSI: %s",
			state, t, imsi)
	}
	uc.SetState(aka.StateIdentity)

	var a eap.Attribute

	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		if a.Type() == aka.AT_AUTS {
			auts := a.Value()
			if len(auts) < 14 {
				s.RemoveSession(ctx.SessionId)
				s.DeleteUserCtx(uc)
				return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Invalid AT_AUTS LKen: %d", len(auts))
			}
			// Resync Info = RAND | AUTS
			resyncInfo := append(append(make([]byte, 0, len(uc.Rand)+len(auts)), uc.Rand...), auts...)
			p, err := createChallengeRequest(s, uc, identifier, resyncInfo)
			if err == nil {
				// Update state
				uc.SetState(aka.StateChallenge)
				s.UpdateSession(uc, ctx.SessionId, aka.ChallengeTimeout())
			} else {
				s.RemoveSession(ctx.SessionId)
				s.DeleteUserCtx(uc)
			}
			return p, err
		}
	}

	return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing AT_AUTS")
}
