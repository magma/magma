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

package eap

import (
	"fbc/cwf/radius/monitoring"
)

var (
	// ExtractEapPacket extraction EAP-Message from RADIUS message
	ExtractEapPacket = monitoring.NewOperation("eap_extract_packet_from_radius")

	// RestoreProtocolState restore eap state from storage
	RestoreProtocolState = monitoring.NewOperation("eap_restore_state")

	// HandleEapPacket handling EAP packet
	HandleEapPacket = monitoring.NewOperation("eap_handle")

	// PersistProtocolState writing new state, after handling, to storage
	PersistProtocolState = monitoring.NewOperation("eap_persist_state")
)
