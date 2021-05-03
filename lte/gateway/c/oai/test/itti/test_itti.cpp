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
#include <chrono>
#include <string.h>
#include <gtest/gtest.h>
#include <thread>

extern "C" {
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

task_zmq_ctx_t task_zmq_ctx_main, task_zmq_ctx_test1, task_zmq_ctx_test2;

typedef struct {
  task_id_t this_task;
  task_id_t task_id_list[3];
  int list_size;
  task_zmq_ctx_t* task_zmq_ctx;
} task_thread_args_t;

long msg_latency;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
    } break;

    case TEST_MESSAGE: {
      msg_latency = ITTI_MSG_LATENCY(received_message_p);
    } break;

    default: { } break; }
  itti_free_msg_content(received_message_p);
  free(received_message_p);
  // Add sleep to introduce delay in pulling the next message
  std::this_thread::sleep_for(std::chrono::milliseconds(1500));
  return 0;
}

void* task_thread(task_thread_args_t* args) {
  init_task_context(
      args->this_task, args->task_id_list, args->list_size, handle_message,
      args->task_zmq_ctx);

  free(args);

  zloop_start(args->task_zmq_ctx->event_loop);

  return NULL;
}

class ITTIApiTest : public ::testing::Test {
  virtual void SetUp() {
    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);

    task_id_t task_id_list[4] = {TASK_TEST_1, TASK_TEST_2};
    init_task_context(TASK_MAIN, task_id_list, 1, NULL, &task_zmq_ctx_main);

    task_thread_args_t* task1_thread_args =
        (task_thread_args_t*) calloc(1, sizeof(task_thread_args_t));
    task1_thread_args->this_task       = TASK_TEST_1;
    task1_thread_args->task_id_list[0] = TASK_TEST_2;
    task1_thread_args->list_size       = 1;
    task1_thread_args->task_zmq_ctx    = &task_zmq_ctx_test1;

    task_thread_args_t* task2_thread_args =
        (task_thread_args_t*) calloc(1, sizeof(task_thread_args_t));
    task2_thread_args->this_task       = TASK_TEST_2;
    task2_thread_args->task_id_list[0] = TASK_TEST_1;
    task2_thread_args->list_size       = 1;
    task2_thread_args->task_zmq_ctx    = &task_zmq_ctx_test2;

    std::thread task1(task_thread, task1_thread_args);
    std::thread task2(task_thread, task2_thread_args);
    task1.detach();
    task2.detach();
    std::this_thread::sleep_for(std::chrono::seconds(2));
  }

  virtual void TearDown() {
    send_terminate_message(&task_zmq_ctx_main);

    // Sleep 100 msec to allow message to be received before
    // destroying zmq context
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
    // Destroy zmq contexts
    destroy_task_context(&task_zmq_ctx_test1);
    destroy_task_context(&task_zmq_ctx_test2);
    destroy_task_context(&task_zmq_ctx_main);
  }
};

TEST_F(ITTIApiTest, TestMessageLatency) {
  MessageDef* test_message_p;
  test_message_p =
      itti_alloc_new_message(task_zmq_ctx_test1.task_id, TEST_MESSAGE);
  send_msg_to_task(&task_zmq_ctx_test1, TASK_TEST_2, test_message_p);
  // Sleep 100 msec to allow message to be received on time
  std::this_thread::sleep_for(std::chrono::milliseconds(100));
  ASSERT_LE(msg_latency, 1000);

  test_message_p =
      itti_alloc_new_message(task_zmq_ctx_test1.task_id, TEST_MESSAGE);
  send_msg_to_task(&task_zmq_ctx_test1, TASK_TEST_2, test_message_p);
  // Sleep 2 seconds to allow message to be received and processed
  std::this_thread::sleep_for(std::chrono::seconds(2));
  ASSERT_GE(msg_latency, 1000000);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
