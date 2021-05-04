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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file s6a_fd_iface.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
*/

#ifdef __cplusplus
extern "C" {
#endif

#include "common_defs.h"
#include "s6a_messages.h"
#include "s6a_messages_types.h"
#include "intertask_interface.h"
#include "mme_config.h"
#include "dynamic_memory_check.h"

#ifdef __cplusplus
}
#endif

#include "s6a_fd_iface.h"

#include <iostream>
#include <exception>

using namespace std;

extern task_zmq_ctx_t s6a_task_zmq_ctx;

static int gnutls_log_level = 9;
static long timer_id        = 0;

static void fd_gnutls_debug(int level, const char* str);
static void oai_fd_logger(int loglevel, const char* format, va_list args);

#define S6A_PEER_CONNECT_TIMEOUT_MSEC (1000)

// TODO Mohit should be member of S6aFdIface (hide this global var).
s6a_fd_cnf_t s6a_fd_cnf;

//------------------------------------------------------------------------------
static void fd_gnutls_debug(int loglevel, const char* str) {
  OAILOG_EXTERNAL(loglevel, LOG_S6A, "[GTLS] %s", str);
}

//------------------------------------------------------------------------------
// callback for freeDiameter logs
static void oai_fd_logger(int loglevel, const char* format, va_list args) {
#define FD_LOG_MAX_MESSAGE_LENGTH 8192
  char buffer[FD_LOG_MAX_MESSAGE_LENGTH];
  int rv = 0;

  rv = vsnprintf(buffer, sizeof(buffer), format, args);
  if ((0 > rv) || ((FD_LOG_MAX_MESSAGE_LENGTH) < rv)) {
    return;
  }
  OAILOG_EXTERNAL(OAILOG_LEVEL_TRACE - loglevel, LOG_S6A, "%s\n", buffer);
}

static int handle_timer(zloop_t* loop, int id, void* arg) {
  /*
   * Trying to connect to peers
   */
  if (s6a_fd_new_peer() != RETURNok) {
    /*
     * On failure, reschedule timer.
     * * Preferred over TIMER_REPEAT_FOREVER because if s6a_fd_new_peer takes
     * * longer to return than the period, the timer will schedule while
     * * the previous one is active, causing a seg fault.
     */
    increment_counter("s6a_subscriberdb_connection_failure", 1, NO_LABELS);
    OAILOG_ERROR(
        LOG_S6A, "s6a_fd_new_peer has failed (%s:%d)\n", __FILE__, __LINE__);
    timer_id = start_timer(
        &s6a_task_zmq_ctx, S6A_PEER_CONNECT_TIMEOUT_MSEC, TIMER_REPEAT_ONCE,
        handle_timer, NULL);
  }
  return 0;
}

//------------------------------------------------------------------------------
S6aFdIface::S6aFdIface(const s6a_config_t* const config) {
  int ret = RETURNok;
  memset(&s6a_fd_cnf, 0, sizeof(s6a_fd_cnf_t));

  // TODO check for a preprocessor difinition in freeDiameter repo hosted on OAI
  // Github

  /*
   * Initializing freeDiameter logger
   */
  ret = fd_log_handler_register(oai_fd_logger);
  if (ret) {
    OAILOG_ERROR(
        LOG_S6A,
        "An error occurred during freeDiameter log handler registration: %d\n",
        ret);
    std::runtime_error(
        "An error occurred during freeDiameter log handler "
        "registration");
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
        "An error occurred during freeDiameter core library initialization: "
        "%d\n",
        ret);
    std::runtime_error(
        "An error occurred during freeDiameter core library "
        "initialization");
  } else {
    OAILOG_DEBUG(LOG_S6A, "Initializing freeDiameter core done\n");
  }

  OAILOG_DEBUG(LOG_S6A, "Default ext path: %s\n", DEFAULT_EXTENSIONS_PATH);

  ret = fd_core_parseconf(bdata(config->conf_file));
  if (ret) {
    OAILOG_ERROR(
        LOG_S6A, "An error occurred during fd_core_parseconf file : %s.\n",
        bdata(config->conf_file));
    std::runtime_error("An error occurred during fd_core_parseconf file");
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
    std::runtime_error(
        "An error occurred during freeDiameter core library "
        "start");
  } else {
    OAILOG_DEBUG(LOG_S6A, "fd_core_start done\n");
  }

  ret = fd_core_waitstartcomplete();
  if (ret) {
    OAILOG_ERROR(
        LOG_S6A, "An error occurred during freeDiameter core library start\n");
    std::runtime_error(
        "An error occurred during freeDiameter core library "
        "start\n");
  } else {
    OAILOG_DEBUG(LOG_S6A, "fd_core_waitstartcomplete done\n");
  }

  ret = s6a_fd_init_dict_objs();
  if (ret) {
    OAILOG_ERROR(LOG_S6A, "An error occurred during s6a_fd_init_dict_objs.\n");
    std::runtime_error("An error occurred during s6a_fd_init_dict_obj\n");
  } else {
    OAILOG_DEBUG(LOG_S6A, "s6a_fd_init_dict_objs done\n");
  }

  OAILOG_DEBUG(
      LOG_S6A,
      "Initializing S6a interface over free-diameter:"
      "DONE\n");

  /* Add timer here to connect to peer */
  timer_id = start_timer(
      &s6a_task_zmq_ctx, S6A_PEER_CONNECT_TIMEOUT_MSEC, TIMER_REPEAT_ONCE,
      handle_timer, NULL);
}

//------------------------------------------------------------------------------
bool S6aFdIface::update_location_req(s6a_update_location_req_t* ulr_p) {
  if (s6a_generate_update_location(ulr_p))
    return false;
  else
    return true;
}
//------------------------------------------------------------------------------
bool S6aFdIface::authentication_info_req(s6a_auth_info_req_t* air_p) {
  if (s6a_generate_authentication_info_req(air_p))
    return false;
  else
    return true;
}
//------------------------------------------------------------------------------
bool S6aFdIface::send_cancel_location_ans(s6a_cancel_location_ans_t* cla_pP) {
  if (s6a_send_cancel_location_ans(cla_pP))
    return false;
  else
    return true;
}
//------------------------------------------------------------------------------
bool S6aFdIface::purge_ue(const char* imsi) {
  if (s6a_generate_purge_ue_req(imsi))
    return false;
  else
    return true;
}
//------------------------------------------------------------------------------
S6aFdIface::~S6aFdIface() {
  stop_timer(&s6a_task_zmq_ctx, timer_id);
  // Release all resources
  free_wrapper((void**) &fd_g_config->cnf_diamid);
  fd_g_config->cnf_diamid_len = 0;
  int rv                      = RETURNok;
  /* Initialize shutdown of the framework */
  rv = fd_core_shutdown();
  if (rv) {
    OAI_FPRINTF_ERR("An error occurred during fd_core_shutdown().\n");
  }

  /* Wait for the shutdown to be complete -- this should always be called after
   * fd_core_shutdown */
  rv = fd_core_wait_shutdown_complete();
  if (rv) {
    OAI_FPRINTF_ERR(
        "An error occurred during fd_core_wait_shutdown_complete().\n");
  }
}
