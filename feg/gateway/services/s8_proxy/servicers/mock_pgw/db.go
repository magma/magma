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

type DeleteBearerRequest struct {
	Imsi        string
	EpsBearerId uint8
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
		err = fmt.Errorf("Error, couldnt find teid on Dedicated Bearer Request: %w", err)
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
		ie.NewEPSBearerID(session.GetDefaultBearer().EBI).WithInstance(0),
		// dedicated bearers
		ie.NewEPSBearerID(req.EpsBearerId).WithInstance(1),
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
		dbResFromSGW := msg.(*message.DeleteBearerResponse)

		// check Cause value first.
		if causeIE := dbResFromSGW.Cause; causeIE != nil {
			cause, err := causeIE.Cause()
			if err != nil {
				return fmt.Errorf("Delete Bearer Request couldn't check cause of %s: %s", msg.MessageTypeName(), err)
			}
			if cause != gtpv2.CauseRequestAccepted {
				return fmt.Errorf("Delete Bearer Request not accepcted. Cause: %d", cause)
			}
		} else {
			return fmt.Errorf("Create Bearer Request has missing cause")
		}

		if linkedEBI := dbResFromSGW.LinkedEBI; linkedEBI != nil {
			lEBI, err := linkedEBI.EPSBearerID()
			if err != nil {
				return fmt.Errorf("Create Bearer Request can't parse EBI")
			}
			if session.GetDefaultBearer().EBI != lEBI {
				return fmt.Errorf(
					"Create Bearer Request LinkedEBI different than Default Bearer id")
			}
		} else {
			return fmt.Errorf("Create Bearer Request has missing linked EBI")
		}

		// collect bearer
		var dCause uint8
		var dEBI uint8
		// TODO: handle multiple bearers
		if brCtxIE := dbResFromSGW.BearerContexts[0]; brCtxIE != nil {
			for _, childIE := range brCtxIE.ChildIEs {
				switch childIE.Type {
				case ie.EPSBearerID:
					dEBI, err = childIE.EPSBearerID()
					if err != nil {
						return err
					}
				case ie.Cause:
					dCause, err = childIE.Cause()
					if err != nil {
						return err
					}
				}
			}
			if dCause == 0 || dCause != gtpv2.CauseRequestAccepted {
				return fmt.Errorf("Create Beaerer Reuqest has a bad default bearer cause %d", dCause)
			}
			if dEBI == 0 {
				return fmt.Errorf("Create Beaerer Reuqest has a bad default bearer EBI")
			}
		} else {
			return &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
		}

		// pass message to same session
		if err = gtpv2.PassMessageTo(session, msg, mPgw.GtpTimeout); err != nil {
			return fmt.Errorf("Mock PGW could not pass the DeleteBererRequest %s", err)
		}
		return nil
	}
}
