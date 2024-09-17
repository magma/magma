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
	"context"
	"io"

	"google.golang.org/grpc/codes"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/metrics"
)

// Handle implements SIM handler RPC
func (s *EapSimSrv) Handle(_ context.Context, req *protos.Eap) (*protos.Eap, error) {
	return s.HandleImpl(req)
}

// Handle implements SIM handler API
func (s *EapSimSrv) HandleImpl(req *protos.Eap) (*protos.Eap, error) {
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
		return sim.EapErrorRes(0, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx, "Nil Request")
	}
	err := p.Validate()
	if err != nil {
		identifier := byte(0)
		if err != io.ErrShortBuffer {
			identifier = p.Identifier()
		}
		return sim.EapErrorRes(identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx, err.Error())
	}
	identifier := p.Identifier()
	method := p.Type()
	if method == eap.MethodIdentity {
		return &protos.Eap{Payload: sim.NewStartReq(identifier+1, sim.AT_PERMANENT_ID_REQ), Ctx: eapCtx}, nil
	}
	if method != sim.TYPE {
		return sim.EapErrorRes(
			identifier, sim.NOTIFICATION_FAILURE, codes.Unimplemented, eapCtx, "Wrong EAP Method: %d", method)
	}
	if len(p) < sim.MIN_PACKET_LEN {
		return sim.EapErrorRes(
			identifier, sim.NOTIFICATION_FAILURE, codes.InvalidArgument, eapCtx,
			"EAP-SIM Packet is too short: %d", len(p))
	}
	h := GetHandler(sim.Subtype(p[eap.EapSubtype]))
	if h == nil {
		return sim.EapErrorRes(
			identifier, sim.NOTIFICATION_FAILURE, codes.NotFound, eapCtx,
			"Unsuported Subtype: %d", p[eap.EapSubtype])
	}
	rp, err := h(s, eapCtx, p)
	failure = err != nil
	return &protos.Eap{Payload: rp, Ctx: eapCtx}, err
}
