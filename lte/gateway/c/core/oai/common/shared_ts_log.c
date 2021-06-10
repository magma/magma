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

/*! \file shared_ts_log.c
   \brief
   \author  Lionel GAUTHIER
   \date 2016
   \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <sys/time.h>
#include <pthread.h>

#include "bstrlib.h"
#include "hashtable.h"
#include "intertask_interface.h"
#include "log.h"
#include "shared_ts_log.h"
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "intertask_interface_types.h"
#include "itti_types.h"

//-------------------------------
#define LOG_MAX_QUEUE_ELEMENTS 2048
#define LOG_MESSAGE_MIN_ALLOC_SIZE 256

#define LOG_FLUSH_PERIOD_MSEC 50
//-------------------------------

typedef unsigned long log_message_number_t;

/*! \struct  oai_shared_log_t
 * \brief Structure containing all the logging utility internal variables.
 */
typedef struct oai_shared_log_s {
  // may be good to use stream instead of file descriptor when
  // logging somewhere else of the console.

  int log_start_time_second; /*!< \brief Logging utility reference time */

  log_message_number_t
      log_message_number; /*!< \brief Counter of log message        */
  struct lfds710_queue_bmm_element* qbmme;
  struct lfds710_queue_bmm_state
      log_message_queue; /*!< \brief Thread safe log message queue */
  struct lfds710_stack_state
      log_free_message_queue; /*!< \brief Thread safe memory pool       */

  hash_table_ts_t*
      thread_context_htbl; /*!< \brief Container for log_thread_ctxt_t */

  void (*logger_callback[MAX_SH_TS_LOG_CLIENT])(shared_log_queue_item_t*);
  bool running;
} oai_shared_log_t;

static oai_shared_log_t g_shared_log = {
    0}; /*!< \brief  logging utility internal variables global var definition*/

static void shared_log_exit(void);

task_zmq_ctx_t shared_log_task_zmq_ctx;
static int timer_id = -1;

//------------------------------------------------------------------------------
int shared_log_get_start_time_sec(void) {
  return g_shared_log.log_start_time_second;
}

//------------------------------------------------------------------------------
void shared_log_reuse_item(shared_log_queue_item_t* item_p) {
#if defined(SHARED_LOG_PREALLOC_STRING_BUFFERS)
  btrunc(item_p->bstr, 0);
#else
  bdestroy_wrapper(&item_p->bstr);
#endif
  LFDS710_STACK_SET_VALUE_IN_ELEMENT(item_p->se, item_p);
  lfds710_stack_push(&g_shared_log.log_free_message_queue, &item_p->se);
}

//------------------------------------------------>-------------------------------
static shared_log_queue_item_t* create_new_log_queue_item(
    sh_ts_log_app_id_t app_id) {
  shared_log_queue_item_t* item_p = calloc(1, sizeof(shared_log_queue_item_t));
  AssertFatal((item_p), "Allocation of log container failed");
  AssertFatal(
      (app_id >= MIN_SH_TS_LOG_CLIENT), "Allocation of log container failed");
  AssertFatal(
      (app_id < MAX_SH_TS_LOG_CLIENT), "Allocation of log container failed");
  item_p->app_id = app_id;
#if defined(SHARED_LOG_PREALLOC_STRING_BUFFERS)
  item_p->bstr = bfromcstralloc(LOG_MESSAGE_MIN_ALLOC_SIZE, "");
  AssertFatal((item_p->bstr), "Allocation of buf in log container failed");
#endif
  return item_p;
}

//------------------------------------------------------------------------------
shared_log_queue_item_t* get_new_log_queue_item(sh_ts_log_app_id_t app_id) {
  shared_log_queue_item_t* item_p  = NULL;
  struct lfds710_stack_element* se = NULL;

  lfds710_stack_pop(&g_shared_log.log_free_message_queue, &se);
  if (!se) {
    shared_log_flush_messages();
    lfds710_stack_pop(&g_shared_log.log_free_message_queue, &se);
  }
  if (se) {
    item_p = LFDS710_STACK_GET_VALUE_FROM_ELEMENT(*se);

    if (!item_p) {
      item_p = create_new_log_queue_item(app_id);
      AssertFatal(item_p, "Out of memory error");
    } else {
      item_p->app_id = app_id;
#if defined(SHARED_LOG_PREALLOC_STRING_BUFFERS)
      btrunc(item_p->bstr, 0);
#endif
    }
#if !defined(SHARED_LOG_PREALLOC_STRING_BUFFERS)
    item_p->bstr = bfromcstralloc(LOG_MESSAGE_MIN_ALLOC_SIZE, "");
    AssertFatal((item_p->bstr), "Allocation of buf in log container failed");
#endif
  } else {
    OAI_FPRINTF_ERR("Could not get memory for logging\n");
  }
  return item_p;
}

