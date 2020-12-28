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

	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

// TODO: to be implemented
func getHandleDeleteSessionRequest() gtpv2.HandlerFunc {
	return func(c *gtpv2.Conn, sgwAddr net.Addr, msg message.Message) error {
		fmt.Print("mock PGW received DeleteSessionRequest")
		return nil
	}
}
