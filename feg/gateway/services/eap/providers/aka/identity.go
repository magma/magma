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

// package aka implements EAP-AKA EAP Method
package aka

import (
	"magma/feg/gateway/services/eap"
)

func NewIdentityReq(identifier uint8, attr eap.AttrType) eap.Packet {
	return []byte{
		eap.RequestCode,
		identifier,
		0, 12, // EAP Len
		TYPE,
		byte(SubtypeIdentity),
		0, 0,
		byte(attr),
		1,
		0, 0} // padding
}
