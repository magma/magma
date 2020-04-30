/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

#define _GNU_SOURCE
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>
#include <signal.h>
#include <sys/epoll.h>
#include <sys/eventfd.h>
#include <malloc.h>
#include <stdint.h>
#include <sys/time.h>

#include "assertions.h"
#include "intertask_interface.h"
#include "memory_pools.h"
#include "intertask_interface_conf.h"

/* Includes "intertask_interface_init.h" to check prototype coherence, but
   disable threads and messages information generation.
*/
#define CHECK_PROTOTYPE_ONLY
#include "intertask_interface_init.h"

#undef CHECK_PROTOTYPE_ONLY

#include "signals.h"
#include "timer.h"
#include "dynamic_memory_check.h"
#include "shared_ts_log.h"
#include "log.h"

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
    if ((m) &itti_debug) OAILOG_DEBUG(LOG_ITTI, x, ##args);                    \
  } while (0);

/* Global message size */
#define MESSAGE_SIZE(mESSAGEiD)                                                \
  (sizeof(MessageHeader) + itti_desc.messages_info[mESSAGEiD].size)

typedef volatile enum task_state_s {
  TASK_STATE_NOT_CONFIGURED,
  TASK_STATE_STARTING,
  TASK_STATE_READY,
  TASK_STATE_ENDED,
  TASK_STATE_MAX,
} task_state_t;

/* This list acts as a FIFO of messages received by tasks (RRC, NAS, ...) */
typedef struct message_list_s {
  MessageDef* msg; ///< Pointer to the message

  message_number_t message_number; ///< Unique message number
  uint32_t message_priority;       ///< Message priority
} message_list_t;

typedef struct thread_desc_s {
  /*
   * pthread associated with the thread
   */
  pthread_t task_thread;

  /*
   * State of the thread
   */
  volatile task_state_t task_state;
/*
   * This fd is used internally by ITTI.
   */
  int epoll_fd;

  /*
   * The thread fd
   */
  int task_event_fd;

  /*
   * Number of events to monitor
   */
  uint16_t nb_events;
  
  int epoll_nb_events;

  /*
   * Array of events monitored by the task.
   * * * By default only one fd is monitored (the one used to received messages
   * * * from other tasks).
   * * * More events can be suscribed later by the task itself.
   */
  struct epoll_event *events;

} thread_desc_t;

typedef struct task_desc_s {
  /*
   * Queue of messages belonging to the task
   */
  struct lfds710_queue_bmm_state message_queue
    __attribute__((aligned(LFDS710_PAL_ATOMIC_ISOLATION_IN_BYTES)));
  struct lfds710_queue_bmm_element* qbmme;
} task_desc_t;

typedef struct itti_desc_s {
  thread_desc_t* threads;
  task_desc_t* tasks;

  /*
   * Current message number. Incremented every call to send_msg_to_task
   */
  message_number_t message_number __attribute__((aligned(8)));

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

  memory_pools_handle_t memory_pools_handle;
} itti_desc_t;

static itti_desc_t itti_desc;

/** \brief Alloc and memset(0) a new itti message.
 * \param origin_task_id Task ID of the sending task
 * \param message_id Message ID
 * \param size size of the payload to send
 * @returns NULL in case of failure or newly allocated mesage ref
 **/
MessageDef* itti_alloc_new_message_sized(
  task_id_t origin_task_id,
  MessagesIds message_id,
  MessageHeaderSize size);

/** \brief Send a broadcast message to every task
 \param message_p Pointer to the message to send
 @returns < 0 on failure, 0 otherwise
 **/
int itti_send_broadcast_message(MessageDef* message_p);

void* itti_malloc(
  task_id_t origin_task_id,
  task_id_t destination_task_id,
  ssize_t size)
{
  void* ptr = NULL;

  ptr = memory_pools_allocate(
    itti_desc.memory_pools_handle, size, origin_task_id, destination_task_id);

  if (ptr == NULL) {
    char* statistics = memory_pools_statistics(itti_desc.memory_pools_handle);

    OAILOG_ERROR(LOG_ITTI, " Memory pools statistics:\n%s", statistics);
    free_wrapper((void**) &statistics);

    Fatal(
      "Memory allocation of %d bytes failed (%d -> %d)!\n",
      (int) size,
      origin_task_id,
      destination_task_id);
  }

  return ptr;
}

