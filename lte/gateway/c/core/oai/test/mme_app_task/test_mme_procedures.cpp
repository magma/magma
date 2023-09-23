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
#include <chrono>
#include <gtest/gtest.h>
#include <cstdint>
#include <thread>
#include <mutex>
#include <condition_variable>
#include <stdio.h>

#include "feg/protos/s6a_proxy.pb.h"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_ip_imsi.hpp"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/proto_msg_to_itti_msg.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_app_test_util.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_procedure_test_fixture.hpp"

#include "lte/gateway/c/core/oai/include/mme_app_state.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"

namespace magma {
namespace lte {

TEST_F(MmeAppProcedureTest, TestInitialUeMessageFaultyNasMsg) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  // The following buffer just includes an attach request
  uint8_t nas_msg_faulty[29] = {0x72, 0x08, 0x09, 0x10, 0x10, 0x00, 0x00, 0x00,
                                0x00, 0x10, 0x02, 0xe0, 0xe0, 0x00, 0x04, 0x02,
                                0x01, 0xd0, 0x11, 0x40, 0x08, 0x04, 0x02, 0x60,
                                0x04, 0x00, 0x02, 0x1c, 0x00};
  send_mme_app_initial_ue_msg(nas_msg_faulty, sizeof(nas_msg_faulty), plmn,
                              guti, 1);

  // Wait for DL NAS Transport for once
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Check MME state: at this point, MME should be sending
  // EMM_STATUS NAS message and holding onto the UE context
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest, TestGutiAttachEpsOnlyDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_guti_attach_req,
                              sizeof(nas_msg_guti_attach_req), plmn, guti, 1);

