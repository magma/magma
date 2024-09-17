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

#define SERVICE303
#define SERVICE303_TASK_C

#include <stdio.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/service303.hpp"

static void service303_server_exit(void);
static void service303_message_exit(void);

task_zmq_ctx_t service303_server_task_zmq_ctx;
task_zmq_ctx_t service303_message_task_zmq_ctx;
static long display_stats_timer_id;
static int handle_display_timer(zloop_t*, int, void*);
static void start_display_stats_timer(size_t);

static int handle_service303_server_message(zloop_t* loop, zsock_t* reader,
                                            void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE:
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      service303_server_exit();
      break;
    default: {
      OAILOG_DEBUG(LOG_UTIL, "Unknown message ID %d: %s\n",
                   ITTI_MSG_ID(received_message_p),
                   ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

static void* service303_server_thread(__attribute__((unused)) void* args) {
  service303_data_t* service303_data = (service303_data_t*)args;

  start_service303_server(service303_data->name, service303_data->version);

  itti_mark_task_ready(TASK_SERVICE303_SERVER);
  const task_id_t peer_task_ids[] = {};
  init_task_context(TASK_SERVICE303_SERVER, peer_task_ids, 0,
                    handle_service303_server_message,
                    &service303_server_task_zmq_ctx);

  zloop_start(service303_server_task_zmq_ctx.event_loop);
  AssertFatal(
      0,
      "Asserting as service303_server_thread should not be exiting on its "
      "own!");
  return NULL;
}

static int handle_service_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case APPLICATION_HEALTHY_MSG: {
      service303_set_application_health(APP_HEALTHY);
    } break;
    case APPLICATION_UNHEALTHY_MSG: {
      service303_set_application_health(APP_UNHEALTHY);
    } break;
    case APPLICATION_MME_APP_STATS_MSG: {
      service303_mme_app_statistics_read(
          &received_message_p->ittiMsg.application_mme_app_stats_msg);
    } break;
    case APPLICATION_S1AP_STATS_MSG: {
      service303_s1ap_statistics_read(
          &received_message_p->ittiMsg.application_s1ap_stats_msg);
    } break;
    case APPLICATION_AMF_APP_STATS_MSG: {
      service303_amf_app_statistics_read(
          &received_message_p->ittiMsg.application_amf_app_stats_msg);
    } break;
    case APPLICATION_NGAP_STATS_MSG: {
      service303_ngap_statistics_read(
          &received_message_p->ittiMsg.application_ngap_stats_msg);
    } break;
    case TERMINATE_MESSAGE:
      free(received_message_p);
      service303_message_exit();
      break;
    default: {
      OAILOG_DEBUG(LOG_UTIL, "Unknown message ID %d: %s\n",
                   ITTI_MSG_ID(received_message_p),
                   ITTI_MSG_NAME(received_message_p));
    } break;
  }

  free(received_message_p);
  return 0;
}

static void* service303_thread(void* args) {
  service303_data_t* service303_data = (service303_data_t*)args;
  itti_mark_task_ready(TASK_SERVICE303);
  const task_id_t peer_task_ids[] = {};
  init_task_context(TASK_SERVICE303, peer_task_ids, 0, handle_service_message,
                    &service303_message_task_zmq_ctx);
  start_display_stats_timer((size_t)service303_data->stats_display_timer_sec);
  zloop_start(service303_message_task_zmq_ctx.event_loop);
  AssertFatal(
      0, "Asserting as service303_thread should not be exiting on its own!");
  return NULL;
}

status_code_e service303_init(service303_data_t* service303_data) {
  OAILOG_DEBUG(LOG_UTIL, "Initializing Service303 task interface\n");

  if (itti_create_task(TASK_SERVICE303_SERVER, &service303_server_thread,
                       service303_data) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_UTIL, "Initializing Service303 server: ERROR\n");
    return RETURNerror;
  }

  if (itti_create_task(TASK_SERVICE303, &service303_thread, service303_data) <
      0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_UTIL,
                 "Initializing Service303 message interface: ERROR\n");
    return RETURNerror;
  }

  OAILOG_DEBUG(LOG_UTIL, "Initializing Service303 task interface: DONE\n");
  return RETURNok;
}

static void service303_server_exit(void) {
  destroy_task_context(&service303_server_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_SERVICE303_SERVER terminated\n");
  pthread_exit(NULL);
}

static void service303_message_exit(void) {
  stop_timer(&service303_message_task_zmq_ctx, display_stats_timer_id);
  destroy_task_context(&service303_message_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_SERVICE303 terminated\n");
  pthread_exit(NULL);
}

static int handle_display_timer(zloop_t* loop, int id, void* arg) {
  service303_statistics_display();
  service303_statistics_display_5G();
  return 0;
}

static void start_display_stats_timer(size_t stats_display_timer_sec) {
  display_stats_timer_id = start_timer(
      &service303_message_task_zmq_ctx, 1000 * stats_display_timer_sec,
      TIMER_REPEAT_FOREVER, handle_display_timer, NULL);
}