//------------------------------------------------------------------------------
static int handle_timer(zloop_t* loop, int id, void* arg) {
  shared_log_flush_messages();
  return 0;
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case TERMINATE_MESSAGE: {
      free(received_message_p);
      shared_log_exit();
    } break;

    default: { } break; }

  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* shared_log_thread(__attribute__((unused)) void* args_p) {
  itti_mark_task_ready(TASK_SHARED_TS_LOG);
  init_task_context(
      TASK_SHARED_TS_LOG, (task_id_t[]){}, 0, handle_message,
      &shared_log_task_zmq_ctx);

  timer_id = start_timer(
      &shared_log_task_zmq_ctx, LOG_FLUSH_PERIOD_MSEC, TIMER_REPEAT_FOREVER,
      handle_timer, NULL);
  shared_log_start_use();

  zloop_start(shared_log_task_zmq_ctx.event_loop);
  shared_log_exit();
  return NULL;
}

//------------------------------------------------------------------------------
void shared_log_get_elapsed_time_since_start(
    struct timeval* const elapsed_time) {
  // no thread safe but do not matter a lot
  gettimeofday(elapsed_time, NULL);
  // no timersub call for fastest operations
  elapsed_time->tv_sec =
      elapsed_time->tv_sec - g_shared_log.log_start_time_second;
}

//------------------------------------------------------------------------------
int shared_log_init(const int max_threadsP) {
  shared_log_queue_item_t* item_p = NULL;
  struct timeval start_time       = {.tv_sec = 0, .tv_usec = 0};

  OAI_FPRINTF_INFO("Initializing shared logging\n");
  gettimeofday(&start_time, NULL);
  g_shared_log.log_start_time_second = start_time.tv_sec;

  g_shared_log.logger_callback[SH_TS_LOG_TXT] = log_flush_message;

  bstring b = bfromcstr("Logging thread context hashtable");
  g_shared_log.thread_context_htbl =
      hashtable_ts_create(LOG_MESSAGE_MIN_ALLOC_SIZE, NULL, free_wrapper, b);
  bdestroy_wrapper(&b);
  AssertFatal(
      NULL != g_shared_log.thread_context_htbl,
      "Could not create hashtable for Log!\n");
  g_shared_log.thread_context_htbl->log_enabled = false;

  log_thread_ctxt_t* thread_ctxt = calloc(1, sizeof(log_thread_ctxt_t));
  AssertFatal(
      NULL != thread_ctxt, "Error Could not create log thread context\n");
  pthread_t p            = pthread_self();
  thread_ctxt->tid       = p;
  hashtable_rc_t hash_rc = hashtable_ts_insert(
      g_shared_log.thread_context_htbl, (hash_key_t) p, thread_ctxt);
  if (HASH_TABLE_OK != hash_rc) {
    OAI_FPRINTF_ERR("Error Could not register log thread context\n");
    free_wrapper((void**) &thread_ctxt);
  }

  lfds710_stack_init_valid_on_current_logical_core(
      &g_shared_log.log_free_message_queue, NULL);
  g_shared_log.qbmme =
      calloc(LOG_MAX_QUEUE_ELEMENTS, sizeof(*g_shared_log.qbmme));
  lfds710_queue_bmm_init_valid_on_current_logical_core(
      &g_shared_log.log_message_queue, g_shared_log.qbmme,
      LOG_MAX_QUEUE_ELEMENTS, NULL);

  shared_log_start_use();

  for (int i = 0; i < max_threadsP * 30; i++) {
    item_p = create_new_log_queue_item(MIN_SH_TS_LOG_CLIENT);  // any logger
    LFDS710_STACK_SET_VALUE_IN_ELEMENT(item_p->se, item_p);
    lfds710_stack_push(&g_shared_log.log_free_message_queue, &item_p->se);
  }

  OAI_FPRINTF_INFO("Initializing shared logging Done\n");

  g_shared_log.running = true;

  return 0;
}

