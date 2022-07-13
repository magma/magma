/*
 Copyright 2022 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package storage

import (
	"fmt"

	"magma/orc8r/cloud/go/services/state"
)

func CreateTestGatewaySubscriberState(imsis ...string) GatewaySubscriberState {
	states := GatewaySubscriberState{Subscribers: map[string]state.ArbitraryJSON{}}
	for i, imsi := range imsis {
		states.Subscribers[imsi] = map[string]interface{}{
			"magma.ipv4": []interface{}{
				map[string]interface{}{
					"active_policy_rules": []interface{}{},
					"active_duration_sec": float64(i),
					"lifecycle_state":     "SESSION_RELEASED",
					"session_start_time":  1653484144.,
					"apn":                 "magma.ipv4",
					"ipv4":                "192.168.128.12",
					"msisdn":              "",
					"session_id":          fmt.Sprintf("%s-1234", imsi),
				},
			},
		}
	}
	return states
}
