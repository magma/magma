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

#ifdef __cplusplus
extern "C" {
#endif

#include <cstdint>

#include "assertions.h"
#include "dynamic_memory_check.h"
#include "hashtable.h"
#include "s1ap_types.h"

#ifdef __cplusplus
}
#endif

#include "state_converter.h"
#include "lte/protos/oai/s1ap_state.pb.h"
#include "s1ap_state.h"

namespace magma {
namespace lte {

class S1apStateConverter : StateConverter {
 public:
  static void state_to_proto(s1ap_state_t* state, oai::S1apState* proto);

  static void proto_to_state(const oai::S1apState& proto, s1ap_state_t* state);

  /**
   * Serializes s1ap_imsi_map_t to S1apImsiMap proto
   */
  static void s1ap_imsi_map_to_proto(
      const s1ap_imsi_map_t* s1ap_imsi_map, oai::S1apImsiMap* s1ap_imsi_proto);

  /**
   * Deserializes s1ap_imsi_map_t from S1apImsiMap proto
   */
  static void proto_to_s1ap_imsi_map(
      const oai::S1apImsiMap& s1ap_imsi_proto, s1ap_imsi_map_t* s1ap_imsi_map);

  /**
   * Serializes supported_ta_list_t to SupportedTaList proto
   */
  static void supported_ta_list_to_proto(
      const supported_ta_list_t* supported_ta_list,
      oai::SupportedTaList* supported_ta_list_proto);

  static void proto_to_supported_ta_list(
      supported_ta_list_t* supported_ta_list_state,
      const oai::SupportedTaList& supported_ta_list_proto);

  /**
   * Serializes supported_tai_items_t to supported_tai_item proto
   */
  static void supported_tai_item_to_proto(
      const supported_tai_items_t* state_supported_tai_item,
      oai::SupportedTaiItems* supported_tai_item_proto);

  static void proto_to_supported_tai_items(
      supported_tai_items_t* supported_tai_item_state,
      const oai::SupportedTaiItems& supported_tai_item_proto);

  static void enb_to_proto(enb_description_t* enb, oai::EnbDescription* proto);

  static void proto_to_enb(
      const oai::EnbDescription& proto, enb_description_t* enb);

  static void ue_to_proto(
      const ue_description_t* ue, oai::UeDescription* proto);

  static void proto_to_ue(
      const oai::UeDescription& proto, ue_description_t* ue);

 private:
  S1apStateConverter();
  ~S1apStateConverter();
};

}  // namespace lte
}  // namespace magma