//------------------------------------------------------------------------------
void shared_log_itti_connect(void) {
  int rv = 0;
  rv     = itti_create_task(TASK_SHARED_TS_LOG, &shared_log_thread, NULL);
  AssertFatal(rv == 0, "Create task for OAI logging failed!\n");
}

//------------------------------------------------------------------------------
void shared_log_start_use(void) {
  pthread_t p            = pthread_self();
  hashtable_rc_t hash_rc = hashtable_ts_is_key_exists(
      g_shared_log.thread_context_htbl, (hash_key_t) p);
  if (HASH_TABLE_KEY_NOT_EXISTS == hash_rc) {
    LFDS710_MISC_MAKE_VALID_ON_CURRENT_LOGICAL_CORE_INITS_COMPLETED_BEFORE_NOW_ON_ANY_OTHER_LOGICAL_CORE;

    log_thread_ctxt_t* thread_ctxt = calloc(1, sizeof(log_thread_ctxt_t));
    if (thread_ctxt) {
      thread_ctxt->tid = p;
      hash_rc          = hashtable_ts_insert(
          g_shared_log.thread_context_htbl, (hash_key_t) p, thread_ctxt);
      if (HASH_TABLE_OK != hash_rc) {
        OAI_FPRINTF_ERR("Error Could not register log thread context\n");
        free_wrapper((void**) &thread_ctxt);
      }
    } else {
      OAI_FPRINTF_ERR("Error Could not create log thread context\n");
    }
  }
}

//------------------------------------------------------------------------------
void shared_log_flush_messages(void) {
  shared_log_queue_item_t* item_p = NULL;

  while (lfds710_queue_bmm_dequeue(
             &g_shared_log.log_message_queue, NULL, (void**) &item_p) == 1) {
    if ((item_p->app_id >= MIN_SH_TS_LOG_CLIENT) &&
        (item_p->app_id < MAX_SH_TS_LOG_CLIENT)) {
      (*g_shared_log.logger_callback[item_p->app_id])(item_p);
    } else {
      OAI_FPRINTF_ERR("Error bad logger identifier: %d\n", item_p->app_id);
    }
    shared_log_reuse_item(item_p);
  }
}

//------------------------------------------------------------------------------
static void shared_log_element_dequeue_cleanup_callback(
    struct lfds710_queue_bmm_state* qbmms, void* key, void* value) {
  shared_log_queue_item_t* item_p = (shared_log_queue_item_t*) value;

  if (item_p) {
    if (item_p->bstr) {
      bdestroy_wrapper(&item_p->bstr);
    }
    free_wrapper((void**) &item_p);
  }
}
//------------------------------------------------------------------------------
static void shared_log_element_pop_cleanup_callback(
    struct lfds710_stack_state* ss, struct lfds710_stack_element* se) {
  shared_log_queue_item_t* item_p = (shared_log_queue_item_t*) se->value;

  if (item_p) {
    if (item_p->bstr) {
      bdestroy_wrapper(&item_p->bstr);
    }
    free_wrapper((void**) &item_p);
  }
}
//------------------------------------------------------------------------------
static void shared_log_exit(void) {
  OAI_FPRINTF_INFO("[TRACE] Entering %s\n", __FUNCTION__);
  stop_timer(&shared_log_task_zmq_ctx, timer_id);
  destroy_task_context(&shared_log_task_zmq_ctx);
  shared_log_flush_messages();
  hashtable_ts_destroy(g_shared_log.thread_context_htbl);
  lfds710_queue_bmm_cleanup(
      &g_shared_log.log_message_queue,
      shared_log_element_dequeue_cleanup_callback);
  lfds710_stack_cleanup(
      &g_shared_log.log_free_message_queue,
      shared_log_element_pop_cleanup_callback);
  free_wrapper((void**) &g_shared_log.qbmme);
  OAI_FPRINTF_INFO("[TRACE] Leaving %s\n", __FUNCTION__);

  OAI_FPRINTF_INFO("TASK_SHARED_TS_LOG terminated\n");
  pthread_exit(NULL);
}

//------------------------------------------------------------------------------
void shared_log_item(shared_log_queue_item_t* messageP) {
  if (messageP) {
    if (g_shared_log.running) {
      shared_log_start_use();
      lfds710_queue_bmm_enqueue(
          &g_shared_log.log_message_queue, NULL, messageP);
    } else {
      if (messageP->bstr) {
        bdestroy_wrapper(&messageP->bstr);
      }
      free_wrapper((void**) &messageP);
    }
  }
}
