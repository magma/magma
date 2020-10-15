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

// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package servicers

import (
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
)

func (s *s6aProxy) addDiamOriginAVPs(m *diam.Message) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Host))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Realm))
	if s.originStateID != 0 {
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(s.originStateID))
	}
}
