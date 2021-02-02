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

func (mPgw *MockPgw) getHandleDeleteSessionRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		session, err := c.GetSessionByTEID(msg.TEID(), sgwAddr)
		if err != nil {
			dsr := message.NewDeleteSessionResponse(
				0, 0,
				ie.NewCause(gtpv2.CauseIMSIIMEINotKnown, 0, 0, 0, nil),
			)
			if err := c.RespondTo(sgwAddr, msg, dsr); err != nil {
				return err
			}

			return err
		}

		// respond to S-GW with DeleteSessionResponse.
		teid, err := session.GetTEID(gtpv2.IFTypeS5S8SGWGTPC)
		if err != nil {
			err = errors.Wrap(err, "Error")
			return err
		}
		dsr := message.NewDeleteSessionResponse(
			teid, 0,
			ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
		)
		if err := c.RespondTo(sgwAddr, msg, dsr); err != nil {
			return err
		}

		fmt.Printf("Session deleted for Subscriber: %s", session.IMSI)
		c.RemoveSession(session)
		return nil
	}
}
