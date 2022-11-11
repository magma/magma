
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
#include <chrono>
#include <gtest/gtest.h>
#include <cstdint>
#include <thread>
#include <mutex>
#include <condition_variable>
#include <stdio.h>

#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_app_test_util.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_procedure_test_fixture.hpp"

namespace magma {
namespace lte {

TEST_F(MmeAppProcedureTest, TestNwInitiatedDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  uint8_t ebi_to_be_deactivated = 5;
  // Constructing and sending deactivate bearer request
  // for default bearer that should trigger session termination
  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, true);

  // Wait for DL NAS Transport for Detach Request
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending Detach Accept to mme_app
  // mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_detach_accept,
                               sizeof(nas_msg_detach_accept), plmn);

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

TEST_F(MmeAppProcedureTest, TestNwInitiatedExpiredDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3422 duration for testing
  mme_config.nas_config.t3422_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(8, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  uint8_t ebi_to_be_deactivated = 5;
  // Constructing and sending deactivate bearer request
  // for default bearer that should trigger session termination
  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, true);

  // Wait for DL NAS Transport for Detach Request up to retransmission limit.
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Constructing and sending Detach Accept to mme_app
  // mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_detach_accept,
                               sizeof(nas_msg_detach_accept), plmn);

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

TEST_F(MmeAppProcedureTest, TestNwInitiatedDetachRetxFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3422 duration for testing
  mme_config.nas_config.t3422_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(8, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  uint8_t ebi_to_be_deactivated = 5;
  // Constructing and sending deactivate bearer request
  // for default bearer that should trigger session termination
  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, true);

  // Wait for DL NAS Transport for Detach Request up to retransmission limit.
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // We are not sending Detach Accept here, so timer T3422 will expire
  // once more. This should trigger implicit detach.

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

TEST_F(MmeAppProcedureTest, TestNwInitiatedActivateDedicatedBearerRej) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 1, 0, 1, 3);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send activate dedicated bearer request for lbi 5 mimicing SPGW
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(1);
  send_s11_create_bearer_req(5);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(6);

  // Constructing and sending Activate Dedicated Bearer Reject to mme_app
  // mimicing S1AP
  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_actv_bearer_rsp())
      .Times(1);
  send_mme_app_uplink_data_ind(nas_msg_activate_ded_bearer_reject,
                               sizeof(nas_msg_activate_ded_bearer_reject),
                               plmn);

  // Check MME state after Bearer Activation Reject
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest,
       TestNwInitiatedDedicatedBearerActivationFailureWithT3485Expiry) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3485 durations for testing
  mme_config.nas_config.t3485_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(7, 1, 1, 1, 1, 1, 1, 1, 0, 1, 3);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send activate dedicated bearer request for lbi 5 mimicing SPGW
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(1);
  send_s11_create_bearer_req(5);

  // Send ERAB Setup Response mimicing S1AP
  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_actv_bearer_rsp())
      .Times(1);
  send_erab_setup_rsp(6);

  // Wait for timer expiry.
  for (int i = 1; i < NAS_RETX_LIMIT + 1; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Check MME state after Bearer Activation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest,
       TestNwInitiatedBearerDeactivationFailureWithT3495Expiry) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3495 duration for testing
  mme_config.nas_config.t3495_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 1, 0, 1, 4);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send activate dedicated bearer request for lbi 5 mimicing SPGW
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(1);
  send_s11_create_bearer_req(5);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(6);

  // Constructing and sending Activate Dedicated Bearer Accept to mme_app
  // mimicing S1AP
  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_actv_bearer_rsp())
      .Times(1);
  send_mme_app_uplink_data_ind(nas_msg_activate_ded_bearer_accept,
                               sizeof(nas_msg_activate_ded_bearer_accept),
                               plmn);

  // Check MME state after Bearer Activation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  uint8_t ebi_to_be_deactivated = 6;
  // Constructing and sending deactivate bearer request
  // for dedicated bearer that should trigger ERAB Release Command
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_rel_cmd())
      .Times(5)
      .WillRepeatedly(ReturnFromAsyncTask(&cv));
  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, false);

  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_deactv_bearer_rsp())
      .Times(1);
  // Wait for timer expiry.
  for (int i = 0; i < NAS_RETX_LIMIT + 1; ++i) {
    // Send ERAB Release Response mimicing S1AP
    send_erab_release_rsp();
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Check MME state after Bearer Deactivation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest,
       TestNwInitiatedActivateDeactivateDedicatedBearerRequest) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 1, 0, 1, 4);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send activate dedicated bearer request for lbi 5 mimicing SPGW
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(1);
  send_s11_create_bearer_req(5);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(6);

  // Constructing and sending Activate Dedicated Bearer Accept to mme_app
  // mimicing S1AP
  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_actv_bearer_rsp())
      .Times(1);
  send_mme_app_uplink_data_ind(nas_msg_activate_ded_bearer_accept,
                               sizeof(nas_msg_activate_ded_bearer_accept),
                               plmn);

  // Check MME state after Bearer Activation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  uint8_t ebi_to_be_deactivated = 6;
  // Constructing and sending deactivate bearer request
  // for dedicated bearer that should trigger ERAB Release Command
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_rel_cmd()).Times(1);
  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, false);

  // Send ERAB Release Response mimicing S1AP
  send_erab_release_rsp();

  // Constructing and sending Deactivate Dedicated Bearer Accept to mme_app
  // mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_deactivate_ded_bearer_accept,
                               sizeof(nas_msg_deactivate_ded_bearer_accept),
                               plmn);

  // Check MME state after Bearer Deactivation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(
    MmeAppProcedureTest,
    TestNwInitiatedActivateDeactivateDedicatedBearerWithT3485T3495Expirations) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3485 and 3495 durations for testing
  mme_config.nas_config.t3485_msec = MME_APP_TIMER_TO_MSEC;
  mme_config.nas_config.t3495_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(7, 1, 1, 1, 1, 1, 1, 1, 0, 1, 4);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send activate dedicated bearer request for lbi 5 mimicing SPGW
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(1);
  send_s11_create_bearer_req(5);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(6);

  // Wait for DL NAS Transport up to retransmission limit.
  // The first transmission was piggybacked on ERAB Setup Request
  for (int i = 1; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Constructing and sending Activate Dedicated Bearer Accept to mme_app
  // mimicing S1AP
  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_actv_bearer_rsp())
      .Times(1);
  send_mme_app_uplink_data_ind(nas_msg_activate_ded_bearer_accept,
                               sizeof(nas_msg_activate_ded_bearer_accept),
                               plmn);

  // Check MME state after Bearer Activation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  uint8_t ebi_to_be_deactivated = 6;
  // Constructing and sending deactivate bearer request
  // for dedicated bearer that should trigger ERAB Release Command
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_rel_cmd())
      .Times(5)
      .WillRepeatedly(ReturnFromAsyncTask(&cv));

  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, false);

  // Wait for ERAB Release up to retransmission limit;
  // Deactivate EPS bearer context requests are piggybacked on ERAB Release
  // Command.
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    // Send ERAB Release Response mimicing S1AP
    send_erab_release_rsp();
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Constructing and sending Deactivate Dedicated Bearer Accept to mme_app
  // mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_deactivate_ded_bearer_accept,
                               sizeof(nas_msg_deactivate_ded_bearer_accept),
                               plmn);

  // Check MME state after Bearer Deactivation
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

}  // namespace lte
}  // namespace magma
