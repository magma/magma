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

package servicers

import (
	"fmt"
	"net"
	"time"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
)

// buildCreateSessionRequestIE creates a Message with all the IE needed for a Create Session Request
func buildCreateSessionRequestMsg(cPgwUDPAddr *net.UDPAddr, apnSuffix string, req *protos.CreateSessionRequestPgw) (message.Message, error) {
	// Create session needs two FTEIDs:
	// - S8 control plane FTEID will be built using local address and control TEID
	//	 passed by MME
	// - S8 user plane FTEID, provided by MME in the requested bearer

	// TODO: look for a better way to find the local ip (avoid pinging on each request)
	// (obtain the IP that is going to send the packet first)
	ip, err := gtp.GetLocalOutboundIP(cPgwUDPAddr)
	if err != nil {
		return nil, err
	}

	// Control plane TEID
	cFegFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPC,
		req.CAgwTeid, ip.String(), "").WithInstance(0)

	// User plane TEID (ip belongs to pipelined GTP-U interface)
	uAgwFTeidReq := req.BearerContext.GetUserPlaneFteid()
	uAgwFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPU,
		uAgwFTeidReq.Teid, uAgwFTeidReq.Ipv4Address, uAgwFTeidReq.Ipv6Address).WithInstance(2)

	// Qos
	qos := req.BearerContext.GetQos()
	ieQos := ie.NewBearerQoS(uint8(qos.Pci), uint8(qos.PriorityLevel), uint8(qos.PreemptionVulnerability),
		uint8(qos.Qci), qos.Mbr.BrUl, qos.Mbr.BrDl, qos.Gbr.BrUl, qos.Gbr.BrDl)

	// bearer
	bearerId := ie.NewEPSBearerID(uint8(req.BearerContext.Id))
	bearer := ie.NewBearerContext(bearerId, uAgwFTeid, ieQos)

	// APN
	apnWithSuffix := fmt.Sprintf("%s%s", req.Apn, apnSuffix)

	//timezone
	offset := time.Duration(req.TimeZone.DeltaSeconds) * time.Second
	daylightSavingTime := uint8(req.TimeZone.DaylightSavingTime)

	ies := []*ie.IE{
		ie.NewIMSI(req.GetImsi()),
		bearer,
		cFegFTeid,
		getUserLocationIndication(req.ServingNetwork, req.Uli),
		getPdnType(req.PdnType),
		getPDNAddressAllocation(req),
		getRatType(req.RatType),
		getSelectionModeType(req.SelectionMode),
		getProtocolConfigurationOptions(req.ProtocolConfigurationOptions),
		ie.NewMSISDN(req.Msisdn[:]),
		ie.NewMobileEquipmentIdentity(req.Mei),
		ie.NewServingNetwork(req.ServingNetwork.Mcc, req.ServingNetwork.Mnc),
		ie.NewAccessPointName(apnWithSuffix),
		ie.NewAggregateMaximumBitRate(uint32(req.Ambr.BrUl), uint32(req.Ambr.BrDl)),
		ie.NewUETimeZone(offset, daylightSavingTime),
		// TODO: Hardcoded values
		ie.NewIndicationFromOctets(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		ie.NewAPNRestriction(gtpv2.APNRestrictionNoExistingContextsorRestriction),
		ie.NewChargingCharacteristics(0),
	}
	msg := message.NewCreateSessionRequest(0, 0, ies...)
	return msg, nil
}

func buildDeleteSessionRequestMsg(cPgwUDPAddr *net.UDPAddr, req *protos.DeleteSessionRequestPgw) (message.Message, error) {
	// TODO: look for a better way to find the local ip (avoid pinging on each request)
	// (obtain the IP that is going to send the packet first)
	ip, err := gtp.GetLocalOutboundIP(cPgwUDPAddr)
	if err != nil {
		return nil, err
	}
	// Control plane TEID
	cFegFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPC, req.CAgwTeid, ip.String(), "").WithInstance(0)

	ies := []*ie.IE{
		ie.NewEPSBearerID(uint8(req.BearerId)),
		cFegFTeid,
		getUserLocationIndication(req.ServingNetwork, req.Uli),
	}
	return message.NewDeleteSessionRequest(req.CPgwTeid, 0, ies...), nil
}

