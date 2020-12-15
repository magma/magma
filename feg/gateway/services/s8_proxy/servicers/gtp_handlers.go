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

//gtp_handlers contains the handlers that will take care of messages received by the gtp server

package servicers

import (
	"fmt"
	"github.com/golang/glog"
	"magma/feg/gateway/gtp/enriched_message"
	"net"
	"time"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp"
)

func addS8GtpHandlers(c *gtp.Client) {
	c.AddHandlers(
		map[uint8]gtpv2.HandlerFunc{
			message.MsgTypeCreateSessionResponse: getHandle_CreateSessionResponse(),
			message.MsgTypeModifyBearerRequest:   getHandle_ModifyBearerRequest(),
			message.MsgTypeDeleteSessionResponse: getHandle_DeleteSessionResponse(),
			message.MsgTypeDeleteBearerRequest:   getHandle_DeleteBearerRequest(),
		})
}

func getHandle_CreateSessionResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {

		session, err := c.GetSessionByTEID(msg.TEID(), senderAddr)

		if err != nil {
			return fmt.Errorf("couldn't find session with TEID %d: %s", msg.TEID(), err)
		}

		csResGtp := msg.(*message.CreateSessionResponse)
		csRes := &protos.CreateSessionResponsePgw{}
		glog.V(2).Infof("Received Create Session Response:\n%s", csResGtp.String())

		// check Cause value first.
		if causeIE := csResGtp.Cause; causeIE != nil {
			cause, err := causeIE.Cause()
			if err != nil {
				return fmt.Errorf("Couldn't check cause of csRes: %s", err)
			}
			if cause != gtpv2.CauseRequestAccepted {
				c.RemoveSession(session)
				return &gtpv2.CauseNotOKError{
					MsgType: csResGtp.MessageTypeName(),
					Cause:   cause,
					Msg:     fmt.Sprintf("subscriber: %s", session.IMSI),
				}
			}
		} else {
			c.RemoveSession(session)
			return &gtpv2.RequiredIEMissingError{
				Type: ie.Cause,
			}
		}

		// TODO: remove this, this is just for GTP-U
		// get values sent by pgw
		bearer := session.GetDefaultBearer()
		if paaIE := csResGtp.PAA; paaIE != nil {
			ip, err := paaIE.IPAddress()
			if err != nil {
				return err
			}
			bearer.SubscriberIP = ip
		} else {
			c.RemoveSession(session)
			return &gtpv2.RequiredIEMissingError{Type: ie.PDNAddressAllocation}
		}

		if fteidcIE := csResGtp.PGWS5S8FTEIDC; fteidcIE != nil {
			it, err := fteidcIE.InterfaceType()
			if err != nil {
				return err
			}
			teid, err := fteidcIE.TEID()
			if err != nil {
				return err
			}
			session.AddTEID(it, teid)
		} else {
			c.RemoveSession(session)
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		if brCtxIE := csResGtp.BearerContextsCreated; brCtxIE != nil {
			for _, childIE := range brCtxIE.ChildIEs {
				switch childIE.Type {
				case ie.Cause:
					cause, err := childIE.Cause()
					if err != nil {
						return err
					}
					if cause != gtpv2.CauseRequestAccepted {
						c.RemoveSession(session)
						return &gtpv2.CauseNotOKError{
							MsgType: csResGtp.MessageTypeName(),
							Cause:   cause,
							Msg:     fmt.Sprintf("subscriber: %s", session.IMSI),
						}
					}
				case ie.EPSBearerID:
					ebi, err := childIE.EPSBearerID()
					if err != nil {
						return err
					}
					bearer.EBI = ebi
				case ie.FullyQualifiedTEID:
					if err := handleFTEIDU(childIE, session, bearer); err != nil {
						return err
					}
				case ie.ChargingID:
					cid, err := childIE.ChargingID()
					if err != nil {
						return err
					}
					bearer.ChargingID = cid
				}
			}
		} else {
			c.RemoveSession(session)
			return &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
		}

		if err := session.Activate(); err != nil {
			c.RemoveSession(session)
			return fmt.Errorf("couldn't activate the session with IMSI %s: %s", session.IMSI, err)
		}
		// TODO: validate message before passing
		enrichedMsg := enriched_message.NewMessageWithGrpc(msg, csRes)

		// pass message to same session
		if err := gtpv2.PassMessageTo(session, enrichedMsg, 5*time.Second); err != nil {
			return err
		}
		return nil
	}
}

// TODO
func getHandle_ModifyBearerRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		return nil
	}
}

// TODO
func getHandle_DeleteSessionResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		return nil
	}
}

// TODO
func getHandle_DeleteBearerRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
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
		return err
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
