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

// GyClient is a client to send Credit Control Request messages over diameter
// And receive Credit Control Answer messages in response
package gy

import (
	"net"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/gateway/service_registry"
)

const (
	RetryCount = 2
)

type CreditClient interface {
	SendCreditControlRequest(
		server *diameter.DiameterServerConfig,
		done chan interface{},
		request *CreditControlRequest,
	) error
	IgnoreAnswer(request *CreditControlRequest)
	EnableConnections() error
	DisableConnections(period time.Duration)
}

// ReAuthHandler defines a function that responds to a RAR message with an RAA
type ChargingReAuthHandler func(request *ChargingReAuthRequest) *ChargingReAuthAnswer

// GyClient holds the relevant state for sending and receiving diameter calls
// over Gy
type GyClient struct {
	diamClient   *diameter.Client
	serverCfg    *diameter.DiameterServerConfig
	globalConfig *GyGlobalConfig
}

type GyGlobalConfig struct {
	OCSOverwriteApn      string
	OCSServiceIdentifier string
	DisableGy            bool
	VirtualApnRules      []*credit_control.VirtualApnRule
}

var (
	serviceIdentifier int64 = -1
)

// NewGyClient constructs a new GyClient with the magma diameter settings
func NewConnectedGyClient(
	diamClient *diameter.Client,
	serverCfg *diameter.DiameterServerConfig,
	reAuthHandler ChargingReAuthHandler,
	cloudRegistry service_registry.GatewayRegistry,
	globalConfig *GyGlobalConfig,
) *GyClient {
	diamClient.RegisterAnswerHandlerForAppID(diam.CreditControl, diam.CHARGING_CONTROL_APP_ID, getCCAHandler())
	registerReAuthHandler(reAuthHandler, diamClient)
	if globalConfig != nil && len(globalConfig.OCSServiceIdentifier) > 0 {
		var err error
		serviceIdentifier, err = strconv.ParseInt(globalConfig.OCSServiceIdentifier, 10, 0)
		if err != nil {
			serviceIdentifier = -1
		}
	}
	if cloudRegistry != nil {
		diamClient.RegisterHandler(
			diam.AbortSession,
			diam.CHARGING_CONTROL_APP_ID,
			true,
			credit_control.NewASRHandler(diamClient, cloudRegistry))
	}
	return &GyClient{
		diamClient:   diamClient,
		serverCfg:    serverCfg,
		globalConfig: globalConfig,
	}
}

// NewGyClient constructs a new GyClient with the magma diameter settings
func NewGyClient(
	clientCfg *diameter.DiameterClientConfig,
	serverCfg *diameter.DiameterServerConfig,
	reAuthHandler ChargingReAuthHandler,
	cloudRegistry service_registry.GatewayRegistry,
	globalConfig *GyGlobalConfig,
) *GyClient {
	diamClient := diameter.NewClient(clientCfg)
	diamClient.BeginConnection(serverCfg)
	return NewConnectedGyClient(diamClient, serverCfg, reAuthHandler, cloudRegistry, globalConfig)
}

// SendCreditControlRequest sends a Credit Control Request to the
// given connection. Multiple requests can be sent in a row without waiting for
// the answer
// Example use:
// 	client := NewGyClient()
//  done :=  make(chan *CreditControlAnswer, 1)
// 	client.SendCreditControlRequest(server, done, requests)
// 	answer := <- done
// Input: DiameterServerConfig containing info about where to send messages
//				chan<- *CreditControlAnswer to send answers to
//			  CreditControlRequest containing the request to send
//
// Output: error if server connection failed
func (gyClient *GyClient) SendCreditControlRequest(
	server *diameter.DiameterServerConfig,
	done chan interface{},
	request *CreditControlRequest,
) error {
	additionalAVPs, err := getAdditionalAvps(request)
	if err != nil {
		return err
	}

	message, err := gyClient.createCreditControlMessage(
		server,
		request,
		additionalAVPs...)
	if err != nil {
		return err
	}

	glog.V(2).Infof("Sending Gy CCR message:\n%s\n", message)
	key := credit_control.GetRequestKey(credit_control.Gy, request.SessionID, request.RequestNumber)
	return gyClient.diamClient.SendRequest(server, done, message, key)
}

// GetAnswer returns a *CreditControlAnswer from the given interface channel
func GetAnswer(done <-chan interface{}) *CreditControlAnswer {
	answer := <-done
	return answer.(*CreditControlAnswer)
}

// IgnoreAnswers removes tracked requests in the request manager to ensure the
// request mapping does not leak. For example, if 10 requests are sent out, and
// 2 time out given the user's timeout duration, then those 2 requests should be
// ignored so that they don't leak
func (gyClient *GyClient) IgnoreAnswer(request *CreditControlRequest) {
	gyClient.diamClient.IgnoreAnswer(
		credit_control.GetRequestKey(credit_control.Gy, request.SessionID, request.RequestNumber),
	)
}

