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
#include "common_types.h"
#include "s11_messages_types.h"
#include "intertask_interface.h"
#include "spgw_types.h"

void sgw_s8_handle_s11_create_session_request(
    sgw_state_t* sgw_state, itti_s11_create_session_request_t* session_req_p,
    imsi64_t imsi64);

void sgw_s8_handle_create_session_response(
    sgw_state_t* sgw_state, s8_create_session_response_t* session_rsp_p,
    imsi64_t imsi64);

void sgw_s8_handle_modify_bearer_request(
    sgw_state_t* state,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64);

void sgw_s8_handle_delete_session_response(
    sgw_state_t* sgw_state, s8_delete_session_response_t* session_rsp_p,
    imsi64_t imsi64);

void sgw_s8_handle_release_access_bearers_request(
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    imsi64_t imsi64);

void sgw_s8_handle_s11_delete_session_request(
    sgw_state_t* sgw_state,
    const itti_s11_delete_session_request_t* const delete_session_req_p,
    imsi64_t imsi64);

teid_t sgw_s8_generate_new_cp_teid(void);

uint32_t sgw_get_new_s5s8u_teid(sgw_state_t* state);

status_code_e sgw_update_teid_in_ue_context(
    sgw_state_t* sgw_state, imsi64_t imsi64, teid_t teid);

sgw_eps_bearer_context_information_t*
sgw_create_bearer_context_information_in_collection(teid_t teid);

sgw_eps_bearer_context_information_t* sgw_get_sgw_eps_bearer_context(
    teid_t teid);

int sgw_update_bearer_context_information_on_csrsp(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    const s8_create_session_response_t* const session_rsp_p);

int sgw_update_bearer_context_information_on_csreq(
    sgw_state_t* sgw_state,
    sgw_eps_bearer_context_information_t* new_sgw_eps_context,
    mme_sgw_tunnel_t sgw_s11_tunnel,
    itti_s11_create_session_request_t* session_req_pP, imsi64_t imsi64);
