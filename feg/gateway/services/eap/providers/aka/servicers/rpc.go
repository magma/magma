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

// package servcers implements EAP-AKA GRPC service
package servicers

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
)

// Handle implements AKA handler RPC
func (s *EapAkaSrv) Handle(_ context.Context, req *protos.Eap) (*protos.Eap, error) {
	return s.HandleImpl(req)
}

// Handle implements AKA handler API
func (s *EapAkaSrv) HandleImpl(req *protos.Eap) (*protos.Eap, error) {
	failure := true
	metrics.Requests.Inc()
	defer func() {
		if failure {
			metrics.FailedRequests.Inc()
		}
	}()

	p := eap.Packet(req.GetPayload())
	eapCtx := req.GetCtx()
	if eapCtx == nil {
		eapCtx = &protos.Context{}
	}
	if p == nil {
		return aka.EapErrorRes(0, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx, "Nil Request")
	}
	err := p.Validate()
	if err != nil {
		identifier := byte(0)
		if err != io.ErrShortBuffer {
			identifier = p.Identifier()
		}
		return aka.EapErrorRes(identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx, err.Error())
	}
	identifier := p.Identifier()
	method := p.Type()
	if method == eap.MethodIdentity {
		return &protos.Eap{Payload: aka.NewIdentityReq(identifier+1, aka.AT_PERMANENT_ID_REQ), Ctx: eapCtx}, nil
	}
	if method != aka.TYPE {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.Unimplemented, eapCtx, "Wrong EAP Method: %d", method)
	}
	if len(p) < aka.MIN_PACKET_LEN {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx,
			"EAP-AKA Packet is too short: %d", len(p))
	}
	h := GetHandler(aka.Subtype(p[eap.EapSubtype]))
	if h == nil {
		return aka.EapErrorRes(
			identifier, aka.NOTIFICATION_FAILURE, codes.NotFound, eapCtx,
			"Unsuported Subtype: %d", p[eap.EapSubtype])
	}
	rp, err := h(s, eapCtx, p)
	failure = err != nil
	return &protos.Eap{Payload: rp, Ctx: eapCtx}, err
}