func (gyClient *GyClient) EnableConnections() error {
	if gyClient.globalConfig.DisableGy {
		return nil
	}
	gyClient.diamClient.EnableConnectionCreation()
	return gyClient.diamClient.BeginConnection(gyClient.serverCfg)
}

func (gyClient *GyClient) DisableConnections(period time.Duration) {
	if gyClient.globalConfig.DisableGy {
		return
	}
	gyClient.diamClient.DisableConnectionCreation(period)
}

// RegisterReAuthHandler adds a handler to the client for responding to RAR
// messages received from the OCS
func registerReAuthHandler(reAuthHandler ChargingReAuthHandler, diamClient *diameter.Client) {
	reqHandler := func(conn diam.Conn, message *diam.Message) {
		glog.V(2).Infof("Received Gy reauth message:\n%s\n", message)
		rar := &ChargingReAuthRequest{}
		if err := message.Unmarshal(rar); err != nil {
			glog.Errorf("Received unparseable RAR over Gy %s\n%s", message, err)
			return
		}
		go func() {
			raa := reAuthHandler(rar)
			raaMsg := createReAuthAnswerMessage(message, raa)
			raaMsg = diamClient.AddOriginAVPsToMessage(raaMsg)
			glog.V(2).Infof("Sending (responding) Gy reauth message:\n%s\n", raaMsg)
			_, err := raaMsg.WriteToWithRetry(conn, diamClient.Retries())
			if err != nil {
				glog.Errorf(
					"Gy RAA Write Failed for %s->%s, SessionID: %s - %v",
					conn.LocalAddr(), conn.RemoteAddr(), rar.SessionID, err)
				conn.Close() // close connection on error
			}
		}()
	}
	diamClient.RegisterRequestHandlerForAppID(diam.ReAuth, diam.CHARGING_CONTROL_APP_ID, reqHandler)
}

func createReAuthAnswerMessage(requestMsg *diam.Message, answer *ChargingReAuthAnswer) *diam.Message {
	ansMsg := requestMsg.Answer(answer.ResultCode)
	ansMsg.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(answer.SessionID)))
	return ansMsg
}

// getAdditionalAvps retrieves any extra AVPs based on the type of request.
// For update and terminate, it returns the used credit AVPs
func getAdditionalAvps(request *CreditControlRequest) ([]*diam.AVP, error) {
	avpList := make([]*diam.AVP, 0, len(request.Credits)+1)
	if len(request.TgppCtx.GetGyDestHost()) > 0 {
		avpList = append(avpList,
			diam.NewAVP(avp.DestinationHost, avp.Mbit, 0, datatype.DiameterIdentity(request.TgppCtx.GetGyDestHost())))
	}
	for _, credit := range request.Credits {
		avpList = append(avpList, getMSCCAVP(request.Type, credit))
	}
	return avpList, nil
}

// createCreditControlMessage creates a base message to be used for any Credit
// Control Request message. Init will just use this, and update and terminate
// pass in extra AVPs through additionalAVPs
// Input: context.Context which has information on where to send to,
//				CreditControlRequest with relevant request info
//			  ...*diam.AVP with any AVPs to add on
// Output: *diam.Message with all AVPs filled in, error if there was an issue
func (gyClient *GyClient) createCreditControlMessage(
	server *diameter.DiameterServerConfig,
	request *CreditControlRequest,
	additionalAVPs ...*diam.AVP,
) (*diam.Message, error) {
	m := diameter.NewProxiableRequest(diam.CreditControl, diam.CHARGING_CONTROL_APP_ID, nil)
	// m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("584144187966")) // TODO
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.CHARGING_CONTROL_APP_ID))
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(request.Type))
	m.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String(gyClient.diamClient.ServiceContextId()))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(request.RequestNumber))

	// Always add MSISDN if it's provided by AGW
	if len(request.Msisdn) > 0 {
		m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(credit_control.EndUserE164)),
				diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(request.Msisdn)),
			},
		})
	}

	m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(credit_control.EndUserIMSI)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(request.IMSI)),
		},
	})

	if len(request.Imei) > 0 {
		m.NewAVP(avp.UserEquipmentInfo, 0, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.UserEquipmentInfoType, 0, 0, datatype.Enumerated(0)), // imeisv
				diam.NewAVP(avp.UserEquipmentInfoValue, 0, 0, datatype.OctetString(request.Imei)),
			},
		})
	}
	m.InsertAVP(getServiceInfoAvp(server, request, gyClient.globalConfig))

	m.NewAVP(avp.MultipleServicesIndicator, avp.Mbit, 0, datatype.Enumerated(0x01))

	for _, additionalAVP := range additionalAVPs {
		m.InsertAVP(additionalAVP)
	}

	// SessionID must be the first AVP
	m.InsertAVP(diam.NewAVP(
		avp.SessionID,
		avp.Mbit,
		0,
		datatype.UTF8String(diameter.EncodeSessionID(gyClient.diamClient.OriginHost(), request.SessionID))))

	return m, nil
}

