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

#include "spgw_state_converter.h"

using magma::lte::gateway::spgw::BearerQos;
using magma::lte::gateway::spgw::CreateSessionMessage;
using magma::lte::gateway::spgw::GTPV1uData;
using magma::lte::gateway::spgw::MmeSgwTunnel;
using magma::lte::gateway::spgw::PacketFilter;
using magma::lte::gateway::spgw::Parameter;
using magma::lte::gateway::spgw::PccRule;
using magma::lte::gateway::spgw::PgwCbrProcedure;
using magma::lte::gateway::spgw::PgwState;
using magma::lte::gateway::spgw::PortRange;
using magma::lte::gateway::spgw::S11BearerContext;
using magma::lte::gateway::spgw::SgwEpsBearerContext;
using magma::lte::gateway::spgw::SgwEpsBearerContextInfo;
using magma::lte::gateway::spgw::SgwPdnConnection;
using magma::lte::gateway::spgw::SgwState;
using magma::lte::gateway::spgw::SpgwState;
using magma::lte::gateway::spgw::TrafficFlowTemplate;

namespace magma {
namespace lte {

SpgwStateConverter::SpgwStateConverter() = default;
SpgwStateConverter::~SpgwStateConverter() = default;

void SpgwStateConverter::spgw_state_to_proto(const spgw_state_t* spgw_state,
                                             SpgwState* proto) {
  proto->Clear();

  sgw_state_to_proto(&spgw_state->sgw_state, proto->mutable_sgw_state());
  pgw_state_to_proto(&spgw_state->pgw_state, proto->mutable_pgw_state());
}

/**********************************************************/
/*                SGW State <-> Proto                    */
/**********************************************************/
void SpgwStateConverter::sgw_state_to_proto(const sgw_state_t* sgw_state,
                                            SgwState* proto) {
  proto->Clear();

  s11teid_mme_ht_to_proto(sgw_state->s11teid2mme, proto->mutable_s11teid_mme());
  s11bearer_context_ht_to_proto(sgw_state->s11_bearer_context_information,
                                proto->mutable_s11_bearer_context_info());

  proto->set_sgw_ip_address_s1u_s12_s4_up(
      sgw_state->sgw_ip_address_S1u_S12_S4_up.s_addr);

  gtpv1u_data_to_proto(&sgw_state->gtpv1u_data, proto->mutable_gtpv1u_data());

  proto->set_last_tunnel_id(sgw_state->tunnel_id);
  proto->set_gtpv1u_teid(sgw_state->gtpv1u_teid);
}

void SpgwStateConverter::mme_sgw_tunnel_to_proto(const mme_sgw_tunnel_t* tunnel,
                                                 MmeSgwTunnel* proto) {
  proto->Clear();

  proto->set_local_teid(tunnel->local_teid);
  proto->set_remote_teid(tunnel->remote_teid);
}

void SpgwStateConverter::s11teid_mme_ht_to_proto(
    hash_table_ts_t* const state_map,
    google::protobuf::Map<unsigned int, MmeSgwTunnel>* proto_map) {
  hashtable_ts_to_proto<mme_sgw_tunnel_t, MmeSgwTunnel>(
      state_map, proto_map, mme_sgw_tunnel_to_proto, LOG_SPGW_APP);
}

void SpgwStateConverter::spgw_bearer_context_to_proto(
    const s_plus_p_gw_eps_bearer_context_information_t* spgw_bearer_state,
    S11BearerContext* spgw_bearer_proto) {
  spgw_bearer_proto->Clear();

  auto* sgw_eps_bearer_proto =
      spgw_bearer_proto->mutable_sgw_eps_bearer_context();
  auto* sgw_eps_bearer_state =
      &spgw_bearer_state->sgw_eps_bearer_context_information;

  sgw_eps_bearer_proto->set_imsi((char*)sgw_eps_bearer_state->imsi.digit);
  sgw_eps_bearer_proto->set_imsi_unauth_indicator(
      sgw_eps_bearer_state->imsi_unauthenticated_indicator);
  sgw_eps_bearer_proto->set_msisdn(sgw_eps_bearer_state->msisdn);
  ecgi_to_proto(sgw_eps_bearer_state->last_known_cell_Id,
                sgw_eps_bearer_proto->mutable_last_known_cell_id());

  if (sgw_eps_bearer_state->trxn != nullptr) {
    sgw_eps_bearer_proto->set_trxn((char*)sgw_eps_bearer_state->trxn);
  }

  sgw_eps_bearer_proto->set_mme_teid_s11(sgw_eps_bearer_state->mme_teid_S11);
  sgw_eps_bearer_proto->set_mme_ip_address_s11(
      bdata(ip_address_to_bstring(&sgw_eps_bearer_state->mme_ip_address_S11)));

  sgw_eps_bearer_proto->set_sgw_teid_s11_s4(
      sgw_eps_bearer_state->s_gw_teid_S11_S4);
  sgw_eps_bearer_proto->set_sgw_ip_address_s11_s4(bdata(
      ip_address_to_bstring(&sgw_eps_bearer_state->s_gw_ip_address_S11_S4)));

  sgw_pdn_connection_to_proto(&sgw_eps_bearer_state->pdn_connection,
                              sgw_eps_bearer_proto->mutable_pdn_connection());

  sgw_create_session_message_to_proto(
      &sgw_eps_bearer_state->saved_message,
      sgw_eps_bearer_proto->mutable_saved_message());
  sgw_pending_procedures_to_proto(sgw_eps_bearer_state->pending_procedures,
                                  sgw_eps_bearer_proto);

  auto* pgw_eps_bearer_proto =
      spgw_bearer_proto->mutable_pgw_eps_bearer_context();
  auto* pgw_eps_bearer_state =
      &spgw_bearer_state->pgw_eps_bearer_context_information;

  pgw_eps_bearer_proto->set_imsi((char*)pgw_eps_bearer_state->imsi.digit);
  pgw_eps_bearer_proto->set_imsi_unauth_indicator(
      pgw_eps_bearer_state->imsi_unauthenticated_indicator);
  pgw_eps_bearer_proto->set_msisdn(pgw_eps_bearer_state->msisdn);
}

void SpgwStateConverter::s11bearer_context_ht_to_proto(
    hash_table_ts_t* const state_map,
    google::protobuf::Map<unsigned int, S11BearerContext>* proto_map) {
  hashtable_ts_to_proto<s_plus_p_gw_eps_bearer_context_information_t,
                        S11BearerContext>(
      state_map, proto_map, spgw_bearer_context_to_proto, LOG_SPGW_APP);
}

void SpgwStateConverter::sgw_pdn_connection_to_proto(
    const sgw_pdn_connection_t* state_pdn, SgwPdnConnection* proto_pdn) {
  proto_pdn->Clear();

  proto_pdn->set_apn_in_use(
      strndup(state_pdn->apn_in_use, strlen(state_pdn->apn_in_use)));
  proto_pdn->set_pgw_address_in_use_cp(
      (char*)ip_address_to_bstring(&state_pdn->p_gw_address_in_use_cp)->data);
  proto_pdn->set_pgw_address_in_use_up(
      (char*)ip_address_to_bstring(&state_pdn->p_gw_address_in_use_up)->data);
  proto_pdn->set_default_bearer(state_pdn->default_bearer);
  proto_pdn->set_ue_suspended_for_ps_handover(
      state_pdn->ue_suspended_for_ps_handover);

  for (auto& eps_bearer : state_pdn->sgw_eps_bearers_array) {
    auto* proto_eps_bearer = proto_pdn->add_eps_bearer_list();
    if (eps_bearer != nullptr) {
      sgw_eps_bearer_to_proto(eps_bearer, proto_eps_bearer);
    }
  }
}

void SpgwStateConverter::sgw_create_session_message_to_proto(
    const itti_s11_create_session_request_t* session_request,
    CreateSessionMessage* proto) {
  proto->Clear();

  if (session_request->trxn != nullptr) {
    proto->set_trxn((char*)session_request->trxn);
  }

  proto->set_teid(session_request->teid);
  proto->set_imsi((char*)session_request->imsi.digit);
  proto->set_msisdn((char*)session_request->msisdn.digit);

  if (MEI_IMEISV) {
    memcpy(proto->mutable_mei(), &session_request->mei.choice.imeisv,
           session_request->mei.choice.imeisv.length);
  } else if (MEI_IMEI) {
    memcpy(proto->mutable_mei(), &session_request->mei.choice.imei,
           session_request->mei.choice.imei.length);
  }

  if (session_request->uli.present) {
    char uli[sizeof(Uli_t)];
    memcpy(&uli, &session_request->uli, sizeof(Uli_t));
    proto->set_uli(uli);
  }

  proto->mutable_serving_network()->set_mcc(
      (char*)session_request->serving_network.mcc, 3);
  proto->mutable_serving_network()->set_mnc(
      (char*)session_request->serving_network.mnc, 3);

  proto->set_rat_type(session_request->rat_type);
  proto->set_pdn_type(session_request->pdn_type);
  proto->mutable_ambr()->set_br_ul(session_request->ambr.br_ul);
  proto->mutable_ambr()->set_br_dl(session_request->ambr.br_dl);

  proto->set_apn(session_request->apn, APN_MAX_LENGTH + 1);
  proto->set_paa(bdata(paa_to_bstring(&session_request->paa)));
  proto->set_peer_ip(session_request->peer_ip.s_addr);

  proto->mutable_pco()->set_ext(session_request->pco.ext);
  proto->mutable_pco()->set_spare(session_request->pco.spare);
  proto->mutable_pco()->set_configuration_protocol(
      session_request->pco.configuration_protocol);
  proto->mutable_pco()->set_num_protocol_or_container_id(
      session_request->pco.num_protocol_or_container_id);

  if (session_request->sender_fteid_for_cp.ipv4) {
    proto->mutable_sender_fteid_for_cp()->set_ipv4_address(
        session_request->sender_fteid_for_cp.ipv4_address.s_addr);
  } else if (session_request->sender_fteid_for_cp.ipv6) {
    memcpy(proto->mutable_sender_fteid_for_cp()->mutable_ipv6_address(),
           &session_request->sender_fteid_for_cp.ipv6_address, 16);
  }

  proto->mutable_sender_fteid_for_cp()->set_interface_type(
      session_request->sender_fteid_for_cp.interface_type);
  proto->mutable_sender_fteid_for_cp()->set_teid(
      session_request->sender_fteid_for_cp.teid);

  proto->mutable_ue_time_zone()->set_time_zone(
      session_request->ue_time_zone.time_zone);
  proto->mutable_ue_time_zone()->set_daylight_saving_time(
      session_request->ue_time_zone.daylight_saving_time);

  for (uint32_t i = 0; i < session_request->pco.num_protocol_or_container_id;
       i++) {
    auto* pco_protocol = &session_request->pco.protocol_or_container_ids[i];
    auto* pco_protocol_proto = proto->mutable_pco()->add_pco_protocol();
    if (pco_protocol->contents) {
      pco_protocol_proto->set_id(pco_protocol->id);
      pco_protocol_proto->set_length(pco_protocol->length);
      pco_protocol_proto->set_contents(bdata(pco_protocol->contents));
    }
  }
  for (uint32_t i = 0;
       i < session_request->bearer_contexts_to_be_created.num_bearer_context;
       i++) {
    auto* bearer =
        &session_request->bearer_contexts_to_be_created.bearer_contexts[i];
    auto* bearer_proto = proto->add_bearer_contexts_to_be_created();
    bearer_proto->set_eps_bearer_id(bearer->eps_bearer_id);
    traffic_flow_template_to_proto(&bearer->tft, bearer_proto->mutable_tft());
    eps_bearer_qos_to_proto(&bearer->bearer_level_qos,
                            bearer_proto->mutable_bearer_level_qos());
  }
}

void SpgwStateConverter::sgw_eps_bearer_to_proto(
    const sgw_eps_bearer_ctxt_t* eps_bearer,
    SgwEpsBearerContext* eps_bearer_proto) {
  eps_bearer_proto->Clear();

  eps_bearer_proto->set_eps_bearer_id(eps_bearer->eps_bearer_id);

  eps_bearer_proto->set_pgw_address_in_use_up(
      bdata(ip_address_to_bstring(&eps_bearer->p_gw_address_in_use_up)));
  eps_bearer_proto->set_pgw_teid_s5_s8_up(eps_bearer->p_gw_teid_S5_S8_up);

  eps_bearer_proto->set_sgw_ip_address_s5_s8_up(
      bdata(ip_address_to_bstring(&eps_bearer->s_gw_ip_address_S5_S8_up)));
  eps_bearer_proto->set_sgw_teid_s5_s8_up(eps_bearer->s_gw_teid_S5_S8_up);

  eps_bearer_proto->set_sgw_ip_address_s1u_s12_s4_up(
    bdata(
      ip_address_to_bstring(&eps_bearer->s_gw_ip_address_S1u_S12_S4_up)));
  eps_bearer_proto->set_sgw_teid_s1u_s12_s4_up(
    eps_bearer->s_gw_teid_S1u_S12_S4_up);

  eps_bearer_proto->set_enb_ip_address_s1u(
    bdata(ip_address_to_bstring(&eps_bearer->enb_ip_address_S1u)));
  eps_bearer_proto->set_enb_teid_s1u(eps_bearer->enb_teid_S1u);
  eps_bearer_proto->set_paa(
    bdata(paa_to_bstring(&eps_bearer->paa)));

  eps_bearer_qos_to_proto(
    &eps_bearer->eps_bearer_qos, eps_bearer_proto->mutable_eps_bearer_qos());
  traffic_flow_template_to_proto(
    &eps_bearer->tft, eps_bearer_proto->mutable_tft());

  eps_bearer_proto->set_num_sdf(eps_bearer->num_sdf);

  for (const auto &sdf_id : eps_bearer->sdf_id) {
    eps_bearer_proto->add_sdf_ids(sdf_id);
  }
}

void SpgwStateConverter::traffic_flow_template_to_proto(
    const traffic_flow_template_t* tft_state, TrafficFlowTemplate* tft_proto) {
  tft_proto->Clear();

  tft_proto->set_tft_operation_code(tft_state->tftoperationcode);
  tft_proto->set_number_of_packet_filters(tft_state->numberofpacketfilters);
  tft_proto->set_ebit(tft_state->ebit);

  // parameters_list member conversion
  tft_proto->mutable_parameters_list()->set_num_parameters(
      tft_state->parameterslist.num_parameters);
  for (uint32_t i = 0; i < tft_state->parameterslist.num_parameters; i++) {
    auto* parameter = &tft_state->parameterslist.parameter[i];
    if (parameter->contents) {
      auto* param_proto =
          tft_proto->mutable_parameters_list()->add_parameters();
      param_proto->set_parameter_identifier(parameter->parameteridentifier);
      param_proto->set_length(parameter->length);
      param_proto->set_contents(bdata(parameter->contents));
    }
  }

  // traffic_flow_template.packet_filter list member conversions
  auto* pft_proto = tft_proto->mutable_packet_filter_list();
  auto pft_state = tft_state->packetfilterlist;
  switch (tft_state->tftoperationcode) {
  case TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET_FILTERS_FROM_EXISTING_TFT:
    for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
      pft_proto->add_delete_packet_filter_identifier(
          pft_state.deletepacketfilter[i].identifier);
    }
    break;
  case TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT:
    for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
      packet_filter_to_proto(&pft_state.createnewtft[i],
                             pft_proto->add_create_new_tft());
    }
    break;
  case TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET_FILTER_TO_EXISTING_TFT:
    for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
      packet_filter_to_proto(&pft_state.createnewtft[i],
                             pft_proto->add_add_packet_filter());
    }
    break;
  case TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET_FILTERS_IN_EXISTING_TFT:
    for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
      packet_filter_to_proto(&pft_state.createnewtft[i],
                             pft_proto->add_replace_packet_filter());
    }
    break;
  default:
    break;
  };
}

