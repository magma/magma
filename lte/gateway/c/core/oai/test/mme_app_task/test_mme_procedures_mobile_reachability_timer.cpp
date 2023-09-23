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

TEST_F(MmeAppProcedureTest, TestMobileReachabilityTimer) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  // Reduce timer 3412 duration for testing since mobile reachability
  // and implicit detach timers are calculated from that
  mme_config.nas_config.t3412_min = 1;
  mme_config.nas_config.t3412_msec =
      50 * MME_APP_TIMER_TO_MSEC;  // implicit detach after 2x of this value

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

  // Wait for delete session request
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  // Constructing and sending Delete Session Response to mme_app
  // mimicing SPGW task
  send_delete_session_resp(DEFAULT_LBI);

  // Wait for context release command
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));

  // Check MME state after implicit detach complete
  send_activate_message_to_mme_app();
  cv.wait_for(lock, std::chrono::milliseconds(STATE_MAX_WAIT_MS));
  EXPECT_EQ(mme_state_p->nb_ue_attached, 0);
  EXPECT_EQ(mme_state_p->nb_ue_connected, 0);
  EXPECT_EQ(mme_state_p->nb_default_eps_bearers, 0);
  EXPECT_EQ(mme_state_p->nb_ue_idle, 0);
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 0);
}

}  // namespace lte
}  // namespace magma
