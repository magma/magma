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

/** @defgroup _intertask_interface_impl_ Intertask Interface Mechanisms
 * Implementation
 * @ingroup _ref_implementation_
 * @{
 */

#ifndef INTERTASK_INTERFACE_H_
#define INTERTASK_INTERFACE_H_

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <sys/types.h>
#include <czmq.h>

#include "intertask_interface_conf.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "common_defs.h"

#define ITTI_MSG_ID(mSGpTR) ((mSGpTR)->ittiMsgHeader.messageId)
#define ITTI_MSG_ORIGIN_ID(mSGpTR) ((mSGpTR)->ittiMsgHeader.originTaskId)
#define ITTI_MSG_DESTINATION_ID(mSGpTR)                                        \
  ((mSGpTR)->ittiMsgHeader.destinationTaskId)
#define ITTI_MSG_INSTANCE(mSGpTR) ((mSGpTR)->ittiMsgHeader.instance)
#define ITTI_MSG_NAME(mSGpTR) itti_get_message_name(ITTI_MSG_ID(mSGpTR))
#define ITTI_MSG_ORIGIN_NAME(mSGpTR)                                           \
  itti_get_task_name(ITTI_MSG_ORIGIN_ID(mSGpTR))
#define ITTI_MSG_DESTINATION_NAME(mSGpTR)                                      \
  itti_get_task_name(ITTI_MSG_DESTINATION_ID(mSGpTR))
#define ITTI_MSG_LATENCY(mSGpTR)                                               \
  itti_get_message_latency((mSGpTR)->ittiMsgHeader.timestamp)
#define ITTI_MSG_LASTHOP_LATENCY(mSGpTR)                                       \
  ((mSGpTR)->ittiMsgHeader.last_hop_latency)

/* Make the message number platform specific */
typedef unsigned long message_number_t;
#define MESSAGE_NUMBER_SIZE (sizeof(unsigned long))

typedef struct task_zmq_ctx_s {
  task_id_t task_id;
  zloop_t* event_loop;
  zsock_t* pull_sock;
  zsock_t* push_socks[TASK_MAX];
  pthread_mutex_t send_mutex;
  bool ready;
} task_zmq_ctx_t;

typedef struct message_info_s {
  MessagesIds id;
  /* Message payload size */
  MessageHeaderSize size;
  /* Printable name */
  const char* const name;
} message_info_t;

typedef struct task_info_s {
  thread_id_t thread;
  /* Printable name */
  const char* const name;
  /* Socket endpoint */
  const char* const uri;
} task_info_t;

typedef enum timer_repeat_s {
  TIMER_REPEAT_FOREVER = 0,
  TIMER_REPEAT_ONCE,
} timer_repeat_t;

/** \brief Send a message to a task
 \param task_zmq_ctx_p Pointer to task ZMQ context
 \param destination_task_id Destination task ID
 \param message Pointer to the message to send
 @returns status_code_e
 **/
status_code_e send_msg_to_task(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t destination_task_id,
    MessageDef* message);

/** \brief Receive a message from zsock
 \param reader Pointer to ZMQ socket
 @returns Pointer to the message read (caller to free)
 **/
MessageDef* receive_msg(zsock_t* reader);

/** \brief Start timer on the ZMQ loop
 \param task_zmq_ctx_p Pointer to task ZMQ context
 \param msec Timer duration in millisecond
 \param repeat Timer type
 \param handler Callback function on timer expiry
 \param arg Data to pass to handler
 @returns -1 on failure, timer ID otherwise
 **/
int start_timer(
    task_zmq_ctx_t* task_zmq_ctx_p, size_t msec, timer_repeat_t repeat,
    zloop_timer_fn handler, void* arg);

/** \brief Stop timer on the ZMQ loop
 \param task_zmq_ctx_p Pointer to task ZMQ context
 \param timer_id Timer ID
 **/
void stop_timer(task_zmq_ctx_t* task_zmq_ctx_p, int timer_id);

/** \brief Initialize task ZMQ context
 \param task_id Task ID
 \param remote_task_ids Array of destination task IDs
 \param remote_tasks_count Size of remote_task_ids
 \param msg_handler message handler for pull socket
 \param task_zmq_ctx_p Pointer to task ZMQ context
 **/
void init_task_context(
    task_id_t task_id, const task_id_t* remote_task_ids,
    uint8_t remote_tasks_count, zloop_reader_fn msg_handler,
    task_zmq_ctx_t* task_zmq_ctx_p);

/** \brief Destroy task ZMQ context
 \param task_zmq_ctx_p Pointer to task ZMQ context
 **/
void destroy_task_context(task_zmq_ctx_t* task_zmq_ctx_p);

/** \brief Send broadcast message to all push sockets part of context
 * \param task_zmq_ctx_p Pointer to task ZMQ context
 * \param message Pointer to message
 **/
void send_broadcast_msg(task_zmq_ctx_t* task_zmq_ctx_p, MessageDef* message);

/** \brief Start thread associated to the task
 * \param task_id task to start
 * \param start_routine entry point for the task
 * \param args_p Optional argument to pass to the start routine
 * @returns status_code_e
 * @note Asserts that task is created
 **/
status_code_e itti_create_task(
    task_id_t task_id, void* (*start_routine)(void*), void* args_p);

/** \brief Mark the task as in ready state
 * \param task_id task to mark as ready
 **/
void itti_mark_task_ready(task_id_t task_id);

/** \brief Exit the current task.
 **/
void itti_exit_task(void);

/** \brief Return the printable string associated with the message
 * \param message_id Id of the message
 **/
const char* itti_get_message_name(MessagesIds message_id);

/** \brief Return the printable string associated with a task id
 * \param thread_id Id of the task
 **/
const char* itti_get_task_name(task_id_t task_id);

/** \brief Alloc and memset(0) a new itti message.
 * \param origin_task_id Task ID of the sending task
 * \param message_id Message ID
 * @returns NULL in case of failure or newly allocated mesage ref
 **/
MessageDef* itti_alloc_new_message(
    task_id_t origin_task_id, MessagesIds message_id);

/**
 * \brief Returns IMSI of ITTI task
 * @param msg MessageDef struct
 * @return uint64 IMSI
 */
imsi64_t itti_get_associated_imsi(MessageDef* msg);

/** \brief handle signals and wait for all threads to join when the process
 *complete. This function should be called from the main thread after having
 *created all ITTI tasks.
 **/
void itti_wait_tasks_end(task_zmq_ctx_t* task_ctx);

/** \brief Send a termination message to all tasks.
 * \param task_id task that is broadcasting the message.
 **/
void send_terminate_message(task_zmq_ctx_t* task_zmq_ctx);

/**
 * \brief Returns the latency of the message
 * @param timestamp timespec struct
 * @return long Message Latency in micro seconds
 */
long itti_get_message_latency(struct timespec timestamp);
#endif /* INTERTASK_INTERFACE_H_ */
/* @} */
