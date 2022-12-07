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

#include <gtest/gtest.h>
#include <cstdint>
#include <string>
#include <thread>

#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_procedures_test_fixture.hpp"

namespace magma {
namespace lte {
TEST_F(SPGWAppProcedureTest, TestModifyBearerFailure) {
  status_code_e return_code = RETURNerror;

  // create sample modify default bearer request
  itti_s11_modify_bearer_request_t sample_modify_bearer_req = {};
  fill_modify_bearer_request(&sample_modify_bearer_req, DEFAULT_MME_S11_TEID,
                             ERROR_SGW_S11_TEID, DEFAULT_ENB_GTP_TEID,
                             DEFAULT_BEARER_INDEX, DEFAULT_EPS_BEARER_ID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_modify_bearer_rsp()).Times(1);
  return_code =
      sgw_handle_modify_bearer_request(&sample_modify_bearer_req, test_imsi64);

  ASSERT_EQ(return_code, RETURNok);

  // verify that no session exists in SPGW state
  ASSERT_TRUE(is_num_ue_contexts_valid(0));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestReleaseBearerSuccess) {
  spgw_state_t* spgw_state = get_spgw_state(false);

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  // send release access bearer request
  itti_s11_release_access_bearers_request_t sample_release_bearer_req = {};
  fill_release_access_bearer_request(&sample_release_bearer_req,
                                     DEFAULT_MME_S11_TEID, ue_sgw_teid);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_release_access_bearers_resp())
      .Times(1);

  sgw_handle_release_access_bearers_request(&sample_release_bearer_req,
                                            test_imsi64);

  // verify that eNB information has been cleared
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestReleaseBearerWithInvalidImsi64) {
  spgw_state_t* spgw_state = get_spgw_state(false);

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  // send release access bearer request
  itti_s11_release_access_bearers_request_t sample_release_bearer_req = {};
  fill_release_access_bearer_request(&sample_release_bearer_req,
                                     DEFAULT_MME_S11_TEID, ERROR_SGW_S11_TEID);

  EXPECT_CALL(*mme_app_handler, mme_app_handle_release_access_bearers_resp())
      .Times(1);

  // Send wrong IMSI so that spgw will not be able to fetch and delete
  // the context
  sgw_handle_release_access_bearers_request(&sample_release_bearer_req,
                                            test_invalid_imsi64);

  // verify that eNB information has not been cleared
  ASSERT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestSuspendNotification) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // trigger suspend notification to SPGW task
  itti_s11_suspend_notification_t sample_suspend_notification = {};
  fill_s11_suspend_notification(&sample_suspend_notification, ue_sgw_teid,
                                test_imsi_str, DEFAULT_EPS_BEARER_ID);

  // verify that mock MME app task receives an acknowledgement with
  // REQUEST_ACCEPTED
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_suspend_acknowledge(check_params_in_suspend_ack(
                  REQUEST_ACCEPTED,
                  spgw_eps_bearer_ctxt_info_p->sgw_eps_bearer_context()
                      .mme_teid_s11())))
      .Times(1);
  return_code = sgw_handle_suspend_notification(&sample_suspend_notification,
                                                test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDeleteBearerCommand) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // Activate dedicated bearer
  ebi_t ded_eps_bearer_id = activate_dedicated_bearer(
      spgw_state, spgw_eps_bearer_ctxt_info_p, ue_sgw_teid);

  // create and send delete bearer command to SPGW task
  itti_s11_delete_bearer_command_t s11_delete_bearer_command = {};
  fill_s11_delete_bearer_command(&s11_delete_bearer_command, ue_sgw_teid,
                                 DEFAULT_MME_S11_TEID, ded_eps_bearer_id);

  // check that MME gets a bearer deactivation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_bearer_deactv_req(
                  check_params_in_deactv_bearer_req(
                      1, s11_delete_bearer_command.ebi_list.ebis)))
      .Times(1);

  // Trigger delete bearer command
  sgw_handle_delete_bearer_cmd(&s11_delete_bearer_command, test_imsi64);

  // send a delete dedicated bearer response from MME
  itti_s11_nw_init_deactv_bearer_rsp_t sample_nw_init_ded_bearer_deactv_resp =
      {};
  int num_bearers_to_delete = 1;
  ebi_t eps_bearer_id_array[] = {ded_eps_bearer_id};

  fill_nw_initiated_deactivate_bearer_response(
      &sample_nw_init_ded_bearer_deactv_resp, test_imsi64, false,
      REQUEST_ACCEPTED, eps_bearer_id_array, num_bearers_to_delete,
      ue_sgw_teid);
  return_code = sgw_handle_nw_initiated_deactv_bearer_rsp(
      spgw_state, &sample_nw_init_ded_bearer_deactv_resp, test_imsi64);
  EXPECT_EQ(return_code, RETURNok);

  // check that the dedicated bearer is deleted
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}
}  // namespace lte
}  // namespace magma
