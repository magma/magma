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
#include "dynamic_memory_check.h"
#define CHECK_PROTOTYPE_ONLY
#include "intertask_interface_init.h"
#undef CHECK_PROTOTYPE_ONLY
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_free_defined_msg.h"
}

const task_info_t tasks_info[] = {
    {THREAD_NULL, "TASK_UNKNOWN", "ipc://IPC_TASK_UNKNOWN"},
#define TASK_DEF(tHREADiD)                                                     \
  {THREAD_##tHREADiD, #tHREADiD, "ipc://IPC_" #tHREADiD},
#include <tasks_def.h>
#undef TASK_DEF
};

/* Map message id to message information */
const message_info_t messages_info[] = {
#define MESSAGE_DEF(iD, sTRUCT, fIELDnAME) {iD, sizeof(sTRUCT), #iD},
#include <messages_def.h>
#undef MESSAGE_DEF
};

class MockS1apHandler {
 public:
  MOCK_METHOD0(s1ap_generate_downlink_nas_transport, void());
  MOCK_METHOD0(s1ap_handle_conn_est_cnf, void());
  MOCK_METHOD0(s1ap_handle_ue_context_release_command, void());
};

class MockMmeAppHandler {
 public:
  MOCK_METHOD0(mme_app_handle_initial_ue_message, void());
  MOCK_METHOD0(mme_app_handle_s1ap_ue_context_release_req, void());
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
};

class MockService303Handler {
 public:
  MOCK_METHOD0(service303_set_application_health, void());
};

void start_mock_ha_task();
void start_mock_s1ap_task(std::shared_ptr<MockS1apHandler>);
void start_mock_sctp_task(std::shared_ptr<MockSctpHandler>);
void start_mock_mme_app_task(std::shared_ptr<MockMmeAppHandler>);
void start_mock_s6a_task(std::shared_ptr<MockS6aHandler>);
void start_mock_s11_task();
void start_mock_service303_task(std::shared_ptr<MockService303Handler>);
void start_mock_sgs_task();
void start_mock_sgw_s8_task();
void start_mock_sms_orc8r_task();
void start_mock_spgw_task(std::shared_ptr<MockSpgwHandler>);
