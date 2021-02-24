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
	"fmt"
	"io"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/metrics"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
)

func init() {
	servicers.AddHandler(sim.SubtypeStart, startResponse)
}

// startResponse implements handler for SIM Challenge, see https://tools.ietf.org/html/rfc4186#section-9.3 for reference
func startResponse(s *servicers.EapSimSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error) {
	var success bool
	metrics.StartRequests.Inc()
	defer func() {
		if !success {
			metrics.FailedStartRequests.Inc()
		}
	}()
	identifier := req.Identifier()
	if ctx == nil {
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, "Nil CTX")
	}
	if len(ctx.SessionId) == 0 {
		ctx.SessionId = eap.CreateSessionId()
		glog.Warningf("Missing Session ID for EAP: %x; Generated new SID: %s", req, ctx.SessionId)
	}
	scanner, err := eap.NewAttributeScanner(req)
	if err != nil {
		s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.Aborted, err.Error())
	}
	var (
		a              eap.Attribute
		identity       string
		imsi           sim.IMSI
		nonce, version []byte
	)

	for a, err = scanner.Next(); err == nil; a, err = scanner.Next() {
		switch a.Type() {
		// Find first valid AT_IDENTITY attribute to get UE IMSI
		case sim.AT_IDENTITY:
			identity, imsi, err = getIMSIIdentity(a)
			if err != nil {
				return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
			}
			if imsi[0] != '1' {
				glog.Warningf("SIM AT_IDENTITY '%s' (IMSI: %s) is non-permanent type", identity, imsi)
			} else {
				imsi = imsi[1:]
			}
			if !s.CheckPlmnId(imsi) {
				s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
				return sim.EapErrorResPacket(
					identifier,
					sim.NOTIFICATION_FAILURE,
					codes.PermissionDenied,
					"PLMN ID of IMSI: %s is not permitted", imsi)
			}
		case sim.AT_NONCE_MT:
			if a.AttrLen() != 5 {
				return sim.EapErrorResPacket(
					identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Invalid AT_NONCE_MT attribute length: %d (expected: 5)", a.Len())
			}
			nonce = a.Value()
			if len(nonce) < 18 {
				return sim.EapErrorResPacket(
					identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Malformed AT_NONCE_MT attribute byte length: %d (expected: 18)", len(nonce))
			}
			nonce = nonce[2:]
		case sim.AT_SELECTED_VERSION:
			if a.AttrLen() != 1 {
				return sim.EapErrorResPacket(
					identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Invalid AT_SELECTED_VERSION attribute length: %d (expected: 1)", a.Len())
			}
			version = a.Value()
			if len(version) < 2 {
				return sim.EapErrorResPacket(
					identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Malformed AT_SELECTED_VERSION attribute byte length: %d (expected: 2)", len(version))
			}
			if version[1] != sim.Version {
				return sim.EapErrorResPacket(
					identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument,
					"Unsupported SIM VERSION: %d (supported: 1)", version[1])
			}
		}
	}

	if err != nil && err != io.EOF {
		return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, err.Error())
	}
	if len(imsi) == 0 {
		return sim.EapErrorResPacket(
			identifier, sim.NOTIFICATION_FAILURE, codes.FailedPrecondition, "Missing AT_IDENTITY Attribute")
	}
	if len(nonce) == 0 {
		return sim.EapErrorResPacket(
			identifier, sim.NOTIFICATION_FAILURE, codes.FailedPrecondition, "Missing AT_NONCE_MT Attribute")
	}
	if len(version) == 0 {
		return sim.EapErrorResPacket(
			identifier, sim.NOTIFICATION_FAILURE, codes.FailedPrecondition, "Missing AT_SELECTED_VERSION Attribute")
	}

	if !s.CheckPlmnId(imsi) {
		s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
		return sim.EapErrorResPacket(
			identifier,
			sim.NOTIFICATION_FAILURE,
			codes.PermissionDenied,
			"PLMN ID of IMSI: %s is not in served PLMN ID list", imsi)
	}

	ctx.Imsi = string(imsi)                  // set IMSI
	uc := s.InitSession(ctx.SessionId, imsi) // we have Locked User Ctx after this call
	state, t := uc.State()
	if state > sim.StateCreated {
		glog.Errorf(
			"EAP SIM StartResponse: Unexpected user state: %d,%s for IMSI: %s, CTX Identity: %s",
			state, t, imsi, uc.Identity)
		if state == sim.StateRedirected {
			sim.EapErrorResPacket(
				identifier, sim.NOTIFICATION_FAILURE, codes.FailedPrecondition,
				"IMSI: %s is redirected to another method", imsi)
		}
	}
	uc.Identity = identity
	uc.SetState(sim.StateIdentity)
	p, err := createChallengeRequest(s, uc, identifier, nonce, []byte{0, sim.Version}, version)
	if err == nil {
		// Update state
		uc.SetState(sim.StateChallenge)
		s.UpdateSessionUnlockCtx(uc, s.ChallengeTimeout())
	} else {
		s.UpdateSessionUnlockCtx(uc, s.NotificationTimeout())
	}
	s.UpdateSessionTimeout(ctx.SessionId, s.NotificationTimeout())
	return p, err
}

// see https://tools.ietf.org/html/rfc4187#section-4.1.1.4
func getIMSIIdentity(a eap.Attribute) (string, sim.IMSI, error) {
	if a.Type() != sim.AT_IDENTITY {
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
	var imsi sim.IMSI
	if atIdx > 0 {
		imsi = sim.IMSI(fullIdentity[:atIdx])
	} else {
		imsi = sim.IMSI(fullIdentity)
	}
	return fullIdentity, imsi, imsi.Validate()
}
