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
#include "mme_config.h"
#include "s1ap_types.h"

int s1ap_state_init(uint32_t max_ues, uint32_t max_enbs, bool use_stateless);

void s1ap_state_exit(void);

s1ap_state_t* get_s1ap_state(bool read_from_db);

void put_s1ap_state(void);

enb_description_t* s1ap_state_get_enb(
  s1ap_state_t* state,
  sctp_assoc_id_t assoc_id);

ue_description_t* s1ap_state_get_ue_enbid(
  enb_description_t* enb,
  enb_ue_s1ap_id_t enb_ue_s1ap_id);

ue_description_t* s1ap_state_get_ue_mmeid(
  s1ap_state_t* state,
  mme_ue_s1ap_id_t mme_ue_s1ap_id);

bool s1ap_enb_find_ue_by_mme_ue_id_cb(
  __attribute__((unused)) hash_key_t keyP,
  void* elementP,
  void* parameterP,
  void** resultP);

bool s1ap_ue_compare_by_mme_ue_id_cb(
  __attribute__((unused)) hash_key_t keyP,
  void* elementP,
  void* parameterP,
  void** resultP);

#ifdef __cplusplus
}
#endif
