/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file log.c
   \brief Thread safe logging utility, log output can be redirected to stdout,
   file or remote host through TCP. \author  Lionel GAUTHIER \date 2015 \email:
   lionel.gauthier@eurecom.fr
*/

#include <errno.h>
#include <fcntl.h>
#include <inttypes.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <sys/time.h>
#include <unistd.h>
#include <stdarg.h>
#include <pthread.h>
#include <syslog.h>
#include <assert.h>
#include <netinet/in.h>
#include <signal.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <time.h>

#include "intertask_interface.h"
#include "log.h"
#include "shared_ts_log.h"
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "asn_system.h"
#include "hashtable.h"
#include "intertask_interface_types.h"
#include "itti_types.h"

#if HAVE_CONFIG_H
#include "config.h"
#endif

//-------------------------------
#define LOG_MAX_QUEUE_ELEMENTS 1024
#define LOG_MAX_PROTO_NAME_LENGTH 16
#define LOG_MESSAGE_MIN_ALLOC_SIZE 256

#define LOG_CONNECT_PERIOD_MSEC 2000
#define LOG_FLUSH_PERIOD_MSEC 50

#define LOG_DISPLAYED_FILENAME_MAX_LENGTH 32
#define LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH 5
#define LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH 6
#define LOG_FUNC_INDENT_SPACES 3
#define LOG_INDENT_MAX 30
#define LOG_LEVEL_NAME_MAX_LENGTH 10

#define LOG_CTXT_INFO_FMT                                                      \
  "%06" PRIu64 " %s %08lX %-*.*s %-*.*s %-*.*s:%04u   %*s"
#define LOG_CTXT_INFO_ID_FMT                                                   \
  "%06" PRIu64 " %s %08lX %-*.*s %-*.*s %-*.*s:%04u   [%lu]%*s"

#define LOG_MAGMA_REPO_ROOT "/oai/"
#define MAX_TIME_STR_LEN 32
//-------------------------------

typedef unsigned long log_message_number_t;

typedef enum {
  MIN_LOG_TCP_STATE      = 0,
  LOG_TCP_STATE_DISABLED = MIN_LOG_TCP_STATE,
  LOG_TCP_STATE_NOT_CONNECTED,
  LOG_TCP_STATE_CONNECTING,
  LOG_TCP_STATE_CONNECTED,
  MAX_LOG_TCP_STATE
} log_tcp_state_t;

typedef struct oai_log_handler_s {
  void (*log_start_use)(void);
  log_queue_item_t* (*get_log_queue_item)(void);
  void (*log)(log_queue_item_t* new_item_p);
  void (*free_log_queue_item)(log_queue_item_t** item_p);
} oai_log_handler_t;

typedef struct oai_shared_log_handler_s {
  shared_log_queue_item_t* (*get_log_queue_item)(sh_ts_log_app_id_t);
  void (*log)(shared_log_queue_item_t* new_item_p);
  void (*free_log_queue_item)(shared_log_queue_item_t* item_p);
} oai_shared_log_handler_t;

#define ANSI_CODE_MAX_LENGTH 32

/*! \struct  oai_log_t
 * \brief Structure containing all the logging utility internal variables.
 */
typedef struct oai_log_s {
  // may be good to use stream instead of file descriptor when
  // logging somewhere else of the console.
  FILE* log_fd;         /*!< \brief output stream */
  bool is_output_is_fd; /* We may want to not use syslog even if exe is a daemon
                         */
  bool is_async;        /* We way want no buffering */
  bool is_ansi_codes;   /* ANSI codes for color in console output */
  bstring bserver_address; /*!< \brief TCP remote (or local) server hostname */
  bstring bserver_port;    /*!< \brief TCP remote (or local) server port     */
  log_tcp_state_t tcp_state; /*!< \brief State of the client TCP connection */

  char log_proto2str[MAX_LOG_PROTOS]
                    [LOG_MAX_PROTO_NAME_LENGTH]; /*!< \brief Convert log client
                                                    (protocol/layer) id into
                                                    human readable log user name
                                                  */
  char log_level2str[MAX_LOG_LEVEL]
                    [LOG_LEVEL_NAME_MAX_LENGTH]; /*!< \brief Convert log level
                                                    id into human readable log
                                                    level string */
  char log_level2ansi[MAX_LOG_LEVEL]
                     [ANSI_CODE_MAX_LENGTH]; /*!< \brief Convert log level id
                                                into human readable log level
                                                string */
  int log_start_time_second; /*!< \brief Logging utility reference time */
  log_level_t log_level[MAX_LOG_PROTOS]; /*!< \brief Loglevel id of each client
                                            (protocol/layer) */
  int log_level2syslog[MAX_LOG_LEVEL];
  log_message_number_t
      log_message_number; /*!< \brief Counter of log message        */
  hash_table_ts_t*
      thread_context_htbl; /*!< \brief Container for log_thread_ctxt_t */
  int max_threads;         /*!< \brief Maximum number of log threads */
  const char* app_name;    /*!< \brief Application name for log context */
  oai_log_handler_t
      log_handler; /*!< \brief Logging handler function pointers */
  oai_shared_log_handler_t
      shared_log_handler; /*!< \brief Logging handler function pointers */
} oai_log_t;

#define LOG_START_USE g_oai_log.log_handler.log_start_use
#define LOG g_oai_log.log_handler.log
#define LOG_GET_ITEM g_oai_log.log_handler.get_log_queue_item
#define LOG_FREE_ITEM g_oai_log.log_handler.free_log_queue_item

#define LOG_ASYNC g_oai_log.shared_log_handler.log
#define LOG_GET_ITEM_ASYNC g_oai_log.shared_log_handler.get_log_queue_item
#define LOG_FREE_ITEM_ASYNC g_oai_log.shared_log_handler.free_log_queue_item
static oai_log_t g_oai_log = {
    0}; /*!< \brief  logging utility internal variables global var definition*/

static void log_connect_to_server(void);
static void log_message_finish_sync(log_queue_item_t* messageP);
static void log_exit(void);
void log_message_finish_async(struct shared_log_queue_item_s* messageP);

task_zmq_ctx_t log_task_zmq_ctx;
static int timer_id = -1;

//------------------------------------------------------------------------------
static log_queue_item_t* new_queue_item(void) {
  log_queue_item_t* item_p = calloc(1, sizeof(log_queue_item_t));
  AssertFatal((item_p), "Allocation of log container failed");
  item_p->bstr = bfromcstralloc(LOG_MESSAGE_MIN_ALLOC_SIZE, "");
  AssertFatal((item_p->bstr), "Allocation of buf in log container failed");
  return item_p;
}
//------------------------------------------------------------------------------
// Get a new queue item from the free lfds stack for async logging or from heap.
static shared_log_queue_item_t* get_log_queue_item_async(
    sh_ts_log_app_id_t app_id) {
  shared_log_queue_item_t* new_item_p = NULL;
  assert(g_oai_log.is_async);
  new_item_p = get_new_log_queue_item(SH_TS_LOG_TXT);
  return new_item_p;
}
//------------------------------------------------------------------------------
static log_queue_item_t* get_log_queue_item_sync(void) {
  log_queue_item_t* new_item_p = NULL;
  new_item_p                   = new_queue_item();
  AssertFatal(new_item_p, "Out of memory error");
  return new_item_p;
}
//------------------------------------------------------------------------------

