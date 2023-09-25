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

#include <czmq.h>

#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

namespace magma {
namespace lte {

typedef struct s1ap_timer_arg_s {
  mme_ue_s1ap_id_t ue_id;
} s1ap_timer_arg_t;

int s1ap_start_timer(size_t msec, timer_repeat_t repeat, zloop_timer_fn handler,
                     mme_ue_s1ap_id_t ue_id);

void s1ap_stop_timer(int timer_id);

// The *_pop_timer_* functions also removes the timer_id from the map.
// These functions are supposed to be used only by expired timers.
bool s1ap_pop_timer_arg_ue_id(int timer_id, mme_ue_s1ap_id_t* ue_id);

}  // namespace lte
}  // namespace magma
