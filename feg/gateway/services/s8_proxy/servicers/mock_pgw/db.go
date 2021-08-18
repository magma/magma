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

type DeleteBearerRequest struct {
	Imsi           string
	LinkedBearerId uint8
}

type DBReq struct {
	Res *message.DeleteBearerResponse
	Err error
}

func (mPgw *MockPgw) DeleteBearerRequest(req DeleteBearerRequest) (chan DBReq, error) {
	session, err := mPgw.GetSessionByIMSI(req.Imsi)
	if err != nil {
		return nil, err
	}

	sgwTeidC, err := session.GetTEID(gtpv2.IFTypeS5S8SGWGTPC)
	if err != nil {
		err = errors.Wrap(err, "Error, couldnt find teid on Dedicated Bearer Request")
		return nil, err
	}

	ies, err := buildDeleteBearerRequestIEs(session, req)
	if err != nil {
		return nil, err
	}

	cbReqMsg := message.NewDeleteBearerRequest(sgwTeidC, 0, ies...)

	sequence, err := mPgw.SendMessageTo(cbReqMsg, session.PeerAddr())
	if err != nil {
		return nil, err
	}
	mPgw.LastSequenceNumber = sequence

	out := make(chan DBReq)
	// this routine is needed due to the fact AGW req/res is split into two grpc servers
	go func() {
		// wait for getHandleDeleteBearerRequest to process the response of sgw
		incomingMsg, err := session.WaitMessage(sequence, mPgw.GtpTimeout)
		if err != nil {
			fmt.Printf("mockPgw couldn't process received DeleteBearerResponse: %s\n", err)
			out <- DBReq{Err: err}
			return
		}
		var dbRspFromSGW *message.DeleteBearerResponse
		switch m := incomingMsg.(type) {
		case *message.DeleteBearerResponse:
			// move forward
			dbRspFromSGW = m
		default:
			errMsg := "mockPgw couldn't parse DeleteBearerResponse"
			fmt.Println(errMsg)
			out <- DBReq{Err: errors.New(errMsg)}
			return
		}
		fmt.Printf("mockPGW received DeleteBearerResponse: %s\n", dbRspFromSGW.String())
		out <- DBReq{
			Res: dbRspFromSGW,
			Err: nil,
		}
	}()
	return out, nil
}

func buildDeleteBearerRequestIEs(session *gtpv2.Session, req DeleteBearerRequest) ([]*ie.IE, error) {
	ies := []*ie.IE{
		ie.NewEPSBearerID(req.LinkedBearerId).WithInstance(0),
	}
	return ies, nil
}

// getHandleDeleteBearerRequest just handle Delete Bearer Response and return it back to
// DeleteBearerRequest function so it can return it s result
func (mPgw *MockPgw) getHandleDeleteBearerResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Println("mock PGW received a DeleteBearerResponse")
		session, err := c.GetSessionByTEID(msg.TEID(), sgwAddr)
		if err != nil {
			return fmt.Errorf("Mock PGW could not hanlde DeleteBearerRequest: %s", err)
		}
		// pass message to same session
		if err = gtpv2.PassMessageTo(session, msg, mPgw.GtpTimeout); err != nil {
			return fmt.Errorf("Mock PGW could not pass the DeleteBererRequest %s", err)
		}
		return nil
	}
}
