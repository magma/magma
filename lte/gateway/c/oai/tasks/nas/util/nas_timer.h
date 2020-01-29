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

/*****************************************************************************
Source      nas_timer.h

Version     0.1

Date        2012/11/22

Product     NAS stack

Subsystem   Utilities

Author      Frederic Maurel

Description Timer utilities

*****************************************************************************/
#ifndef FILE_NAS_TIMER_SEEN
#define FILE_NAS_TIMER_SEEN

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * Timer identifier returned when in inactive state (timer is stopped or has
 * failed to be started)
 */
#define NAS_TIMER_INACTIVE_ID (-1)

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/* Timer structure */
typedef struct nas_timer_s {
  long int id; /* The timer identifier                 */
  long sec;    /* The timer interval value in seconds  */
} nas_timer_t;

/* Type of the callback executed when the timer expired */
typedef void (*nas_timer_callback_t)(void *);

typedef struct nas_itti_timer_arg_s {
  nas_timer_callback_t nas_timer_callback;
  void *nas_timer_callback_arg;
} nas_itti_timer_arg_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

int nas_timer_init(void);
void nas_timer_cleanup(void);
long int nas_timer_start(
  long sec,
  long usec,
  nas_timer_callback_t nas_timer_callback,
  void *nas_timer_callback_args);
long int nas_timer_stop(long int timer_id, void **nas_timer_callback_arg);
void mme_app_nas_timer_handle_signal_expiry(
  long timer_id,
  nas_itti_timer_arg_t *nas_itti_timer_arg);

#endif /* FILE_NAS_TIMER_SEEN */
