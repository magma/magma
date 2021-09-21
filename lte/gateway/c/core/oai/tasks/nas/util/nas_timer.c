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
  Source      nas_timer.c

  Version     0.1

  Date        2012/10/09

  Product     NAS stack

  Subsystem   Utilities

  Author      Frederic Maurel

  Description Timer utilities

*****************************************************************************/

#include <string.h>  // memset
#include <stdlib.h>  // malloc, free

#include "nas_timer.h"
#include "common_defs.h"
#include "dynamic_memory_check.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
#include "log.h"
#include "mme_app_timer.h"

//------------------------------------------------------------------------------
status_code_e nas_timer_init(void) {
  return (RETURNok);
}

//------------------------------------------------------------------------------
void nas_timer_cleanup(void) {}

//------------------------------------------------------------------------------
void nas_timer_start(
    struct nas_timer_s* const timer, time_out_t time_out_cb,
    timer_arg_t* time_out_cb_args) {
  if ((timer) && (timer->id == NAS_TIMER_INACTIVE_ID)) {
    timer->id = mme_app_start_timer_arg(
        timer->sec * 1000, TIMER_REPEAT_ONCE, time_out_cb, time_out_cb_args);
    if (NAS_TIMER_INACTIVE_ID != timer->id) {
      OAILOG_DEBUG(
          LOG_NAS_EMM, "NAS EBR Timer started UE " MME_UE_S1AP_ID_FMT "\n",
          time_out_cb_args->ue_id);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "Could not start NAS EBR Timer for UE " MME_UE_S1AP_ID_FMT " ",
          time_out_cb_args->ue_id);
    }
  }
}

//------------------------------------------------------------------------------
void nas_timer_stop(struct nas_timer_s* const timer) {
  if ((timer) && (timer->id != NAS_TIMER_INACTIVE_ID)) {
    mme_app_stop_timer(timer->id);
    timer->id = NAS_TIMER_INACTIVE_ID;
  }
}

