/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

const (
	MON_KEY1  = "mon_key_1"
	MON_KEY2  = "mon_key_2"
	POLICY_ID = "12345"
)

var (
	UsageTx1    uint64 = 3000000
	UsageRx1    uint64 = 7000000
	UsageTotal1 uint64 = UsageTx1 + UsageRx1
	UsageTx2    uint64 = 8000000
	UsageRx2    uint64 = 14000000
	UsageTotal2 uint64 = UsageTx2 + UsageRx2
)

func TestTerminateSession(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	expectedArg := n7_sbi.PostSmPoliciesSmPolicyIdDeleteJSONRequestBody{
		AccuUsageReports: &[]n7_sbi.AccuUsageReport{
			{
				RefUmIds:         MON_KEY1,
				VolUsageUplink:   n7.GetSbiVolume(UsageTx1),
				VolUsageDownlink: n7.GetSbiVolume(UsageRx1),
				VolUsage:         n7.GetSbiVolume(UsageTotal1),
			},
			{
				RefUmIds:         MON_KEY2,
				VolUsageUplink:   n7.GetSbiVolume(UsageTx2),
				VolUsageDownlink: n7.GetSbiVolume(UsageRx2),
				VolUsage:         n7.GetSbiVolume(UsageTotal2),
			},
		},
	}

	mockN7.On("PostSmPoliciesSmPolicyIdDeleteWithResponse", mock.Anything, POLICY_ID, expectedArg).
		Return(&n7_sbi.PostSmPoliciesSmPolicyIdDeleteResponse{
			HTTPResponse: &http.Response{StatusCode: 204},
		}, nil).Once()

	termSessProto := defaultTerminateSessionRequest(IMSI1)
	response, err := srv.TerminateSession(context.Background(), termSessProto)
	require.NoError(t, err)
	mockN7.AssertExpectations(t)
	assert.Equal(t, defaultTerminateSessionResponse(IMSI1), response)
}

func TestTerminateSessionTimeout(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesSmPolicyIdDeleteWithResponse", mock.Anything, POLICY_ID, mock.Anything).
		Return(nil, &url.Error{Err: context.DeadlineExceeded}).Once()

	termSessProto := defaultTerminateSessionRequest(IMSI1)
	response, err := srv.TerminateSession(context.Background(), termSessProto)
	require.Error(t, err)
	mockN7.AssertExpectations(t)
	assert.Nil(t, response)
}

func TestTerminateSessionErrResp(t *testing.T) {
	srv, _, mockN7 := createCentralSessionControllerForTest(t, false)
	defer srv.Close()

	mockN7.On("PostSmPoliciesSmPolicyIdDeleteWithResponse", mock.Anything, POLICY_ID, mock.Anything).
		Return(&n7_sbi.PostSmPoliciesSmPolicyIdDeleteResponse{
			HTTPResponse: &http.Response{StatusCode: 404},
		}, nil).Once()

	termSessProto := defaultTerminateSessionRequest(IMSI1)
	response, err := srv.TerminateSession(context.Background(), termSessProto)
	require.Error(t, err)
	mockN7.AssertExpectations(t)
	assert.Nil(t, response)
}

func defaultTerminateSessionRequest(imsi string) *protos.SessionTerminateRequest {
	return &protos.SessionTerminateRequest{
		SessionId: SESS_ID1,
		CommonContext: &protos.CommonSessionContext{
			Sid:     &protos.SubscriberID{Id: IMSI1},
			RatType: protos.RATType_TGPP_NR,
			UeIpv4:  UE_IPV4,
		},
		TgppCtx: &protos.TgppContext{GxDestHost: SmPolicyUrl},
		MonitorUsages: []*protos.UsageMonitorUpdate{
			{
				MonitoringKey: []byte(MON_KEY1),
				BytesTx:       UsageTx1,
				BytesRx:       UsageRx1,
			},
			{
				MonitoringKey: []byte(MON_KEY2),
				BytesTx:       UsageTx2,
				BytesRx:       UsageRx2,
			},
		},
	}
}

func defaultTerminateSessionResponse(imsi string) *protos.SessionTerminateResponse {
	return &protos.SessionTerminateResponse{
		Sid:       imsi,
		SessionId: SESS_ID1,
	}
}