void SpgwStateConverter::port_range_to_proto(const port_range_t* port_range,
                                             PortRange* port_range_proto) {
  port_range_proto->Clear();

  port_range_proto->set_low_limit(port_range->lowlimit);
  port_range_proto->set_high_limit(port_range->highlimit);
}

void SpgwStateConverter::packet_filter_to_proto(
    const packet_filter_t* packet_filter, PacketFilter* packet_filter_proto) {
  packet_filter_proto->Clear();

  packet_filter_proto->set_spare(packet_filter->spare);
  packet_filter_proto->set_direction(packet_filter->direction);
  packet_filter_proto->set_identifier(packet_filter->identifier);
  packet_filter_proto->set_eval_precedence(packet_filter->eval_precedence);

  auto* packet_filter_contents =
      packet_filter_proto->mutable_packet_filter_contents();
  packet_filter_contents->set_flags(packet_filter->packetfiltercontents.flags);
  packet_filter_contents->set_protocol_identifier_nextheader(
      packet_filter->packetfiltercontents.protocolidentifier_nextheader);
  packet_filter_contents->set_single_local_port(
      packet_filter->packetfiltercontents.singlelocalport);
  packet_filter_contents->set_single_remote_port(
      packet_filter->packetfiltercontents.singleremoteport);
  packet_filter_contents->set_security_parameter_index(
      packet_filter->packetfiltercontents.securityparameterindex);
  packet_filter_contents->set_flow_label(
      packet_filter->packetfiltercontents.flowlabel);

  for (auto& ip : packet_filter->packetfiltercontents.ipv4remoteaddr) {
    auto* ipv4_proto = packet_filter_contents->add_ipv4_remote_addresses();
    ipv4_proto->set_addr(ip.addr);
    ipv4_proto->set_mask(ip.mask);
  }

  for (auto& ip : packet_filter->packetfiltercontents.ipv6remoteaddr) {
    auto* ipv6_proto = packet_filter_contents->add_ipv6_remote_addresses();
    ipv6_proto->set_addr(ip.addr);
    ipv6_proto->set_mask(ip.mask);
  }

  port_range_to_proto(&packet_filter->packetfiltercontents.localportrange,
                      packet_filter_contents->mutable_local_port_range());
  port_range_to_proto(&packet_filter->packetfiltercontents.remoteportrange,
                      packet_filter_contents->mutable_remote_port_range());

  packet_filter_contents->mutable_type_of_service_traffic_class()->set_value(
      packet_filter->packetfiltercontents.typdeofservice_trafficclass.value);
  packet_filter_contents->mutable_type_of_service_traffic_class()->set_mask(
      packet_filter->packetfiltercontents.typdeofservice_trafficclass.mask);
}

