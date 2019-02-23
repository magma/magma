/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gx

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/golang/glog"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
)

// PolicyClient is an interface to define something that sends requests over Gx.
// This can be used to stub out requests
type PolicyClient interface {
	SendCreditControlRequest(
		server *diameter.DiameterServerConfig,
		done chan interface{},
		request *CreditControlRequest,
	) error
	IgnoreAnswer(request *CreditControlRequest)
	EnableConnections()
	DisableConnections(period time.Duration)
}

// GxClient is a client to send Gx Credit Control Request messages over diameter
// And receive Gx Credit Control Answer messages in response
// Although Gy and Gx both send Credit Control Requests, their Application IDs,
// allowed AVPs, and purposes are different
type GxClient struct {
	diamClient      *diameter.Client
	pcrf91Compliant bool // to support PCRF which is 29.212 release 9.1 compliant
}

// NewConnectedGxClient contructs a new GxClient with the magma diameter settings
func NewConnectedGxClient(
	diamClient *diameter.Client,
	reAuthHandler ReAuthHandler,
) *GxClient {
	diamClient.RegisterAnswerHandlerForAppID(diam.CreditControl, diam.GX_CHARGING_CONTROL_APP_ID, ccaHandler)
	registerReAuthHandler(reAuthHandler, diamClient)

	return &GxClient{
		diamClient:      diamClient,
		pcrf91Compliant: *pcrf91Compliant || isThruthy(os.Getenv(PCRF91CompliantEnv))}
}

func isThruthy(value string) bool {
	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return false
	}
	value = strings.ToLower(value)
	if value == "0" || strings.HasPrefix(value, "false") || strings.HasPrefix(value, "no") {
		return false
	}
	return true
}

// NewGxClient contructs a new GxClient with the magma diameter settings
func NewGxClient(
	clientCfg *diameter.DiameterClientConfig,
	servers []*diameter.DiameterServerConfig,
	reAuthHandler ReAuthHandler,
) *GxClient {
	diamClient := diameter.NewClient(clientCfg)
	for _, server := range servers {
		diamClient.BeginConnection(server)
	}
	return NewConnectedGxClient(diamClient, reAuthHandler)
}

// SendCreditControlRequest sends a Gx Credit Control Requests to the
// given connection
// Input: DiameterServerConfig containing info about where to send messages
//				chan<- *CreditControlAnswer to send answers to
//			  CreditControlRequest with the request to send
//
// Output: error if server connection failed
func (gxClient *GxClient) SendCreditControlRequest(
	server *diameter.DiameterServerConfig,
	done chan interface{},
	request *CreditControlRequest,
) error {
	additionalAVPs, err := gxClient.getAdditionalAvps(request)
	if err != nil {
		return err
	}

	message, err := gxClient.createCreditControlMessage(request, additionalAVPs...)
	if err != nil {
		return err
	}

	glog.V(2).Infof("Sending Gx CCR message\n%s\n", message)
	key := credit_control.GetRequestKey(credit_control.Gx, request.SessionID, request.RequestNumber)
	return gxClient.diamClient.SendRequest(server, done, message, key)
}

// GetAnswer returns a *CreditControlAnswer from the given interface channel
func GetAnswer(done <-chan interface{}) *CreditControlAnswer {
	answer := <-done
	return answer.(*CreditControlAnswer)
}

// IgnoreAnswer removes tracked requests in the request manager to ensure the
// request mapping does not leak. For example, if 10 requests are sent out, and
// 2 time out given the user's timeout duration, then those 2 requests should be
// ignored so that they don't leak
func (gxClient *GxClient) IgnoreAnswer(request *CreditControlRequest) {
	gxClient.diamClient.IgnoreAnswer(
		credit_control.GetRequestKey(credit_control.Gx, request.SessionID, request.RequestNumber),
	)
}

func (gxClient *GxClient) EnableConnections() {
	gxClient.diamClient.EnableConnectionCreation()
}

func (gxClient *GxClient) DisableConnections(period time.Duration) {
	gxClient.diamClient.DisableConnectionCreation(period)
}

