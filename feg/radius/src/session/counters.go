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

package session

import (
	"fbc/cwf/radius/monitoring"
)

var (
	// ReadSessionState counts reading session state from storage
	ReadSessionState = monitoring.NewOperation("read_session_state")

	// WriteSessionState counts writing session state from storage
	WriteSessionState = monitoring.NewOperation("write_session_state")

	// ResetSessionState counts reseting session state from storage
	ResetSessionState = monitoring.NewOperation("reset_session_state")
)
