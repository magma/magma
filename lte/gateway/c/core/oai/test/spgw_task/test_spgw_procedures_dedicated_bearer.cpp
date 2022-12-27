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

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <cstdint>
#include <string>
#include <thread>

extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/queue.h"
#include "lte/gateway/c/core/oai/include/gx_messages_types.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
}

#include "lte/gateway/c/core/oai/include/s11_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_procedures_test_fixture.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"

namespace magma {
namespace lte {
TEST_F(SPGWAppProcedureTest, TestDedicatedBearerActivation) {
  spgw_state_t* spgw_state = get_spgw_state(false);

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // Activate dedicated bearer
  activate_dedicated_bearer(spgw_state, spgw_eps_bearer_ctxt_info_p,
                            ue_sgw_teid);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDedicatedBearerDeactivation) {
  spgw_state_t* spgw_state = get_spgw_state(false);

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // Activate dedicated bearer
  ebi_t ded_eps_bearer_id = activate_dedicated_bearer(
      spgw_state, spgw_eps_bearer_ctxt_info_p, ue_sgw_teid);

  // check that MME gets a bearer deactivation request
  uint32_t expected_no_of_bearers = 1;
  ebi_t expected_ebi[] = {6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};
  EXPECT_CALL(*mme_app_handler, mme_app_handle_nw_init_bearer_deactv_req(
                                    check_params_in_deactv_bearer_req(
                                        expected_no_of_bearers, expected_ebi)))
      .Times(1);

  // Deactivate dedicated bearer
  deactivate_dedicated_bearer(spgw_state, ue_sgw_teid, ded_eps_bearer_id);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest,
       TestDedicatedBearerDeactivationDeleteDefaultBearer) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  sgw_cm_get_eps_bearer_entry(
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection(),
      DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt);

  // Activate dedicated bearer
  ebi_t ded_eps_bearer_id = activate_dedicated_bearer(
      spgw_state, spgw_eps_bearer_ctxt_info_p, ue_sgw_teid);

  // send deactivate request for dedicated bearer from Session Manager
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  // check that MME gets a bearer deactivation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_bearer_deactv_req(
                  check_params_in_deactv_bearer_req(
                      sample_gx_nw_init_ded_bearer_deactv_req.no_of_bearers,
                      sample_gx_nw_init_ded_bearer_deactv_req.ebi)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_deactv_req(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // send a delete dedicated bearer response from MME
  itti_s11_nw_init_deactv_bearer_rsp_t sample_nw_init_ded_bearer_deactv_resp =
      {};
  int num_bearers_to_delete = 2;
  ebi_t eps_bearer_id_array[] = {DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id};

  fill_nw_initiated_deactivate_bearer_response(
      &sample_nw_init_ded_bearer_deactv_resp, test_imsi64, true,
      REQUEST_ACCEPTED, eps_bearer_id_array, num_bearers_to_delete,
      ue_sgw_teid);
  return_code = sgw_handle_nw_initiated_deactv_bearer_rsp(
      spgw_state, &sample_nw_init_ded_bearer_deactv_resp, test_imsi64);
  EXPECT_EQ(return_code, RETURNok);

  // check that session is removed
  EXPECT_TRUE(is_num_ue_contexts_valid(0));
  EXPECT_TRUE(is_num_cp_teids_valid(test_imsi64, 0));

  free(sample_nw_init_ded_bearer_deactv_resp.lbi);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDedicatedBearerActivationInvalidImsiLbi) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  EXPECT_EQ((sgw_cm_get_eps_bearer_entry(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
                    ->mutable_pdn_connection(),
                DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt)),
            magma::PROTO_MAP_OK);

  // send network initiated dedicated bearer activation request with
  // invalid imsi
  itti_gx_nw_init_actv_bearer_request_t sample_gx_nw_init_ded_bearer_actv_req =
      {};
  gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, invalid_imsi_str,
      DEFAULT_EPS_BEARER_ID, sample_dedicated_bearer_qos);

  // check that MME bearer activation request is not sent to MME
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_ded_bearer_actv_req(
                  check_params_in_actv_bearer_req(
                      sample_gx_nw_init_ded_bearer_actv_req.lbi,
                      sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(0);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNerror);
  EXPECT_EQ(failed_cause, REQUEST_REJECTED);

  // check number of pending procedures
  EXPECT_EQ(get_num_pending_create_bearer_procedures(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()),
            0);

  // check that dedicated bearer is not created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // send network initiated dedicated bearer activation request from Session
  // Manager with invalid lbi
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID + 1, sample_dedicated_bearer_qos);

