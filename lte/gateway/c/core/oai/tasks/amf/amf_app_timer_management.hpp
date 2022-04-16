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

#pragma once
//--C includes -----------------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
}
////--C++ includes
///---------------------------------------------------------------
#include "lte/gateway/c/core/oai/common/common_utility_funs.hpp"
////--Other includes
///-------------------------------------------------------------
#include <czmq.h>
#include <map>
#include <stddef.h>
#include <stdint.h>

namespace magma5g {

typedef uint64_t timer_arg_t;

typedef struct ue_pdu_id {
  amf_ue_ngap_id_t ue_id;
  uint8_t pdu_id;
} ue_pdu_id_t;

#define AMF_APP_TIMER_INACTIVE_ID (-1)

int amf_app_start_timer(size_t msec, timer_repeat_t repeat,
                        zloop_timer_fn handler, timer_arg_t id);

void amf_app_stop_timer(int timer_id);

int amf_pdu_start_timer(size_t msec, timer_repeat_t repeat,
                        zloop_timer_fn handler, ue_pdu_id_t id);

void amf_pdu_stop_timer(int timer_id);

/*void amf_app_resume_timer(
    struct ue_mm_context_s* const ue_mm_context_pP, time_t start_time,
    struct amf_app_timer_t* timer, zloop_timer_fn timer_expiry_handler,
    char* timer_name);
*/

// The *_pop*timer_* functions also removes the timer_id from the map.
// These functions are supposed to be used only by expired timers.
bool amf_pop_timer_arg(int timer_id, timer_arg_t* arg);
bool amf_pop_pdu_timer_arg(int timer_id, ue_pdu_id_t* arg);

struct AmfUeContext : public magma::utils::AppTimerUeContext<timer_arg_t> {
  static AmfUeContext& Instance();
  explicit AmfUeContext(task_zmq_ctx_s* zctx)
      : magma::utils::AppTimerUeContext<timer_arg_t>{zctx} {}
};

struct AmfPduUeContext : public magma::utils::AppTimerUeContext<ue_pdu_id_t> {
  static AmfPduUeContext& Instance();
  explicit AmfPduUeContext(task_zmq_ctx_s* zctx)
      : magma::utils::AppTimerUeContext<ue_pdu_id_t>{zctx} {}
};

}  // namespace magma5g