int itti_free(task_id_t task_id, void* ptr)
{
  int rc = EXIT_SUCCESS;

  //if (ptr == NULL) return;

  rc = memory_pools_free(itti_desc.memory_pools_handle, ptr, task_id);

  AssertFatal(
    rc == EXIT_SUCCESS, "Failed to free memory at %p (%d)\n", ptr, task_id);
return (rc);
}

static inline message_number_t itti_increment_message_number(void)
{
  /*
   * Atomic operation supported by GCC: returns the current message number
   * * * and then increment it by 1.
   * * * This can be done without mutex.
   */
  return __sync_fetch_and_add(&itti_desc.message_number, 1);
}

static inline uint32_t itti_get_message_priority(MessagesIds message_id)
{
  AssertFatal(
    message_id < itti_desc.messages_id_max,
    "Message id (%d) is out of range (%d)!\n",
    message_id,
    itti_desc.messages_id_max);
  return (itti_desc.messages_info[message_id].priority);
}

const char* itti_get_message_name(MessagesIds message_id)
{
  AssertFatal(
    message_id < itti_desc.messages_id_max,
    "Message id (%d) is out of range (%d)!\n",
    message_id,
    itti_desc.messages_id_max);
  return (itti_desc.messages_info[message_id].name);
}

const char* itti_get_task_name(task_id_t task_id)
{
  if (itti_desc.task_max > 0) {
    AssertFatal(
      task_id < itti_desc.task_max,
      "Task id (%d) is out of range (%d)!\n",
      task_id,
      itti_desc.task_max);
  } else {
    return ("ITTI NOT INITIALIZED !!!");
  }

  return (itti_desc.tasks_info[task_id].name);
}

