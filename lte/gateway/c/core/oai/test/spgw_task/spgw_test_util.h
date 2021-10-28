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
#include "intertask_interface.h"
#include "sgw_ie_defs.h"
}

namespace magma {
namespace lte {

#define END_OF_TEST_SLEEP_MS 1000
#define STATE_MAX_WAIT_MS 1000

#define DEFAULT_BEARER_INDEX 0
#define DEFAULT_EPS_BEARER_ID 5
#define DEFAULT_UE_IP 0xc0a8800a  // 192.168.128.10
#define DEFAULT_VLAN 0

void fill_create_session_request(
    itti_s11_create_session_request_t* session_request_p,
    const std::string& imsi_str, int bearer_id,
    bearer_context_to_be_created_t sample_bearer_context, plmn_t sample_plmn);

void fill_ip_allocation_response(
    itti_ip_allocation_response_t* ip_alloc_resp_p, SGIStatus_t status,
    teid_t context_teid, ebi_t eps_bearer_id, unsigned long ue_ip, int vlan);

void send_create_session_request(
    const std::string& imsi_str, int beareri_id,
    bearer_context_to_be_created_t sample_bearer_context, plmn_t sample_plmn);
}  // namespace lte
}  // namespace magma
