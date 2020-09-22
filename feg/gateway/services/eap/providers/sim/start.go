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

// package sim implements EAP-SIM EAP Method
package sim

import (
	"magma/feg/gateway/services/eap"
)

func NewStartReq(identifier uint8, identityTypeAttr eap.AttrType) eap.Packet {
	return []byte{
		eap.RequestCode,
		identifier,
		0, 20, // EAP Len
		TYPE,
		byte(SubtypeStart),
		0, 0,
		byte(identityTypeAttr),
		1,    // attr len / 4
		0, 0, // padding
		byte(AT_VERSION_LIST),
		2,    // attr len / 4
		0, 2, // Actual Version List Length
		0, Version, // Our Supported Version
		0, 0, // padding
	}
}
