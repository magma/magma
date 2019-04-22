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

/*! \file pgw_handlers.h
* \brief
* \author Lionel Gauthier
* \company Eurecom
* \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_PGW_HANDLERS_SEEN
#define FILE_PGW_HANDLERS_SEEN
#include "s5_messages_types.h"

int pgw_handle_create_bearer_request(
  const itti_s5_create_bearer_request_t *const bearer_req_p);
uint32_t pgw_handle_activate_ded_bearer_rsp(
  const itti_s5_activate_dedicated_bearer_rsp_t *const act_ded_bearer_rsp);
uint32_t pgw_handle_dedicated_bearer_actv_req(
  Imsi_t *imsi,
  ip_address_t *ue_ip,
  traffic_flow_template_t *tft,
  bearer_qos_t *eps_bearer_qos);
uint32_t pgw_handle_deactivate_ded_bearer_rsp(
  const itti_s5_deactivate_dedicated_bearer_rsp_t *const deact_ded_bearer_rsp);
uint32_t pgw_handle_deactivate_ded_bearer_req(
  Imsi_t *imsi,
  uint32_t no_of_bearers,
  ebi_t ebi[]);
#endif /* FILE_PGW_HANDLERS_SEEN */
