/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package servicers implements Swx GRPC proxy service which sends MAR/SAR messages over
// diameter connection, waits (blocks) for diameter's MAA/SAAs returns their RPC representation
package servicers

import (
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"
	"github.com/fiorix/go-diameter/diam/sm"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	TIMEOUT_SECONDS  = 10
	MAX_DIAM_RETRIES = 1
)

type swxProxy struct {
	clientCfg      *diameter.DiameterClientConfig
	serverCfg      *diameter.DiameterServerConfig
	smClient       *sm.Client
	connMan        *diameter.ConnectionManager
	requestTracker *diameter.RequestTracker
	originStateID  uint32
}

func NewSwxProxy(
	clientCfg *diameter.DiameterClientConfig,
	serverCfg *diameter.DiameterServerConfig,
) (*swxProxy, error) {
	err := clientCfg.Validate()
	if err != nil {
		return nil, err
	}
	clientCfg = clientCfg.FillInDefaults()

	err = serverCfg.Validate()
	if err != nil {
		return nil, err
	}
	originStateID := uint32(time.Now().Unix())

	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(clientCfg.Host),
		OriginRealm:      datatype.DiameterIdentity(clientCfg.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(clientCfg.ProductName),
		OriginStateID:    datatype.Unsigned32(originStateID),
		FirmwareRevision: 1,
	})

	mux.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		if m != nil {
			glog.Infof("Unhandled SWx message: %s", m)
		}
	}) // Catch all.

	if clientCfg.WatchdogInterval == 0 {
		clientCfg.WatchdogInterval = diameter.DefaultWatchdogIntervalSeconds
	}

	smClient := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     clientCfg.Retransmits,
		RetransmitInterval: time.Second,
		EnableWatchdog:     clientCfg.WatchdogInterval > 0,
		WatchdogInterval:   time.Second * time.Duration(clientCfg.WatchdogInterval),
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_SWX_APP_ID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
				},
			}),
		},
	}

	connMan := diameter.NewConnectionManager()
	// create connection in connection map
	connMan.GetConnection(smClient, serverCfg)

	proxy := &swxProxy{
		clientCfg:      clientCfg,
		serverCfg:      serverCfg,
		smClient:       smClient,
		connMan:        connMan,
		requestTracker: diameter.NewRequestTracker(),
		originStateID:  originStateID,
	}
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.MultimediaAuthentication, Request: false},
		handleMAA(proxy))
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.TGPP_SWX_APP_ID, Code: diam.ServerAssignment, Request: false},
		handleSAA(proxy))

	return proxy, nil
}

// SwxProxyServer implementation
//
// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s *swxProxy) Authenticate(
	ctx context.Context,
	req *protos.AuthenticationRequest,
) (*protos.AuthenticationAnswer, error) {
	return s.AuthenticateImpl(req)
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s *swxProxy) Register(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	return s.RegisterImpl(req)
}

func (s *swxProxy) genSID() string {
	return s.clientCfg.GenSessionID("swx")
}
