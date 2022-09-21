/*
 *
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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <cstdint>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"
#include "lte/gateway/c/core/oai/include/state_converter.hpp"
#include "lte/protos/oai/s1ap_state.pb.h"
#include "lte/gateway/c/core/oai/include/s1ap_state.hpp"

namespace magma {
namespace lte {

class S1apStateConverter : StateConverter {
 public:
  static void state_to_proto(oai::S1apState* state, oai::S1apState* proto);

  static void proto_to_state(const oai::S1apState& proto, oai::S1apState* state);

  /**
   * Serializes s1ap_imsi_map_t to S1apImsiMap proto
   */
  static void s1ap_imsi_map_to_proto(const s1ap_imsi_map_t* s1ap_imsi_map,
                                     oai::S1apImsiMap* s1ap_imsi_proto);

  /**
   * Deserializes s1ap_imsi_map_t from S1apImsiMap proto
   */
  static void proto_to_s1ap_imsi_map(const oai::S1apImsiMap& s1ap_imsi_proto,
                                     s1ap_imsi_map_t* s1ap_imsi_map);

  static void ue_to_proto(const oai::UeDescription* ue,
                          oai::UeDescription* proto);

  static void proto_to_ue(const oai::UeDescription& proto,
                          oai::UeDescription* ue);

 private:
  S1apStateConverter();
  ~S1apStateConverter();
};

}  // namespace lte
}  // namespace magma
