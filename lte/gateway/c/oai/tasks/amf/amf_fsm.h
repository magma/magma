/**
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

#pragma once

namespace magma5g {
// TODO in upcoming PR of FSM this enum and got changed
enum amf_fsm_state_t {
  AMF_STATE_MIN = 0,
  AMF_DEREGISTERED,
  AMF_REGISTERED,
  AMF_DEREGISTERED_INITIATED,
  AMF_COMMON_PROCEDURE_INITIATED,
  AMF_STATE_MAX
};
}  // namespace magma5g
