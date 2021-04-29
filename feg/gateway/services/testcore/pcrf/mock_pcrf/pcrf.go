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

package mock_pcrf

import (
	"context"
	"fmt"
	"net"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/testcore/mock_driver"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// This is a temporary struct to handle unsupported AVP Charging-Rule-Report in go-diameter
type ReAuthAnswer struct {
	SessionID  string `avp:"Session-Id"`
	ResultCode uint32 `avp:"Result-Code"`
}

type creditByMkey map[string]*protos.UsageMonitor

type SubscriberSessionState struct {
	SessionID  string
	Connection diam.Conn
}

type subscriberAccount struct {
	RuleNames       []string
	RuleBaseNames   []string
	RuleDefinitions []*protos.RuleDefinition
	UsageMonitors   creditByMkey
	CurrentState    *SubscriberSessionState
}

// PCRFServer wraps an PCRF storing subscribers and their rules
type PCRFServer struct {
	diameterClientConfig    *diameter.DiameterClientConfig // Gx Client Config
	diameterServerConfig    *diameter.DiameterServerConfig
	serviceConfig           *protos.PCRFConfigs
	subscribers             map[string]*subscriberAccount // map of imsi to to rules
	mux                     *sm.StateMachine
	lastDiamMessageReceived *diam.Message
	mockDriver              *mock_driver.MockDriver
}

// NewPCRFServer initializes an PCRF with an empty rule map
// Input: *sm.Settings containing the diameter related parameters
//				*TestPCRFConfig containing the server address
//
// Output: a new PCRFServer
func NewPCRFServer(clientConfig *diameter.DiameterClientConfig, serverConfig *diameter.DiameterServerConfig) *PCRFServer {
	return &PCRFServer{
		diameterClientConfig: clientConfig,
		diameterServerConfig: serverConfig,
		subscribers:          map[string]*subscriberAccount{},
		serviceConfig:        &protos.PCRFConfigs{},
	}
}

// Reset is an GRPC procedure which configure the server to the default status.
// It will be called from the gateway.
func (srv *PCRFServer) Reset(
	_ context.Context,
	_ *orcprotos.Void,
) (*orcprotos.Void, error) {
	return nil, errors.New("Not implemented")
}

// ConfigServer is an GRPC procedure which configure the server to respond
// to requests. It will be called from the gateway
func (srv *PCRFServer) ConfigServer(
	_ context.Context,
	_ *protos.ServerConfiguration,
) (*orcprotos.Void, error) {
	return nil, errors.New("Not implemented")
}

// Start begins the server and blocks, listening to the network
// Output: error if the server could not be started
func (srv *PCRFServer) Start(lis net.Listener) error {
	srv.mux = sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(srv.diameterClientConfig.Host),
		OriginRealm:      datatype.DiameterIdentity(srv.diameterClientConfig.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(srv.diameterClientConfig.ProductName),
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	})
	srv.mux.HandleIdx(
		diam.CommandIndex{AppID: diam.GX_CHARGING_CONTROL_APP_ID, Code: diam.CreditControl, Request: true},
		getCCRHandler(srv))
	go logErrors(srv.mux.ErrorReports())
	serverConfig := srv.diameterServerConfig
	server := &diam.Server{
		Network: serverConfig.Protocol,
		Addr:    serverConfig.Addr,
		Handler: srv.mux,
		Dict:    nil,
	}
	return server.Serve(lis)
}