void SpgwStateConverter::eps_bearer_qos_to_proto(
    const bearer_qos_t* eps_bearer_qos_state, BearerQos* eps_bearer_qos_proto) {
  eps_bearer_qos_proto->Clear();

  eps_bearer_qos_proto->set_pci(eps_bearer_qos_state->pci);
  eps_bearer_qos_proto->set_pl(eps_bearer_qos_state->pl);
  eps_bearer_qos_proto->set_pvi(eps_bearer_qos_state->pvi);
  eps_bearer_qos_proto->set_qci(eps_bearer_qos_state->qci);

  eps_bearer_qos_proto->mutable_gbr()->set_br_ul(
      eps_bearer_qos_state->gbr.br_ul);
  eps_bearer_qos_proto->mutable_gbr()->set_br_dl(
      eps_bearer_qos_state->gbr.br_dl);

  eps_bearer_qos_proto->mutable_mbr()->set_br_ul(
      eps_bearer_qos_state->mbr.br_ul);
  eps_bearer_qos_proto->mutable_mbr()->set_br_dl(
      eps_bearer_qos_state->mbr.br_dl);
}

void SpgwStateConverter::gtpv1u_data_to_proto(const gtpv1u_data_t* gtp_data,
                                              GTPV1uData* gtp_proto) {
  gtp_proto->Clear();

  if (gtp_data->ip_addr != nullptr) {
    gtp_proto->set_ip_address(gtp_data->ip_addr);
  }

  gtp_proto->set_seq_num(gtp_data->seq_num);
  gtp_proto->set_restart_counter(gtp_data->restart_counter);
  gtp_proto->set_fd0(gtp_data->fd0);
  gtp_proto->set_fd1u(gtp_data->fd1u);
}

