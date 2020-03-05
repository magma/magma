/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package store provides an implementation for AAA Session and SessionTable interfaces
package store

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/metrics"
	"magma/feg/gateway/services/aaa/protos"
)

// Session - struct to save an authenticated session state
type memSession struct {
	*protos.Context
	imsi            string
	cleanupTimerCtx unsafe.Pointer // *cleanupTimerCtx
	mu              sync.Mutex
}

// Lock - locks the Session's mutex
func (s *memSession) Lock() {
	if s != nil {
		s.mu.Lock()
	}
}

// Unlock - unlocks the Session's mutex
func (s *memSession) Unlock() {
	if s != nil {
		s.mu.Unlock()
	}
}

// GetCtx returns AAA Session Context
func (s *memSession) GetCtx() *protos.Context {
	if s != nil {
		return s.Context
	}
	return nil
}

// SetCtx sets AAA Session Context - must be called on a Locked session
func (s *memSession) SetCtx(pc *protos.Context) {
	if s != nil {
		s.Context = pc
	}
}

// StopTimeout - stops the session's timeout if possible, returns if the timeout was successfully stopped
func (s *memSession) StopTimeout() bool {
	if s != nil {
		if ctx := atomic.SwapPointer(&s.cleanupTimerCtx, nil); ctx != nil {
			if t := atomic.SwapPointer(&((*cleanupTimerCtx)(ctx).sessionTimerPtr), nil); t != nil {
				return (*time.Timer)(t).Stop()
			}
		}
	}
	return false
}

// SessionTable - synchronized map of authenticated sessions
type memSessionTable struct {
	sm   map[string]*memSession
	sids map[string]string // Session IDs by IMSI: SID[IMSI]
	rwl  sync.RWMutex      // R/W lock synchronizing maps access
}

// NewSessionTable - returns a new initialized session table
func NewMemorySessionTable() aaa.SessionTable {
	return &memSessionTable{sm: map[string]*memSession{}, sids: map[string]string{}}
}

// AddSession - adds a new session to the table & returns the newly created session pointer.
// If a session with the same ID already is in the table - returns "Session with SID: XYZ already exist" as well as the
// existing session.
func (st *memSessionTable) AddSession(
	pc *protos.Context, tout time.Duration, notifier aaa.TimeoutNotifier, overwrite ...bool) (aaa.Session, error) {

	if st == nil {
		return nil, fmt.Errorf("Nil SessionTable")
	}
	if pc == nil {
		return nil, fmt.Errorf("Nil Session Context")
	}
	sid := strings.TrimSpace(pc.SessionId)
	if len(sid) == 0 {
		return nil, fmt.Errorf("Empty Session Id")
	}
	if tout < aaa.MinimalSessionTimeout {
		tout = aaa.MinimalSessionTimeout
	}

	imsi := pc.GetImsi()
	s := &memSession{Context: pc, imsi: imsi}
	st.rwl.Lock()
	// Handle the case of old session with the same radius session ID
	if oldSession, ok := st.sm[sid]; ok {
		if len(overwrite) > 0 && overwrite[0] {
			if oldSession != nil {
				oldImsi := oldSession.imsi
				log.Printf("Session with SID: %s already exist, will overwrite. Old IMSI: %s, New IMSI: %s",
					sid, oldImsi, imsi)

				oldSession.StopTimeout()

				if oldImsi != imsi {
					if oldSid, ok := st.sids[oldImsi]; ok && oldSid == sid {
						delete(st.sids, oldImsi)
					}
				}
			}
		} else {
			st.rwl.Unlock() // return old session is "best effort", done outside of the table lock
			return oldSession, fmt.Errorf("Session with SID: %s already exist", sid)
		}
	}
	// Handle the case of old session with the same IMSI and diferent radius session ID (roaming?)
	if oldSessionId, ok := st.sids[imsi]; ok && oldSessionId != sid {
		if oldImsiSession, ok := st.sm[oldSessionId]; ok {
			if oldImsiSession != nil {
				oldImsiSession.StopTimeout()
			}
			delete(st.sids, oldSessionId)
			log.Printf("old session with SID: %s found for IMSI: %s, will remove", oldSessionId, imsi)
		}
	}
	st.sm[sid] = s
	st.sids[imsi] = sid
	apn := s.GetApn()
	st.rwl.Unlock()

	setTimeoutUnsafe(st, sid, tout, s, notifier)
	imsi = metrics.DecorateIMSI(imsi)
	metrics.Sessions.WithLabelValues(apn, imsi, sid).Inc()
	metrics.SessionStart.WithLabelValues(apn, imsi, sid).Inc()

	return s, nil
}