  // check that MME bearer activation request is not sent to MME
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_ded_bearer_actv_req(
                  check_params_in_actv_bearer_req(
                      sample_gx_nw_init_ded_bearer_actv_req.lbi,
                      sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(0);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNerror);
  EXPECT_EQ(failed_cause, REQUEST_REJECTED);

  // check number of pending procedures
  EXPECT_EQ(get_num_pending_create_bearer_procedures(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()),
            0);

  // check that dedicated bearer is not created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDedicatedBearerDeactivationInvalidImsiLbi) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  EXPECT_EQ((sgw_cm_get_eps_bearer_entry(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
                    ->mutable_pdn_connection(),
                DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt)),
            magma::PROTO_MAP_OK);

  // Activate dedicated bearer
  ebi_t ded_eps_bearer_id = activate_dedicated_bearer(
      spgw_state, spgw_eps_bearer_ctxt_info_p, ue_sgw_teid);

  // send deactivate request for dedicated bearer from Session Manager
  // with invalid imsi
  itti_gx_nw_init_deactv_bearer_request_t
      sample_gx_nw_init_ded_bearer_deactv_req = {};
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, invalid_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id);

  // check that MME does not get bearer deactivation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_bearer_deactv_req(
                  check_params_in_deactv_bearer_req(
                      sample_gx_nw_init_ded_bearer_deactv_req.no_of_bearers,
                      sample_gx_nw_init_ded_bearer_deactv_req.ebi)))
      .Times(0);

  return_code = spgw_handle_nw_initiated_bearer_deactv_req(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // check that there are 2 bearers
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // send deactivate request for dedicated bearer from Session Manager
  // with invalid bearer id
  fill_nw_initiated_deactivate_bearer_request(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, ded_eps_bearer_id + 1);

  // check that MME does not get bearer deactivation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_bearer_deactv_req(
                  check_params_in_deactv_bearer_req(
                      sample_gx_nw_init_ded_bearer_deactv_req.no_of_bearers,
                      sample_gx_nw_init_ded_bearer_deactv_req.ebi)))
      .Times(0);

  return_code = spgw_handle_nw_initiated_bearer_deactv_req(
      &sample_gx_nw_init_ded_bearer_deactv_req, test_imsi64);

  EXPECT_EQ(return_code, RETURNok);

  // verify that the dedicated bearer is not deleted
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 2));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDedicatedBearerActivationReject) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  // send network initiated dedicated bearer activation request from Session
  // Manager
  itti_gx_nw_init_actv_bearer_request_t sample_gx_nw_init_ded_bearer_actv_req =
      {};
  gtpv2c_cause_value_t failed_cause = REQUEST_ACCEPTED;
  fill_nw_initiated_activate_bearer_request(
      &sample_gx_nw_init_ded_bearer_actv_req, test_imsi_str,
      DEFAULT_EPS_BEARER_ID, sample_dedicated_bearer_qos);

  // check that MME gets a bearer activation request
  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_nw_init_ded_bearer_actv_req(
                  check_params_in_actv_bearer_req(
                      sample_gx_nw_init_ded_bearer_actv_req.lbi,
                      sample_gx_nw_init_ded_bearer_actv_req.ul_tft)))
      .Times(1);

  return_code = spgw_handle_nw_initiated_bearer_actv_req(
      spgw_state, &sample_gx_nw_init_ded_bearer_actv_req, test_imsi64,
      &failed_cause);

  EXPECT_EQ(return_code, RETURNok);

  magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p =
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context();
  // check number of pending procedures
  EXPECT_EQ(get_num_pending_create_bearer_procedures(
                spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()),
            1);

  // fetch new SGW teid for the pending bearer procedure
  teid_t ue_ded_bearer_sgw_teid = 0;
  for (uint8_t proc_index = 0;
       proc_index < sgw_context_p->pending_procedures_size(); proc_index++) {
    magma::lte::oai::PgwCbrProcedure pgw_ni_cbr_proc =
        sgw_context_p->pending_procedures(proc_index);
    magma::lte::oai::SgwEpsBearerContext bearer_context_proto =
        pgw_ni_cbr_proc.pending_eps_bearers(0);
    ue_ded_bearer_sgw_teid = bearer_context_proto.sgw_teid_s1u_s12_s4_up();
  }

  // send bearer activation response from MME with cause=REQUEST_REJECTED
  ebi_t ded_eps_bearer_id = DEFAULT_EPS_BEARER_ID + 1;
  itti_s11_nw_init_actv_bearer_rsp_t sample_nw_init_ded_bearer_actv_resp = {};
  fill_nw_initiated_activate_bearer_response(
      &sample_nw_init_ded_bearer_actv_resp, DEFAULT_MME_S11_TEID, ue_sgw_teid,
      ue_ded_bearer_sgw_teid, DEFAULT_ENB_GTP_TEID + 1, ded_eps_bearer_id,
      REQUEST_REJECTED, test_plmn);
  return_code = sgw_handle_nw_initiated_actv_bearer_rsp(
      &sample_nw_init_ded_bearer_actv_resp, test_imsi64);

  EXPECT_EQ(return_code, RETURNerror);

  // check that the dedicated bearer is not created
  EXPECT_TRUE(is_num_s1_bearers_valid(ue_sgw_teid, 1));

  // check that no pending procedure is left
  EXPECT_EQ(get_num_pending_create_bearer_procedures(sgw_context_p), 0);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}
}  // namespace lte
}  // namespace magma
