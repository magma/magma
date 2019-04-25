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
	"sync/atomic"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
)

type UserCtx struct {
	mu         sync.Mutex
	created    time.Time
	state      aka.AkaState
	stateTime  time.Time
	locked     bool
	Identity   string
	Imsi       aka.IMSI
	Profile    *protos.AuthenticationAnswer_UserProfile
	Identifier uint8
	Rand,
	K_aut,
	MSK,
	Xres []byte
	SessionId string
}

type SessionCtx struct {
	*UserCtx
	CleanupTimer *time.Timer
}

type touts struct {
	challengeTimeout,
	errorNotificationTimeout,
	sessionTimeout,
	sessionAuthenticatedTimeout time.Duration
}

type plmnIdVal struct {
	l5 bool
	b6 byte
}

type EapAkaSrv struct {
	rwl sync.RWMutex // R/W lock synchronizing maps access
	// Map of UE Sessions keyed by sessionId
	sessions map[string]*SessionCtx

	// PLMN IDs map, if not empty -> serve only IMSIs with specified PLMN IDs - Read Only
	plmnIds map[string]plmnIdVal

	timeouts touts
}

var defaultTimeouts = touts{
	challengeTimeout:            aka.DefaultChallengeTimeout,
	errorNotificationTimeout:    aka.DefaultErrorNotificationTimeout,
	sessionTimeout:              aka.DefaultSessionTimeout,
	sessionAuthenticatedTimeout: aka.DefaultSessionAuthenticatedTimeout,
}

func (s *EapAkaSrv) ChallengeTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.challengeTimeout)))
}

func (s *EapAkaSrv) SetChallengeTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.challengeTimeout), int64(tout))
}

func (s *EapAkaSrv) NotificationTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.errorNotificationTimeout)))
}

func (s *EapAkaSrv) SetNotificationTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.errorNotificationTimeout), int64(tout))
}

func (s *EapAkaSrv) SessionTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.sessionTimeout)))
}

func (s *EapAkaSrv) SetSessionTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.sessionTimeout), int64(tout))
}

func (s *EapAkaSrv) SessionAuthenticatedTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.sessionAuthenticatedTimeout)))
}

func (s *EapAkaSrv) SetSessionAuthenticatedTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.sessionAuthenticatedTimeout), int64(tout))
}

// NewEapAkaService creates new Aka Service 'object'
func NewEapAkaService(config *mconfig.EapAkaConfig) (*EapAkaSrv, error) {
	service := &EapAkaSrv{
		sessions: map[string]*SessionCtx{},
		plmnIds:  map[string]plmnIdVal{},
		timeouts: defaultTimeouts,
	}
	if config != nil {
		if config.Timeout != nil {
			if config.Timeout.ChallengeMs > 0 {
				service.SetChallengeTimeout(time.Millisecond * time.Duration(config.Timeout.ChallengeMs))
			}
			if config.Timeout.ErrorNotificationMs > 0 {
				service.SetNotificationTimeout(time.Millisecond * time.Duration(config.Timeout.ErrorNotificationMs))
			}
			if config.Timeout.SessionMs > 0 {
				service.SetSessionTimeout(time.Millisecond * time.Duration(config.Timeout.SessionMs))
			}
			if config.Timeout.SessionAuthenticatedMs > 0 {
				service.SetSessionAuthenticatedTimeout(
					time.Millisecond * time.Duration(config.Timeout.SessionAuthenticatedMs))
			}
		}
		for _, plmnid := range config.PlmnIds {
			l := len(plmnid)
			switch l {
			case 5:
				service.plmnIds[plmnid] = plmnIdVal{l5: true}
			case 6:
				plmnid5 := plmnid[:5]
				val, _ := service.plmnIds[plmnid5]
				val.b6 = plmnid[5]
				service.plmnIds[plmnid5] = val
			}
		}
	}
	return service, nil
}

// CheckPlmnId returns true either if there is no PLMN ID filters (whitelist) configured or
// one the configured PLMN IDs matches passed IMSI
func (s *EapAkaSrv) CheckPlmnId(imsi aka.IMSI) bool {
	if len(s.plmnIds) == 0 {
		return true
	}
	if val, ok := s.plmnIds[string(imsi)[:5]]; ok && (val.l5 || (len(imsi) > 5 && val.b6 == imsi[6])) {
		return true
	}
	return false
}

