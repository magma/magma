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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/sgw_s8_state.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw_s8/sgw_s8_s11_handlers.hpp"

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg);
static void sgw_s8_exit(void);
task_zmq_ctx_t sgw_s8_task_zmq_ctx;

static void* sgw_s8_thread(void* args) {
  itti_mark_task_ready(TASK_SGW_S8);
  const task_id_t peer_task_id[] = {TASK_MME_APP};
  init_task_context(TASK_SGW_S8, peer_task_id, 1, handle_message,
                    &sgw_s8_task_zmq_ctx);

  zloop_start(sgw_s8_task_zmq_ctx.event_loop);
  AssertFatal(0,
              "Asserting as sgw_s8_thread should not be exiting on its own! "
              "This is likely due to a timer handler function returning -1 "
              "(RETURNerror) on one of the conditions.");
  return NULL;
}

status_code_e sgw_s8_init(sgw_config_t* sgw_config_p) {
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

  imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
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
          sgw_state,
          &received_message_p->ittiMsg.s11_release_access_bearers_request,
          imsi64);
    } break;
    case S8_CREATE_BEARER_REQ: {
      gtpv2c_cause_value_t cause_value = REQUEST_REJECTED;
      s8_create_bearer_request_t* cb_req =
          &received_message_p->ittiMsg.s8_create_bearer_req;
      imsi64_t imsi64 =
          sgw_s8_handle_create_bearer_request(sgw_state, cb_req, &cause_value);
      Imsi_t imsi = {0};
      if (imsi64 == INVALID_IMSI64) {
        sgw_s8_send_failed_create_bearer_response(
            sgw_state, cb_req->sequence_number, cb_req->pgw_cp_address,
            cause_value, imsi, cb_req->bearer_context[0].pgw_s8_up.teid);
      }
    } break;
    case S11_NW_INITIATED_ACTIVATE_BEARER_RESP: {
      sgw_s8_handle_s11_create_bearer_response(
          sgw_state, &received_message_p->ittiMsg.s11_nw_init_actv_bearer_rsp,
          imsi64);
    } break;
    case S8_DELETE_BEARER_REQ: {
      sgw_s8_handle_delete_bearer_request(
          sgw_state, &received_message_p->ittiMsg.s8_delete_bearer_req);
    } break;
    case S11_NW_INITIATED_DEACTIVATE_BEARER_RESP: {
      sgw_s8_handle_s11_delete_bearer_response(
          sgw_state, &received_message_p->ittiMsg.s11_nw_init_deactv_bearer_rsp,
          imsi64);
    } break;

    default: {
      OAILOG_DEBUG(LOG_SGW_S8, "Unknown message ID %d: %s\n",
                   ITTI_MSG_ID(received_message_p),
                   ITTI_MSG_NAME(received_message_p));
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
