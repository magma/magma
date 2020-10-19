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
	"fmt"
	"io"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
)

func init() {
	servicers.AddHandler(aka.SubtypeIdentity, identityResponse)
}

// identityResponse implements handler for AKA Challenge, see https://tools.ietf.org/html/rfc4187#page-49 for reference
func identityResponse(s *servicers.EapAkaSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error) {
	var success bool
	metrics.IdentityRequests.Inc()
	defer func() {
		if !success {
			metrics.FailedIdentityRequests.Inc()
		}
	}()
	identifier := req.Identifier()
	if ctx == nil {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		ctx.SessionId = eap.CreateSessionId()
		glog.Warningf("Missing Session ID for EAP: %x; Generated new SID: %s", req, ctx.SessionId)
	}
	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}
	var a eap.Attribute

	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		// Find first valid AT_IDENTITY attribute to get UE IMSI
		if a.Type() == aka.AT_IDENTITY {
			identity, imsi, err := getIMSIIdentity(a)
			if err == nil {
				if imsi[0] != '0' {
					glog.Warningf("AKA AT_IDENTITY '%s' (IMSI: %s) is non-permanent type", identity, imsi)
				} else {
					imsi = imsi[1:]
				}
				if !s.CheckPlmnId(imsi) {
					s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
					return aka.EapErrorResPacket(
						identifier,
						aka.NOTIFICATION_FAILURE,
						codes.PermissionDenied,
						"PLMN ID of IMSI: %s is not permitted", imsi)
				}
				ctx.Imsi = string(imsi)                  // set IMSI
				uc := s.InitSession(ctx.SessionId, imsi) // we have Locked User Ctx after this call
				state, t := uc.State()
				if state > aka.StateCreated {
					glog.Errorf(
						"EAP AKA IdentityResponse: Unexpected user state: %d,%s for IMSI: %s, CTX Identity: %s",
						state, t, imsi, uc.Identity)
					if state == aka.StateRedirected {
						aka.EapErrorResPacket(
							identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition,
							"IMSI: %s is redirected to another method", imsi)
					}
				}
				uc.Identity = identity
				uc.SetState(aka.StateIdentity)
				p, err := createChallengeRequest(s, uc, identifier, nil)
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
	}
	s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
	if err != nil && err != io.EOF {
		return aka.EapErrorResPacket(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}
	return aka.EapErrorResPacket(
		identifier, aka.NOTIFICATION_FAILURE, codes.FailedPrecondition, "Missing AT_IDENTITY Attribute")
}

// see https://tools.ietf.org/html/rfc4187#section-4.1.1.4
func getIMSIIdentity(a eap.Attribute) (string, aka.IMSI, error) {
	if a.Type() != aka.AT_IDENTITY {
		return "", "", fmt.Errorf("Unexpected Attr Type: %d, AT_IDENTITY expected", a.Type())
	}
	if a.Len() <= 4 {
		return "", "", fmt.Errorf("AT_IDENTITY is too short: %d", a.Len())
	}
	val := a.Value()
	actualLen2 := int(val[0])<<8 + int(val[1]) + 2
	if actualLen2 > len(val) {
		return "", "", fmt.Errorf("Corrupt AT_IDENTITY Attribute: actual len %d > data len %d", actualLen2-2, len(val))
	}
	fullIdentity := string(val[2:actualLen2])
	atIdx := strings.Index(fullIdentity, "@")
	var imsi aka.IMSI
	if atIdx > 0 {
		imsi = aka.IMSI(fullIdentity[:atIdx])
	} else {
		imsi = aka.IMSI(fullIdentity)
	}
	return fullIdentity, imsi, imsi.Validate()
}
