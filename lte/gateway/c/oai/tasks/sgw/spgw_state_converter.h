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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "assertions.h"
#include "common_types.h"
#include "hashtable.h"

#ifdef __cplusplus
}
#endif

#include "state_converter.h"
#include "lte/gateway/c/oai/protos/spgw_state.pb.h"
#include "sgw_types.h"
#include "pgw_types.h"
#include "pgw_procedures.h"
#include "spgw_state.h"

namespace magma {
namespace lte {

/**
 * Class for SGW / PGW tasks state conversion helper functions.
 */
class SpgwStateConverter : StateConverter {
 public:
  /**
   * Main function to convert SPGW state to proto definition
   * @param spgw_state const pointer to spgw_state struct
   * @param spgw_proto SpgwState proto object to be written to
   * Memory is owned by the caller
   */
  static void spgw_state_to_proto(const spgw_state_t* spgw_state,
                                  gateway::spgw::SpgwState* spgw_proto);

 private:
  SpgwStateConverter();
  ~SpgwStateConverter();

  /**
   * Converts SGW state to proto, memory is owned by the caller
   * @param sgw_state sgw state struct
   * @param proto object to write to
   */
  static void sgw_state_to_proto(const sgw_state_t* sgw_state,
                                 gateway::spgw::SgwState* proto);

  /**
   * Converts PGW state to proto, memory is owned by the caller
   * @param pgw_state pgw state struct
   * @param proto object to write to
   */
  static void pgw_state_to_proto(const pgw_state_t* pgw_state,
                                 gateway::spgw::PgwState* proto);

  /**
   * Converts s11teid_mme hashtable object to proto, memory owned by the caller
   * @param state_map
   * @param proto_map Map protobuf object to write to
   */
  static void s11teid_mme_ht_to_proto(
      hash_table_ts_t* state_map,
      google::protobuf::Map<unsigned int, gateway::spgw::MmeSgwTunnel>*
          proto_map);

  /**
   * Converts s11bearer_context hashtable object to proto,
   * memory owned by the caller
   * @param state_map
   * @param proto_map Map protobuf object to write to
   */
  static void s11bearer_context_ht_to_proto(
      hash_table_ts_t* state_map,
      google::protobuf::Map<unsigned int, gateway::spgw::S11BearerContext>*
          proto_map);

  /**
   * Converts mme sgw tunnel struct to proto, memory is owned by the caller
   * @param tunnel
   * @param proto
   */
  static void mme_sgw_tunnel_to_proto(const mme_sgw_tunnel_t* tunnel,
                                      gateway::spgw::MmeSgwTunnel* proto);

  /**
   * Converts spgw bearer context struct to proto, memory is owned by the caller
   * @param spgw_bearer_state
   * @param spgw_bearer_proto
   */
  static void spgw_bearer_context_to_proto(
      const s_plus_p_gw_eps_bearer_context_information_t* spgw_bearer_state,
      gateway::spgw::S11BearerContext* spgw_bearer_proto);

  /**
   * Converts sgw eps bearer struct to proto, memory is owned by the caller
   * @param eps_bearer
   * @param eps_bearer_proto
   */
  static void
  sgw_eps_bearer_to_proto(const sgw_eps_bearer_ctxt_t* eps_bearer,
                          gateway::spgw::SgwEpsBearerContext* eps_bearer_proto);

  /**
   * Converts sgw pdn connection struct to proto, memory is owned by the caller
   * @param state_pdn
   * @param proto_pdn
   */
  static void
  sgw_pdn_connection_to_proto(const sgw_pdn_connection_t* state_pdn,
                              gateway::spgw::SgwPdnConnection* proto_pdn);

  /**
   * Converts itti_s11_create_session_request struct to proto, memory is
   * owned by the caller
   * @param session_request
   * @param proto
   */
  static void sgw_create_session_message_to_proto(
      const itti_s11_create_session_request_t* session_request,
      gateway::spgw::CreateSessionMessage* proto);

  /**
   * Converts sgw pending procedures entries list to proto, memory is
   * owned by the caller
   * @param procedures LIST entries of bearer pending procedures
   * @param proto
   */
  static void sgw_pending_procedures_to_proto(
      const sgw_eps_bearer_context_information_t::pending_procedures_s*
          procedures,
      gateway::spgw::SgwEpsBearerContextInfo* proto);

  /**
   * Converts traffic flow template struct to proto, memory is owned by the
   * caller
   * @param tft_state
   * @param tft_proto
   */
  static void
  traffic_flow_template_to_proto(const traffic_flow_template_t* tft_state,
                                 gateway::spgw::TrafficFlowTemplate* tft_proto);

  /**
   * Converts eps bearer QOS struct to proto, memory is owned by the caller
   * @param eps_bearer_qos_state
   * @param eps_bearer_qos_proto
   */
  static void
  eps_bearer_qos_to_proto(const bearer_qos_t* eps_bearer_qos_state,
                          gateway::spgw::BearerQos* eps_bearer_qos_proto);

  /**
   * Converts gtpv1u_data struct to proto, memory is owned by the caller
   * @param gtp_data
   * @param gtp_proto
   */
  static void gtpv1u_data_to_proto(const gtpv1u_data_t* gtp_data,
                                   gateway::spgw::GTPV1uData* gtp_proto);

  /**
   * Converts port range struct to proto, memory is owned by the caller
   * @param port_range
   * @param port_range_proto
   */
  static void port_range_to_proto(const port_range_t* port_range,
                                  gateway::spgw::PortRange* port_range_proto);

  /**
   * Converts packet filter struct to proto, memory is owned by the caller
   * @param packet_filter
   * @param packet_filter_proto
   */
  static void
  packet_filter_to_proto(const packet_filter_t* packet_filter,
                         gateway::spgw::PacketFilter* packet_filter_proto);

  /**
   * Converts pcc_rules hashtable to proto
   * @param state_map
   * @param proto_map
   */
  static void pcc_rule_ht_to_proto(
      hash_table_ts_t* state_map,
      google::protobuf::Map<unsigned int, gateway::spgw::PccRule>* proto_map);

  /**
   * Converts pcc rule object to proto, memory is owned by the caller
   * @param pcc_rule_state
   * @param proto
   */
  static void pcc_rule_to_proto(const pcc_rule_t* pcc_rule_state,
                                gateway::spgw::PccRule* proto);
  // TODO: Implement proto to state struct conversions
};
} // namespace lte
} // namespace magma
