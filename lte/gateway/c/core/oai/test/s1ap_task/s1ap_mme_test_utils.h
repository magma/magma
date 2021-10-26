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
#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "s1ap_mme.h"
#include "s1ap_mme_decoder.h"
#include "s1ap_mme_handlers.h"
#include "s1ap_mme_nas_procedures.h"
}

namespace magma {
namespace lte {

status_code_e setup_new_association(
    s1ap_state_t* state, sctp_assoc_id_t assoc_id);
status_code_e generate_s1_setup_request_pdu(S1ap_S1AP_PDU_t* pdu_s1);

void handle_mme_ue_id_notification(s1ap_state_t* s, sctp_assoc_id_t assoc_id);

}  // namespace lte
}  // namespace magma