  // Wait for DL NAS Transport for once
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending Identity Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_ident_resp, sizeof(nas_msg_ident_resp),
                               plmn);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

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

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyAirFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Sending AIA to mme_app mimicing negative S6A response for AIR
  send_authentication_info_resp(imsi, false);

  // Wait for context release command; MME should be sending attach reject
  // as well as context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyAirTimeout) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce the S6a timer duration for testing
  mme_config.nas_config.ts6a_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Wait for context release command; MME should be sending attach reject
  // as well as context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyAuthMacFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 0, 2, 1, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi, true);

  // Wait for DL NAS Transport for Auth Req
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending Authentication Response with Mac Failure to
  // mme_app mimicing S1AP
  // Message Type 0x5c = Auth Failure
  // Cause 0x14 = MAC failure
  const uint8_t nas_msg_auth_resp_mac_fail[19] = {
      0x07, 0x5c, 0x14, 0x30, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};

  send_mme_app_uplink_data_ind(nas_msg_auth_resp_mac_fail,
                               sizeof(nas_msg_auth_resp_mac_fail), plmn);

  // Wait for DL NAS Transport for Auth Reject
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for DL NAS Transport for Attach Reject
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after UE context release
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyAuthSynchFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 0, 2, 2, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi, true);

  // Wait for DL NAS Transport for Auth Req
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending Authentication Response with Synchronization
  // Failure to mme_app mimicing S1AP
  // Message Type 0x5c = Auth Failure
  // Cause 0x15 = Synchronization failure
  const uint8_t nas_msg_auth_resp_synch_fail[19] = {
      0x07, 0x5c, 0x15, 0x30, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00};

  send_mme_app_uplink_data_ind(nas_msg_auth_resp_synch_fail,
                               sizeof(nas_msg_auth_resp_synch_fail), plmn);

  // Sending AIA to mme_app mimicing failed S6A response for AIR
  send_authentication_info_resp(imsi, false);

  // Wait for DL NAS Transport for Auth Reject
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for DL NAS Transport for Attach Reject
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after UE context release
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachEpsOnlyUlaFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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

  // Sending ULA to mme_app mimicing negative S6A response for ULR
  send_s6a_ula(imsi, false);

  // Wait for context release command; MME should be sending attach reject
  // as well as context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);

  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachRejectAuthRetxFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3460 duration for testing
  mme_config.nas_config.t3460_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(6, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi, true);

  // Wait for DL NAS Transport to max out retransmission limit
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Wait for context release command; MME should be sending attach reject
  // as well as context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachRejectSmcRetxFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3460 duration for testing
  mme_config.nas_config.t3460_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(6, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi, true);

  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending Authentication Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_auth_resp, sizeof(nas_msg_auth_resp),
                               plmn);

  // Wait for DL NAS Transport to max out retransmission limit
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Wait for context release command; MME should be performing
  // implicit detach
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestGutiAttachExpiredIdentity) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3470 duration for testing
  mme_config.nas_config.t3470_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(8, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_guti_attach_req,
                              sizeof(nas_msg_guti_attach_req), plmn, guti, 1);

  // Wait for DL NAS Transport up to retransmission limit
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }
  // Constructing and sending Identity Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_ident_resp, sizeof(nas_msg_ident_resp),
                               plmn);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

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

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachRejectIdentRetxFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3470 duration for testing
  mme_config.nas_config.t3470_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(5, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1);

  // Constructing and sending Initial Attach Request with GUTI to
  // mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_guti_attach_req,
                              sizeof(nas_msg_guti_attach_req), plmn, guti, 1);

  // Wait for DL NAS Transport to max out retransmission limit
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }

  // Wait for context release command; MME should be sending attach reject
  // as well as context release command.
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestIcsRequestTimeout) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);

  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce the ICS timeout for testing
  mme_config.nas_config.tics_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(2, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Wait for ICS Request timeout
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  send_delete_session_resp(DEFAULT_LBI);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after delete session processing
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestImsiAttachIcsFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(2, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Send ICS failure to mme_app mimicing S1AP
  send_ics_failure();

  // Wait for Delete Session request
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  send_delete_session_resp(DEFAULT_LBI);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after delete session processing
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestCreateSessionFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Context release request is triggered twice once during Attach Reject
  // and once during processing the response for Delete Session Request
  MME_APP_EXPECT_CALLS(3, 0, 2, 1, 1, 1, 1, 0, 0, 1, 1);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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
  send_create_session_resp(M_PDN_APN_NOT_ALLOWED, DEFAULT_LBI);

  // Wait for context release command; MME should be sending attach reject
  // as well as context release command.
  // This should be unnecessary but a delete session request is also
  // triggered in the current code and need to wait for that event at
  // spgw handler.
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  send_ue_ctx_release_complete();

  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  send_delete_session_resp(DEFAULT_LBI);

  // Waiting for the receptiopn of the second context release request
  // which is triggered after receiving delete session response.
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Check if the context is properly released
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
}

TEST_F(MmeAppProcedureTest, TestAttachIdleDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 2, 1, 1, 1, 1, 1, 1, 1, 3);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest, TestAttachIdleServiceReqDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 2, 2, 1, 1, 1, 1, 2, 1, 1, 4);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending Service Request
  send_mme_app_initial_ue_msg(nas_msg_service_req, sizeof(nas_msg_service_req),
                              plmn, guti, 1);

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest, TestPagingMaxRetx) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce paging retransmission timer duration for testing
  mme_config.nas_config.tpaging_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(3, 2, 2, 1, 1, 1, 1, 2, 1, 1, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending Paging Request
  EXPECT_CALL(*s1ap_handler, s1ap_handle_paging_request())
      .Times(2)
      .WillRepeatedly(ReturnFromAsyncTask(&cv));
  send_paging_request();

  // wait for s1ap_handle_paging_request to arrive twice
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending Service Request
  send_mme_app_initial_ue_msg(nas_msg_service_req, sizeof(nas_msg_service_req),
                              plmn, guti, 1);

  // Wait for ICS request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

// Periodic TAU with active flag set to true
TEST_F(MmeAppProcedureTest, TestAttachIdlePeriodicTauReqWithActiveFlag) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 2, 2, 1, 1, 1, 1, 2, 1, 1, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending periodic TAU Request with active flag
  send_mme_app_initial_ue_msg(nas_msg_periodic_tau_req_with_actv_flag,
                              sizeof(nas_msg_periodic_tau_req_with_actv_flag),
                              plmn, guti, 1);

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

// Periodic TAU with active flag set to false
TEST_F(MmeAppProcedureTest, TestAttachIdlePeriodicTauReqWithoutActiveFlag) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 3, 1, 1, 1, 1, 1, 1, 1, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending periodic TAU Request without active flag
  send_mme_app_initial_ue_msg(
      nas_msg_periodic_tau_req_without_actv_flag,
      sizeof(nas_msg_periodic_tau_req_without_actv_flag), plmn, guti, 1);

  // Wait for UE context release cmd and DL NAS Transport
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  detach_ue(cv, lock, mme_state_p, guti, true);
}