// Register reauth request handler
func registerReAuthHandler(reAuthHandler ReAuthHandler, diamClient *diameter.Client) {
	handler := func(conn diam.Conn, message *diam.Message) {
		rar := &ReAuthRequest{}
		if err := message.Unmarshal(rar); err != nil {
			glog.Errorf("Received unparseable RAR over Gx %s\n%s", message, err)
			return
		}
		go func() {
			ans := reAuthHandler(rar)
			ansMsg := createReAuthAnswerMessage(message, ans, diamClient)
			ansMsg = diamClient.AddOriginAVPsToMessage(ansMsg)
			_, err := ansMsg.WriteToWithRetry(conn, diamClient.Retries())
			if err != nil {
				glog.Errorf(
					"Gx RAA Write Failed for %s->%s, SessionID: %s - %v",
					conn.LocalAddr(), conn.RemoteAddr(), rar.SessionID, err)
				conn.Close() // close connection on error
			}
		}()
	}
	diamClient.RegisterRequestHandlerForAppID(diam.ReAuth, diam.GX_CHARGING_CONTROL_APP_ID, handler)
}

func createReAuthAnswerMessage(
	requestMsg *diam.Message, answer *ReAuthAnswer, diamClient *diameter.Client) *diam.Message {

	ret := requestMsg.Answer(answer.ResultCode)
	ret.InsertAVP(
		diam.NewAVP(
			avp.SessionID,
			avp.Mbit,
			0,
			datatype.UTF8String(diameter.EncodeSessionID(diamClient.OriginHost(), answer.SessionID))))
	return ret
}

// createCreditControlMessage creates a base message to be used for any Credit
// Control Request message. Init will just use this, and update and terminate
// pass in extra AVPs through additionalAVPs
// Input: context.Context which has information on where to send to,
//				CreditControlRequest with relevant request info
//			  ...*diam.AVP with any AVPs to add on
// Output: *diam.Message with all AVPs filled in, error if there was an issue
func (gxClient *GxClient) createCreditControlMessage(
	request *CreditControlRequest,
	additionalAVPs ...*diam.AVP,
) (*diam.Message, error) {
	m := diameter.NewProxiableRequest(diam.CreditControl, diam.GX_CHARGING_CONTROL_APP_ID, nil)
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.GX_CHARGING_CONTROL_APP_ID))
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(request.Type))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(request.RequestNumber))
	m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(credit_control.EndUserIMSI)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(request.IMSI)),
		},
	})
	if len(request.Msisdn) > 0 {
		m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
				diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(request.Msisdn)),
			},
		})
	}
	m.NewAVP(avp.IPCANType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(5))
	m.NewAVP(avp.RATType, avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(1004))
	m.NewAVP(avp.FramedIPAddress, avp.Mbit, 0, datatype.IPv4(net.ParseIP(request.IPAddr)))
	if len(request.Apn) > 0 {
		m.NewAVP(avp.CalledStationID, avp.Mbit, 0, datatype.UTF8String(request.Apn))
	}

	if request.Type == credit_control.CRTInit {
		gxClient.getInitAvps(m, request)
	}

	if request.Type == credit_control.CRTTerminate {
		// TODO support more than DIAMETER_LOGOUT
		m.NewAVP(avp.TerminationCause, avp.Mbit, 0, datatype.Enumerated(1))
	}

	for _, avp := range additionalAVPs {
		m.InsertAVP(avp)
	}

	// SessionID must be the first AVP
	m.InsertAVP(diam.NewAVP(
		avp.SessionID,
		avp.Mbit,
		0,
		datatype.UTF8String(diameter.EncodeSessionID(gxClient.diamClient.OriginHost(), request.SessionID))))

	return m, nil
}