// getServiceInfoAvp() Fills the Service-Information AVP
func getServiceInfoAvp(server *diameter.DiameterServerConfig, request *CreditControlRequest, gyGlobalConfig *GyGlobalConfig) *diam.AVP {

	svcInfoGrp := []*diam.AVP{}
	csAddr, _, _ := net.SplitHostPort(server.Addr)

	ratType := request.RatType
	if len(ratType) == 0 {
		ratType = RAT_TYPE_EUTRAN
	}
	psInfoAvps := []*diam.AVP{
		// Set PDP Type as IPV4(0)
		diam.NewAVP(avp.TGPPPDPType, avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(0)),
		// Set RAT Type as EUTRAN(6). See 3GPP TS 29.274, 8.17 "Table 8.17-1: RAT Type values"
		diam.NewAVP(avp.TGPPRATType, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(ratType)),
		// Set it to 0
		diam.NewAVP(avp.TGPPSelectionMode, avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("0")),
		diam.NewAVP(avp.TGPPNSAPI, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString("5")),
		diam.NewAVP(avp.CGAddress, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Address(net.ParseIP(csAddr))),
	}
	if pdpAddr := net.ParseIP(request.UeIPV4); pdpAddr != nil {
		psInfoAvps = append(
			psInfoAvps,
			diam.NewAVP(avp.PDPAddress, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Address(pdpAddr)))
	}
	psInfoGrp := &diam.GroupedAVP{AVP: psInfoAvps}

	if len(request.SpgwIPV4) > 0 {
		psInfoGrp.AddAVP(diam.NewAVP(avp.SGSNAddress, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Address(net.ParseIP(request.SpgwIPV4))))
		psInfoGrp.AddAVP(diam.NewAVP(avp.GGSNAddress, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Address(net.ParseIP(request.SpgwIPV4))))
	}
	if len(request.PlmnID) > 0 {
		psInfoGrp.AddAVP(diam.NewAVP(avp.TGPPSGSNMCCMNC, avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(request.PlmnID)))
		psInfoGrp.AddAVP(diam.NewAVP(avp.TGPPGGSNMCCMNC, avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(request.PlmnID)))
	}

	apn := getAPNFromConfig(gyGlobalConfig, request.Apn, request.ChargingCharacteristics)
	if len(apn) > 0 {
		psInfoGrp.AddAVP(diam.NewAVP(avp.CalledStationID, avp.Mbit, 0, apn))
	}

	if len(request.UserLocation) > 0 {
		psInfoGrp.AddAVP(diam.NewAVP(avp.TGPPUserLocationInfo, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(string(request.UserLocation))))
	}
	if len(request.GcID) > 0 {
		psInfoGrp.AddAVP(diam.NewAVP(avp.TGPPChargingID, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(request.GcID)))
	}
	/********************** TBD - doesn't work with some OCSes *****************
	if request.Qos != nil {
		qosGrp := &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.APNAggregateMaxBitrateDL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(request.Qos.ApnAggMaxBitRateDL)),
				diam.NewAVP(avp.APNAggregateMaxBitrateUL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(request.Qos.ApnAggMaxBitRateUL)),
			},
		}
		psInfoGrp.AddAVP(diam.NewAVP(avp.QoSInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, qosGrp))
	}
	***************************************************************************/

	svcInfoGrp = append(
		svcInfoGrp,
		diam.NewAVP(avp.PSInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, psInfoGrp),
	)
	return diam.NewAVP(avp.ServiceInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{AVP: svcInfoGrp})
}

// getAPNFromConfig returns a new apn value to overwrite the one in the request based on list of regex definied in Gy config.
// If Virtual APN config is not defined, the function returnis OCSOverwriteApn instead.
// Input: GyGlobalConfig and the APN received from the request
// Output: Overwritten apn value
func getAPNFromConfig(gyGlobalConfig *GyGlobalConfig, requestAPN, chargingCharacteristics string) datatype.UTF8String {
	apn := datatype.UTF8String(requestAPN)
	if gyGlobalConfig != nil {
		if len(gyGlobalConfig.VirtualApnRules) > 0 {
			apn = datatype.UTF8String(credit_control.MatchAndGetOverwriteApn(requestAPN, chargingCharacteristics, gyGlobalConfig.VirtualApnRules))
		} else if len(gyGlobalConfig.OCSOverwriteApn) > 0 {
			// OverwriteApn is deprecated transition to VirtualApnRules
			apn = datatype.UTF8String(gyGlobalConfig.OCSOverwriteApn)
		}
	}
	return apn
}

