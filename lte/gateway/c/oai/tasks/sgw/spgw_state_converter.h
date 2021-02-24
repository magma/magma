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
#include "lte/protos/oai/std_3gpp_types.pb.h"
#include "lte/protos/oai/spgw_state.pb.h"
#include "spgw_types.h"
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
  static void state_to_proto(
      const spgw_state_t* spgw_state, oai::SpgwState* spgw_proto);

  /**
   * Main function to convert SPGW proto to state definition
   * @param spgw_proto SpgwState proto object to be written to
   * @param spgw_state const pointer to spgw_state struct
   * Memory is owned by the caller
   */
  static void proto_to_state(
      const oai::SpgwState& proto, spgw_state_t* spgw_state);

  static void ue_to_proto(
      const spgw_ue_context_t* ue_state, oai::SpgwUeContext* ue_proto);

  static void proto_to_ue(
      const oai::SpgwUeContext& ue_proto, spgw_ue_context_t* ue_context_p);

 private:
  SpgwStateConverter();
  ~SpgwStateConverter();

  /**
   * Converts spgw bearer context struct to proto, memory is owned by the caller
   * @param spgw_bearer_state
   * @param spgw_bearer_proto
   */
  static void spgw_bearer_context_to_proto(
      const s_plus_p_gw_eps_bearer_context_information_t* spgw_bearer_state,
      oai::S11BearerContext* spgw_bearer_proto);

  /**
   * Converts proto to spgw bearer context struct
   * @param spgw_bearer_proto
   * @param spgw_bearer_state
   */
  static void proto_to_spgw_bearer_context(
      const oai::S11BearerContext& spgw_bearer_proto,
      s_plus_p_gw_eps_bearer_context_information_t* spgw_bearer_state);

  /**
   * Converts sgw eps bearer struct to proto, memory is owned by the caller
   * @param eps_bearer
   * @param eps_bearer_proto
   */
  static void sgw_eps_bearer_to_proto(
      const sgw_eps_bearer_ctxt_t* eps_bearer,
      oai::SgwEpsBearerContext* eps_bearer_proto);

  /**
   * Converts proto to sgw eps bearer struct to proto
   * @param eps_bearer_proto
   * @param eps_bearer
   */
  static void proto_to_sgw_eps_bearer(
      const oai::SgwEpsBearerContext& eps_bearer_proto,
      sgw_eps_bearer_ctxt_t* eps_bearer);

  /**
   * Converts sgw pdn connection struct to proto, memory is owned by the caller
   * @param state_pdn
   * @param proto_pdn
   */
  static void sgw_pdn_connection_to_proto(
      const sgw_pdn_connection_t* state_pdn, oai::SgwPdnConnection* proto_pdn);

  /**
   * Converts proto to sgw pdn connection struct
   * @param proto
   * @param state_pdn
   */
  static void proto_to_sgw_pdn_connection(
      const oai::SgwPdnConnection& proto, sgw_pdn_connection_t* state_pdn);

  /**
   * Converts itti_s11_create_session_request struct to proto, memory is
   * owned by the caller
   * @param session_request
   * @param proto
   */
  static void sgw_create_session_message_to_proto(
      const itti_s11_create_session_request_t* session_request,
      oai::CreateSessionMessage* proto);

  /**
   * Converts proto to itti_s11_create_session_request struct
   * @param proto
   * @param session_request
   */
  static void proto_to_sgw_create_session_message(
      const oai::CreateSessionMessage& proto,
      itti_s11_create_session_request_t* session_request);

  /**
   * Converts sgw pending procedures entries list to proto, memory is
   * owned by the caller
   * @param procedures LIST entries of bearer pending procedures
   * @param proto
   */
  static void sgw_pending_procedures_to_proto(
      const sgw_eps_bearer_context_information_t::pending_procedures_s*
          procedures,
      oai::SgwEpsBearerContextInfo* proto);

  /**
   * Converts proto to sgw pending procedures entries list
   * @param proto sgw eps bearer context info proto to read from
   * @param procedures LIST entries of bearer pending procedures
   */
  static void proto_to_sgw_pending_procedures(
      const oai::SgwEpsBearerContextInfo& proto,
      sgw_eps_bearer_context_information_t::pending_procedures_s** procedures);

  /**
   * Inserts new procedure struct to eps bearer pending procedures list
   * @param proto
   * @param pending_procedures entries list to insert new proc
   */
  static void insert_proc_into_sgw_pending_procedures(
      const oai::PgwCbrProcedure& proto,
      sgw_eps_bearer_context_information_t::pending_procedures_s*
          pending_procedures);

  /**
   * Converts traffic flow template struct to proto, memory is owned by the
   * caller
   * @param tft_state
   * @param tft_proto
   */
  static void traffic_flow_template_to_proto(
      const traffic_flow_template_t* tft_state,
      oai::TrafficFlowTemplate* tft_proto);

  /**
   * Converts proto to traffic flow template struct
   * @param tft_proto
   * @param tft_state
   */
  static void proto_to_traffic_flow_template(
      const oai::TrafficFlowTemplate& tft_proto,
      traffic_flow_template_t* tft_state);

  /**
   * Converts eps bearer QOS struct to proto, memory is owned by the caller
   * @param eps_bearer_qos_state
   * @param eps_bearer_qos_proto
   */
  static void eps_bearer_qos_to_proto(
      const bearer_qos_t* eps_bearer_qos_state,
      oai::SgwBearerQos* eps_bearer_qos_proto);

  /**
   * Converts proto to eps bearer QOS struct
   * @param eps_bearer_qos_proto
   * @param eps_bearer_qos_state
   */
  static void proto_to_eps_bearer_qos(
      const oai::SgwBearerQos& eps_bearer_qos_proto,
      bearer_qos_t* eps_bearer_qos_state);

  /**
   * Converts gtpv1u_data struct to proto, memory is owned by the caller
   * @param gtp_data
   * @param gtp_proto
   */
  static void gtpv1u_data_to_proto(
      const gtpv1u_data_t* gtp_data, oai::GTPV1uData* gtp_proto);

  /**
   * Converts proto to gtpv1u_data struct
   * @param gtp_proto
   * @param gtp_data
   */
  static void proto_to_gtpv1u_data(
      const oai::GTPV1uData& gtp_proto, gtpv1u_data_t* gtp_data);

  /**
   * Converts port range struct to proto, memory is owned by the caller
   * @param port_range
   * @param port_range_proto
   */
  static void port_range_to_proto(
      const port_range_t* port_range, oai::PortRange* port_range_proto);

  /**
   * Converts proto to port range struct
   * @param port_range_proto
   * @param port_range
   */
  static void proto_to_port_range(
      const oai::PortRange& port_range_proto, port_range_t* port_range);

  /**
   * Converts packet filter struct to proto, memory is owned by the caller
   * @param packet_filter
   * @param packet_filter_proto
   */
  static void packet_filter_to_proto(
      const packet_filter_t* packet_filter,
      oai::PacketFilter* packet_filter_proto);

  /**
   * Converts proto to packet filter struct
   * @param packet_filter_proto
   * @param packet_filter
   */
  static void proto_to_packet_filter(
      const oai::PacketFilter& packet_filter_proto,
      packet_filter_t* packet_filter);
};
}  // namespace lte
}  // namespace magma
