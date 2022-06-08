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
	"errors"
	"fmt"
	"net"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

type DedicatedBearerContext struct {
	DedicatedBearereID uint8
	Pgw_u_ip           string
	Pgw_u_teid         uint32
}

type CreateBearerRequest struct {
	Imsi               string
	QosQCI             uint8
	ChargingID         uint32
	BiFilterProtocolId uint8
	BiLocalFilterPort  uint16
	BiRemoteFilterPort uint16
	BearerContext      DedicatedBearerContext
}

type CBReq struct {
	Res *message.CreateBearerResponse
	Err error
}

func (mPgw *MockPgw) CreateBearerRequest(req CreateBearerRequest) (chan CBReq, error) {
	session, err := mPgw.GetSessionByIMSI(req.Imsi)
	if err != nil {
		return nil, err
	}

	sgwTeidC, err := session.GetTEID(gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		err = fmt.Errorf("Error, couldnt find teid con Create Bearer Request: %w", err)
		return nil, err
	}

	ies, err := buildCreateBearerRequestIEs(session, req)
	if err != nil {
		return nil, err
	}

	cbReqMsg := message.NewCreateBearerRequest(sgwTeidC, 0, ies...)

	sequence, err := mPgw.SendMessageTo(cbReqMsg, session.PeerAddr())
	if err != nil {
		return nil, err
	}
	mPgw.LastSequenceNumber = sequence

	out := make(chan CBReq)
	// this routine is needed due to the fact AGW req/res is split into two grpc servers
	go func() {
		// wait for getHandleCreateBearerRequest to process the response of sgw
		incomingMsg, err := session.WaitMessage(sequence, mPgw.GtpTimeout)
		if err != nil {
			fmt.Printf("mockPgw couldn't process received CreateBearerResponse: %s\n", err)
			out <- CBReq{Err: err}
			return
		}
		var cbRspFromSGW *message.CreateBearerResponse
		switch m := incomingMsg.(type) {
		case *message.CreateBearerResponse:
			// move forward
			cbRspFromSGW = m
		default:
			errMsg := "mockPgw couldn't parse CreateBearerResponse"
			fmt.Println(errMsg)
			out <- CBReq{Err: errors.New(errMsg)}
			return
		}
		fmt.Printf("mockPGW received CreateBearerResponse: %s\n", cbRspFromSGW.String())

		out <- CBReq{
			Res: cbRspFromSGW,
			Err: nil,
		}
	}()
	return out, nil
}

func buildCreateBearerRequestIEs(session *gtpv2.Session, req CreateBearerRequest) ([]*ie.IE, error) {
	pgwFTEIDu := ie.NewFullyQualifiedTEID(
		gtpv2.IFTypeS5S8PGWGTPU, req.BearerContext.Pgw_u_teid, req.BearerContext.Pgw_u_ip, "").WithInstance(1)

	ies := []*ie.IE{
		ie.NewEPSBearerID(session.GetDefaultBearer().EBI),
		ie.NewBearerContext(
			ie.NewEPSBearerID(req.BearerContext.DedicatedBearereID),
			buildNewBearerTFTCreateNewTFT(req),
			ie.NewBearerQoS(1, 0, 1, req.QosQCI, 0x1111111111, 0x2222222222, 0x1111111111, 0x2222222222),
			ie.NewChargingID(req.ChargingID),
			pgwFTEIDu,
		),
	}
	return ies, nil
}

// buildNewBearerTFTCreateNewTFT creates a default TFT
func buildNewBearerTFTCreateNewTFT(req CreateBearerRequest) *ie.IE {

	return ie.NewBearerTFTCreateNewTFT(
		[]*ie.TFTPacketFilter{
			ie.NewTFTPacketFilter(
				ie.TFTPFBidirectional, 0, 0,
				// component 0.0
				ie.NewTFTPFComponentSecurityParameterIndex(0xdeadbeef),
				// component 0.1
				ie.NewTFTPFComponentIPv4RemoteAddress(net.IP{127, 0, 0, 1}, net.IPMask{255, 255, 255, 0}),
				// component 0.2
				ie.NewTFTPFComponentProtocolIdentifierNextHeader(req.BiFilterProtocolId),
				// component 0.3
				ie.NewTFTPFComponentTypeOfServiceTrafficClass(1, 2),
				// component 0.4
				ie.NewTFTPFComponentSingleLocalPort(req.BiLocalFilterPort),
				// component 0.5
				ie.NewTFTPFComponentSingleRemotePort(req.BiRemoteFilterPort),
			),
			ie.NewTFTPacketFilter(
				ie.TFTPFDownlinkOnly, 1, 0,
				// component 1.0
				ie.NewTFTPFComponentProtocolIdentifierNextHeader(1),
				// component 1.1
				ie.NewTFTPFComponentSecurityParameterIndex(0xdeadbeef),
				// component 1.2
				ie.NewTFTPFComponentLocalPortRange(req.BiLocalFilterPort, req.BiLocalFilterPort+10),
				// component 1.3
				ie.NewTFTPFComponentRemotePortRange(req.BiRemoteFilterPort, req.BiRemoteFilterPort+10),
			),
		},
		[]*ie.TFTParameter{
			ie.NewTFTParameter(ie.TFTParamIDAuthorizationToken, []byte{0xde, 0xad, 0xbe, 0xef}),
			ie.NewTFTParameter(ie.TFTParamIDFlowIdentifier, []byte{0x11, 0x11, 0x22, 0x22}),
			ie.NewTFTParameter(ie.TFTParamIDPacketFileterIdentifier, []byte{0x01, 0x02, 0x03, 0x04}),
		},
	)
}

// getHandleCreateBearerResponse just handle Create Bearer Response and return it back to
// CreateBearerRequest function so it can return it s result
func (mPgw *MockPgw) getHandleCreateBearerResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Println("mock PGW received a CreateBearerResponse")
		session, err := c.GetSessionByTEID(msg.TEID(), sgwAddr)
		if err != nil {
			return fmt.Errorf("Mock PGW could not hanlde CreateBearerRequest: %s", err)
		}
		// pass message to same session
		if err = gtpv2.PassMessageTo(session, msg, mPgw.GtpTimeout); err != nil {
			return fmt.Errorf("Mock PGW could not pass the CreateBererRequest %s", err)
		}
		return nil
	}
}
