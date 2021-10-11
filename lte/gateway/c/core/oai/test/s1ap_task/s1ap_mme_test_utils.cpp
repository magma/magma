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
#include "s1ap_mme_test_utils.h"

extern "C" {
#include "intertask_interface.h"
}

namespace magma {
namespace lte {

extern task_zmq_ctx_t task_zmq_ctx_main_s1ap;

status_code_e setup_new_association(
    s1ap_state_t* state, sctp_assoc_id_t assoc_id) {
  bstring ran_cp_ipaddr = bfromcstr("\xc0\xa8\x3c\x8d");
  sctp_new_peer_t p     = {
      .instreams     = 1,
      .outstreams    = 2,
      .assoc_id      = assoc_id,
      .ran_cp_ipaddr = ran_cp_ipaddr,
  };
  status_code_e rc = s1ap_handle_new_association(state, &p);
  bdestroy(ran_cp_ipaddr);
  return rc;
}
status_code_e generate_s1_setup_request_pdu(S1ap_S1AP_PDU_t* pdu_s1) {
  uint8_t packet_bytes[] = {
      0x00, 0x11, 0x00, 0x2f, 0x00, 0x00, 0x04, 0x00, 0x3b, 0x00, 0x09,
      0x00, 0x00, 0xf1, 0x10, 0x40, 0x00, 0x00, 0x00, 0xa0, 0x00, 0x3c,
      0x40, 0x0b, 0x80, 0x09, 0x22, 0x52, 0x41, 0x44, 0x49, 0x53, 0x59,
      0x53, 0x22, 0x00, 0x40, 0x00, 0x07, 0x00, 0x00, 0x00, 0x40, 0x00,
      0xf1, 0x10, 0x00, 0x89, 0x40, 0x01, 0x00};

  bstring payload_s1_setup;
  payload_s1_setup = blk2bstr(&packet_bytes, sizeof(packet_bytes));

  status_code_e pdu_rc = s1ap_mme_decode_pdu(pdu_s1, payload_s1_setup);
  bdestroy_wrapper(&payload_s1_setup);
  return pdu_rc;
}

void handle_mme_ue_id_notification(s1ap_state_t* s, sctp_assoc_id_t assoc_id) {
  MessageDef* message_p =
      itti_alloc_new_message(TASK_MME_APP, MME_APP_S1AP_MME_UE_ID_NOTIFICATION);
  itti_mme_app_s1ap_mme_ue_id_notification_t* notification_p =
      &message_p->ittiMsg.mme_app_s1ap_mme_ue_id_notification;
  memset(notification_p, 0, sizeof(itti_mme_app_s1ap_mme_ue_id_notification_t));
  notification_p->enb_ue_s1ap_id = 1;
  notification_p->mme_ue_s1ap_id = 7;
  notification_p->sctp_assoc_id  = assoc_id;
  s1ap_handle_mme_ue_id_notification(s, notification_p);
  free(message_p);
}

}  // namespace lte
}  // namespace magma