static void free_log_queue_item_sync(log_queue_item_t** item_p) {
  btrunc((*item_p)->bstr, 0);
  bdestroy((*item_p)->bstr);
  free_wrapper((void**) item_p);
}
//------------------------------------------------------------------------------

static void free_log_queue_item_async(shared_log_queue_item_t* item_p) {
  assert(g_oai_log.is_async);
  shared_log_reuse_item(item_p);
}

//------------------------------------------------------------------------------
static void log_start_use_sync(void) {
  pthread_t p = pthread_self();
  hashtable_rc_t hash_rc =
      hashtable_ts_is_key_exists(g_oai_log.thread_context_htbl, (hash_key_t) p);
  if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
    log_thread_ctxt_t* thread_ctxt = calloc(1, sizeof(log_thread_ctxt_t));
    if (thread_ctxt) {
      thread_ctxt->tid = p;
      hash_rc          = hashtable_ts_insert(
          g_oai_log.thread_context_htbl, (hash_key_t) p, thread_ctxt);
      if (HASH_TABLE_OK != hash_rc) {
        OAI_FPRINTF_ERR("Error Could not register log thread context\n");
        free_wrapper((void**) &thread_ctxt);
      }
    } else {
      OAI_FPRINTF_ERR("Error Could not create log thread context\n");
    }
  }
}

static void log_start_use_async(void) {
  assert(g_oai_log.is_async);
  shared_log_start_use();
}
static void init_syslog(void) {
  // Initialize syslog params
  // Log to console on failure, log with PID, log without buffering, user log
  openlog(g_oai_log.app_name, LOG_CONS | LOG_PID | LOG_NDELAY, LOG_USER);
  g_oai_log.log_fd          = NULL;
  g_oai_log.is_output_is_fd = false;
}

static void init_console(void) {
  setvbuf(stdout, NULL, _IONBF, 0);
  g_oai_log.log_fd          = stdout;
  g_oai_log.is_output_is_fd = true;
}

//------------------------------------------------------------------------------
// Check if we should be logging the message for the said module
static bool log_is_enabled(
    const log_level_t log_levelP, const log_proto_t protoP) {
  if ((MIN_LOG_PROTOS > protoP) || (MAX_LOG_PROTOS <= protoP)) {
    return false;
  }
  if ((MIN_LOG_LEVEL > log_levelP) || (MAX_LOG_LEVEL <= log_levelP)) {
    return false;
  }
  if (log_levelP > g_oai_log.log_level[protoP]) {
    return false;
  }
  return true;
}

//------------------------------------------------------------------------------
// Get the associated thread context for the current thread allocating if
// required
static void get_thread_context(log_thread_ctxt_t** thread_ctxt) {
  hashtable_rc_t hash_rc = HASH_TABLE_OK;

  if (NULL == *thread_ctxt) {
    pthread_t p = pthread_self();
    hash_rc     = hashtable_ts_get(
        g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) thread_ctxt);
    if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
      // Initialize thread context
      LOG_START_USE();
      hash_rc = hashtable_ts_get(
          g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) thread_ctxt);
      AssertFatal(
          HASH_TABLE_KEY_NOT_EXISTS != hash_rc,
          "Could not get new log thread context\n");
    }
  }
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      free(received_message_p);
      log_exit();
    } break;

    default: { } break; }

  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static int handle_timer(zloop_t* loop, int id, void* arg) {
  timer_id = -1;
  if (LOG_TCP_STATE_NOT_CONNECTED == g_oai_log.tcp_state) {
    log_connect_to_server();
    timer_id = start_timer(
        &log_task_zmq_ctx, LOG_CONNECT_PERIOD_MSEC, TIMER_REPEAT_ONCE,
        handle_timer, NULL);
  } else {
    // log_flush_messages ();
    timer_id = start_timer(
        &log_task_zmq_ctx, LOG_FLUSH_PERIOD_MSEC, TIMER_REPEAT_ONCE,
        handle_timer, NULL);
  }
  return 0;
}

//------------------------------------------------------------------------------
static void* log_thread(__attribute__((unused)) void* args_p) {
  itti_mark_task_ready(TASK_LOG);
  init_task_context(
      TASK_LOG, (task_id_t[]){}, 0, handle_message, &log_task_zmq_ctx);

  timer_id = start_timer(
      &log_task_zmq_ctx, LOG_FLUSH_PERIOD_MSEC, TIMER_REPEAT_ONCE, handle_timer,
      NULL);

  LOG_START_USE();

  zloop_start(log_task_zmq_ctx.event_loop);
  log_exit();
  return NULL;
}

//------------------------------------------------------------------------------
static void log_connect_to_server(void) {
  struct addrinfo hints;
  struct addrinfo *result, *rp;
  int sfd = 0, s;

  g_oai_log.tcp_state = LOG_TCP_STATE_CONNECTING;

  // man getaddrinfo:
  /* Obtain address(es) matching host/port */
  memset(&hints, 0, sizeof(struct addrinfo));
  hints.ai_family   = AF_UNSPEC;   /* Allow IPv4 or IPv6 */
  hints.ai_socktype = SOCK_STREAM; /* Stream socket */
  hints.ai_flags    = 0;
  hints.ai_protocol = IPPROTO_TCP; /* TCP protocol */

  s = getaddrinfo(
      bdata(g_oai_log.bserver_address), bdata(g_oai_log.bserver_port), &hints,
      &result);
  if (s != 0) {
    g_oai_log.tcp_state = LOG_TCP_STATE_NOT_CONNECTED;
    OAI_FPRINTF_ERR(
        "Could not connect to log server: getaddrinfo: %s\n", gai_strerror(s));
    return;
  }

  /* getaddrinfo() returns a list of address structures.
      Try each address until we successfully connect(2).
      If socket(2) (or connect(2)) fails, we (close the socket
      and) try the next address. */
  for (rp = result; rp != NULL; rp = rp->ai_next) {
    sfd = socket(rp->ai_family, rp->ai_socktype, rp->ai_protocol);
    if (sfd == -1) continue;

    if (connect(sfd, rp->ai_addr, rp->ai_addrlen) != -1) break; /* Success */

    close(sfd);
  }

  freeaddrinfo(result); /* No longer needed */

  if (rp == NULL) { /* No address succeeded */
    g_oai_log.tcp_state = LOG_TCP_STATE_NOT_CONNECTED;
    OAI_FPRINTF_ERR(
        "Could not connect to log server %s:%s\n",
        bdata(g_oai_log.bserver_address), bdata(g_oai_log.bserver_port));
    return;
  }

  fcntl(sfd, F_SETFL, fcntl(sfd, F_GETFL, 0) | O_NONBLOCK);

  g_oai_log.log_fd = fdopen(sfd, "w");
  if (NULL == g_oai_log.log_fd) {
    g_oai_log.tcp_state = LOG_TCP_STATE_NOT_CONNECTED;
    close(sfd);
    OAI_FPRINTF_ERR(
        "ERROR: Could not associate a stream with the TCP socket file "
        "descriptor\n");
    OAI_FPRINTF_ERR("ERROR: errno %d: %s\n", errno, strerror(errno));
    return;
  }
  OAI_FPRINTF_INFO(
      "Connected to log server %s:%s\n", bdata(g_oai_log.bserver_address),
      bdata(g_oai_log.bserver_port));
  g_oai_log.tcp_state = LOG_TCP_STATE_CONNECTED;
}

