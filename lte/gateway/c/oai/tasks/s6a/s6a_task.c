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

/*! \file s6a_task.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <stdint.h>
#include <gnutls/gnutls.h>
#include <stdarg.h>
#include <string.h>

#include "bstrlib.h"
#include "3gpp_23.003.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "s6a_messages_types.h"
#include "service303.h"
#include "timer_messages_types.h"
#if HAVE_CONFIG_H
#include "config.h"
#endif
#include <freeDiameter/freeDiameter-host.h>
#include <freeDiameter/libfdcore.h>

#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "itti_free_defined_msg.h"
#include "common_defs.h"
#include "s6a_defs.h"
#include "s6a_messages.h"
#include "mme_config.h"
#include "timer.h"
#include "s6a_client_api.h"

#define S6A_PEER_CONNECT_TIMEOUT_MICRO_SEC (0)
#define S6A_PEER_CONNECT_TIMEOUT_SEC (1)

static void fd_gnutls_debug(int level, const char *str);
static void s6a_exit(void);

static int gnutls_log_level = 9;
static long timer_id = 0;
struct session_handler *ts_sess_hdl;
s6a_fd_cnf_t s6a_fd_cnf;
task_zmq_ctx_t s6a_task_zmq_ctx;

//------------------------------------------------------------------------------
static void fd_gnutls_debug(int loglevel, const char *str)
{
  OAILOG_EXTERNAL(loglevel, LOG_S6A, "[GTLS] %s", str);
}

//------------------------------------------------------------------------------
// callback for freeDiameter logs
static void oai_fd_logger(int loglevel, const char *format, va_list args)
{
#define FD_LOG_MAX_MESSAGE_LENGTH 8192
  char buffer[FD_LOG_MAX_MESSAGE_LENGTH];
  int rv = 0;

  rv = vsnprintf(buffer, sizeof(buffer), format, args);
  if ((0 > rv) || ((FD_LOG_MAX_MESSAGE_LENGTH) < rv)) {
    return;
  }
  OAILOG_EXTERNAL(OAILOG_LEVEL_TRACE - loglevel, LOG_S6A, "%s\n", buffer);
}

static int handle_message (zloop_t* loop, zsock_t* reader, void* arg)
{
  zframe_t* msg_frame = zframe_recv(reader);
  assert(msg_frame);
  MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);
  int rc = RETURNerror;

  switch (ITTI_MSG_ID(received_message_p)) {
    case MESSAGE_TEST: {
      OAI_FPRINTF_INFO("TASK_S6A received MESSAGE_TEST\n");
    } break;
    case S6A_UPDATE_LOCATION_REQ: {
#if S6A_OVER_GRPC
      rc = s6a_update_location_req(
        &received_message_p->ittiMsg.s6a_update_location_req);
#else
      rc = s6a_generate_update_location(
        &received_message_p->ittiMsg.s6a_update_location_req);
#endif
      if (rc) {
        OAILOG_DEBUG(
          LOG_S6A,
          "Sending s6a ULR for imsi=%s\n",
          received_message_p->ittiMsg.s6a_update_location_req.imsi);
      } else {
        OAILOG_ERROR(
          LOG_S6A,
          "Failure in sending s6a ULR for imsi=%s\n",
          received_message_p->ittiMsg.s6a_update_location_req.imsi);
      }
    } break;
    case S6A_AUTH_INFO_REQ: {
#if S6A_OVER_GRPC
      rc = s6a_authentication_info_req(
        &received_message_p->ittiMsg.s6a_auth_info_req);
#else
      rc = s6a_generate_authentication_info_req(
        &received_message_p->ittiMsg.s6a_auth_info_req);
#endif
      if (rc) {
        OAILOG_DEBUG(
          LOG_S6A,
          "Sending s6a AIR for imsi=%s\n",
          received_message_p->ittiMsg.s6a_auth_info_req.imsi);
      } else {
        OAILOG_ERROR(
          LOG_S6A,
          "Failure in sending s6a AIR for imsi=%s\n",
          received_message_p->ittiMsg.s6a_auth_info_req.imsi);
      }
    } break;
    case TIMER_HAS_EXPIRED: {
      if (!timer_exists(
            received_message_p->ittiMsg.timer_has_expired.timer_id)) {
        break;
      }
      /*
        * Trying to connect to peers
        */
      timer_id = 0;
      if (s6a_fd_new_peer() != RETURNok) {
        /*
          * On failure, reschedule timer.
          * * Preferred over TIMER_PERIODIC because if s6a_fd_new_peer takes
          * * longer to return than the period, the timer will schedule while
          * * the previous one is active, causing a seg fault.
          */
        increment_counter(
          "s6a_subscriberdb_connection_failure", 1, NO_LABELS);
        OAILOG_ERROR(
          LOG_S6A,
          "s6a_fd_new_peer has failed (%s:%d)\n",
          __FILE__,
          __LINE__);
        timer_setup(
          S6A_PEER_CONNECT_TIMEOUT_SEC,
          S6A_PEER_CONNECT_TIMEOUT_MICRO_SEC,
          TASK_S6A,
          INSTANCE_DEFAULT,
          TIMER_ONE_SHOT,
          NULL,
          0,
          &timer_id);
      }
      timer_handle_expired(
        received_message_p->ittiMsg.timer_has_expired.timer_id);
    } break;
    case S6A_CANCEL_LOCATION_ANS: {
      s6a_send_cancel_location_ans(
        &received_message_p->ittiMsg.s6a_cancel_location_ans);
    } break;
    case S6A_PURGE_UE_REQ: {
#if S6A_OVER_GRPC
      uint8_t imsi_length;
      imsi_length = received_message_p->ittiMsg.s6a_purge_ue_req.imsi_length;
      if (imsi_length > IMSI_BCD_DIGITS_MAX) {
        OAILOG_ERROR(
          LOG_S6A, "imsi length exceeds IMSI_BCD_DIGITS_MAX length \n");
      }
      received_message_p->ittiMsg.s6a_purge_ue_req.imsi[imsi_length] = '\0';
      rc = s6a_purge_ue(received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
#else
      rc = s6a_generate_purge_ue_req(
        &received_message_p->ittiMsg.s6a_purge_ue_req);
#endif
      if (rc) {
        OAILOG_DEBUG(
          LOG_S6A,
          "Sending s6a pur for imsi=%s\n",
          received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
      } else {
        OAILOG_ERROR(
          LOG_S6A,
          "Failure in sending s6a pur for imsi=%s\n",
          received_message_p->ittiMsg.s6a_purge_ue_req.imsi);
      }
    } break;
    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      zframe_destroy(&msg_frame);
      s6a_exit();
    } break;
    default: {
      OAILOG_DEBUG(
        LOG_S6A,
        "Unkwnon message ID %d: %s\n",
        ITTI_MSG_ID(received_message_p),
        ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  zframe_destroy(&msg_frame);
  return 0;
}

//------------------------------------------------------------------------------
static void* s6a_thread(__attribute__((unused)) void* args)
{
  itti_mark_task_ready(TASK_S6A);
  init_task_context(
      TASK_S6A,
      (task_id_t []) {TASK_MME_APP, TASK_S1AP},
      2,
      handle_message,
      &s6a_task_zmq_ctx);

#if S6A_OVER_GRPC
  send_activate_messages();
  OAILOG_DEBUG(LOG_S6A, "Initializing S6a interface over gRPC: DONE\n");
#endif

  zloop_start(s6a_task_zmq_ctx.event_loop);
  s6a_exit();
  return NULL;
}

//------------------------------------------------------------------------------
int s6a_init(const mme_config_t *mme_config_p)
{
  int ret;

  OAILOG_DEBUG(LOG_S6A, "Initializing S6a interface\n");

  if (itti_create_task(TASK_S6A, &s6a_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_S6A, "s6a create task\n");
    return RETURNerror;
  }

#if S6A_OVER_GRPC
  return RETURNok;
#endif
  memset(&s6a_fd_cnf, 0, sizeof(s6a_fd_cnf_t));

  /*
   * if (strcmp(fd_core_version(), free_wrapper_DIAMETER_MINIMUM_VERSION) ) {
   * S6A_ERROR("Freediameter version %s found, expecting %s\n", fd_core_version(),
   * free_wrapper_DIAMETER_MINIMUM_VERSION);
   * return RETURNerror;
   * } else {
   * S6A_DEBUG("Freediameter version %s\n", fd_core_version());
   * }
   */

  /*
   * Initializing freeDiameter logger
   */
  ret = fd_log_handler_register(oai_fd_logger);
  if (ret) {
    OAILOG_ERROR(
      LOG_S6A,
      "An error occurred during freeDiameter log handler registration: %d\n",
      ret);
    return ret;
  } else {
    OAILOG_DEBUG(LOG_S6A, "Initializing freeDiameter log handler done\n");
  }

  /*
   * Initializing freeDiameter core
   */
  OAILOG_DEBUG(LOG_S6A, "Initializing freeDiameter core...\n");
  ret = fd_core_initialize();
  if (ret) {
    OAILOG_ERROR(
      LOG_S6A,
      "An error occurred during freeDiameter core library initialization: %d\n",
      ret);
    return ret;
  } else {
    OAILOG_DEBUG(LOG_S6A, "Initializing freeDiameter core done\n");
  }

  OAILOG_DEBUG(LOG_S6A, "Default ext path: %s\n", DEFAULT_EXTENSIONS_PATH);

  ret = fd_core_parseconf(bdata(mme_config_p->s6a_config.conf_file));
  if (ret) {
    OAILOG_ERROR(
      LOG_S6A,
      "An error occurred during fd_core_parseconf file : %s.\n",
      bdata(mme_config_p->s6a_config.conf_file));
    return ret;
  } else {
    OAILOG_DEBUG(LOG_S6A, "fd_core_parseconf done\n");
  }

  /*
   * Set gnutls debug level ?
   */
  if (gnutls_log_level) {
    gnutls_global_set_log_function((gnutls_log_func) fd_gnutls_debug);
    gnutls_global_set_log_level(gnutls_log_level);
    OAILOG_DEBUG(
      LOG_S6A, "Enabled GNUTLS debug at level %d\n", gnutls_log_level);
  }

  /*
   * Starting freeDiameter core
   */
  ret = fd_core_start();
  if (ret) {
    OAILOG_ERROR(
      LOG_S6A, "An error occurred during freeDiameter core library start\n");
    return ret;
  } else {
    OAILOG_DEBUG(LOG_S6A, "fd_core_start done\n");
  }

  ret = fd_core_waitstartcomplete();
  if (ret) {
    OAILOG_ERROR(
      LOG_S6A, "An error occurred during freeDiameter core library start\n");
    return ret;
  } else {
    OAILOG_DEBUG(LOG_S6A, "fd_core_waitstartcomplete done\n");
  }

  ret = s6a_fd_init_dict_objs();
  if (ret) {
    OAILOG_ERROR(LOG_S6A, "An error occurred during s6a_fd_init_dict_objs.\n");
    return ret;
  } else {
    OAILOG_DEBUG(LOG_S6A, "s6a_fd_init_dict_objs done\n");
  }

  OAILOG_DEBUG(
    LOG_S6A,
    "Initializing S6a interface over free-diameter:"
    "DONE\n");

  /* Add timer here to send message to connect to peer */
  timer_setup(
    S6A_PEER_CONNECT_TIMEOUT_SEC,
    S6A_PEER_CONNECT_TIMEOUT_MICRO_SEC,
    TASK_S6A,
    INSTANCE_DEFAULT,
    TIMER_ONE_SHOT,
    NULL,
    0,
    &timer_id);

  return RETURNok;
}

//------------------------------------------------------------------------------
static void s6a_exit(void)
{
  destroy_task_context(&s6a_task_zmq_ctx);

  if (timer_id) {
    timer_remove(timer_id, NULL);
  }
#if (!S6A_OVER_GRPC)
  // Release all resources
  free_wrapper((void **) &fd_g_config->cnf_diamid);
  fd_g_config->cnf_diamid_len = 0;
  int rv = RETURNok;
  /* Initialize shutdown of the framework */
  rv = fd_core_shutdown();
  if (rv) {
    OAI_FPRINTF_ERR("An error occurred during fd_core_shutdown().\n");
  }

  /* Wait for the shutdown to be complete -- this should always be called after fd_core_shutdown */
  rv = fd_core_wait_shutdown_complete();
  if (rv) {
    OAI_FPRINTF_ERR(
      "An error occurred during fd_core_wait_shutdown_complete().\n");
  }
#endif

  OAI_FPRINTF_INFO("TASK_S6A terminated\n");
  pthread_exit(NULL);
}
