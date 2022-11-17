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
#include "lte/gateway/c/core/oai/common/conversions.h"
#define CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_init.h"
#undef CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
}

const task_info_t tasks_info[] = {
    {THREAD_NULL, "TASK_UNKNOWN", "ipc://IPC_TASK_UNKNOWN"},
#define TASK_DEF(tHREADiD) \
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

    default: {
    } break;
  }
  itti_free_msg_content(received_message_p);
  free(received_message_p);
  // Add sleep to introduce delay in pulling the next message
  std::this_thread::sleep_for(std::chrono::milliseconds(1500));
  return 0;
}

void* task_thread(task_thread_args_t* args) {
  init_task_context(args->this_task, args->task_id_list, args->list_size,
                    handle_message, args->task_zmq_ctx);

  zloop_start(args->task_zmq_ctx->event_loop);

  return NULL;
}

class ITTIMessagePassingTest : public ::testing::Test {
  virtual void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    task_id_t task_id_list[4] = {TASK_TEST_1, TASK_TEST_2};
    init_task_context(TASK_MAIN, task_id_list, 1, NULL, &task_zmq_ctx_main);

    task_thread_args_t task1_thread_args = {};
    task1_thread_args.this_task = TASK_TEST_1;
    task1_thread_args.task_id_list[0] = TASK_TEST_2;
    task1_thread_args.list_size = 1;
    task1_thread_args.task_zmq_ctx = &task_zmq_ctx_test1;

    task_thread_args_t task2_thread_args = {};
    task2_thread_args.this_task = TASK_TEST_2;
    task2_thread_args.task_id_list[0] = TASK_TEST_1;
    task2_thread_args.list_size = 1;
    task2_thread_args.task_zmq_ctx = &task_zmq_ctx_test2;

    std::thread task1(task_thread, &task1_thread_args);
    std::thread task2(task_thread, &task2_thread_args);
    task1.detach();
    task2.detach();
    std::this_thread::sleep_for(std::chrono::seconds(2));
  }

  virtual void TearDown() {
    send_terminate_message_fatal(&task_zmq_ctx_main);

    // Sleep 100 msec to allow message to be received before
    // destroying zmq context
    std::this_thread::sleep_for(std::chrono::milliseconds(100));
    // Destroy zmq contexts
    destroy_task_context(&task_zmq_ctx_test1);
    destroy_task_context(&task_zmq_ctx_test2);
    destroy_task_context(&task_zmq_ctx_main);
    itti_free_desc_threads();
  }
};

TEST_F(ITTIMessagePassingTest, TestMessageLatency) {
  MessageDef* test_message_p;
  test_message_p = DEPRECATEDitti_alloc_new_message_fatal(
      task_zmq_ctx_test1.task_id, TEST_MESSAGE);
  send_msg_to_task(&task_zmq_ctx_test1, TASK_TEST_2, test_message_p);
  // Sleep 100 msec to allow message to be received on time
  std::this_thread::sleep_for(std::chrono::milliseconds(100));
  ASSERT_LE(msg_latency, 1000);

  test_message_p = DEPRECATEDitti_alloc_new_message_fatal(
      task_zmq_ctx_test1.task_id, TEST_MESSAGE);
  send_msg_to_task(&task_zmq_ctx_test1, TASK_TEST_2, test_message_p);
  // Sleep 2 seconds to allow message to be received and processed
  std::this_thread::sleep_for(std::chrono::seconds(2));
  ASSERT_GE(msg_latency, 1000000);
}

class ITTIApiTest : public ::testing::Test {
  virtual void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);
  }

  virtual void TearDown() { itti_free_desc_threads(); }
};

