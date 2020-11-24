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

	"magma/feg/gateway/services/eap/providers/sim"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	swx_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	"magma/feg/gateway/services/swx_proxy"
)

var (
	challengeReqTemplate eap.Packet
	challengeReqTemplateLen,
	// Offsets in the challenge template of corresponding attribute values
	atRandOffset,
	atAutnOffset,
	atMacOffset int
)

func init() {
	var err error
	p := eap.NewPacket(eap.RequestCode, 0, []byte{aka.TYPE, byte(aka.SubtypeChallenge), 0, 0})
	atRandOffset = len(p) + aka.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		aka.AT_RAND, append(
			[]byte{0, 0}, // reserved
			make([]byte, aka.RAND_LEN)...)))
	if err != nil {
		panic(err)
	}
	atAutnOffset = len(p) + aka.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		aka.AT_AUTN, append(
			[]byte{0, 0}, // reserved
			make([]byte, aka.AUTN_LEN)...)))
	if err != nil {
		panic(err)
	}
	atMacOffset = len(p) + aka.ATT_HDR_LEN
	p, err = p.Append(eap.NewAttribute(
		aka.AT_MAC, append(
			[]byte{0, 0}, // reserved
			make([]byte, aka.MAC_LEN)...)))
	if err != nil {
		panic(err)
	}
	challengeReqTemplateLen = len(p)
	challengeReqTemplate = p
}

type tgppAuthResult struct {
	rand, autn, xres, ck, ik []byte
	sid                      string
	profile                  *swx_protos.AuthenticationAnswer_UserProfile
}

func getSwxVector(s *servicers.EapAkaSrv, imsi string, resyncInfo []byte) (*tgppAuthResult, error) {
	metrics.SwxRequests.Inc()
	swxStartTime := time.Now()

	ans, err := swx_proxy.Authenticate(
		&swx_protos.AuthenticationRequest{
			UserName:             imsi,
			SipNumAuthVectors:    1,
			AuthenticationScheme: swx_protos.AuthenticationScheme_EAP_AKA,
			ResyncInfo:           resyncInfo,
			RetrieveUserProfile:  true,
		})

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
	av := ans.SipAuthVectors[0] // Use first vector for now
	ra := av.GetRandAutn()
	if len(ra) < aka.RandAutnLen {
		return nil, status.Errorf(codes.Internal,
			"Invalid SWx RandAutn len (%d, expected: %d) in Response: %+v", len(ra), aka.RandAutnLen, *ans)
	}
	return &tgppAuthResult{
		rand:    ra[:aka.RAND_LEN],
		autn:    ra[aka.RAND_LEN:aka.RandAutnLen],
		xres:    av.GetXres(),
		ck:      av.GetConfidentialityKey(),
		ik:      av.GetIntegrityKey(),
		sid:     ans.GetSessionId(),
		profile: ans.GetUserProfile(),
	}, nil
}

func createChallengeRequest(
	s *servicers.EapAkaSrv,
	lockedCtx *servicers.UserCtx,
	identifier uint8,
	resyncInfo []byte) (eap.Packet, error) {

	var (
		authRes *tgppAuthResult
		err     error
	)
	if s.UseS6a() {
		authRes, err = getS6aVector(s, string(lockedCtx.Imsi), resyncInfo)
	} else {
		authRes, err = getSwxVector(s, string(lockedCtx.Imsi), resyncInfo)
	}
	if err != nil {
		if err == GERANOnlyVectorsErr { // User only has GERAN Vectors (2G) try to steer UE toward EAP-SIM
			glog.Warningf("GERAN only Vectors for IMSI: %s, will try EAP-SIM ID Request", lockedCtx.Imsi)
			lockedCtx.SetState(aka.StateRedirected)
			return sim.NewStartReq(identifier+1, sim.AT_PERMANENT_ID_REQ), nil
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
		glog.Errorf("AKA RPC [%s] %s", code, msg)
		return aka.NewAKANotificationReq(identifier, aka.NOTIFICATION_FAILURE), nil
	}
	identifier++

	lockedCtx.Identifier = identifier
	lockedCtx.Rand = authRes.rand
	lockedCtx.Xres = authRes.xres
	lockedCtx.AuthSessionId = authRes.sid
	lockedCtx.Profile = authRes.profile

	// Clone EAP Challenge packet
	p := eap.Packet(make([]byte, challengeReqTemplateLen))
	copy(p, challengeReqTemplate)

	// Set current identifier
	p[eap.EapMsgIdentifier] = identifier

	// Set AT_RAND
	copy(p[atRandOffset:], lockedCtx.Rand)

	// Set AT_AUTN
	copy(p[atAutnOffset:], authRes.autn)

	// Calculate AT_MAC
	_, lockedCtx.K_aut, lockedCtx.MSK, _ = aka.MakeAKAKeys([]byte(lockedCtx.Identity), authRes.ik, authRes.ck)
	mac := aka.GenMac(p, lockedCtx.K_aut)
	// Set AT_MAC
	copy(p[atMacOffset:], mac)
	return p, nil
}
