/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    This list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF ;
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF
 * THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */
#ifndef FILE_SPGW_TYPES_SEEN
#define FILE_SPGW_TYPES_SEEN

#include "3gpp_23.401.h"
#include "ip_forward_messages_types.h"
#include "sgw_ie_defs.h"
#include "gtpv1u_types.h"

typedef struct s5_create_session_request_s {
  teid_t context_teid; ///< local SGW S11 Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
} s5_create_session_request_t;

enum s5_failure_cause { S5_OK = 0, PCEF_FAILURE };

typedef struct s5_create_session_response_s {
  teid_t context_teid; ///< local SGW S11 Tunnel Endpoint Identifier
  ebi_t eps_bearer_id;
  itti_sgi_create_end_point_response_t sgi_create_endpoint_resp;
  enum s5_failure_cause failure_cause;
} s5_create_session_response_t;

typedef struct s5_nw_init_actv_bearer_request_s {
  ebi_t lbi;
  teid_t mme_teid_S11;
  teid_t s_gw_teid_S11_S4;
  bearer_qos_t eps_bearer_qos;          ///< Bearer QoS
  traffic_flow_template_t ul_tft;       ///< UL TFT will be sent to UE
  traffic_flow_template_t dl_tft;       ///< DL TFT will be stored at SPGW
  protocol_configuration_options_t pco; ///< PCO protocol_configuration_options
} s5_nw_init_actv_bearer_request_t;

// Data entry for SGW UE context
typedef struct s_plus_p_gw_eps_bearer_context_information_s {
  sgw_eps_bearer_context_information_t sgw_eps_bearer_context_information;
  pgw_eps_bearer_context_information_t pgw_eps_bearer_context_information;
} s_plus_p_gw_eps_bearer_context_information_t;

// Data entry for s11teid2mme
typedef struct mme_sgw_tunnel_s {
  uint32_t local_teid;  ///< Local tunnel endpoint Identifier
  uint32_t remote_teid; ///< Remote tunnel endpoint Identifier
} mme_sgw_tunnel_t;

// AGW-wide state for SPGW task
typedef struct spgw_state_s {
  STAILQ_HEAD(ipv4_list_allocated_s, ipv4_list_elm_s) ipv4_list_allocated;
  hash_table_ts_t* deactivated_predefined_pcc_rules;
  hash_table_ts_t* predefined_pcc_rules;
  gtpv1u_data_t gtpv1u_data;
  teid_t tunnel_id;
  uint32_t gtpv1u_teid;
  struct in_addr sgw_ip_address_S1u_S12_S4_up;
  hash_table_uint64_ts_t* imsi_teid_htbl;
} spgw_state_t;

void handle_s5_create_session_response(
  spgw_state_t* state,
  s_plus_p_gw_eps_bearer_context_information_t *new_bearer_ctxt_info_p,
  s5_create_session_response_t session_resp);
#endif /* FILE_SPGW_TYPES_SEEN */
