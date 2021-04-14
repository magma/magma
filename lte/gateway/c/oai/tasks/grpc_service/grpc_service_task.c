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

#define grpc_service
#define grpc_service_TASK_C

#include <stddef.h>

#include "assertions.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "log.h"
#include "mme_default_values.h"
#include "grpc_service.h"

static void grpc_service_exit(void);

static grpc_service_data_t* grpc_service_config;
task_zmq_ctx_t grpc_service_task_zmq_ctx;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE:
      free(received_message_p);
      grpc_service_exit();
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

static void* grpc_service_thread(__attribute__((unused)) void* args) {
  itti_mark_task_ready(TASK_GRPC_SERVICE);
  init_task_context(
      TASK_GRPC_SERVICE,
      (task_id_t[]){TASK_SPGW_APP, TASK_HA, TASK_AMF_APP, TASK_SGW_S8}, 4,
      handle_message, &grpc_service_task_zmq_ctx);

  start_grpc_service(grpc_service_config->server_address);
  zloop_start(grpc_service_task_zmq_ctx.event_loop);
  grpc_service_exit();
  return NULL;
}

int grpc_service_init(void) {
  OAILOG_DEBUG(LOG_UTIL, "Initializing grpc_service task interface\n");
  grpc_service_config                 = calloc(1, sizeof(grpc_service_data_t));
  grpc_service_config->server_address = bfromcstr(GRPCSERVICES_SERVER_ADDRESS);

  if (itti_create_task(TASK_GRPC_SERVICE, &grpc_service_thread, NULL) < 0) {
    OAILOG_ALERT(LOG_UTIL, "Initializing grpc_service: ERROR\n");
    return RETURNerror;
  }
  return RETURNok;
}

static void grpc_service_exit(void) {
  bdestroy_wrapper(&grpc_service_config->server_address);
  free_wrapper((void**) &grpc_service_config);
  stop_grpc_service();
  destroy_task_context(&grpc_service_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_GRPC_SERVICE terminated\n");
  pthread_exit(NULL);
}
