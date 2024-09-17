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
#include <netinet/in.h>
#include <chrono>
#include <memory>
#include <thread>

#include "lte/gateway/c/core/oai/include/s11_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_procedures_test_fixture.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"

extern "C" {
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/queue.h"
#include "lte/gateway/c/core/oai/include/gx_messages_types.h"
#include "lte/gateway/c/core/oai/include/ip_forward_messages_types.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"
}

namespace magma {
namespace lte {

TEST_F(SPGWAppProcedureTest, TestCreateSessionSuccess) {
  spgw_state_t* spgw_state = get_spgw_state(false);

  // Create session
  create_default_session(spgw_state);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestCreateSessionIPAllocFailure) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(&sample_session_req_p, test_imsi_str,
                              DEFAULT_MME_S11_TEID, DEFAULT_BEARER_INDEX,
                              sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  status_code_e create_session_rc = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(create_session_rc, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  magma::proto_map_rc_t rc = sgw_cm_get_eps_bearer_entry(
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection(),
      DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt);

  ASSERT_TRUE(eps_bearer_ctxt.ue_ip_paa().ipv4_addr().size() ==
              UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW with status as failure
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(
      &test_ip_alloc_resp, SGI_STATUS_ERROR_ALL_DYNAMIC_ADDRESSES_OCCUPIED,
      ue_sgw_teid, DEFAULT_EPS_BEARER_ID, DEFAULT_UE_IP, DEFAULT_VLAN);
  sgw_handle_ip_allocation_rsp(spgw_state, &test_ip_alloc_resp, test_imsi64);

  // Verify that UE context was removed in SPGW state after CSR failure
  ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p == nullptr);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestCreateSessionPCEFFailure) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  // expect call to MME create session response
  itti_s11_create_session_request_t sample_session_req_p = {};
  fill_create_session_request(&sample_session_req_p, test_imsi_str,
                              DEFAULT_MME_S11_TEID, DEFAULT_BEARER_INDEX,
                              sample_default_bearer_context, test_plmn);

  // trigger create session req to SPGW
  status_code_e create_session_rc = sgw_handle_s11_create_session_request(
      spgw_state, &sample_session_req_p, test_imsi64);

  ASSERT_EQ(create_session_rc, RETURNok);

  // Verify that a UE context exists in SPGW state after CSR is received
  spgw_ue_context_t* ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p != nullptr);

  // Verify that teid is created
  ASSERT_FALSE(LIST_EMPTY(&ue_context_p->sgw_s11_teid_list));
  teid_t ue_sgw_teid =
      LIST_FIRST(&ue_context_p->sgw_s11_teid_list)->sgw_s11_teid;

  // Verify that no IP address is allocated for this UE
  magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p =
      sgw_cm_get_spgw_context(ue_sgw_teid);

  magma::lte::oai::SgwEpsBearerContext eps_bearer_ctxt;
  magma::proto_map_rc_t rc = sgw_cm_get_eps_bearer_entry(
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection(),
      DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt);

  ASSERT_TRUE(eps_bearer_ctxt.ue_ip_paa().ipv4_addr().size() ==
              UNASSIGNED_UE_IP);

  // send an IP alloc response to SPGW
  itti_ip_allocation_response_t test_ip_alloc_resp = {};
  fill_ip_allocation_response(&test_ip_alloc_resp, SGI_STATUS_OK, ue_sgw_teid,
                              DEFAULT_EPS_BEARER_ID, DEFAULT_UE_IP,
                              DEFAULT_VLAN);
  sgw_handle_ip_allocation_rsp(spgw_state, &test_ip_alloc_resp, test_imsi64);

  // check if IP address is allocated after this message is done
  sgw_cm_get_eps_bearer_entry(
      spgw_eps_bearer_ctxt_info_p->mutable_sgw_eps_bearer_context()
          ->mutable_pdn_connection(),
      DEFAULT_EPS_BEARER_ID, &eps_bearer_ctxt);
  struct in_addr ue_ipv4 = {};
  int ue_ip = DEFAULT_UE_IP;
  inet_pton(AF_INET, eps_bearer_ctxt.ue_ip_paa().ipv4_addr().c_str(), &ue_ipv4);
  ASSERT_TRUE(!(memcmp(&ue_ipv4, &ue_ip, sizeof(DEFAULT_UE_IP))));

  // send pcef create session response to SPGW
  itti_pcef_create_session_response_t sample_pcef_csr_resp;
  fill_pcef_create_session_response(&sample_pcef_csr_resp, PCEF_STATUS_FAILED,
                                    ue_sgw_teid, DEFAULT_EPS_BEARER_ID,
                                    SGI_STATUS_OK);

  // check if MME gets a create session response
  EXPECT_CALL(*mme_app_handler, mme_app_handle_create_sess_resp()).Times(1);

  spgw_handle_pcef_create_session_response(spgw_state, &sample_pcef_csr_resp,
                                           test_imsi64);

  // verify that spgw context for IMSI has been cleared
  ue_context_p = spgw_get_ue_context(test_imsi64);
  ASSERT_TRUE(ue_context_p == nullptr);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

TEST_F(SPGWAppProcedureTest, TestDeleteSessionSuccess) {
  spgw_state_t* spgw_state = get_spgw_state(false);
  status_code_e return_code = RETURNerror;

  // Create session
  teid_t ue_sgw_teid = create_default_session(spgw_state);

  // create sample delete session request
  itti_s11_delete_session_request_t sample_delete_session_request = {};
  fill_delete_session_request(&sample_delete_session_request,
                              DEFAULT_MME_S11_TEID, ue_sgw_teid,
                              DEFAULT_EPS_BEARER_ID, test_plmn);

  EXPECT_CALL(*mme_app_handler,
              mme_app_handle_delete_sess_rsp(check_cause_in_ds_rsp(
                  REQUEST_ACCEPTED, DEFAULT_MME_S11_TEID)))
      .Times(1);

  return_code = sgw_handle_delete_session_request(
      &sample_delete_session_request, test_imsi64);
  ASSERT_EQ(return_code, RETURNok);

  // verify SPGW state is cleared
  ASSERT_TRUE(is_num_ue_contexts_valid(0));
  ASSERT_TRUE(is_num_cp_teids_valid(test_imsi64, 0));

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

}  // namespace lte
}  // namespace magma
