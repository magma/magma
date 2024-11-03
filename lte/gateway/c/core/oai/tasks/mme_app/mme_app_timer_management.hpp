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
// C includes --------------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
}
// C++ includes ------------------------------------------------------------
#include <czmq.h>
#include <map>
#include <utility>
#include <stddef.h>
#include <stdint.h>
// Other includes ----------------------------------------------------------
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.hpp"

namespace magma {
namespace lte {

typedef timer_arg_t TimerArgType;

class MmeUeContext {
 private:
  std::map<int, TimerArgType> mme_app_timers;
  MmeUeContext() : mme_app_timers() {};

 public:
  static MmeUeContext& Instance() {
    static MmeUeContext instance;
    return instance;
  }

  MmeUeContext(MmeUeContext const&) = delete;
  void operator=(MmeUeContext const&) = delete;

  int StartTimer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                 const TimerArgType& arg);
  void StopTimer(int timer_id);

  /**
   * Pop timer, save arguments and return existence.
   *
   * @param timer_id Unique timer id for active timers
   * @param arg Timer arguments to be stored in this parameter
   * @return True if timer_id exists, False otherwise.
   */
  bool PopTimerById(const int timer_id, TimerArgType* arg);
};

}  // namespace lte
}  // namespace magma
