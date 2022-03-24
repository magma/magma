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

#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_handlers_test_fixture.h"

#include "lte/gateway/c/core/oai/test/s1ap_task/s1ap_mme_test_utils.h"
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_manager.h"
#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"

namespace magma {
namespace lte {

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

void S1apMmeHandlersTest::SetUp() {
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
                    &task_zmq_ctx_main_s1ap);

  std::thread task_mme_app(start_mock_mme_app_task, mme_app_handler);
  std::thread task_sctp(start_mock_sctp_task, sctp_handler);
  task_mme_app.detach();
  task_sctp.detach();

  s1ap_mme_init(&mme_config);

  // Setup new association for testing
  state = S1apStateManager::getInstance().get_state(false);
  assoc_id = 1;
  stream_id = 0;
  setup_new_association(state, assoc_id);
  std::this_thread::sleep_for(std::chrono::milliseconds(500));
}

void S1apMmeHandlersTest::TearDown() {
  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(500));

  send_terminate_message_fatal(&task_zmq_ctx_main_s1ap);
  destroy_task_context(&task_zmq_ctx_main_s1ap);
  itti_free_desc_threads();

  free_mme_config(&mme_config);

  // Sleep to ensure that messages are received and contexts are released
  std::this_thread::sleep_for(std::chrono::milliseconds(200));
}
}  // namespace lte
}  // namespace magma