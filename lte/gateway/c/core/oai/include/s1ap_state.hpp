/*
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
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"
#include "lte/protos/oai/s1ap_state.pb.h"

namespace magma {
namespace lte {

int s1ap_state_init(bool use_stateless);

void s1ap_state_exit(void);

oai::S1apState* get_s1ap_state(bool read_from_db);

void put_s1ap_state(void);

proto_map_rc_t s1ap_state_get_enb(oai::S1apState* state,
                                  sctp_assoc_id_t assoc_id,
                                  oai::EnbDescription* enb);

oai::UeDescription* s1ap_state_get_ue_enbid(sctp_assoc_id_t sctp_assoc_id,
                                            enb_ue_s1ap_id_t enb_ue_s1ap_id);

oai::UeDescription* s1ap_state_get_ue_mmeid(mme_ue_s1ap_id_t mme_ue_s1ap_id);

oai::UeDescription* s1ap_state_get_ue_imsi(imsi64_t imsi64);

/**
 * Return unique composite id for S1AP UE context
 * @param sctp_assoc_id unique SCTP assoc id
 * @param enb_ue_s1ap_id unique UE s1ap ID on eNB
 * @return uint64_t of composite id
 */
#define S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id) \
  (uint64_t) enb_ue_s1ap_id << 32 | sctp_assoc_id

/**
 * Converts s1ap_imsi_map to protobuf and saves it into data store
 */
void put_s1ap_imsi_map(void);

/**
 * @return pointer to oai::S1apImsiMap
 */
oai::S1apImsiMap* get_s1ap_imsi_map(void);

map_uint64_ue_description_t* get_s1ap_ue_state(void);

int read_s1ap_ue_state_db(void);

void put_s1ap_ue_state(imsi64_t imsi64);

void delete_s1ap_ue_state(imsi64_t imsi64);

bool s1ap_ue_compare_by_mme_ue_id_cb(__attribute__((unused)) uint64_t keyP,
                                     oai::UeDescription* elementP,
                                     void* parameterP, void** resultP);

bool s1ap_ue_compare_by_imsi(__attribute__((unused)) uint64_t keyP,
                             oai::UeDescription* elementP, void* parameterP,
                             void** resultP);

void remove_ues_without_imsi_from_ue_id_coll(void);

void clean_stale_enb_state(oai::S1apState* state,
                           oai::EnbDescription* new_enb_association);

proto_map_rc_t s1ap_state_update_enb_map(oai::S1apState* state,
                                         sctp_assoc_id_t assoc_id,
                                         oai::EnbDescription* enb);

void get_s1ap_ueid_imsi_map(magma::proto_map_uint32_uint64_t* ueid_imsi_map);

}  // namespace lte
}  // namespace magma
