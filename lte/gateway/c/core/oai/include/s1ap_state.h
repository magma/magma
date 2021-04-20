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

#ifdef __cplusplus
extern "C" {
#endif

#include "hashtable.h"
#include "mme_config.h"
#include "s1ap_types.h"

int s1ap_state_init(uint32_t max_ues, uint32_t max_enbs, bool use_stateless);

void s1ap_state_exit(void);

s1ap_state_t* get_s1ap_state(bool read_from_db);

void put_s1ap_state(void);

enb_description_t* s1ap_state_get_enb(
    s1ap_state_t* state, sctp_assoc_id_t assoc_id);

ue_description_t* s1ap_state_get_ue_enbid(
    sctp_assoc_id_t sctp_assoc_id, enb_ue_s1ap_id_t enb_ue_s1ap_id);

ue_description_t* s1ap_state_get_ue_mmeid(mme_ue_s1ap_id_t mme_ue_s1ap_id);

ue_description_t* s1ap_state_get_ue_imsi(imsi64_t imsi64);

/**
 * Return unique composite id for S1AP UE context
 * @param sctp_assoc_id unique SCTP assoc id
 * @param enb_ue_s1ap_id unique UE s1ap ID on eNB
 * @return uint64_t of composite id
 */
#define S1AP_GENERATE_COMP_S1AP_ID(sctp_assoc_id, enb_ue_s1ap_id)              \
  (uint64_t) enb_ue_s1ap_id << 32 | sctp_assoc_id

/**
 * Converts s1ap_imsi_map to protobuf and saves it into data store
 */
void put_s1ap_imsi_map(void);

/**
 * @return s1ap_imsi_map_t pointer
 */
s1ap_imsi_map_t* get_s1ap_imsi_map(void);

hash_table_ts_t* get_s1ap_ue_state(void);

int read_s1ap_ue_state_db(void);

void put_s1ap_ue_state(imsi64_t imsi64);

void delete_s1ap_ue_state(imsi64_t imsi64);

bool s1ap_ue_compare_by_mme_ue_id_cb(
    __attribute__((unused)) hash_key_t keyP, void* elementP, void* parameterP,
    void** resultP);

bool s1ap_ue_compare_by_imsi(
    __attribute__((unused)) hash_key_t keyP, void* elementP, void* parameterP,
    void** resultP);

bool get_mme_ue_ids_no_imsi(
    const hash_key_t keyP, uint64_t const dataP,
    __attribute__((unused)) void* argP, void** resultP);

void remove_ues_without_imsi_from_ue_id_coll(void);

#ifdef __cplusplus
}
#endif
