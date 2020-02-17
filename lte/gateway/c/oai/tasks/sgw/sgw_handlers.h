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

/*! \file sgw_handlers.h
* \brief
* \author Lionel Gauthier
* \company Eurecom
* \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_SGW_HANDLERS_SEEN
#define FILE_SGW_HANDLERS_SEEN

#include "common_types.h"
#include "gtpv1_u_messages_types.h"
#include "ip_forward_messages_types.h"
#include "s11_messages_types.h"
#include "s5_messages_types.h"
#include "spgw_state.h"

int sgw_handle_s11_create_session_request(
  spgw_state_t* state,
  const itti_s11_create_session_request_t* const session_req_p,
  imsi64_t imsi64);
int sgw_handle_sgi_endpoint_created(
  spgw_state_t *state,
  itti_sgi_create_end_point_response_t *const resp_p,
  imsi64_t imsi64);
int sgw_handle_sgi_endpoint_updated(
  spgw_state_t *state,
  const itti_sgi_update_end_point_response_t *const resp_p,
  imsi64_t imsi64);
int sgw_handle_gtpv1uCreateTunnelResp(
  spgw_state_t *state,
  const Gtpv1uCreateTunnelResp *const endpoint_created_p, imsi64_t imsi64);
int sgw_handle_gtpv1uUpdateTunnelResp(
  spgw_state_t *state,
  const Gtpv1uUpdateTunnelResp *const endpoint_updated_p, imsi64_t imsi64);
int sgw_handle_gtpv1uDeleteTunnelResp(
  const Gtpv1uDeleteTunnelResp *const endpoint_deleted_p);
int sgw_handle_modify_bearer_request(
  spgw_state_t *state,
  const itti_s11_modify_bearer_request_t *const modify_bearer_p,
  imsi64_t imsi64);
int sgw_handle_delete_session_request(
  spgw_state_t *state,
  const itti_s11_delete_session_request_t *const delete_session_p,
  imsi64_t imsi64);
int sgw_handle_release_access_bearers_request(
  spgw_state_t *state,
  const itti_s11_release_access_bearers_request_t
    *const release_access_bearers_req_pP,
    imsi64_t imsi64);
int sgw_handle_suspend_notification(
  spgw_state_t *state,
  const itti_s11_suspend_notification_t *const suspend_notification_pP,
  imsi64_t imsi64);
int sgw_no_pcef_create_dedicated_bearer(spgw_state_t *state, s11_teid_t teid,
  imsi64_t imsi64);
int sgw_handle_create_bearer_response(
  spgw_state_t *state,
  const itti_s11_create_bearer_response_t *const create_bearer_response_pP);
int spgw_handle_nw_initiated_bearer_actv_req(
  spgw_state_t* state,
  const itti_spgw_nw_init_actv_bearer_request_t* const bearer_req_p,
  imsi64_t imsi64);
int sgw_handle_nw_initiated_actv_bearer_rsp(
  spgw_state_t *state,
  const itti_s11_nw_init_actv_bearer_rsp_t *const s11_actv_bearer_rsp,
  imsi64_t imsi64);
uint32_t sgw_handle_nw_initiated_deactv_bearer_req(
  const itti_s5_nw_init_deactv_bearer_request_t
    *const itti_s5_deactiv_ded_bearer_req, imsi64_t imsi64);
int sgw_handle_nw_initiated_deactv_bearer_rsp(
  spgw_state_t *state,
  const itti_s11_nw_init_deactv_bearer_rsp_t
    *const s11_pcrf_ded_bearer_deactv_rsp,
    imsi64_t imsi64);
bool is_enb_ip_address_same(const fteid_t *fte_p, ip_address_t *ip_p);
#endif /* FILE_SGW_HANDLERS_SEEN */
