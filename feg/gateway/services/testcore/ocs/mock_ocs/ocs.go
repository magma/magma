/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package mock_ocs

import (
	"errors"
	"fmt"
	"net"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/testcore/mock_driver"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const DiameterCreditLimitReached = 4012

type CreditBucket struct {
	Unit   protos.CreditInfo_UnitType
	Volume *protos.Octets
}

type SubscriberSessionState struct {
	SessionID  string
	Connection diam.Conn
}

type SubscriberAccount struct {
	ChargingCredit map[uint32]*CreditBucket // map of charging key to credit bucket
	CurrentState   *SubscriberSessionState
}

type OCSConfig struct {
	MaxUsageOctets  *protos.Octets
	MaxUsageTime    uint32
	ValidityTime    uint32
	ServerConfig    *diameter.DiameterServerConfig
	GyInitMethod    gy.InitMethod
	UseMockDriver   bool
	RedirectAddress string
	FinalUnitAction protos.FinalUnitAction
}

// OCSDiamServer wraps an OCS storing subscriber accounts and their credit
type OCSDiamServer struct {
	diameterSettings    *diameter.DiameterClientConfig
	ocsConfig           *OCSConfig
	accounts            map[string]*SubscriberAccount // map of IMSI to subscriber account
	mux                 *sm.StateMachine
	LastMessageReceived *ccrMessage
	mockDriver          *mock_driver.MockDriver
}

// NewOCSDiamServer initializes an OCS with an empty account map
// Input: *sm.Settings containing the diameter related parameters
//				*TestOCSConfig containing the server address, and standard OCS settings
//					like how many bytes to allocate to users
//
// Output: a new OCSDiamServer
func NewOCSDiamServer(
	diameterSettings *diameter.DiameterClientConfig,
	ocsConfig *OCSConfig,
) *OCSDiamServer {
	return &OCSDiamServer{
		diameterSettings: diameterSettings,
		ocsConfig:        ocsConfig,
		accounts:         make(map[string]*SubscriberAccount),
	}
}

// Reset is an GRPC procedure which resets the server to its default state.
// It will be called from the gateway.
func (srv *OCSDiamServer) Reset(
	_ context.Context,
	req *orcprotos.Void,
) (*orcprotos.Void, error) {
	return nil, errors.New("Not implemented")
}

// ConfigServer is an GRPC procedure which configure the server to respond
// to requests. It will be called from the gateway
func (srv *OCSDiamServer) ConfigServer(
	_ context.Context,
	config *protos.ServerConfiguration,
) (*orcprotos.Void, error) {
	return nil, errors.New("Not implemented")
}

// Start begins the server and blocks, listening to the network
// Output: error if the server could not be started
func (srv *OCSDiamServer) Start(lis net.Listener) error {
	srv.mux = sm.New(&sm.Settings{
		OriginHost:       datatype.DiameterIdentity(srv.diameterSettings.Host),
		OriginRealm:      datatype.DiameterIdentity(srv.diameterSettings.Realm),
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      datatype.UTF8String(srv.diameterSettings.ProductName),
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	})
	srv.mux.Handle(diam.CCR, getCCRHandler(srv))
	serverConfig := srv.ocsConfig.ServerConfig
	server := &diam.Server{
		Network: serverConfig.Protocol,
		Addr:    serverConfig.Addr,
		Handler: srv.mux,
		Dict:    nil,
	}
	return server.Serve(lis)
}

func (srv *OCSDiamServer) StartListener() (net.Listener, error) {
	serverConfig := srv.ocsConfig.ServerConfig

	network := serverConfig.Protocol
	if len(network) == 0 {
		network = "tcp"
	}
	addr := serverConfig.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	l, e := diam.Listen(network, addr)
	if e != nil {
		return nil, e
	}
	return l, nil
}

// NewAccount adds a subscriber to the OCS to be tracked
// Input: string containing the subscriber IMSI (can be in any form)
func (srv *OCSDiamServer) CreateAccount(
	_ context.Context,
	subscriberID *lteprotos.SubscriberID,
) (*orcprotos.Void, error) {
	srv.accounts[subscriberID.Id] = &SubscriberAccount{
		ChargingCredit: make(map[uint32]*CreditBucket),
	}
	glog.V(2).Infof("New account %s added", subscriberID.Id)
	return &orcprotos.Void{}, nil
}

// SetOCSSettings changes the standard OCS return values. All parameters are
// optional, and this only sets the non-nil ones.
// Input: *uint32 optional maximum bytes to return in a CCA
//			  *uint32 optional maximum time to return in a CCA
//			  *uint32 optional credit validity time to return in a CCA
func (srv *OCSDiamServer) SetOCSSettings(
	_ context.Context,
	ocsConfig *protos.OCSConfig,
) (*orcprotos.Void, error) {
	config := srv.ocsConfig
	config.MaxUsageOctets = ocsConfig.MaxUsageOctets
	config.MaxUsageTime = ocsConfig.MaxUsageTime
	config.ValidityTime = ocsConfig.ValidityTime
	config.UseMockDriver = ocsConfig.UseMockDriver
	config.RedirectAddress = ocsConfig.RedirectAddress
	config.FinalUnitAction = ocsConfig.FinalUnitAction
	return &orcprotos.Void{}, nil
}

// SetCredit sets or overrides the prepaid credit allocated for an account
// Input: string IMSI for the account
//			  uint32 charging key to add credit to
//			  uint64 volume (in any units) to set this bucket to
//		    UnitType dictating which unit the volume represents
// Output: error if account could not be found
func (srv *OCSDiamServer) SetCredit(
	_ context.Context,
	creditInfo *protos.CreditInfo,
) (*orcprotos.Void, error) {
	account, ok := srv.accounts[creditInfo.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", creditInfo.Imsi)
	}
	account.ChargingCredit[creditInfo.ChargingKey] = &CreditBucket{
		Unit:   creditInfo.UnitType,
		Volume: creditInfo.Volume,
	}
	return &orcprotos.Void{}, nil
}

// GetCredits returns all the credits allocated for an account
// Input: string IMSI for the account
// Output: map[uint32]*CreditBucket a map of charging key to credit bucket
//			   error if account could not be found
func (srv *OCSDiamServer) GetCredits(
	_ context.Context,
	subscriberID *lteprotos.SubscriberID,
) (*protos.CreditInfos, error) {
	account, ok := srv.accounts[subscriberID.Id]
	if !ok {
		return &protos.CreditInfos{}, fmt.Errorf("Could not find imsi %s", subscriberID.Id)
	}
	infos := make(map[uint32]*protos.CreditInfo)
	for id, chargingCredit := range account.ChargingCredit {
		infos[id] =
			&protos.CreditInfo{
				UnitType: chargingCredit.Unit,
				Volume:   chargingCredit.Volume,
			}
	}
	return &protos.CreditInfos{CreditInformation: infos}, nil
}

// Reset eliminates all the accounts allocated for the system.
func (srv *OCSDiamServer) ClearSubscribers(_ context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	glog.V(2).Infof("Accounts (%d) will be deleted from OCS:", len(srv.accounts))
	for imsi, subs := range srv.accounts {
		glog.V(2).Infof("\tRemaing credit for IMSI: %s", imsi)
		for key, credits := range subs.ChargingCredit {
			glog.V(2).Infof("\t - key %d, Total:%d Tx:%d Rx:%d",
				key,
				credits.Volume.TotalOctets,
				credits.Volume.OutputOctets,
				credits.Volume.InputOctets,
			)
		}
	}
	srv.accounts = make(map[string]*SubscriberAccount)
	glog.V(2).Info("All accounts deleted.")
	return &orcprotos.Void{}, nil
}

func (srv *OCSDiamServer) SetExpectations(_ context.Context, req *protos.GyCreditControlExpectations) (*orcprotos.Void, error) {
	es := []mock_driver.Expectation{}
	for _, e := range req.Expectations {
		es = append(es, mock_driver.Expectation(GyExpectation{e}))
	}
	srv.mockDriver = mock_driver.NewMockDriver(es, req.UnexpectedRequestBehavior, GyAnswer{req.GyDefaultCca})
	return &orcprotos.Void{}, nil
}

func (srv *OCSDiamServer) AssertExpectations(_ context.Context, void *orcprotos.Void) (*protos.GyCreditControlResult, error) {
	srv.mockDriver.Lock()
	defer srv.mockDriver.Unlock()

	results, errs := srv.mockDriver.AggregateResults()
	return &protos.GyCreditControlResult{Results: results, Errors: errs}, nil
}

// ReAuth initiates a reauth call for a subscriber and optional rating group.
// It waits for any answer from the OCS
func (srv *OCSDiamServer) ReAuth(
	ctx context.Context,
	target *protos.ChargingReAuthTarget,
) (*protos.ChargingReAuthAnswer, error) {
	account, ok := srv.accounts[target.Imsi]
	if !ok {
		return nil, fmt.Errorf("Could not find imsi %s", target.Imsi)
	}
	if account.CurrentState == nil {
		return nil, fmt.Errorf("Credit client location unknown for imsi %s", target.Imsi)
	}
	done := make(chan *gy.ChargingReAuthAnswer)
	srv.mux.Handle(diam.RAA, handleRAA(done))
	err := sendRAR(account.CurrentState, &target.RatingGroup, srv.mux.Settings())
	if err != nil {
		glog.Errorf("Error sending RaR for target IMSI=%v, RG=%v: %v", target.GetImsi(), target.GetRatingGroup(), err)
		return nil, err
	}
	select {
	case raa := <-done:
		return &protos.ChargingReAuthAnswer{SessionId: diameter.DecodeSessionID(raa.SessionID), ResultCode: raa.ResultCode}, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("No RAA received")
	}
}

func sendRAR(state *SubscriberSessionState, ratingGroup *uint32, cfg *sm.Settings) error {
	meta, ok := smpeer.FromContext(state.Connection.Context())
	if !ok {
		return fmt.Errorf("peer metadata unavailable")
	}
	m := diameter.NewProxiableRequest(diam.ReAuth, diam.CHARGING_CONTROL_APP_ID, nil)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(state.SessionID))
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, cfg.OriginHost)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, cfg.OriginRealm)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, meta.OriginRealm)
	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, meta.OriginHost)
	if ratingGroup != nil {
		m.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(*ratingGroup))
	}
	glog.V(2).Infof("Sending RAR to %s\n%s", state.Connection.RemoteAddr(), m)
	_, err := m.WriteTo(state.Connection)
	return err
}

func handleRAA(done chan *gy.ChargingReAuthAnswer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var raa gy.ChargingReAuthAnswer
		if err := m.Unmarshal(&raa); err != nil {
			glog.Errorf("Received unparseable RAA over Gy %s", m)
			return
		}
		done <- &raa
	}
}
