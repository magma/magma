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

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include <stdint.h>
#include <signal.h>
#include <time.h>
#include <errno.h>
#include <sys/time.h>

#include "intertask_interface.h"
#include "timer.h"
#include "log.h"
#include "queue.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "timer_messages_types.h"

int timer_handle_signal(siginfo_t* info, task_zmq_ctx_t* task_ctx);

struct timer_elm_s {
  task_id_t task_id;  ///< Task ID which has requested the timer
  int32_t instance;   ///< Instance of the task which has requested the timer
  timer_t timer;      ///< Unique timer id
  timer_type_t type;  ///< Timer type
  void*
      timer_arg;  ///< Optional argument that will be passed when timer expires
  STAILQ_ENTRY(timer_elm_s) entries;  ///< Pointer to next element
};

typedef struct timer_desc_s {
  STAILQ_HEAD(timer_list_head, timer_elm_s) timer_queue;
  pthread_mutex_t timer_list_mutex;
  struct timespec timeout;
} timer_desc_t;

static timer_desc_t timer_desc;

static int timer_delete_helper(struct timer_elm_s* timer_p);
static struct timer_elm_s* find_timer(long timer_id);

#define TIMER_SEARCH(vAR, tIMERfIELD, tIMERvALUE, tIMERqUEUE)                  \
  do {                                                                         \
    STAILQ_FOREACH(vAR, tIMERqUEUE, entries) {                                 \
      if (((vAR)->tIMERfIELD == tIMERvALUE)) break;                            \
    }                                                                          \
  } while (0)

int timer_handle_signal(siginfo_t* info, task_zmq_ctx_t* task_ctx) {
  struct timer_elm_s* timer_p;
  MessageDef* message_p;
  timer_has_expired_t* timer_expired_p;
  task_id_t task_id;

  /*
   * Get back pointer to timer list element
   */
  timer_p = (struct timer_elm_s*) info->si_ptr;
  // LG: To many traces for msc timer:
  // TMR_DEBUG("Timer with id 0x%lx has expired", (long)timer_p->timer);
  task_id         = timer_p->task_id;
  message_p       = itti_alloc_new_message(TASK_MAIN, TIMER_HAS_EXPIRED);
  timer_expired_p = &message_p->ittiMsg.timer_has_expired;
  timer_expired_p->timer_id = (long) timer_p->timer;
  timer_expired_p->arg      = timer_p->timer_arg;

  /*
   * Notify task of timer expiry
   */
  if (send_msg_to_task(task_ctx, task_id, message_p) < 0) {
    OAILOG_DEBUG(
        LOG_ITTI, "Failed to send msg TIMER_HAS_EXPIRED to task %u\n", task_id);
    return -1;
  }

  return 0;
}

int timer_setup(
    uint32_t interval_sec, uint32_t interval_us, task_id_t task_id,
    int32_t instance, timer_type_t type, void* timer_arg, size_t arg_size,
    long* timer_id) {
  struct sigevent se;
  struct itimerspec its;
  struct timer_elm_s* timer_p;
  timer_t timer;

  if (timer_id == NULL) {
    return -1;
  }

  AssertFatal(
      type < TIMER_TYPE_MAX, "Invalid timer type (%d/%d)!\n", type,
      TIMER_TYPE_MAX);
  /*
   * Allocate new timer list element
   */
  timer_p = calloc(1, sizeof(struct timer_elm_s));

  if (timer_p == NULL) {
    OAILOG_ERROR(LOG_ITTI, "Failed to create new timer element\n");
    return -1;
  }

  memset(&timer, 0, sizeof(timer_t));
  memset(&se, 0, sizeof(struct sigevent));
  timer_p->task_id  = task_id;
  timer_p->instance = instance;
  timer_p->type     = type;
  // copy timer_arg if it exists
  if (timer_arg != NULL) {
    void* arg_copy = calloc(1, arg_size);
    if (arg_copy == NULL) {
      OAILOG_ERROR(LOG_ITTI, "Failed to copy timer argument\n");
      free_wrapper((void**) &timer_p);
      return -1;
    }
    memcpy(arg_copy, timer_arg, arg_size);
    timer_p->timer_arg = arg_copy;
  }

  /*
   * Setting up alarm
   */
  /*
   * Set and enable alarm
   */
  se.sigev_notify          = SIGEV_SIGNAL;
  se.sigev_signo           = SIGTIMER;
  se.sigev_value.sival_ptr = timer_p;

  /*
   * At the timer creation, the timer structure will be filled in with timer_id,
   * * * which is unique for this process. This id is allocated by kernel and
   * the
   * * * value might be used to distinguish timers.
   */
  if (timer_create(CLOCK_REALTIME, &se, &timer) < 0) {
    OAILOG_ERROR(
        LOG_ITTI, "Failed to create timer: (%s:%d)\n", strerror(errno), errno);
    free_wrapper((void**) &timer_p);
    return -1;
  }

  /*
   * Fill in the first expiration value.
   */
  its.it_value.tv_sec  = interval_sec;
  its.it_value.tv_nsec = interval_us * 1000;

  if (type == TIMER_PERIODIC) {
    /*
     * Asked for periodic timer. We set the interval time
     */
    its.it_interval.tv_sec  = interval_sec;
    its.it_interval.tv_nsec = interval_us * 1000;
  } else {
    /*
     * Asked for one-shot timer. Do not set the interval field
     */
    its.it_interval.tv_sec  = 0;
    its.it_interval.tv_nsec = 0;
  }

  timer_settime(timer, 0, &its, NULL);
  /*
   * Simply set the timer_id argument. so it can be used by caller
   */
  *timer_id = (long) timer;
  OAILOG_INFO(
      LOG_ITTI,
      "Requesting new %s timer with id 0x%lx that expires within "
      "%d sec and %d usec\n",
      type == TIMER_PERIODIC ? "periodic" : "single shot", *timer_id,
      interval_sec, interval_us);
  timer_p->timer = timer;
  /*
   * Lock the queue and insert the timer at the tail
   */
  pthread_mutex_lock(&timer_desc.timer_list_mutex);
  STAILQ_INSERT_TAIL(&timer_desc.timer_queue, timer_p, entries);
  pthread_mutex_unlock(&timer_desc.timer_list_mutex);
  return 0;
}