// GetSession returns session corresponding to the given sid or nil if not found
func (st *memSessionTable) GetSession(sid string) aaa.Session {
	var s *memSession
	if st != nil {
		st.rwl.RLock()
		s, _ = st.sm[sid]
		st.rwl.RUnlock()
	}
	return s
}

// FindSession returns session corresponding to the given sid or nil if not found
func (st *memSessionTable) FindSession(imsi string) (sid string) {
	if st != nil {
		st.rwl.RLock()
		sid, _ = st.sids[imsi]
		st.rwl.RUnlock()
	}
	return sid
}

// RemoveSession - removes the session with the given SID and returns it
func (st *memSessionTable) RemoveSession(sid string) aaa.Session {
	if st != nil {
		var (
			found bool
			s     *memSession
		)
		var apn, imsi string
		st.rwl.Lock()
		if s, found = st.sm[sid]; found {
			apn, imsi = s.GetApn(), s.GetImsi()
			delete(st.sm, sid)
			if oldSid, ok := st.sids[s.imsi]; ok && oldSid == sid {
				delete(st.sids, s.imsi)
			}
		}
		st.rwl.Unlock()
		if found {
			s.StopTimeout()
			imsi = metrics.DecorateIMSI(imsi)
			metrics.Sessions.WithLabelValues(apn, imsi, sid).Dec()
			metrics.SessionStop.WithLabelValues(apn, imsi, sid).Inc()
			return s
		}
	}
	return nil
}

// SetTimeout - [Re]sets the session's cleanup timeout to fire after tout duration
func (st *memSessionTable) SetTimeout(sid string, tout time.Duration, notifier aaa.TimeoutNotifier) bool {
	var res bool
	if tout > 0 && st != nil && len(sid) > 0 {
		st.rwl.Lock()
		if s, ok := st.sm[sid]; ok && s != nil {
			setTimeoutUnsafe(st, sid, tout, s, notifier)
			res = true
		}
		st.rwl.Unlock()
	}
	return res
}

type cleanupTimerCtx struct {
	owner           *memSessionTable
	sidKey          string
	s               *memSession
	notifyRoutine   aaa.TimeoutNotifier
	sessionTimerPtr unsafe.Pointer
}

func setTimeoutUnsafe(st *memSessionTable, sid string, tout time.Duration, s *memSession, notifier aaa.TimeoutNotifier) {
	var ctx = &cleanupTimerCtx{owner: st, sidKey: sid, s: s, notifyRoutine: notifier}
	newTimer := time.AfterFunc(tout, func() { cleanupTimer(ctx) })
	atomic.StorePointer(&ctx.sessionTimerPtr, unsafe.Pointer(newTimer))
	atomic.StorePointer(&s.cleanupTimerCtx, unsafe.Pointer(ctx))
}

func cleanupTimer(ctx *cleanupTimerCtx) {
	if ctx != nil && ctx.s != nil && ctx.owner != nil {
		var deleted bool

		ctx.owner.rwl.Lock()
		if ctx.owner.sm != nil {
			if ms, ok := ctx.owner.sm[ctx.sidKey]; ok && ms == ctx.s {
				if atomic.CompareAndSwapPointer((*unsafe.Pointer)(&ms.cleanupTimerCtx), unsafe.Pointer(ctx), nil) {
					delete(ctx.owner.sm, ctx.sidKey)
					if oldSid, ok := ctx.owner.sids[ms.imsi]; ok && oldSid == ctx.sidKey {
						delete(ctx.owner.sids, ms.imsi)
					}
					deleted = true
				}
			}
		}
		ctx.owner.rwl.Unlock()

		if deleted {
			var notifyResult error
			s := ctx.s
			if ctx.notifyRoutine != nil {
				notifyResult = ctx.notifyRoutine(s)
			}
			log.Printf(
				"Timed out session '%s' for SessionId: %s; IMSI: %s; Identity: %s; MAC: %s; IP: %s; notify result: %v",
				ctx.sidKey, s.GetSessionId(), s.GetImsi(), s.GetIdentity(), s.GetMacAddr(), s.GetIpAddr(), notifyResult)

			metrics.SessionTimeouts.WithLabelValues(s.GetApn(), metrics.DecorateIMSI(s.GetImsi())).Inc()
		}
	}
}
