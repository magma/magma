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

package handlers

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	swx_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/metrics"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	"magma/feg/gateway/services/swx_proxy"
)

var (
	challengeReqTemplate eap.Packet
	challengeReqTemplateLen,
	// Offsets in the challenge template of corresponding attribute values
	atRandOffset,
	atMacOffset int
)

func init() {
	var err error
	p := eap.NewPacket(eap.RequestCode, 0, []byte{sim.TYPE, byte(sim.SubtypeChallenge), 0, 0})
	atRandOffset = len(p) + sim.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		sim.AT_RAND, append(
			[]byte{0, 0}, // reserved
			make([]byte, sim.RAND_LEN*sim.GsmTripletsNumber)...)))
	if err != nil {
		panic(err)
	}
	atMacOffset = len(p) + sim.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		sim.AT_MAC, append(
			[]byte{0, 0}, // reserved
			make([]byte, sim.MAC_LEN)...)))
	if err != nil {
		panic(err)
	}
	challengeReqTemplateLen = len(p)
	challengeReqTemplate = p
}

type umtsVector struct{ rand, autn, xres, ck, ik []byte }

func createChallengeRequest(
	lockedCtx *servicers.UserCtx,
	identifier uint8,
	nonce, versionList, selectedVersion []byte) (eap.Packet, error) {

	var (
		err  error
		ans  *swx_protos.AuthenticationAnswer
		rand [sim.GsmTripletsNumber][]byte
		Kc   [sim.GsmTripletsNumber][]byte
		sres [sim.GsmTripletsNumber][]byte
	)
	metrics.SwxRequests.Inc()
	swxStartTime := time.Now()
	vectors := make([]*umtsVector, sim.GsmTripletsNumber)

	lockedCtx.Profile = nil
	for vlen := 0; vlen < sim.GsmTripletsNumber; {
		swxReq := &swx_protos.AuthenticationRequest{
			UserName:             string(lockedCtx.Imsi),
			SipNumAuthVectors:    sim.GsmTripletsNumber - uint32(vlen),
			AuthenticationScheme: swx_protos.AuthenticationScheme_EAP_AKA, // we are getting UMTS vectors, so - use AKA
			RetrieveUserProfile:  lockedCtx.Profile == nil,
		}
		ans, err = swx_proxy.Authenticate(swxReq)

		metrics.SWxLatency.Observe(time.Since(swxStartTime).Seconds())

		if err != nil {
			metrics.SwxFailures.Inc()
			errCode := codes.Internal
			if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
				errCode = se.GRPCStatus().Code()
			}
			return sim.EapErrorResPacket(identifier, sim.NOTIFICATION_FAILURE, errCode, err.Error())
		}
		if ans == nil {
			return sim.EapErrorResPacket(
				identifier, sim.NOTIFICATION_FAILURE, codes.Internal, "Error: Nil SWx Response")
		}
		if len(ans.SipAuthVectors) == 0 {
			return sim.EapErrorResPacket(
				identifier, sim.NOTIFICATION_FAILURE, codes.Internal, "Error: Missing/empty SWx Auth Vector: %+v", *ans)
		}
		for _, v := range ans.GetSipAuthVectors() {
			ra := v.GetRandAutn()
			if len(ra) < sim.RandAutnLen {
				return sim.EapErrorResPacket(
					identifier,
					sim.NOTIFICATION_FAILURE,
					codes.Internal,
					"Invalid SWx RandAutn len (%d, expected: %d) in Response: %+v",
					len(ra), sim.RandAutnLen, *ans)
			}
			vectors[vlen] = &umtsVector{
				rand: ra[:sim.RAND_LEN],
				autn: ra[sim.RAND_LEN:sim.RandAutnLen],
				xres: v.GetXres(),
				ck:   v.GetConfidentialityKey(),
				ik:   v.GetIntegrityKey()}
			vlen++
			if vlen >= sim.GsmTripletsNumber {
				break
			}
		}
		if swxReq.RetrieveUserProfile {
			lockedCtx.Profile = ans.GetUserProfile()
		}
	}
	for i, v := range vectors {
		rand[i] = v.rand
		Kc[i], sres[i] = sim.GsmFromUmts1(v.ck, v.ik, v.xres)
	}
	identifier++
	lockedCtx.AuthSessionId = ans.GetSessionId()
	lockedCtx.Identifier = identifier
	lockedCtx.Rand = rand[:]
	lockedCtx.Sres = sres[:]
	_, lockedCtx.K_aut, lockedCtx.MSK, _ =
		sim.MakeKeys([]byte(lockedCtx.Identity), nonce, versionList, selectedVersion, Kc[:])

	// Clone EAP Challenge packet
	p := eap.Packet(make([]byte, challengeReqTemplateLen))
	copy(p, challengeReqTemplate)
	// Set current identifier
	p[eap.EapMsgIdentifier] = identifier
	// Set AT_RAND
	for i, offset := 0, atRandOffset; i < sim.GsmTripletsNumber; i, offset = i+1, offset+sim.RAND_LEN {
		copy(p[offset:], rand[i])
	}
	// Calculate AT_MAC
	mac := sim.GenMac(p, nonce, lockedCtx.K_aut)
	// Set AT_MAC
	copy(p[atMacOffset:], mac)
	return p, nil
}