// Helper function to delete a timer from queue and cleanup associated resources
static int timer_delete_helper(struct timer_elm_s* timer_p) {
  int rc = TIMER_OK;
  pthread_mutex_lock(&timer_desc.timer_list_mutex);
  STAILQ_REMOVE(&timer_desc.timer_queue, timer_p, timer_elm_s, entries);
  pthread_mutex_unlock(&timer_desc.timer_list_mutex);

  if (timer_delete(timer_p->timer) < 0) {
    OAILOG_ERROR(
        LOG_ITTI, "Failed to delete timer 0x%lx\n", (long) timer_p->timer);
    rc = TIMER_ERR;
  }

  free_wrapper(&timer_p->timer_arg);
  free_wrapper((void**) &timer_p);
  timer_p = NULL;
  return rc;
}

// Helper function to find a timer in the queue
static struct timer_elm_s* find_timer(long timer_id) {
  struct timer_elm_s* timer_p = NULL;
  pthread_mutex_lock(&timer_desc.timer_list_mutex);
  TIMER_SEARCH(timer_p, timer, ((timer_t) timer_id), &timer_desc.timer_queue);

  if (timer_p == NULL) {
    OAILOG_ERROR(LOG_ITTI, "Didn't find timer 0x%lx in list\n", timer_id);
  }
  pthread_mutex_unlock(&timer_desc.timer_list_mutex);
  return timer_p;
}

/**
 * Called when another actor gets a message that a timer has expired.
 * If the timer is a one shot timer, then the timer is removed. If it is
 * periodic, then nothing is done
 */
int timer_handle_expired(long timer_id) {
  OAILOG_INFO(LOG_ITTI, "timer 0x%lx expired \n", timer_id);
  struct timer_elm_s* timer_p = find_timer(timer_id);
  if (timer_p == NULL) {
    return TIMER_NOT_FOUND;
  }

  if (timer_p->type == TIMER_ONE_SHOT) {
    OAILOG_INFO(
        LOG_ITTI, "Timer 0x%lx expiry signal received, deleting\n", timer_id);
    return timer_delete_helper(timer_p);
  }

  OAILOG_INFO(
      LOG_ITTI, "Timer 0x%lx expired but is not one shot, not deleting\n",
      timer_id);
  return TIMER_OK;
}

bool timer_exists(long timer_id) {
  return find_timer(timer_id) != NULL;
}

int timer_remove(long timer_id, void** arg) {
  int rc = 0;
  struct timer_elm_s* timer_p;

  OAILOG_DEBUG(LOG_ITTI, "Removing timer 0x%lx\n", timer_id);
  pthread_mutex_lock(&timer_desc.timer_list_mutex);
  TIMER_SEARCH(timer_p, timer, ((timer_t) timer_id), &timer_desc.timer_queue);

  /*
   * We didn't find the timer in list
   */
  if (timer_p == NULL) {
    pthread_mutex_unlock(&timer_desc.timer_list_mutex);
    if (arg) *arg = NULL;
    OAILOG_ERROR(LOG_ITTI, "Didn't find timer 0x%lx in list\n", timer_id);
    return -1;
  }

  STAILQ_REMOVE(&timer_desc.timer_queue, timer_p, timer_elm_s, entries);
  pthread_mutex_unlock(&timer_desc.timer_list_mutex);

  // let user of API get back arg that can be an allocated memory (memory leak).

  if (arg) *arg = timer_p->timer_arg;
  if (timer_delete(timer_p->timer) < 0) {
    OAILOG_ERROR(
        LOG_ITTI, "Failed to delete timer 0x%lx\n", (long) timer_p->timer);
    rc = -1;
  }

  free_wrapper((void**) &timer_p);
  return rc;
}

int timer_init(void) {
  OAI_FPRINTF_INFO("Initializing TIMER module\n");
  memset(&timer_desc, 0, sizeof(timer_desc_t));
  STAILQ_INIT(&timer_desc.timer_queue);
  pthread_mutex_init(&timer_desc.timer_list_mutex, NULL);
  OAI_FPRINTF_INFO("Initializing TIMER module: DONE\n");
  return 0;
}
