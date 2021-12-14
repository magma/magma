/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/s11_messages_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/include/spgw_types.h"

status_code_e sgw_s8_handle_s11_create_session_request(
    sgw_state_t* sgw_state, itti_s11_create_session_request_t* session_req_p,
    imsi64_t imsi64);

status_code_e sgw_s8_handle_create_session_response(
    sgw_state_t* sgw_state, s8_create_session_response_t* session_rsp_p,
    imsi64_t imsi64);

void sgw_s8_handle_modify_bearer_request(
    sgw_state_t* state,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64);

status_code_e sgw_s8_handle_delete_session_response(
    sgw_state_t* sgw_state, s8_delete_session_response_t* session_rsp_p,
    imsi64_t imsi64);

void sgw_s8_handle_release_access_bearers_request(
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    imsi64_t imsi64);

status_code_e sgw_s8_handle_s11_delete_session_request(
    sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64);

imsi64_t sgw_s8_handle_create_bearer_request(
    sgw_state_t* sgw_state, const s8_create_bearer_request_t* const cb_req,
    gtpv2c_cause_value_t* cause_value);

void sgw_s8_handle_s11_create_bearer_response(
    sgw_state_t* sgw_state,
    itti_s11_nw_init_actv_bearer_rsp_t* s11_actv_bearer_rsp, imsi64_t imsi64);

int sgw_s8_handle_delete_bearer_request(
    sgw_state_t* sgw_state, const s8_delete_bearer_request_t* const db_req);

status_code_e sgw_s8_handle_s11_delete_bearer_response(
    sgw_state_t* sgw_state,
    const itti_s11_nw_init_deactv_bearer_rsp_t* const
        s11_delete_bearer_response_p,
    imsi64_t imsi64);

void sgw_s8_send_failed_create_bearer_response(
    sgw_state_t* sgw_state, uint32_t sequence_number, char* pgw_cp_address,
    gtpv2c_cause_value_t cause_value, Imsi_t imsi, teid_t pgw_s8_teid);
teid_t sgw_s8_generate_new_cp_teid(void);

uint32_t sgw_get_new_s5s8u_teid(sgw_state_t* state);

status_code_e sgw_update_teid_in_ue_context(
    sgw_state_t* sgw_state, imsi64_t imsi64, teid_t teid);

sgw_eps_bearer_context_information_t*
sgw_create_bearer_context_information_in_collection(
    sgw_state_t* sgw_state, uint32_t* temporary_create_session_procedure_id);

sgw_eps_bearer_context_information_t* sgw_get_sgw_eps_bearer_context(
    teid_t teid);

int sgw_update_bearer_context_information_on_csrsp(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const s8_create_session_response_t* const session_rsp_p);

int sgw_update_bearer_context_information_on_csreq(
    sgw_state_t* sgw_state,
    sgw_eps_bearer_context_information_t* new_sgw_eps_context,
    itti_s11_create_session_request_t* session_req_pP, imsi64_t imsi64);

uint32_t sgw_get_new_s1u_teid(sgw_state_t* state);

int update_pgw_info_to_temp_dedicated_bearer_context(
    sgw_eps_bearer_context_information_t* sgw_context_p, teid_t s1_u_sgw_fteid,
    s8_bearer_context_t* bc_cbreq, sgw_state_t* sgw_state,
    char* pgw_cp_ip_port);

void sgw_s8_proc_s11_create_bearer_rsp(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    bearer_context_within_create_bearer_response_t* bc_cbrsp,
    itti_s11_nw_init_actv_bearer_rsp_t* s11_actv_bearer_rsp, imsi64_t imsi64,
    sgw_state_t* sgw_state);

void print_bearer_ids_helper(const ebi_t* ebi, uint32_t no_of_bearers);

void sgw_s8_send_failed_delete_bearer_response(
    const s8_delete_bearer_request_t* const db_req,
    gtpv2c_cause_value_t cause_value, Imsi_t imsi, teid_t pgw_s8_teid);
