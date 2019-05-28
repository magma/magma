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

#ifndef SGS_MESSAGES_H_
#define SGS_MESSAGES_H_

int sgs_send_eps_detach_indication(
  itti_sgsap_eps_detach_ind_t *sgs_eps_detach_ind_p);

int sgs_send_imsi_detach_indication(
  itti_sgsap_imsi_detach_ind_t *sgs_imsi_detach_ind_p);

int sgs_send_tmsi_reallocation_complete(
  itti_sgsap_tmsi_reallocation_comp_t *sgs_tmsi_realloc_comp_p);

int sgs_send_service_request(
  itti_sgsap_service_request_t *const sgs_service_request_p);

int sgs_send_paging_reject(
  itti_sgsap_paging_reject_t *const sgs_paging_reject_p);

int sgs_send_ue_unreachable(
  itti_sgsap_ue_unreachable_t *const sgs_ue_unreachable_p);

#endif