//------------------------------------------------------------------------------
static void log_sync(log_queue_item_t* new_item_p) {
  log_string(new_item_p->log_level, bdata(new_item_p->bstr));
  flush_log(MIN_LOG_LEVEL);

  // Release the log_item
  free_log_queue_item_sync(&new_item_p);
}
static void log_async(shared_log_queue_item_t* new_item_p) {
  log_string(new_item_p->log.log_level, bdata(new_item_p->bstr));
}
//------------------------------------------------------------------------------
// for sync or async logging
static void log_init_handler(bool async) {
  if (async) {
    g_oai_log.log_handler.log_start_use             = log_start_use_async;
    g_oai_log.shared_log_handler.log                = log_async;
    g_oai_log.shared_log_handler.get_log_queue_item = get_log_queue_item_async;
    g_oai_log.shared_log_handler.free_log_queue_item =
        free_log_queue_item_async;
  } else {
    g_oai_log.log_handler.log_start_use       = log_start_use_sync;
    g_oai_log.log_handler.log                 = log_sync;
    g_oai_log.log_handler.get_log_queue_item  = get_log_queue_item_sync;
    g_oai_log.log_handler.free_log_queue_item = free_log_queue_item_sync;
  }
}

//------------------------------------------------------------------------------
static void log_signal_callback_handler(int signum) {
  OAI_FPRINTF_ERR("Caught signal SIGPIPE %d\n", signum);
  if (LOG_TCP_STATE_DISABLED != g_oai_log.tcp_state) {
    // Let ITTI LOG Timer do the reconnection
    g_oai_log.tcp_state = LOG_TCP_STATE_NOT_CONNECTED;
    return;
  }
}
//------------------------------------------------------------------------------
void log_configure(const log_config_t* const config) {
  if (NULL == config) {
    log_message(
        NULL, OAILOG_LEVEL_WARNING, LOG_UTIL, __FILE__, __LINE__,
        "Log config unset, defaulting to syslog\n");
    return;
  }
  if ((MAX_LOG_LEVEL > config->udp_log_level) &&
      (MIN_LOG_LEVEL <= config->udp_log_level))
    g_oai_log.log_level[LOG_UDP] = config->udp_log_level;
  if ((MAX_LOG_LEVEL > config->gtpv1u_log_level) &&
      (MIN_LOG_LEVEL <= config->gtpv1u_log_level))
    g_oai_log.log_level[LOG_GTPV1U] = config->gtpv1u_log_level;
  if ((MAX_LOG_LEVEL > config->gtpv2c_log_level) &&
      (MIN_LOG_LEVEL <= config->gtpv2c_log_level))
    g_oai_log.log_level[LOG_GTPV2C] = config->gtpv2c_log_level;
  if ((MAX_LOG_LEVEL > config->sctp_log_level) &&
      (MIN_LOG_LEVEL <= config->sctp_log_level))
    g_oai_log.log_level[LOG_SCTP] = config->sctp_log_level;
  if ((MAX_LOG_LEVEL > config->s1ap_log_level) &&
      (MIN_LOG_LEVEL <= config->s1ap_log_level))
    g_oai_log.log_level[LOG_S1AP] = config->s1ap_log_level;
  if ((MAX_LOG_LEVEL > config->mme_app_log_level) &&
      (MIN_LOG_LEVEL <= config->mme_app_log_level)) {
    g_oai_log.log_level[LOG_MME_APP] = config->mme_app_log_level;
    g_oai_log.log_level[LOG_AMF_APP] = config->mme_app_log_level;
    g_oai_log.log_level[LOG_NGAP]    = config->mme_app_log_level;
    g_oai_log.log_level[LOG_NAS_AMF] = config->mme_app_log_level;
  }

  if ((MAX_LOG_LEVEL > config->nas_log_level) &&
      (MIN_LOG_LEVEL <= config->nas_log_level)) {
    g_oai_log.log_level[LOG_NAS]     = config->nas_log_level;
    g_oai_log.log_level[LOG_NAS_EMM] = config->nas_log_level;
    g_oai_log.log_level[LOG_NAS_ESM] = config->nas_log_level;
  }
  if ((MAX_LOG_LEVEL > config->spgw_app_log_level) &&
      (MIN_LOG_LEVEL <= config->spgw_app_log_level))
    g_oai_log.log_level[LOG_SPGW_APP] = config->spgw_app_log_level;
  if ((MAX_LOG_LEVEL > config->s11_log_level) &&
      (MIN_LOG_LEVEL <= config->s11_log_level))
    g_oai_log.log_level[LOG_S11] = config->s11_log_level;
  if ((MAX_LOG_LEVEL > config->s6a_log_level) &&
      (MIN_LOG_LEVEL <= config->s6a_log_level))
    g_oai_log.log_level[LOG_S6A] = config->s6a_log_level;
  if ((MAX_LOG_LEVEL > config->util_log_level) &&
      (MIN_LOG_LEVEL <= config->util_log_level))
    g_oai_log.log_level[LOG_UTIL] = config->util_log_level;
  if ((MAX_LOG_LEVEL > config->itti_log_level) &&
      (MIN_LOG_LEVEL <= config->itti_log_level))
    g_oai_log.log_level[LOG_ITTI] = config->itti_log_level;
  if ((MAX_LOG_LEVEL > config->async_system_log_level) &&
      (MIN_LOG_LEVEL <= config->async_system_log_level))
    g_oai_log.log_level[LOG_ASYNC_SYSTEM] = config->async_system_log_level;
  g_oai_log.is_async      = config->is_output_thread_safe;
  g_oai_log.is_ansi_codes = config->color;
  log_init_handler(g_oai_log.is_async);

  if (config->output) {
    if (1 ==
        biseqcstrcaseless(config->output, LOG_CONFIG_STRING_OUTPUT_SYSLOG)) {
      // Output to syslog
      init_syslog();
      return;
    }
    if (1 ==
        biseqcstrcaseless(config->output, LOG_CONFIG_STRING_OUTPUT_CONSOLE)) {
      init_console();
      return;
    }
    // if seems to be a file path
    if (('.' == bchar(config->output, 0)) ||
        ('/' == bchar(config->output, 0))) {
      g_oai_log.log_fd = fopen(bdata(config->output), "w");
      AssertFatal(
          NULL != g_oai_log.log_fd, "Could not open log file %s : %s",
          bdata(config->output), strerror(errno));
      g_oai_log.is_output_is_fd = true;
    } else {
      // may be a TCP server address host:portnum
      g_oai_log.bserver_address = bstrcpy(config->output);
      int pos                   = bstrchr(g_oai_log.bserver_address, ':');
      if (BSTR_ERR != pos) {
        g_oai_log.bserver_port =
            bmidstr(g_oai_log.bserver_address, pos + 1, 1024);
        btrunc(g_oai_log.bserver_address, pos);
      }
      int server_port = atoi((const char*) g_oai_log.bserver_port->data);
      AssertFatal(
          1024 <= server_port, "Invalid Server TCP port %d/%s", server_port,
          bdata(g_oai_log.bserver_port));
      AssertFatal(
          65535 >= server_port, "Invalid Server TCP port %d/%s", server_port,
          bdata(g_oai_log.bserver_port));
      g_oai_log.tcp_state       = LOG_TCP_STATE_NOT_CONNECTED;
      g_oai_log.is_output_is_fd = true;
      log_connect_to_server();
    }
  }
}
/*
 * Disabling below function to get actual time of the day in the logs
 */
