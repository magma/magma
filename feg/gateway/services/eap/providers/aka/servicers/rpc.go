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
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/client"
	"magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers/aka"
)

type UserCtx struct {
	mu         sync.Mutex
	state      aka.AkaState
	stateTime  time.Time
	locked     bool
	Identity   string
	Imsi       aka.IMSI
	Identifier uint8
	Rand,
	K_aut,
	MSK,
	Xres []byte
	SessionId string
}
type SessionCtx struct {
	Imsi         aka.IMSI
	CleanupTimer *time.Timer
}

type EapAkaSrv struct {
	rwl   sync.RWMutex
	users map[aka.IMSI]*UserCtx

	// Map of UE Sessions to IMSIs
	sessionsMu sync.Mutex
	sessions   map[string]*SessionCtx
}

// NewEapAkaService creates new Aka Service 'object'
func NewEapAkaService() (*EapAkaSrv, error) {
	return &EapAkaSrv{users: map[aka.IMSI]*UserCtx{}, sessions: map[string]*SessionCtx{}}, nil
}

// Handle implements AKA handler RPC
func (s *EapAkaSrv) Handle(ctx context.Context, req *protos.Eap) (*protos.Eap, error) {
	p := eap.Packet(req.GetPayload())
	eapCtx := req.GetCtx()
	if p == nil {
		return aka.EapErrorRes(0, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx, "Nil Request")
	}
	err := p.Validate()
	if err != nil {
		identifier := byte(0)
		if err != io.ErrShortBuffer {
			identifier = p.Identifier()
		}
		return aka.EapErrorRes(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx, err.Error())
	}
	identifier := p.Identifier()
	method := p.Type()
	if method == client.EapMethodIdentity {
		return &protos.Eap{Payload: aka.NewIdentityReq(identifier, aka.AT_PERMANENT_ID_REQ)}, nil
	}
	if method != aka.TYPE {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.Unimplemented, eapCtx, "Wrong EAP Method: %d", method)
	}
	if len(p) < aka.MIN_PACKET_LEN {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx,
			"EAP-AKA Packet is too short: %d", len(p))
	}
	h := GetHandler(aka.Subtype(p[eap.EapSubtype]))
	if h == nil {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.NotFound, eapCtx,
			"Unsuported Subtype: %d", p[eap.EapSubtype])
	}
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
		if res.locked {
			panic("Expected unlocked")
		}
		if res.Imsi != imsi {
			panic("IMSI Mismatch")
		}
		res.locked = true
		return res
	}
	s.rwl.RUnlock()
	s.rwl.Lock()
	// check again after locking
	if res, ok := s.users[imsi]; ok {
		res.mu.Lock()
		s.rwl.Unlock()
		if res.locked {
			panic("Expected unlocked")
		}
		res.locked = true
		return res
	}
	res = &UserCtx{Imsi: imsi, state: aka.StateCreated, stateTime: time.Now(), locked: true}
	res.mu.Lock()
	if s.users == nil {
		s.users = map[aka.IMSI]*UserCtx{}
	}
	s.users[imsi] = res
	s.rwl.Unlock()
	return res
}

// FindLockedUserCtx finds, locks & returns the CTX associated with given IMSI
func (s *EapAkaSrv) FindLockedUserCtx(imsi aka.IMSI) *UserCtx {
	s.rwl.RLock()
	defer s.rwl.RUnlock()
	if res, ok := s.users[imsi]; ok {
		res.mu.Lock()
		if res.locked {
			panic("Expected unlocked")
		}
		if res.Imsi != imsi {
			panic("IMSI Mismatch")
		}
		res.locked = true
		return res
	}
	return nil
}

// Unlock - unlocks the CTX
func (lockedCtx *UserCtx) Unlock() {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	lockedCtx.locked = false
	lockedCtx.mu.Unlock()
}

// DeleteUserCtx deletes unlocked CTX
func (s *EapAkaSrv) DeleteUserCtx(ctx *UserCtx) bool {
	key := ctx.Imsi
	s.rwl.Lock()
	_, ok := s.users[key]
	if ok {
		delete(s.users, key)
	}
	s.rwl.Unlock()
	return ok
}

