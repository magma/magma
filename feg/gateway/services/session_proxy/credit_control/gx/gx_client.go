/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gx

import (
	"math/rand"
	"net"
	"os"
	"time"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/gateway/service_registry"
	"magma/orc8r/lib/go/util"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
)

const (
	defaultFramedIpv4Addr = "10.10.10.10"
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
	EnableConnections() error
	DisableConnections(period time.Duration)
}

// GxClient is a client to send Gx Credit Control Request messages over diameter
// And receive Gx Credit Control Answer messages in response
// Although Gy and Gx both send Credit Control Requests, their Application IDs,
// allowed AVPs, and purposes are different
type GxClient struct {
	diamClient             *diameter.Client
	serverCfg              *diameter.DiameterServerConfig
	pcrf91Compliant        bool // to support PCRF which is 29.212 release 9.1 compliant
	dontUseEUIIpIfEmpty    bool // Disable using MAC derived EUI-64 IPv6 address for CCR if IP is not provided
	framedIpv4AddrRequired bool // PCRF requires FramedIpv4Addr to be included
	globalConfig           *GxGlobalConfig
}

type GxGlobalConfig struct {
	PCFROverwriteApn string
}

// NewConnectedGxClient contructs a new GxClient with the magma diameter settings
func NewConnectedGxClient(
	diamClient *diameter.Client,
	serverCfg *diameter.DiameterServerConfig,
	reAuthHandler PolicyReAuthHandler,
	cloudRegistry service_registry.GatewayRegistry,
	gxGlobalConfig *GxGlobalConfig,
) *GxClient {
	diamClient.RegisterAnswerHandlerForAppID(diam.CreditControl, diam.GX_CHARGING_CONTROL_APP_ID, ccaHandler)
	registerReAuthHandler(reAuthHandler, diamClient)
	if cloudRegistry != nil {
		diamClient.RegisterHandler(
			diam.AbortSession,
			diam.GX_CHARGING_CONTROL_APP_ID,
			true,
			credit_control.NewASRHandler(diamClient, cloudRegistry))
	}
	return &GxClient{
		diamClient:             diamClient,
		serverCfg:              serverCfg,
		pcrf91Compliant:        *pcrf91Compliant || util.IsTruthyEnv(PCRF91CompliantEnv),
		dontUseEUIIpIfEmpty:    *disableEUIIpIfEmpty || util.IsTruthyEnv(DisableEUIIPv6IfNoIPEnv),
		framedIpv4AddrRequired: util.IsTruthyEnv(FramedIPv4AddrRequiredEnv),
		globalConfig:           gxGlobalConfig,
	}

}

