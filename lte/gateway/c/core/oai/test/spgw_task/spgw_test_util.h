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
#include <string>

extern "C" {
#include "sgw_ie_defs.h"
}

namespace magma {
namespace lte {

#define END_OF_TEST_SLEEP_MS 1000
#define STATE_MAX_WAIT_MS 1000

void send_create_session_request(
    const std::string& imsi_str, int beareri_id,
    bearer_context_to_be_created_t sample_bearer_context, plmn_t sample_plmn);
}  // namespace lte
}  // namespace magma
