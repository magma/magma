
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

  Source      amf_message.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_fsm.h"
#include "amf_app_ue_context_and_proc.h"
#include "M5gNasMessage.h"
namespace magma5g {
int AMFMsg::amf_msg_decode(AMFMsg* msg, uint8_t* buffer, uint32_t len) {
  int header_result = 0;
  int decode_result = 0;
  header_result     = amf_msg_decode_header(&msg->header, buffer, len);
  if (header_result < 0) {
    // some error msg put in log file.
  }
  buffer += header_result;
  len -= header_result;
  if (decode_result < 0) {
    // TODO some error logs
  }
}
}  // namespace magma5g
