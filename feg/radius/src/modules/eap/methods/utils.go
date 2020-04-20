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

package methods

import (
	"fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
)

// ToRadiusCode returs the RADIUS packet code which, as per RFCxxxx
// should carry the EAP payload of the given EAP Code
func ToRadiusCode(eapCode packet.Code) radius.Code {
	switch eapCode {
	case packet.CodeFAILURE:
		return radius.CodeAccessReject
	case packet.CodeSUCCESS:
		return radius.CodeAccessAccept
	case packet.CodeRESPONSE:
	case packet.CodeREQUEST:
		return radius.CodeAccessChallenge
	}
	return radius.CodeAccessReject
}
