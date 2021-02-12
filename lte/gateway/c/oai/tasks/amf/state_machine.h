/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#ifndef AMF_AS_SEEN
#define AMF_AS_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"

#ifdef __cplusplus
};
#endif
#include "amf_app_ue_context_and_proc.h"  // "amf_data.h" included in it
using namespace std;
namespace magma5g {

typedef struct UE_Handlers_s {
  const char* name;
  void (*func)(void);
} UE_Handlers_t;

typedef struct ue_state_transition_s {
  UE_Handlers_t handler;
  m5gmm_state_t next_state;
  SMSessionFSMState next_sess_state;
} ue_state_transition_t;

void create_state_matrix();
}  // namespace magma5g
#endif
