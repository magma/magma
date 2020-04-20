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
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
)

func init() {
	servicers.AddHandler(aka.SubtypeSynchronizationFailure, resyncResponse)
}

// resyncResponse implements handler for EAP-Response/AKA-Synchronization-Failure,
// see https://tools.ietf.org/html/rfc4187#section-9.6 for details
func resyncResponse(s *servicers.EapAkaSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error) {
	var success bool
	metrics.ResyncRequests.Inc()
	defer func() {
		if !success {
			metrics.FailedResyncRequests.Inc()
		}
	}()
	identifier := req.Identifier()
	if ctx == nil {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing Session ID")
	}
	imsi, uc, ok := s.FindSession(ctx.SessionId)
	if !ok {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No Session found for ID: %s", ctx.SessionId)
	}
	if uc == nil {
		s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
			"No IMSI '%s' found for SessionID: %s", imsi, ctx.SessionId)
	}
	ctx.Imsi = string(imsi) // set IMSI

	p := make([]byte, len(req))
	copy(p, req)
	scanner, err := eap.NewAttributeScanner(p)
	if err != nil {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}

	state, t := uc.State()
	if state != aka.StateChallenge {
		glog.Errorf(
			"AKA-Synchronization-Failure: Overwriting unexpected user state: %d,%s for IMSI: %s",
			state, t, imsi)
	}
	uc.SetState(aka.StateIdentity)

	var a eap.Attribute

	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		if a.Type() == aka.AT_AUTS {
			auts := a.Value()
			if len(auts) < 14 {
				s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
				return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Invalid AT_AUTS LKen: %d", len(auts))
			}
			// Resync Info = RAND | AUTS
			resyncInfo := append(append(make([]byte, 0, len(uc.Rand)+len(auts)), uc.Rand...), auts...)
			p, err := createChallengeRequest(s, uc, identifier, resyncInfo)
			if success = err == nil; success {
				// Update state
				uc.SetState(aka.StateChallenge)
				s.UpdateSessionUnlockCtx(uc, s.ChallengeTimeout())
			} else {
				s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
			}
			return p, err
		}
	}

	s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
	return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Missing AT_AUTS")
}
