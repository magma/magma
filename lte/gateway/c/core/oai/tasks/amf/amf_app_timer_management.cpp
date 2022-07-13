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
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
}
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.hpp"
//--C++ includes ---------------------------------------------------------------
#include <utility>
#include <stdexcept>
//--Other includes -------------------------------------------------------------

namespace magma5g {

extern task_zmq_ctx_t amf_app_task_zmq_ctx;

//------------------------------------------------------------------------------
int amf_app_start_timer(size_t msec, timer_repeat_t repeat,
                        zloop_timer_fn handler, timer_arg_t id) {
  return magma5g::AmfUeContext::Instance().StartTimer(msec, repeat, handler,
                                                      id);
}

//------------------------------------------------------------------------------
void amf_app_stop_timer(int timer_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  magma5g::AmfUeContext::Instance().StopTimer(timer_id);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//------------------------------------------------------------------------------
bool amf_pop_timer_arg(int timer_id, timer_arg_t* arg) {
  return magma5g::AmfUeContext::Instance().PopTimerArgById(timer_id, arg);
}

//------------------------------------------------------------------------------
int AmfUeContext::StartTimer(size_t msec, timer_repeat_t repeat,
                             zloop_timer_fn handler, timer_arg_t arg) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#if !MME_UNIT_TEST
  int timer_id = -1;
  if ((timer_id = start_timer(&amf_app_task_zmq_ctx, msec, repeat, handler,
                              nullptr)) != -1) {
    amf_app_timers.insert(std::pair<int, timer_arg_t>(timer_id, arg));
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, timer_id);
#else
  OAILOG_FUNC_RETURN(LOG_AMF_APP, 0);
#endif /* !MME_UNIT_TEST */
}
//------------------------------------------------------------------------------
void AmfUeContext::StopTimer(int timer_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
#if !MME_UNIT_TEST
  stop_timer(&amf_app_task_zmq_ctx, timer_id);
  amf_app_timers.erase(timer_id);
#endif /* !MME_UNIT_TEST */
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
//------------------------------------------------------------------------------
bool AmfUeContext::PopTimerArgById(const int timer_id, timer_arg_t* arg) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  try {
    *arg = amf_app_timers.at(timer_id);
    amf_app_timers.erase(timer_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, true);
  } catch (std::out_of_range& e) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, false);
  }
}

//------------------------------------------------------------------------------
int amf_pdu_start_timer(size_t msec, timer_repeat_t repeat,
                        zloop_timer_fn handler, ue_pdu_id_t id) {
  return magma5g::AmfUeContext::Instance().StartPduTimer(msec, repeat, handler,
                                                         id);
}

//------------------------------------------------------------------------------
void amf_pdu_stop_timer(int timer_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  magma5g::AmfUeContext::Instance().StopPduTimer(timer_id);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//------------------------------------------------------------------------------
bool amf_pop_pdu_timer_arg(int timer_id, ue_pdu_id_t* arg) {
  return magma5g::AmfUeContext::Instance().PopPduTimerArgById(timer_id, arg);
}

//------------------------------------------------------------------------------
int AmfUeContext::StartPduTimer(size_t msec, timer_repeat_t repeat,
                                zloop_timer_fn handler, ue_pdu_id_t arg) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  int timer_id = -1;
  if ((timer_id = start_timer(&amf_app_task_zmq_ctx, msec, repeat, handler,
                              nullptr)) != -1) {
    amf_pdu_timers[timer_id] = arg;
  }
  OAILOG_FUNC_RETURN(LOG_AMF_APP, timer_id);
}
//------------------------------------------------------------------------------
void AmfUeContext::StopPduTimer(int timer_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  stop_timer(&amf_app_task_zmq_ctx, timer_id);
  std::map<int, ue_pdu_id_t>::iterator it = amf_pdu_timers.find(timer_id);
  amf_pdu_timers.erase(it);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//------------------------------------------------------------------------------
bool AmfUeContext::PopPduTimerArgById(const int timer_id, ue_pdu_id_t* arg) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  try {
    *arg = amf_pdu_timers.at(timer_id);
    amf_pdu_timers.erase(timer_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, true);
  } catch (std::out_of_range& e) {
    OAILOG_FUNC_RETURN(LOG_AMF_APP, false);
  }
}

}  // namespace magma5g
