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
#include <string>
#include <vector>
#include <gmock/gmock-matchers.h>

extern "C" {
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsMobileIdentity.h"
}

using std::vector;

namespace magma {
namespace lte {

#define MME_APP_TIMER_TO_MSEC 10
#define STATE_MAX_WAIT_MS 10000
#define NAS_RETX_LIMIT 5

#define MME_APP_EXPECT_CALLS(                                                  \
    dlNas, connEstConf, ctxRel, air, ulr, purgeReq, csr, mbr, relBearer, dsr,  \
    setAppHealth)                                                              \
  do {                                                                         \
    EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport)           \
        .Times(dlNas)                                                          \
        .WillRepeatedly(                                                       \
            DoAll(SaveArg<0>(&msg_nas_dl_data), ReturnFromAsyncTask(&cv)));    \
    EXPECT_CALL(*s1ap_handler, s1ap_handle_conn_est_cnf(_))                    \
        .Times(connEstConf)                                                    \
        .WillRepeatedly(SaveArg<0>(&nas_msg));                                 \
    EXPECT_CALL(*s1ap_handler, s1ap_handle_ue_context_release_command())       \
        .Times(ctxRel)                                                         \
        .WillRepeatedly(ReturnFromAsyncTask(&cv));                             \
    EXPECT_CALL(*s6a_handler, s6a_viface_authentication_info_req())            \
        .Times(air);                                                           \
    EXPECT_CALL(*s6a_handler, s6a_viface_update_location_req()).Times(ulr);    \
    EXPECT_CALL(*s6a_handler, s6a_viface_purge_ue()).Times(purgeReq);          \
    EXPECT_CALL(*spgw_handler, sgw_handle_s11_create_session_request())        \
        .Times(csr);                                                           \
    EXPECT_CALL(*spgw_handler, sgw_handle_modify_bearer_request()).Times(mbr); \
    EXPECT_CALL(*spgw_handler, sgw_handle_release_access_bearers_request())    \
        .Times(relBearer);                                                     \
    EXPECT_CALL(*spgw_handler, sgw_handle_delete_session_request())            \
        .Times(dsr)                                                            \
        .WillRepeatedly(ReturnFromAsyncTask(&cv));                             \
    EXPECT_CALL(*service303_handler, service303_set_application_health())      \
        .Times(setAppHealth)                                                   \
        .WillRepeatedly(ReturnFromAsyncTask(&cv));                             \
  } while (0)

#define EXPECT_ARRAY_EQ(orig_array, expected_array, len)                       \
  ASSERT_THAT(                                                                 \
      vector<uint8_t>(expected_array, expected_array + len),                   \
      ::testing::ElementsAreArray(orig_array));

void nas_config_timer_reinit(nas_config_t* nas_conf, uint32_t timeout_msec);

void send_sctp_mme_server_initialized();

void send_activate_message_to_mme_app();

void send_mme_app_initial_ue_msg(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn,
    guti_eps_mobile_identity_t& guti, tac_t tac);

void send_mme_app_uplink_data_ind(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn);

void send_authentication_info_resp(const std::string& imsi, bool success);

void send_s6a_ula(const std::string& imsi, bool success);

void send_create_session_resp(gtpv2c_cause_value_t cause_value);

void send_delete_session_resp();

void send_ics_response();

void send_ics_failure();

void send_ue_ctx_release_complete();

void send_ue_capabilities_ind();

void send_context_release_req(s1cause rel_cause, task_id_t TASK_ID);

void send_modify_bearer_resp(
    const std::vector<int>& bearer_to_modify,
    const std::vector<int>& bearer_to_remove);

void sgw_send_release_access_bearer_response(gtpv2c_cause_value_t cause);

void send_s11_deactivate_bearer_req(
    uint8_t no_of_bearers_to_be_deact, uint8_t* ebi_to_be_deactivated,
    bool delete_default_bearer);

void send_s11_create_bearer_req();

void send_erab_setup_rsp();

void send_erab_release_rsp();

void send_paging_request();

}  // namespace lte
}  // namespace magma
