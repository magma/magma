
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

TEST_F(MmeAppProcedureTest, TestX2HandoverPathSwitchFailure) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 2, 0, 1, 3);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send path switch request to mme_app mimicing S1AP, use 1 for all new ids,
  // since Initial UE message is using 0
  uint32_t new_enb_ue_s1ap_id = DEFAULT_eNB_S1AP_UE_ID + 1;
  uint32_t new_sctp_assoc_id = DEFAULT_SCTP_ASSOC_ID + 1;
  uint32_t new_enb_id = DEFAULT_ENB_ID + 1;

  // Expect that Path Switch Failure is sent to S1AP from mme_app after modify
  // bearer response is received
  EXPECT_CALL(*s1ap_handler, s1ap_handle_path_switch_req_failure(
                                 check_params_in_path_switch_req_failure(
                                     new_enb_ue_s1ap_id, new_sctp_assoc_id)))
      .Times(1);

  send_s1ap_path_switch_req(new_sctp_assoc_id, new_enb_id, new_enb_ue_s1ap_id,
                            plmn);

  // Constructing and sending Modify Bearer Response to mme_app
  // mimicing SPGW, with empty modify bearer list to trigger failure
  std::vector<int> b_modify = {};
  std::vector<int> b_rm = {};
  // b_modify = {};
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

TEST_F(MmeAppProcedureTest, TestX2HandoverPathSwitchSuccess) {
  mme_app_desc_t* mme_state_p =
      magma::lte::MmeNasStateManager::getInstance().get_state(false);
  std::condition_variable cv;
  std::mutex mx;
  std::unique_lock<std::mutex> lock(mx);

  MME_APP_EXPECT_CALLS(3, 1, 1, 1, 1, 1, 1, 2, 0, 1, 3);

  // Attach the UE
  guti = {0};
  attach_ue(cv, lock, mme_state_p, &guti);

  // Send path switch request to mme_app mimicing S1AP, use 1 for all new ids,
  // since Initial UE message is using 0
  uint32_t new_enb_ue_s1ap_id = DEFAULT_eNB_S1AP_UE_ID + 1;
  uint32_t new_sctp_assoc_id = DEFAULT_SCTP_ASSOC_ID + 1;
  uint32_t new_enb_id = DEFAULT_ENB_ID + 1;

  // Expect that Path Switch ACK is sent to S1AP from mme_app after modify
  // bearer response is received
  EXPECT_CALL(*s1ap_handler, s1ap_handle_path_switch_req_ack(
                                 check_params_in_path_switch_req_ack(
                                     new_enb_ue_s1ap_id, new_sctp_assoc_id)))
      .Times(1);

  send_s1ap_path_switch_req(new_sctp_assoc_id, new_enb_id, new_enb_ue_s1ap_id,
                            plmn);

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
  EXPECT_EQ(mme_state_p->nb_s1u_bearers, 1);

  detach_ue(cv, lock, mme_state_p, guti, false);
}

}  // namespace lte
}  // namespace magma
