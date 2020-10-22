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

// package servcers implements EAP-SIM GRPC service
package servicers

import (
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/plmn_filter"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/metrics"
)

type UserCtx struct {
	mu         sync.Mutex
	created    time.Time
	state      sim.SimState
	stateTime  time.Time
	locked     bool
	Identity   string
	Imsi       sim.IMSI
	Profile    *protos.AuthenticationAnswer_UserProfile
	Identifier uint8
	Rand       [][]byte
	Sres       [][]byte
	K_aut,
	MSK []byte
	SessionId     string
	AuthSessionId string
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

type EapSimSrv struct {
	rwl sync.RWMutex // R/W lock synchronizing maps access
	// Map of UE Sessions keyed by sessionId
	sessions map[string]*SessionCtx

	// PLMN IDs map, if not empty -> serve only IMSIs with specified PLMN IDs - Read Only
	plmnFilter plmn_filter.PlmnIdVals
	timeouts   touts
	useS6a     bool
	mncLen     int32
}

var defaultTimeouts = touts{
	challengeTimeout:            sim.DefaultChallengeTimeout,
	errorNotificationTimeout:    sim.DefaultErrorNotificationTimeout,
	sessionTimeout:              sim.DefaultSessionTimeout,
	sessionAuthenticatedTimeout: sim.DefaultSessionAuthenticatedTimeout,
}

func (s *EapSimSrv) ChallengeTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.challengeTimeout)))
}

func (s *EapSimSrv) SetChallengeTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.challengeTimeout), int64(tout))
}

func (s *EapSimSrv) NotificationTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.errorNotificationTimeout)))
}

func (s *EapSimSrv) SetNotificationTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.errorNotificationTimeout), int64(tout))
}

func (s *EapSimSrv) SessionTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.sessionTimeout)))
}

func (s *EapSimSrv) SetSessionTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.sessionTimeout), int64(tout))
}

func (s *EapSimSrv) SessionAuthenticatedTimeout() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&s.timeouts.sessionAuthenticatedTimeout)))
}

func (s *EapSimSrv) SetSessionAuthenticatedTimeout(tout time.Duration) {
	atomic.StoreInt64((*int64)(&s.timeouts.sessionAuthenticatedTimeout), int64(tout))
}

// NewEapSimService creates new Sim Service 'object'
func NewEapSimService(config *mconfig.EapSimConfig) (*EapSimSrv, error) {
	service := &EapSimSrv{
		sessions:   map[string]*SessionCtx{},
		plmnFilter: plmn_filter.PlmnIdVals{},
		timeouts:   defaultTimeouts,
		mncLen:     3,
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
		service.plmnFilter = plmn_filter.GetPlmnVals(config.PlmnIds, "EAP-SIM")
		service.useS6a = config.GetUseS6A()
		if mncLn := config.GetMncLen(); mncLn >= 2 && mncLn <= 3 {
			service.mncLen = mncLn
		}
	}
	if useS6aStr, isset := os.LookupEnv("USE_S6A_BASED_AUTH"); isset {
		service.useS6a, _ = strconv.ParseBool(useS6aStr)
	}
	if service.useS6a {
		glog.Info("EAP-SIM: Using S6a Auth Vectors")
	} else {
		glog.Info("EAP-SIM: Using SWx Auth Vectors")
	}
	return service, nil
}

// CheckPlmnId returns true either if there is no PLMN ID filters configured or
// one the configured PLMN IDs matches given IMSI
func (s *EapSimSrv) CheckPlmnId(imsi sim.IMSI) bool {
	return s == nil || s.plmnFilter.Check(string(imsi))
}

func (s *EapSimSrv) UseS6a() bool {
	if s != nil {
		return s.useS6a
	}
	return false
}
func (s *EapSimSrv) MncLen() int {
	if s != nil {
		return int(s.mncLen)
	}
	return 3
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
func (lockedCtx *UserCtx) State() (sim.SimState, time.Time) {
	if !lockedCtx.locked {
		panic("Expected locked")
	}
	return lockedCtx.state, lockedCtx.stateTime
}

// SetState updates current CTX state (CTX must be locked)
func (lockedCtx *UserCtx) SetState(s sim.SimState) {
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
func (s *EapSimSrv) InitSession(sessionId string, imsi sim.IMSI) (lockedUserContext *UserCtx) {
	var (
		oldSessionTimer *time.Timer
		oldSessionState sim.SimState
	)
	// create new session with long session wide timeout
	t := time.Now()
	newSession := &SessionCtx{UserCtx: &UserCtx{
		created: t, Imsi: imsi, state: sim.StateCreated, stateTime: t, locked: true, SessionId: sessionId}}

	newSession.mu.Lock()

	newSession.CleanupTimer = time.AfterFunc(s.SessionTimeout(), func() {
		sessionTimeoutCleanup(s, sessionId, newSession)
	})
	uc := newSession.UserCtx

	s.rwl.Lock()
	if oldSession, ok := s.sessions[sessionId]; ok && oldSession != nil {
		oldSessionTimer, oldSessionState, oldSession.CleanupTimer = oldSession.CleanupTimer, oldSession.state, nil
	}
	s.sessions[sessionId] = newSession
	s.rwl.Unlock()

	if oldSessionTimer != nil {
		oldSessionTimer.Stop()
		// Copy Redirected state to a new session to avoid auth thrashing between EAP methods
		if oldSessionState == sim.StateRedirected {
			newSession.state = sim.StateRedirected
		}
	}
	return uc
}

// UpdateSessionUnlockCtx sets session ID into the CTX and initializes session map & session timeout
func (s *EapSimSrv) UpdateSessionUnlockCtx(lockedCtx *UserCtx, timeout time.Duration) {
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
func (s *EapSimSrv) UpdateSessionTimeout(sessionId string, timeout time.Duration) bool {
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

func sessionTimeoutCleanup(s *EapSimSrv, sessionId string, mySessionCtx *SessionCtx) {
	metrics.SessionTimeouts.Inc()
	if s == nil {
		glog.Errorf("nil EAP-SIM Server for session ID: %s", sessionId)
		return
	}
	var (
		imsi sim.IMSI
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
		if state != sim.StateAuthenticated {
			glog.Warningf("EAP-SIM Session %s timeout for IMSI: %s", sessionId, imsi)
		}
	}
}

// FindSession finds and returns IMSI of a session and a flag indication if the find succeeded
// If found, FindSession tries to stop outstanding session timer
func (s *EapSimSrv) FindSession(sessionId string) (sim.IMSI, *UserCtx, bool) {
	var (
		imsi      sim.IMSI
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
func (s *EapSimSrv) RemoveSession(sessionId string) sim.IMSI {
	var (
		timer *time.Timer
		imsi  sim.IMSI
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
func (s *EapSimSrv) FindAndRemoveSession(sessionId string) (sim.IMSI, bool) {
	var (
		imsi  sim.IMSI
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
func (s *EapSimSrv) ResetSessionTimeout(sessionId string, newTimeout time.Duration) {
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
