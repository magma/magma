/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gx_test

import (
	"errors"
	"fmt"
	"testing"

	policydb_mocks "magma/feg/gateway/policydb/mocks"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	relay_mocks "magma/feg/gateway/services/session_proxy/relay/mocks"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests data flow from feg to/from SessionProxyResponder. Protobuf conversion
// is tested separately.
func TestReAuthRelay(t *testing.T) {
	sm, cloudRegistry := relay_mocks.StartMockSessionProxyResponder(t)
	mockPolicyClient := &policydb_mocks.PolicyDBClient{}
	handler := gx.GetGxReAuthHandler(cloudRegistry, mockPolicyClient)

	imsi := "IMSI000000000000001"
	sessionID := fmt.Sprintf("%s-%d", imsi, 1234)
	// We're not putting any rules in the request, those code paths are covered
	// by model conversion tests
	mockPolicyClient.On("GetRuleIDsForBaseNames", mock.Anything).Return([]string{})

	// Happy path
	req := &gx.ReAuthRequest{SessionID: sessionID}
	sm.On("PolicyReAuth", mock.Anything, &protos.PolicyReAuthRequest{SessionId: sessionID, Imsi: imsi}).
		Return(&protos.PolicyReAuthAnswer{SessionId: "mock_ret"}, nil).Once()
	actual := handler(req)
	// Handler should use session ID from request instead of response
	assert.Equal(t, &gx.ReAuthAnswer{SessionID: sessionID, ResultCode: diam.Success, RuleReports: []*gx.ChargingRuleReport{}}, actual)

	// Bad session ID
	req = &gx.ReAuthRequest{SessionID: "bad"}
	actual = handler(req)
	assert.Equal(t, &gx.ReAuthAnswer{SessionID: "bad", ResultCode: diam.UnknownSessionID}, actual)

	// Error from client
	req = &gx.ReAuthRequest{SessionID: sessionID}
	sm.On("PolicyReAuth", mock.Anything, &protos.PolicyReAuthRequest{SessionId: sessionID, Imsi: imsi}).
		Return(nil, errors.New("oops")).Once()
	actual = handler(req)
	assert.Equal(t, &gx.ReAuthAnswer{SessionID: sessionID, ResultCode: diam.UnableToDeliver}, actual)
	sm.AssertExpectations(t)
}
