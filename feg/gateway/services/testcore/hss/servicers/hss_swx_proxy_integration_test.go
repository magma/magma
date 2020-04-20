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

package servicers_test

import (
	"context"
	"testing"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/cache"
	swx "magma/feg/gateway/services/swx_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/lte/cloud/go/crypto"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestMAR_Successful(t *testing.T) {
	testMARSuccessful(t, true, false)
	testMARSuccessful(t, false, false)
	testMARSuccessful(t, true, true)
	testMARSuccessful(t, false, true)
}

func testMARSuccessful(t *testing.T, verifyAuthorization bool, clearAAAserver bool) {
	hss := getTestHSSDiameterServer(t)
	subscriber, err := hss.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	if clearAAAserver {
		subscriber.State.TgppAaaServerName = ""
	}
	_, err = hss.UpdateSubscriber(context.Background(), subscriber)
	assert.NoError(t, err)

	swxProxy := getTestSwxProxy(t, hss, verifyAuthorization, false, true)
	mar := &fegprotos.AuthenticationRequest{
		UserName:             "sub1",
		SipNumAuthVectors:    5,
		AuthenticationScheme: fegprotos.AuthenticationScheme_EAP_AKA,
	}
	maa, err := swxProxy.Authenticate(context.Background(), mar)
	assert.NoError(t, err)

	assert.Equal(t, "sub1", maa.GetUserName())
	assert.Equal(t, 5, len(maa.GetSipAuthVectors()))
	for _, vector := range maa.GetSipAuthVectors() {
		assert.Equal(t, fegprotos.AuthenticationScheme_EAP_AKA, vector.AuthenticationScheme)
		assert.Equal(t, crypto.ConfidentialityKeyBytes, len(vector.ConfidentialityKey))
		assert.Equal(t, crypto.IntegrityKeyBytes, len(vector.IntegrityKey))
		assert.Equal(t, crypto.RandChallengeBytes+crypto.AutnBytes, len(vector.RandAutn))
		assert.Equal(t, crypto.XresBytes, len(vector.Xres))
	}
}

func TestMAR_AuthRejected(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	subscriber, err := hss.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	subscriber.Non_3Gpp.Non_3GppIpAccess = lteprotos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_BARRED
	_, err = hss.UpdateSubscriber(context.Background(), subscriber)
	assert.NoError(t, err)

	swxProxy := getTestSwxProxy(t, hss, true, true, true)
	mar := &fegprotos.AuthenticationRequest{
		UserName:             "sub1",
		SipNumAuthVectors:    5,
		AuthenticationScheme: fegprotos.AuthenticationScheme_EAP_AKA,
	}
	maa, err := swxProxy.Authenticate(context.Background(), mar)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = User sub1 is not authorized for Non-3GPP Subscription Access")
	assert.Equal(t, "sub1", maa.GetUserName())
	assert.Equal(t, 0, len(maa.GetSipAuthVectors()))
}

func TestMAR_UnknownIMSI(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, false, true, true)
	mar := &fegprotos.AuthenticationRequest{
		UserName:             "sub_unknown",
		SipNumAuthVectors:    1,
		AuthenticationScheme: fegprotos.AuthenticationScheme_EAP_AKA,
	}
	maa, err := swxProxy.Authenticate(context.Background(), mar)
	assert.EqualError(t, err, "rpc error: code = Code(5001) desc = Diameter Error: 5001 (USER_UNKNOWN)")
	assert.Equal(t, "sub_unknown", maa.UserName)
	assert.Equal(t, 0, len(maa.SipAuthVectors))
}

func TestSAR_SuccessfulRegistration(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, false, true, true)
	sar := &fegprotos.RegistrationRequest{UserName: "sub1"}
	_, err := swxProxy.Register(context.Background(), sar)
	assert.NoError(t, err)
}

func TestSAR_UnknownIMSI(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, false, true, true)
	sar := &fegprotos.RegistrationRequest{UserName: "sub_unknown"}
	_, err := swxProxy.Register(context.Background(), sar)
	assert.EqualError(t, err, "rpc error: code = Code(5001) desc = Diameter Error: 5001 (USER_UNKNOWN)")
}

