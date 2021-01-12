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

package swx_proxy_test

import (
	"context"
	"strconv"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/swx_proxy"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/swx_proxy/servicers/test"
	"magma/feg/gateway/services/swx_proxy/test_init"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestSwxProxyClient_VerifyAuthorization(t *testing.T) {
	err := test_init.InitTestMconfig(t, "127.0.0.1:0", true)
	assert.NoError(t, err)
	srv, err := test_init.StartTestService(t)
	if err != nil {
		t.Fatal(err)
	}
	standardSwxProxyTest(t)
	_, err = srv.StopService(context.Background(), &orcprotos.Void{})
	assert.NoError(t, err)
}

func TestSwxProxyClient_VerifyAuthorizationOff(t *testing.T) {
	err := test_init.InitTestMconfig(t, "127.0.0.1:0", false)
	assert.NoError(t, err)
	srv, err := test_init.StartTestService(t)
	if err != nil {
		t.Fatal(err)
	}
	standardSwxProxyTest(t)
	_, err = srv.StopService(context.Background(), &orcprotos.Void{})
	assert.NoError(t, err)
}

func standardSwxProxyTest(t *testing.T) {
	expectedUsername := test.BASE_IMSI
	numVectors := 1
	expectedAuthScheme := protos.AuthenticationScheme_EAP_AKA
	authReq := &protos.AuthenticationRequest{
		UserName:             expectedUsername,
		SipNumAuthVectors:    uint32(numVectors),
		AuthenticationScheme: expectedAuthScheme,
	}

	// Authentication Request - MAR
	// with cache numVectors will be ignored & the proxy will always ask for MinRequestedVectors
	// and always will return 1 vector
	for i := uint32(0); i < servicers.MinRequestedVectors; i++ {
		authRes, err := swx_proxy.Authenticate(authReq)
		if err != nil {
			t.Fatalf("GRPC MAR Error: %v", err)
			return
		}
		t.Logf("GRPC MAA: %#+v", *authRes)
		assert.Equal(t, expectedUsername, authRes.GetUserName())
		assert.Equal(t, 1, len(authRes.GetSipAuthVectors()))
		v := authRes.SipAuthVectors[0]
		assert.Equal(t, protos.AuthenticationScheme_EAP_AKA, v.GetAuthenticationScheme())
		assert.Equal(t, []byte(test.DefaultSIPAuthenticate+strconv.Itoa(int(i+14))), v.GetRandAutn())
		assert.Equal(t, []byte(test.DefaultSIPAuthorization), v.GetXres())
		assert.Equal(t, []byte(test.DefaultCK), v.GetConfidentialityKey())
		assert.Equal(t, []byte(test.DefaultIK), v.GetIntegrityKey())
	}

	// Registration request - SAR
	regReq := &protos.RegistrationRequest{
		UserName: expectedUsername,
	}
	regRes, err := swx_proxy.Register(regReq)
	if err != nil {
		t.Fatalf("GRPC SAR Register Error: %v", err)
		return
	}
	assert.Equal(t, &protos.RegistrationAnswer{SessionId: regRes.GetSessionId()}, regRes)
	t.Logf("GRPC Register SAA: %#+v", *regRes)

	regReq.SessionId = regRes.GetSessionId()
	deregRes, err := swx_proxy.Deregister(regReq)
	if err != nil {
		t.Fatalf("GRPC SAR Deregister Error: %v", err)
		return
	}
	assert.Equal(t, &protos.RegistrationAnswer{SessionId: regRes.GetSessionId()}, deregRes)
	t.Logf("GRPC Deregister SAA: %#+v", *deregRes)

	// Test client error handling
	authRes, err := swx_proxy.Authenticate(nil)
	assert.EqualError(t, err, "Invalid AuthenticationRequest provided: request is nil")
	assert.Nil(t, authRes)

	regRes, err = swx_proxy.Register(nil)
	assert.EqualError(t, err, "Invalid RegistrationRequest provided: request is nil")
	assert.Nil(t, regRes)
}