TEST_F(ITTIApiTest, TestInitialContextSetupRsp) {
  uint16_t erab_no_of_items = 3;
  uint16_t failed_erabs = 4;
  char ipv4string[16] = "192.168.101.234";
  MessageDef* message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_S1AP, MME_APP_INITIAL_CONTEXT_SETUP_RSP);
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).ue_id = 10;
  MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p).e_rab_setup_list.no_of_items =
      erab_no_of_items;
  for (int item = 0; item < erab_no_of_items; item++) {
    MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
        .e_rab_setup_list.item[item]
        .e_rab_id = (uint8_t)item;  // Arbitrary RAB ids
    MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
        .e_rab_setup_list.item[item]
        .gtp_teid = (uint32_t)item;
    MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p)
        .e_rab_setup_list.item[item]
        .transport_layer_address = blk2bstr(ipv4string, 15);
  }
  itti_mme_app_initial_context_setup_rsp_t* ics_rsp =
      &(MME_APP_INITIAL_CONTEXT_SETUP_RSP(message_p));
  for (int index = 0; index < failed_erabs; index++) {
    ics_rsp->e_rab_failed_to_setup_list.item[index].e_rab_id =
        (uint8_t)index;  // Arbitrary RAB ids
    ics_rsp->e_rab_failed_to_setup_list.item[index].cause.present =
        S1ap_Cause_PR_radioNetwork;
    ics_rsp->e_rab_failed_to_setup_list.item[index].cause.choice.radioNetwork =
        S1ap_CauseRadioNetwork_radio_resources_not_available;
  }

  // Now free the message
  itti_free_msg_content(message_p);
  free(message_p);
}

TEST_F(ITTIApiTest, TestHandoverRequest) {
  char arbitrary_src_tgt_container[20] = "This is arbitrary";
  MessageDef* message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_MME_APP, MME_APP_HANDOVER_REQUEST);
  itti_mme_app_handover_request_t* ho_request_p =
      &message_p->ittiMsg.mme_app_handover_request;

  // fill in arbitrary values
  ho_request_p->encryption_algorithm_capabilities = 1;
  ho_request_p->integrity_algorithm_capabilities = 2;
  ho_request_p->mme_ue_s1ap_id = 10;
  ho_request_p->target_sctp_assoc_id = 1;
  ho_request_p->target_enb_id = 2;
  ho_request_p->cause.present = S1ap_Cause_PR_radioNetwork;
  ho_request_p->cause.choice.radioNetwork =
      S1ap_CauseRadioNetwork_handover_desirable_for_radio_reason;
  ho_request_p->handover_type = S1ap_HandoverType_intralte;
  ho_request_p->src_tgt_container = blk2bstr(arbitrary_src_tgt_container, 10);
  ho_request_p->ue_ambr.br_unit = KBPS;
  ho_request_p->ue_ambr.br_ul = 1000;
  ho_request_p->ue_ambr.br_dl = 10000;
  ho_request_p->e_rab_list.no_of_items = 2;
  fteid_t s_gw_fteid_s1u = {1};

  for (int i = 0; i < ho_request_p->e_rab_list.no_of_items; ++i) {
    ho_request_p->e_rab_list.item[i].e_rab_id = 1;
    ho_request_p->e_rab_list.item[i].transport_layer_address =
        fteid_ip_address_to_bstring(&s_gw_fteid_s1u);
    ho_request_p->e_rab_list.item[i].gtp_teid = 1;
    ho_request_p->e_rab_list.item[i].e_rab_level_qos_parameters.qci = 9;
    ho_request_p->e_rab_list.item[i]
        .e_rab_level_qos_parameters.allocation_and_retention_priority
        .priority_level = 0;
    ho_request_p->e_rab_list.item[i]
        .e_rab_level_qos_parameters.allocation_and_retention_priority
        .pre_emption_capability =
        (pre_emption_capability_t)PRE_EMPTION_CAPABILITY_ENABLED;
    ho_request_p->e_rab_list.item[i]
        .e_rab_level_qos_parameters.allocation_and_retention_priority
        .pre_emption_vulnerability =
        (pre_emption_vulnerability_t)PRE_EMPTION_VULNERABILITY_DISABLED;
  }
  for (int i = 0; i < AUTH_NEXT_HOP_SIZE; ++i) {
    ho_request_p->nh[i] = 0x11;
  }
  ho_request_p->ncc = 2;
  itti_free_msg_content(message_p);
  free(message_p);
}