#if 0
//------------------------------------------------------------------------------
static void log_get_elapsed_time_since_start(struct timeval * const elapsed_time)
{
  // no thread safe but do not matter a lot
  gettimeofday(elapsed_time, NULL);
  // no timersub call for fastest operations
  elapsed_time->tv_sec = elapsed_time->tv_sec - g_oai_log.log_start_time_second;
}
#endif
//------------------------------------------------------------------------------
static void log_get_readable_cur_time(time_t* cur_time, char* time_str) {
  // get the current local time
  time(cur_time);
  struct tm* cur_local_time;
  cur_local_time = localtime(cur_time);
  // get the current local time in readable string format
  strftime(
      time_str, MAX_TIME_STR_LEN, "%a %b %d %H:%M:%S %Y",
      (const struct tm*) cur_local_time);
}

//------------------------------------------------------------------------------
const char* log_level_int2str(const log_level_t log_level) {
  if ((MAX_LOG_LEVEL > log_level) && (MIN_LOG_LEVEL <= log_level)) {
    return g_oai_log.log_level2str[log_level];
  }
  return "INVALID_LOG_LEVEL";
}

//------------------------------------------------------------------------------
log_level_t log_level_str2int(const char* const log_level_str) {
  log_level_t log_level;

  if (log_level_str) {
    for (log_level = MIN_LOG_LEVEL; log_level < MAX_LOG_LEVEL; log_level++) {
      if (0 ==
          strcasecmp(log_level_str, &g_oai_log.log_level2str[log_level][0])) {
        return log_level;
      }
    }
  }
  // By default
  return MAX_LOG_LEVEL;  // == invalid
}
//------------------------------------------------------------------------------
int log_init(
    const char* app_name, const log_level_t default_log_levelP,
    const int max_threadsP) {
  // init glog logging
  init_logging(app_name, default_log_levelP);

  int i                     = 0;
  struct timeval start_time = {.tv_sec = 0, .tv_usec = 0};

  OAI_FPRINTF_INFO("Initializing OAI Logging\n");
  signal(SIGPIPE, log_signal_callback_handler);

  g_oai_log.log_fd = NULL;

  gettimeofday(&start_time, NULL);
  g_oai_log.log_start_time_second = (int) start_time.tv_sec;

  OAI_FPRINTF_INFO("Initializing OAI Logging to syslog\n");
  bstring b = bfromcstr("Logging thread context hashtable");
  g_oai_log.thread_context_htbl =
      hashtable_ts_create(LOG_MESSAGE_MIN_ALLOC_SIZE, NULL, free_wrapper, b);
  bdestroy_wrapper(&b);
  AssertFatal(
      NULL != g_oai_log.thread_context_htbl,
      "Could not create hashtable for Log!\n");
  g_oai_log.thread_context_htbl->log_enabled = false;
  g_oai_log.max_threads                      = max_threadsP;
  g_oai_log.app_name                         = app_name;
  g_oai_log.is_async                         = false;

  snprintf(
      &g_oai_log.log_proto2str[LOG_SCTP][0], LOG_MAX_PROTO_NAME_LENGTH, "SCTP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_UDP][0], LOG_MAX_PROTO_NAME_LENGTH, "UDP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_GTPV1U][0], LOG_MAX_PROTO_NAME_LENGTH,
      "GTPv1-U");
  snprintf(
      &g_oai_log.log_proto2str[LOG_GTPV2C][0], LOG_MAX_PROTO_NAME_LENGTH,
      "GTPv2-C");
  snprintf(
      &g_oai_log.log_proto2str[LOG_S1AP][0], LOG_MAX_PROTO_NAME_LENGTH, "S1AP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_MME_APP][0], LOG_MAX_PROTO_NAME_LENGTH,
      "MME-APP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_AMF_APP][0], LOG_MAX_PROTO_NAME_LENGTH,
      "AMF-APP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_NGAP][0], LOG_MAX_PROTO_NAME_LENGTH, "NGAP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_NAS_AMF][0], LOG_MAX_PROTO_NAME_LENGTH,
      "NAS-AMF");
  snprintf(
      &g_oai_log.log_proto2str[LOG_NAS][0], LOG_MAX_PROTO_NAME_LENGTH, "NAS");
  snprintf(
      &g_oai_log.log_proto2str[LOG_NAS_EMM][0], LOG_MAX_PROTO_NAME_LENGTH,
      "NAS-EMM");
  snprintf(
      &g_oai_log.log_proto2str[LOG_NAS_ESM][0], LOG_MAX_PROTO_NAME_LENGTH,
      "NAS-ESM");
  snprintf(
      &g_oai_log.log_proto2str[LOG_SPGW_APP][0], LOG_MAX_PROTO_NAME_LENGTH,
      "SPGW-APP");
  snprintf(
      &g_oai_log.log_proto2str[LOG_S11][0], LOG_MAX_PROTO_NAME_LENGTH, "S11");
  snprintf(
      &g_oai_log.log_proto2str[LOG_S6A][0], LOG_MAX_PROTO_NAME_LENGTH, "S6A");
  snprintf(
      &g_oai_log.log_proto2str[LOG_SGW_S8][0], LOG_MAX_PROTO_NAME_LENGTH,
      "SGW_S8");
  snprintf(
      &g_oai_log.log_proto2str[LOG_SECU][0], LOG_MAX_PROTO_NAME_LENGTH, "SECU");
  snprintf(
      &g_oai_log.log_proto2str[LOG_UTIL][0], LOG_MAX_PROTO_NAME_LENGTH, "UTIL");
  snprintf(
      &g_oai_log.log_proto2str[LOG_CONFIG][0], LOG_MAX_PROTO_NAME_LENGTH,
      "CONFIG");
  snprintf(
      &g_oai_log.log_proto2str[LOG_ITTI][0], LOG_MAX_PROTO_NAME_LENGTH, "ITTI");
  snprintf(
      &g_oai_log.log_proto2str[LOG_ASYNC_SYSTEM][0], LOG_MAX_PROTO_NAME_LENGTH,
      "CMD");
  snprintf(
      &g_oai_log.log_proto2str[LOG_ASSERT][0], LOG_MAX_PROTO_NAME_LENGTH,
      "ASSERT");

  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_TRACE][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "TRACE");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_DEBUG][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "DEBUG");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_INFO][0], LOG_LEVEL_NAME_MAX_LENGTH,
      "INFO");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_NOTICE][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "NOTICE");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_WARNING][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "WARNING");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_ERROR][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "ERROR");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_CRITICAL][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "CRITICAL");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_ALERT][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "ALERT");
  snprintf(
      &g_oai_log.log_level2str[OAILOG_LEVEL_EMERGENCY][0],
      LOG_LEVEL_NAME_MAX_LENGTH, "EMERGENCY");

  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_TRACE][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_WHITE);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_DEBUG][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_GREEN);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_INFO][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_CYAN);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_NOTICE][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_BLUE);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_WARNING][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_YELLOW);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_ERROR][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_RED);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_CRITICAL][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_REV_RED);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_ALERT][0], ANSI_CODE_MAX_LENGTH,
      ANSI_COLOR_FG_REV_RED);
  snprintf(
      &g_oai_log.log_level2ansi[OAILOG_LEVEL_EMERGENCY][0],
      ANSI_CODE_MAX_LENGTH, ANSI_COLOR_FG_REV_RED);

  for (i = MIN_LOG_PROTOS; i < MAX_LOG_PROTOS; i++) {
    g_oai_log.log_level[i] = default_log_levelP;
  }

  // Map OAI log levels to syslog
  g_oai_log.log_level2syslog[OAILOG_LEVEL_EMERGENCY] = LOG_EMERG;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_ALERT]     = LOG_ALERT;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_CRITICAL]  = LOG_CRIT;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_ERROR]     = LOG_ERR;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_WARNING]   = LOG_WARNING;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_NOTICE]    = LOG_NOTICE;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_INFO]      = LOG_INFO;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_DEBUG]     = LOG_DEBUG;
  g_oai_log.log_level2syslog[OAILOG_LEVEL_TRACE]     = LOG_DEBUG;

  log_init_handler(g_oai_log.is_async);

  // We should probably initialize default here which is console but we log a
  // lot before we actually set the log config, so treat syslog as default for
  // now.
  init_syslog();
  log_message(
      NULL, OAILOG_LEVEL_INFO, LOG_UTIL, __FILE__, __LINE__,
      "Initializing OAI logging Done\n");
  return 0;
}

