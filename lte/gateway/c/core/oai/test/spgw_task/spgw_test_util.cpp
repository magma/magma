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

#include "spgw_test_util.h"

extern "C" {
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
}

namespace magma {
namespace lte {

extern task_zmq_ctx_t task_zmq_ctx_main_spgw;

void send_create_session_request(
    const std::string& imsi_str, int bearer_id,
    bearer_context_to_be_created_t sample_bearer_context, plmn_t sample_plmn) {
  MessageDef* message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_MME_APP, S11_CREATE_SESSION_REQUEST);
  itti_s11_create_session_request_t* session_request_p =
      &message_p->ittiMsg.s11_create_session_request;

  session_request_p->teid = 0;
  strncpy(
      (char*) session_request_p->imsi.digit, imsi_str.c_str(), imsi_str.size());
  session_request_p->imsi.length = imsi_str.size();
  session_request_p->sender_fteid_for_cp.teid == 1;
  session_request_p->sender_fteid_for_cp.interface_type = S11_MME_GTP_C;

  // Fill User Location Information
  session_request_p->uli.present = 0;  // initialize the presencemask
  session_request_p->rat_type    = RAT_EUTRAN;

  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .eps_bearer_id = sample_bearer_context.eps_bearer_id;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.pci = sample_bearer_context.bearer_level_qos.pci;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.pl = sample_bearer_context.bearer_level_qos.pl;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.pvi = sample_bearer_context.bearer_level_qos.pvi;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.qci = sample_bearer_context.bearer_level_qos.qci;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.mbr.br_ul =
      sample_bearer_context.bearer_level_qos.mbr.br_ul;
  session_request_p->bearer_contexts_to_be_created.bearer_contexts[bearer_id]
      .bearer_level_qos.mbr.br_dl =
      sample_bearer_context.bearer_level_qos.mbr.br_dl;
  session_request_p->bearer_contexts_to_be_created.num_bearer_context = 1;

  session_request_p->sender_fteid_for_cp.teid           = (teid_t) 1;
  session_request_p->sender_fteid_for_cp.interface_type = S11_MME_GTP_C;
  session_request_p->sender_fteid_for_cp.ipv4_address.s_addr =
      0xc0a83c8e;  // 192.168.60.142
  session_request_p->sender_fteid_for_cp.ipv4 = 1;

  const char default_apn[] = "magma.ipv4";
  strncpy(session_request_p->apn, default_apn, 10);
  session_request_p->ambr.br_dl = 100000000;
  session_request_p->ambr.br_ul = 200000000;

  session_request_p->pdn_type                = IPv4;
  session_request_p->paa.pdn_type            = IPv4;
  session_request_p->paa.ipv4_address.s_addr = INADDR_ANY;
  session_request_p->paa.ipv6_address        = in6addr_any;

  session_request_p->serving_network.mcc[0] = sample_plmn.mcc_digit1;
  session_request_p->serving_network.mcc[1] = sample_plmn.mcc_digit2;
  session_request_p->serving_network.mcc[2] = sample_plmn.mcc_digit3;
  session_request_p->serving_network.mnc[0] = sample_plmn.mnc_digit1;
  session_request_p->serving_network.mnc[1] = sample_plmn.mnc_digit2;
  session_request_p->serving_network.mnc[2] = sample_plmn.mnc_digit3;

  send_msg_to_task(&task_zmq_ctx_main_spgw, TASK_SPGW_APP, message_p);
  // itti_free_msg_content(message_p);
  // free(message_p);
}
}  // namespace lte
}  // namespace magma