// Normal TAU sent in idle mode with active flag set to true
TEST_F(MmeAppProcedureTest, TestAttachIdleNormalTauReqWithActiveFlag) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 2, 2, 1, 1, 1, 1, 2, 1, 1, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending Normal TAU Request with active flag
  send_mme_app_initial_ue_msg(nas_msg_normal_tau_req_with_actv_flag,
                              sizeof(nas_msg_normal_tau_req_with_actv_flag),
                              plmn, guti, 1);

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

// Normal TAU sent in idle mode with active flag set to false
TEST_F(MmeAppProcedureTest, TestAttachIdleNormalTauReqWithoutActiveFlag) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 3, 1, 1, 1, 1, 1, 1, 1, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending TAU Request without active flag
  send_mme_app_initial_ue_msg(nas_msg_normal_tau_req_without_actv_flag,
                              sizeof(nas_msg_normal_tau_req_without_actv_flag),
                              plmn, guti, 1);

  // Wait for UE context release cmd and DL NAS Transport
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  detach_ue(cv, lock, mme_state_p, guti, true);
}

// Normal TAU sent in connected mode with active flag set to false
TEST_F(MmeAppProcedureTest, TestAttachNormalTauReqInConnectedState) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Constructing and sending TAU Request without active flag
  send_mme_app_uplink_data_ind(nas_msg_normal_tau_req_without_actv_flag,
                               sizeof(nas_msg_normal_tau_req_without_actv_flag),
                               plmn);
  // Wait for MME to send TAU accept
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

// TAU reject due to service area restriction
TEST_F(MmeAppProcedureTest, TestTauRejDueToInvalidTac) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 3, 1, 1, 1, 1, 1, 1, 1, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending TAU Request with invalid TAC value 2
  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg_nas_dl_data.mme_ue_s1ap_id);
  EXPECT_FALSE(ue_mm_context == nullptr);
  ue_mm_context->num_reg_sub = 1;
  send_mme_app_initial_ue_msg(nas_msg_normal_tau_req_with_actv_flag,
                              sizeof(nas_msg_normal_tau_req_with_actv_flag),
                              plmn, guti, 2);

  // Waiting for context release request & DL NAS TRANSPORT
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  detach_ue(cv, lock, mme_state_p, guti, true);
}

TEST_F(MmeAppProcedureTest, TestFailedPagingForPendingBearers) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce the paging retransmission timer duration for testing
  mme_config.nas_config.tpaging_msec = MME_APP_TIMER_TO_MSEC;

  MME_APP_EXPECT_CALLS(3, 1, 2, 1, 1, 1, 1, 1, 1, 1, 4);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Force switching to IDLE Mode
  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Trigger paging via bearer request in control plane
  // Send activate dedicated bearer request mimicing SPGW
  send_s11_create_bearer_req(DEFAULT_LBI);
  EXPECT_CALL(*s1ap_handler, s1ap_handle_paging_request())
      .Times(MAX_PAGING_RETRY_COUNT + 1)
      .WillRepeatedly(ReturnFromAsyncTask(&cv));
  // Force paging failure via ignoring paging requests
  for (int i = 0; i <= MAX_PAGING_RETRY_COUNT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }
  EXPECT_CALL(*spgw_handler, sgw_handle_nw_initiated_actv_bearer_rsp())
      .Times(1)
      .WillOnce(ReturnFromAsyncTask(&cv));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Check expected MME state after failure
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

