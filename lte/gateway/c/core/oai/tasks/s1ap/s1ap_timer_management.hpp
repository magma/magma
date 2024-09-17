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

// C++ includes ------------------------------------------------------------
#include <stddef.h>
#include <stdint.h>
#include <czmq.h>
#include <map>
#include <utility>

// C includes --------------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
}
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_timer.hpp"

// Other includes ----------------------------------------------------------

namespace magma {
namespace lte {

class S1apUeContext {
 private:
  std::map<int, s1ap_timer_arg_t> s1ap_timers;
  S1apUeContext() : s1ap_timers() {}

 public:
  static S1apUeContext& Instance() {
    static S1apUeContext instance;
    return instance;
  }

  S1apUeContext(S1apUeContext const&) = delete;
  void operator=(S1apUeContext const&) = delete;

  int StartTimer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                 const s1ap_timer_arg_t arg);
  void StopTimer(int timer_id);

  /**
   * Pop timer, save arguments and return existence.
   *
   * @param timer_id Unique timer id for active timers
   * @param arg Timer arguments to be stored in this parameter
   * @return True if timer_id exists, False otherwise.
   */
  bool PopTimerById(const int timer_id, s1ap_timer_arg_t* arg);
};

}  // namespace lte
}  // namespace magma
