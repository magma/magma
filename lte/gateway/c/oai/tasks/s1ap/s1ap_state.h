/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#include "s1ap_types.h"

typedef struct s1ap_state_s {
  // contains eNB_description_s, key is eNB_description_s.enb_id (uint32_t)
  hash_table_ts_t enbs;
  // contains sctp association id, key is mme_ue_s1ap_id
  hash_table_ts_t mmeid2associd;
} s1ap_state_t;

int s1ap_state_init(void);
void s1ap_state_exit(void);

s1ap_state_t *s1ap_state_get(void);
void s1ap_state_put(s1ap_state_t *state);

enb_description_t *s1ap_state_get_enb(
  s1ap_state_t *state,
  sctp_assoc_id_t assoc_id);
ue_description_t *s1ap_state_get_ue_enbid(
  s1ap_state_t *state,
  enb_description_t *enb,
  enb_ue_s1ap_id_t enb_ue_s1ap_id);
ue_description_t *s1ap_state_get_ue_mmeid(
  s1ap_state_t *state,
  mme_ue_s1ap_id_t mme_ue_s1ap_id);

#ifdef __cplusplus
}
#endif
