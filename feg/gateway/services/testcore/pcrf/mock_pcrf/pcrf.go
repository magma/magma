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
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/golang/glog"
)

// PCRFConfig defines the configuration for a PCRF server, which for now is just
// the host ip
type PCRFConfig struct {
	ServerConfig *diameter.DiameterServerConfig
}

type subscriberAccount struct {
	RuleNames       []string
	RuleBaseNames   []string
	RuleDefinitions []*protos.RuleDefinition
	UsageMonitors   map[string]*protos.UsageMonitorCredit
}

// PCRFDiamServer wraps an PCRF storing subscribers and their rules
type PCRFDiamServer struct {
	diameterSettings *diameter.DiameterClientConfig
	pcrfConfig       *PCRFConfig
	subscribers      map[string]*subscriberAccount // map of imsi to to rules
}

type ccrMessage struct {
	SessionID        datatype.UTF8String       `avp:"Session-Id"`
	OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
	DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
	RequestType      datatype.Enumerated       `avp:"CC-Request-Type"`
	RequestNumber    datatype.Unsigned32       `avp:"CC-Request-Number"`
	SubscriptionIDs  []*subscriptionIDDiam     `avp:"Subscription-Id"`
	IPAddr           datatype.OctetString      `avp:"Framed-IP-Address"`
	UsageMonitors    []*usageMonitorRequestAVP `avp:"Usage-Monitoring-Information"`
}

type subscriptionIDDiam struct {
	IDType credit_control.SubscriptionIDType `avp:"Subscription-Id-Type"`
	IDData string                            `avp:"Subscription-Id-Data"`
}

type usageMonitorRequestAVP struct {
	MonitoringKey   string             `avp:"Monitoring-Key"`
	UsedServiceUnit usedServiceUnitAVP `avp:"Used-Service-Unit"`
	Level           gx.MonitoringLevel `avp:"Usage-Monitoring-Level"`
}

type usedServiceUnitAVP struct {
	InputOctets  uint64 `avp:"CC-Input-Octets"`
	OutputOctets uint64 `avp:"CC-Output-Octets"`
	TotalOctets  uint64 `avp:"CC-Total-Octets"`
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

// NewSubscriber adds a subscriber to the PCRF to be tracked
// Input: string containing the subscriber IMSI (can be in any form)
func (srv *PCRFDiamServer) CreateAccount(
	ctx context.Context,
	subscriberID *lteprotos.SubscriberID,
) (*orcprotos.Void, error) {
	srv.subscribers[subscriberID.Id] = &subscriberAccount{
		RuleNames:     []string{},
		UsageMonitors: make(map[string]*protos.UsageMonitorCredit),
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
	account.RuleNames = accountRules.RuleNames
	account.RuleBaseNames = accountRules.RuleBaseNames
	account.RuleDefinitions = accountRules.RuleDefinitions
	return &orcprotos.Void{}, nil
}

func (srv *PCRFDiamServer) SetUsageMonitors(
	ctx context.Context,
	usageMonitorInfo *protos.UsageMonitorInfo,
) (*orcprotos.Void, error) {
	account, ok := srv.subscribers[usageMonitorInfo.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", usageMonitorInfo.Imsi)
	}
	account.UsageMonitors = make(map[string]*protos.UsageMonitorCredit)
	for _, monitor := range usageMonitorInfo.UsageMonitorCredits {
		account.UsageMonitors[monitor.MonitoringKey] = monitor
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

// getCCRHandler returns a handler to be called when the server receives a CCR
func getCCRHandler(srv *PCRFDiamServer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received CCR from %s\n", c.RemoteAddr())
		var ccr ccrMessage
		if err := m.Unmarshal(&ccr); err != nil {
			glog.Errorf("Failed to unmarshal CCR %s", err)
			return
		}
		imsi, err := getIMSI(ccr)
		if err != nil {
			glog.Errorf("Could not parse CCR: %s", err.Error())
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		account, found := srv.subscribers[imsi]
		if !found {
			glog.Error("IMSI not found in subscribers")
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}

		if credit_control.CreditRequestType(ccr.RequestType) == credit_control.CRTInit {
			ruleInstalls := getRuleInstallAVPs(account.RuleNames, account.RuleBaseNames, account.RuleDefinitions)
			usageMonitors := getInitialUsageMonitoringAVPs(account.UsageMonitors)
			avps := append(ruleInstalls, usageMonitors...)
			sendAnswer(ccr, c, m, diam.Success, avps...)
			return
		}

		returnAVPs, err := getUsageMonitorUpdates(account, ccr.UsageMonitors)
		if err != nil {
			sendAnswer(ccr, c, m, diam.InvalidAVPValue)
			return
		}
		sendAnswer(ccr, c, m, diam.Success, returnAVPs...)
	}
}

func shouldReturnRules(requestType credit_control.CreditRequestType) bool {
	return requestType == credit_control.CRTInit
}

// getIMSI finds the account IMSI in a CCR message
func getIMSI(message ccrMessage) (string, error) {
	for _, subID := range message.SubscriptionIDs {
		if subID.IDType == credit_control.EndUserIMSI {
			return subID.IDData, nil
		}
	}
	return "", errors.New("Could not obtain IMSI from CCR message")
}

func getRuleInstallAVPs(
	ruleNames []string,
	ruleBaseNames []string,
	ruleDefs []*protos.RuleDefinition,
) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(ruleNames)+len(ruleBaseNames)+len(ruleDefs))
	for _, rule := range ruleNames {
		avps = append(avps, diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule)),
			},
		}))
	}

	for _, rule := range ruleBaseNames {
		avps = append(avps, diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.ChargingRuleBaseName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(rule)),
			},
		}))
	}

	for _, rule := range ruleDefs {
		avps = append(avps, getRuleDefinitionAVP(rule))
	}

	return avps
}