// getMSCCAVP retrieves the MultipleServicesCreditControl AVP for the
// given used Credits. This is used for terminate and update CCRs credit updates
// and terminations
// Input: UsedCredits with input/output/total octets used
// Output: *diam.Message with all AVPs filled in, error if there was an issue
func getMSCCAVP(requestType credit_control.CreditRequestType, credits *UsedCredits) *diam.AVP {
	avpGroup := []*diam.AVP{
		diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(credits.RatingGroup)),
	}

	// Requested-Service-Unit can only be send in CCR-I and CCR-U
	if requestType != credit_control.CRTTerminate {
		var usuGrp []*diam.AVP
		if credits.RequestedUnits == nil {
			glog.Errorf("Not adding AVP Requested-Service-Unit. Not found on credit request for session %+v", credits)
			usuGrp = []*diam.AVP{}
		} else {
			usuGrp = []*diam.AVP{
				diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(credits.RequestedUnits.Rx)),
				diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(credits.RequestedUnits.Tx)),
				diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(credits.RequestedUnits.Total)),
			}
		}
		avpGroup = append(
			avpGroup, diam.NewAVP(avp.RequestedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{AVP: usuGrp}))

	}

	if serviceIdentifier >= 0 {
		avpGroup = append(
			avpGroup,
			diam.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(serviceIdentifier)))
	} else if credits.ServiceIdentifier != nil {
		avpGroup = append(
			avpGroup,
			diam.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(*credits.ServiceIdentifier)))
	}

	// Used credits can only be sent on updates and terminates
	if requestType != credit_control.CRTInit {
		usuGrp := []*diam.AVP{
			diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(credits.InputOctets)),
			diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(credits.OutputOctets)),
			diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(credits.TotalOctets)),
		}

		// For documentation on where the reporting reason AVP is placed, see section 7.2.175 on
		// https://www.etsi.org/deliver/etsi_ts/132200_132299/132299/12.06.00_60/ts_132299v120600p.pdf
		switch credits.Type {
		case FINAL, VALIDITY_TIMER_EXPIRED, FORCED_REAUTHORISATION, QHT, RATING_CONDITION_CHANGE:
			avpGroup = append(
				avpGroup,
				diam.NewAVP(
					avp.ReportingReason, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(credits.Type)))
		case QUOTA_EXHAUSTED, THRESHOLD, OTHER_QUOTA_TYPE, POOL_EXHAUSTED:
			usuGrp = append(
				usuGrp,
				diam.NewAVP(
					avp.ReportingReason, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(credits.Type)))

		}
		avpGroup = append(
			avpGroup, diam.NewAVP(avp.UsedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{AVP: usuGrp}))
	}

	return diam.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: avpGroup,
	})
}

// getReceivedCredits gets the received octets if applicable from the unmarshalled
// diameter message,
func getReceivedCredits(cca *CCADiameterMessage) []*ReceivedCredits {
	creditList := make([]*ReceivedCredits, 0, len(cca.CreditControl))
	for _, mscc := range cca.CreditControl {
		receivedCredits := &ReceivedCredits{
			ResultCode:          mscc.ResultCode,
			GrantedUnits:        &mscc.GrantedServiceUnit,
			ValidityTime:        mscc.ValidityTime,
			RatingGroup:         mscc.RatingGroup,
			ServiceIdentifier:   mscc.ServiceIdentifier,
			FinalUnitIndication: mscc.FinalUnitIndication,
		}
		creditList = append(creditList, receivedCredits)
	}
	return creditList
}

// getCCAHandler returns a callback function to use when an answer is received
// over Gy. Parses the message for relevant information.
// Input: requestTracker *diameter.RequestTracker to get channel to send answers out
//				when one is received
// Output: diam.HandlerFunc
func getCCAHandler() diameter.AnswerHandler {
	return func(message *diam.Message) diameter.KeyAndAnswer {
		glog.V(2).Infof("Received Gy CCA message:\n%s\n", message)
		var cca CCADiameterMessage
		if err := message.Unmarshal(&cca); err != nil {
			metrics.GyUnparseableMsg.Inc()
			glog.Errorf("Received unparseable CCA over Gy")
			return diameter.KeyAndAnswer{}
		}
		sid := diameter.DecodeSessionID(cca.SessionID)
		return diameter.KeyAndAnswer{
			Key: credit_control.GetRequestKey(credit_control.Gy, sid, cca.RequestNumber),
			Answer: &CreditControlAnswer{
				ResultCode:    cca.ResultCode,
				SessionID:     sid,
				OriginHost:    cca.OriginHost,
				RequestNumber: cca.RequestNumber,
				Credits:       getReceivedCredits(&cca),
			},
		}
	}
}
