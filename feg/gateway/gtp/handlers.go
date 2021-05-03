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

package gtp

import (
	"net"

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

func (c *Client) AddDefaultHandlers() {
	// Add default
	c.AddHandlers(
		map[uint8]gtpv2.HandlerFunc{
			message.MsgTypeEchoResponse: c.echoResponseHandler(),
		})
}

// echoResponseHandler is a special handler that does not use gtpv2.PassMessageTo.
// It instead uses echoChannel to pass the error if any
func (c *Client) echoResponseHandler() gtpv2.HandlerFunc {
	return func(_ *gtpv2.Conn, senderAddr net.Addr, msg message.Message) error {
		if _, ok := msg.(*message.EchoResponse); !ok {
			err := &gtpv2.UnexpectedTypeError{Msg: msg}
			c.echoChannel <- err
			return err
		}
		c.echoChannel <- nil
		return nil
	}
}
