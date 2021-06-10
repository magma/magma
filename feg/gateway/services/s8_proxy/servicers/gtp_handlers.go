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
	"fmt"
	"net"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func addS8GtpHandlers(s8p *S8Proxy) {
	// echo hanlders added by gtp_client. Use echoChannel for errors
	s8p.gtpClient.AddHandlers(
		map[uint8]gtpv2.HandlerFunc{
			message.MsgTypeCreateSessionResponse: s8p.createSessionResponseHandler(),
			message.MsgTypeDeleteSessionResponse: s8p.deleteSessionResponseHandler(),
			message.MsgTypeCreateBearerRequest:   s8p.createBearerRequestHandler(),
		})
}

func (s *S8Proxy) createSessionResponseHandler() gtpv2.HandlerFunc {
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

func (s *S8Proxy) createBearerRequestHandler() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, senderAddr net.Addr, msg message.Message) (err error) {

		cbReq, gtpErr, err := parseCreateBearerRequest(msg)
		if err != nil {
			return err
		}
		if gtpErr != nil {
			return fmt.Errorf(gtpErr.Msg)
		}
		cbRes, err := GWS8ProxyCreateBearerRequest(cbReq)
		if err != nil {
			return fmt.Errorf("Failed while CreateBearerRequest to feg relay: %s", err)
		}

		cbResMsg, _ := buildCreateBearerResMsg(msg.Sequence(), cbRes)
		err = c.RespondTo(senderAddr, msg, cbResMsg)
		if err != nil {
			return err
		}

		return nil
	}
}
