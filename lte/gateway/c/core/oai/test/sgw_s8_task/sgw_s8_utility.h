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

#include <gtest/gtest.h>
#include <string>
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_state_manager.h"
#include "lte/gateway/c/core/oai/include/sgw_s8_state.h"
#include "../mock_tasks/mock_tasks.h"

    extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_s11_handlers.h"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.h"
#include "lte/gateway/c/core/oai/include/spgw_types.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/sgw_config.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.h"
#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u.h"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_defs.h"
}

#define SLEEP_AT_INITIALIZATION_TIME_MS 500
void fill_imsi(char* imsi);
void fill_itti_csreq(
    itti_s11_create_session_request_t* session_req_pP,
    uint8_t default_eps_bearer_id);
void fill_itti_csrsp(s8_create_session_response_t* csr_resp, uint32_t teid);

void fill_create_bearer_request(
    s8_create_bearer_request_t* cb_req, uint32_t teid,
    uint8_t default_eps_bearer_id);

void fill_create_bearer_response(
    itti_s11_nw_init_actv_bearer_rsp_t* cb_response, uint32_t teid,
    uint8_t eps_bearer_id, teid_t s1_u_sgw_fteid, gtpv2c_cause_value_t cause);

void fill_delete_bearer_response(
    itti_s11_nw_init_deactv_bearer_rsp_t* db_response,
    uint32_t s_gw_teid_s11_s4, uint8_t eps_bearer_id,
    gtpv2c_cause_value_t cause);
void fill_delete_bearer_request(
    s8_delete_bearer_request_t* db_req, uint32_t teid, uint8_t eps_bearer_id);
// Initialize config params
class SgwS8Config : public ::testing::Test {
 public:
  sgw_state_t* create_ue_context(mme_sgw_tunnel_t* sgw_s11_tunnel);

 protected:
  sgw_config_t* config =
      reinterpret_cast<sgw_config_t*>(calloc(1, sizeof(sgw_config_t)));
  uint64_t imsi64               = 1010000000001;
  uint8_t default_eps_bearer_id = 5;
  virtual void SetUp();
  void sgw_s8_config_init();
  virtual void TearDown() {
    bdestroy_wrapper(&config->itti_config.log_file);
    bdestroy_wrapper(&config->ipv4.if_name_S1u_S12_S4_up);
    bdestroy_wrapper(&config->ipv4.if_name_S5_S8_up);
    bdestroy_wrapper(&config->ipv4.if_name_S11);
    bdestroy_wrapper(&config->config_file);
    free(config);
  }
};

// Initialize config params
class SgwS8ConfigAndCreateMock : public ::testing::Test {
 public:
  sgw_state_t* create_ue_context(mme_sgw_tunnel_t* sgw_s11_tunnel);

 protected:
  sgw_config_t* config =
      reinterpret_cast<sgw_config_t*>(calloc(1, sizeof(sgw_config_t)));
  uint64_t imsi64               = 1010000000001;
  uint8_t default_eps_bearer_id = 5;
  virtual void SetUp();
  void sgw_s8_config_init();
  virtual void TearDown();
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
};
