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

namespace magma {
namespace lte {

typedef uint32_t TimerArgType;

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
      TimerArgType id);
  void StopTimer(int timer_id);

  bool GetTimerArg(const int timer_id, TimerArgType* arg) const;
};

}  // namespace lte
}  // namespace magma
#endif /* FILE_MME_UE_CONTEXT_H_SEEN */
