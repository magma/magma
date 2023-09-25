/**
 * Copyright 2022 The Magma Authors.
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

#include "lte/gateway/c/core/oai/test/mme_app_task/mme_procedure_test_fixture.hpp"

namespace magma {
namespace lte {

void MmeAppProcedureTest ::attach_ue(std::condition_variable& cv,
                                     std::unique_lock<std::mutex>& lock,
                                     mme_app_desc_t* mme_state_p,
                                     guti_eps_mobile_identity_t* guti) {
  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, *guti, 1);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi, true);

  // Wait for DL NAS Transport for once
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending Authentication Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_auth_resp, sizeof(nas_msg_auth_resp),
                               plmn);

  // Wait for DL NAS Transport for once
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending SMC Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_smc_resp, sizeof(nas_msg_smc_resp),
                               plmn);

  // Sending ULA to mme_app mimicing successful S6A response for ULR
  send_s6a_ula(imsi, true);

  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  send_create_session_resp(REQUEST_ACCEPTED, DEFAULT_LBI);

  // Wait for ICS request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  nas_message_t nas_msg_decoded = {0};
  emm_security_context_t emm_security_context;
  nas_message_decode_status_t decode_status;
  int decoder_rc = 0;
  decoder_rc = nas_message_decode(
      nas_msg->data, &nas_msg_decoded, nas_msg->slen,
      reinterpret_cast<void*>(&emm_security_context), &decode_status);
  EXPECT_EQ(nas_msg->slen, 67);
  EXPECT_EQ(decoder_rc, nas_msg->slen);
  *guti = nas_msg_decoded.plain.emm.attach_accept.guti.guti;
  bdestroy_wrapper(
      &nas_msg_decoded.plain.emm.attach_accept.esmmessagecontainer);

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

  // Constructing UE Capability Indication message to mme_app
  // mimicing S1AP with dummy radio capabilities
  send_ue_capabilities_ind();

  // Constructing and sending Attach Complete to mme_app
  // mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_attach_comp, sizeof(nas_msg_attach_comp),
                               plmn);

  // Wait for DL NAS Transport for EMM Information
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Destruction at tear down is not sufficient as nas_msg might be used
  // again in the TC
  bdestroy_wrapper(&nas_msg);
  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW
  std::vector<int> b_modify = {5};
  std::vector<int> b_rm = {};
  send_modify_bearer_resp(b_modify, b_rm);

  // Check MME state after Modify Bearer Response
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

void MmeAppProcedureTest ::detach_ue(std::condition_variable& cv,
                                     std::unique_lock<std::mutex>& lock,
                                     mme_app_desc_t* mme_state_p,
                                     guti_eps_mobile_identity_t guti,
                                     bool is_initial_ue) {
  // Constructing and sending Detach Request to mme_app
  // mimicing S1AP
  if (is_initial_ue) {
    send_mme_app_initial_ue_msg(nas_msg_detach_req, sizeof(nas_msg_detach_req),
                                plmn, guti, 1);
  } else {
    send_mme_app_uplink_data_ind(nas_msg_detach_req, sizeof(nas_msg_detach_req),
                                 plmn);
  }
  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  send_delete_session_resp(DEFAULT_LBI);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after detach complete
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}
}  // namespace lte
}  // namespace magma