// State returns current CTX state (CTX must be locked)
func (lockedCtx *UserCtx) State() (aka.AkaState, time.Time) {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	return lockedCtx.state, lockedCtx.stateTime
}

// SetState updates current CTX state (CTX must be locked)
func (lockedCtx *UserCtx) SetState(s aka.AkaState) {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	lockedCtx.state, lockedCtx.stateTime = s, time.Now()
}

// UpdateSession sets session ID into the CTX and initializes session map & session timeout
func (s *EapAkaSrv) UpdateSession(lockedCtx *UserCtx, sessionId string, timeout time.Duration) {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	var oldSession, newSession *SessionCtx
	newSession = &SessionCtx{Imsi: lockedCtx.Imsi}
	oldSid := lockedCtx.SessionId
	s.sessionsMu.Lock()
	oldSession, exist := s.sessions[sessionId]
	s.sessions[sessionId] = newSession
	if len(oldSid) > 0 && oldSid != sessionId {
		delete(s.sessions, oldSid)
	}
	lockedCtx.SessionId = sessionId
	if exist && oldSession != nil && oldSession.CleanupTimer != nil {
		oldSession.CleanupTimer.Stop()
	}
	newSession.CleanupTimer = time.AfterFunc(timeout, func() {
		sessionTimeoutCleanup(s, sessionId, newSession)
	})
	s.sessionsMu.Unlock()

}

func sessionTimeoutCleanup(s *EapAkaSrv, sessionId string, mySessionCtx *SessionCtx) {
	if s == nil {
		log.Printf("ERROR: Nil EAP-AKA Server for session ID: %s", sessionId)
		return
	}
	var imsi aka.IMSI
	s.sessionsMu.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist && sessionCtx != nil {
		imsi = sessionCtx.Imsi
		if sessionCtx == mySessionCtx {
			delete(s.sessions, sessionId)
		}
	} else {
		exist = false
	}
	s.sessionsMu.Unlock()

	if exist {
		uc := s.FindLockedUserCtx(imsi)
		if uc != nil {
			if uc.SessionId == sessionId {
				s.DeleteUserCtx(uc)
				uc.Unlock()
				log.Printf("Timed out User Context Removed for IMSI: %s", imsi)
			} else {
				uc.Unlock()
			}
		}
	}
	log.Printf("EAP-AKA Session %s timeout for IMSI: %s", sessionId, imsi)
}

// FindSession finds and returns IMSI of a session and a flag indication if the find succeeded
func (s *EapAkaSrv) FindSession(sessionId string) (aka.IMSI, bool) {
	var (
		imsi  aka.IMSI
		timer *time.Timer
	)
	s.sessionsMu.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist && sessionCtx != nil {
		imsi, timer = sessionCtx.Imsi, sessionCtx.CleanupTimer
	}
	s.sessionsMu.Unlock()

	if timer != nil {
		timer.Stop()
	}
	return imsi, exist
}

// RemoveSession removes session ID from the session map and attempts to cancel corresponding timer
func (s *EapAkaSrv) RemoveSession(sessionId string) {
	var timer *time.Timer
	s.sessionsMu.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist {
		delete(s.sessions, sessionId)
		if sessionCtx != nil {
			timer = sessionCtx.CleanupTimer
		}
	}
	s.sessionsMu.Unlock()

	if timer != nil {
		timer.Stop()
	}
}

// FindAndRemoveSession finds returns IMSI of a session and a flag indication if the find succeeded
// then it deletes the session ID from the map
func (s *EapAkaSrv) FindAndRemoveSession(sessionId string) (aka.IMSI, bool) {
	var (
		imsi  aka.IMSI
		timer *time.Timer
	)
	s.sessionsMu.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist {
		delete(s.sessions, sessionId)
		if sessionCtx != nil {
			imsi, timer = sessionCtx.Imsi, sessionCtx.CleanupTimer
		}
	}
	s.sessionsMu.Unlock()
	if timer != nil {
		timer.Stop()
	}
	return imsi, exist
}
