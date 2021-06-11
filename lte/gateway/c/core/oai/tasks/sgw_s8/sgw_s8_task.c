/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#define SGW_S8
#define SGW_S8_TASK_C

#include <stdio.h>
#include "log.h"
#include "assertions.h"
#include "common_defs.h"
#include "itti_free_defined_msg.h"
#include "sgw_s8_defs.h"
#include "sgw_s8_s11_handlers.h"
#include "sgw_s8_state.h"

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg);
static void sgw_s8_exit(void);
task_zmq_ctx_t sgw_s8_task_zmq_ctx;

static void* sgw_s8_thread(void* args) {
  itti_mark_task_ready(TASK_SGW_S8);
  init_task_context(
      TASK_SGW_S8, (task_id_t[]){TASK_MME_APP}, 1, handle_message,
      &sgw_s8_task_zmq_ctx);

  zloop_start(sgw_s8_task_zmq_ctx.event_loop);
  sgw_s8_exit();
  return NULL;
}

int sgw_s8_init(sgw_config_t* sgw_config_p) {
  OAILOG_DEBUG(LOG_SGW_S8, "Initializing SGW-S8 interface\n");
  if (sgw_state_init(false, sgw_config_p) < 0) {
    OAILOG_CRITICAL(LOG_SGW_S8, "Error while initializing SGW_S8 state\n");
    return RETURNerror;
  }

  if (itti_create_task(TASK_SGW_S8, &sgw_s8_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_SGW_S8, "Failed to create sgw_s8 task\n");
    return RETURNerror;
  }
  OAILOG_DEBUG(LOG_SGW_S8, "Done initialization of SGW_S8 interface\n");
  return RETURNok;
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  imsi64_t imsi64        = itti_get_associated_imsi(received_message_p);
  sgw_state_t* sgw_state = get_sgw_state(false);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      sgw_s8_exit();
    } break;

    case S11_CREATE_SESSION_REQUEST: {
      sgw_s8_handle_s11_create_session_request(
          sgw_state, &received_message_p->ittiMsg.s11_create_session_request,
          imsi64);
    } break;
    case S8_CREATE_SESSION_RSP: {
      sgw_s8_handle_create_session_response(
          sgw_state, &received_message_p->ittiMsg.s8_create_session_rsp,
          imsi64);
    } break;
    case S11_MODIFY_BEARER_REQUEST: {
      sgw_s8_handle_modify_bearer_request(
          sgw_state, &received_message_p->ittiMsg.s11_modify_bearer_request,
          imsi64);
    } break;

    case S11_DELETE_SESSION_REQUEST: {
      sgw_s8_handle_s11_delete_session_request(
          sgw_state, &received_message_p->ittiMsg.s11_delete_session_request,
          imsi64);
    } break;

    case S8_DELETE_SESSION_RSP: {
      sgw_s8_handle_delete_session_response(
          sgw_state, &received_message_p->ittiMsg.s8_delete_session_rsp,
          imsi64);
    } break;
    case S11_RELEASE_ACCESS_BEARERS_REQUEST: {
      sgw_s8_handle_release_access_bearers_request(
          &received_message_p->ittiMsg.s11_release_access_bearers_request,
          imsi64);
    } break;

    default: {
      OAILOG_DEBUG(
          LOG_SGW_S8, "Unknown message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void sgw_s8_exit(void) {
  destroy_task_context(&sgw_s8_task_zmq_ctx);
  OAILOG_DEBUG(LOG_SGW_S8, "Finished cleaning up SGW_S8 task \n");
  OAI_FPRINTF_INFO("TASK_SGW_S8 terminated\n");
  pthread_exit(NULL);
}
