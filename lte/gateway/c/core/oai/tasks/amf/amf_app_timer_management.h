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
#include "intertask_interface.h"
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

typedef uint32_t timer_arg_t;
// typedef std::pair<uint8_t, uint8_t> ue_pdu_id;
typedef struct ue_pdu_id {
  uint8_t ue_id;
  uint8_t pdu_id;
} ue_pdu_id_t;

#define AMF_APP_TIMER_INACTIVE_ID (-1)

int amf_app_start_timer(
    size_t msec, timer_repeat_t repeat, zloop_timer_fn handler, timer_arg_t id);

void amf_app_stop_timer(int timer_id);

int amf_pdu_start_timer(
    size_t msec, timer_repeat_t repeat, zloop_timer_fn handler, ue_pdu_id_t id);

void amf_pdu_stop_timer(int timer_id);

/*void amf_app_resume_timer(
    struct ue_mm_context_s* const ue_mm_context_pP, time_t start_time,
    struct amf_app_timer_t* timer, zloop_timer_fn timer_expiry_handler,
    char* timer_name);
*/

bool amf_app_get_timer_arg(int timer_id, timer_arg_t* arg);
bool amf_pdu_get_timer_arg(int timer_id, ue_pdu_id_t* arg);

class AmfUeContext {
 private:
  std::map<int, timer_arg_t> amf_app_timers;
  std::map<int, ue_pdu_id_t> amf_pdu_timers;
  AmfUeContext() : amf_app_timers(), amf_pdu_timers(){};

 public:
  static AmfUeContext& Instance() {
    static AmfUeContext instance;
    return instance;
  }

  AmfUeContext(AmfUeContext const&) = delete;
  void operator=(AmfUeContext const&) = delete;

  int StartTimer(
      size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
      timer_arg_t id);
  void StopTimer(int timer_id);

  bool GetTimerArg(const int timer_id, timer_arg_t* arg) const;

  int StartPduTimer(
      size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
      ue_pdu_id_t id);
  void StopPduTimer(int timer_id);

  bool GetPduTimerArg(const int timer_id, ue_pdu_id_t* arg) const;
};

}  // namespace magma5g
