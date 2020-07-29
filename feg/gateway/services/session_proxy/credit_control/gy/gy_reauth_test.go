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

package gy_test

import (
	"fmt"
	"testing"

	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/relay/mocks"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReAuthRelay(t *testing.T) {
	sm, cloudRegistry := mocks.StartMockSessionProxyResponder(t)
	handler := gy.GetGyReAuthHandler(cloudRegistry)

	var rg uint32 = 1
	assertReAuth(t, handler, sm, &rg, protos.ReAuthResult_UPDATE_INITIATED, diam.LimitedSuccess)
	assertReAuth(t, handler, sm, nil, protos.ReAuthResult_UPDATE_INITIATED, diam.LimitedSuccess)
	assertReAuth(t, handler, sm, &rg, protos.ReAuthResult_UPDATE_NOT_NEEDED, diam.Success)
	assertReAuth(t, handler, sm, &rg, protos.ReAuthResult_SESSION_NOT_FOUND, diam.UnknownSessionID)
	assertReAuth(t, handler, sm, &rg, protos.ReAuthResult_OTHER_FAILURE, diam.UnableToComply)
}

func assertReAuth(
	t *testing.T,
	handler gy.ChargingReAuthHandler,
	sm *mocks.SessionProxyResponderServer,
	ratingGroup *uint32,
	protoResult protos.ReAuthResult,
	expectedResultCode int,
) {
	imsi := "IMSI000000000000001"
	sessionID := fmt.Sprintf("%s-%d", imsi, 1234)

	var matchFunc interface{}
	if ratingGroup == nil {
		matchFunc = getRAAMatcher(imsi, 0, protos.ChargingReAuthRequest_ENTIRE_SESSION)
	} else {
		matchFunc = getRAAMatcher(imsi, *ratingGroup, protos.ChargingReAuthRequest_SINGLE_SERVICE)
	}

	sm.On("ChargingReAuth", mock.Anything, mock.MatchedBy(matchFunc)).Return(
		&protos.ChargingReAuthAnswer{Result: protoResult},
		nil,
	).Once()
	raa := handler(&gy.ChargingReAuthRequest{SessionID: sessionID, RatingGroup: ratingGroup})
	sm.AssertExpectations(t)
	assert.Equal(t, expectedResultCode, int(raa.ResultCode))
	assert.Equal(t, sessionID, raa.SessionID)
}

func getRAAMatcher(imsi string, ratingGroup uint32, reqType protos.ChargingReAuthRequest_Type) interface{} {
	return func(request *protos.ChargingReAuthRequest) bool {
		return request.Sid == imsi && request.ChargingKey == ratingGroup && request.Type == reqType
	}
}