// Unlock - unlocks the CTX
func (lockedCtx *UserCtx) Unlock() {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	lockedCtx.locked = false
	lockedCtx.mu.Unlock()
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

// CreatedTime returns time of CTX creation
func (lockedCtx *UserCtx) CreatedTime() time.Time {
	return lockedCtx.created
}

// Lifetime returns duration in seconds of the CTX existence
func (lockedCtx *UserCtx) Lifetime() float64 {
	return time.Since(lockedCtx.created).Seconds()
}

// InitSession either creates new or updates existing session & user ctx,
// it session ID into the CTX and initializes session map as well as users map
// Returns Locked User Ctx
func (s *EapAkaSrv) InitSession(sessionId string, imsi aka.IMSI) (lockedUserContext *UserCtx) {
	var (
		oldSessionTimer *time.Timer
	)
	// create new session with long session wide timeout
	t := time.Now()
	newSession := &SessionCtx{UserCtx: &UserCtx{
		created: t, Imsi: imsi, state: aka.StateCreated, stateTime: t, locked: true, SessionId: sessionId}}

	newSession.mu.Lock()

	newSession.CleanupTimer = time.AfterFunc(s.SessionTimeout(), func() {
		sessionTimeoutCleanup(s, sessionId, newSession)
	})
	uc := newSession.UserCtx

	s.rwl.Lock()
	if oldSession, ok := s.sessions[sessionId]; ok && oldSession != nil {
		oldSessionTimer, oldSession.CleanupTimer = oldSession.CleanupTimer, nil
	}
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
	newSession = &SessionCtx{UserCtx: lockedCtx}
	sessionId := lockedCtx.SessionId
	lockedCtx.Unlock()

	newSession.CleanupTimer = time.AfterFunc(timeout, func() {
		sessionTimeoutCleanup(s, sessionId, newSession)
	})

	s.rwl.Lock()

	oldSession, exist = s.sessions[sessionId]
	s.sessions[sessionId] = newSession
	if exist && oldSession != nil {
		oldSession.UserCtx = nil
		if oldSession.CleanupTimer != nil {
			oldTimer, oldSession.CleanupTimer = oldSession.CleanupTimer, nil
		}
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
			newSession, oldSession.UserCtx = &SessionCtx{UserCtx: oldSession.UserCtx}, nil
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
	metrics.SessionTimeouts.Inc()
	if s == nil {
		log.Printf("ERROR: Nil EAP-AKA Server for session ID: %s", sessionId)
		return
	}
	var (
		imsi aka.IMSI
		uc   *UserCtx
	)

	s.rwl.Lock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist {
		if sessionCtx != nil {
			imsi = sessionCtx.Imsi
			if sessionCtx == mySessionCtx {
				delete(s.sessions, sessionId)
				uc = sessionCtx.UserCtx
			}
		} else {
			exist = false
		}
	}
	s.rwl.Unlock()

	if exist && uc != nil {
		uc.mu.Lock()
		state := uc.state
		uc.mu.Unlock()
		if state != aka.StateAuthenticated {
			log.Printf("EAP-AKA Session %s timeout for IMSI: %s", sessionId, imsi)
		}
	}
}

// FindSession finds and returns IMSI of a session and a flag indication if the find succeeded
// If found, FindSession tries to stop outstanding session timer
func (s *EapAkaSrv) FindSession(sessionId string) (aka.IMSI, *UserCtx, bool) {
	var (
		imsi      aka.IMSI
		lockedCtx *UserCtx
		timer     *time.Timer
	)
	s.rwl.RLock()
	sessionCtx, exist := s.sessions[sessionId]
	if exist && sessionCtx != nil {
		lockedCtx, timer, sessionCtx.CleanupTimer = sessionCtx.UserCtx, sessionCtx.CleanupTimer, nil
	}
	s.rwl.RUnlock()

	if lockedCtx != nil {
		lockedCtx.mu.Lock()
		lockedCtx.SessionId = sessionId // just in case - should always match
		imsi = lockedCtx.Imsi
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
			imsi, timer, sessionCtx.CleanupTimer, sessionCtx.UserCtx =
				sessionCtx.Imsi, sessionCtx.CleanupTimer, nil, nil
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
		if session != nil {
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