func TestRTR_SuccessfulDeregistration(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, false, false, true)
	sar := &fegprotos.RegistrationRequest{
		UserName: "sub1",
	}
	_, err := swxProxy.Register(context.Background(), sar)
	assert.NoError(t, err)

	sub := &lteprotos.SubscriberID{Id: "sub1"}
	_, err = hss.DeregisterSubscriber(context.Background(), sub)
	assert.NoError(t, err)

	subData, err := hss.GetSubscriberData(context.Background(), sub)
	assert.NoError(t, err)
	assert.False(t, subData.GetState().GetTgppAaaServerRegistered())
}

func TestRTR_UnsuccessfulDeregistration(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, false, false, false)
	sar := &fegprotos.RegistrationRequest{
		UserName: "sub1",
	}
	_, err := swxProxy.Register(context.Background(), sar)
	assert.NoError(t, err)

	sub := &lteprotos.SubscriberID{Id: "sub1"}
	_, err = hss.DeregisterSubscriber(context.Background(), sub)
	assert.Error(t, err)

	subData, err := hss.GetSubscriberData(context.Background(), sub)
	assert.NoError(t, err)
	assert.True(t, subData.GetState().GetTgppAaaServerRegistered())
}

func TestRTR_UnknownIMSI(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	_, err := hss.DeregisterSubscriber(context.Background(), &lteprotos.SubscriberID{Id: "sub_unknown"})
	assert.Error(t, err)
}

// getTestSwxProxy creates a SWx Proxy server and test HSS Diameter
// server which are configured to communicate with each other.
func getTestSwxProxy(t *testing.T, hss *hss.HomeSubscriberServer, verifyAuthr, wCache bool, successfulRelay bool) fegprotos.SwxProxyServer {
	serverCfg := hss.Config.Server

	// Create an swx proxy server.
	clientCfg := &diameter.DiameterClientConfig{
		Host:             serverCfg.DestHost,
		Realm:            serverCfg.DestRealm,
		ProductName:      "magma",
		AppID:            0,
		AuthAppID:        0,
		Retransmits:      3,
		WatchdogInterval: 10,
		RetryCount:       3,
	}
	diameterServerCfg := &diameter.DiameterServerConfig{
		DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:      serverCfg.Address,
			Protocol:  serverCfg.Protocol,
			LocalAddr: serverCfg.LocalAddress},
		DestHost:  serverCfg.DestHost,
		DestRealm: serverCfg.DestRealm,
	}
	swxProxyConfig := &swx.SwxProxyConfig{
		ClientCfg:           clientCfg,
		ServerCfg:           diameterServerCfg,
		VerifyAuthorization: verifyAuthr,
	}
	var vc *cache.Impl
	if wCache {
		vc = cache.New()
	}
	swxProxy, err := swx.NewSwxProxyWithCache(swxProxyConfig, vc)
	if successfulRelay {
		swxProxy.Relay = &successfulMockRelay{}
	} else {
		swxProxy.Relay = &unsuccessfulMockRelay{}
	}
	assert.NoError(t, err)
	return swxProxy
}

type successfulMockRelay struct{}

func (s *successfulMockRelay) RelayRTR(*swx.RTR) (fegprotos.ErrorCode, error) {
	return fegprotos.ErrorCode_SUCCESS, nil
}

func (s *successfulMockRelay) RelayASR(*diameter.ASR) (fegprotos.ErrorCode, error) {
	return fegprotos.ErrorCode_SUCCESS, nil
}

type unsuccessfulMockRelay struct{}

func (s *unsuccessfulMockRelay) RelayRTR(*swx.RTR) (fegprotos.ErrorCode, error) {
	return fegprotos.ErrorCode_UNABLE_TO_DELIVER, nil
}

func (s *unsuccessfulMockRelay) RelayASR(*diameter.ASR) (fegprotos.ErrorCode, error) {
	return fegprotos.ErrorCode_UNABLE_TO_DELIVER, nil
}
