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

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func addS8GtpHandlers(s8p *S8Proxy) {
	s8p.gtpClient.AddHandlers(
		map[uint8]gtpv2.HandlerFunc{
			message.MsgTypeCreateSessionResponse: s8p.createSessionResponseHander(),
			message.MsgTypeDeleteSessionResponse: s8p.deleteSessionResponseHandler(),
			message.MsgTypeEchoResponse:          s8p.echoResponseHandler(),
		})
}

func (s *S8Proxy) createSessionResponseHander() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) (err error) {
		csRes, err := parseCreateSessionResponse(msg)
		return s.gtpClient.PassMessage(msg.TEID(), senderAddr, msg, csRes, err)
	}
}

func (s *S8Proxy) deleteSessionResponseHandler() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		dsRes, err := parseDeleteSessionResponse(msg)
		return s.gtpClient.PassMessage(msg.TEID(), senderAddr, msg, dsRes, err)
	}
}

// echoResponseHandler handles echo request received in S8_proxy. This is a special handler
// that does not use gtpv2.PassMessageTo. It instead uses S8proxy echoChannel to pass the error if any
func (s *S8Proxy) echoResponseHandler() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		if _, ok := msg.(*message.EchoResponse); !ok {
			err := &gtpv2.UnexpectedTypeError{Msg: msg}
			s.echoChannel <- err
			return err
		}
		s.echoChannel <- nil
		return nil
	}
}
