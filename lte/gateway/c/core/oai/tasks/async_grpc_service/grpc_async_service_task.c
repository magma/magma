/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#define grpc_async_service
#define grpc_async_service_TASK_C

#include "assertions.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "grpc_async_service_task.h"
#include "S6aAsyncGrpc.h"

static void grpc_async_service_exit(void);

task_zmq_ctx_t grpc_async_service_task_zmq_ctx;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE:
      free(received_message_p);
      grpc_async_service_exit();
      break;
    default:
      OAILOG_DEBUG(
          LOG_UTIL, "Unknown message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
      break;
  }
  free(received_message_p);
  return 0;
}

static void* grpc_async_service_thread(__attribute__((unused)) void* args) {
  task_zmq_ctx_t* grpc_async_service_task_zmq_ctx_p =
      &grpc_async_service_task_zmq_ctx;

  itti_mark_task_ready(TASK_ASYNC_GRPC_SERVICE);

  init_task_context(
      TASK_ASYNC_GRPC_SERVICE, (task_id_t[]){TASK_MME_APP, TASK_S6A}, 2,
      handle_message, grpc_async_service_task_zmq_ctx_p);
  init_async_grpc_server();
  zloop_start(grpc_async_service_task_zmq_ctx.event_loop);
  AssertFatal(
      0, "Asserting as grpc_service_thread should not be exiting on its own!");
  return NULL;
}

status_code_e grpc_async_service_init(void) {
  OAILOG_DEBUG(LOG_UTIL, "Initializing async_grpc_service task interface\n");

  if (itti_create_task(
          TASK_ASYNC_GRPC_SERVICE, &grpc_async_service_thread, NULL) < 0) {
    OAILOG_ALERT(LOG_UTIL, "Initializing async_grpc_service: ERROR\n");
    return RETURNerror;
  }
  return RETURNok;
}

static void grpc_async_service_exit(void) {
  stop_async_grpc_service();
  destroy_task_context(&grpc_async_service_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_ASYNC_GRPC_SERVICE terminated\n");
  pthread_exit(NULL);
}
