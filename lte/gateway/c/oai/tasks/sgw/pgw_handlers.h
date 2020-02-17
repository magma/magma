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
#include "sgw_messages_types.h"
#include "spgw_state.h"

void handle_s5_create_session_request(
  spgw_state_t* spgw_state,
  teid_t context_teid,
  ebi_t eps_bearer_id);
uint32_t pgw_handle_nw_init_deactivate_bearer_rsp(
  const itti_s5_nw_init_deactv_bearer_rsp_t *const deact_ded_bearer_rsp);
uint32_t pgw_handle_nw_initiated_bearer_deactv_req(
  spgw_state_t *spgw_state,
  const itti_pgw_nw_init_deactv_bearer_request_t *const bearer_req_p,
  imsi64_t imsi64);
#endif /* FILE_PGW_HANDLERS_SEEN */
