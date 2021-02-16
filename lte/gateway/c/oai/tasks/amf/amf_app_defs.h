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
/*****************************************************************************

  Source      amf_app_defs.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#pragma once
#include <sstream>
#include "amf_as_message.h"
#include "ngap_messages_types.h"
#include "n11_messages_types.h"  //pdu_change
#include "amf_app_ue_context_and_proc.h"

using namespace std;
namespace magma5g {
class amf_app_desc_t {
 public:
  /* UE contexts */
  amf_ue_context_t amf_ue_contexts;
  amf_ue_ngap_id_t amf_app_ue_ngap_id_generator;
  long m5_statistic_timer_id;
  uint32_t m5_statistic_timer_period;
  uint32_t nb_gnb_connected;
  uint32_t nb_ue_registered;
  uint32_t nb_ue_connected;
};

class amf_app_defs : public amf_app_ue_context, public amf_app_desc_t {
 public:
  uint64_t amf_app_handle_initial_ue_message(
      amf_app_desc_t* amf_app_desc_p,
      itti_ngap_initial_ue_message_t* conn_est_ind_pP);
  int amf_app_handle_nas_dl_req(
      amf_ue_ngap_id_t ue_id, bstring nas_msg,
      nas5g_error_code_t transaction_status);
  int amf_app_handle_uplink_nas_message(
      amf_app_desc_t* amf_app_desc_p, bstring msg);
  void amf_app_handle_pdu_session_response(
      itti_n11_create_pdu_session_response_t* pdu_session_resp);
};
}  // namespace magma5g