func buildCreateBearerResMsg(seq uint32, res *protos.CreateBearerResponsePgw) (message.Message, error) {
	if res.Cause != uint32(gtpv2.CauseRequestAccepted) {
		return buildCreateBearerResWithErrorCauseMsg(res.Cause, res.CPgwTeid, seq), nil
	}
	if res.BearerContext == nil {
		return nil, fmt.Errorf("CreateBearerResponse could not be sent. Missing Bearer Contex")
	}

	// bearer
	bearerId := ie.NewEPSBearerID(uint8(res.BearerContext.Id))
	bearer := ie.NewBearerContext(bearerId)

	//timezone
	offset := time.Duration(res.TimeZone.DeltaSeconds) * time.Second
	daylightSavingTime := uint8(res.TimeZone.DaylightSavingTime)

	return message.NewCreateBearerResponse(
		res.CPgwTeid, seq,
		ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
		bearer,
		getUserLocationIndication(res.ServingNetwork, res.Uli),
		getProtocolConfigurationOptions(res.ProtocolConfigurationOptions),
		ie.NewUETimeZone(offset, daylightSavingTime),
	), nil
}

func buildCreateBearerResWithErrorCauseMsg(cause uint32, cPgwTeid uint32, seq uint32) message.Message {
	return message.NewCreateBearerResponse(
		cPgwTeid, seq, ie.NewCause(uint8(cause), 0, 0, 0, nil))
}

func getPDNAddressAllocation(req *protos.CreateSessionRequestPgw) *ie.IE {
	var (
		res        *ie.IE
		ipv4       string
		ipv6       string
		ipv6Prefix uint8
	)
	// extract ips of default values
	if req.Paa == nil || req.Paa.Ipv4Address == "" {
		ipv4 = "0.0.0.0"
	} else {
		ipv4 = req.Paa.Ipv4Address
	}

	if req.Paa == nil || req.Paa.Ipv6Address == "" {
		ipv6 = "::"
		ipv6Prefix = 0
	} else {
		ipv6 = req.Paa.Ipv6Address
		ipv6Prefix = uint8(req.Paa.Ipv6Prefix)
	}

	// create the IE based on the type
	if req.PdnType == protos.PDNType_IPV4 {
		res = ie.NewPDNAddressAllocation(ipv4)
	}
	if req.PdnType == protos.PDNType_IPV6 {
		res = ie.NewPDNAddressAllocationIPv6(ipv6, ipv6Prefix)
	}
	if req.PdnType == protos.PDNType_IPV4V6 {
		res = ie.NewPDNAddressAllocationDual(ipv4, ipv6, ipv6Prefix)
	}
	return res
}

// getPdnType convert proto PDNType into GTP PDN type
func getPdnType(pdnType protos.PDNType) *ie.IE {
	var res = uint8(0)
	switch pdnType {
	case protos.PDNType_IPV4:
		res = gtpv2.PDNTypeIPv4 // v4
	case protos.PDNType_IPV6:
		res = gtpv2.PDNTypeIPv6 // v6
	case protos.PDNType_IPV4V6:
		res = gtpv2.PDNTypeIPv4 // v4v6
	case protos.PDNType_NonIP:
		res = gtpv2.PDNTypeNonIP // nonIP
	default:
		panic(fmt.Sprintf("PdnType %d does not exist", pdnType))
	}
	return ie.NewPDNType(res)
}

