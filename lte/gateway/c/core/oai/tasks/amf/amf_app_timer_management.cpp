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
#if !MME_UNIT_TEST
  return magma5g::AmfUeContext::Instance().StartTimer(msec, repeat, handler,
                                                      id);
#else
  return 0;
#endif /* !MME_UNIT_TEST */
}

//------------------------------------------------------------------------------
void amf_app_stop_timer(int timer_id) {
#if !MME_UNIT_TEST
  magma5g::AmfUeContext::Instance().StopTimer(timer_id);
#endif /* !MME_UNIT_TEST */
}

//------------------------------------------------------------------------------
bool amf_pop_timer_arg(int timer_id, timer_arg_t* arg) {
  return magma5g::AmfUeContext::Instance().PopTimerById(timer_id, arg);
}

//------------------------------------------------------------------------------
int amf_pdu_start_timer(size_t msec, timer_repeat_t repeat,
                        zloop_timer_fn handler, ue_pdu_id_t id) {
  return magma5g::AmfPduUeContext::Instance().StartTimer(msec, repeat, handler,
                                                         id);
}

//------------------------------------------------------------------------------
void amf_pdu_stop_timer(int timer_id) {
  magma5g::AmfPduUeContext::Instance().StopTimer(timer_id);
}

//------------------------------------------------------------------------------
bool amf_pop_pdu_timer_arg(int timer_id, ue_pdu_id_t* arg) {
  return magma5g::AmfPduUeContext::Instance().PopTimerById(timer_id, arg);
}

//------------------------------------------------------------------------------
AmfPduUeContext& AmfPduUeContext::Instance() {
  static AmfPduUeContext instance{&amf_app_task_zmq_ctx};
  return instance;
}

//------------------------------------------------------------------------------
AmfUeContext& AmfUeContext::Instance() {
  static AmfUeContext instance{&amf_app_task_zmq_ctx};
  return instance;
}

}  // namespace magma5g