// StartListener starts a listener based on ServerConfig
// If ServerConfig did not have valid values, default values would be used
func (srv *PCRFServer) StartListener() (net.Listener, error) {
	serverConfig := srv.diameterServerConfig
	network := serverConfig.Protocol
	addr := serverConfig.Addr
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

func (srv *PCRFServer) SetPCRFConfigs(
	_ context.Context,
	configs *protos.PCRFConfigs,
) (*orcprotos.Void, error) {
	srv.serviceConfig = configs
	return &orcprotos.Void{}, nil
}

// NewSubscriber adds a subscriber to the PCRF to be tracked
// Input: string containing the subscriber IMSI (can be in any form)
func (srv *PCRFServer) CreateAccount(
	_ context.Context,
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
func (srv *PCRFServer) SetRules(
	_ context.Context,
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

func (srv *PCRFServer) SetUsageMonitors(
	_ context.Context,
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
func (srv *PCRFServer) GetRuleNames(imsi string) ([]string, error) {
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
func (srv *PCRFServer) GetRuleBaseNames(imsi string) ([]string, error) {
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
func (srv *PCRFServer) GetRuleDefinitions(
	imsi string,
) ([]*protos.RuleDefinition, error) {
	account, ok := srv.subscribers[imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", imsi)
	}
	return account.RuleDefinitions, nil
}

// Reset eliminates all the subscribers allocated for the system.
func (srv *PCRFServer) ClearSubscribers(_ context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	srv.subscribers = map[string]*subscriberAccount{}
	glog.V(2).Info("All accounts deleted.")
	return &orcprotos.Void{}, nil
}

func (srv *PCRFServer) SetExpectations(_ context.Context, req *protos.GxCreditControlExpectations) (*orcprotos.Void, error) {
	es := []mock_driver.Expectation{}
	for _, e := range req.Expectations {
		es = append(es, mock_driver.Expectation(GxExpectation{e}))
	}
	srv.mockDriver = mock_driver.NewMockDriver(es, req.UnexpectedRequestBehavior, GxAnswer{req.GxDefaultCca})
	return &orcprotos.Void{}, nil
}

func (srv *PCRFServer) AssertExpectations(_ context.Context, void *orcprotos.Void) (*protos.GxCreditControlResult, error) {
	srv.mockDriver.Lock()
	defer srv.mockDriver.Unlock()

	results, errs := srv.mockDriver.AggregateResults()
	return &protos.GxCreditControlResult{Results: results, Errors: errs}, nil
}

// AbortSession call for a subscriber
// Initiate a Abort session request and provide the response
func (srv *PCRFServer) AbortSession(
	_ context.Context,
	req *protos.AbortSessionRequest,
) (*protos.AbortSessionAnswer, error) {
	glog.V(1).Infof("AbortSession: imsi %s", req.GetImsi())
	account, ok := srv.subscribers[req.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", req.Imsi)
	}
	if account.CurrentState == nil {
		return nil, fmt.Errorf("Credit client State unknown for imsi %s", req.Imsi)
	}

	var asaHandler diam.HandlerFunc
	resp := make(chan *diameter.ASA)
	asaHandler = func(conn diam.Conn, msg *diam.Message) {
		var asa diameter.ASA
		if err := msg.Unmarshal(&asa); err != nil {
			glog.Errorf("Received unparseable ASA over Gx, %s\n%s", err, msg)
			return
		}
		glog.V(2).Infof("Received ASA \n%s", msg)
		resp <- &diameter.ASA{SessionID: asa.SessionID, ResultCode: asa.ResultCode}
	}
	srv.mux.Handle(diam.ASA, asaHandler)
	err := sendASR(account.CurrentState, srv.mux.Settings())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send Gx ASR")
	}
	select {
	case asa := <-resp:
		return &protos.AbortSessionAnswer{SessionId: diameter.DecodeSessionID(asa.SessionID), ResultCode: asa.ResultCode}, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("No ASA received")
	}
}

// GetLastAVPreceived gets the last message in diam format received
// Message gets overwriten every time a new CCR is sent
func (srv *PCRFServer) GetLastAVPreceived() (*diam.Message, error) {
	if srv.lastDiamMessageReceived == nil {
		return nil, fmt.Errorf("No AVP message received")
	}
	return srv.lastDiamMessageReceived, nil
}

func sendASR(state *SubscriberSessionState, cfg *sm.Settings) error {
	meta, ok := smpeer.FromContext(state.Connection.Context())
	if !ok {
		return fmt.Errorf("peer metadata unavailable")
	}
	m := diameter.NewProxiableRequest(diam.AbortSession, diam.GX_CHARGING_CONTROL_APP_ID, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(state.SessionID))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	fmt.Printf("Sending Abort Session to %s\n%s", state.Connection.RemoteAddr(), m)
	glog.V(2).Infof("Sending Abort Session to %s\n%s", state.Connection.RemoteAddr(), m)
	_, err := m.WriteTo(state.Connection)
	return err
}

// ReAuth call for a subscriber
// Initiate a RAR requenst and handle a response
func (srv *PCRFServer) ReAuth(
	ctx context.Context,
	target *protos.PolicyReAuthTarget,
) (*protos.PolicyReAuthAnswer, error) {
	account, ok := srv.subscribers[target.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", target.Imsi)
	}
	if account.CurrentState == nil {
		return nil, fmt.Errorf("Credit client State unknown for imsi %s", target.Imsi)
	}

	var raaHandler diam.HandlerFunc
	done := make(chan *gx.PolicyReAuthAnswer)
	raaHandler = func(conn diam.Conn, msg *diam.Message) {
		// TODO Remove ReAuthAnswer and use PolicyReAuthAnswer once go-diameter supports Charging-Rule-Report
		var raa ReAuthAnswer
		if err := msg.Unmarshal(&raa); err != nil {
			glog.Errorf("Received unparseable RAA over Gx,  %s\n%s", err, msg)
			return
		}
		glog.V(2).Infof("Received RAA \n%s", msg)
		done <- &gx.PolicyReAuthAnswer{SessionID: raa.SessionID, ResultCode: raa.ResultCode}
	}
	srv.mux.Handle(diam.RAA, raaHandler)
	err := sendRAR(account.CurrentState, target, srv.mux.Settings())
	if err != nil {
		glog.Errorf("Error sending RaR for target %v: %v", target.GetImsi(), err)
		return nil, err
	}
	select {
	case raa := <-done:
		return &protos.PolicyReAuthAnswer{SessionId: diameter.DecodeSessionID(raa.SessionID), ResultCode: raa.ResultCode}, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("No RAA received")
	}
}

func sendRAR(state *SubscriberSessionState, target *protos.PolicyReAuthTarget, cfg *sm.Settings) error {
	meta, ok := smpeer.FromContext(state.Connection.Context())
	if !ok {
		return fmt.Errorf("peer metadata unavailable")
	}
	m := diameter.NewProxiableRequest(diam.ReAuth, diam.GX_CHARGING_CONTROL_APP_ID, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(state.SessionID))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)

	additionalAVPs := []*diam.AVP{}
	// Construct AVPs for Rules to Install
	ruleInstalls := target.GetRulesToInstall()
	if ruleInstalls != nil {
		additionalAVPs = append(additionalAVPs,
			toRuleInstallAVPs(
				ruleInstalls.GetRuleNames(),
				ruleInstalls.GetRuleBaseNames(),
				ruleInstalls.GetRuleDefinitions(),
				nil,
				nil)...,
		)
	}

	// Construct AVPs for Rules to Remove
	ruleRemovals := target.GetRulesToRemove()
	if ruleRemovals != nil {
		additionalAVPs = append(additionalAVPs,
			toRuleRemovalAVPs(
				ruleRemovals.GetRuleNames(),
				ruleRemovals.GetRuleBaseNames(),
			)...,
		)
	}

	// Construct AVPs for UsageMonitoring
	monitorInstalls := target.GetUsageMonitoringInfos()
	for _, monitor := range monitorInstalls {
		octets := monitor.GetOctets()
		if octets == nil {
			glog.Errorf("Monitor Octets is nil, skipping.")
			continue
		}
		additionalAVPs = append(additionalAVPs, toUsageMonitoringInfoAVP(string(monitor.MonitoringKey), octets, monitor.MonitoringLevel))
	}

	for _, avp := range additionalAVPs {
		m.InsertAVP(avp)
	}
	glog.V(2).Infof("Sending RAR to %s\n%s", state.Connection.RemoteAddr(), m)
	_, err := m.WriteTo(state.Connection)
	return err
}
