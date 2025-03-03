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

class AmfUeContext {
 private:
  std::map<int, timer_arg_t> amf_app_timers;
  std::map<int, ue_pdu_id_t> amf_pdu_timers;
  AmfUeContext() : amf_app_timers(), amf_pdu_timers() {};

 public:
  static AmfUeContext& Instance() {
    static AmfUeContext instance;
    return instance;
  }

  AmfUeContext(AmfUeContext const&) = delete;
  void operator=(AmfUeContext const&) = delete;

  int StartTimer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                 timer_arg_t id);
  void StopTimer(int timer_id);

  int StartPduTimer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                    ue_pdu_id_t id);
  void StopPduTimer(int timer_id);

  /**
   * Pop timer, save arguments and return existence.
   *
   * @param timer_id Unique timer id for active timers
   * @param arg Timer arguments to be stored in this parameter
   * @return True if timer_id exists, False otherwise.
   */
  bool PopTimerArgById(const int timer_id, timer_arg_t* arg);
  bool PopPduTimerArgById(const int timer_id, ue_pdu_id_t* arg);
};

}  // namespace magma5g
