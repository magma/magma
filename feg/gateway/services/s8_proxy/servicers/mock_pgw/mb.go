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

package mock_pgw

import (
	"fmt"
	"net"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

// TODO
func (mPgw *MockPgw) getHandleModifyBearerRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		session, err := c.GetSessionByTEID(msg.TEID(), sgwAddr)
		if err != nil {
			dsr := message.NewDeleteSessionResponse(
				0, 0,
				ie.NewCause(gtpv2.CauseIMSIIMEINotKnown, 0, 0, 0, nil),
			)
			if err := c.RespondTo(sgwAddr, msg, dsr); err != nil {
				return err
			}
			return err
		}
		bearer := session.GetDefaultBearer()

		mbReqFromMME := msg.(*message.ModifyBearerRequest)
		if brCtxIE := mbReqFromMME.BearerContextsToBeModified; brCtxIE != nil {
			for _, childIE := range brCtxIE.ChildIEs {
				switch childIE.Type {
				case ie.Indication:
					// TODO:
					// do nothing in this example.
					// S-GW should change its behavior based on indication flags like;
					//  - pass Modify Bearer Request to P-GW if handover is indicated.
					//  - XXX...
				case ie.FullyQualifiedTEID:
					if err := handleFTEIDU(childIE, session, bearer); err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

func handleFTEIDU(fteiduIE *ie.IE, session *gtpv2.Session, bearer *gtpv2.Bearer) error {
	if fteiduIE.Type != ie.FullyQualifiedTEID {
		return &gtpv2.UnexpectedIEError{IEType: fteiduIE.Type}
	}
	ip, err := fteiduIE.IPAddress()
	if err != nil {
		return err
	}
	addr, err := net.ResolveUDPAddr("udp", ip+gtpv2.GTPUPort)
	if err != nil {
		fmt.Printf("Warning, Mock PGW couldnt resolve user plane ip (you may ignore it)")
	}
	bearer.SetRemoteAddress(addr)

	teid, err := fteiduIE.TEID()
	if err != nil {
		return err
	}
	bearer.SetOutgoingTEID(teid)

	it, err := fteiduIE.InterfaceType()
	if err != nil {
		return err
	}
	session.AddTEID(it, teid)
	return nil
}