func getInitialUsageMonitoringAVPs(monitors map[string]*protos.UsageMonitorCredit) []*diam.AVP {
	avps := make([]*diam.AVP, 0, len(monitors))
	for key, monitor := range monitors {
		avps = append(avps, getUsageMonitoringResponseAVP(key, monitor.ReturnBytes, monitor.MonitoringLevel))
	}
	return avps
}

func getUsageMonitorUpdates(account *subscriberAccount, monitorUpdates []*usageMonitorRequestAVP) ([]*diam.AVP, error) {
	avps := make([]*diam.AVP, 0, len(monitorUpdates))
	for _, update := range monitorUpdates {
		monitorCredit, ok := account.UsageMonitors[update.MonitoringKey]
		if !ok {
			return []*diam.AVP{}, fmt.Errorf("unknown monitoring key %s", update.MonitoringKey)
		}
		returnBytes := decrementUsageAndGetReturnBytes(monitorCredit, update.UsedServiceUnit.TotalOctets)
		avps = append(avps, getUsageMonitoringResponseAVP(monitorCredit.MonitoringKey, returnBytes, monitorCredit.MonitoringLevel))
	}
	return avps, nil
}

func getUsageMonitoringResponseAVP(monitoringKey string, returnBytes uint64, level protos.UsageMonitorCredit_MonitoringLevel) *diam.AVP {
	return diam.NewAVP(avp.UsageMonitoringInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MonitoringKey, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(monitoringKey)),
			diam.NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(returnBytes)),
				},
			}),
			diam.NewAVP(avp.UsageMonitoringLevel, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(level)),
		},
	})
}

func decrementUsageAndGetReturnBytes(monitorCredit *protos.UsageMonitorCredit, usedVolume uint64) uint64 {
	monitorCredit.Volume -= usedVolume
	if monitorCredit.Volume < 0 {
		monitorCredit.Volume = 0
	}
	if monitorCredit.ReturnBytes > monitorCredit.Volume {
		return monitorCredit.Volume
	}
	return monitorCredit.ReturnBytes
}

func getQosAVP(qos *lteprotos.FlowQos) *diam.AVP {
	qosAVPs := []*diam.AVP{}
	if qos.MaxReqBwUl != 0 {
		qosAVPs = append(
			qosAVPs,
			diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(qos.MaxReqBwUl)),
		)
	}
	if qos.MaxReqBwDl != 0 {
		qosAVPs = append(
			qosAVPs,
			diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(qos.MaxReqBwDl)),
		)
	}
	return diam.NewAVP(avp.QoSInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: qosAVPs,
	})
}

func getRuleDefinitionAVP(rule *protos.RuleDefinition) *diam.AVP {
	installAVPs := []*diam.AVP{
		diam.NewAVP(avp.ChargingRuleName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule.RuleName)),
		diam.NewAVP(avp.Precedence, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(rule.Precedence)),
	}
	if rule.RatingGroup != 0 {
		installAVPs = append(installAVPs, diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(rule.RatingGroup)))
	}
	if rule.MonitoringKey != "" {
		installAVPs = append(
			installAVPs,
			diam.NewAVP(avp.MonitoringKey, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(rule.MonitoringKey)),
		)
	}
	if rule.QosInformation != nil && (rule.QosInformation.MaxReqBwUl != 0 || rule.QosInformation.MaxReqBwDl != 0) {
		installAVPs = append(
			installAVPs,
			getQosAVP(rule.QosInformation),
		)
	}
	for _, flowDescription := range rule.FlowDescriptions {
		installAVPs = append(
			installAVPs,
			diam.NewAVP(avp.FlowDescription, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.IPFilterRule(flowDescription)),
		)
	}
	if rule.RedirectInformation != nil {
		RedirectInformation := []*diam.AVP{
			diam.NewAVP(avp.RedirectSupport, avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(rule.RedirectInformation.Support)),
			diam.NewAVP(avp.RedirectAddressType, avp.Mbit, 0, datatype.Enumerated(rule.RedirectInformation.AddressType)),
			diam.NewAVP(avp.RedirectServerAddress, avp.Mbit, 0, datatype.UTF8String(rule.RedirectInformation.ServerAddress)),
		}
		installAVPs = append(
			installAVPs,
			diam.NewAVP(avp.RedirectInformation, avp.Mbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: RedirectInformation,
			}),
		)
	}
	return diam.NewAVP(avp.ChargingRuleInstall, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ChargingRuleDefinition, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: installAVPs,
			}),
		},
	})
}

// sendAnswer sends a CCA to the connection given
func sendAnswer(
	ccr ccrMessage,
	conn diam.Conn,
	message *diam.Message,
	statusCode uint32,
	additionalAVPs ...*diam.AVP,
) {
	a := message.Answer(statusCode)
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, ccr.DestinationHost)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, ccr.DestinationRealm)
	a.NewAVP(avp.DestinationRealm, avp.Mbit, 0, ccr.OriginRealm)
	a.NewAVP(avp.DestinationHost, avp.Mbit, 0, ccr.OriginHost)
	a.NewAVP(avp.CCRequestType, avp.Mbit, 0, ccr.RequestType)
	a.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, ccr.RequestNumber)
	a.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID)
	for _, avp := range additionalAVPs {
		a.InsertAVP(avp)
	}
	// SessionID must be the first AVP
	a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID))

	_, err := a.WriteTo(conn)
	if err != nil {
		glog.V(2).Infof("Failed to write message to %s: %s\n%s\n",
			conn.RemoteAddr(), err, a)
		return
	}
	glog.V(2).Infof("Sent CCA to %s:\n", conn.RemoteAddr())
}
