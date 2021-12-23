/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package n7

import (
	"strings"

	n7_sbi "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControl"
)

// NewN7Client creates a N7 oapi client and sets the OAuth2 client credentiatls for authorizing requests
func NewN7Client(cfg *N7Config) (*n7_sbi.ClientWithResponses, error) {
	n7Options := n7_sbi.WithHTTPClient(cfg.ServerConfig.BuildHttpClient())
	serverString := cfg.ServerConfig.BuildServerString()
	return n7_sbi.NewClientWithResponses(serverString, n7Options)
}

func removeIMSIPrefix(imsi string) string {
	return strings.TrimPrefix(imsi, "IMSI")
}
