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
	testMARSuccessful(t, true)
	testMARSuccessful(t, false)
}

func testMARSuccessful(t *testing.T, verifyAuthorization bool) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, verifyAuthorization, false)
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

	swxProxy := getTestSwxProxy(t, hss, true, true)
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
	swxProxy := getTestSwxProxy(t, hss, false, true)
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
	swxProxy := getTestSwxProxy(t, hss, false, true)
	sar := &fegprotos.RegistrationRequest{UserName: "sub1"}
	_, err := swxProxy.Register(context.Background(), sar)
	assert.NoError(t, err)
}

func TestSAR_UnknownIMSI(t *testing.T) {
	hss := getTestHSSDiameterServer(t)
	swxProxy := getTestSwxProxy(t, hss, false, true)
	sar := &fegprotos.RegistrationRequest{UserName: "sub_unknown"}
	_, err := swxProxy.Register(context.Background(), sar)
	assert.EqualError(t, err, "rpc error: code = Code(5001) desc = Diameter Error: 5001 (USER_UNKNOWN)")
}

// getTestSwxProxy creates a SWx Proxy server and test HSS Diameter
// server which are configured to communicate with each other.
func getTestSwxProxy(t *testing.T, hss *hss.HomeSubscriberServer, verifyAuthr, wCache bool) fegprotos.SwxProxyServer {
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
	assert.NoError(t, err)
	return swxProxy
}
