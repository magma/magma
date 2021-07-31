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
#include "conversions.h"
#include "intertask_interface.h"
#include "common_types.h"
}
#include "amf_app_timer_management.h"
//--C++ includes ---------------------------------------------------------------
#include <stdexcept>
//--Other includes -------------------------------------------------------------

namespace magma5g {

extern task_zmq_ctx_t amf_app_task_zmq_ctx;
typedef uint32_t timer_arg_t;

//------------------------------------------------------------------------------
int amf_app_start_timer(
    size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
    timer_arg_t id) {
  return magma5g::AmfUeContext::Instance().StartTimer(
      msec, repeat, handler, id);
}

//------------------------------------------------------------------------------
void amf_app_stop_timer(int timer_id) {
  magma5g::AmfUeContext::Instance().StopTimer(timer_id);
}

//------------------------------------------------------------------------------
bool amf_app_get_timer_arg(int timer_id, timer_arg_t* arg) {
  return magma5g::AmfUeContext::Instance().GetTimerArg(timer_id, arg);
}

//------------------------------------------------------------------------------
int AmfUeContext::StartTimer(
    size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
    timer_arg_t arg) {
  int timer_id = -1;
  if ((timer_id = start_timer(
           &amf_app_task_zmq_ctx, msec, repeat, handler, nullptr)) != -1) {
    amf_app_timers.insert(std::pair<int, uint32_t>(timer_id, arg));
  }
  return timer_id;
}
//------------------------------------------------------------------------------
void AmfUeContext::StopTimer(int timer_id) {
  stop_timer(&amf_app_task_zmq_ctx, timer_id);
  amf_app_timers.erase(timer_id);
}
//------------------------------------------------------------------------------
bool AmfUeContext::GetTimerArg(const int timer_id, timer_arg_t* arg) const {
  try {
    *arg = amf_app_timers.at(timer_id);
    return true;
  } catch (std::out_of_range& e) {
    return false;
  }
}

}  // namespace magma5g
