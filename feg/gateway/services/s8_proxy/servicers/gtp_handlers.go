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
	"net"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp/enriched_message"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func addS8GtpHandlers(s8p *S8Proxy) {
	s8p.gtpClient.AddHandlers(
		map[uint8]gtpv2.HandlerFunc{
			message.MsgTypeCreateSessionResponse: getHandle_CreateSessionResponse(),
			message.MsgTypeModifyBearerRequest:   getHandle_ModifyBearerRequest(),
			message.MsgTypeDeleteSessionResponse: getHandle_DeleteSessionResponse(),
			message.MsgTypeDeleteBearerRequest:   getHandle_DeleteBearerRequest(),
			message.MsgTypeEchoResponse:          getHandle_EchoResponse(s8p.echoChannel),
		})
}

func getHandle_CreateSessionResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		csResGtp := msg.(*message.CreateSessionResponse)
		csRes := &protos.CreateSessionResponsePgw{}
		glog.V(2).Infof("Received Create Session Response (gtp):\n%s", csResGtp.String())

		session, err := c.GetSessionByTEID(msg.TEID(), senderAddr)
		if err != nil {
			return fmt.Errorf("couldn't find session with TEID %d: %s", msg.TEID(), err)
		}

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

		// get values sent by pgw
		if paaIE := csResGtp.PAA; paaIE != nil {
			ip, err := paaIE.IPAddress()
			if err != nil {
				return err
			}
			csRes.SubscriberIp = ip
		} else {
			c.RemoveSession(session)
			return &gtpv2.RequiredIEMissingError{Type: ie.PDNAddressAllocation}
		}

		// control plane fteid
		if fteidcIE := csResGtp.PGWS5S8FTEIDC; fteidcIE != nil {
			fteidc, interfaceType, err := handleFTEID(fteidcIE)
			if err != nil {
				return err
			}
			session.AddTEID(interfaceType, fteidc.GetTeid())
		} else {
			c.RemoveSession(session)
			return &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		}

		// TODO: handle more than one bearer
		if brCtxIE := csResGtp.BearerContextsCreated; brCtxIE != nil {
			bearerCtx := &protos.BearerContext{}
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
					if ebi != session.GetDefaultBearer().EBI {
						return fmt.Errorf("Create Session Response bearer id different than "+
							"default bearer id (%d != %d)", ebi, session.GetDefaultBearer().EBI)
					}
					bearerCtx.Id = uint32(ebi)
				case ie.FullyQualifiedTEID:
					uFteid, typeIf, err := handleFTEID(childIE)
					if err != nil {
						return err
					}
					bearerCtx.UserPlaneFteid = uFteid
					// save uFteid in session and default bearer
					session.AddTEID(typeIf, uFteid.GetTeid())
					session.GetDefaultBearer().SetOutgoingTEID(uFteid.GetTeid())
				case ie.ChargingID:
					bearerCtx.ChargingId, err = childIE.ChargingID()
					if err != nil {
						return err
					}
					session.GetDefaultBearer().ChargingID = bearerCtx.ChargingId
				}
			}
			csRes.BearerContext = bearerCtx
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

func getHandle_DeleteSessionResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		cdResGtp := msg.(*message.DeleteSessionResponse)
		cdRes := &protos.DeleteSessionResponsePgw{}
		glog.V(2).Infof("Received Delete Session Response (gtp):\n%s", cdResGtp.String())

		session, err := c.GetSessionByTEID(msg.TEID(), senderAddr)
		if err != nil {
			return fmt.Errorf("couldn't find session with TEID %d: %s", msg.TEID(), err)
		}

		// check Cause value first.
		if causeIE := cdResGtp.Cause; causeIE != nil {
			cause, err := causeIE.Cause()
			if err != nil {
				return fmt.Errorf("Couldn't check cause of delete session response: %s", err)
			}
			if cause != gtpv2.CauseRequestAccepted {
				return &gtpv2.CauseNotOKError{
					MsgType: cdResGtp.MessageTypeName(),
					Cause:   cause,
					Msg:     fmt.Sprintf("Delete Session Response not accepted"),
				}
			}
		} else {
			return &gtpv2.RequiredIEMissingError{
				Type: ie.Cause,
			}
		}

		// TODO: validate message before passing
		enrichedMsg := enriched_message.NewMessageWithGrpc(msg, cdRes)

		// pass message to same session
		if err := gtpv2.PassMessageTo(session, enrichedMsg, 5*time.Second); err != nil {
			return err
		}
		return nil
	}
}

// TODO
func getHandle_DeleteBearerRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		return nil
	}
}

// getHandle_EchoResponse handles echo request received in S8_proxy. This is a special handler
// hat does not use gtpv2.PassMessageTo. It instead uses S8proxy echoChannel to pass the error if any
func getHandle_EchoResponse(echoCh chan error) gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		if _, ok := msg.(*message.EchoResponse); !ok {
			err := &gtpv2.UnexpectedTypeError{Msg: msg}
			echoCh <- err
			return err
		}
		echoCh <- nil
		return nil
	}
}

func handleFTEID(fteidIE *ie.IE) (*protos.Fteid, uint8, error) {
	interfaceType, err := fteidIE.InterfaceType()
	if err != nil {
		return nil, interfaceType, err
	}
	teid, err := fteidIE.TEID()
	if err != nil {
		return nil, interfaceType, err
	}

	fteid := &protos.Fteid{Teid: teid}

	if !fteidIE.HasIPv4() && !fteidIE.HasIPv6() {
		return nil, interfaceType, fmt.Errorf("Error: fteid %+v has no ips", fteidIE.String())
	}
	if fteidIE.HasIPv4() {
		ipv4, err := fteidIE.IPv4()
		if err != nil {
			return nil, interfaceType, err
		}
		fteid.Ipv4Address = ipv4.String()
	}
	if fteidIE.HasIPv6() {
		ipv6, err := fteidIE.IPv6()
		if err != nil {
			return nil, interfaceType, err
		}
		fteid.Ipv6Address = ipv6.String()
	}
	return fteid, interfaceType, nil
}
