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

	"github.com/pkg/errors"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

type CreateBearerRequest struct {
	Imsi               string
	DedicatedBearereID uint8
	QosQCI             uint8
	ChargingID         uint32
	BiFilterProtocolId uint8
	BiFilterPort       uint16
}

func (mPgw *MockPgw) CreateBearerRequest(req CreateBearerRequest) (*message.CreateBearerResponse, error) {
	session, err := mPgw.GetSessionByIMSI(req.Imsi)
	if err != nil {
		return nil, err
	}

	sgwTeidC, err := session.GetTEID(gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		err = errors.Wrap(err, "Error, couldnt find teid con Create Bearer Request")
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

	// wait for getHandleCreateBearerRequest to process the response of sgw
	incomingMsg, err := session.WaitMessage(sequence, mPgw.GtpTimeout)
	if err != nil {
		fmt.Printf("mockPgw couldn't process received CreateBearerResponse: %s\n", err)
		return nil, err
	}
	var cbRspFromSGW *message.CreateBearerResponse
	switch m := incomingMsg.(type) {
	case *message.CreateBearerResponse:
		// move forward
		cbRspFromSGW = m
	default:
		errMsg := "mockPgw couldn't parse CreateBearerResponse"
		fmt.Println(errMsg)
		return nil, errors.New(errMsg)
	}
	fmt.Printf("mockPGW received GreateBearerResponse: %s\n", cbRspFromSGW.String())
	return cbRspFromSGW, nil
}

func buildCreateBearerRequestIEs(session *gtpv2.Session, req CreateBearerRequest) ([]*ie.IE, error) {
	ies := []*ie.IE{
		ie.NewEPSBearerID(session.GetDefaultBearer().EBI),
		ie.NewBearerContext(
			ie.NewEPSBearerID(req.DedicatedBearereID),
			buildNewBearerTFTCreateNewTFT(req),
			ie.NewBearerQoS(1, 0, 1, req.QosQCI, 0x1111111111, 0x2222222222, 0x1111111111, 0x2222222222),
			ie.NewChargingID(req.ChargingID),
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
				ie.NewTFTPFComponentSecurityParameterIndex(0xdeadbeef),
				ie.NewTFTPFComponentIPv4RemoteAddress(net.IP{127, 0, 0, 1}, net.IPMask{255, 255, 255, 0}),
				ie.NewTFTPFComponentProtocolIdentifierNextHeader(req.BiFilterProtocolId),
				ie.NewTFTPFComponentSingleLocalPort(req.BiFilterPort),
			),
			ie.NewTFTPacketFilter(
				ie.TFTPFDownlinkOnly, 1, 0,
				ie.NewTFTPFComponentProtocolIdentifierNextHeader(1),
				ie.NewTFTPFComponentLocalPortRange(20, 21),
			),
		},
		[]*ie.TFTParameter{
			ie.NewTFTParameter(ie.TFTParamIDAuthorizationToken, []byte{0xde, 0xad, 0xbe, 0xef}),
			ie.NewTFTParameter(ie.TFTParamIDFlowIdentifier, []byte{0x11, 0x11, 0x22, 0x22}),
			ie.NewTFTParameter(ie.TFTParamIDPacketFileterIdentifier, []byte{0x01, 0x02, 0x03, 0x04}),
		},
	)
}

// getHandleCreateBearerRequest just handle Create Bearer Response and return it back to
// CreateBearerRequest function so it can return it s result
func (mPgw *MockPgw) getHandleCreateBearerRequest() gtpv2.HandlerFunc {
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
