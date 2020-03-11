/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mock_pcrf

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/golang/glog"
)

// PCRFConfig defines the configuration for a PCRF server, which for now is just
// the host ip
type PCRFConfig struct {
	ServerConfig *diameter.DiameterServerConfig
}

type creditByMkey map[string]*protos.UsageMonitor

type subscriberAccount struct {
	RuleNames       []string
	RuleBaseNames   []string
	RuleDefinitions []*protos.RuleDefinition
	UsageMonitors   creditByMkey
}

// PCRFDiamServer wraps an PCRF storing subscribers and their rules
type PCRFDiamServer struct {
	diameterSettings *diameter.DiameterClientConfig
	pcrfConfig       *PCRFConfig
	serviceConfig    *protos.PCRFConfigs
	subscribers      map[string]*subscriberAccount // map of imsi to to rules
}

// NewPCRFDiamServer initializes an PCRF with an empty rule map
// Input: *sm.Settings containing the diameter related parameters
//				*TestPCRFConfig containing the server address
//
// Output: a new PCRFDiamServer
func NewPCRFDiamServer(
	diameterSettings *diameter.DiameterClientConfig,
	pcrfConfig *PCRFConfig,
) *PCRFDiamServer {
	return &PCRFDiamServer{
		diameterSettings: diameterSettings,
		pcrfConfig:       pcrfConfig,
		subscribers:      map[string]*subscriberAccount{},
	}
}

// Reset is an GRPC procedure which configure the server to the default status.
// It will be called from the gateway.
func (srv *PCRFDiamServer) Reset(
	ctx context.Context,
	req *orcprotos.Void,
) (*orcprotos.Void, error) {
	return nil, errors.New("Not implemented")
}

// ConfigServer is an GRPC procedure which configure the server to respond
// to requests. It will be called from the gateway
func (srv *PCRFDiamServer) ConfigServer(
	ctx context.Context,
	config *protos.ServerConfiguration,
) (*orcprotos.Void, error) {
	return nil, errors.New("Not implemented")
}

// Start begins the server and blocks, listening to the network
// Output: error if the server could not be started
func (srv *PCRFDiamServer) Start(lis net.Listener) error {
	mux := sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(srv.diameterSettings.Host),
		OriginRealm:      datatype.DiameterIdentity(srv.diameterSettings.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(srv.diameterSettings.ProductName),
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	})
	mux.HandleIdx(
		diam.CommandIndex{AppID: diam.GX_CHARGING_CONTROL_APP_ID, Code: diam.CreditControl, Request: true},
		getCCRHandler(srv))
	go logErrors(mux.ErrorReports())
	serverConfig := srv.pcrfConfig.ServerConfig
	server := &diam.Server{
		Network: serverConfig.Protocol,
		Addr:    serverConfig.Addr,
		Handler: mux,
		Dict:    nil,
	}
	return server.Serve(lis)
}

// StartListener starts a listener based on ServerConfig
// If ServerConfig did not have valid values, default values would be used
func (srv *PCRFDiamServer) StartListener() (net.Listener, error) {
	serverConfig := srv.pcrfConfig.ServerConfig

	network := serverConfig.Protocol
	if len(network) == 0 {
		network = "tcp"
	}
	addr := serverConfig.Addr
	if len(addr) == 0 {
		addr = ":3870"
	}
	l, e := diam.Listen(network, addr)
	if e != nil {
		return nil, e
	}
	return l, nil
}

// logErrors logs errors received during transmission
func logErrors(ec <-chan *diam.ErrorReport) {
	for err := range ec {
		glog.Errorf("PCRF transmit error: %s", err)
	}
}

func (srv *PCRFDiamServer) SetPCRFConfigs(
	ctx context.Context,
	configs *protos.PCRFConfigs,
) (*orcprotos.Void, error) {
	srv.serviceConfig = configs
	return &orcprotos.Void{}, nil
}

