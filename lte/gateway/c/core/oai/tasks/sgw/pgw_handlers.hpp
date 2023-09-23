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

/*! \file pgw_handlers.hpp
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/gx_messages_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/spgw_state.hpp"

void handle_s5_create_session_request(
    spgw_state_t* spgw_state,
    magma::lte::oai::S11BearerContext* new_bearer_ctxt_info_p,
    teid_t context_teid, ebi_t eps_bearer_id);

void spgw_handle_pcef_create_session_response(
    spgw_state_t* spgw_state,
    const itti_pcef_create_session_response_t* const pcef_csr_resp_p,
    imsi64_t imsi64);

status_code_e spgw_handle_nw_init_deactivate_bearer_rsp(gtpv2c_cause_t cause,
                                                        ebi_t lbi);

status_code_e spgw_handle_nw_initiated_bearer_actv_req(
    spgw_state_t* state,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    imsi64_t imsi64, gtpv2c_cause_value_t* failed_cause);

status_code_e spgw_handle_nw_initiated_bearer_deactv_req(
    const itti_gx_nw_init_deactv_bearer_request_t* const bearer_req_p,
    imsi64_t imsi64);

status_code_e spgw_send_nw_init_activate_bearer_rsp(
    gtpv2c_cause_value_t cause, imsi64_t imsi64,
    bearer_context_within_create_bearer_response_t* bearer_ctx,
    uint8_t default_bearer_id, char* policy_rule_name);

status_code_e spgw_build_and_send_s11_deactivate_bearer_req(
    imsi64_t imsi64, uint8_t no_of_bearers_to_be_deact,
    ebi_t* ebi_to_be_deactivated, bool delete_default_bearer,
    teid_t mme_teid_S11, log_proto_t module);