//------------------------------------------------------------------------------
// listen to ITTI events
void log_itti_connect(void) {
  if (g_oai_log.is_async) {
    int rv = 0;
    rv     = itti_create_task(TASK_LOG, &log_thread, NULL);
    AssertFatal(rv == 0, "Create task for OAI logging failed!\n");
  }
}
//------------------------------------------------------------------------------
void log_flush_message(struct shared_log_queue_item_s* item_p) {
  int rv     = 0;
  int rv_put = 0;

  if (blength(item_p->bstr) > 0) {
    if (g_oai_log.is_output_is_fd) {
      if (g_oai_log.log_fd) {
        rv_put = fputs((const char*) item_p->bstr->data, g_oai_log.log_fd);
        if (rv_put < 0) {
          // error occured
          OAI_FPRINTF_ERR("Error while writing log %d\n", rv_put);
          rv = fclose(g_oai_log.log_fd);
          if (rv != 0) {
            OAI_FPRINTF_ERR(
                "Error while closing Log file stream: %s\n", strerror(errno));
          }
          // do not exit
          if (LOG_TCP_STATE_DISABLED != g_oai_log.tcp_state) {
            // Let ITTI LOG Timer do the reconnection
            g_oai_log.tcp_state = LOG_TCP_STATE_NOT_CONNECTED;
            return;
          }
        }
        fflush(g_oai_log.log_fd);
      }
    } else {
      syslog(item_p->log.log_level, "%s", bdata(item_p->bstr));
    }
  }
}

//------------------------------------------------------------------------------
static void log_exit(void) {
  assert(g_oai_log.is_async);

  OAI_FPRINTF_INFO("[TRACE] Entering %s\n", __FUNCTION__);
  stop_timer(&log_task_zmq_ctx, timer_id);
  destroy_task_context(&log_task_zmq_ctx);
  if (g_oai_log.log_fd) {
    int rv = fflush(g_oai_log.log_fd);

    if (rv != 0) {
      OAI_FPRINTF_ERR(
          "Error while flushing stream of Log file: %s", strerror(errno));
    }

    rv = fclose(g_oai_log.log_fd);

    if (rv != 0) {
      OAI_FPRINTF_ERR("Error while closing Log file: %s", strerror(errno));
    }
  }
  if (!g_oai_log.is_output_is_fd) {
    closelog();
  }
  hashtable_ts_destroy(g_oai_log.thread_context_htbl);
  bdestroy_wrapper(&g_oai_log.bserver_address);
  bdestroy_wrapper(&g_oai_log.bserver_port);
  OAI_FPRINTF_INFO("[TRACE] Leaving %s\n", __FUNCTION__);

  OAI_FPRINTF_INFO("TASK_LOG terminated\n");
  pthread_exit(NULL);
}
//------------------------------------------------------------------------------
static void log_stream_hex_sync(
    const log_level_t log_levelP, const log_proto_t protoP,
    const char* const source_fileP, const unsigned int line_numP,
    const char* const messageP, const char* const streamP, const size_t sizeP) {
  log_queue_item_t* message      = NULL;
  size_t octet_index             = 0;
  int rv                         = 0;
  log_thread_ctxt_t* thread_ctxt = NULL;

  get_thread_context(&thread_ctxt);
  if (messageP) {
    log_message_start_sync(
        thread_ctxt, log_levelP, protoP, &message, source_fileP, line_numP,
        "%s (%ld bytes)", messageP, sizeP);
  } else {
    log_message_start_sync(
        thread_ctxt, log_levelP, protoP, &message, source_fileP, line_numP,
        "%p dumped(%ld bytes):", streamP, sizeP);
  }
  if ((streamP) && (message)) {
    for (octet_index = 0; octet_index < sizeP; octet_index++) {
      // do not call log_message_add_sync(), too much overhead for sizeP*3chars
      rv = bformata(
          message->bstr, " %02x", (streamP[octet_index]) & (uint) 0x00ff);

      if (BSTR_ERR == rv) {
        OAI_FPRINTF_ERR("Error while logging message\n");
      }
    }
    log_message_finish_sync(message);
  }
}
//------------------------------------------------------------------------------
static void log_stream_hex_async(
    const log_level_t log_levelP, const log_proto_t protoP,
    const char* const source_fileP, const unsigned int line_numP,
    const char* const messageP, const char* const streamP, const size_t sizeP) {
  struct shared_log_queue_item_s* message = NULL;
  size_t octet_index                      = 0;
  int rv                                  = 0;
  log_thread_ctxt_t* thread_ctxt          = NULL;
  hashtable_rc_t hash_rc                  = HASH_TABLE_OK;

  pthread_t p = pthread_self();
  hash_rc     = hashtable_ts_get(
      g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
  if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
    // make the thread safe LFDS collections usable by this thread
    LOG_START_USE();
  }
  hash_rc = hashtable_ts_get(
      g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
  AssertFatal(NULL != thread_ctxt, "Could not get new log thread context\n");
  if (messageP) {
    log_message_start_async(
        thread_ctxt, log_levelP, protoP, &message, source_fileP, line_numP,
        "hex stream ");
    if (!message) return;
    rv = bformata(message->bstr, "%s", messageP);
  } else {
    log_message_start_async(
        thread_ctxt, log_levelP, protoP, &message, source_fileP, line_numP,
        "hex stream (%ld bytes):", sizeP);
    if (!message) return;
  }
  if ((streamP) && (message)) {
    for (octet_index = 0; octet_index < sizeP; octet_index++) {
      // do not call log_message_add_async(), too much overhead for sizeP*3chars
      rv = bformata(
          message->bstr, " %02x", (streamP[octet_index]) & (uint) 0x00ff);

      if (BSTR_ERR == rv) {
        OAI_FPRINTF_ERR("Error while logging message\n");
      }
    }
    log_message_finish_async(message);
  }
}

//------------------------------------------------------------------------------
void log_stream_hex(
    const log_level_t log_levelP, const log_proto_t protoP,
    const char* const source_fileP, const unsigned int line_numP,
    const char* const messageP, const char* const streamP, const size_t sizeP) {
  if (g_oai_log.is_async) {
    log_stream_hex_async(
        log_levelP, protoP, source_fileP, line_numP, messageP, streamP, sizeP);
  } else {
    log_stream_hex_sync(
        log_levelP, protoP, source_fileP, line_numP, messageP, streamP, sizeP);
  }
}
//------------------------------------------------------------------------------
void log_stream_hex_array(
    const log_level_t log_levelP, const log_proto_t protoP,
    const char* const source_fileP, const unsigned int line_numP,
    const char* const messageP, const char* const streamP, const size_t sizeP) {
  struct shared_log_queue_item_s* message = NULL;
  unsigned long octet_index               = 0;
  unsigned long index                     = 0;
  log_thread_ctxt_t* thread_ctxt          = NULL;

  get_thread_context(&thread_ctxt);

  if (messageP) {
    log_message(
        thread_ctxt, log_levelP, protoP, source_fileP, line_numP, "%s\n",
        messageP);
  }
  log_message(
      thread_ctxt, log_levelP, protoP, source_fileP, line_numP,
      "------+-------------------------------------------------|\n");
  log_message(
      thread_ctxt, log_levelP, protoP, source_fileP, line_numP,
      "      |  0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f |\n");
  log_message(
      thread_ctxt, log_levelP, protoP, source_fileP, line_numP,
      "------+-------------------------------------------------|\n");

  if (streamP) {
    for (octet_index = 0; octet_index < sizeP; octet_index++) {
      if ((octet_index % 16) == 0) {
        if (octet_index != 0) {
          log_message_add_async(message, " |");
          log_message_finish_async(message);
        }
        log_message_start_async(
            thread_ctxt, log_levelP, protoP, &message, source_fileP, line_numP,
            " %04ld |", octet_index);
      }

      /*
       * Print every single octet in hexadecimal form
       */
      log_message_add_async(
          message, " %02x", ((unsigned char*) streamP)[octet_index]);
    }
    /*
     * Append enough spaces and put final pipe
     */
    for (index = octet_index % 16; index < 16; ++index) {
      log_message_add_async(message, "   ");
    }
    log_message_add_async(message, " |");
    log_message_finish(message);
  }
}

//------------------------------------------------------------------------------
void log_message_add_async(
    struct shared_log_queue_item_s* messageP, char* format, ...) {
  va_list args;
  int rv = 0;

  if (messageP) {
    va_start(args, format);
    rv = bvcformata(
        messageP->bstr, 4096, format, args);  // big number, see bvcformata
    va_end(args);

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR("Error while logging message\n");
    }
  }
}
//------------------------------------------------------------------------------
void log_message_add_sync(log_queue_item_t* messageP, char* format, ...) {
  va_list args;
  int rv = 0;

  if (messageP) {
    va_start(args, format);
    rv = bvcformata(
        messageP->bstr, 4096, format, args);  // big number, see bvcformata
    va_end(args);

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR("Error while logging message\n");
    }
  }
}
//------------------------------------------------------------------------------
static void log_message_finish_sync(log_queue_item_t* messageP) {
  // flush everything
  flush_log(MIN_LOG_LEVEL);
  int rv = 0;

  if (NULL == messageP) {
    log_message(
        NULL, OAILOG_LEVEL_WARNING, LOG_UTIL, __FILE__, __LINE__,
        "Calling finish on a NULL message\n");
    return;
  }

  rv = bcatcstr(messageP->bstr, "\n");

  if (BSTR_ERR == rv) {
    OAI_FPRINTF_ERR("Error while logging message\n");
    goto error_event;
  }
  LOG(messageP);
  return;

error_event:
  LOG_FREE_ITEM(&messageP);
}
//------------------------------------------------------------------------------
void log_message_finish_async(struct shared_log_queue_item_s* messageP) {
  // flush everything
  flush_log(MIN_LOG_LEVEL);
  int rv = 0;

  if (messageP) {
    if (g_oai_log.is_ansi_codes) {
      rv = bformata(messageP->bstr, "%s\n", ANSI_COLOR_RESET);
    } else {
      rv = bcatcstr(messageP->bstr, "\n");
    }

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR("Error while logging message\n");
    }
    // send message
    shared_log_item(messageP);
  }
}
//------------------------------------------------------------------------------
void log_message_finish(void* messageP) {
  if (NULL == messageP) {
    log_message(
        NULL, OAILOG_LEVEL_WARNING, LOG_UTIL, __FILE__, __LINE__,
        "Calling finish on a NULL message\n");
    return;
  }
  if (g_oai_log.is_async) {
    log_message_finish_async((struct shared_log_queue_item_s*) messageP);
  } else {
    log_message_finish_sync((log_queue_item_t*) messageP);
  }
}
//------------------------------------------------------------------------------
void log_message_start_sync(
    log_thread_ctxt_t* thread_ctxtP, const log_level_t log_levelP,
    const log_proto_t protoP,
    log_queue_item_t** messageP,  // Out parameter
    const char* const source_fileP, const unsigned int line_numP, char* format,
    ...) {
  va_list args;

  va_start(args, format);
  log_message_int(
      thread_ctxtP, log_levelP, protoP, (void**) messageP, source_fileP,
      line_numP, format, args);
  va_end(args);
}

