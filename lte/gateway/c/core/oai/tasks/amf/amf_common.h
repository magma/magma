/**
 * Copyright 2021 The Magma Authors.
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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GNasEnums.h"

#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH) \
  LENGTH += QUADLET - (LENGTH % QUADLET)

namespace magma5g {
std::string uint8_to_hex_string(const uint8_t* v, const size_t s);
status_code_e amf_send_msg_to_task(task_zmq_ctx_t* task_zmq_ctx_p,
                                   task_id_t destination_task_id,
                                   MessageDef* message);
int paa_to_address_info(const paa_t* paa, uint8_t* pdu_address_info,
                        uint8_t* pdu_address_length);
}  // namespace magma5g