// init message
func (gxClient *GxClient) getInitAvps(m *diam.Message, request *CreditControlRequest) {
	m.NewAVP(avp.SupportedFeatures, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
			diam.NewAVP(avp.FeatureListID, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1)),
			// Set Bit 0 and Bit 1
			diam.NewAVP(avp.FeatureList, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(3)),
		},
	})
	// NETWORK_REQUEST_NOT_SUPPORTED(0)
	m.NewAVP(avp.NetworkRequestSupport, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(0))
	// DISABLE_OFFLINE(0)
	m.NewAVP(avp.Offline, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(0))
	// ENABLE_ONLINE(1)
	m.NewAVP(avp.Online, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(1))
	// Bearer-Usage - GENERAL(0)
	m.NewAVP(avp.BearerUsage, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(0))
	if len(request.SpgwIPV4) > 0 {
		m.NewAVP(avp.TGPPSGSNAddress, avp.Vbit, diameter.Vendor3GPP, datatype.IPv4(net.ParseIP(request.SpgwIPV4)))
		m.NewAVP(avp.TGPPGGSNAddress, avp.Vbit, diameter.Vendor3GPP, datatype.IPv4(net.ParseIP(request.SpgwIPV4)))
		m.NewAVP(avp.ANGWAddress, avp.Vbit, diameter.Vendor3GPP, datatype.Address(net.ParseIP(request.SpgwIPV4)))
		m.NewAVP(avp.AccessNetworkChargingAddress, avp.Mbit|avp.Vbit,
			diameter.Vendor3GPP, datatype.Address(net.ParseIP(request.SpgwIPV4)))

	}
	m.NewAVP(avp.TGPPSelectionMode, avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("0")) // IMEISV
	if len(request.Imei) > 0 {
		m.NewAVP(avp.UserEquipmentInfo, 0, 0, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.UserEquipmentInfoType, 0, 0, datatype.Enumerated(0)), // imeisv
				diam.NewAVP(avp.UserEquipmentInfoValue, 0, 0, datatype.OctetString(request.Imei)),
			},
		})
	}
	if len(request.PlmnID) > 0 {
		m.NewAVP(avp.TGPPSGSNMCCMNC, avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(request.PlmnID))
	}
	if len(request.UserLocation) > 0 {
		m.NewAVP(avp.TGPPUserLocationInfo, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(string(request.UserLocation)))
	}
	if len(request.GcID) > 0 {
		m.NewAVP(avp.AccessNetworkChargingIdentifierGx, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.AccessNetworkChargingIdentifierValue, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.OctetString(request.GcID)),
			},
		})
	}
	if request.Qos != nil {
		m.NewAVP(avp.QoSInformation, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.APNAggregateMaxBitrateDL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(request.Qos.ApnAggMaxBitRateDL)),
				diam.NewAVP(avp.APNAggregateMaxBitrateUL, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(request.Qos.ApnAggMaxBitRateUL)),
			},
		})

		var arpAVP *diam.AVP
		if gxClient.pcrf91Compliant {
			// PCRF is 29.212 release 9.1 compliant
			arpAVP = diam.NewAVP(avp.AllocationRetentionPriority, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.PriorityLevel, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(request.Qos.PriLevel)),
					diam.NewAVP(avp.PreemptionCapability, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.Qos.PreCapability)),
					diam.NewAVP(avp.PreemptionVulnerability, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.Qos.PreVulnerability)),
				},
			})
		} else {
			// PCRF is NOT 29.212 release 9.1 compliant
			arpAVP = diam.NewAVP(avp.AllocationRetentionPriority, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.PriorityLevel, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(request.Qos.PriLevel)),
					diam.NewAVP(avp.PreemptionCapability, avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.Qos.PreCapability)),
					diam.NewAVP(avp.PreemptionVulnerability, avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.Qos.PreVulnerability)),
				},
			})
		}
		m.NewAVP(avp.DefaultEPSBearerQoS, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.QoSClassIdentifier, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.Qos.QosClassIdentifier)),
				arpAVP,
			},
		})
	}
	// Argentina TZ (UTC-3hrs) TODO: Make it configurable
	m.NewAVP(avp.TGPPMSTimeZone, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(string([]byte{0x29, 0})))
}

// getAdditionalAvps retrieves any extra AVPs based on the type of request.
// For update and terminate, it returns the used credit AVPs
func (gxClient *GxClient) getAdditionalAvps(request *CreditControlRequest) ([]*diam.AVP, error) {
	if request.Type == credit_control.CRTInit || len(request.UsageReports) == 0 {
		return []*diam.AVP{}, nil
	}
	avpList := make([]*diam.AVP, 0, len(request.UsageReports)+1)
	for _, usage := range request.UsageReports {
		avpList = append(avpList, getUsageMonitoringAVP(usage))
	}
	if request.Type == credit_control.CRTUpdate {
		avpList = append(avpList, gxClient.getUsageReportEventTrigger())
	}

	return avpList, nil
}

func getUsageMonitoringAVP(usage *UsageReport) *diam.AVP {
	return diam.NewAVP(avp.UsageMonitoringInformation, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MonitoringKey, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(usage.MonitoringKey)),
			diam.NewAVP(avp.UsedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(usage.InputOctets)),
					diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(usage.OutputOctets)),
					diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(usage.TotalOctets)),
				},
			}),
			diam.NewAVP(avp.UsageMonitoringLevel, avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(usage.Level)),
		},
	})
}

func (gxClient *GxClient) getUsageReportEventTrigger() *diam.AVP {
	var urt = UsageReportTrigger
	if gxClient.pcrf91Compliant {
		urt = PCRF91UsageReportTrigger
	}
	return diam.NewAVP(avp.EventTrigger, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(urt))
}
