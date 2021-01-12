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
package radius

import (
	"magma/feg/cloud/go/protos/mconfig"
)

const (
	defaultNetwork  = "udp"
	defaultAuthAddr = ":1812"
	defaultAcctAddr = ":1813"
	defaultDAEAddr  = ":3799"
	defaultSecret   = "123456"

	defaultAcctInterimUpdateInterval = 180 // in seconds
)

var defaultConfigs = &mconfig.RadiusConfig{
	Secret:   []byte(defaultSecret),
	Network:  defaultNetwork,
	AuthAddr: defaultAuthAddr,
	AcctAddr: defaultAcctAddr,
	DAEAddr:  defaultDAEAddr,
}

func ValidateConfigs(cfg *mconfig.RadiusConfig) *mconfig.RadiusConfig {
	res := &mconfig.RadiusConfig{}
	if cfg != nil {
		*res = *cfg
	}
	if len(res.Secret) == 0 {
		res.Secret = []byte(defaultSecret)
	}
	if len(res.Network) == 0 {
		res.Network = defaultNetwork
	}
	if len(res.AuthAddr) == 0 {
		res.AuthAddr = defaultAuthAddr
	}
	if len(res.AcctAddr) == 0 {
		res.AcctAddr = defaultAcctAddr
	}
	if len(res.DAEAddr) == 0 {
		res.DAEAddr = defaultDAEAddr
	}
	return res
}
