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

using magma::lte::oai::EnbDescription;
using magma::lte::oai::S1apState;
using magma::lte::oai::UeDescription;

namespace magma {
namespace lte {

S1apStateConverter::~S1apStateConverter() = default;
S1apStateConverter::S1apStateConverter() = default;

void S1apStateConverter::state_to_proto(s1ap_state_t* state, S1apState* proto) {
  proto->Clear();

  // copy over enbs
  state_map_to_proto<map_uint32_enb_description_t, EnbDescription,
                     EnbDescription>(state->enbs, proto->mutable_enbs(),
                                     enb_to_proto, LOG_S1AP);

  // copy over mmeid2associd
  mme_ue_s1ap_id_t mmeid;
  sctp_assoc_id_t sctp_assoc_id = 0;
  auto mmeid2associd = proto->mutable_mmeid2associd();

  if (state->mmeid2associd.isEmpty()) {
    OAILOG_DEBUG(LOG_S1AP, "No entries in mmeid2associd map");
  } else {
    *(proto->mutable_mmeid2associd()) = *(state->mmeid2associd.map);
  }

  uint32_t expected_enb_count = state->enbs.size();
  if (expected_enb_count != state->num_enbs) {
    OAILOG_ERROR(LOG_S1AP,
                 "Updating num_eNBs from maintained to actual count %u->%u",
                 state->num_enbs, expected_enb_count);
    state->num_enbs = expected_enb_count;
  }
  proto->set_num_enbs(state->num_enbs);
}

void S1apStateConverter::proto_to_state(const S1apState& proto,
                                        s1ap_state_t* state) {
  proto_to_state_map<map_uint32_enb_description_t, EnbDescription,
                     EnbDescription>(proto.enbs(), state->enbs, proto_to_enb,
                                     LOG_S1AP);

  *(state->mmeid2associd.map) = proto.mmeid2associd();
  state->num_enbs = proto.num_enbs();
  uint32_t expected_enb_count = state->enbs.size();
  OAILOG_WARNING(LOG_S1AP, "expected_enb_count:%d state->num_enbs :%d \n",
                 expected_enb_count, state->num_enbs);
  if (expected_enb_count != state->num_enbs) {
    OAILOG_WARNING(LOG_S1AP,
                   "Updating num_eNBs from maintained to actual count %u->%u",
                   state->num_enbs, expected_enb_count);
    state->num_enbs = expected_enb_count;
  }
}

void S1apStateConverter::enb_to_proto(oai::EnbDescription* enb,
                                      oai::EnbDescription* proto) {
  proto->Clear();
  proto->MergeFrom(*enb);
}

void S1apStateConverter::proto_to_enb(const oai::EnbDescription& proto,
                                      oai::EnbDescription* enb) {
  enb->Clear();
  enb->MergeFrom(proto);
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
