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
#include <thread>

#include "../mock_tasks/mock_tasks.h"
#include "spgw_test_util.h"
#include "spgw_state.h"

extern "C" {
#include "mme_config.h"
#include "spgw_config.h"
#include "sgw_defs.h"
}

extern bool hss_associated;

namespace magma {
namespace lte {

task_zmq_ctx_t task_zmq_ctx_main_spgw;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    default: {
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

class SPGWAppProcedureTest : public ::testing::Test {
  virtual void SetUp() {
    // setup mock MME app task
    mme_app_handler = std::make_shared<MockMmeAppHandler>();
    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    // initialize configs
    mme_config_init(&mme_config);
    spgw_config_init(&spgw_config);
    create_partial_lists(&mme_config);
    mme_config.use_stateless = false;
    hss_associated           = true;

    task_id_t task_id_list[2] = {TASK_MME_APP, TASK_SPGW_APP};
    init_task_context(
        TASK_MAIN, task_id_list, 2, handle_message, &task_zmq_ctx_main_spgw);

    std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
    task_mme_app.detach();

    std::cout << "Running setup" << std::endl;
    // initialize the SPGW task
    spgw_app_init(&spgw_config, mme_config.use_stateless);
    std::this_thread::sleep_for(std::chrono::milliseconds(5000));
  }

  virtual void TearDown() {
    send_terminate_message_fatal(&task_zmq_ctx_main_spgw);
    destroy_task_context(&task_zmq_ctx_main_spgw);
    itti_free_desc_threads();

    free_mme_config(&mme_config);
    free_spgw_config(&spgw_config);

    // Sleep to ensure that messages are received and contexts are released
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
  }

 protected:
  std::shared_ptr<MockMmeAppHandler> mme_app_handler;
  std::string test_imsi_str = "001010000000001";
  plmn_t test_plmn          = {
               .mcc_digit2 = 0,
               .mcc_digit1 = 0,
               .mnc_digit3 = 0x0f,
               .mcc_digit3 = 1,
               .mnc_digit2 = 1,
               .mnc_digit1 = 0};
  bearer_context_to_be_created_t sample_default_bearer_context = {
      .eps_bearer_id    = 5,
      .bearer_level_qos = {
          .pci = 1,
          .pl  = 15,
          .pvi = 0,
          .qci = 9,
          .gbr = {},
          .mbr = {.br_ul = 200000000, .br_dl = 100000000}}};
}

TEST_F(SPGWAppProcedureTest, TestIPAllocSuccess) {
  // send a create session req to SPGW
  spgw_state_t* spgw_state = get_spgw_state(false);
  send_create_session_request(
      test_imsi_str, 0, sample_default_bearer_context, test_plmn);
  ASSERT_EQ(0, 0);
  // Verify that a UE context exists in SPGW state after CSR is received

  // Verify that no IP address is allocated for this UE

  // send an IP alloc response to SPGW
  // check that IP address is allocated after this message is done

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(10000));
};

TEST_F(SPGWAppProcedureTest, TestIPAllocFailure) {
  // send a create session req to SPGW
  spgw_state_t* spgw_state = get_spgw_state(false);
  // send_create_session_request( test_imsi_str, 0,
  // sample_default_bearer_context, test_plmn);
  ASSERT_EQ(0, 0);
  // Verify that a UE context exists in SPGW state after CSR is received

  // Verify that no IP address is allocated for this UE

  // send an IP alloc response to SPGW with status as failure
  // check that IP address is not allocated after this message is done

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(10000));
}

TEST_F(SPGWAppProcedureTest, TestCreateSessionRequest) {
  // expect call to MME create session response

  // send a create session req to SPGW
  // send a IP alloc response to SPGW
  // check if IP address is allocated after this message is done
  // send s5 response to SPGW
  // check if ambr has been returned

  ASSERT_EQ(0, 0);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(END_OF_TEST_SLEEP_MS));
}

}  // namespace lte
}  // namespace magma
