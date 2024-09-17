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

package diameter

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// NewProxiableRequest creates a request with the proxy bit set, meaning that
// a diameter server can relay the request to another server. This is required
// by certain servers
func NewProxiableRequest(cmd uint32, appid uint32, dictionary *dict.Parser) *diam.Message {
	req := diam.NewRequest(cmd, appid, dictionary)
	req.Header.CommandFlags |= diam.ProxiableFlag
	return req
}
