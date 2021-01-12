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

// Package handlers provided SIM Response handlers for supported SIM subtypes
package handlers

import (
	"io"
	"reflect"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/metrics"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
)

func init() {
	servicers.AddHandler(sim.SubtypeChallenge, challengeResponse)
}

// challengeResponse implements handler for SIM Challenge Response,
// see https://tools.ietf.org/html/rfc4186 for details
func challengeResponse(s *servicers.EapSimSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error) {
	var (
		success    bool
		ctxCreated time.Time
	)
	metrics.ChallengeRequests.Inc()
	defer func() {
		if !ctxCreated.IsZero() {
			metrics.AuthLatency.Observe(time.Since(ctxCreated).Seconds())
		}
		if !success {
			metrics.FailedChallengeRequests.Inc()
		}
	}()
	identifier := req.Identifier()
	if ctx == nil {
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing Session ID")
	}
	sessionId := ctx.SessionId
	imsi, uc, ok := s.FindSession(sessionId)
	if !ok {
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No Session found for ID: %s", ctx.SessionId)
	}
	if uc == nil {
		s.UpdateSessionTimeout(sessionId, s.NotificationTimeout())
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No IMSI '%s' found for SessionID: %s", imsi, ctx.SessionId)
	}
	ctxCreated = uc.CreatedTime()

	state, _ := uc.State()
	if state != sim.StateChallenge {
		glog.Errorf(
			"SIM Challenge Response: Unexpected user state: %d for IMSI: %s, Session: %s", state, imsi, ctx.SessionId)
	}
	p := make([]byte, len(req))
	copy(p, req)
	scanner, err := eap.NewAttributeScanner(p)
	if err != nil {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}

	var a, atMac eap.Attribute
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		if a.Type() == sim.AT_MAC {
			atMac = a
			break
		}
		glog.Infof("Unexpected EAP-SIM Challenge Response Attribute type %d", a.Type())
	}
	if err != nil {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		if err == io.EOF {
			return sim.EapErrorResPacket(
				identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing AT_MAC")
		}
		return sim.EapErrorResPacket(
			identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}
	// Verify MAC
	macBytes := atMac.Marshaled()
	if len(macBytes) < sim.ATT_HDR_LEN+sim.MAC_LEN {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		return sim.EapErrorResPacket(
			identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, "Malformed AT_MAC")
	}
	ueMac := make([]byte, len(macBytes)-sim.ATT_HDR_LEN)
	copy(ueMac, macBytes[sim.ATT_HDR_LEN:])
	// Set MAC value to zeros before calculating it for verification
	for i := sim.ATT_HDR_LEN; i < len(macBytes); i++ {
		macBytes[i] = 0
	}
	mac := sim.GenChallengeMac(p, uc.Sres, uc.K_aut)
	if !reflect.DeepEqual(ueMac, mac) {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		glog.Errorf(
			"Invalid MAC for Session ID: %s; IMSI: %s; UE MAC: %x; Expected MAC: %x; EAP: %x",
			ctx.SessionId, imsi, ueMac, mac, req)
		return sim.EapErrorResPacket(
			identifier, sim.NOTIFICATION_FAILURE, codes.Unauthenticated,
			"Invalid MAC for Session ID: %s; IMSI: %s", ctx.SessionId, imsi)
	}
	// All good, set IMSI, MSK & Identity for farther use by Radius and return SuccessCode
	ctx.Imsi = string(imsi)
	if uc.Profile != nil {
		ctx.Msisdn = uc.Profile.Msisdn
	}
	ctx.Msk = uc.MSK
	ctx.Identity = uc.Identity
	ctx.AuthSessionId = uc.AuthSessionId
	uc.SetState(sim.StateAuthenticated)

	// Keep session & User Ctx around for some time after authentication and then clean them up
	uc.Unlock()
	s.ResetSessionTimeout(sessionId, s.SessionAuthenticatedTimeout())

	// RFC 3748 p4.2 EAP Success packet
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |     Code      |  Identifier   |            Length             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	return []byte{
			eap.SuccessCode, // Code
			identifier,      // Identifier
			0, 4},           // Length
		nil
}