//------------------------------------------------------------------------------
void log_message_start_async(
    log_thread_ctxt_t* thread_ctxtP, const log_level_t log_levelP,
    const log_proto_t protoP,
    struct shared_log_queue_item_s** messageP,  // Out parameter
    const char* const source_fileP, const unsigned int line_numP, char* format,
    ...) {
  va_list args;
  int rv                         = 0;
  int filename_length            = 0;
  log_thread_ctxt_t* thread_ctxt = thread_ctxtP;
  hashtable_rc_t hash_rc         = HASH_TABLE_OK;

  if ((MIN_LOG_PROTOS > protoP) || (MAX_LOG_PROTOS <= protoP)) {
    return;
  }
  if ((MIN_LOG_LEVEL > log_levelP) || (MAX_LOG_LEVEL <= log_levelP)) {
    return;
  }
  if (log_levelP > g_oai_log.log_level[protoP]) {
    return;
  }

  if (NULL == thread_ctxt) {
    pthread_t p = pthread_self();
    hash_rc     = hashtable_ts_get(
        g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
    if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
      // make the thread safe LFDS collections usable by this thread
      LOG_START_USE();
      hash_rc = hashtable_ts_get(
          g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
      AssertFatal(
          NULL != thread_ctxt, "Could not get new log thread context\n");
    }
  }

  if (!*messageP) {
    *messageP = get_new_log_queue_item(SH_TS_LOG_TXT);
  }

  if (*messageP) {
    struct timeval elapsed_time;
    (*messageP)->log.log_level = log_levelP;
    shared_log_get_elapsed_time_since_start(&elapsed_time);

    // get the short file name to use for printing in log
    const char* const short_source_fileP = get_short_file_name(source_fileP);

    filename_length = strlen(short_source_fileP);
    if (g_oai_log.is_ansi_codes) {
      rv = bformata(
          (*messageP)->bstr, "%s", &g_oai_log.log_level2ansi[log_levelP][0]);
    }
    if (filename_length > LOG_DISPLAYED_FILENAME_MAX_LENGTH) {
      rv = bformata(
          (*messageP)->bstr,
          "%06" PRIu64 " %05ld:%06ld %08lX %-*.*s %-*.*s %-*.*s:%04u   %*s",
          __sync_fetch_and_add(&g_oai_log.log_message_number, 1),
          elapsed_time.tv_sec, elapsed_time.tv_usec, thread_ctxt->tid,
          LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
          LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
          &g_oai_log.log_level2str[log_levelP][0],
          LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH,
          LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH,
          &g_oai_log.log_proto2str[protoP][0],
          LOG_DISPLAYED_FILENAME_MAX_LENGTH, LOG_DISPLAYED_FILENAME_MAX_LENGTH,
          &short_source_fileP
              [filename_length - LOG_DISPLAYED_FILENAME_MAX_LENGTH],
          line_numP, thread_ctxt->indent, " ");
    } else {
      rv = bformata(
          (*messageP)->bstr,
          "%06" PRIu64 " %05ld:%06ld %08lX %-*.*s %-*.*s %-*.*s:%04u   %*s",
          __sync_fetch_and_add(&g_oai_log.log_message_number, 1),
          elapsed_time.tv_sec, elapsed_time.tv_usec, thread_ctxt->tid,
          LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
          LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
          &g_oai_log.log_level2str[log_levelP][0],
          LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH,
          LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH,
          &g_oai_log.log_proto2str[protoP][0],
          LOG_DISPLAYED_FILENAME_MAX_LENGTH, LOG_DISPLAYED_FILENAME_MAX_LENGTH,
          short_source_fileP, line_numP, thread_ctxt->indent, " ");
    }

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event_start;
    }

    va_start(args, format);
    rv = bformata((*messageP)->bstr, format, args);
    va_end(args);

    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event_start;
    }
    return;
  }
  return;
