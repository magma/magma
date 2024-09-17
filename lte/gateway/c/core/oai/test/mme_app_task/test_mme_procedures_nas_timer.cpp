
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

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_app_test_util.hpp"
#include "lte/gateway/c/core/oai/test/mme_app_task/mme_procedure_test_fixture.hpp"

namespace magma {
namespace lte {

TEST_F(MmeAppProcedureTest, TestImsiAttachExpiredNasTimers) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  nas_config_timer_reinit(&mme_config.nas_config, MME_APP_TIMER_TO_MSEC);

  MME_APP_EXPECT_CALLS(15, 1, 1, 1, 1, 1, 1, 1, 0, 1, 2);

  // Constructing and sending Initial Attach Request to mme_app mimicing S1AP
  send_mme_app_initial_ue_msg(nas_msg_imsi_attach_req,
                              sizeof(nas_msg_imsi_attach_req), plmn, guti, 1);

  // Sending AIA to mme_app mimicing successful S6A response for AIR
  send_authentication_info_resp(imsi, true);

  // Wait for DL NAS Transport up to retransmission limit
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }
  // Constructing and sending Authentication Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_auth_resp, sizeof(nas_msg_auth_resp),
                               plmn);

  // Wait for DL NAS Transport up to retransmission limit
  for (int i = 0; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }
  // Constructing and sending SMC Response to mme_app mimicing S1AP
  send_mme_app_uplink_data_ind(nas_msg_smc_resp, sizeof(nas_msg_smc_resp),
                               plmn);

  // Sending ULA to mme_app mimicing successful S6A response for ULR
  send_s6a_ula(imsi, true);

  // Constructing and sending Create Session Response to mme_app mimicing SPGW
  send_create_session_resp(REQUEST_ACCEPTED, DEFAULT_LBI);

  // Wait for ICS request to be sent
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Constructing and sending ICS Response to mme_app mimicing S1AP
  send_ics_response();

  // Constructing UE Capability Indication message to mme_app
  // mimicing S1AP with dummy radio capabilities
  send_ue_capabilities_ind();

  // Wait for DL NAS Transport up to retransmission limit.
  // The first Attach Accept was piggybacked on ICS Request.
  for (int i = 1; i < NAS_RETX_LIMIT; ++i) {
    cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  }
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
}  // namespace lte
}  // namespace magma
