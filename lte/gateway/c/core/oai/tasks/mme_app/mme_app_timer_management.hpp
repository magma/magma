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
#ifndef FILE_MME_UE_CONTEXT_H_SEEN
#define FILE_MME_UE_CONTEXT_H_SEEN
// C includes --------------------------------------------------------------
extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.h"
}
// C++ includes ------------------------------------------------------------
#include <czmq.h>
#include <map>
#include <utility>
#include <stddef.h>
#include <stdint.h>
// Other includes ----------------------------------------------------------

namespace magma {
namespace lte {

typedef timer_arg_t TimerArgType;

class MmeUeContext {
 private:
  std::map<int, TimerArgType> mme_app_timers;
  MmeUeContext() : mme_app_timers(){};

 public:
  static MmeUeContext& Instance() {
    static MmeUeContext instance;
    return instance;
  }

  MmeUeContext(MmeUeContext const&) = delete;
  void operator=(MmeUeContext const&) = delete;

  int StartTimer(
      size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
      const TimerArgType& arg);
  void StopTimer(int timer_id);

  bool GetTimerArg(const int timer_id, TimerArgType* arg) const;
};

}  // namespace lte
}  // namespace magma
#endif /* FILE_MME_UE_CONTEXT_H_SEEN */