// NewGxClient contructs a new GxClient with the magma diameter settings
func NewGxClient(
	clientCfg *diameter.DiameterClientConfig,
	serverCfg *diameter.DiameterServerConfig,
	reAuthHandler PolicyReAuthHandler,
	cloudRegistry service_registry.GatewayRegistry,
	globalConfig *GxGlobalConfig,
) *GxClient {
	diamClient := diameter.NewClient(clientCfg)
	diamClient.BeginConnection(serverCfg)
	return NewConnectedGxClient(diamClient, serverCfg, reAuthHandler, cloudRegistry, globalConfig)
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

	message, err := gxClient.createCreditControlMessage(request, gxClient.globalConfig, additionalAVPs...)
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

func (gxClient *GxClient) EnableConnections() error {
	gxClient.diamClient.EnableConnectionCreation()
	return gxClient.diamClient.BeginConnection(gxClient.serverCfg)
}

func (gxClient *GxClient) DisableConnections(period time.Duration) {
	gxClient.diamClient.DisableConnectionCreation(period)
}

// Register reauth request handler
func registerReAuthHandler(reAuthHandler PolicyReAuthHandler, diamClient *diameter.Client) {
	reqHandler := func(conn diam.Conn, message *diam.Message) {
		rar := &PolicyReAuthRequest{}
		if err := message.Unmarshal(rar); err != nil {
			glog.Errorf("Received unparseable RAR over Gx %s\n%s", message, err)
			return
		}
		go func() {
			raa := reAuthHandler(rar)
			raaMsg := createReAuthAnswerMessage(message, raa, diamClient)
			raaMsg = diamClient.AddOriginAVPsToMessage(raaMsg)
			_, err := raaMsg.WriteToWithRetry(conn, diamClient.Retries())
			if err != nil {
				glog.Errorf(
					"Gx RAA Write Failed for %s->%s, SessionID: %s - %v",
					conn.LocalAddr(), conn.RemoteAddr(), rar.SessionID, err)
				conn.Close() // close connection on error
			}
		}()
	}
	diamClient.RegisterRequestHandlerForAppID(diam.ReAuth, diam.GX_CHARGING_CONTROL_APP_ID, reqHandler)
}

func createReAuthAnswerMessage(
	requestMsg *diam.Message, answer *PolicyReAuthAnswer, diamClient *diameter.Client) *diam.Message {
	ret := requestMsg.Answer(answer.ResultCode)
	ret.InsertAVP(
		diam.NewAVP(
			avp.SessionID,
			avp.Mbit,
			0,
			datatype.UTF8String(diameter.EncodeSessionID(diamClient.OriginRealm(), answer.SessionID))))
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
	globalConfig *GxGlobalConfig,
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

	m.NewAVP(avp.IPCANType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.IPCANType))
	m.NewAVP(avp.RATType, avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(request.RATType))

	if ip := net.ParseIP(request.IPAddr); ipNotZeros(ip) {
		if ipV4 := ip.To4(); ipV4 != nil {
			m.NewAVP(avp.FramedIPAddress, avp.Mbit, 0, datatype.IPv4(ipV4))
		} else if gxClient.framedIpv4AddrRequired {
			defaultIp := getDefaultFramedIpv4Addr()
			m.NewAVP(avp.FramedIPAddress, avp.Mbit, 0, datatype.IPv4(defaultIp))
		} else if ipV6 := ip.To16(); ipV6 != nil {
			m.NewAVP(avp.FramedIPv6Prefix, avp.Mbit, 0, datatype.OctetString(ipV6))
		}
	} else if gxClient.framedIpv4AddrRequired {
		defaultIp := getDefaultFramedIpv4Addr()
		m.NewAVP(avp.FramedIPAddress, avp.Mbit, 0, datatype.IPv4(defaultIp))
	} else if (!gxClient.dontUseEUIIpIfEmpty) && len(request.HardwareAddr) >= 6 {
		m.NewAVP(avp.FramedIPv6Prefix, avp.Mbit, 0, datatype.OctetString(Ipv6PrefixFromMAC(request.HardwareAddr)))
	}

	apn := datatype.UTF8String(request.Apn)
	if globalConfig != nil && len(globalConfig.PCFROverwriteApn) > 0 {
		apn = datatype.UTF8String(globalConfig.PCFROverwriteApn)
	}
	if len(apn) > 0 {
		m.NewAVP(avp.CalledStationID, avp.Mbit, 0, apn)
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
		datatype.UTF8String(diameter.EncodeSessionID(gxClient.diamClient.OriginRealm(), request.SessionID))))

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
	if request.RATType != credit_control.RAT_WLAN {
		// Bearer-Usage - GENERAL(0)
		m.NewAVP(avp.BearerUsage, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(0))
		m.NewAVP(avp.TGPPSelectionMode, avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("0")) // IMEISV
	}
	if len(request.SpgwIPV4) > 0 {
		m.NewAVP(avp.TGPPSGSNAddress, avp.Vbit, diameter.Vendor3GPP, datatype.IPv4(net.ParseIP(request.SpgwIPV4)))
		m.NewAVP(avp.TGPPGGSNAddress, avp.Vbit, diameter.Vendor3GPP, datatype.IPv4(net.ParseIP(request.SpgwIPV4)))
		m.NewAVP(avp.ANGWAddress, avp.Vbit, diameter.Vendor3GPP, datatype.Address(net.ParseIP(request.SpgwIPV4)))
		m.NewAVP(avp.AccessNetworkChargingAddress, avp.Mbit|avp.Vbit,
			diameter.Vendor3GPP, datatype.Address(net.ParseIP(request.SpgwIPV4)))

	}

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
	// Argentina TZ (UTC-3hrs) TODO: Make it so that it takes the FeG's timezone
	//m.NewAVP(avp.TGPPMSTimeZone, avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(string([]byte{0x29, 0})))
}

