/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
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

  std::pair<TimerArgType, bool> GetTimerArg(int timer_id) const;
};

}  // namespace lte
}  // namespace magma
#endif /* FILE_MME_UE_CONTEXT_H_SEEN */