static task_id_t itti_get_current_task_id(void)
{
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

int itti_send_broadcast_message(MessageDef* message_p)
{
  task_id_t destination_task_id;
  task_id_t origin_task_id;
  thread_id_t origin_thread_id;
  uint32_t thread_id;
  int ret = 0;
  int result;

  AssertFatal(message_p != NULL, "Trying to broadcast a NULL message!\n");
  origin_task_id = message_p->ittiMsgHeader.originTaskId;
  origin_thread_id = TASK_GET_THREAD_ID(origin_task_id);
  destination_task_id = TASK_FIRST;

  for (thread_id = THREAD_FIRST; thread_id < itti_desc.thread_max;
       thread_id++) {
    MessageDef* new_message_p;

    while (thread_id != TASK_GET_THREAD_ID(destination_task_id)) {
      destination_task_id++;
    }

    /*
     * Skip task that broadcast the message
     */
    if (thread_id != origin_thread_id) {
      /*
       * Skip tasks which are not running
       */
      if (itti_desc.threads[thread_id].task_state == TASK_STATE_READY) {
        size_t size =
          sizeof(MessageHeader) + message_p->ittiMsgHeader.ittiMsgSize;

        new_message_p = itti_malloc(origin_task_id, destination_task_id, size);
        AssertFatal(new_message_p != NULL, "New message allocation failed!\n");
        memcpy(new_message_p, message_p, size);
        result = itti_send_msg_to_task(
          destination_task_id, INSTANCE_DEFAULT, new_message_p);
        AssertFatal(
          result >= 0,
          "Failed to send message %d to thread %d (task %d)!\n",
          message_p->ittiMsgHeader.messageId,
          thread_id,
          destination_task_id);
      }
    }
  }

  itti_free(ITTI_MSG_ORIGIN_ID(message_p), message_p);

  return ret;
}

MessageDef* itti_alloc_new_message_sized(
  task_id_t origin_task_id,
  MessagesIds message_id,
  MessageHeaderSize size)
{
  MessageDef* new_msg = NULL;

  AssertFatal(
    message_id < itti_desc.messages_id_max,
    "Message id (%d) is out of range (%d)!\n",
    message_id,
    itti_desc.messages_id_max);

  if (origin_task_id == TASK_UNKNOWN) {
    /*
     * Try to identify real origin task ID
     */
    origin_task_id = itti_get_current_task_id();
  }

  new_msg =
    itti_malloc(origin_task_id, TASK_UNKNOWN, sizeof(MessageHeader) + size);

  // better to do it here than in client code
  memset(&new_msg->ittiMsg, 0, size);

  new_msg->ittiMsgHeader.messageId = message_id;
  new_msg->ittiMsgHeader.originTaskId = origin_task_id;
  new_msg->ittiMsgHeader.ittiMsgSize = size;

  return new_msg;
}

MessageDef* itti_alloc_new_message(
  task_id_t origin_task_id,
  MessagesIds message_id)
{
  return itti_alloc_new_message_sized(
    origin_task_id, message_id, itti_desc.messages_info[message_id].size);
}

int itti_send_msg_to_task(
  task_id_t destination_task_id,
  instance_t instance,
  MessageDef* message)
{
  thread_id_t destination_thread_id;
  task_id_t origin_task_id;
  message_list_t* new;
  uint32_t priority;
  message_number_t message_number;
  uint32_t message_id;

  AssertFatal(message != NULL, "Message is NULL!\n");
  AssertFatal(
    destination_task_id < itti_desc.task_max,
    "Destination task id (%d) is out of range (%d)\n",
    destination_task_id,
    itti_desc.task_max);
  destination_thread_id = TASK_GET_THREAD_ID(destination_task_id);
  message->ittiMsgHeader.destinationTaskId = destination_task_id;
  message->ittiMsgHeader.instance = instance;
  message_id = message->ittiMsgHeader.messageId;
  AssertFatal(
    message_id < itti_desc.messages_id_max,
    "Message id (%d) is out of range (%d)!\n",
    message_id,
    itti_desc.messages_id_max);
  origin_task_id = ITTI_MSG_ORIGIN_ID(message);
  priority = itti_get_message_priority(message_id);
  /*
   * Increment the global message number
   */
  message_number = itti_increment_message_number();

  if (destination_task_id != TASK_UNKNOWN) {
    memory_pools_set_info(
      itti_desc.memory_pools_handle, message, 1, destination_task_id);

    if (
      itti_desc.threads[destination_thread_id].task_state == TASK_STATE_ENDED) {
      ITTI_DEBUG(
        ITTI_DEBUG_ISSUES,
        " Message %s, number %lu with priority %d can not be sent from %s to "
        "queue (%u:%s), ended destination task!\n",
        itti_desc.messages_info[message_id].name,
        message_number,
        priority,
        itti_get_task_name(origin_task_id),
        destination_task_id,
        itti_get_task_name(destination_task_id));
      itti_free(
        origin_task_id,
        message); // In case of issues free the memory allocated for message
    } else {
      /*
       * We cannot send a message if the task is not running
       */
      AssertFatal(
        itti_desc.threads[destination_thread_id].task_state == TASK_STATE_READY,
        "Task %s Cannot send message %s (%d) to thread %s (%d), it is not in "
        "ready state (%d)!\n",
        itti_get_task_name(origin_task_id),
        itti_desc.messages_info[message_id].name,
        message_id,
        itti_desc.tasks_info[destination_thread_id].name,
        destination_thread_id,
        itti_desc.threads[destination_thread_id].task_state);
      /*
       * Allocate new list element
       */
      new = (message_list_t*) itti_malloc(
        origin_task_id, destination_task_id, sizeof(struct message_list_s));
      /*
       * Fill in members
       */
      new->msg = message;
      new->message_number = message_number;
      new->message_priority = priority;
      /*
       * Enqueue message in destination task queue
       */
      lfds710_queue_bmm_enqueue(
        &itti_desc.tasks[destination_task_id].message_queue, NULL, new);

      /*
        * Only use event fd for tasks, subtasks will pool the queue
        */
      if (TASK_GET_PARENT_TASK_ID(destination_task_id) == TASK_UNKNOWN) {
        ssize_t write_ret;
        eventfd_t sem_counter = 1;

        /*
          * Call to write for an event fd must be of 8 bytes
          */
        write_ret = write(
          itti_desc.threads[destination_thread_id].task_event_fd,
          &sem_counter,
          sizeof(sem_counter));
        AssertFatal(
          write_ret == sizeof(sem_counter),
          "Write to task message FD (%d) failed (%d/%d)\n",
          destination_thread_id,
          (int) write_ret,
          (int) sizeof(sem_counter));
      }

      ITTI_DEBUG(
        ITTI_DEBUG_SEND,
        " Message %s, number %lu with priority %d successfully sent from %s to "
        "queue (%u:%s)\n",
        itti_desc.messages_info[message_id].name,
        message_number,
        priority,
        itti_get_task_name(origin_task_id),
        destination_task_id,
        itti_get_task_name(destination_task_id));
    }
  } else {
    /*
     * This is a debug message to TASK_UNKNOWN, we can release safely release it
     */
    itti_free(origin_task_id, message);
  }

  return 0;
}

void itti_subscribe_event_fd(task_id_t task_id, int fd) {
  thread_id_t thread_id;
  struct epoll_event event;

  AssertFatal(task_id < itti_desc.task_max,
              "Task id (%d) is out of range (%d)!\n", task_id,
              itti_desc.task_max);
  thread_id = TASK_GET_THREAD_ID(task_id);
  itti_desc.threads[thread_id].nb_events++;
  /*
   * Reallocate the events
   */
  itti_desc.threads[thread_id].events = realloc(
      itti_desc.threads[thread_id].events,
      itti_desc.threads[thread_id].nb_events * sizeof(struct epoll_event));
  event.events = EPOLLIN | EPOLLERR;
  event.data.u64 = 0;
  event.data.fd = fd;

  /*
   * Add the event fd to the list of monitored events
   */
  if (epoll_ctl(itti_desc.threads[thread_id].epoll_fd, EPOLL_CTL_ADD, fd,
                &event) != 0) {
    /*
     * Always assert on this condition
     */
    AssertFatal(0, "epoll_ctl (EPOLL_CTL_ADD) failed for task %s, fd %d: %s!\n",
                itti_get_task_name(task_id), fd, strerror(errno));
  }

  ITTI_DEBUG(ITTI_DEBUG_EVEN_FD, " Successfully subscribed fd %d for task %s\n",
             fd, itti_get_task_name(task_id));
}

void itti_unsubscribe_event_fd(task_id_t task_id, int fd) {
  thread_id_t thread_id;

  AssertFatal(task_id < itti_desc.task_max,
              "Task id (%d) is out of range (%d)!\n", task_id,
              itti_desc.task_max);
  AssertFatal(fd >= 0, "File descriptor (%d) is invalid!\n", fd);
  thread_id = TASK_GET_THREAD_ID(task_id);

  /*
   * Add the event fd to the list of monitored events
   */
  if (epoll_ctl(itti_desc.threads[thread_id].epoll_fd, EPOLL_CTL_DEL, fd,
                NULL) != 0) {
    /*
     * Always assert on this condition
     */
    AssertFatal(0, "epoll_ctl (EPOLL_CTL_DEL) failed for task %s, fd %d: %s!\n",
                itti_get_task_name(task_id), fd, strerror(errno));
  }

  itti_desc.threads[thread_id].nb_events--;
  itti_desc.threads[thread_id].events = realloc(
      itti_desc.threads[thread_id].events,
      itti_desc.threads[thread_id].nb_events * sizeof(struct epoll_event));
}

int itti_get_events(task_id_t task_id, struct epoll_event **events) {
  thread_id_t thread_id;

  AssertFatal(task_id < itti_desc.task_max,
              "Task id (%d) is out of range (%d)\n", task_id,
              itti_desc.task_max);
  thread_id = TASK_GET_THREAD_ID(task_id);
  *events = itti_desc.threads[thread_id].events;
  return itti_desc.threads[thread_id].epoll_nb_events;
}

static inline void itti_receive_msg_internal_event_fd(
    task_id_t task_id, uint8_t polling, MessageDef **received_msg) {
  thread_id_t thread_id;
  int epoll_ret = 0;
  int epoll_timeout = 0;
  int i;

  AssertFatal(task_id < itti_desc.task_max,
              "Task id (%d) is out of range (%d)!\n", task_id,
              itti_desc.task_max);
  AssertFatal(received_msg != NULL, "Received message is NULL!\n");
  thread_id = TASK_GET_THREAD_ID(task_id);
  *received_msg = NULL;

  if (polling) {
    /*
     * In polling mode we set the timeout to 0 causing epoll_wait to return
     * * * immediately.
     */
    epoll_timeout = 0;
  } else {
    /*
     * timeout = -1 causes the epoll_wait to wait indefinitely.
     */
    epoll_timeout = -1;
  }

  do {
    epoll_ret =
        epoll_wait(itti_desc.threads[thread_id].epoll_fd,
                   itti_desc.threads[thread_id].events,
                   itti_desc.threads[thread_id].nb_events, epoll_timeout);
  } while (epoll_ret < 0 && errno == EINTR);

  if (epoll_ret < 0) {
    AssertFatal(0, "epoll_wait failed for task %s: %s!\n",
                itti_get_task_name(task_id), strerror(errno));
  }

  if (epoll_ret == 0 && polling) {
    /*
     * No data to read -> return
     */
    return;
  }

  itti_desc.threads[thread_id].epoll_nb_events = epoll_ret;

  for (i = 0; i < epoll_ret; i++) {
    /*
     * Check if there is an event for ITTI for the event fd
     */
    if ((itti_desc.threads[thread_id].events[i].events & EPOLLIN) &&
        (itti_desc.threads[thread_id].events[i].data.fd ==
         itti_desc.threads[thread_id].task_event_fd)) {
      struct message_list_s *message = NULL;
      eventfd_t sem_counter;
      ssize_t read_ret;
      int result = EXIT_SUCCESS;

      /*
       * Read will always return 1
       */
      read_ret = read(itti_desc.threads[thread_id].task_event_fd, &sem_counter,
                      sizeof(sem_counter));
      AssertFatal(read_ret == sizeof(sem_counter),
                  "Read from task message FD (%d) failed (%d/%d)!\n", thread_id,
                  (int)read_ret, (int)sizeof(sem_counter));

      if (lfds710_queue_bmm_dequeue(&itti_desc.tasks[task_id].message_queue,
                                    NULL, (void **)&message) == 0) {
        /*
         * No element in list -> this should not happen
         */
        AssertFatal(0,
                    "No message in queue for task %d while there are %d events "
                    "and some for the messages queue!\n",
                    task_id, epoll_ret);
        return;
      }

      AssertFatal(message != NULL, "Message from message queue is NULL!\n");
      *received_msg = message->msg;
      result = itti_free(ITTI_MSG_ORIGIN_ID(message->msg), message);
      AssertFatal(result == EXIT_SUCCESS, "Failed to free memory (%d)!\n",
                  result);
      /*
       * Mark that the event has been processed
       */
      itti_desc.threads[thread_id].events[i].events &= ~EPOLLIN;
      return;
    }
  }
}


void itti_receive_msg(task_id_t task_id, MessageDef** received_msg)
{
  //thread_id_t thread_id;
  //struct message_list_s* message = NULL;
  //eventfd_t sem_counter;
  //ssize_t n_read;


  AssertFatal(
    task_id < itti_desc.task_max,
    "Task id (%d) is out of range (%d)!\n",
    task_id,
    itti_desc.task_max);
  AssertFatal(received_msg != NULL, "Received message is NULL!\n");

  //thread_id = TASK_GET_THREAD_ID(task_id);
  *received_msg = NULL;

itti_receive_msg_internal_event_fd(task_id, 0, received_msg);

/*
  n_read = read(
    itti_desc.threads[thread_id].task_event_fd,
    &sem_counter,
    sizeof(sem_counter));
  AssertFatal(
    n_read == sizeof(sem_counter),
    "Read from task message FD (%d) failed (%zu/%zu)!\n",
    thread_id,
    n_read,
    sizeof(sem_counter));
//
  if (
    lfds710_queue_bmm_dequeue(
      &itti_desc.tasks[task_id].message_queue, NULL, (void**) &message) == 0) {
    OAILOG_WARNING(
      LOG_ITTI,
      "No message in queue for task %d while there are %zu and some "
      "for the messages queue!\n",
      task_id,
      sem_counter);
  }
*/
  //AssertFatal(message != NULL, "Message from message queue is NULL!\n");

  //*received_msg = message->msg;

  //itti_free(ITTI_MSG_ORIGIN_ID(message->msg), message);
}

int itti_create_task(
  task_id_t task_id,
  void* (*start_routine)(void*),
  void* args_p)
{
  thread_id_t thread_id = TASK_GET_THREAD_ID(task_id);
  int result = 0;

  AssertFatal(start_routine != NULL, "Start routine is NULL!\n");
  AssertFatal(
    thread_id < itti_desc.thread_max,
    "Thread id (%d) is out of range (%d)!\n",
    thread_id,
    itti_desc.thread_max);
  AssertFatal(
    itti_desc.threads[thread_id].task_state == TASK_STATE_NOT_CONFIGURED,
    "Task %d, thread %d state is not correct (%d)!\n",
    task_id,
    thread_id,
    itti_desc.threads[thread_id].task_state);

  itti_desc.threads[thread_id].task_state = TASK_STATE_STARTING;

  ITTI_DEBUG(
    ITTI_DEBUG_INIT,
    " Creating thread for task %s ...\n",
    itti_get_task_name(task_id));

  result = pthread_create(
    &itti_desc.threads[thread_id].task_thread, NULL, start_routine, args_p);

  AssertFatal(
    result >= 0,
    "Thread creation for task %d, thread %d failed (%d)!\n",
    task_id,
    thread_id,
    result);

  char name[16];

  snprintf(name, sizeof(name), "ITTI %d", thread_id);
  pthread_setname_np(itti_desc.threads[thread_id].task_thread, name);
  itti_desc.created_tasks++;

  /*
   * Wait till the thread is completely ready
   */
  while (itti_desc.threads[thread_id].task_state != TASK_STATE_READY)
    usleep(1000);

  return 0;
}

void itti_mark_task_ready(task_id_t task_id)
{
  thread_id_t thread_id = TASK_GET_THREAD_ID(task_id);

  AssertFatal(
    thread_id < itti_desc.thread_max,
    "Thread id (%d) is out of range (%d)!\n",
    thread_id,
    itti_desc.thread_max);

  /*
   * Mark the thread as using LFDS queue
   */
  LFDS710_MISC_MAKE_VALID_ON_CURRENT_LOGICAL_CORE_INITS_COMPLETED_BEFORE_NOW_ON_ANY_OTHER_LOGICAL_CORE;
  itti_desc.threads[thread_id].task_state = TASK_STATE_READY;
  itti_desc.ready_tasks++;

  ITTI_DEBUG(
    ITTI_DEBUG_INIT, " task %s started\n", itti_get_task_name(task_id));
}

void itti_exit_task(void)
{
  pthread_exit(NULL);
}

int itti_init(
  task_id_t task_max,
  thread_id_t thread_max,
  MessagesIds messages_id_max,
  const task_info_t* tasks_info,
  const message_info_t* messages_info,
  const char* const messages_definition_xml,
  const char* const dump_file_name)
{
  task_id_t task_id;
  thread_id_t thread_id;

  itti_desc.message_number = 1;
  ITTI_DEBUG(
    ITTI_DEBUG_INIT,
    " Init: %d tasks, %d threads, %d messages\n",
    task_max,
    thread_max,
    messages_id_max);
  CHECK_INIT_RETURN(signal_mask());
  /*
   * Saves threads and messages max values
   */
  itti_desc.task_max = task_max;
  itti_desc.thread_max = thread_max;
  itti_desc.messages_id_max = messages_id_max;
  itti_desc.thread_handling_signals = false;
  itti_desc.tasks_info = tasks_info;
  itti_desc.messages_info = messages_info;
  /*
   * Allocates memory for tasks info
   */
  itti_desc.tasks = memalign(
    LFDS710_PAL_ATOMIC_ISOLATION_IN_BYTES,
    itti_desc.task_max * sizeof(task_desc_t));
  memset(itti_desc.tasks, 0, itti_desc.task_max * sizeof(task_desc_t));
  /*
   * Allocates memory for threads info
   */
  itti_desc.threads = calloc(itti_desc.thread_max, sizeof(thread_desc_t));

  /*
   * Initializing each queue and related stuff
   */
  for (task_id = TASK_FIRST; task_id < itti_desc.task_max; task_id++) {
    ITTI_DEBUG(
      ITTI_DEBUG_INIT,
      " Initializing %stask %s%s%s\n",
      itti_desc.tasks_info[task_id].parent_task != TASK_UNKNOWN ? "sub-" : "",
      itti_desc.tasks_info[task_id].name,
      itti_desc.tasks_info[task_id].parent_task != TASK_UNKNOWN ?
        " with parent " :
        "",
      itti_desc.tasks_info[task_id].parent_task != TASK_UNKNOWN ?
        itti_get_task_name(itti_desc.tasks_info[task_id].parent_task) :
        "");
    ITTI_DEBUG(
      ITTI_DEBUG_INIT,
      " Creating queue of message of size %u\n",
      itti_desc.tasks_info[task_id].queue_size);
    printf(
      " Creating queue of message of size %u\n",
      itti_desc.tasks_info[task_id].queue_size);

    itti_desc.tasks[task_id].qbmme = calloc(
      itti_desc.tasks_info[task_id].queue_size,
      sizeof(struct lfds710_queue_bmm_element));
    lfds710_queue_bmm_init_valid_on_current_logical_core(
      &itti_desc.tasks[task_id].message_queue,
      itti_desc.tasks[task_id].qbmme,
      itti_desc.tasks_info[task_id].queue_size,
      NULL);
  }

  /*
   * Initializing each thread
   */
  for (thread_id = THREAD_FIRST; thread_id < itti_desc.thread_max;
       thread_id++) {
    itti_desc.threads[thread_id].task_state = TASK_STATE_NOT_CONFIGURED;
    itti_desc.threads[thread_id].epoll_fd = epoll_create1(0);
    if (itti_desc.threads[thread_id].epoll_fd == -1) {
      /*
       * Always assert on this condition
       */
      AssertFatal(0, "Failed to create new epoll fd: %s!\n", strerror(errno));
    }

    itti_desc.threads[thread_id].task_event_fd = eventfd(0, EFD_SEMAPHORE);
    
    if (itti_desc.threads[thread_id].task_event_fd == -1) {
      Fatal("eventfd failed: %s!\n", strerror(errno));
    }

itti_desc.threads[thread_id].nb_events = 1;
    itti_desc.threads[thread_id].events = calloc(1, sizeof(struct epoll_event));
    itti_desc.threads[thread_id].events->events = EPOLLIN | EPOLLERR;
    itti_desc.threads[thread_id].events->data.fd =
        itti_desc.threads[thread_id].task_event_fd;

    /*
     * Add the event fd to the list of monitored events
     */
    if (epoll_ctl(itti_desc.threads[thread_id].epoll_fd, EPOLL_CTL_ADD,
                  itti_desc.threads[thread_id].task_event_fd,
                  itti_desc.threads[thread_id].events) != 0) {
      /*
       * Always assert on this condition
       */
      AssertFatal(0, " epoll_ctl (EPOLL_CTL_ADD) failed: %s!\n",
                  strerror(errno));
    }


    ITTI_DEBUG(
      ITTI_DEBUG_EVEN_FD,
      " Successfully subscribed fd %d for thread %d\n",
      itti_desc.threads[thread_id].task_event_fd,
      thread_id);
  }

  itti_desc.running = 1;
  itti_desc.created_tasks = 0;
  itti_desc.ready_tasks = 0;

  itti_desc.memory_pools_handle = memory_pools_create(5);
  memory_pools_add_pool(
    itti_desc.memory_pools_handle, 1000 + ITTI_QUEUE_MAX_ELEMENTS, 50);
  memory_pools_add_pool(
    itti_desc.memory_pools_handle, 1000 + (2 * ITTI_QUEUE_MAX_ELEMENTS), 100);
  memory_pools_add_pool(itti_desc.memory_pools_handle, 10000, 1000);
  memory_pools_add_pool(itti_desc.memory_pools_handle, 400, 20050);
  memory_pools_add_pool(itti_desc.memory_pools_handle, 100, 30050);
  {
    char* statistics = memory_pools_statistics(itti_desc.memory_pools_handle);

    ITTI_DEBUG(
      ITTI_DEBUG_MP_STATISTICS, " Memory pools statistics:\n%s", statistics);
    free_wrapper((void**) &statistics);
  }

  CHECK_INIT_RETURN(timer_init());
  // Could not be launched before ITTI initialization
  shared_log_itti_connect();
  OAILOG_ITTI_CONNECT();
  return 0;
}

imsi64_t itti_get_associated_imsi(MessageDef* msg)
{
  if (msg->ittiMsgHeader.imsi == 0) {
    OAILOG_DEBUG(
      LOG_ITTI,
      "IMSI associated to msg: %d, origin task id: %d, dest task id: %d is not "
      "set",
      msg->ittiMsgHeader.messageId,
      msg->ittiMsgHeader.originTaskId,
      msg->ittiMsgHeader.destinationTaskId);
  }
  return msg->ittiMsgHeader.imsi;
}

void itti_wait_tasks_end(void)
{
  int end = 0;
  int thread_id;
  task_id_t task_id;
  int ready_tasks;
  int result;
  int retries = 10;

  itti_desc.thread_handling_signals = true;
  itti_desc.thread_ref = pthread_self();

  /*
   * Handle signals here
   */
  while (end == 0) {
    signal_handle(&end);
  }

  OAILOG_INFO(LOG_ITTI, "Closing all tasks");
  sleep(1);

  do {
    ready_tasks = 0;
    task_id = TASK_FIRST;

    for (thread_id = THREAD_FIRST; thread_id < itti_desc.thread_max;
         thread_id++) {
      /*
       * Skip tasks which are not running
       */
      if (itti_desc.threads[thread_id].task_state == TASK_STATE_READY) {
        while (thread_id != TASK_GET_THREAD_ID(task_id)) {
          task_id++;
        }

        result =
          pthread_tryjoin_np(itti_desc.threads[thread_id].task_thread, NULL);
        ITTI_DEBUG(
          ITTI_DEBUG_EXIT,
          " Thread %s join status %d\n",
          itti_get_task_name(task_id),
          result);

        if (result == 0) {
          /*
           * Thread has terminated
           */
          itti_desc.threads[thread_id].task_state = TASK_STATE_ENDED;
        } else {
          /*
           * Thread is still running, count it
           */
          ready_tasks++;
        }
      }
    }

    if (ready_tasks > 0) {
      usleep(100 * 1000);
    }
  } while ((ready_tasks > 0) && (retries--));

  OAILOG_INFO(LOG_ITTI, "ready_tasks %d", ready_tasks);
  itti_desc.running = 0;
  {
    char* statistics = memory_pools_statistics(itti_desc.memory_pools_handle);

    ITTI_DEBUG(
      ITTI_DEBUG_MP_STATISTICS, " Memory pools statistics:\n%s\n", statistics);
    free_wrapper((void**) &statistics);
  }

  for (task_id = TASK_FIRST; task_id < itti_desc.task_max; task_id++) {
    free_wrapper((void**) &itti_desc.tasks[task_id].qbmme);
  }

  free_wrapper((void**) &itti_desc.tasks);
  free_wrapper((void**) &itti_desc.threads);

  if (ready_tasks > 0) {
    ITTI_DEBUG(
      ITTI_DEBUG_ISSUES, " Some threads are still running, force exit\n");
    return;
  }
}

void itti_send_terminate_message(task_id_t task_id)
{
  MessageDef* terminate_message_p;

  terminate_message_p = itti_alloc_new_message(task_id, TERMINATE_MESSAGE);
  itti_send_broadcast_message(terminate_message_p);
}
