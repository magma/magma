/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
//--C includes -----------------------------------------------------------------
extern "C" {
#include "log.h"
#include "mme_app_timer.h"
#include "conversions.h"
#include "intertask_interface.h"
#include "common_types.h"
}
#include "mme_app_timer_management.h"
//--C++ includes ---------------------------------------------------------------
#include <stdexcept>
//--Other includes -------------------------------------------------------------

extern task_zmq_ctx_t mme_app_task_zmq_ctx;

//------------------------------------------------------------------------------
int mme_app_start_timer(
    size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
    timer_arg_t id) {
  return magma::lte::MmeUeContext::Instance().StartTimer(
      msec, repeat, handler, id);
}

//------------------------------------------------------------------------------
void mme_app_stop_timer(int timer_id) {
  magma::lte::MmeUeContext::Instance().StopTimer(timer_id);
}
//------------------------------------------------------------------------------
void mme_app_resume_timer(
    struct ue_mm_context_s* const ue_mm_context_pP, time_t start_time,
    struct mme_app_timer_t* timer, zloop_timer_fn timer_expiry_handler,
    char* timer_name) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  time_t current_time = time(NULL);
  time_t lapsed_time  = current_time - start_time;
  OAILOG_DEBUG(LOG_MME_APP, "Handling :%s timer \n", timer_name);

  /* Below condition validates whether timer has expired before MME recovers
   * from restart, so MME shall handle as timer expiry
   */
  if (timer->sec <= lapsed_time) {
    timer_expiry_handler(mme_app_task_zmq_ctx.event_loop, timer->id, NULL);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  uint32_t remaining_time_in_seconds = timer->sec - lapsed_time;
  OAILOG_DEBUG(
      LOG_MME_APP,
      "Current_time :%ld %s timer start time :%ld "
      "lapsed time:%ld remaining time:%d \n",
      current_time, timer_name, start_time, lapsed_time,
      remaining_time_in_seconds);

  // Start timer only for remaining duration
  if ((timer->id = magma::lte::MmeUeContext::Instance().StartTimer(
           remaining_time_in_seconds * 1000, TIMER_REPEAT_ONCE,
           timer_expiry_handler, ue_mm_context_pP->mme_ue_s1ap_id)) == -1) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context_pP->emm_context._imsi64,
        "Failed to start %s timer for UE id "
        "" MME_UE_S1AP_ID_FMT "\n",
        timer_name, ue_mm_context_pP->mme_ue_s1ap_id);
    timer->id = MME_APP_TIMER_INACTIVE_ID;
  } else {
    OAILOG_DEBUG_UE(
        LOG_MME_APP, ue_mm_context_pP->emm_context._imsi64,
        "Started %s timer for UE id " MME_UE_S1AP_ID_FMT "\n", timer_name,
        ue_mm_context_pP->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_OUT(LOG_MME_APP);
}
//------------------------------------------------------------------------------
bool mme_app_get_timer_arg(int timer_id, timer_arg_t* arg) {
  return magma::lte::MmeUeContext::Instance().GetTimerArg(timer_id, arg);
}

namespace magma {
namespace lte {
//------------------------------------------------------------------------------
int MmeUeContext::StartTimer(
    size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
    TimerArgType arg) {
  int timer_id = -1;
  if ((timer_id = start_timer(
           &mme_app_task_zmq_ctx, msec, repeat, handler, nullptr)) != -1) {
    mme_app_timers.insert(std::pair<int, uint32_t>(timer_id, arg));
  }
  return timer_id;
}
//------------------------------------------------------------------------------
void MmeUeContext::StopTimer(int timer_id) {
  stop_timer(&mme_app_task_zmq_ctx, timer_id);
  mme_app_timers.erase(timer_id);
}
//------------------------------------------------------------------------------
bool MmeUeContext::GetTimerArg(const int timer_id, TimerArgType* arg) const {
  try {
    *arg = mme_app_timers.at(timer_id);
    return true;
  } catch (std::out_of_range& e) {
    return false;
  }
}

}  // namespace lte
}  // namespace magma
