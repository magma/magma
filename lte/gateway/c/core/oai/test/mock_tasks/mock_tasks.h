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
#include <gmock/gmock.h>
#include <gtest/gtest.h>

extern "C" {
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#define CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_init.h"
#undef CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
}

const task_info_t tasks_info[] = {
    {THREAD_NULL, "TASK_UNKNOWN", "ipc://IPC_TASK_UNKNOWN"},
#define TASK_DEF(tHREADiD)                                                     \
  {THREAD_##tHREADiD, #tHREADiD, "ipc://IPC_" #tHREADiD},
#include "lte/gateway/c/core/oai/include/tasks_def.h"
#undef TASK_DEF
};

/* Map message id to message information */
const message_info_t messages_info[] = {
#define MESSAGE_DEF(iD, sTRUCT, fIELDnAME) {iD, sizeof(sTRUCT), #iD},
#include "lte/gateway/c/core/oai/include/messages_def.h"
#undef MESSAGE_DEF
};

#define TEST_GRPCSERVICES_SERVER_ADDRESS "127.0.0.1:50095"

#define END_OF_TESTCASE_SLEEP_MS 100
#define SLEEP_AT_INITIALIZATION_TIME_MS 500
class MockS1apHandler {
 public:
  MOCK_METHOD1(
      s1ap_generate_downlink_nas_transport,
      void(itti_s1ap_nas_dl_data_req_t cb_req));
  MOCK_METHOD1(s1ap_handle_conn_est_cnf, void(bstring nas_pdu));
  MOCK_METHOD0(s1ap_handle_ue_context_release_command, void());
  MOCK_METHOD0(s1ap_generate_s1ap_e_rab_setup_req, void());
  MOCK_METHOD0(s1ap_generate_s1ap_e_rab_rel_cmd, void());
  MOCK_METHOD0(s1ap_handle_paging_request, void());
  MOCK_METHOD1(
      s1ap_handle_path_switch_req_ack,
      bool(itti_s1ap_path_switch_request_ack_t path_switch_ack));
  MOCK_METHOD1(
      s1ap_handle_path_switch_req_failure,
      bool(itti_s1ap_path_switch_request_failure_t path_switch_fail));
};

class MockMmeAppHandler {
 public:
  MOCK_METHOD0(mme_app_handle_initial_ue_message, void());
  MOCK_METHOD0(mme_app_handle_s1ap_ue_context_release_req, void());
  MOCK_METHOD0(mme_app_handle_create_sess_resp, void());
  MOCK_METHOD1(
      mme_app_handle_nw_init_ded_bearer_actv_req,
      bool(itti_s11_nw_init_actv_bearer_request_t cb_req));
  MOCK_METHOD1(
      mme_app_handle_nw_init_bearer_deactv_req,
      bool(itti_s11_nw_init_deactv_bearer_request_t db_req));
  MOCK_METHOD0(mme_app_handle_modify_bearer_rsp, void());
  MOCK_METHOD1(
      mme_app_handle_delete_sess_rsp,
      bool(itti_s11_delete_session_response_t ds_rsp));
  MOCK_METHOD0(nas_proc_dl_transfer_rej, void());
  MOCK_METHOD0(mme_app_handle_release_access_bearers_resp, void());
  MOCK_METHOD0(mme_app_handle_handover_required, void());
  MOCK_METHOD0(mme_app_handle_initial_context_setup_failure, void());
  MOCK_METHOD0(mme_app_handle_enb_reset_req, void());
  MOCK_METHOD0(mme_app_handle_e_rab_setup_rsp, void());
  MOCK_METHOD0(mme_app_handle_path_switch_request, void());
  MOCK_METHOD1(
      mme_app_handle_suspend_acknowledge,
      bool(itti_s11_suspend_acknowledge_t suspend_ack));
};

class MockSctpHandler {
 public:
  MOCK_METHOD0(sctpd_send_dl, void());
};

class MockS6aHandler {
 public:
  MOCK_METHOD0(s6a_viface_authentication_info_req, void());
  MOCK_METHOD0(s6a_viface_update_location_req, void());
  MOCK_METHOD0(s6a_viface_purge_ue, void());
};

class MockSpgwHandler {
 public:
  MOCK_METHOD0(sgw_handle_s11_create_session_request, void());
  MOCK_METHOD0(sgw_handle_delete_session_request, void());
  MOCK_METHOD0(sgw_handle_modify_bearer_request, void());
  MOCK_METHOD0(sgw_handle_release_access_bearers_request, void());
  MOCK_METHOD0(sgw_handle_nw_initiated_actv_bearer_rsp, void());
};

class MockService303Handler {
 public:
  MOCK_METHOD0(service303_set_application_health, void());
};

class MockS8Handler {
 public:
  MOCK_METHOD1(
      sgw_s8_handle_create_bearer_request,
      bool(s8_create_bearer_request_t cb_req));
  MOCK_METHOD1(
      sgw_s8_handle_delete_bearer_request,
      bool(s8_delete_bearer_request_t db_req));
  MOCK_METHOD1(
      sgw_s8_handle_create_session_response,
      bool(s8_create_session_response_t cs_rsp));
  MOCK_METHOD1(
      sgw_s8_handle_delete_session_response,
      bool(s8_delete_session_response_t ds_rsp));
};

void start_mock_ha_task();
void start_mock_s1ap_task(std::shared_ptr<MockS1apHandler>);
void start_mock_sctp_task(std::shared_ptr<MockSctpHandler>);
void start_mock_mme_app_task(std::shared_ptr<MockMmeAppHandler>);
void start_mock_s6a_task(std::shared_ptr<MockS6aHandler>);
void start_mock_s11_task();
void start_mock_service303_task(std::shared_ptr<MockService303Handler>);
void start_mock_sgs_task();
void start_mock_sgw_s8_task(std::shared_ptr<MockS8Handler>);
void start_mock_sms_orc8r_task();
void start_mock_spgw_task(std::shared_ptr<MockSpgwHandler>);
void start_mock_grpc_task();
