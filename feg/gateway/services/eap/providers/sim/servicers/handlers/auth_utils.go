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

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	swx_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
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

type tgppAuthResult struct {
	rand, Kc, sres [sim.GsmTripletsNumber][]byte
	sid            string
	profile        *swx_protos.AuthenticationAnswer_UserProfile
}

func getSwxVectors(_ *servicers.EapSimSrv, imsi string) (*tgppAuthResult, error) {
	var (
		err error
		ans *swx_protos.AuthenticationAnswer
		res tgppAuthResult
	)
	for vlen := 0; vlen < sim.GsmTripletsNumber; {
		metrics.SwxRequests.Inc()
		swxStartTime := time.Now()
		swxReq := &swx_protos.AuthenticationRequest{
			UserName:             imsi,
			SipNumAuthVectors:    sim.GsmTripletsNumber - uint32(vlen),
			AuthenticationScheme: swx_protos.AuthenticationScheme_EAP_AKA, // we are getting UMTS vectors, so - use AKA
			RetrieveUserProfile:  res.profile == nil,
		}
		ans, err = swx_proxy.Authenticate(swxReq)
		metrics.SWxLatency.Observe(time.Since(swxStartTime).Seconds())
		if err != nil {
			metrics.SwxFailures.Inc()
			errCode := codes.Internal
			if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
				errCode = se.GRPCStatus().Code()
			}
			return nil, status.Errorf(errCode, "%v; IMSI: %s", err, imsi)
		}
		if ans == nil {
			return nil, status.Error(codes.Internal, "Error: Nil SWx Response")
		}
		if len(ans.SipAuthVectors) == 0 {
			return nil, status.Errorf(codes.Internal, "Error: Missing/empty SWx Auth Vector: %+v", *ans)
		}
		for _, v := range ans.GetSipAuthVectors() {
			ra := v.GetRandAutn()
			if len(ra) < sim.RAND_LEN {
				return nil, status.Errorf(codes.Internal,
					"Invalid SWx RandAutn len (%d, expected: %d) in Response: %+v", len(ra), sim.RandAutnLen, *ans)
			}
			res.rand[vlen] = ra[:sim.RAND_LEN]
			res.Kc[vlen], res.sres[vlen] = sim.GsmFromUmts1(v.GetConfidentialityKey(), v.GetIntegrityKey(), v.GetXres())
			vlen++
			if vlen >= sim.GsmTripletsNumber {
				break
			}
		}
		if swxReq.RetrieveUserProfile {
			res.profile = ans.GetUserProfile()
		}
	}
	res.sid = ans.GetSessionId()
	return &res, nil
}

func createChallengeRequest(
	s *servicers.EapSimSrv,
	lockedCtx *servicers.UserCtx,
	identifier uint8,
	nonce, versionList, selectedVersion []byte) (eap.Packet, error) {

	var (
		err     error
		authRes *tgppAuthResult
	)
	if s.UseS6a() {
		authRes, err = getS6aVectors(s, string(lockedCtx.Imsi))
	} else {
		authRes, err = getSwxVectors(s, string(lockedCtx.Imsi))
	}
	if err != nil {
		if err == EUTRANOnlyVectorsErr {
			// User may only have EUTRAN Vectors (4G) try to steer UE toward EAP-AKA auth instead
			glog.Warningf("EUTRAN only Vectors for IMSI: %s, will try EAP-AKA ID Request", lockedCtx.Imsi)
			lockedCtx.SetState(sim.StateRedirected)
			return aka.NewIdentityReq(identifier+1, aka.AT_PERMANENT_ID_REQ), nil
		}
		var (
			code codes.Code
			msg  string
		)
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			code = se.GRPCStatus().Code()
			msg = se.GRPCStatus().Message()
		} else {
			code = codes.Internal
			msg = err.Error()
		}
		glog.Errorf("SIM RPC [%s] %s", code, msg)
		return sim.NewSIMNotificationReq(identifier, sim.NOTIFICATION_FAILURE), nil
	}
	identifier++
	lockedCtx.AuthSessionId = authRes.sid
	lockedCtx.Identifier = identifier
	lockedCtx.Rand = authRes.rand[:]
	lockedCtx.Sres = authRes.sres[:]
	_, lockedCtx.K_aut, lockedCtx.MSK, _ =
		sim.MakeKeys([]byte(lockedCtx.Identity), nonce, versionList, selectedVersion, authRes.Kc[:])

	// Clone EAP Challenge packet
	p := eap.Packet(make([]byte, challengeReqTemplateLen))
	copy(p, challengeReqTemplate)
	// Set current identifier
	p[eap.EapMsgIdentifier] = identifier
	// Set AT_RAND
	for i, offset := 0, atRandOffset; i < sim.GsmTripletsNumber; i, offset = i+1, offset+sim.RAND_LEN {
		copy(p[offset:], authRes.rand[i])
	}
	// Calculate AT_MAC
	mac := sim.GenMac(p, nonce, lockedCtx.K_aut)
	// Set AT_MAC
	copy(p[atMacOffset:], mac)
	return p, nil
}
