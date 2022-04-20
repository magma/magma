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
#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#ifdef __cplusplus
}
#endif
#include <memory.h>
#include <string.h>
#include <vector>
#include "lte/gateway/c/core/oai/include/map.h"

namespace magma5g {

class NGAPClientServicerBase {
 public:
  virtual ~NGAPClientServicerBase() = default;
  virtual status_code_e send_message_to_amf(task_zmq_ctx_t* task_zmq_ctx_p,
                                            task_id_t destination_task_id,
                                            MessageDef* message) = 0;
};

class NGAPClientServicer : public NGAPClientServicerBase {
 public:
  NGAPClientServicer();
  std::vector<MessagesIds>
      msgtype_stack;  // stack maintains type of msgs sent to amf
  static NGAPClientServicer& getInstance();

  NGAPClientServicer(NGAPClientServicer const&) = delete;
  void operator=(NGAPClientServicer const&) = delete;

  magma::map_string_string_t map_ngap_state_proto_str;
  magma::map_string_string_t map_ngap_uestate_proto_str;
  magma::map_string_string_t map_imsi_table_proto_str;

  status_code_e send_message_to_amf(task_zmq_ctx_t* task_zmq_ctx_p,
                                    task_id_t destination_task_id,
                                    MessageDef* message);
};

}  // namespace magma5g
