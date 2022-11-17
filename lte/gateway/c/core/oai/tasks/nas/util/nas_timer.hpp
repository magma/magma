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
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*****************************************************************************
Source      nas_timer.hpp

Version     0.1

Date        2012/11/22

Product     NAS stack

Subsystem   Utilities

Author      Frederic Maurel

Description Timer utilities

*****************************************************************************/
#pragma once

#include <czmq.h>

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * Timer identifier returned when in inactive state (timer is stopped or has
 * failed to be started)
 */
#define NAS_TIMER_INACTIVE_ID (-1)
typedef int (*time_out_t)(zloop_t* loop, int timer_id, void* args);

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/* Timer structure */
typedef struct nas_timer_s {
  long int id;   /* The timer identifier                     */
  uint32_t msec; /* The timer interval value in milliseconds */
} nas_timer_t;

typedef struct timer_arg_s {
  mme_ue_s1ap_id_t ue_id;
  ebi_t ebi;
} timer_arg_t;

/* Type of the callback executed when the timer expired */
typedef void (*nas_timer_callback_t)(void*, imsi64_t* imsi64);

typedef struct nas_itti_timer_arg_s {
  nas_timer_callback_t nas_timer_callback;
  void* nas_timer_callback_arg;
} nas_itti_timer_arg_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
status_code_e nas_timer_init(void);
void nas_timer_cleanup(void);

void nas_timer_start(nas_timer_t* const timer, time_out_t time_out_cb,
                     timer_arg_t* time_out_cb_args);
void nas_timer_stop(nas_timer_t* const timer);
