/**
 * Copyright 2021 The Magma Authors.
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
#include <thread>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
}

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/mme_init.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_state.hpp"
#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_test_utils.h"
#include "lte/gateway/c/core/oai/test/s1ap_task/mock_s1ap_op.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.hpp"

extern bool hss_associated;

namespace magma {
namespace lte {

task_zmq_ctx_t task_zmq_ctx_main_s1ap_with_injected_states;

// mocking the message handler for the ITTI
static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    // TODO: adding the message handler for different types of message
    default: {
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

class S1apMmeHandlersWithInjectedStatesTest : public ::testing::Test {
  virtual void SetUp() {
    mme_app_handler = std::make_shared<MockMmeAppHandler>();
    sctp_handler = std::make_shared<MockSctpHandler>();

    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    // initialize mme config
    mme_config_init(&mme_config);
    create_partial_lists(&mme_config);
    mme_config.use_stateless = false;
    hss_associated = true;

    task_id_t task_id_list[4] = {TASK_MME_APP, TASK_S1AP, TASK_SCTP,
                                 TASK_SERVICE303};
    init_task_context(TASK_MAIN, task_id_list, 4, handle_message,
                      &task_zmq_ctx_main_s1ap_with_injected_states);

    std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
    std::thread task_sctp(start_mock_sctp_task, sctp_handler);
    task_mme_app.detach();
    task_sctp.detach();

    s1ap_mme_init(&mme_config);

    // add injection of state loaded in S1AP state manager
    std::string magma_root = std::getenv("MAGMA_ROOT");
    std::string state_data_path =
        magma_root + "/" + DEFAULT_S1AP_STATE_DATA_PATH;
    std::string data_folder_path =
        magma_root + "/" + DEFAULT_S1AP_CONTEXT_DATA_PATH;
    std::string data_list_path =
        magma_root + "/" + DEFAULT_S1AP_CONTEXT_DATA_PATH + "data_list.txt";
    assoc_id = 37;
    stream_id = 1;
    number_attached_ue = 2;

    mock_read_s1ap_state_db(state_data_path);
    name_of_ue_samples =
        load_file_into_vector_of_line_content(data_folder_path, data_list_path);
    mock_read_s1ap_ue_state_db(name_of_ue_samples);

    state = S1apStateManager::getInstance().get_state(false);
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
  }

  virtual void TearDown() {
    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(2000));

    send_terminate_message_fatal(&task_zmq_ctx_main_s1ap_with_injected_states);
    destroy_task_context(&task_zmq_ctx_main_s1ap_with_injected_states);
    itti_free_desc_threads();

    free_mme_config(&mme_config);

    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(1500));
  }

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::shared_ptr<MockSctpHandler> sctp_handler;
  oai::S1apState* state;
  sctp_assoc_id_t assoc_id;
  sctp_stream_id_t stream_id;
  std::vector<std::string> name_of_ue_samples;
  int number_attached_ue;
};

TEST_F(S1apMmeHandlersWithInjectedStatesTest, GenerateUEContextReleaseCommand) {
  oai::UeDescription ue_ref_p;
  ue_ref_p.Clear();
  ue_ref_p.set_enb_ue_s1ap_id(1);
  ue_ref_p.set_mme_ue_s1ap_id(99);
  ue_ref_p.set_sctp_assoc_id(assoc_id);
  ue_ref_p.set_comp_s1ap_id(S1AP_GENERATE_COMP_S1AP_ID(assoc_id, 1));
  ue_ref_p.mutable_s1ap_ue_context_rel_timer()->set_id(-1);
  ue_ref_p.mutable_s1ap_ue_context_rel_timer()->set_msec(1000);

  S1ap_S1AP_PDU_t pdu_s1;
  memset(&pdu_s1, 0, sizeof(pdu_s1));
  ASSERT_EQ(RETURNok, generate_s1_setup_request_pdu(&pdu_s1));

  // State validation
  ASSERT_TRUE(
      is_enb_state_valid(state, assoc_id, oai::S1AP_READY, number_attached_ue));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));

  // Invalid S1 Cause returns error
  ASSERT_EQ(RETURNerror, s1ap_mme_generate_ue_context_release_command(
                             state, &ue_ref_p, S1AP_IMPLICIT_CONTEXT_RELEASE,
                             INVALID_IMSI64, assoc_id, stream_id, 99, 1));
  // Valid S1 Causes passess successfully
  ASSERT_EQ(RETURNok, s1ap_mme_generate_ue_context_release_command(
                          state, &ue_ref_p, S1AP_INITIAL_CONTEXT_SETUP_FAILED,
                          INVALID_IMSI64, assoc_id, stream_id, 99, 1));

  EXPECT_NE(ue_ref_p.s1ap_ue_context_rel_timer().id(), S1AP_TIMER_INACTIVE_ID);

  // State validation
  ASSERT_TRUE(
      is_enb_state_valid(state, assoc_id, oai::S1AP_READY, number_attached_ue));

  // Freeing pdu and payload data
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, &pdu_s1);
}

TEST_F(S1apMmeHandlersWithInjectedStatesTest, HandleS1apPathSwitchRequest) {
  ASSERT_EQ(task_zmq_ctx_main_s1ap_with_injected_states.ready, true);

  // State validation
  ASSERT_TRUE(
      is_enb_state_valid(state, assoc_id, oai::S1AP_READY, number_attached_ue));
  ASSERT_TRUE(is_num_enbs_valid(state, 1));
  ASSERT_EQ(state->mmeid2associd_size(), number_attached_ue);

  // Send S1AP_PATH_SWITCH_REQUEST_ACK mimicing MME_APP
  ASSERT_EQ(send_s1ap_path_switch_req(assoc_id, 1, 7), RETURNok);

  // verify number of ues after sending S1AP_PATH_SWITCH_REQUEST_ACK
  ASSERT_TRUE(
      is_enb_state_valid(state, assoc_id, oai::S1AP_READY, number_attached_ue));
}

}  // namespace lte
}  // namespace magma
