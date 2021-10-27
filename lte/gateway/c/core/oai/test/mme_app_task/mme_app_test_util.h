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
extern "C" {
#include "mme_config.h"
}

namespace magma {
namespace lte {

#define MME_APP_TIMER_TO_MSEC 10
#define END_OF_TEST_SLEEP_MS 500
#define STATE_MAX_WAIT_MS 1000
#define NAS_RETX_LIMIT 5

#define MME_APP_EXPECT_CALLS(                                                  \
    dlNas, connEstConf, ctxRel, air, ulr, purgeReq, csr, dsr, setAppHealth)    \
  do {                                                                         \
    EXPECT_CALL(*s1ap_handler, s1ap_generate_downlink_nas_transport())         \
        .Times(dlNas)                                                          \
        .WillRepeatedly(ReturnFromAsyncTask(&cv));                             \
    EXPECT_CALL(*s1ap_handler, s1ap_handle_conn_est_cnf()).Times(connEstConf); \
    EXPECT_CALL(*s1ap_handler, s1ap_handle_ue_context_release_command())       \
        .Times(ctxRel);                                                        \
    EXPECT_CALL(*s6a_handler, s6a_viface_authentication_info_req())            \
        .Times(air);                                                           \
    EXPECT_CALL(*s6a_handler, s6a_viface_update_location_req()).Times(ulr);    \
    EXPECT_CALL(*s6a_handler, s6a_viface_purge_ue()).Times(purgeReq);          \
    EXPECT_CALL(*spgw_handler, sgw_handle_s11_create_session_request())        \
        .Times(csr);                                                           \
    EXPECT_CALL(*spgw_handler, sgw_handle_delete_session_request())            \
        .Times(dsr)                                                            \
        .WillRepeatedly(ReturnFromAsyncTask(&cv));                             \
    EXPECT_CALL(*service303_handler, service303_set_application_health())      \
        .Times(setAppHealth)                                                   \
        .WillRepeatedly(ReturnFromAsyncTask(&cv));                             \
  } while (0)

void nas_config_timer_reinit(nas_config_t* nas_conf, uint32_t timeout_msec);

void send_sctp_mme_server_initialized();

void send_activate_message_to_mme_app();

void send_mme_app_initial_ue_msg(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn);

void send_mme_app_uplink_data_ind(
    const uint8_t* nas_msg, uint8_t nas_msg_length, const plmn_t& plmn);

void send_authentication_info_resp(const std::string& imsi, bool success);

void send_s6a_ula(const std::string& imsi, bool success);

void send_create_session_resp();

void send_delete_session_resp();

void send_ics_response();

void send_ue_ctx_release_complete();

void send_ue_capabilities_ind();

}  // namespace lte
}  // namespace magma
