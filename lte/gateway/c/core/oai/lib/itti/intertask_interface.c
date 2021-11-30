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

#define _GNU_SOURCE
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>
#include <malloc.h>
#include <stdint.h>
#include <liblfds710.h>

#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"

/* Includes "intertask_interface_init.h" to check prototype coherence, but
   disable threads and messages information generation.
*/
#define CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_init.h"

#undef CHECK_PROTOTYPE_ONLY

#include "lte/gateway/c/core/oai/lib/itti/signals.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"

/* ITTI DEBUG groups */
#define ITTI_DEBUG_POLL (1 << 0)
#define ITTI_DEBUG_SEND (1 << 1)
#define ITTI_DEBUG_EVEN_FD (1 << 2)
#define ITTI_DEBUG_INIT (1 << 3)
#define ITTI_DEBUG_EXIT (1 << 4)
#define ITTI_DEBUG_ISSUES (1 << 5)
#define ITTI_DEBUG_MP_STATISTICS (1 << 6)

const int itti_debug = ITTI_DEBUG_ISSUES | ITTI_DEBUG_MP_STATISTICS;

#define ITTI_DEBUG(m, x, args...)                                              \
  do {                                                                         \
    /* stdout is redirected to syslog when MME is run via systemd */           \
    if ((m) &itti_debug) fprintf(stdout, "[ITTI][D]" x, ##args);               \
  } while (0);

/* Global message size */
#define MESSAGE_SIZE(mESSAGEiD)                                                \
  (sizeof(MessageHeader) + itti_desc.messages_info[mESSAGEiD].size)

#define likely(x) __builtin_expect(!!(x), 1)
#define unlikely(x) __builtin_expect(!!(x), 0)

typedef volatile enum task_state_s {
  TASK_STATE_NOT_CONFIGURED,
  TASK_STATE_STARTING,
  TASK_STATE_READY,
  TASK_STATE_ENDED,
  TASK_STATE_MAX,
} task_state_t;

typedef struct thread_desc_s {
  pthread_t task_thread;  // pthread associated with the thread

  volatile task_state_t task_state;  // State of the thread

} thread_desc_t;

typedef struct itti_desc_s {
  thread_desc_t* threads;
  thread_id_t thread_max;
  task_id_t task_max;
  MessagesIds messages_id_max;

  bool thread_handling_signals;
  pthread_t thread_ref;

  const task_info_t* tasks_info;
  const message_info_t* messages_info;

  int running;

  volatile uint32_t created_tasks;
  volatile uint32_t ready_tasks;
} itti_desc_t;

static itti_desc_t itti_desc;

status_code_e send_msg_to_task(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message) {
  if (likely(task_zmq_ctx_p->ready)) {
    AssertFatal(
        task_zmq_ctx_p->push_socks[destination_task_id],
        "Sending to task without push socket. id: %s to %s!\n",
        itti_get_message_name(message->ittiMsgHeader.messageId),
        itti_get_task_name(destination_task_id));

    // TODO: can we use zframe_frommem to avoid memcopy
    zframe_t* frame = zframe_new(
        message, sizeof(MessageHeader) + message->ittiMsgHeader.ittiMsgSize);
    assert(frame);

    // Protect against multiple threads using this context
    pthread_mutex_lock(&task_zmq_ctx_p->send_mutex);
    int rc =
        zframe_send(&frame, task_zmq_ctx_p->push_socks[destination_task_id], 0);
    assert(rc == 0);
    pthread_mutex_unlock(&task_zmq_ctx_p->send_mutex);
  } else {
    ITTI_DEBUG(
        ITTI_DEBUG_SEND, "Sending msg using uninitialized context. %s to %s!\n",
        itti_get_message_name(message->ittiMsgHeader.messageId),
        itti_get_task_name(destination_task_id));
  }

  free(message);
  return RETURNok;
}

MessageDef* receive_msg(zsock_t* reader) {
  zframe_t* msg_frame = zframe_recv(reader);
  assert(msg_frame);

  // Copy message to avoid memory alignment problems
  MessageDef* msg = (MessageDef*) malloc(zframe_size(msg_frame));
  AssertFatal(msg != NULL, "Message memory allocation failed!\n");
  memcpy(msg, zframe_data(msg_frame), zframe_size(msg_frame));

  zframe_destroy(&msg_frame);
  return msg;
}

void send_broadcast_msg(task_zmq_ctx_t* task_zmq_ctx_p, MessageDef* message) {
  zframe_t* frame = zframe_new(
      message, sizeof(MessageHeader) + message->ittiMsgHeader.ittiMsgSize);
  assert(frame);

  for (int i = 0; i < TASK_MAX; i++) {
    if (task_zmq_ctx_p->push_socks[i]) {
      // Reuse the same frame
      int rc = zframe_send(&frame, task_zmq_ctx_p->push_socks[i], ZFRAME_REUSE);
      assert(rc == 0);
    }
  }

  // Destroy frame as zframe_send did not destroy it because of ZFRAME_REUSE
  zframe_destroy(&frame);
  free(message);
}

int start_timer(
    task_zmq_ctx_t* task_zmq_ctx_p, size_t msec, timer_repeat_t repeat,
    zloop_timer_fn handler, void* arg) {
  int timer_id = zloop_timer(
      task_zmq_ctx_p->event_loop, msec, repeat == TIMER_REPEAT_FOREVER ? 0 : 1,
      handler, arg);

  AssertFatal(
      timer_id != -1, "Error starting timer for Task: %s\n",
      itti_get_task_name(task_zmq_ctx_p->task_id));

  return timer_id;
}

void stop_timer(task_zmq_ctx_t* task_zmq_ctx_p, int timer_id) {
  zloop_timer_end(task_zmq_ctx_p->event_loop, timer_id);
}

void init_task_context(
    task_id_t task_id, const task_id_t* remote_task_ids,
    uint8_t remote_tasks_count, zloop_reader_fn msg_handler,
    task_zmq_ctx_t* task_zmq_ctx_p) {
  task_zmq_ctx_p->task_id = task_id;

  task_zmq_ctx_p->event_loop = zloop_new();
  assert(task_zmq_ctx_p->event_loop);

  pthread_mutex_init(&task_zmq_ctx_p->send_mutex, NULL);

  for (int i = 0; i < remote_tasks_count; i++) {
    task_zmq_ctx_p->push_socks[remote_task_ids[i]] =
        zsock_new_push(itti_desc.tasks_info[remote_task_ids[i]].uri);
    AssertFatal(
        task_zmq_ctx_p->push_socks[remote_task_ids[i]],
        "remote task id: %d uri: %s", remote_task_ids[i],
        itti_desc.tasks_info[remote_task_ids[i]].uri);
  }

  if (msg_handler) {
    task_zmq_ctx_p->pull_sock =
        zsock_new_pull(itti_desc.tasks_info[task_id].uri);
    AssertFatal(
        task_zmq_ctx_p->pull_sock, "task id: %d uri: %s", task_id,
        itti_desc.tasks_info[task_id].uri);

    int rc = zloop_reader(
        task_zmq_ctx_p->event_loop, task_zmq_ctx_p->pull_sock, msg_handler,
        NULL);
    assert(rc == 0);
  }

  task_zmq_ctx_p->ready = true;
}

void destroy_task_context(task_zmq_ctx_t* task_zmq_ctx_p) {
  task_zmq_ctx_p->ready = false;
  zloop_destroy(&task_zmq_ctx_p->event_loop);
  zsock_destroy(&task_zmq_ctx_p->pull_sock);
  for (int i = 0; i < TASK_MAX; i++) {
    if (task_zmq_ctx_p->push_socks[i]) {
      zsock_destroy(&task_zmq_ctx_p->push_socks[i]);
    }
  }
}

const char* itti_get_message_name(MessagesIds message_id) {
  AssertFatal(
      message_id < itti_desc.messages_id_max,
      "Message id (%d) is out of range (%d)!\n", message_id,
      itti_desc.messages_id_max);
  return (itti_desc.messages_info[message_id].name);
}

const char* itti_get_task_name(task_id_t task_id) {
  if (itti_desc.task_max > 0) {
    AssertFatal(
        task_id < itti_desc.task_max, "Task id (%d) is out of range (%d)!\n",
        task_id, itti_desc.task_max);
  } else {
    return ("ITTI NOT INITIALIZED !!!");
  }

  return (itti_desc.tasks_info[task_id].name);
}

static task_id_t itti_get_current_task_id(void) {
  task_id_t task_id;
  thread_id_t thread_id;
  pthread_t thread = pthread_self();

  for (task_id = TASK_FIRST; task_id < itti_desc.task_max; task_id++) {
    thread_id = TASK_GET_THREAD_ID(task_id);

    if (itti_desc.threads[thread_id].task_thread == thread) {
      return task_id;
    }
  }

  return TASK_UNKNOWN;
}

static MessageDef* itti_alloc_new_message_sized(
    task_id_t origin_task_id, MessagesIds message_id, MessageHeaderSize size) {
  MessageDef* new_msg = NULL;

  AssertFatal(
      message_id < itti_desc.messages_id_max,
      "Message id (%d) is out of range (%d)!\n", message_id,
      itti_desc.messages_id_max);

  if (origin_task_id == TASK_UNKNOWN) {
    origin_task_id =
        itti_get_current_task_id();  // Try to identify real origin task ID
  }

  new_msg = (MessageDef*) malloc(sizeof(MessageHeader) + size);
  AssertFatal(new_msg != NULL, "Message memory allocation failed!\n");

  // better to do it here than in client code
  memset(&new_msg->ittiMsg, 0, size);

  new_msg->ittiMsgHeader.messageId    = message_id;
  new_msg->ittiMsgHeader.originTaskId = origin_task_id;
  new_msg->ittiMsgHeader.ittiMsgSize  = size;
  new_msg->ittiMsgHeader.imsi         = 0;
  clock_gettime(CLOCK_MONOTONIC_RAW, &new_msg->ittiMsgHeader.timestamp);

  return new_msg;
}

MessageDef* itti_alloc_new_message(
    task_id_t origin_task_id, MessagesIds message_id) {
  return itti_alloc_new_message_sized(
      origin_task_id, message_id, itti_desc.messages_info[message_id].size);
}

MessageDef* DEPRECATEDitti_alloc_new_message_fatal(
    task_id_t origin_task_id, MessagesIds message_id) {
  MessageDef* message_p = itti_alloc_new_message_sized(
      origin_task_id, message_id, itti_desc.messages_info[message_id].size);
  AssertFatal(message_p, "DEPRECATEDitti_alloc_new_message_fatal Failed");
  return message_p;
}

status_code_e itti_create_task(
    task_id_t task_id, void* (*start_routine)(void*), void* args_p) {
  thread_id_t thread_id = TASK_GET_THREAD_ID(task_id);

  AssertFatal(start_routine != NULL, "Start routine is NULL!\n");
  AssertFatal(
      thread_id < itti_desc.thread_max,
      "Thread id (%d) is out of range (%d)!\n", thread_id,
      itti_desc.thread_max);
  AssertFatal(
      itti_desc.threads[thread_id].task_state == TASK_STATE_NOT_CONFIGURED,
      "Task %d, thread %d state is not correct (%d)!\n", task_id, thread_id,
      itti_desc.threads[thread_id].task_state);

  itti_desc.threads[thread_id].task_state = TASK_STATE_STARTING;

  ITTI_DEBUG(
      ITTI_DEBUG_INIT, " Creating thread for task %s ...\n",
      itti_get_task_name(task_id));

  int result = pthread_create(
      &itti_desc.threads[thread_id].task_thread, NULL, start_routine, args_p);

  AssertFatal(
      result >= 0, "Thread creation for task %d, thread %d failed (%d)!\n",
      task_id, thread_id, result);

  pthread_setname_np(
      itti_desc.threads[thread_id].task_thread, itti_get_task_name(task_id));
  itti_desc.created_tasks++;

  // Wait till the thread is completely ready

  while (itti_desc.threads[thread_id].task_state != TASK_STATE_READY)
    usleep(1000);

  return RETURNok;
}

void itti_mark_task_ready(task_id_t task_id) {
  thread_id_t thread_id = TASK_GET_THREAD_ID(task_id);

  AssertFatal(
      thread_id < itti_desc.thread_max,
      "Thread id (%d) is out of range (%d)!\n", thread_id,
      itti_desc.thread_max);

  // Mark the thread as using LFDS queue

  LFDS710_MISC_MAKE_VALID_ON_CURRENT_LOGICAL_CORE_INITS_COMPLETED_BEFORE_NOW_ON_ANY_OTHER_LOGICAL_CORE;
  itti_desc.threads[thread_id].task_state = TASK_STATE_READY;
  itti_desc.ready_tasks++;

  ITTI_DEBUG(
      ITTI_DEBUG_INIT, " task %s started\n", itti_get_task_name(task_id));
}

void itti_exit_task(void) {
  pthread_exit(NULL);
}

int itti_init(
    task_id_t task_max, thread_id_t thread_max, MessagesIds messages_id_max,
    const task_info_t* tasks_info, const message_info_t* messages_info,
    const char* const messages_definition_xml,
    const char* const dump_file_name) {
  thread_id_t thread_id;

  ITTI_DEBUG(
      ITTI_DEBUG_INIT, " Init: %d tasks, %d threads, %d messages\n", task_max,
      thread_max, messages_id_max);
  CHECK_INIT_RETURN(signal_mask());

  // This assert make sure \ref ittiMsg directly following \ref ittiMsgHeader.
  // See \ref MessageDef definition for details.
  assert(sizeof(MessageHeader) == offsetof(MessageDef, ittiMsg));
  // Saves threads and messages max values

  itti_desc.task_max                = task_max;
  itti_desc.thread_max              = thread_max;
  itti_desc.messages_id_max         = messages_id_max;
  itti_desc.thread_handling_signals = false;
  itti_desc.tasks_info              = tasks_info;
  itti_desc.messages_info           = messages_info;

  // Allocates memory for threads info
  itti_desc.threads = calloc(itti_desc.thread_max, sizeof(thread_desc_t));

  // Initializing each thread

  for (thread_id = THREAD_FIRST; thread_id < itti_desc.thread_max;
       thread_id++) {
    itti_desc.threads[thread_id].task_state = TASK_STATE_NOT_CONFIGURED;
  }

  itti_desc.running       = 1;
  itti_desc.created_tasks = 0;
  itti_desc.ready_tasks   = 0;

  return 0;
}

imsi64_t itti_get_associated_imsi(MessageDef* msg) {
  return msg != NULL ? msg->ittiMsgHeader.imsi : 0;
}

void itti_wait_tasks_end(task_zmq_ctx_t* task_ctx) {
  int end = 0;
  int thread_id;
  task_id_t task_id;
  int ready_tasks;
  int result;
  int retries = 10;

  itti_desc.thread_handling_signals = true;
  itti_desc.thread_ref              = pthread_self();

  // Handle signals here

  while (end == 0) {
    signal_handle(&end, task_ctx);
  }

  ITTI_DEBUG(ITTI_DEBUG_EXIT, "Closing all tasks");
  sleep(1);

  do {
    ready_tasks = 0;
    task_id     = TASK_FIRST;

    for (thread_id = THREAD_FIRST; thread_id < itti_desc.thread_max;
         thread_id++) {
      // Skip tasks which are not running

      if (itti_desc.threads[thread_id].task_state == TASK_STATE_READY) {
        while (thread_id != TASK_GET_THREAD_ID(task_id)) {
          task_id++;
        }

        result =
            pthread_tryjoin_np(itti_desc.threads[thread_id].task_thread, NULL);
        ITTI_DEBUG(
            ITTI_DEBUG_EXIT, " Thread %s join status %d\n",
            itti_get_task_name(task_id), result);

        if (result == 0) {
          // Thread has terminated

          itti_desc.threads[thread_id].task_state = TASK_STATE_ENDED;
        } else {
          // Thread is still running, count it

          ready_tasks++;
        }
      }
    }

    if (ready_tasks > 0) {
      usleep(100 * 1000);
    }
  } while ((ready_tasks > 0) && (retries--));

  ITTI_DEBUG(ITTI_DEBUG_EXIT, "ready_tasks %d", ready_tasks);
  itti_desc.running = 0;

  free_wrapper((void**) &itti_desc.threads);

  if (ready_tasks > 0) {
    ITTI_DEBUG(
        ITTI_DEBUG_ISSUES, " Some threads are still running, force exit\n");
    return;
  }
}

void itti_free_desc_threads() {
  free_wrapper((void**) &itti_desc.threads);
  return;
}

void send_terminate_message_fatal(task_zmq_ctx_t* task_zmq_ctx) {
  MessageDef* terminate_message_p;

  terminate_message_p = DEPRECATEDitti_alloc_new_message_fatal(
      task_zmq_ctx->task_id, TERMINATE_MESSAGE);
  send_broadcast_msg(task_zmq_ctx, terminate_message_p);
}

long itti_get_message_latency(struct timespec timestamp) {
  struct timespec current_time;
  clock_gettime(CLOCK_MONOTONIC_RAW, &current_time);
  return (
      1000000 * (current_time.tv_sec - timestamp.tv_sec) +
      (current_time.tv_nsec - timestamp.tv_nsec) / 1000);
}
