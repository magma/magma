/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package servcers implements EAP-AKA GRPC service
package servicers

import (
	"log"
	"sync"
	"time"

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
	rwl   sync.RWMutex // R/W lock synchronizing maps access
	users map[aka.IMSI]*UserCtx

	// Map of UE Sessions to IMSIs
	sessions map[string]*SessionCtx
}

// NewEapAkaService creates new Aka Service 'object'
func NewEapAkaService() (*EapAkaSrv, error) {
	return &EapAkaSrv{users: map[aka.IMSI]*UserCtx{}, sessions: map[string]*SessionCtx{}}, nil
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
func (s *EapAkaSrv) DeleteUserCtx(imsi aka.IMSI) bool {
	s.rwl.Lock()
	_, ok := s.users[imsi]
	if ok {
		delete(s.users, imsi)
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

// InitSession either creates new or updates existing session & user ctx,
// it session ID into the CTX and initializes session map as well as users map
// Returns Locked User Ctx
func (s *EapAkaSrv) InitSession(sessionId string, imsi aka.IMSI) (lockedUserContext *UserCtx) {
	var (
		oldSessionTimer *time.Timer
		oldImsi         aka.IMSI
	)
	// create new session with long session wide timeout
	newSession := &SessionCtx{Imsi: imsi}
	newSession.CleanupTimer = time.AfterFunc(aka.SessionTimeout(), func() {
		sessionTimeoutCleanup(s, sessionId, newSession)
	})
	uc := &UserCtx{Imsi: imsi, state: aka.StateCreated, stateTime: time.Now(), locked: true, SessionId: sessionId}
	uc.mu.Lock()

	s.rwl.Lock()

	if oldSession, ok := s.sessions[sessionId]; ok && oldSession != nil {
		oldSessionTimer, oldImsi, oldSession.CleanupTimer = oldSession.CleanupTimer, oldSession.Imsi, nil
	}
	if len(oldImsi) > 0 && oldImsi != imsi {
		delete(s.users, oldImsi)
	}
	s.users[imsi] = uc // overwrite previous ctx on session init
	s.sessions[sessionId] = newSession

	s.rwl.Unlock()

	if oldSessionTimer != nil {
		oldSessionTimer.Stop()
	}
	return uc
}

// UpdateSessionUnlockCtx sets session ID into the CTX and initializes session map & session timeout
func (s *EapAkaSrv) UpdateSessionUnlockCtx(lockedCtx *UserCtx, timeout time.Duration) {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	var (
		oldSession, newSession *SessionCtx
		exist                  bool
		oldTimer               *time.Timer
	)
	newSession = &SessionCtx{Imsi: lockedCtx.Imsi}
	oldSid := lockedCtx.SessionId
	sessionId := lockedCtx.SessionId
	lockedCtx.Unlock()

	newSession.CleanupTimer = time.AfterFunc(timeout, func() {
		sessionTimeoutCleanup(s, sessionId, newSession)
	})

	s.rwl.Lock()
	oldSession, exist = s.sessions[sessionId]
	s.sessions[sessionId] = newSession
	if len(oldSid) > 0 && oldSid != sessionId {
		delete(s.sessions, oldSid)
	}
	if exist && oldSession != nil && oldSession.CleanupTimer != nil {
		oldTimer, oldSession.CleanupTimer = oldSession.CleanupTimer, nil
	}
	s.rwl.Unlock()

	if oldTimer != nil {
		oldTimer.Stop()
	}
}

// UpdateSessionTimeout finds a session with specified ID, if found - cancels its current timeout
// & schedules the new one. Returns true if the session was found
func (s *EapAkaSrv) UpdateSessionTimeout(sessionId string, timeout time.Duration) bool {
	var (
		newSession *SessionCtx
		exist      bool
		oldTimer   *time.Timer
	)

	s.rwl.Lock()

	oldSession, exist := s.sessions[sessionId]
	if exist {
		if oldSession == nil {
			exist = false
		} else {
			oldTimer, oldSession.CleanupTimer = oldSession.CleanupTimer, nil
			newSession = &SessionCtx{Imsi: oldSession.Imsi}
			s.sessions[sessionId] = newSession
			newSession.CleanupTimer = time.AfterFunc(timeout, func() {
				sessionTimeoutCleanup(s, sessionId, newSession)
			})
		}
	}
	s.rwl.Unlock()

	if oldTimer != nil {
		oldTimer.Stop()
	}
	return exist
}

func sessionTimeoutCleanup(s *EapAkaSrv, sessionId string, mySessionCtx *SessionCtx) {
	if s == nil {
		log.Printf("ERROR: Nil EAP-AKA Server for session ID: %s", sessionId)
		return
	}
	var (
		imsi  aka.IMSI
		state aka.AkaState
		uc    *UserCtx
		ok    bool
	)

	s.rwl.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist {
		if sessionCtx != nil {
			imsi = sessionCtx.Imsi
			if sessionCtx == mySessionCtx {
				if uc, ok = s.users[imsi]; ok {
					delete(s.users, imsi)
				}
				delete(s.sessions, sessionId)
			}
		} else {
			exist = false
		}
	}
	s.rwl.Unlock()

	if uc != nil {
		uc.mu.Lock()
		state = uc.state
		uc.mu.Unlock()
	}
	if exist && state != aka.StateAuthenticated {
		log.Printf("EAP-AKA Session %s timeout for IMSI: %s", sessionId, imsi)
	}
}

// FindSession finds and returns IMSI of a session and a flag indication if the find succeeded
// If found, FindSession tries to stop outstanding session timer
func (s *EapAkaSrv) FindSession(sessionId string) (aka.IMSI, *UserCtx, bool) {
	var (
		imsi      aka.IMSI
		lockedCtx *UserCtx
		ok        bool
		timer     *time.Timer
	)
	s.rwl.RLock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist && sessionCtx != nil {
		imsi, timer, sessionCtx.CleanupTimer = sessionCtx.Imsi, sessionCtx.CleanupTimer, nil
		lockedCtx, ok = s.users[imsi]
	}
	s.rwl.RUnlock()

	if ok && lockedCtx != nil {
		lockedCtx.mu.Lock()
		lockedCtx.SessionId = sessionId // just in case - should always match
		lockedCtx.locked = true
	}

	if timer != nil {
		timer.Stop()
	}
	return imsi, lockedCtx, exist
}

// RemoveSession removes session ID from the session map and attempts to cancel corresponding timer
// It also removes associated with the session user CTX if any
// returns associated with the session IMSI or an empty string
func (s *EapAkaSrv) RemoveSession(sessionId string) aka.IMSI {
	var (
		timer *time.Timer
		imsi  aka.IMSI
	)
	s.rwl.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist {
		delete(s.sessions, sessionId)
		if sessionCtx != nil {
			imsi, timer, sessionCtx.CleanupTimer = sessionCtx.Imsi, sessionCtx.CleanupTimer, nil
			delete(s.users, imsi)
		}
	}
	s.rwl.Unlock()

	if timer != nil {
		timer.Stop()
	}
	return imsi
}

// FindAndRemoveSession finds returns IMSI of a session and a flag indication if the find succeeded
// then it deletes the session ID from the map
func (s *EapAkaSrv) FindAndRemoveSession(sessionId string) (aka.IMSI, bool) {
	var (
		imsi  aka.IMSI
		timer *time.Timer
	)
	s.rwl.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist {
		delete(s.sessions, sessionId)
		if sessionCtx != nil {
			imsi, timer, sessionCtx.CleanupTimer = sessionCtx.Imsi, sessionCtx.CleanupTimer, nil
		}
	}
	s.rwl.Unlock()
	if timer != nil {
		timer.Stop()
	}
	return imsi, exist
}

// ResetSessionTimeout finds a session with specified ID, if found - attempts to cancel its current timeout
// (best effort) & schedules the new one. ResetSessionTimeout does not guarantee that the old timeout cleanup
// won't be executed
func (s *EapAkaSrv) ResetSessionTimeout(sessionId string, newTimeout time.Duration) {
	var oldTimer *time.Timer

	s.rwl.Lock()
	session, exist := s.sessions[sessionId]
	if exist {
		if session == nil {
			oldTimer, session.CleanupTimer = session.CleanupTimer, time.AfterFunc(newTimeout, func() {
				sessionTimeoutCleanup(s, sessionId, session)
			})
		}
	}
	s.rwl.Unlock()

	if oldTimer != nil {
		oldTimer.Stop()
	}
}
