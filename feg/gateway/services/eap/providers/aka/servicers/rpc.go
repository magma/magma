/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package servce implements EAP-AKA GRPC service
package servicers

import (
	"io"
	"sync"
	"time"

	"magma/feg/gateway/services/eap/client"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
)

type UserCtx struct {
	mu        sync.Mutex
	imsi      aka.IMSI
	state     aka.AkaState
	stateTime time.Time
	Rand,
	Mac,
	Xres []byte
}

type EapAkaSrv struct {
	rwl   sync.RWMutex
	users map[aka.IMSI]*UserCtx
}

// NewEapAkaService creates new Aka Service 'object'
func NewEapAkaService() (*EapAkaSrv, error) {
	return &EapAkaSrv{}, nil
}

// Handle implements AKA handler RPC
func (s *EapAkaSrv) Handle(ctx context.Context, req *protos.Eap) (*protos.Eap, error) {
	p := eap.Packet(req.GetPayload())
	if p == nil {
		return aka.EapErrorRes(0, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil Request")
	}
	err := p.Validate()
	if err != nil {
		identifier := byte(0)
		if err != io.ErrShortBuffer {
			identifier = p.Identifier()
		}
		return aka.EapErrorRes(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}
	identifier := p.Identifier()
	method := p.Type()
	if method == client.EapMethodIdentity {
		return &protos.Eap{Payload: aka.NewIdentityReq(identifier, aka.AT_PERMANENT_ID_REQ)}, nil
	}
	if method != aka.TYPE {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.Unimplemented, "Wrong EAP Method: %d", method)
	}
	if len(p) < aka.MIN_PACKET_LEN {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "EAP-AKA Packet is too short: %d", len(p))
	}
	h := GetHandler(aka.Subtype(p[eap.EapSubtype]))
	if h == nil {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.NotFound, "Unsuported Subtype: %d", p[eap.EapSubtype])
	}
	eapCtx := req.GetCtx()
	rp, err := h(s, eapCtx, p)
	return &protos.Eap{Payload: rp, Ctx: eapCtx}, err
}

// GetLockedUserCtx finds, locks & returns the CTX associated with given IMSI, creates the new state if needed
func (s *EapAkaSrv) GetLockedUserCtx(imsi aka.IMSI) *UserCtx {
	var res *UserCtx
	s.rwl.RLock()
	if res, ok := s.users[imsi]; ok {
		res.mu.Lock()
		s.rwl.RUnlock()
		return res
	}
	s.rwl.RUnlock()
	s.rwl.Lock()
	// check again after locking
	if res, ok := s.users[imsi]; ok {
		res.mu.Lock()
		s.rwl.Unlock()
		return res
	}
	res = &UserCtx{imsi: imsi, state: aka.StateCreated, stateTime: time.Now()}
	res.mu.Lock()
	if s.users == nil {
		s.users = map[aka.IMSI]*UserCtx{}
	}
	s.users[imsi] = res
	s.rwl.Unlock()
	return res
}

// Unlock - unlocks the CTX
func (ctx *UserCtx) Unlock() {
	ctx.mu.Unlock()
}

// DeleteUserCtx deletes unlocked CTX
func (s *EapAkaSrv) DeleteUserCtx(ctx *UserCtx) bool {
	key := ctx.imsi
	s.rwl.Lock()
	_, ok := s.users[key]
	if ok {
		delete(s.users, key)
	}
	s.rwl.Unlock()
	return ok
}

// State returns current CTX state (CTX must be locked)
func (ctx *UserCtx) State() (aka.AkaState, time.Time) {
	return ctx.state, ctx.stateTime
}

// SetState updates current CTX state (CTX must be locked)
func (ctx *UserCtx) SetState(s aka.AkaState) {
	ctx.state, ctx.stateTime = s, time.Now()
}
