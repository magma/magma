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
/*****************************************************************************

  Source      amf_fsm.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef AMF_FSM_SEEN
#define AMF_FSM_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif

#include "bstrlib.h"
#ifdef __cplusplus
};
#endif
using namespace std;

namespace magma5g {
enum amf_fsm_state_t {

  AMF_STATE_MIN = 0,
  // AMF_INVALID = EMM_STATE_MIN,
  AMF_DEREGISTERED,
  AMF_REGISTERED,
  AMF_DEREGISTERED_INITIATED,
  AMF_COMMON_PROCEDURE_INITIATED,
  AMF_STATE_MAX
};

}
#endif
