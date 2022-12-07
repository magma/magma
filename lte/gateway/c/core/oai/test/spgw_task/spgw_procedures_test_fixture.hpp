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

#pragma once

#include <gtest/gtest.h>
#include <stdlib.h>
#include <string.h>
#include <chrono>
#include <cstdint>
#include <iostream>
#include <memory>
#include <string>
#include <thread>

extern "C" {
#include "lte/gateway/c/core/oai/include/spgw_config.h"
}

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"
#include "lte/gateway/c/core/oai/test/spgw_task/spgw_test_util.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_handlers.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"

extern bool hss_associated;

namespace magma {
namespace lte {

class SPGWAppProcedureTest : public ::testing::Test {
  virtual void SetUp();

  virtual void TearDown();

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::string test_imsi_str = "001010000000001";
  std::string invalid_imsi_str = "001010000000002";
  uint64_t test_imsi64 = 1010000000001;
  uint64_t test_invalid_imsi64 = 1010000000002;
  plmn_t test_plmn = {.mcc_digit2 = 0,
                      .mcc_digit1 = 0,
                      .mnc_digit3 = 0x0f,
                      .mcc_digit3 = 1,
                      .mnc_digit2 = 1,
                      .mnc_digit1 = 0};
  bearer_context_to_be_created_t sample_default_bearer_context = {
      .eps_bearer_id = 5,
      .bearer_level_qos = {.pci = 1,
                           .pl = 15,
                           .pvi = 0,
                           .qci = 9,
                           .gbr = {},
                           .mbr = {.br_ul = 200000000, .br_dl = 100000000}}};

  bearer_qos_t sample_dedicated_bearer_qos = {
      .pci = 1,
      .pl = 1,
      .pvi = 0,
      .qci = 1,
      .gbr = {.br_ul = 200000000, .br_dl = 100000000},
      .mbr = {.br_ul = 200000000, .br_dl = 100000000}};

  teid_t create_default_session(spgw_state_t* spgw_state);
  ebi_t activate_dedicated_bearer(
      spgw_state_t* spgw_state,
      magma::lte::oai::S11BearerContext* spgw_eps_bearer_ctxt_info_p,
      teid_t ue_sgw_teid);
  void deactivate_dedicated_bearer(spgw_state_t* spgw_state, teid_t ue_sgw_teid,
                                   ebi_t ded_eps_bearer_id);
};

MATCHER_P2(check_params_in_actv_bearer_req, lbi, tft, "") {
  auto cb_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_actv_bearer_request_t>(arg);
  if (cb_req_rcvd_at_mme.lbi != lbi) {
    return false;
  }
  if (!(cb_req_rcvd_at_mme.s1_u_sgw_fteid.teid)) {
    return false;
  }
  if ((memcmp(&cb_req_rcvd_at_mme.tft, &tft,
              sizeof(traffic_flow_template_t)))) {
    return false;
  }
  return true;
}

MATCHER_P2(check_params_in_deactv_bearer_req, num_bearers, eps_bearer_id_array,
           "") {
  auto db_req_rcvd_at_mme =
      static_cast<itti_s11_nw_init_deactv_bearer_request_t>(arg);
  if (db_req_rcvd_at_mme.no_of_bearers != num_bearers) {
    return false;
  }
  if (memcmp(db_req_rcvd_at_mme.ebi, eps_bearer_id_array,
             sizeof(db_req_rcvd_at_mme.ebi))) {
    return false;
  }
  return true;
}

MATCHER_P2(check_cause_in_ds_rsp, cause, teid, "") {
  auto ds_rsp_rcvd_at_mme =
      static_cast<itti_s11_delete_session_response_t>(arg);
  if (ds_rsp_rcvd_at_mme.cause.cause_value == cause) {
    return true;
  }
  if (ds_rsp_rcvd_at_mme.teid == teid) {
    return true;
  }
  return false;
}

MATCHER_P2(check_params_in_suspend_ack, return_val, teid, "") {
  auto suspend_ack_rcvd_at_mme =
      static_cast<itti_s11_suspend_acknowledge_t>(arg);
  if ((suspend_ack_rcvd_at_mme.cause.cause_value == return_val) &&
      (suspend_ack_rcvd_at_mme.teid == teid)) {
    return true;
  }
  return false;
}
}  // namespace lte
}  // namespace magma
