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

	s6a_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap/providers/aka/metrics"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
	"magma/feg/gateway/services/s6a_proxy"
	"magma/feg/gateway/tgpp"
)

const DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION codes.Code = 5420

var GERANOnlyVectorsErr = status.Errorf(codes.Internal, "Only GERAN vectors available")

func getS6aVector(s *servicers.EapAkaSrv, imsi string, resyncInfo []byte) (*tgppAuthResult, error) {
	metrics.S6aRequests.Inc()
	s6aStartTime := time.Now()
	air := &s6a_protos.AuthenticationInformationRequest{
		UserName:                      imsi,
		VisitedPlmn:                   tgpp.GetPlmnID(imsi, s.MncLen()),
		ImmediateResponsePreferred:    true,
		NumRequestedUtranGeranVectors: 1,
		UtranGeranResyncInfo:          resyncInfo,
	}
	ans, err := s6a_proxy.AuthenticationInformation(air)
checkAIError:
	metrics.S6aLatency.Observe(time.Since(s6aStartTime).Seconds())
	if err != nil {
		metrics.S6aFailures.Inc()
		errCode := codes.Internal
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			errCode = se.GRPCStatus().Code()
			if errCode == DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION && air.NumRequestedUtranGeranVectors > 0 {
				glog.Warningf("No UTRAN/GERAN profile for IMSI: %s, error: %v", imsi, err)
				// No UTRAN/GERAN profile, try E-UTRAN
				air.NumRequestedEutranVectors = 1
				air.ResyncInfo = resyncInfo
				air.UtranGeranResyncInfo = nil
				air.NumRequestedUtranGeranVectors = 0

				metrics.S6aRequests.Inc()
				s6aStartTime = time.Now()

				ans, err = s6a_proxy.AuthenticationInformation(air)
				goto checkAIError
			}
		}
		return nil, status.Errorf(errCode, "AI Failure: %v; IMSI: %s", err, imsi)
	}
	if ans == nil {
		return nil, status.Error(codes.Internal, "AI Error: Nil S6a Response")
	}
	if len(ans.UtranVectors) == 0 && len(ans.EutranVectors) == 0 {
		if len(ans.GeranVectors) > 0 {
			return nil, GERANOnlyVectorsErr
		}
		return nil, status.Errorf(codes.Internal, "Error: Missing/empty S6a Auth Vector: %s", ans.String())
	}
	metrics.S6aULRequests.Inc()
	s6aStartTime = time.Now()
	ula, err := s6a_proxy.UpdateLocation(
		&s6a_protos.UpdateLocationRequest{
			UserName:    imsi,
			VisitedPlmn: air.VisitedPlmn,
		})
	metrics.S6aULLatency.Observe(time.Since(s6aStartTime).Seconds())
	if err != nil {
		metrics.S6aULFailures.Inc()
		errCode := codes.Internal
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			errCode = se.GRPCStatus().Code()
		}
		return nil, status.Errorf(errCode, "UL Failure: %v; IMSI: %s", err, imsi)
	}
	if ula == nil {
		return nil, status.Error(codes.Internal, "UI Error: Nil S6a Response")
	}
	if len(ans.UtranVectors) > 0 { // preferable auth method - UTRAN Vector
		v := ans.UtranVectors[0]
		return &tgppAuthResult{
			rand: v.GetRand(),
			autn: v.GetAutn(),
			xres: v.GetXres(),
			ck:   v.GetConfidentialityKey(),
			ik:   v.GetIntegrityKey(),
			profile: &s6a_protos.AuthenticationAnswer_UserProfile{
				Msisdn: tgpp.DecodeMsisdn(ula.GetMsisdn()),
			},
		}, nil
	} else {
		v := ans.EutranVectors[0]
		ck, ik := []byte{}, []byte{}
		if len(v.GetKasme()) > 128 {
			ck, ik = v.GetKasme()[:128], v.GetKasme()[128:] // 3GPP TS 33.401 A.8
		}
		return &tgppAuthResult{
			rand: v.GetRand(),
			autn: v.GetAutn(),
			xres: v.GetXres(),
			ck:   ck,
			ik:   ik,
			profile: &s6a_protos.AuthenticationAnswer_UserProfile{
				Msisdn: tgpp.DecodeMsisdn(ula.GetMsisdn()),
			},
		}, nil
	}
}
