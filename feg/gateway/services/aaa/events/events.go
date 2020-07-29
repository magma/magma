/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package events

import (
	"encoding/json"
	"fmt"

	"magma/feg/gateway/services/aaa/protos"
	"magma/gateway/eventd"
	"magma/gateway/status"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

type sessionEventBase struct {
	Imsi      string `json:"imsi"`
	SessionID string `json:"session_id"`
	MacAddr   string `json:"mac_addr"`
	Apn       string `json:"apn"`
}

// AuthenticationSucceeded event
type AuthenticationSucceeded struct {
	sessionEventBase
}

// AuthenticationFailed event
type AuthenticationFailed struct {
	sessionEventBase
	FailureReason string `json:"failure_reason"`
}

// SessionTerminationSucceeded event
type SessionTerminationSucceeded struct {
	sessionEventBase
	TerminationReason SessionTerminationReason `json:"session_termination_reason"`
}

// SessionTerminationFailed event
type SessionTerminationFailed struct {
	sessionEventBase
	TerminationReason SessionTerminationReason `json:"session_termination_reason"`
	FailureReason     string                   `json:"failure_reason"`
}

type SessionTerminationReason string

const (
	sessionsStreamName               = "sessions"
	authenticationSucceedEvent       = "authentication_succeeded"
	authenticationFailedEvent        = "authentication_failed"
	sessionTerminationSucceededEvent = "session_termination_succeeded"
	sessionTerminationFailedEvent    = "session_termination_failed"

	AccountingStop          SessionTerminationReason = "accounting_stop"
	SessionTimeout          SessionTerminationReason = "session_timeout"
	AbortSession            SessionTerminationReason = "abort_session"
	RegistrationTermination SessionTerminationReason = "registration_termination"
)

func LogAuthenticationSuccessEvent(ctx *protos.Context) {
	if ctx == nil {
		glog.Errorf("Could not log authentication succeeded event; ctx is nil")
		return
	}
	authSucceededEvent := AuthenticationSucceeded{
		sessionEventBase: sessionEventBase{
			Imsi:      ctx.GetImsi(),
			SessionID: ctx.GetSessionId(),
			MacAddr:   ctx.GetMacAddr(),
			Apn:       ctx.GetApn(),
		},
	}
	serializedAuthEvent, err := json.Marshal(authSucceededEvent)
	if err != nil {
		glog.Errorf("Could not serialize authentication succeeded event: %s", err)
		return
	}
	logEvent(sessionsStreamName, authenticationSucceedEvent, serializedAuthEvent)
}

func LogAuthenticationFailedEvent(ctx *protos.Context, failureReason string) {
	if ctx == nil {
		glog.Errorf("Could not log authentication failed event; ctx is nil")
		return
	}
	authFailedEvent := AuthenticationFailed{
		sessionEventBase: sessionEventBase{
			Imsi:      ctx.GetImsi(),
			SessionID: ctx.GetSessionId(),
			MacAddr:   ctx.GetMacAddr(),
			Apn:       ctx.GetApn(),
		},
		FailureReason: failureReason,
	}
	serializedAuthEvent, err := json.Marshal(authFailedEvent)
	if err != nil {
		glog.Errorf("Could not serialize authentication failed event: %s", err)
		return
	}
	logEvent(sessionsStreamName, authenticationFailedEvent, serializedAuthEvent)
}

func LogSessionTerminationSucceededEvent(ctx *protos.Context, terminationReason SessionTerminationReason) {
	if ctx == nil {
		glog.Errorf("Could not log session termination succeeded event; ctx is nil")
		return
	}
	sessionTerminationSucceded := SessionTerminationSucceeded{
		sessionEventBase: sessionEventBase{
			Imsi:      ctx.GetImsi(),
			SessionID: ctx.GetSessionId(),
			MacAddr:   ctx.GetMacAddr(),
			Apn:       ctx.GetApn(),
		},
		TerminationReason: terminationReason,
	}
	serializedTermEvent, err := json.Marshal(sessionTerminationSucceded)
	if err != nil {
		glog.Errorf("Could not serialize session termination succeeded event: %s", err)
		return
	}
	logEvent(sessionsStreamName, sessionTerminationSucceededEvent, serializedTermEvent)
}

func LogSessionTerminationFailedEvent(ctx *protos.Context, terminationReason SessionTerminationReason, failureReason string) {
	if ctx == nil {
		glog.Errorf("Could not log session termination failed event; ctx is nil")
		return
	}
	sessionTerminationFailed := SessionTerminationFailed{
		sessionEventBase: sessionEventBase{
			Imsi:      ctx.GetImsi(),
			SessionID: ctx.GetSessionId(),
			MacAddr:   ctx.GetMacAddr(),
			Apn:       ctx.GetApn(),
		},
		TerminationReason: terminationReason,
		FailureReason:     failureReason,
	}
	serializedTermEvent, err := json.Marshal(sessionTerminationFailed)
	if err != nil {
		glog.Errorf("Could not serialize session termination failed event: %s", err)
		return
	}
	logEvent(sessionsStreamName, sessionTerminationFailedEvent, serializedTermEvent)
}

func logEvent(streamName string, eventType string, serializedEvent []byte) {
	hwid := status.GetHwId()
	event := &orcprotos.Event{
		StreamName: streamName,
		EventType:  eventType,
		Tag:        hwid,
		Value:      fmt.Sprintf("%s", serializedEvent),
	}
	err := eventd.V(eventd.DefaultVerbosity).Log(event)
	if err != nil {
		glog.Errorf("Sending %s event failed: %s", eventType, err)
	}
}
