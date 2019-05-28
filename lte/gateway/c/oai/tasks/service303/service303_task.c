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
#define SERVICE303
#define SERVICE303_TASK_C

#include <stdio.h>

#include "log.h"
#include "intertask_interface.h"
#include "timer.h"
#include "common_defs.h"
#include "service303.h"
#include "bstrlib.h"
#include "intertask_interface_types.h"
#include "timer_messages_types.h"

static long service303_epc_stats_timer_id;

static void *service303_server_task(void *args)
{
  service303_data_t *service303_data = (service303_data_t *) args;
  MessageDef *received_message_p = NULL;

  start_service303_server(service303_data->name, service303_data->version);

  itti_mark_task_ready(TASK_SERVICE303_SERVER);

  while (1) {
    itti_receive_msg(TASK_SERVICE303_SERVER, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case TERMINATE_MESSAGE:
        stop_service303_server();
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

static void *service303_message_task(void *args)
{
  bstring pkg_name = bfromcstr(SERVICE303_MME_PACKAGE_NAME);
  service303_data_t *service303_data = (service303_data_t *) args;

  itti_mark_task_ready(TASK_SERVICE303);

  if (bstricmp(service303_data->name, pkg_name) == 0) {
    /* NOTE : Above check for MME package is added since SPGW does not support stats at present
     * TODO : Whenever SPGW implements stats,remove the above "if" check so that timer is started
     * in SPGW also and SPGW stats can also be read as part of timer expiry handling
     */

    /*
     * Check if this thread is started by MME service if so start a timer
     * to trigger reading the mme stats so that it cen be sent to server
     * for display
     * Start periodic timer
     */
    if (
      timer_setup(
        EPC_STATS_TIMER_VALUE,
        0,
        TASK_SERVICE303,

        TIMER_PERIODIC,
        NULL,
        0,
        &service303_epc_stats_timer_id) < 0) {
      OAILOG_ALERT(
        LOG_UTIL,
        " TASK SERVICE303_MESSAGE for EPC: Periodic Stat Timer Start: ERROR\n");
      service303_epc_stats_timer_id = 0;
    }
  }

  bdestroy(pkg_name);

  while (1) {
    MessageDef *received_message_p = NULL;
    /*
     * Trying to fetch a message from the message queue.
     * If the queue is empty, this function will block till a
     * message is sent to the task.
     */
    itti_receive_msg(TASK_SERVICE303, &received_message_p);

    switch (ITTI_MSG_ID(received_message_p)) {
      case TIMER_HAS_EXPIRED: {
        /*
       * Check statistic timer
       */
        if (!timer_exists(
              received_message_p->ittiMsg.timer_has_expired.timer_id)) {
          break;
        }
        if (
          received_message_p->ittiMsg.timer_has_expired.timer_id ==
          service303_epc_stats_timer_id) {
          service303_statistics_read();
        }
        timer_handle_expired(
          received_message_p->ittiMsg.timer_has_expired.timer_id);
        break;
      }
      case TERMINATE_MESSAGE: {
        timer_remove(service303_epc_stats_timer_id, NULL);
        itti_exit_task();
      } break;
      case APPLICATION_HEALTHY_MSG: {
        service303_set_application_health(APP_HEALTHY);
      } break;
      case APPLICATION_UNHEALTHY_MSG: {
        service303_set_application_health(APP_UNHEALTHY);
      } break;
      default: {
        OAILOG_DEBUG(
          LOG_UTIL,
          "Unkwnon message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }
    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }
  return NULL;
}

int service303_init(service303_data_t *service303_data)
{
  OAILOG_DEBUG(LOG_UTIL, "Initializing Service303 task interface\n");

  if (
    itti_create_task(
      TASK_SERVICE303_SERVER, &service303_server_task, service303_data) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_UTIL, "Initializing Service303 server: ERROR\n");
    return RETURNerror;
  }

  if (
    itti_create_task(
      TASK_SERVICE303, &service303_message_task, service303_data) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(
      LOG_UTIL, "Initializing Service303 message interface: ERROR\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_UTIL, "Initializing Service303 task interface: DONE\n");
  return RETURNok;
}
