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

#include "mock_task.h"

#include <pthread.h>
#include <stdlib.h>

#include "assertions.h"
#include "common_defs.h"

typedef struct msg_queue_elem {
  MessageDef *msg;
  struct msg_queue_elem *next;
} msg_queue_elem_t;

// single producer - single consumer queue
typedef struct msg_queue {
  msg_queue_elem_t *head;
  msg_queue_elem_t *tail;
  pthread_mutex_t m;
  pthread_cond_t cv;
} msg_queue_t;

// mock task needed for testing and to spin up standalone tasks
struct mock_task {
  task_id_t task_id;
  int exited;
  msg_queue_t *recv_msgs;
};

// the task thread for a mock task -- loops on recving msgs and enqueueing
// in the mock task queue
void *mock_task_thread(void *arg);

mock_task_t *mock_task_new(task_id_t task_id);
MessageDef *mock_task_recv(mock_task_t *task);
int mock_task_exited(mock_task_t *task);

msg_queue_t *msg_queue_new(void);
void msg_queue_enq(msg_queue_t *msg_queue, MessageDef *msg);
MessageDef *msg_queue_deq(msg_queue_t *msg_queue);

msg_queue_elem_t *msg_queue_elem_new(MessageDef *msg);
void msg_queue_elem_free(msg_queue_elem_t *elem);

void *mock_task_thread(void *arg)
{
  mock_task_t *mock_task = arg;
  MessageDef *msg;

  itti_mark_task_ready(mock_task->task_id);

  while (1) {
    itti_receive_msg(TASK_MME_APP, &msg);

    if (ITTI_MSG_ID(msg) == TERMINATE_MESSAGE) {
      mock_task->exited = 1;
      itti_exit_task();
    } else {
      msg_queue_enq(mock_task->recv_msgs, msg);
    }
  }

  return NULL;
}

mock_task_t *mock_task_new(task_id_t task_id)
{
  mock_task_t *mock_task;

  mock_task = malloc(sizeof(*mock_task));
  if (mock_task == NULL) return NULL;

  mock_task->task_id = task_id;
  mock_task->exited = 0;
  mock_task->recv_msgs = msg_queue_new();

  if (mock_task->recv_msgs == NULL) {
    free(mock_task);
    return NULL;
  }

  if (itti_create_task(task_id, mock_task_thread, mock_task) != RETURNok) {
    free(mock_task);
    // free msg_queue
    return NULL;
  }

  return mock_task;
}

MessageDef *mock_task_recv(mock_task_t *task)
{
  return msg_queue_deq(task->recv_msgs);
}

int mock_task_exited(mock_task_t *task)
{
  return task->exited;
}

msg_queue_t *msg_queue_new(void)
{
  msg_queue_t *msg_queue;

  msg_queue = malloc(sizeof(*msg_queue));
  if (msg_queue == NULL) return NULL;

  msg_queue->head = NULL;
  msg_queue->tail = NULL;

  if (pthread_mutex_init(&msg_queue->m, NULL) != 0) {
    Fatal("failed to init mutex for msg_queue");
  }

  if (pthread_cond_init(&msg_queue->cv, NULL) != 0) {
    Fatal("failed to init condvar for msg_queue");
  }

  return msg_queue;
}

void msg_queue_enq(msg_queue_t *msg_queue, MessageDef *msg)
{
  msg_queue_elem_t *elem;

  elem = msg_queue_elem_new(msg);
  AssertFatal(elem != NULL, "failed to create new elem for msg_queue");

  pthread_mutex_lock(&msg_queue->m);

  if (msg_queue->head == NULL) {
    msg_queue->head = elem;
    msg_queue->tail = elem;
  } else {
    msg_queue->tail->next = elem;
    msg_queue->tail = elem;
  }

  pthread_cond_signal(&msg_queue->cv);
  pthread_mutex_unlock(&msg_queue->m);
}

// blocking deque -- waits until message available
MessageDef *msg_queue_deq(msg_queue_t *msg_queue)
{
  MessageDef *msg;
  msg_queue_elem_t *old_head;

  pthread_mutex_lock(&msg_queue->m);

  while (msg_queue->head == NULL)
    pthread_cond_wait(&msg_queue->cv, &msg_queue->m);

  msg = msg_queue->head->msg;
  old_head = msg_queue->head;

  if (msg_queue->head == msg_queue->tail) {
    msg_queue->head = NULL;
    msg_queue->tail = NULL;
  } else {
    msg_queue->head = msg_queue->head->next;
  }

  msg_queue_elem_free(old_head);

  pthread_mutex_unlock(&msg_queue->m);

  return msg;
}

msg_queue_elem_t *msg_queue_elem_new(MessageDef *msg)
{
  msg_queue_elem_t *elem;

  elem = malloc(sizeof(*elem));
  if (elem == NULL) return NULL;

  elem->msg = msg;
  elem->next = NULL;

  return elem;
}

void msg_queue_elem_free(msg_queue_elem_t *elem)
{
  free(elem);
}
