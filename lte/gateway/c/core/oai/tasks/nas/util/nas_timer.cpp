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
  Source      nas_timer.cpp

  Version     0.1

  Date        2012/10/09

  Product     NAS stack

  Subsystem   Utilities

  Author      Frederic Maurel

  Description Timer utilities

*****************************************************************************/

#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.hpp"

#include <string.h>  // memset
#include <stdlib.h>  // malloc, free

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.hpp"

//------------------------------------------------------------------------------
status_code_e nas_timer_init(void) { return (RETURNok); }

//------------------------------------------------------------------------------
void nas_timer_cleanup(void) {}

//------------------------------------------------------------------------------
void nas_timer_start(struct nas_timer_s* const timer, time_out_t time_out_cb,
                     timer_arg_t* time_out_cb_args) {
  if ((timer) && (timer->id == NAS_TIMER_INACTIVE_ID)) {
    timer->id = mme_app_start_timer_arg(timer->msec, TIMER_REPEAT_ONCE,
                                        time_out_cb, time_out_cb_args);
    if (NAS_TIMER_INACTIVE_ID != timer->id) {
      OAILOG_DEBUG(LOG_NAS_EMM,
                   "NAS EBR Timer started UE " MME_UE_S1AP_ID_FMT "\n",
                   time_out_cb_args->ue_id);
    } else {
      OAILOG_ERROR(LOG_NAS_EMM,
                   "Could not start NAS EBR Timer for UE " MME_UE_S1AP_ID_FMT
                   " ",
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