// getAdditionalAvps retrieves any extra AVPs based on the type of request.
// For update and terminate, it returns the used credit AVPs
func (gxClient *GxClient) getAdditionalAvps(request *CreditControlRequest) ([]*diam.AVP, error) {
	if request.Type == credit_control.CRTInit {
		return []*diam.AVP{}, nil
	}
	avpList := []*diam.AVP{}
	if len(request.UsageReports) > 0 {
		avpList = make([]*diam.AVP, 0, len(request.UsageReports)+2)
		if len(request.TgppCtx.GetGxDestHost()) > 0 {
			avpList = append(avpList,
				diam.NewAVP(avp.DestinationHost, avp.Mbit, 0, datatype.DiameterIdentity(request.TgppCtx.GetGxDestHost())))
		}
		for _, usage := range request.UsageReports {
			avpList = append(avpList, getUsageMonitoringAVP(usage))
		}
	}
	if request.Type == credit_control.CRTUpdate {
		avpList = append(avpList, gxClient.getEventTriggerAVP(request.EventTrigger))
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

func (gxClient *GxClient) getEventTriggerAVP(eventTrigger EventTrigger) *diam.AVP {
	if eventTrigger == UsageReportTrigger {
		if gxClient.pcrf91Compliant {
			eventTrigger = PCRF91UsageReportTrigger
		}
	}
	return diam.NewAVP(avp.EventTrigger, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(eventTrigger))
}

// Is p all zeros?
func ipNotZeros(p net.IP) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return true
		}
	}
	return false
}

var (
	prefix = []byte{0, 0x80, 0xfd, 0xfa, 0xce, 0xb0, 0x0c, 0xab, 0xcd, 0xef}
	_, _   = rand.Read(prefix[6:])
)

// Ipv6PrefixFromMAC creates a unique local EUI-64 based IPv6 address from given MAC address
// see: https://www.rfc-editor.org/rfc/rfc4193.html
func Ipv6PrefixFromMAC(mac net.HardwareAddr) []byte {
	ip := make([]byte, net.IPv6len+2)
	// Copy prefix directly into first 8 bytes of IP address
	copy(ip[0:10], prefix)

	// If MAC is in EUI-48 form, split first three bytes and last three bytes,
	// and inject 0xff and 0xfe between them
	if len(mac) == 6 {
		copy(ip[10:13], mac[0:3])
		// Flip 7th bit
		ip[10] ^= 0x02
		ip[13] = 0xff
		ip[14] = 0xfe
		copy(ip[15:18], mac[3:6])
	} else if len(mac) == 8 {
		// If MAC is in EUI-64 form, directly copy it into output IP address
		copy(ip[10:18], mac)
		// Flip 7th bit
		ip[10] ^= 0x02
	}
	return ip
}

func getDefaultFramedIpv4Addr() net.IP {
	ip := os.Getenv(DefaultFramedIPv4AddrEnv)
	if len(ip) == 0 {
		ip = defaultFramedIpv4Addr
	}
	ipV4V6 := net.ParseIP(ip)
	if ipV4 := ipV4V6.To4(); ipV4 != nil {
		return ipV4
	}
	return ipV4V6
}
