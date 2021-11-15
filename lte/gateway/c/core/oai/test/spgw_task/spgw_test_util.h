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

#include "lte/gateway/c/core/oai/include/spgw_state.h"

extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.h"
#include "lte/gateway/c/core/oai/include/sgw_ie_defs.h"
}

namespace magma {
namespace lte {

#define END_OF_TEST_SLEEP_MS 1000
#define STATE_MAX_WAIT_MS 1000

#define DEFAULT_MME_S11_TEID 1
#define DEFAULT_BEARER_INDEX 0
#define DEFAULT_EPS_BEARER_ID 5
#define UNASSIGNED_UE_IP 0
#define DEFAULT_UE_IP 0xc0a8800a  // 192.168.128.10
#define DEFAULT_VLAN 0
#define DEFAULT_ENB_GTP_TEID 1
#define ERROR_ENB_GTP_TEID 100
#define ERROR_SGW_S11_TEID 100
#define DEFAULT_EDNS_IP 0x7f000001  // localhost
#define DEFAULT_SGW_IP 0x7f000001   // localhost
#define DEFAULT_ENB_IP 0xc0a83c8d // 192.168.60.141

bool is_num_sessions_valid(
    uint64_t imsi64, int expected_num_ue_contexts,
    int expected_num_teids);

bool is_num_s1_bearers_valid(
    teid_t context_teid, int expected_num_active_bearers);

void fill_create_session_request(
    itti_s11_create_session_request_t* session_request_p,
    const std::string& imsi_str, teid_t mme_s11_teid, int bearer_id,
    bearer_context_to_be_created_t sample_bearer_context, plmn_t sample_plmn);

void fill_ip_allocation_response(
    itti_ip_allocation_response_t* ip_alloc_resp_p, SGIStatus_t status,
    teid_t context_teid, ebi_t eps_bearer_id, unsigned long ue_ip, int vlan);

void fill_pcef_create_session_response(
    itti_pcef_create_session_response_t* pcef_csr_resp_p,
    PcefRpcStatus_t csr_status, teid_t context_teid, ebi_t eps_bearer_id,
    SGIStatus_t sgi_status);

void fill_modify_bearer_request(
    itti_s11_modify_bearer_request_t* modify_bearer_req, teid_t mme_s11_teid,
    teid_t sgw_s11_teid, teid_t enb_gtp_teid, int bearer_id,
    ebi_t eps_bearer_id);

void fill_delete_session_request(
    itti_s11_delete_session_request_t* delete_session_req, teid_t mme_s11_teid,
    teid_t sgw_s11_context_teid, ebi_t eps_bearer_id, plmn_t test_plmn);

void fill_release_access_bearer_request(
    itti_s11_release_access_bearers_request_t* release_access_bearers_req,
    teid_t mme_s11_teid, teid_t sgw_s11_context_teid);
}  // namespace lte
}  // namespace magma
