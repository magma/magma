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
extern "C" {
#include "lte/gateway/c/core/oai/include/s1ap_types.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_decoder.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_handlers.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_mme_nas_procedures.h"
}

namespace magma {
namespace lte {

status_code_e setup_new_association(s1ap_state_t* state,
                                    sctp_assoc_id_t assoc_id);

status_code_e generate_s1_setup_request_pdu(S1ap_S1AP_PDU_t* pdu_s1);

status_code_e send_s1ap_erab_rel_cmd(mme_ue_s1ap_id_t ue_id,
                                     enb_ue_s1ap_id_t enb_ue_id);

status_code_e send_conn_establishment_cnf(mme_ue_s1ap_id_t ue_id,
                                          bool sec_capabilities_present,
                                          bool ue_radio_capability);

status_code_e send_s1ap_erab_setup_req(mme_ue_s1ap_id_t ue_id,
                                       enb_ue_s1ap_id_t enb_ue_id, ebi_t ebi);

status_code_e send_s1ap_erab_reset_req(sctp_assoc_id_t assoc_id,
                                       sctp_stream_id_t stream_id,
                                       enb_ue_s1ap_id_t enb_ue_id,
                                       mme_ue_s1ap_id_t ue_id);

bool is_enb_state_valid(s1ap_state_t* state, sctp_assoc_id_t assoc_id,
                        mme_s1_enb_state_s expected_state,
                        uint32_t expected_num_ues);

bool is_num_enbs_valid(s1ap_state_t* state, uint32_t expected_num_enbs);

bool is_ue_state_valid(sctp_assoc_id_t assoc_id, enb_ue_s1ap_id_t enb_ue_id,
                       enum s1_ue_state_s expected_ue_state);

void handle_mme_ue_id_notification(s1ap_state_t* s, sctp_assoc_id_t assoc_id);

}  // namespace lte
}  // namespace magma