func getUserLocationIndication(servingNetwork *protos.ServingNetwork, uli *protos.UserLocationInformation) *ie.IE {
	var (
		cgi    *ie.CGI    = nil
		sai    *ie.SAI    = nil
		rai    *ie.RAI    = nil
		tai    *ie.TAI    = nil
		ecgi   *ie.ECGI   = nil
		lai    *ie.LAI    = nil
		menbi  *ie.MENBI  = nil
		emenbi *ie.EMENBI = nil
	)

	mcc := servingNetwork.Mcc
	mnc := servingNetwork.Mnc

	if uli.Lac != 0 && uli.Ci != 0 {
		cgi = ie.NewCGI(mcc, mnc, uint16(uli.Lac), uint16(uli.Ci))
	}
	if uli.Lac != 0 && uli.Sac != 0 {
		sai = ie.NewSAI(mcc, mnc, uint16(uli.Lac), uint16(uli.Sac))
	}
	if uli.Lac != 0 && uli.Rac != 0 {
		rai = ie.NewRAI(mcc, mnc, uint16(uli.Lac), uint16(uli.Rac))
	}
	if uli.Tac != 0 {
		tai = ie.NewTAI(mcc, mnc, uint16(uli.Tac))
	}
	if uli.Eci != 0 {
		ecgi = ie.NewECGI(mcc, mnc, uli.Eci)
	}
	if uli.Lac != 0 {
		lai = ie.NewLAI(mcc, mnc, uint16(uli.Lac))
	}
	if uli.MeNbi != 0 {
		menbi = ie.NewMENBI(mcc, mnc, uli.MeNbi)
	}
	if uli.EMeNbi != 0 {
		emenbi = ie.NewEMENBI(mcc, mnc, uli.EMeNbi)
	}
	return ie.NewUserLocationInformationStruct(cgi, sai, rai, tai, ecgi, lai, menbi, emenbi)
}

func getRatType(ratType protos.RATType) *ie.IE {
	var rType uint8
	switch ratType {
	case protos.RATType_RESERVED:
		rType = 0
	case protos.RATType_UTRAN:
		rType = gtpv2.RATTypeUTRAN
	case protos.RATType_GERAN:
		rType = gtpv2.RATTypeGERAN
	case protos.RATType_WLAN:
		rType = gtpv2.RATTypeWLAN
	case protos.RATType_GAN:
		rType = gtpv2.RATTypeGAN
	case protos.RATType_HSPA:
		rType = gtpv2.RATTypeHSPAEvolution
	case protos.RATType_EUTRAN:
		rType = gtpv2.RATTypeEUTRAN
	case protos.RATType_VIRTUAL:
		rType = gtpv2.RATTypeVirtual
	case protos.RATType_EUTRAN_NB_IOT:
		rType = gtpv2.RATTypeEUTRANNBIoT
	case protos.RATType_LTE_M:
		rType = gtpv2.RATTypeLTEM
	case protos.RATType_NR:
		rType = gtpv2.RATTypeNR
	default:
		panic(fmt.Sprintf("RatType %d does not exist", ratType))
	}
	return ie.NewRATType(rType)
}

func getSelectionModeType(selMode protos.SelectionModeType) *ie.IE {
	var rType uint8
	switch selMode {
	case protos.SelectionModeType_APN_provided_subscription_verified:
		rType = gtpv2.SelectionModeMSorNetworkProvidedAPNSubscribedVerified
	case protos.SelectionModeType_ms_APN_subscription_not_verified:
		rType = gtpv2.SelectionModeMSProvidedAPNSubscriptionNotVerified
	case protos.SelectionModeType_network_APN_subscription_not_verified:
		rType = gtpv2.SelectionModeNetworkProvidedAPNSubscriptionNotVerified
	default:
		panic(fmt.Sprintf("RatType %d does not exist", selMode))
	}
	return ie.NewSelectionMode(rType)
}

func getProtocolConfigurationOptions(pco *protos.ProtocolConfigurationOptions) *ie.IE {
	if pco == nil {
		return nil
	}
	var options []*ie.PCOContainer
	for _, container := range pco.ProtoOrContainerId {
		options = append(options, ie.NewPCOContainer(uint16(container.Id), container.Contents))
	}
	return ie.NewProtocolConfigurationOptions(uint8(pco.ConfigProtocol), options...)
}
