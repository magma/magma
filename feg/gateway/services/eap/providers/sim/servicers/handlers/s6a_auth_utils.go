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

	s6a_protos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/eap/providers/sim"
	"magma/feg/gateway/services/eap/providers/sim/metrics"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	"magma/feg/gateway/services/s6a_proxy"
	"magma/feg/gateway/tgpp"
)

const DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION codes.Code = 5420

// EUTRANOnlyVectorsErr - this error indicates that the user only has EUTRAN (4G) profile
// and may be able to authenticate with EAP-AKA/AKA'
var EUTRANOnlyVectorsErr = status.Errorf(codes.Internal, "Only GERAN vectors available")

func getS6aVectors(s *servicers.EapSimSrv, imsi string) (*tgppAuthResult, error) {
	var (
		err error
		ans *s6a_protos.AuthenticationInformationAnswer
		res tgppAuthResult
	)
	visitedPlmn := tgpp.GetPlmnID(imsi, s.MncLen())
	for vlen := 0; vlen < sim.GsmTripletsNumber; {
		metrics.S6aRequests.Inc()
		s6aStartTime := time.Now()
		air := &s6a_protos.AuthenticationInformationRequest{
			UserName:                      imsi,
			VisitedPlmn:                   visitedPlmn,
			ImmediateResponsePreferred:    true,
			NumRequestedUtranGeranVectors: sim.GsmTripletsNumber - uint32(vlen),
		}
		ans, err = s6a_proxy.AuthenticationInformation(air)
		metrics.S6aLatency.Observe(time.Since(s6aStartTime).Seconds())
		if err != nil {
			metrics.S6aFailures.Inc()
			errCode := codes.Internal
			if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
				errCode = se.GRPCStatus().Code()
				if errCode == DIAMETER_ERROR_UNKNOWN_EPS_SUBSCRIPTION {
					// No UTRAN/GERAN profile
					return nil, EUTRANOnlyVectorsErr
				}
			}
			return nil, status.Errorf(errCode, "AI Failure: %v; IMSI: %s", err, imsi)
		}
		if ans == nil {
			return nil, status.Error(codes.Internal, "AI Error: Nil S6a Response")
		}
		// first check for GERAN vectors, they are direct mapping for EAP-SIM
		if len(ans.GetGeranVectors()) > 0 {
			for _, v := range ans.GetGeranVectors() {
				res.rand[vlen], res.Kc[vlen], res.sres[vlen] = v.GetRand(), v.GetKc(), v.GetSres()
				vlen++
				if vlen >= sim.GsmTripletsNumber {
					break
				}
			}
		} else {
			// no GERAN vectors/profile, try to use UTRAN vectors
			if len(ans.GetUtranVectors()) == 0 {
				return nil, status.Error(codes.Internal, "Invalid S6a Response, no GERAN/UTRAN vectors")
			}
			for _, v := range ans.GetUtranVectors() {
				res.rand[vlen] = v.GetRand()
				res.Kc[vlen], res.sres[vlen] =
					sim.GsmFromUmts1(v.GetConfidentialityKey(), v.GetIntegrityKey(), v.GetXres())
				vlen++
				if vlen >= sim.GsmTripletsNumber {
					break
				}
			}
		}
	}
	metrics.S6aULRequests.Inc()
	s6aStartTime := time.Now()
	ula, err := s6a_proxy.UpdateLocation(
		&s6a_protos.UpdateLocationRequest{
			UserName:    imsi,
			VisitedPlmn: visitedPlmn,
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
	res.profile = &s6a_protos.AuthenticationAnswer_UserProfile{Msisdn: tgpp.DecodeMsisdn(ula.GetMsisdn())}
	return &res, nil
}
