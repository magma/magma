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
  Source      nas_timer.c

  Version     0.1

  Date        2012/10/09

  Product     NAS stack

  Subsystem   Utilities

  Author      Frederic Maurel

  Description Timer utilities

*****************************************************************************/

#include <string.h>   // memset
#include <stdlib.h>   // malloc, free

#include "timer.h"
#include "nas_timer.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"

//------------------------------------------------------------------------------
int nas_timer_init(void)
{
  return (RETURNok);
}

//------------------------------------------------------------------------------
void nas_timer_cleanup(void) {}

//------------------------------------------------------------------------------
long int nas_timer_start(
  long sec,
  long usec,
  nas_timer_callback_t nas_timer_callback,
  void *nas_timer_callback_args)
{
  long timer_id;
  nas_itti_timer_arg_t cb;

  // do not start null timer
  if (sec == 0 && usec == 0) {
    return NAS_TIMER_INACTIVE_ID;
  }

  memset(&cb, 0, sizeof(cb));
  cb.nas_timer_callback = nas_timer_callback;
  cb.nas_timer_callback_arg = nas_timer_callback_args;

  if (
    timer_setup(
      sec,
      usec,
      TASK_MME_APP,
      INSTANCE_DEFAULT,
      TIMER_ONE_SHOT,
      &cb,
      sizeof(cb),
      &timer_id) == -1) {
    return NAS_TIMER_INACTIVE_ID;
  }

  return timer_id;
}

//------------------------------------------------------------------------------
long int nas_timer_stop(long int timer_id, void **nas_timer_callback_arg)
{
  nas_itti_timer_arg_t *nas_itti_timer_arg = NULL;
  timer_remove(timer_id, (void **) &nas_itti_timer_arg);
  if (nas_itti_timer_arg) {
    *nas_timer_callback_arg = nas_itti_timer_arg->nas_timer_callback_arg;
    free_wrapper((void **) &nas_itti_timer_arg);
  } else {
    *nas_timer_callback_arg = NULL;
  }
  return (NAS_TIMER_INACTIVE_ID);
  /*
  // TODO: Fix me: Amar
  timer_remove (timer_id);
  if (nas_timer_callback_arg) {
    *nas_timer_callback_arg = NULL;
  }
  return (NAS_TIMER_INACTIVE_ID);
	*/
}

//------------------------------------------------------------------------------
void mme_app_nas_timer_handle_signal_expiry(
  long timer_id,
  nas_itti_timer_arg_t *cb)
{
  OAILOG_FUNC_IN(LOG_NAS);
  if ((!timer_exists(timer_id)) || (cb->nas_timer_callback == NULL)) {
    OAILOG_ERROR(
      LOG_NAS, "Invalid timer id %ld \n", timer_id);
    OAILOG_FUNC_OUT(LOG_NAS);
  }
  cb->nas_timer_callback(cb->nas_timer_callback_arg);
  OAILOG_FUNC_OUT(LOG_NAS);
}