error_event_start:
  // put in memory pool the message buffer
  btrunc((*messageP)->bstr, 0);
  shared_log_reuse_item(*messageP);
  *messageP = NULL;
  return;
}

//------------------------------------------------------------------------------
// hard-coded to use LOG_LEVEL_TRACE
void log_func(
    const bool is_enteringP, const log_proto_t protoP,
    const char* const source_fileP, const unsigned int line_numP,
    const char* const functionP) {
  log_thread_ctxt_t* thread_ctxt = NULL;
  hashtable_rc_t hash_rc         = HASH_TABLE_OK;
  pthread_t p                    = pthread_self();

  hash_rc = hashtable_ts_get(
      g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
  if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
    // make the thread safe LFDS collections usable by this thread
    LOG_START_USE();
  }
  hash_rc = hashtable_ts_get(
      g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
  AssertFatal(NULL != thread_ctxt, "Could not get new log thread context\n");
  if (is_enteringP) {
    log_message(
        thread_ctxt, OAILOG_LEVEL_TRACE, protoP, source_fileP, line_numP,
        "Entering %s()\n", functionP);
    thread_ctxt->indent += LOG_FUNC_INDENT_SPACES;
    if (thread_ctxt->indent > LOG_INDENT_MAX) {
      thread_ctxt->indent = LOG_INDENT_MAX;
    }
  } else {
    thread_ctxt->indent -= LOG_FUNC_INDENT_SPACES;
    if (thread_ctxt->indent < 0) {
      thread_ctxt->indent = 0;
    }
    log_message(
        thread_ctxt, OAILOG_LEVEL_TRACE, protoP, source_fileP, line_numP,
        "Leaving %s()\n", functionP);
  }
}
//------------------------------------------------------------------------------
// hard-coded to use LOG_LEVEL_TRACE
void log_func_return(
    const log_proto_t protoP, const char* const source_fileP,
    const unsigned int line_numP, const char* const functionP,
    const long return_codeP) {
  log_thread_ctxt_t* thread_ctxt = NULL;
  hashtable_rc_t hash_rc         = HASH_TABLE_OK;
  pthread_t p                    = pthread_self();

  hash_rc = hashtable_ts_get(
      g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
  if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
    // make the thread safe LFDS collections usable by this thread
    LOG_START_USE();
  }
  hash_rc = hashtable_ts_get(
      g_oai_log.thread_context_htbl, (hash_key_t) p, (void**) &thread_ctxt);
  AssertFatal(NULL != thread_ctxt, "Could not get new log thread context\n");
  thread_ctxt->indent -= LOG_FUNC_INDENT_SPACES;
  if (thread_ctxt->indent < 0) thread_ctxt->indent = 0;
  log_message(
      thread_ctxt, OAILOG_LEVEL_TRACE, protoP, source_fileP, line_numP,
      "Leaving %s() (rc=%ld)\n", functionP, return_codeP);
}
//------------------------------------------------------------------------------
void log_message(
    log_thread_ctxt_t* thread_ctxtP, const log_level_t log_levelP,
    const log_proto_t protoP, const char* const source_fileP,
    const unsigned int line_numP, const char* format, ...) {
  va_list args;
  void* new_item_p                                 = NULL;
  log_queue_item_t* new_item_p_sync                = NULL;
  struct shared_log_queue_item_s* new_item_p_async = NULL;

  va_start(args, format);
  log_message_int(
      thread_ctxtP, log_levelP, protoP, &new_item_p, source_fileP, line_numP,
      format, args);
  va_end(args);

  if (NULL == new_item_p) {
    return;
  }
  if (g_oai_log.is_async) {
    new_item_p_async = (struct shared_log_queue_item_s*) new_item_p;
    if (g_oai_log.is_ansi_codes) {
      bformata(new_item_p_async->bstr, "%s", ANSI_COLOR_RESET);
    }
    LOG_ASYNC(new_item_p_async);
  } else {
    new_item_p_sync = (log_queue_item_t*) new_item_p;
    if (g_oai_log.is_ansi_codes) {
      bformata(new_item_p_sync->bstr, "%s", ANSI_COLOR_RESET);
    }
    LOG(new_item_p_sync);
  }
}

void log_message_prefix_id(
    const log_level_t log_levelP, const log_proto_t protoP,
    const char* const source_fileP, const unsigned int line_numP,
    uint64_t prefix_id, const char* format, ...) {
  va_list args;
  void* new_item_p                                 = NULL;
  log_queue_item_t* new_item_p_sync                = NULL;
  struct shared_log_queue_item_s* new_item_p_async = NULL;

  va_start(args, format);
  log_message_int_prefix_id(
      log_levelP, protoP, &new_item_p, source_fileP, line_numP, prefix_id,
      format, args);
  va_end(args);

  if (new_item_p == NULL) {
    return;
  }
  if (g_oai_log.is_async) {
    new_item_p_async = (struct shared_log_queue_item_s*) new_item_p;
    if (g_oai_log.is_ansi_codes) {
      bformata(new_item_p_async->bstr, "%s", ANSI_COLOR_RESET);
    }
    LOG_ASYNC(new_item_p_async);
  } else {
    new_item_p_sync = (log_queue_item_t*) new_item_p;
    if (g_oai_log.is_ansi_codes) {
      bformata(new_item_p_sync->bstr, "%s", ANSI_COLOR_RESET);
    }
    LOG(new_item_p_sync);
  }
}

void log_message_int(
    log_thread_ctxt_t* const thread_ctxtP, const log_level_t log_levelP,
    const log_proto_t protoP,
    void** contextP,  // Out parameter
    const char* const source_fileP, const unsigned int line_numP,
    const char* format, va_list args) {
  int rv                                    = 0;
  size_t filename_length                    = 0;
  log_thread_ctxt_t* thread_ctxt            = thread_ctxtP;
  log_queue_item_t** sync_context_p         = NULL;
  shared_log_queue_item_t** async_context_p = NULL;
  if (!log_is_enabled(log_levelP, protoP)) {
    return;
  }
  get_thread_context(&thread_ctxt);

  assert(thread_ctxt != NULL);
  *contextP = LOG_GET_ITEM();
#if 0
  struct timeval elapsed_time;
  log_get_elapsed_time_since_start(&elapsed_time);
#endif
  time_t cur_time;

  // get the short file name to use for printing in log
  const char* const short_source_fileP = get_short_file_name(source_fileP);

  filename_length = MIN(
      (strlen(short_source_fileP) - LOG_DISPLAYED_FILENAME_MAX_LENGTH), (0));
  if (!(g_oai_log.is_async)) {
    sync_context_p = (log_queue_item_t**) contextP;
    rv             = append_log_ctx_info(
        (*sync_context_p)->bstr, &log_levelP, &protoP, line_numP,
        filename_length, thread_ctxt, &cur_time, short_source_fileP);
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
    rv = bvcformata((*sync_context_p)->bstr, 4096, format, args);  // big number
    (*sync_context_p)->log_level = g_oai_log.log_level2syslog[log_levelP];
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
  } else {
    async_context_p = (shared_log_queue_item_t**) contextP;
    rv              = append_log_ctx_info(
        (*async_context_p)->bstr, &log_levelP, &protoP, line_numP,
        filename_length, thread_ctxt, &cur_time, short_source_fileP);
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
    rv =
        bvcformata((*async_context_p)->bstr, 4096, format, args);  // big number
    (*async_context_p)->log.log_level = g_oai_log.log_level2syslog[log_levelP];
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
  }
  return;

error_event:
  if (!(g_oai_log.is_async)) {
    LOG_FREE_ITEM(sync_context_p);
  } else if (async_context_p == NULL) {
    // To guard against mutation of is_async during post-init
    // runtime, though this should not be mutated.
    return;
  } else {
    LOG_FREE_ITEM_ASYNC(*async_context_p);
  }
}

void log_message_int_prefix_id(
    const log_level_t log_levelP, const log_proto_t protoP,
    void** contextP,  // Out parameter
    const char* const source_fileP, const unsigned int line_numP,
    const uint64_t prefix_id, const char* format, va_list args) {
  int rv                                    = 0;
  size_t filename_length                    = 0;
  log_thread_ctxt_t* thread_ctxt            = NULL;
  log_queue_item_t** sync_context_p         = NULL;
  shared_log_queue_item_t** async_context_p = NULL;
  if (!log_is_enabled(log_levelP, protoP)) {
    return;
  }
  get_thread_context(&thread_ctxt);

  assert(thread_ctxt != NULL);
  *contextP = LOG_GET_ITEM();
  time_t cur_time;

  // get the short file name to use for printing in log
  const char* const short_source_fileP = get_short_file_name(source_fileP);

  filename_length = MIN(
      (strlen(short_source_fileP) - LOG_DISPLAYED_FILENAME_MAX_LENGTH), (0));
  if (!(g_oai_log.is_async)) {
    sync_context_p = (log_queue_item_t**) contextP;
    rv             = append_log_ctx_info_prefix_id(
        prefix_id, (*sync_context_p)->bstr, &log_levelP, &protoP, line_numP,
        filename_length, thread_ctxt, &cur_time, short_source_fileP);
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
    rv = bvcformata((*sync_context_p)->bstr, 4096, format, args);  // big number
    (*sync_context_p)->log_level = g_oai_log.log_level2syslog[log_levelP];
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
  } else {
    async_context_p = (shared_log_queue_item_t**) contextP;
    rv              = append_log_ctx_info_prefix_id(
        prefix_id, (*async_context_p)->bstr, &log_levelP, &protoP, line_numP,
        filename_length, thread_ctxt, &cur_time, short_source_fileP);
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
    rv =
        bvcformata((*async_context_p)->bstr, 4096, format, args);  // big number
    (*async_context_p)->log.log_level = g_oai_log.log_level2syslog[log_levelP];
    if (BSTR_ERR == rv) {
      OAI_FPRINTF_ERR(
          "Error while logging LOG message : %s",
          &g_oai_log.log_proto2str[protoP][0]);
      goto error_event;
    }
  }
  return;

error_event:
  if (!(g_oai_log.is_async)) {
    LOG_FREE_ITEM(sync_context_p);
  } else if (async_context_p == NULL) {
    return;
  } else {
    LOG_FREE_ITEM_ASYNC(*async_context_p);
  }
}

int append_log_ctx_info(
    bstring bstr, const log_level_t* log_levelP, const log_proto_t* protoP,
    const unsigned int line_numP, size_t filename_length,
    const log_thread_ctxt_t* thread_ctxt, time_t* cur_time,
    const char* short_source_fileP) {
  int rv;
  char time_str[MAX_TIME_STR_LEN];
  log_get_readable_cur_time(cur_time, time_str);
  rv = bformata(
      bstr, LOG_CTXT_INFO_FMT,
      __sync_fetch_and_add(&g_oai_log.log_message_number, 1), time_str,
      thread_ctxt->tid, LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
      LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
      &g_oai_log.log_level2str[(*log_levelP)][0],
      LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH, LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH,
      &g_oai_log.log_proto2str[(*protoP)][0], LOG_DISPLAYED_FILENAME_MAX_LENGTH,
      LOG_DISPLAYED_FILENAME_MAX_LENGTH, &short_source_fileP[filename_length],
      line_numP, thread_ctxt->indent, " ");
  return rv;
}

int append_log_ctx_info_prefix_id(
    const uint64_t prefix_id, bstring bstr, const log_level_t* log_levelP,
    const log_proto_t* protoP, const unsigned int line_numP,
    size_t filename_length, const log_thread_ctxt_t* thread_ctxt,
    time_t* cur_time, const char* short_source_fileP) {
  int rv;
  char time_str[MAX_TIME_STR_LEN];
  log_get_readable_cur_time(cur_time, time_str);
  rv = bformata(
      bstr, LOG_CTXT_INFO_ID_FMT,
      __sync_fetch_and_add(&g_oai_log.log_message_number, 1), time_str,
      thread_ctxt->tid, LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
      LOG_DISPLAYED_LOG_LEVEL_NAME_MAX_LENGTH,
      &g_oai_log.log_level2str[(*log_levelP)][0],
      LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH, LOG_DISPLAYED_PROTO_NAME_MAX_LENGTH,
      &g_oai_log.log_proto2str[(*protoP)][0], LOG_DISPLAYED_FILENAME_MAX_LENGTH,
      LOG_DISPLAYED_FILENAME_MAX_LENGTH, &short_source_fileP[filename_length],
      line_numP, prefix_id, thread_ctxt->indent, " ");
  return rv;
}

//------------------------------------------------------------------------------
// Get the short source file name to print in the log line
// Chop off the prefix string appearing before ROOT (Ex: /oai/) and print
// Ex:
//    input: /home/vagrant/magma/lte/gateway/c/oai/tasks/nas/emm/sap/emm_cn.c
//           Assume root is /oai/
//    output: tasks/nas/emm/sap/emm_cn.c
const char* get_short_file_name(const char* const source_file_nameP) {
  if (!source_file_nameP) return source_file_nameP;

  char* root_startP = strstr(source_file_nameP, LOG_MAGMA_REPO_ROOT);

  if (!root_startP) return source_file_nameP;  // root pattern not found

  return root_startP + strlen(LOG_MAGMA_REPO_ROOT);
}

// Return the hex representation of a char array

char* bytes_to_hex(char* byte_array, int length, char* hex_array) {
  int i;
  for (i = 0; i < length; i++) {
    sprintf(hex_array + i * 3, " %02x", (unsigned char) byte_array[i]);
  }
  hex_array[3 * length + 1] = '\0';
  return hex_array;
}
