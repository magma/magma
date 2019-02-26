package test

import (
	"context"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/crypto"

	"github.com/stretchr/testify/assert"
)

func TestMAR_Successful(t *testing.T) {
	swxProxy := getTestSwxProxy(t)
	mar := &protos.AuthenticationRequest{
		UserName:             "sub1",
		SipNumAuthVectors:    5,
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
	}
	maa, err := swxProxy.Authenticate(context.Background(), mar)
	assert.NoError(t, err)

	assert.Equal(t, "sub1", maa.GetUserName())
	assert.Equal(t, 5, len(maa.GetSipAuthVectors()))
	for _, vector := range maa.GetSipAuthVectors() {
		assert.Equal(t, protos.AuthenticationScheme_EAP_AKA, vector.AuthenticationScheme)
		assert.Equal(t, crypto.ConfidentialityKeyBytes, len(vector.ConfidentialityKey))
		assert.Equal(t, crypto.IntegrityKeyBytes, len(vector.IntegrityKey))
		assert.Equal(t, crypto.RandChallengeBytes+crypto.AutnBytes, len(vector.RandAutn))
		assert.Equal(t, crypto.XresBytes, len(vector.Xres))
	}
}

func TestMAR_UnknownIMSI(t *testing.T) {
	swxProxy := getTestSwxProxy(t)
	mar := &protos.AuthenticationRequest{
		UserName:             "sub_unknown",
		SipNumAuthVectors:    1,
		AuthenticationScheme: protos.AuthenticationScheme_EAP_AKA,
	}
	maa, err := swxProxy.Authenticate(context.Background(), mar)
	assert.EqualError(t, err, "rpc error: code = Code(5001) desc = Diameter Error: 5001 (USER_UNKNOWN)")
	assert.Equal(t, "sub_unknown", maa.UserName)
	assert.Equal(t, 0, len(maa.SipAuthVectors))
}

// getTestSwxProxy creates a SWx Proxy server and test HSS Diameter
// server which are configured to communicate with each other.
func getTestSwxProxy(t *testing.T) protos.SwxProxyServer {
	hss := getTestHSSDiameterServer(t)
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
	swxProxy, err := servicers.NewSwxProxy(clientCfg, diameterServerCfg)
	assert.NoError(t, err)
	return swxProxy
}
