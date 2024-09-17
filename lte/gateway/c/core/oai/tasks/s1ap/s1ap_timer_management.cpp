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

// --C includes
#include <stdexcept>
extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
}

// --C++ includes
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_timer_management.hpp"
// ---------------------------------------------------------------
// --Other includes
// -------------------------------------------------------------

namespace magma {
namespace lte {

extern task_zmq_ctx_t s1ap_task_zmq_ctx;

//------------------------------------------------------------------------------
int s1ap_start_timer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                     mme_ue_s1ap_id_t ue_id) {
  s1ap_timer_arg_t arg;
  arg.ue_id = ue_id;
  return S1apUeContext::Instance().StartTimer(msec, repeat, handler, arg);
}

//------------------------------------------------------------------------------
void s1ap_stop_timer(int timer_id) {
  S1apUeContext::Instance().StopTimer(timer_id);
}

//------------------------------------------------------------------------------
bool s1ap_pop_timer_arg_ue_id(int timer_id, mme_ue_s1ap_id_t* ue_id) {
  s1ap_timer_arg_t arg;
  bool result = S1apUeContext::Instance().PopTimerById(timer_id, &arg);
  *ue_id = arg.ue_id;
  return result;
}

//------------------------------------------------------------------------------
int S1apUeContext::StartTimer(size_t msec, timer_repeat_t repeat,
                              zloop_timer_fn handler,
                              const s1ap_timer_arg_t arg) {
  int timer_id = -1;
  if ((timer_id = start_timer(&s1ap_task_zmq_ctx, msec, repeat, handler,
                              nullptr)) != -1) {
    s1ap_timers.insert(std::pair<int, s1ap_timer_arg_t>(timer_id, arg));
  }
  return timer_id;
}

//------------------------------------------------------------------------------
void S1apUeContext::StopTimer(int timer_id) {
  stop_timer(&s1ap_task_zmq_ctx, timer_id);
  s1ap_timers.erase(timer_id);
}

//------------------------------------------------------------------------------
bool S1apUeContext::PopTimerById(const int timer_id, s1ap_timer_arg_t* arg) {
  try {
    *arg = s1ap_timers.at(timer_id);
    s1ap_timers.erase(timer_id);
    return true;
  } catch (std::out_of_range& e) {
    return false;
  }
}

}  // namespace lte
}  // namespace magma
