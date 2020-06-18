/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

static void* grpc_service_task(void* args)
{
  MessageDef* received_message_p = NULL;
  grpc_service_data_t* grpc_service_data = (grpc_service_data_t*) args;

  itti_mark_task_ready(TASK_GRPC_SERVICE);
  start_grpc_service(grpc_service_data->server_address);

  while (1) {
    /*
     * Trying to fetch a message from the message queue.
     * If the queue is empty, this function will block till a
     * message is sent to the task.
     */
    itti_receive_msg(TASK_GRPC_SERVICE, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case TERMINATE_MESSAGE:
        stop_grpc_service();
        bdestroy_wrapper(&grpc_service_data->server_address);
        free_wrapper((void**) &grpc_service_data);
        OAI_FPRINTF_INFO("TASK_GRPC_SERVICE terminated\n");
        itti_exit_task();
        break;
      default:
        OAILOG_DEBUG(
          LOG_UTIL,
          "Unknown message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
        break;
    }

    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

int grpc_service_init(void)
{
  OAILOG_DEBUG(LOG_UTIL, "Initializing grpc_service task interface\n");
  grpc_service_data_t* grpc_service_config =
    calloc(1, sizeof(grpc_service_config));
  grpc_service_config->server_address = bfromcstr(GRPCSERVICES_SERVER_ADDRESS);

  if (
    itti_create_task(
      TASK_GRPC_SERVICE, &grpc_service_task, grpc_service_config) < 0) {
    OAILOG_ALERT(LOG_UTIL, "Initializing grpc_service: ERROR\n");
    return RETURNerror;
  }
  return RETURNok;
}
