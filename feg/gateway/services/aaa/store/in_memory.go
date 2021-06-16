/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package store provides an implementation for AAA Session and SessionTable interfaces
package store

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/golang/glog"

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
	msisdn := pc.GetMsisdn()
	s := &memSession{Context: pc, imsi: imsi}
	st.rwl.Lock()

	// Handle the case of old session with the same radius session ID
	isExistingSession := false
	if oldSession, ok := st.sm[sid]; ok {
		if len(overwrite) > 0 && overwrite[0] {
			if oldSession != nil {
				isExistingSession = true
				oldImsi := oldSession.imsi
				glog.Warningf("Session with SID: %s already exist, will overwrite. Old IMSI: %s, New IMSI: %s",
					sid, oldImsi, imsi)

				oldSession.StopTimeout()

				if oldImsi != imsi {
					isExistingSession = false
					if oldSid, ok := st.sids[oldImsi]; ok && oldSid == sid {
						delete(st.sids, oldImsi)
					}
					updateSessionMetricsForRemovedSession(oldSession.GetApn(), oldImsi, sid, msisdn)
				}
			}
		} else {
			st.rwl.Unlock() // return old session is "best effort", done outside of the table lock
			return oldSession, fmt.Errorf("Session with SID: %s already exist", sid)
		}
	}
	// Handle the case of old session with the same IMSI and different radius session ID (roaming?)
	if oldSessionId, ok := st.sids[imsi]; ok && oldSessionId != sid {
		if oldImsiSession, ok := st.sm[oldSessionId]; ok {
			if oldImsiSession != nil {
				oldImsiSession.StopTimeout()
			}
			delete(st.sids, oldSessionId)
			updateSessionMetricsForRemovedSession(oldImsiSession.GetApn(), imsi, oldSessionId, msisdn)
			glog.Infof("old session with SID: %s found for IMSI: %s, will remove", oldSessionId, imsi)
		}
	}
	st.sm[sid] = s
	st.sids[imsi] = sid
	apn := s.GetApn()
	st.rwl.Unlock()

	glog.V(1).Infof("setting timeout of %f seconds for session: %s", tout.Seconds(), sid)
	setTimeoutUnsafe(st, sid, tout, s, notifier)
	if !isExistingSession {
		updateSessionMetricsForNewSession(apn, imsi, sid, msisdn)
	}
	return s, nil
}

// GetSession returns session corresponding to the given sid or nil if not found
func (st *memSessionTable) GetSession(sid string) aaa.Session {
	if st != nil {
		st.rwl.RLock()
		s, found := st.sm[sid]
		st.rwl.RUnlock()
		if found {
			return s
		}
	}
	return nil
}

// FindSession returns session corresponding to the given sid or nil if not found
func (st *memSessionTable) FindSession(imsi string) (sid string) {
	if st != nil {
		st.rwl.RLock()
		sid = st.sids[imsi]
		st.rwl.RUnlock()
	}
	return sid
}

// GetSessionByImsi returns session corresponding to the given IMSI or nil if not found
func (st *memSessionTable) GetSessionByImsi(imsi string) aaa.Session {
	if st != nil {
		var s aaa.Session
		st.rwl.RLock()
		defer st.rwl.RUnlock()
		if sid, found := st.sids[imsi]; found {
			if s, found = st.sm[sid]; found {
				return s
			}
		}
	}
	return nil
}

// RemoveSession - removes the session with the given SID and returns it
func (st *memSessionTable) RemoveSession(sid string) aaa.Session {
	if st != nil {
		var (
			found bool
			s     *memSession
		)
		var apn, imsi, msisdn string
		st.rwl.Lock()
		if s, found = st.sm[sid]; found {
			apn, imsi, msisdn = s.GetApn(), s.GetImsi(), s.GetMsisdn()
			delete(st.sm, sid)
			if oldSid, ok := st.sids[s.imsi]; ok && oldSid == sid {
				delete(st.sids, s.imsi)
			}
		}
		st.rwl.Unlock()
		if found {
			s.StopTimeout()
			updateSessionMetricsForRemovedSession(apn, imsi, sid, msisdn)
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
				if atomic.CompareAndSwapPointer(&ms.cleanupTimerCtx, unsafe.Pointer(ctx), nil) {
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
			glog.Infof(
				"Timed out session '%s' for SessionId: %s; IMSI: %s; Identity: %s; MAC: %s; IP: %s; notify result: %v",
				ctx.sidKey, s.GetSessionId(), s.GetImsi(), s.GetIdentity(), s.GetMacAddr(), s.GetIpAddr(), notifyResult)

			updateSessionMetricsForTimedOutSession(s.GetApn(), s.GetImsi(), s.GetSessionId(), s.GetMsisdn())
		}
	}
}

func updateSessionMetricsForNewSession(apn string, imsi string, sid string, msisdn string) {
	imsi = metrics.DecorateIMSI(imsi)
	metrics.Sessions.WithLabelValues(apn, imsi, sid, msisdn).Inc()
	metrics.SessionStart.WithLabelValues(apn, imsi, sid, msisdn).Inc()
}

func updateSessionMetricsForRemovedSession(apn string, imsi string, sid string, msisdn string) {
	imsi = metrics.DecorateIMSI(imsi)
	metrics.Sessions.WithLabelValues(apn, imsi, sid, msisdn).Dec()
	metrics.SessionStop.WithLabelValues(apn, imsi, sid, msisdn).Inc()
}

func updateSessionMetricsForTimedOutSession(apn string, imsi string, sid string, msisdn string) {
	imsi = metrics.DecorateIMSI(imsi)
	metrics.Sessions.WithLabelValues(apn, imsi, sid, msisdn).Dec()
	metrics.SessionStop.WithLabelValues(apn, imsi, sid, msisdn).Inc()
	metrics.SessionTimeouts.WithLabelValues(apn, imsi, msisdn).Inc()
}