// Normal TAU sent in idle mode with eps bearer context status IE with inactive
// default bearer
TEST_F(MmeAppProcedureTest,
       TestNormalTauReqWithEpsBearerCtxStsInactiveDefBearer) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);
  uint8_t ebi_idx = 0;

  MME_APP_EXPECT_CALLS(3, 2, 2, 1, 1, 1, 2, 3, 1, 2, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send PDN connectivity request
  send_mme_app_uplink_data_ind(nas_msg_pdn_connectivity_req_ims,
                               sizeof(nas_msg_pdn_connectivity_req_ims), plmn);

  ebi_idx++;
  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  send_create_session_resp(REQUEST_ACCEPTED, DEFAULT_LBI + ebi_idx);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(DEFAULT_LBI + ebi_idx);

  // Send activate default eps bearer accept
  send_mme_app_uplink_data_ind(nas_msg_act_def_bearer_acc,
                               sizeof(nas_msg_act_def_bearer_acc), plmn);

  // Send modify bearer response for secondary PDN
  std::vector<int> b_modify = {6};
  std::vector<int> b_rm = {};
  send_modify_bearer_resp(b_modify, b_rm);

  // Add dedicated bearer for LBI 6
  // s1ap_generate_s1ap_e_rab_setup_req is called twice. Once for default
  // bearer and once for dedicated bearer
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(2);

  send_s11_create_bearer_req(DEFAULT_LBI + ebi_idx);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(7);

  // Send activate dedicated bearer accept mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_activate_ded_bearer_accept_ebi_7,
                               sizeof(nas_msg_activate_ded_bearer_accept_ebi_7),
                               plmn);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response
  // to mme_app mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending Normal TAU Request with EPS bearer context status
  send_mme_app_initial_ue_msg(
      nas_msg_tau_req_with_eps_bearer_ctx_sts_def_ber,
      sizeof(nas_msg_tau_req_with_eps_bearer_ctx_sts_def_ber), plmn, guti, 1);

  // Wait for spgw to send delete session request
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  send_delete_session_resp(DEFAULT_LBI + ebi_idx);

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW
  send_modify_bearer_resp(b_modify, b_rm);

  // Check MME state after Modify Bearer Response
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

// Normal TAU sent in idle mode with eps bearer context status IE with inactive
// dedicated bearer
TEST_F(MmeAppProcedureTest,
       TestNormalTauReqWithEpsBearerCtxStsInactiveDedBearer) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);
  uint8_t ebi_idx = 0;

  MME_APP_EXPECT_CALLS(3, 2, 2, 1, 1, 1, 2, 4, 1, 2, 4);

  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send PDN connectivity request
  send_mme_app_uplink_data_ind(nas_msg_pdn_connectivity_req_ims,
                               sizeof(nas_msg_pdn_connectivity_req_ims), plmn);

  ebi_idx++;
  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  send_create_session_resp(REQUEST_ACCEPTED, DEFAULT_LBI + ebi_idx);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(DEFAULT_LBI + ebi_idx);

  // Send activate default eps bearer accept
  send_mme_app_uplink_data_ind(nas_msg_act_def_bearer_acc,
                               sizeof(nas_msg_act_def_bearer_acc), plmn);

  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW
  std::vector<int> b_modify = {6};
  std::vector<int> b_rm = {};
  send_modify_bearer_resp(b_modify, b_rm);

  // Add dedicated bearer for LBI 6
  // s1ap_generate_s1ap_e_rab_setup_req is called twice. Once for default
  // bearer and once for dedicated bearer
  EXPECT_CALL(*s1ap_handler, s1ap_generate_s1ap_e_rab_setup_req()).Times(2);

  send_s11_create_bearer_req(DEFAULT_LBI + ebi_idx);

  // Send ERAB Setup Response mimicing S1AP
  send_erab_setup_rsp(7);

  // Send activate dedicated eps bearer accept
  send_mme_app_uplink_data_ind(nas_msg_activate_ded_bearer_accept_ebi_7,
                               sizeof(nas_msg_activate_ded_bearer_accept_ebi_7),
                               plmn);

  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response
  // to mme_app mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  mme_state_p = magma::lte::MmeNasStateManager::getInstance().get_state(false);
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  // Constructing and sending Normal TAU Request with EPS bearer context status
  send_mme_app_initial_ue_msg(
      nas_msg_tau_req_with_eps_bearer_ctx_sts_ded_ber,
      sizeof(nas_msg_tau_req_with_eps_bearer_ctx_sts_ded_ber), plmn, guti, 1);

  // Constructing and sending deactivate bearer request
  uint8_t ebi_to_be_deactivated = 7;
  send_s11_deactivate_bearer_req(1, &ebi_to_be_deactivated, false);

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW
  send_modify_bearer_resp(b_modify, b_rm);
  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW
  b_modify = {6};
  send_modify_bearer_resp(b_modify, b_rm);

  // Check MME state after Modify Bearer Response
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 2);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 2);

  // Constructing and sending Detach Request to mme_app
  // mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_detach_req, sizeof(nas_msg_detach_req),
                               plmn);

  // Constructing and sending Delete Session Response for each session
  // to mme_app mimicing SPGW task
  // Wait for 2 delete session requests
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  send_delete_session_resp(DEFAULT_LBI);
  send_delete_session_resp(6);

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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);
}

