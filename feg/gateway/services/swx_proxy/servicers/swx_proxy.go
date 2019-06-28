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
	"magma/feg/gateway/services/swx_proxy/cache"
	"magma/feg/gateway/services/swx_proxy/metrics"

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
	config         *SwxProxyConfig
	smClient       *sm.Client
	connMan        *diameter.ConnectionManager
	requestTracker *diameter.RequestTracker
	originStateID  uint32
	cache          *cache.Impl
}

type SwxProxyConfig struct {
	ClientCfg           *diameter.DiameterClientConfig
	ServerCfg           *diameter.DiameterServerConfig
	VerifyAuthorization bool // should we verify non-3gpp IP access is enabled for user
	CacheTTLSeconds     uint32
}

// NewSwxProxy creates a new instance of the proxy with configured cache TTL
func NewSwxProxy(config *SwxProxyConfig) (*swxProxy, error) {
	if config.CacheTTLSeconds < uint32(cache.DefaultGcInterval.Seconds()) {
		config.CacheTTLSeconds = uint32(cache.DefaultTtl.Seconds())
	}
	cch, _ := cache.NewExt(cache.DefaultGcInterval, time.Second*time.Duration(config.CacheTTLSeconds))
	return NewSwxProxyWithCache(config, cch)
}

// NewSwxProxyWithCache creates a new instance of the proxy with given cache implementation
func NewSwxProxyWithCache(config *SwxProxyConfig, cache *cache.Impl) (*swxProxy, error) {
	err := ValidateSwxProxyConfig(config)
	if err != nil {
		return nil, err
	}
	config.ClientCfg = config.ClientCfg.FillInDefaults()

	originStateID := uint32(time.Now().Unix())

	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(config.ClientCfg.Host),
		OriginRealm:      datatype.DiameterIdentity(config.ClientCfg.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(config.ClientCfg.ProductName),
		OriginStateID:    datatype.Unsigned32(originStateID),
		FirmwareRevision: 1,
	})

	mux.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		if m != nil {
			glog.Infof("Unhandled SWx message: %s", m)
		}
	}) // Catch all.

	if config.ClientCfg.WatchdogInterval == 0 {
		config.ClientCfg.WatchdogInterval = diameter.DefaultWatchdogIntervalSeconds
	}

	smClient := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     config.ClientCfg.Retransmits,
		RetransmitInterval: time.Second,
		EnableWatchdog:     config.ClientCfg.WatchdogInterval > 0,
		WatchdogInterval:   time.Second * time.Duration(config.ClientCfg.WatchdogInterval),
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
	connMan.GetConnection(smClient, config.ServerCfg)

	proxy := &swxProxy{
		config:         config,
		smClient:       smClient,
		connMan:        connMan,
		requestTracker: diameter.NewRequestTracker(),
		originStateID:  originStateID,
		cache:          cache,
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
	authStartTime := time.Now()
	res, err := s.AuthenticateImpl(req)
	if err == nil {
		metrics.AuthLatency.Observe(time.Since(authStartTime).Seconds())
	}
	return res, err
}

// Register sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s *swxProxy) Register(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	registerStartTime := time.Now()
	res, err := s.RegisterImpl(req, ServerAssignmentType_REGISTRATION)
	if err == nil {
		metrics.RegisterLatency.Observe(time.Since(registerStartTime).Seconds())
	}
	return res, err
}

// Deregister sends SAR (code 301) over diameter connection,
// waits (blocks) for SAA & returns its RPC representation
func (s *swxProxy) Deregister(
	ctx context.Context,
	req *protos.RegistrationRequest,
) (*protos.RegistrationAnswer, error) {
	deregisterStartTime := time.Now()
	res, err := s.RegisterImpl(req, ServerAssignnmentType_USER_DEREGISTRATION)
	if err == nil {
		metrics.DeregisterLatency.Observe(time.Since(deregisterStartTime).Seconds())
	}
	return res, err
}

func (s *swxProxy) genSID() string {
	return s.config.ClientCfg.GenSessionID("swx")
}
