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
#include "sgw_s8_state_manager.h"
#include "sgw_s8_state.h"

extern "C" {
#include "log.h"
#include "sgw_s8_s11_handlers.h"
#include "sgw_handlers.h"
#include "spgw_types.h"
#include "s11_messages_types.h"
#include "common_types.h"
#include "sgw_config.h"
#include "dynamic_memory_check.h"
#include "sgw_context_manager.h"
#include "gtpv1u.h"
}

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
    uint8_t eps_bearer_id, teid_t s1_u_sgw_fteid);

// Initialize config params
class SgwS8Config : public ::testing::Test {
 public:
  sgw_state_t* create_ue_context(mme_sgw_tunnel_t* sgw_s11_tunnel);
  void sgw_initialize_gtpv1u(void);
  void sgw_uninitialize_gtpv1u(void);

 protected:
  sgw_config_t* config =
      reinterpret_cast<sgw_config_t*>(calloc(1, sizeof(sgw_config_t)));
  uint64_t imsi64               = 1010000000001;
  uint8_t default_eps_bearer_id = 5;
  virtual void SetUp() {
    config->itti_config.queue_size     = 0;
    std::string file_string            = "/var/opt/magma/tmp/spgw.conf";
    config->itti_config.log_file       = bfromcstr(file_string.c_str());
    std::string s1u_if_name            = "eth1";
    config->ipv4.if_name_S1u_S12_S4_up = bfromcstr(s1u_if_name.c_str());
    config->ipv4.S1u_S12_S4_up.s_addr  = 0x8e3ca8c0;
    config->ipv4.netmask_S1u_S12_S4_up = 24;
    std::string s5s8u_if_name          = "eth0";
    config->ipv4.if_name_S5_S8_up      = bfromcstr(s5s8u_if_name.c_str());
    config->ipv4.S5_S8_up.s_addr       = 0xf02000a;
    config->ipv4.netmask_S5_S8_up      = 24;
    std::string s11                    = "lo";
    config->ipv4.if_name_S11           = bfromcstr(s11.c_str());
    config->ipv4.S11.s_addr            = 0x100007f;
    config->ipv4.netmask_S11           = 8;
    config->udp_port_S1u_S12_S4_up     = 2152;
    config->config_file                = bfromcstr(file_string.c_str());
  }
  virtual void TearDown() {
    bdestroy_wrapper(&config->itti_config.log_file);
    bdestroy_wrapper(&config->ipv4.if_name_S1u_S12_S4_up);
    bdestroy_wrapper(&config->ipv4.if_name_S5_S8_up);
    bdestroy_wrapper(&config->ipv4.if_name_S11);
    bdestroy_wrapper(&config->config_file);
    free(config);
  }
};
