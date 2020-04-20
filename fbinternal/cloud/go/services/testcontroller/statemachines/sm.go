/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package statemachines

import "time"

// TestMachine is the interface for a test case state machine.
type TestMachine interface {
	// Run will perform the appropriate action for the test case in the
	// given state. This will return the state to transition into, the time
	// delay for the next scheduled run, and an error if one occurred.
	Run(state string, config interface{}, previousErr error) (string, time.Duration, error)
}
