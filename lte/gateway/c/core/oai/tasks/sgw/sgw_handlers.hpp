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

/*! \file sgw_handlers.hpp
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
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/ip_forward_messages_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/gtpv1_u_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/s11_messages_types.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"
#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u.hpp"

extern task_zmq_ctx_t spgw_app_task_zmq_ctx;

status_code_e sgw_handle_s11_create_session_request(
    spgw_state_t* state,
    const itti_s11_create_session_request_t* const session_req_p,
    imsi64_t imsi64);
void sgw_handle_sgi_endpoint_updated(
    const itti_sgi_update_end_point_response_t* const resp_p, imsi64_t imsi64);
int sgw_handle_sgi_endpoint_deleted(
    const itti_sgi_delete_end_point_request_t* const resp_pP, imsi64_t imsi64);
status_code_e sgw_handle_modify_bearer_request(
    const itti_s11_modify_bearer_request_t* const modify_bearer_p,
    imsi64_t imsi64);
status_code_e sgw_handle_delete_session_request(
    const itti_s11_delete_session_request_t* const delete_session_p,
    imsi64_t imsi64);
void sgw_handle_release_access_bearers_request(
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    imsi64_t imsi64);
status_code_e sgw_handle_suspend_notification(
    const itti_s11_suspend_notification_t* const suspend_notification_pP,
    imsi64_t imsi64);
status_code_e sgw_handle_nw_initiated_actv_bearer_rsp(
    const itti_s11_nw_init_actv_bearer_rsp_t* const s11_actv_bearer_rsp,
    imsi64_t imsi64);
status_code_e sgw_handle_nw_initiated_deactv_bearer_rsp(
    spgw_state_t* spgw_state,
    const itti_s11_nw_init_deactv_bearer_rsp_t* const
        s11_pcrf_ded_bearer_deactv_rsp,
    imsi64_t imsi64);
status_code_e sgw_handle_ip_allocation_rsp(
    spgw_state_t* spgw_state,
    const itti_ip_allocation_response_t* ip_allocation_rsp, imsi64_t imsi64);
uint32_t spgw_get_new_s1u_teid(spgw_state_t* state);

bool is_enb_ip_address_same(const fteid_t* fte_p, struct in_addr ipv4,
                            struct in6_addr ipv6);
status_code_e sgw_handle_sgi_endpoint_created(
    spgw_state_t* state, itti_sgi_create_end_point_response_t* const resp_p,
    imsi64_t imsi64);
status_code_e send_mbr_failure(
    log_proto_t module,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    imsi64_t imsi64);
void sgw_populate_mbr_bearer_contexts_removed(
    const itti_sgi_update_end_point_response_t* const resp_pP, imsi64_t imsi64,
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p,
    itti_s11_modify_bearer_response_t* modify_response_p);
void sgw_populate_mbr_bearer_contexts_not_found(
    log_proto_t module,
    const itti_sgi_update_end_point_response_t* const resp_pP,
    itti_s11_modify_bearer_response_t* modify_response_p);

void populate_sgi_end_point_update(
    uint8_t sgi_rsp_idx, uint8_t idx,
    const itti_s11_modify_bearer_request_t* const modify_bearer_pP,
    magma::lte::oai::SgwEpsBearerContext* eps_bearer_ctxt_p,
    itti_sgi_update_end_point_response_t* sgi_update_end_point_resp);

bool does_sgw_bearer_context_hold_valid_enb_ip(ip_address_t enb_ip_address_S1u);
bool does_bearer_context_hold_valid_enb_ip(
    magma::lte::oai::IpTupple enb_ip_address_S1u);
void sgw_process_release_access_bearer_request(
    log_proto_t module, imsi64_t imsi64,
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context);
void sgw_send_release_access_bearer_response(
    log_proto_t module, task_id_t origin_task_id, imsi64_t imsi64,
    gtpv2c_cause_value_t cause,
    const itti_s11_release_access_bearers_request_t* const
        release_access_bearers_req_pP,
    teid_t mme_teid_s11);

status_code_e sgw_build_and_send_s11_create_bearer_request(
    magma::lte::oai::SgwEpsBearerContextInfo*
        sgw_eps_bearer_context_information,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    pdn_type_t pdn_type, uint32_t sgw_ip_address_S1u_S12_S4_up,
    struct in6_addr* sgw_ipv6_address_S1u_S12_S4_up, teid_t s1_u_sgw_fteid,
    log_proto_t module);

status_code_e create_temporary_dedicated_bearer_context(
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_ctxt_p,
    const itti_gx_nw_init_actv_bearer_request_t* const bearer_req_p,
    pdn_type_t pdn_type, uint32_t sgw_ip_address_S1u_S12_S4_up,
    struct in6_addr* sgw_ipv6_address_S1u_S12_S4_up, teid_t s1_u_sgw_fteid,
    uint32_t sequence_number, log_proto_t module);
// TODO(rsarwad): shall be removed while porting sgw_s8 context to protobuf
// github issue: 11191
void handle_failed_s8_create_bearer_response(
    sgw_eps_bearer_context_information_t* sgw_context_p,
    gtpv2c_cause_value_t cause, imsi64_t imsi64,
    bearer_context_within_create_bearer_response_t* bearer_context,
    sgw_eps_bearer_ctxt_t* dedicated_bearer_ctxt_p, log_proto_t module);

void handle_failed_create_bearer_response(
    magma::lte::oai::SgwEpsBearerContextInfo* sgw_context_p,
    gtpv2c_cause_value_t cause, imsi64_t imsi64,
    bearer_context_within_create_bearer_response_t* bearer_context,
    magma::lte::oai::SgwEpsBearerContext* dedicated_bearer_ctxt_p,
    log_proto_t module);

status_code_e spgw_build_and_send_s11_deactivate_bearer_req(
    imsi64_t imsi64, uint8_t no_of_bearers_to_be_deact,
    ebi_t* ebi_to_be_deactivated, bool delete_default_bearer,
    teid_t mme_teid_S11, log_proto_t module);

void generate_dl_flow(packet_filter_contents_t* packet_filter,
                      in_addr_t ipv4_s_addr, struct in6_addr* ue_ipv6,
                      struct ip_flow_dl* dlflow);
void sgw_handle_delete_bearer_cmd(
    itti_s11_delete_bearer_command_t* s11_delete_bearer_command,
    imsi64_t imsi64);

void convert_proto_ip_to_standard_ip_fmt(magma::lte::oai::IpTupple* proto_ip,
                                         struct in_addr* ipv4,
                                         struct in6_addr* ipv6,
                                         bool ipv6_enabled);
void traffic_flow_template_to_proto(
    const traffic_flow_template_t* tft_state,
    magma::lte::oai::TrafficFlowTemplate* tft_proto);

void eps_bearer_qos_to_proto(
    const bearer_qos_t* eps_bearer_qos_state,
    magma::lte::oai::SgwBearerQos* eps_bearer_qos_proto);

void port_range_to_proto(const port_range_t* port_range,
                         magma::lte::oai::PortRange* port_range_proto);

void sgw_create_session_message_to_proto(
    const itti_s11_create_session_request_t* session_request,
    magma::lte::oai::CreateSessionMessage* proto);

void proto_to_packet_filter(
    const magma::lte::oai::PacketFilter& packet_filter_proto,
    packet_filter_t* packet_filter);
