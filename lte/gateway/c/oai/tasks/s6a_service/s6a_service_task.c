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
#define s6a_service
#define s6a_service_TASK_C

#include "log.h"
#include "intertask_interface.h"
#include "assertions.h"
#include "common_defs.h"
#include "s6a_service.h"

static void *s6a_service_server_task(void *args)
{
  MessageDef *received_message_p = NULL;
  s6a_service_data_t *s6a_service_data = (s6a_service_data_t *) args;

  itti_mark_task_ready(TASK_S6A_SERVICE_SERVER);
  start_s6a_service_server(s6a_service_data->server_address);

  while (1) {
    /*
     * Trying to fetch a message from the message queue.
     * If the queue is empty, this function will block till a
     * message is sent to the task.
     */
    itti_receive_msg(TASK_S6A_SERVICE_SERVER, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case TERMINATE_MESSAGE:
        stop_s6a_service_server();
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
  }

  return NULL;
}

int s6a_service_server_init(void)
{
  OAILOG_DEBUG(LOG_UTIL, "Initializing s6a_server task interface\n");
  s6a_service_data_t s6a_config;
  s6a_config.server_address = bfromcstr(S6ASERVICE_SERVER_ADDRESS);

  if (
    itti_create_task(
      TASK_S6A_SERVICE_SERVER, &s6a_service_server_task, &s6a_config) < 0) {
    OAILOG_ALERT(LOG_UTIL, "Initializing s6a_server: ERROR\n");
    return RETURNerror;
  }
  return RETURNok;
}

static void *s6a_service_message_task(void *args)
{
  itti_mark_task_ready(TASK_S6A_SERVICE);
  while (1) {
    MessageDef *received_message_p = NULL;
    /*
     * Trying to fetch a message from the message queue.
     * If the queue is empty, this function will block till a
     * message is sent to the task.
     */
    itti_receive_msg(TASK_S6A_SERVICE, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case APPLICATION_HEALTHY_MSG:
        CHECK_INIT_RETURN(s6a_service_server_init());
        break;
      case TERMINATE_MESSAGE: itti_exit_task(); break;
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

int s6a_service_init()
{
  if (itti_create_task(TASK_S6A_SERVICE, &s6a_service_message_task, NULL) < 0) {
    OAILOG_ALERT(LOG_UTIL, "Initializing s6a_service: ERROR\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_UTIL, "Initializing s6a_service task interface: DONE\n");
  return RETURNok;
}