// NewSubscriber adds a subscriber to the PCRF to be tracked
// Input: string containing the subscriber IMSI (can be in any form)
func (srv *PCRFDiamServer) CreateAccount(
	ctx context.Context,
	subscriberID *lteprotos.SubscriberID,
) (*orcprotos.Void, error) {
	srv.subscribers[subscriberID.Id] = &subscriberAccount{
		RuleNames:     []string{},
		UsageMonitors: make(creditByMkey),
	}
	glog.V(2).Infof("New account %s added", subscriberID.Id)
	return &orcprotos.Void{}, nil
}

// SetRules sets or overrides the rules applicable to the subscriber
// Input: imsi string IMSI for the subscriber
//			  ruleNames []string containing all rule names to apply
//			  ruleBaseNames []string containing all rule base names to apply
//			  ruleDefinitions []*RuleDefinition containing all dynamic rules to apply
// Output: error if subscriber could not be found
func (srv *PCRFDiamServer) SetRules(
	ctx context.Context,
	accountRules *protos.AccountRules,
) (*orcprotos.Void, error) {
	account, ok := srv.subscribers[accountRules.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", accountRules.Imsi)
	}
	account.RuleNames = accountRules.StaticRuleNames
	account.RuleBaseNames = accountRules.StaticRuleBaseNames
	account.RuleDefinitions = accountRules.DynamicRuleDefinitions
	return &orcprotos.Void{}, nil
}

func (srv *PCRFDiamServer) SetUsageMonitors(
	ctx context.Context,
	usageMonitorInfo *protos.UsageMonitorConfiguration,
) (*orcprotos.Void, error) {
	account, ok := srv.subscribers[usageMonitorInfo.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", usageMonitorInfo.Imsi)
	}
	account.UsageMonitors = make(creditByMkey)
	for _, monitor := range usageMonitorInfo.UsageMonitorCredits {
		account.UsageMonitors[string(monitor.MonitorInfoPerRequest.MonitoringKey)] = monitor
	}
	return &orcprotos.Void{}, nil
}

// GetRuleNames returns all the rules set for a subscriber
// Input: string IMSI for the subscriber
// Output: []string containing all applicable rules
//			   error if subscriber could not be found
func (srv *PCRFDiamServer) GetRuleNames(imsi string) ([]string, error) {
	account, ok := srv.subscribers[imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", imsi)
	}
	return account.RuleNames, nil
}

// GetRuleBaseNames returns all the rule base names set for a subscriber
// Input: string IMSI for the subscriber
// Output: []string containing all applicable rule base names
//			   error if subscriber could not be found
func (srv *PCRFDiamServer) GetRuleBaseNames(imsi string) ([]string, error) {
	account, ok := srv.subscribers[imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", imsi)
	}
	return account.RuleBaseNames, nil
}

// GetRuleDefinitions returns all the dynamic rule definitions set for a subscriber
// Input: string IMSI for the subscriber
// Output: []*RuleDefinition containing all applicable rules
//			   error if subscriber could not be found
func (srv *PCRFDiamServer) GetRuleDefinitions(
	imsi string,
) ([]*protos.RuleDefinition, error) {
	account, ok := srv.subscribers[imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", imsi)
	}
	return account.RuleDefinitions, nil
}

// Reset eliminates all the subscribers allocated for the system.
func (srv *PCRFDiamServer) ClearSubscribers(ctx context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	srv.subscribers = map[string]*subscriberAccount{}
	glog.V(2).Info("All accounts deleted.")
	return &orcprotos.Void{}, nil
}

func (srv *PCRFDiamServer) SetExpectations(ctx context.Context, expectations *protos.GxCreditControlExpectations) (*orcprotos.Void, error) {
	return &orcprotos.Void{}, nil
}

func (srv *PCRFDiamServer) AssertExpectations(ctx context.Context, void *orcprotos.Void) (*protos.GxCreditControlResult, error) {
	return &protos.GxCreditControlResult{}, nil
}
