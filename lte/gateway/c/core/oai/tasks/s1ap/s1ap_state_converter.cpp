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
#include "lte/gateway/c/core/oai/tasks/s1ap/s1ap_state_converter.hpp"

namespace magma {
namespace lte {

S1apStateConverter::~S1apStateConverter() = default;
S1apStateConverter::S1apStateConverter() = default;

void S1apStateConverter::state_to_proto(oai::S1apState* state,
                                        oai::S1apState* proto) {
  proto->Clear();
  proto->MergeFrom(*state);
}

void S1apStateConverter::proto_to_state(const oai::S1apState& proto,
                                        oai::S1apState* state) {
  state->Clear();
  state->MergeFrom(proto);
}

void S1apStateConverter::ue_to_proto(const oai::UeDescription* ue,
                                     oai::UeDescription* proto) {
  proto->Clear();
  proto->MergeFrom(*ue);
}

void S1apStateConverter::proto_to_ue(const oai::UeDescription& proto,
                                     oai::UeDescription* ue) {
  ue->Clear();
  ue->MergeFrom(proto);
}

void S1apStateConverter::s1ap_imsi_map_to_proto(
    const s1ap_imsi_map_t* s1ap_imsi_map, oai::S1apImsiMap* s1ap_imsi_proto) {
  *s1ap_imsi_proto->mutable_mme_ue_s1ap_id_imsi_map() =
      *(s1ap_imsi_map->mme_ueid2imsi_map.map);
}

void S1apStateConverter::proto_to_s1ap_imsi_map(
    const oai::S1apImsiMap& s1ap_imsi_proto, s1ap_imsi_map_t* s1ap_imsi_map) {
  *(s1ap_imsi_map->mme_ueid2imsi_map.map) =
      s1ap_imsi_proto.mme_ue_s1ap_id_imsi_map();
}

}  // namespace lte
}  // namespace magma
