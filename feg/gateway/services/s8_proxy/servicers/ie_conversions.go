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
	uFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8SGWGTPU, req.UserPlaneTeid, req.S5S8Ip4UserPane, req.S5S8Ip6UserPane)
	// TODO: add proper to bearer
	ieQos := ie.NewBearerQoS(1, 2, 1, 0xff, 0, 0, 0, 0)

	// TODO: indication flag
	// TODO: set uli
	// TODO: set apn restriction

	return []*ie.IE{
		ie.NewIMSI(req.GetSid().Id),
		ie.NewMSISDN(string(req.MSISDN[:])),
		ie.NewMobileEquipmentIdentity(req.MEI),
		ie.NewServingNetwork(req.MCC, req.MNC),
		ie.NewRATType(uint8(req.RatType)),
		cFTeid,
		ie.NewAccessPointName(req.Apn),
		// TODO: selection mode
		ie.NewSelectionMode(gtpv2.SelectionModeMSorNetworkProvidedAPNSubscribedVerified),
		ie.NewPDNType(uint8(req.PdnType)),
		ie.NewPDNAddressAllocation(req.PdnAddressAllocation),
		ie.NewAggregateMaximumBitRate(req.AmbrUp, req.AmbrDown),
		ie.NewBearerContext(ie.NewEPSBearerID(uint8(req.BearerId)), uFTeid, ieQos),
	}, SessionFTeids{cFTeid, uFTeid}
}
