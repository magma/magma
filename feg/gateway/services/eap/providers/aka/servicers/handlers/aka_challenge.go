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

// Package handlers provided AKA Response handlers for supported AKA subtypes
package handlers

import (
	"io"
	"reflect"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
)

func init() {
	servicers.AddHandler(aka.SubtypeChallenge, challengeResponse)
}

// challengeResponse implements handler for AKA Challenge Response,
// see https://tools.ietf.org/html/rfc4187#page-49 for details
func challengeResponse(s *servicers.EapAkaSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error) {
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
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing Session ID")
	}
	sessionId := ctx.SessionId
	imsi, uc, ok := s.FindSession(sessionId)
	if !ok {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No Session found for ID: %s", ctx.SessionId)
	}
	if uc == nil {
		s.UpdateSessionTimeout(sessionId, s.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No IMSI '%s' found for SessionID: %s", imsi, ctx.SessionId)
	}
	ctxCreated = uc.CreatedTime()

	state, _ := uc.State()
	if state != aka.StateChallenge {
		glog.Errorf(
			"AKA Challenge Response: Unexpected user state: %d for IMSI: %s, Session: %s", state, imsi, ctx.SessionId)
	}

	p := make([]byte, len(req))
	copy(p, req)
	scanner, err := eap.NewAttributeScanner(p)
	if err != nil {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}

	var a, atMac, atRes eap.Attribute

attrLoop:
	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		case aka.AT_MAC:
			atMac = a
			if atRes != nil {
				break attrLoop
			}
		case aka.AT_RES:
			atRes = a
			if atMac != nil {
				break attrLoop
			}
		case aka.AT_CHECKCODE: // Ignore CHECKCODE for now
		default:
			glog.Infof("Unexpected EAP-AKA Challenge Response Attribute type %d", a.Type())
		}
	}

	if err != nil {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		if err == io.EOF {
			return aka.EapErrorResPacket(
				identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing AT_MAC | AT_RES")
		}
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}

	// Verify MAC
	macBytes := atMac.Marshaled()
	if len(macBytes) < aka.ATT_HDR_LEN+aka.MAC_LEN {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Malformed AT_MAC")
	}
	ueMac := make([]byte, len(macBytes)-aka.ATT_HDR_LEN)
	copy(ueMac, macBytes[aka.ATT_HDR_LEN:])

	for i := aka.ATT_HDR_LEN; i < len(macBytes); i++ {
		macBytes[i] = 0
	}
	mac := aka.GenMac(p, uc.K_aut)
	if !reflect.DeepEqual(ueMac, mac) {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		glog.Errorf(
			"Invalid MAC for Session ID: %s; IMSI: %s; UE MAC: %x; Expected MAC: %x; EAP: %x",
			ctx.SessionId, imsi, ueMac, mac, req)
		return aka.EapErrorResPacket(
			identifier, aka.NOTIFICATION_FAILURE, codes.Unauthenticated,
			"Invalid MAC for Session ID: %s; IMSI: %s", ctx.SessionId, imsi)
	}

	// Verify AT_RES
	ueRes := atRes.Marshaled()[aka.ATT_HDR_LEN:]
	if success = reflect.DeepEqual(ueRes, uc.Xres); !success {
		glog.Errorf("Invalid AT_RES for Session ID: %s; IMSI: %s\n\t%.3v !=\n\t%.3v",
			sessionId, imsi, ueRes, uc.Xres)
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		return aka.EapErrorResPacketWithMac(
			identifier, aka.NOTIFICATION_FAILURE_AUTH, uc.K_aut, codes.Unauthenticated,
			"Invalid AT_RES for Session ID: %s; IMSI: %s", ctx.SessionId, imsi)
	}

	// All good, set IMSI, MSK & Identity for farther use by Radius and return SuccessCode
	ctx.Imsi = string(imsi)
	if uc.Profile != nil {
		ctx.Msisdn = uc.Profile.Msisdn
	}
	ctx.AuthSessionId = uc.AuthSessionId
	ctx.Msk = uc.MSK
	ctx.Identity = uc.Identity
	uc.SetState(aka.StateAuthenticated)

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
