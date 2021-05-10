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

		fmt.Println("mock PGW received a DeleteSessionRequest")

		var err error
		dsReqFromSGW := msg.(*message.DeleteSessionRequest)

		pgwTeidC := dsReqFromSGW.TEID()
		session, err := c.GetSessionByTEID(pgwTeidC, sgwAddr)
		mPgw.LastTEIDc = pgwTeidC
		if err != nil {
			return fmt.Errorf("PGW can't find session for PGWC teid %d, %s\n ",
				pgwTeidC, err)
		}

		// get TEUD
		sgwTeidC, err := session.GetTEID(gtpv2.IFTypeS5S8SGWGTPC)
		if err != nil {
			err = errors.Wrap(err, "Error")
			return err
		}

		// check bearer is n there
		if dsReqFromSGW.LinkedEBI == nil {
			dsr := message.NewDeleteSessionResponse(
				sgwTeidC, msg.Sequence(),
				ie.NewCause(gtpv2.CauseMandatoryIEMissing, 0, 0, 0, ie.NewEPSBearerID(0)),
			)
			if err := c.RespondTo(sgwAddr, msg, dsr); err != nil {
				return err
			}

			return err
		}

		// check if bearer associated with EBI exists or not.
		_, err = session.LookupBearerByEBI(dsReqFromSGW.LinkedEBI.MustEPSBearerID())
		if err != nil {
			dsr := message.NewDeleteBearerResponse(
				sgwTeidC, msg.Sequence(),
				ie.NewCause(gtpv2.CauseContextNotFound, 0, 0, 0, nil),
			)
			if err := c.RespondTo(sgwAddr, msg, dsr); err != nil {
				return err
			}
			return err
		}

		// respond to S-GW with DeleteSessionResponse.
		dsr := message.NewDeleteSessionResponse(
			sgwTeidC, msg.Sequence(),
			ie.NewCause(gtpv2.CauseRequestAccepted, 0, 0, 0, nil),
		)
		c.RemoveSession(session)
		if err := c.RespondTo(sgwAddr, msg, dsr); err != nil {
			return err
		}
		fmt.Printf("mock PGW deleted a session for: %s\n", session.IMSI)
		return nil
	}
}
