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
	"net"

	"magma/feg/cloud/go/protos"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
)

type SessionFTeids struct {
	cFTeid *ie.IE
	uFTeid *ie.IE
}

func buildCreateSessionRequestIE(req *protos.CreateSessionRequestPgw, conn *gtpv2.Conn, s8IpAddr net.Addr) ([]*ie.IE, SessionFTeids) {
	// cTEID will be managed by s8_proxy
	cFTeid := conn.NewSenderFTEID(s8IpAddr.String(), "")

	// uTEID will be given by MME (managed by MME)
	uFteidReq := req.BearerContext.GetAgwUserPlaneFteid()
	uFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPU,
		uFteidReq.Teid, uFteidReq.Ipv6Address, uFteidReq.Ipv6Address)

	// Qos
	qos := req.BearerContext.GetQos()
	ieQos := ie.NewBearerQoS(uint8(qos.Pci), uint8(qos.PriorityLevel), uint8(qos.PreemptionVulnerability),
		uint8(qos.Qci), qos.Mbr.BrUl, qos.Mbr.BrDl, qos.Gbr.BrUl, qos.Gbr.BrDl)

	bearer := ie.NewBearerContext(ie.NewEPSBearerID(uint8(req.BearerContext.Id)), uFTeid, ieQos)
	userLoc := getUserLocationIndication(req)
	pdnAllocation := getPDNAddressAllocation(req)

	// TODO: indication flag
	// TODO: set apn restriction

	return []*ie.IE{
		ie.NewIMSI(req.GetImsi()),
		userLoc,
		bearer,
		pdnAllocation,
		cFTeid,
		ie.NewMSISDN(string(req.Msisdn[:])),
		ie.NewMobileEquipmentIdentity(req.Mei),
		ie.NewServingNetwork(req.ServingNetwork.Mcc, req.ServingNetwork.Mnc),
		ie.NewRATType(uint8(req.RatType)),
		ie.NewAccessPointName(req.Apn),
		// TODO: selection mode (hadcoded for now)
		ie.NewSelectionMode(gtpv2.SelectionModeMSorNetworkProvidedAPNSubscribedVerified),
		ie.NewPDNType(uint8(req.PdnType)),
		ie.NewAggregateMaximumBitRate(uint32(req.Ambr.BrUl), uint32(req.Ambr.BrDl)),
	}, SessionFTeids{cFTeid, uFTeid}
}

func getPDNAddressAllocation(req *protos.CreateSessionRequestPgw) *ie.IE {
	var res *ie.IE
	if req.PdnType == protos.PDNType_IPV4 {
		res = ie.NewPDNAddressAllocation(req.Paa.Ipv4Address)
	}
	if req.PdnType == protos.PDNType_IPV6 {
		res = ie.NewPDNAddressAllocationIPv6(req.Paa.Ipv6Address, uint8(req.Paa.Ipv6Prefix))
	}
	if req.PdnType == protos.PDNType_IPV4V6 {
		res = ie.NewPDNAddressAllocationDual(req.Paa.Ipv4Address, req.Paa.Ipv6Address, uint8(req.Paa.Ipv6Prefix))
	}
	return res
}

func getUserLocationIndication(req *protos.CreateSessionRequestPgw) *ie.IE {
	mcc := req.ServingNetwork.Mcc
	mnc := req.ServingNetwork.Mnc
	uliReq := req.GetUli()

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

	if uliReq.Lac != 0 && uliReq.Ci != 0 {
		cgi = ie.NewCGI(mcc, mnc, uint16(uliReq.Lac), uint16(uliReq.Ci))
	}
	if uliReq.Lac != 0 && uliReq.Sac != 0 {
		sai = ie.NewSAI(mcc, mnc, uint16(uliReq.Lac), uint16(uliReq.Sac))
	}
	if uliReq.Lac != 0 && uliReq.Rac != 0 {
		rai = ie.NewRAI(mcc, mnc, uint16(uliReq.Lac), uint16(uliReq.Rac))
	}
	if uliReq.Tac != 0 {
		tai = ie.NewTAI(mcc, mnc, uint16(uliReq.Tac))
	}
	if uliReq.Eci != 0 {
		ecgi = ie.NewECGI(mcc, mnc, uliReq.Eci)
	}
	if uliReq.Lac != 0 {
		lai = ie.NewLAI(mcc, mnc, uint16(uliReq.Lac))
	}
	if uliReq.MeNbi != 0 {
		menbi = ie.NewMENBI(mcc, mnc, uliReq.MeNbi)
	}
	if uliReq.EMeNbi != 0 {
		emenbi = ie.NewEMENBI(mcc, mnc, uliReq.EMeNbi)
	}
	return ie.NewUserLocationInformationStruct(cgi, sai, rai, tai, ecgi, lai, menbi, emenbi)
}
