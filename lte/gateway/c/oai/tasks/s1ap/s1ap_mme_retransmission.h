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

/*! \file s1ap_mme_retransmission.h
  \brief
  \author Sebastien ROUX
  \company Eurecom
*/

#ifndef FILE_S1AP_MME_RETRANSMISSION_SEEN
#define FILE_S1AP_MME_RETRANSMISSION_SEEN

#include "tree.h"
#include "3gpp_36.401.h"
#include "timer_messages_types.h"

typedef struct s1ap_timer_map_s {
  long timer_id;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;

  RB_ENTRY(s1ap_timer_map_s) entries;
} s1ap_timer_map_t;

int s1ap_mme_timer_map_compare_id(
  const struct s1ap_timer_map_s *const p1,
  const struct s1ap_timer_map_s *const p2);

int s1ap_handle_timer_expiry(timer_has_expired_t *timer_has_expired);

// TODO: (amar) unused functions check with OAI.
int s1ap_timer_insert(
  const mme_ue_s1ap_id_t mme_ue_s1ap_id,
  const long timer_id);

int s1ap_timer_remove_ue(const mme_ue_s1ap_id_t mme_ue_s1ap_id);

#endif /* FILE_S1AP_MME_RETRANSMISSION_SEEN */