TEST_F(MmeAppProcedureTest, TestS1HandoverSuccess) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 2, 0, 1, 3);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  uint32_t new_enb_ue_s1ap_id = DEFAULT_eNB_S1AP_UE_ID + 1;
  uint32_t new_sctp_assoc_id = DEFAULT_SCTP_ASSOC_ID + 1;
  uint32_t new_enb_id = DEFAULT_ENB_ID + 1;
  uint32_t mme_ue_s1ap_id = 1;

  // Send Handover Required to mme_app mimicing S1AP, use 1 for all new ids,
  // since Initial UE message is using 0

  // Expect that MME Handover Request is sent to S1AP from mme_app
  EXPECT_CALL(*s1ap_handler, s1ap_mme_handle_handover_request(
                                 check_params_in_mme_app_handover_request(
                                     mme_ue_s1ap_id, DEFAULT_SCTP_ASSOC_ID)))
      .Times(1);

  send_s1ap_handover_required(DEFAULT_SCTP_ASSOC_ID, new_enb_id,
                              new_enb_ue_s1ap_id, mme_ue_s1ap_id);

  // Send Handover Request Ack to mme_app mimicing S1AP
  // Expect that MME Handover Command is sent to S1AP from mme_app
  EXPECT_CALL(*s1ap_handler, s1ap_mme_handle_handover_command(
                                 check_params_in_mme_app_handover_command(
                                     mme_ue_s1ap_id, new_enb_id)))
      .Times(1);
  send_s1ap_handover_request_ack(DEFAULT_SCTP_ASSOC_ID, DEFAULT_ENB_ID,
                                 new_enb_id, DEFAULT_eNB_S1AP_UE_ID,
                                 new_enb_ue_s1ap_id, mme_ue_s1ap_id);

  // Send Handover Notify to mme_app mimicing S1AP
  send_s1ap_handover_notify(new_sctp_assoc_id, DEFAULT_ENB_ID, new_enb_id,
                            DEFAULT_eNB_S1AP_UE_ID, new_enb_ue_s1ap_id,
                            mme_ue_s1ap_id);

  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW, with same parameters as last one
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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 2);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

TEST_F(MmeAppProcedureTest, TestDuplicateAttach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(6, 2, 2, 2, 2, 1, 2, 2, 1, 2, 4);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

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

  // Destruction at tear down is not sufficient as ICS occurs
  // twice in this test case.
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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  /* Move UE to ECM_IDLE mode */
  // Send context release request mimicing S1AP
  send_context_release_req(S1AP_RADIO_EUTRAN_GENERATED_REASON, TASK_S1AP);

  // Constructing and sending Release Access Bearer Response to mme_app
  // mimicing SPGW
  sgw_send_release_access_bearer_response(REQUEST_ACCEPTED);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending CONTEXT RELEASE COMPLETE to mme_app
  // mimicing S1AP task
  send_ue_ctx_release_complete();

  // Check MME state after context release request is processed
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 1);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);

  /* Duplicate Attach Request */
  // Constructing and sending Duplicate Attach Request to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_imsi_attach_req,
                               sizeof(nas_msg_imsi_attach_req), plmn);

  // An implicit detach event should occur leading to Delete Session Request
  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  send_delete_session_resp(DEFAULT_LBI);

  // Now the new attach request should proceed
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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

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

  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW
  send_modify_bearer_resp(b_modify, b_rm);

  // Check MME state after Modify Bearer Response
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 1);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 1);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 1);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}
// Test case validates the handling of cancel location request,
// which initiates network initiated detach
TEST_F(MmeAppProcedureTest, TestCLRNwInitiatedDetach) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(4, 1, 1, 1, 1, 0, 1, 1, 0, 1, 2);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

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
  EXPECT_CALL(*s6a_handler, s6a_cancel_location_ans()).Times(1);
  send_s6a_clr(imsi);
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

  // Wait for context release request
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

// Test case validates the handling of S6a Reset message,
// which sends update location request
TEST_F(MmeAppProcedureTest, TestS6aReset) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 2, 1, 1, 1, 0, 1, 2);
  // Setting the 3422 and 3460 timers to standard duration
  mme_config.nas_config.t3422_msec = 8000;
  mme_config.nas_config.t3460_msec = 8000;

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

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

  // Wait for ICS Request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

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

  send_s6a_reset();
  // ULR is sent in response S6a reset message handling
  // Sending ULA to mme_app mimicing successful S6A response for ULR
  send_s6a_ula(imsi, true);

  detach_ue(cv, lock, mme_state_p, guti, false);
}
}  // namespace lte
}  // namespace magma
