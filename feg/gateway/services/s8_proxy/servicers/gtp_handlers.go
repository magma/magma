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
	"net"

	"magma/feg/gateway/gtp/enriched_message"

	proto "github.com/golang/protobuf/proto"
	"github.com/wmnsk/go-gtp/gtpv2"
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
		session, err := c.GetSessionByTEID(msg.TEID(), senderAddr)
		if err != nil {
			err = fmt.Errorf("couldn't find session with TEID %d: %s", msg.TEID(), err)
			return err
		}
		csRes, err := parseCreateSessionResponse(session, msg)
		return passMessage(session, msg, csRes, err)
	}
}

func getHandle_DeleteSessionResponse() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		session, err := c.GetSessionByTEID(msg.TEID(), senderAddr)
		if err != nil {
			err = fmt.Errorf("couldn't find session with TEID %d: %s", msg.TEID(), err)
			return err
		}
		csRes, err := parseDelteSessionResponse(session, msg)
		return passMessage(session, msg, csRes, err)
	}
}

// TODO
func getHandle_ModifyBearerRequest() gtpv2.HandlerFunc {
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

// getHandle_EchoResponse handles echo request received in S8_proxy. This is a special handler
// that does not use gtpv2.PassMessageTo. It instead uses S8proxy echoChannel to pass the error if any
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

// passMessage will send a valid message to the caller.
func passMessage(session *gtpv2.Session, gtpMessage message.Message, grpcMessage proto.Message, err error) error {
	enrichedMsg := enriched_message.NewMessageWithGrpc(gtpMessage, grpcMessage, err)
	// pass message to same session
	if err := gtpv2.PassMessageTo(session, enrichedMsg, GtpTimeout); err != nil {
		return err
	}
	return nil
}