/**********************************************************/
/*                PGW State <-> Proto                    */
/**********************************************************/

void SpgwStateConverter::pgw_state_to_proto(const pgw_state_t* pgw_state,
                                            PgwState* proto) {
  proto->Clear();

  if (pgw_state->predefined_pcc_rules != nullptr) {
    pcc_rule_ht_to_proto(pgw_state->predefined_pcc_rules,
                         proto->mutable_predefined_pcc_rules());
  }

  if (pgw_state->deactivated_predefined_pcc_rules != nullptr) {
    pcc_rule_ht_to_proto(pgw_state->deactivated_predefined_pcc_rules,
                         proto->mutable_deactivated_predefined_pcc_rules());
  }

}

void SpgwStateConverter::pcc_rule_ht_to_proto(
    hash_table_ts_t* const state_map,
    google::protobuf::Map<unsigned int, PccRule>* proto_map) {
  hashtable_ts_to_proto<pcc_rule_t, PccRule>(state_map, proto_map,
                                             pcc_rule_to_proto, LOG_SPGW_APP);
}

void SpgwStateConverter::pcc_rule_to_proto(const pcc_rule_t* pcc_rule_state,
                                           PccRule* proto) {
  proto->Clear();

  proto->set_name(bdata(pcc_rule_state->name));
  proto->set_is_activated(pcc_rule_state->is_activated);
  proto->set_sdf_id((unsigned int)pcc_rule_state->sdf_id);
  proto->set_precedence(pcc_rule_state->precedence);

  eps_bearer_qos_to_proto(&pcc_rule_state->bearer_qos,
                          proto->mutable_bearer_qos());

  proto->mutable_sdf_template()->set_number_of_packet_filters(
      pcc_rule_state->sdf_template.number_of_packet_filters);
  for (uint32_t i = 0;
       i < pcc_rule_state->sdf_template.number_of_packet_filters; i++) {
    packet_filter_to_proto(&pcc_rule_state->sdf_template.sdf_filter[i],
                           proto->mutable_sdf_template()->add_sdf_filter());
  }
}

void SpgwStateConverter::sgw_pending_procedures_to_proto(
    const sgw_eps_bearer_context_information_t::pending_procedures_s*
        procedures,
    SgwEpsBearerContextInfo* proto) {
  if (procedures != nullptr) {
    pgw_base_proc_t* base_proc = nullptr;

    LIST_FOREACH(base_proc, procedures, entries) {
      if (base_proc->type ==
          PGW_BASE_PROC_TYPE_NETWORK_INITATED_CREATE_BEARER_REQUEST) {
        auto* create_proc = (pgw_ni_cbr_proc_t*)base_proc;
        auto* cbr_procedure_proto = proto->add_pending_procedures();
        cbr_procedure_proto->set_teid(create_proc->teid);
        cbr_procedure_proto->set_sdf_id(create_proc->sdf_id);
        sgw_eps_bearer_entry_wrapper_t* b1 = nullptr;
        LIST_FOREACH(b1, create_proc->pending_eps_bearers, entries) {
          sgw_eps_bearer_to_proto(
              b1->sgw_eps_bearer_entry,
              cbr_procedure_proto->add_pending_eps_bearers());
        }
      }
    }
  }
}

} // namespace lte
} // namespace magma
