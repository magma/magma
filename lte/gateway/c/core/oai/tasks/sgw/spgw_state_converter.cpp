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

#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_converter.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/common/conversions.h"
}

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/lib/pcef/pcef_handlers.hpp"
#include "lte/gateway/c/core/oai/include/sgw_context_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_handlers.hpp"

using magma::lte::oai::CreateSessionMessage;
using magma::lte::oai::GTPV1uData;
using magma::lte::oai::PacketFilter;
using magma::lte::oai::PgwCbrProcedure;
using magma::lte::oai::S11BearerContext;
using magma::lte::oai::SgwBearerQos;
using magma::lte::oai::SgwEpsBearerContext;
using magma::lte::oai::SgwEpsBearerContextInfo;
using magma::lte::oai::SgwPdnConnection;
using magma::lte::oai::SpgwState;
using magma::lte::oai::SpgwUeContext;
using magma::lte::oai::TrafficFlowTemplate;

namespace magma {
namespace lte {

SpgwStateConverter::SpgwStateConverter() = default;
SpgwStateConverter::~SpgwStateConverter() = default;

void SpgwStateConverter::state_to_proto(const spgw_state_t* spgw_state,
                                        SpgwState* proto) {
  proto->Clear();

  gtpv1u_data_to_proto(&spgw_state->gtpv1u_data, proto->mutable_gtpv1u_data());

  proto->set_gtpv1u_teid(spgw_state->gtpv1u_teid);
}

void SpgwStateConverter::proto_to_state(const SpgwState& proto,
                                        spgw_state_t* spgw_state) {
  proto_to_gtpv1u_data(proto.gtpv1u_data(), &spgw_state->gtpv1u_data);
  spgw_state->gtpv1u_teid = proto.gtpv1u_teid();
}

void SpgwStateConverter::spgw_bearer_context_to_proto(
    const S11BearerContext* spgw_bearer_state,
    S11BearerContext* spgw_bearer_proto) {
  spgw_bearer_proto->Clear();
  spgw_bearer_proto->MergeFrom(*spgw_bearer_state);
}

void SpgwStateConverter::proto_to_spgw_bearer_context(
    const S11BearerContext& spgw_bearer_proto,
    S11BearerContext* spgw_bearer_state) {
  spgw_bearer_state->Clear();
  spgw_bearer_state->MergeFrom(spgw_bearer_proto);
}

void SpgwStateConverter::port_range_to_proto(const port_range_t* port_range,
                                             oai::PortRange* port_range_proto) {
  port_range_proto->Clear();

  port_range_proto->set_low_limit(port_range->lowlimit);
  port_range_proto->set_high_limit(port_range->highlimit);
}

void SpgwStateConverter::proto_to_port_range(
    const oai::PortRange& port_range_proto, port_range_t* port_range) {
  port_range->lowlimit = port_range_proto.low_limit();
  port_range->highlimit = port_range_proto.high_limit();
}

void SpgwStateConverter::proto_to_packet_filter(
    const oai::PacketFilter& packet_filter_proto,
    packet_filter_t* packet_filter) {
  packet_filter->spare = packet_filter_proto.spare();
  packet_filter->direction = packet_filter_proto.direction();
  packet_filter->identifier = packet_filter_proto.identifier();
  packet_filter->eval_precedence = packet_filter_proto.eval_precedence();
  packet_filter->length = packet_filter_proto.length();

  auto* packet_filter_contents = &packet_filter->packetfiltercontents;
  for (uint8_t pf_idx = 0;
       pf_idx < packet_filter_proto.packet_filter_contents_size(); pf_idx++) {
    auto& packet_filter_content_proto =
        packet_filter_proto.packet_filter_contents(pf_idx);
    switch (packet_filter_content_proto.flags()) {
      case TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR: {
        if (packet_filter_content_proto.ipv4_remote_addresses_size()) {
          packet_filter_contents->flags |=
              TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
          uint8_t local_idx = TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE - 1;
          for (uint8_t idx = 0; idx < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE;
               idx++) {
            packet_filter_contents->ipv4remoteaddr[local_idx].addr =
                packet_filter_content_proto.ipv4_remote_addresses(0).addr() >>
                (idx * 8);
            packet_filter_contents->ipv4remoteaddr[local_idx].mask =
                packet_filter_content_proto.ipv4_remote_addresses(0).mask() >>
                (idx * 8);
            --local_idx;
          }
        }
      } break;
      case TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR: {
        if (packet_filter_content_proto.ipv6_remote_addresses_size()) {
          packet_filter_contents->flags |=
              TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG;
          uint8_t local_idx = TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE - 1;
          for (uint8_t idx = 0; idx < TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE;
               idx++) {
            packet_filter_contents->ipv6remoteaddr[local_idx].addr =
                packet_filter_content_proto.ipv6_remote_addresses(idx).addr();
            packet_filter_contents->ipv6remoteaddr[local_idx].mask =
                packet_filter_content_proto.ipv6_remote_addresses(idx).mask();
            --local_idx;
          }
        }
      } break;
      case TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG;
        packet_filter_contents->protocolidentifier_nextheader =
            packet_filter_content_proto.protocol_identifier_nextheader();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
        packet_filter_contents->singlelocalport =
            packet_filter_content_proto.single_local_port();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
        packet_filter_contents->singleremoteport =
            packet_filter_content_proto.single_remote_port();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG;
        packet_filter_contents->securityparameterindex =
            packet_filter_content_proto.security_parameter_index();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG;
        packet_filter_contents->typdeofservice_trafficclass.value =
            packet_filter_content_proto.type_of_service_traffic_class().value();
        packet_filter_contents->typdeofservice_trafficclass.mask =
            packet_filter_content_proto.type_of_service_traffic_class().mask();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL: {
        packet_filter_contents->flags |= TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG;
        packet_filter_contents->flowlabel =
            packet_filter_content_proto.flow_label();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG;
        proto_to_port_range(packet_filter_content_proto.local_port_range(),
                            &packet_filter_contents->localportrange);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG;
        proto_to_port_range(packet_filter_content_proto.remote_port_range(),
                            &packet_filter_contents->remoteportrange);
      } break;
    }
  }
}

void SpgwStateConverter::proto_to_traffic_flow_template(
    const oai::TrafficFlowTemplate& tft_proto,
    traffic_flow_template_t* tft_state) {
  tft_state->tftoperationcode = tft_proto.tft_operation_code();
  tft_state->ebit = tft_proto.ebit();

  tft_state->parameterslist.num_parameters =
      tft_proto.parameters_list().num_parameters();

  for (uint32_t i = 0; i < tft_proto.parameters_list().num_parameters(); i++) {
    auto* param_state = &tft_state->parameterslist.parameter[i];
    auto& param_proto = tft_proto.parameters_list().parameters(i);
    param_state->parameteridentifier = param_proto.parameter_identifier();
    param_state->length = param_proto.length();
    param_state->contents = bfromcstr(param_proto.contents().c_str());
  }

  auto& pft_proto = tft_proto.packet_filter_list();
  auto* pft_state = &tft_state->packetfilterlist;
  switch (tft_proto.tft_operation_code()) {
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET_FILTERS_FROM_EXISTING_TFT:
      tft_state->numberofpacketfilters =
          pft_proto.delete_packet_filter_identifier_size();
      for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
        pft_state->deletepacketfilter[i].identifier =
            pft_proto.delete_packet_filter_identifier(i);
      }
      break;
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT:
      tft_state->numberofpacketfilters = pft_proto.create_new_tft_size();
      for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
        proto_to_packet_filter(pft_proto.create_new_tft(i),
                               &pft_state->createnewtft[i]);
      }
      break;
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET_FILTER_TO_EXISTING_TFT:
      tft_state->numberofpacketfilters = pft_proto.add_packet_filter_size();
      for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
        proto_to_packet_filter(pft_proto.add_packet_filter(i),
                               &pft_state->addpacketfilter[i]);
      }
      break;
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET_FILTERS_IN_EXISTING_TFT:
      tft_state->numberofpacketfilters = pft_proto.replace_packet_filter_size();
      for (uint32_t i = 0; i < tft_state->numberofpacketfilters; i++) {
        proto_to_packet_filter(pft_proto.replace_packet_filter(i),
                               &pft_state->replacepacketfilter[i]);
      }
      break;
    default:
      break;
  };
}

void SpgwStateConverter::gtpv1u_data_to_proto(const gtpv1u_data_t* gtp_data,
                                              GTPV1uData* gtp_proto) {
  gtp_proto->Clear();
  gtp_proto->set_fd0(gtp_data->fd0);
  gtp_proto->set_fd1u(gtp_data->fd1u);
}

void SpgwStateConverter::proto_to_gtpv1u_data(const oai::GTPV1uData& gtp_proto,
                                              gtpv1u_data_t* gtp_data) {
  gtp_data->fd0 = gtp_proto.fd0();
  gtp_data->fd1u = gtp_proto.fd1u();
}

void SpgwStateConverter::ue_to_proto(const spgw_ue_context_t* ue_state,
                                     oai::SpgwUeContext* ue_proto) {
  if (ue_state && (!LIST_EMPTY(&ue_state->sgw_s11_teid_list))) {
    sgw_s11_teid_t* s11_teid_p = nullptr;
    LIST_FOREACH(s11_teid_p, &ue_state->sgw_s11_teid_list, entries) {
      if (s11_teid_p) {
        auto spgw_ctxt = sgw_cm_get_spgw_context(s11_teid_p->sgw_s11_teid);
        if (spgw_ctxt) {
          spgw_bearer_context_to_proto(spgw_ctxt,
                                       ue_proto->add_s11_bearer_context());
        }
      }
    }
  }
}

void SpgwStateConverter::proto_to_ue(const oai::SpgwUeContext& ue_proto,
                                     spgw_ue_context_t* ue_context_p) {
  OAILOG_FUNC_IN(LOG_SPGW_APP);
  map_uint64_spgw_ue_context_t* state_ue_map = nullptr;
  state_teid_map_t* state_teid_map = nullptr;
  if (ue_proto.s11_bearer_context_size()) {
    state_teid_map = get_spgw_teid_state();
    if (!state_teid_map) {
      OAILOG_ERROR(LOG_SPGW_APP, "Failed to get state_teid_map \n");
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
    }

    state_ue_map = get_spgw_ue_state();
    if (!state_ue_map) {
      OAILOG_ERROR(LOG_SPGW_APP,
                   "Failed to get state_ue_map from get_spgw_ue_state() \n");
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
    }

    // All s11_bearer_context on this UE context will be of same imsi
    imsi64_t imsi64 =
        ue_proto.s11_bearer_context(0).sgw_eps_bearer_context().imsi64();
    if (ue_context_p) {
      LIST_INIT(&ue_context_p->sgw_s11_teid_list);
      state_ue_map->insert(imsi64, ue_context_p);
    } else {
      OAILOG_ERROR_UE(LOG_SPGW_APP, imsi64,
                      "Failed to allocate memory for UE context \n");
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
    }
  } else {
    OAILOG_ERROR(LOG_SPGW_APP,
                 "There are no spgw_context stored to Redis DB \n");
    OAILOG_FUNC_OUT(LOG_SPGW_APP);
  }
  for (int idx = 0; idx < ue_proto.s11_bearer_context_size(); idx++) {
    oai::S11BearerContext S11BearerContext = ue_proto.s11_bearer_context(idx);
    oai::S11BearerContext* spgw_context_p = new oai::S11BearerContext();

    proto_to_spgw_bearer_context(S11BearerContext, spgw_context_p);
    if ((state_teid_map->insert(
             spgw_context_p->sgw_eps_bearer_context().sgw_teid_s11_s4(),
             spgw_context_p) != magma::PROTO_MAP_OK)) {
      OAILOG_ERROR(LOG_SPGW_APP,
                   "Failed to insert spgw_context_p for teid " TEID_FMT " \n",
                   spgw_context_p->sgw_eps_bearer_context().sgw_teid_s11_s4());
      OAILOG_FUNC_OUT(LOG_SPGW_APP);
    }
    spgw_update_teid_in_ue_context(
        spgw_context_p->sgw_eps_bearer_context().imsi64(),
        spgw_context_p->sgw_eps_bearer_context().sgw_teid_s11_s4());
  }
  OAILOG_FUNC_OUT(LOG_SPGW_APP);
}

}  // namespace lte
}  // namespace magma

void port_range_to_proto(const port_range_t* port_range,
                         magma::lte::oai::PortRange* port_range_proto) {
  port_range_proto->Clear();

  port_range_proto->set_low_limit(port_range->lowlimit);
  port_range_proto->set_high_limit(port_range->highlimit);
}

void packet_filter_to_proto(
    const packet_filter_t* packet_filter,
    magma::lte::oai::PacketFilter* packet_filter_proto) {
  packet_filter_proto->Clear();

  packet_filter_proto->set_spare(packet_filter->spare);
  packet_filter_proto->set_direction(packet_filter->direction);
  packet_filter_proto->set_identifier(packet_filter->identifier);
  packet_filter_proto->set_eval_precedence(packet_filter->eval_precedence);

  uint16_t flag = TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
  while (flag <= TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG) {
    auto* pf_content = packet_filter_proto->add_packet_filter_contents();
    switch (packet_filter->packetfiltercontents.flags & flag) {
      case TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR);
        for (auto& ip : packet_filter->packetfiltercontents.ipv4remoteaddr) {
          auto* ipv4_proto = pf_content->add_ipv4_remote_addresses();
          ipv4_proto->set_addr(ip.addr);
          ipv4_proto->set_mask(ip.mask);
        }
      } break;
      case TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR);
        for (auto& ip : packet_filter->packetfiltercontents.ipv6remoteaddr) {
          auto* ipv6_proto = pf_content->add_ipv6_remote_addresses();
          ipv6_proto->set_addr(ip.addr);
          ipv6_proto->set_mask(ip.mask);
        }
      } break;
      case TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER);
        pf_content->set_protocol_identifier_nextheader(
            packet_filter->packetfiltercontents.protocolidentifier_nextheader);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT);
        pf_content->set_single_local_port(
            packet_filter->packetfiltercontents.singlelocalport);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT);
        pf_content->set_single_remote_port(
            packet_filter->packetfiltercontents.singleremoteport);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX);
        pf_content->set_security_parameter_index(
            packet_filter->packetfiltercontents.securityparameterindex);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG: {
        pf_content->set_flags(
            TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS);
        pf_content->mutable_type_of_service_traffic_class()->set_value(
            packet_filter->packetfiltercontents.typdeofservice_trafficclass
                .value);
        pf_content->mutable_type_of_service_traffic_class()->set_mask(
            packet_filter->packetfiltercontents.typdeofservice_trafficclass
                .mask);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL);
        pf_content->set_flow_label(
            packet_filter->packetfiltercontents.flowlabel);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE);
        port_range_to_proto(&packet_filter->packetfiltercontents.localportrange,
                            pf_content->mutable_local_port_range());
      } break;
      case TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG: {
        pf_content->set_flags(TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE);
        port_range_to_proto(
            &packet_filter->packetfiltercontents.remoteportrange,
            pf_content->mutable_remote_port_range());
      } break;
      default:
        break;
    }
    flag = flag << 1;
  }
}

void eps_bearer_qos_to_proto(
    const bearer_qos_t* eps_bearer_qos_state,
    magma::lte::oai::SgwBearerQos* eps_bearer_qos_proto) {
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
void proto_to_port_range(const magma::lte::oai::PortRange& port_range_proto,
                         port_range_t* port_range) {
  port_range->lowlimit = port_range_proto.low_limit();
  port_range->highlimit = port_range_proto.high_limit();
}

void proto_to_packet_filter(
    const magma::lte::oai::PacketFilter& packet_filter_proto,
    packet_filter_t* packet_filter) {
  packet_filter->spare = packet_filter_proto.spare();
  packet_filter->direction = packet_filter_proto.direction();
  packet_filter->identifier = packet_filter_proto.identifier();
  packet_filter->eval_precedence = packet_filter_proto.eval_precedence();
  packet_filter->length = packet_filter_proto.length();

  auto* packet_filter_contents = &packet_filter->packetfiltercontents;
  for (int32_t i = 0; i < packet_filter_proto.packet_filter_contents_size();
       i++) {
    auto& packet_filter_content_proto =
        packet_filter_proto.packet_filter_contents(i);
    switch (packet_filter_content_proto.flags()) {
      case TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR: {
        if (packet_filter_content_proto.ipv4_remote_addresses_size()) {
          packet_filter_contents->flags |=
              TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
          int local_idx = TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE - 1;
          for (int i = 0; i < TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE; i++) {
            packet_filter_contents->ipv4remoteaddr[local_idx].addr =
                packet_filter_content_proto.ipv4_remote_addresses(i).addr();
            packet_filter_contents->ipv4remoteaddr[local_idx].mask =
                packet_filter_content_proto.ipv4_remote_addresses(0).mask();
            --local_idx;
          }
        }
      } break;
      case TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR: {
        if (packet_filter_content_proto.ipv6_remote_addresses_size()) {
          packet_filter_contents->flags |=
              TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG;
          for (int i = 0; i < TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE; i++) {
            packet_filter_contents->ipv6remoteaddr[i].addr =
                packet_filter_content_proto.ipv6_remote_addresses(i).addr();
            packet_filter_contents->ipv6remoteaddr[i].mask =
                packet_filter_content_proto.ipv6_remote_addresses(i).mask();
          }
        }
      } break;
      case TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG;
        packet_filter_contents->protocolidentifier_nextheader =
            packet_filter_content_proto.protocol_identifier_nextheader();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
        packet_filter_contents->singlelocalport =
            packet_filter_content_proto.single_local_port();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
        packet_filter_contents->singleremoteport =
            packet_filter_content_proto.single_remote_port();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_SECURITY_PARAMETER_INDEX_FLAG;
        packet_filter_contents->securityparameterindex =
            packet_filter_content_proto.security_parameter_index();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_TYPE_OF_SERVICE_TRAFFIC_CLASS_FLAG;
        packet_filter_contents->typdeofservice_trafficclass.value =
            packet_filter_content_proto.type_of_service_traffic_class().value();
        packet_filter_contents->typdeofservice_trafficclass.mask =
            packet_filter_content_proto.type_of_service_traffic_class().mask();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL: {
        packet_filter_contents->flags |= TRAFFIC_FLOW_TEMPLATE_FLOW_LABEL_FLAG;
        packet_filter_contents->flowlabel =
            packet_filter_content_proto.flow_label();
      } break;
      case TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_LOCAL_PORT_RANGE_FLAG;
        proto_to_port_range(packet_filter_content_proto.local_port_range(),
                            &packet_filter_contents->localportrange);
      } break;
      case TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE: {
        packet_filter_contents->flags |=
            TRAFFIC_FLOW_TEMPLATE_REMOTE_PORT_RANGE_FLAG;
        proto_to_port_range(packet_filter_content_proto.remote_port_range(),
                            &packet_filter_contents->remoteportrange);
      } break;
    }
  }
}

void traffic_flow_template_to_proto(
    const traffic_flow_template_t* tft_state,
    magma::lte::oai::TrafficFlowTemplate* tft_proto) {
  tft_proto->Clear();

  tft_proto->set_tft_operation_code(tft_state->tftoperationcode);
  tft_proto->set_number_of_packet_filters(tft_state->numberofpacketfilters);
  tft_proto->set_ebit(tft_state->ebit);

  // parameters_list member conversion
  tft_proto->mutable_parameters_list()->set_num_parameters(
      tft_state->parameterslist.num_parameters);
  for (int idx = 0; idx < tft_state->parameterslist.num_parameters; idx++) {
    auto* parameter = &tft_state->parameterslist.parameter[idx];
    if (parameter->contents) {
      auto* param_proto =
          tft_proto->mutable_parameters_list()->add_parameters();
      param_proto->set_parameter_identifier(parameter->parameteridentifier);
      param_proto->set_length(parameter->length);
      BSTRING_TO_STRING(parameter->contents, param_proto->mutable_contents());
    }
  }

  // traffic_flow_template.packet_filter list member conversions
  auto* pft_proto = tft_proto->mutable_packet_filter_list();
  auto pft_state = tft_state->packetfilterlist;
  switch (tft_state->tftoperationcode) {
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_PACKET_FILTERS_FROM_EXISTING_TFT:
      for (int idx = 0; idx < tft_state->numberofpacketfilters; idx++) {
        pft_proto->add_delete_packet_filter_identifier(
            pft_state.deletepacketfilter[idx].identifier);
      }
      break;
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT:
      for (int idx = 0; idx < tft_state->numberofpacketfilters; idx++) {
        packet_filter_to_proto(&pft_state.createnewtft[idx],
                               pft_proto->add_create_new_tft());
      }
      break;
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_ADD_PACKET_FILTER_TO_EXISTING_TFT:
      for (int idx = 0; idx < tft_state->numberofpacketfilters; idx++) {
        packet_filter_to_proto(&pft_state.createnewtft[idx],
                               pft_proto->add_add_packet_filter());
      }
      break;
    case TRAFFIC_FLOW_TEMPLATE_OPCODE_REPLACE_PACKET_FILTERS_IN_EXISTING_TFT:
      for (int idx = 0; idx < tft_state->numberofpacketfilters; idx++) {
        packet_filter_to_proto(&pft_state.createnewtft[idx],
                               pft_proto->add_replace_packet_filter());
      }
      break;
    default:
      OAILOG_ERROR(LOG_SPGW_APP, "Invalid TFT operation code:%u ",
                   tft_state->tftoperationcode);
      break;
  }
}

void convert_serving_network_to_proto(
    ServingNetwork_t serving_nw,
    magma::lte::oai::ServingNetwork* serving_nw_proto) {
  char mcc[3] = {0};
  char mnc[3] = {0};
  uint8_t mnc_len = 0;

  mcc[0] = convert_digit_to_char(serving_nw.mcc[0]);
  mcc[1] = convert_digit_to_char(serving_nw.mcc[1]);
  mcc[2] = convert_digit_to_char(serving_nw.mcc[2]);
  mnc[0] = convert_digit_to_char(serving_nw.mnc[0]);
  mnc[1] = convert_digit_to_char(serving_nw.mnc[1]);
  if ((serving_nw.mnc[2] & 0xf) != 0xf) {
    mnc[2] = convert_digit_to_char(serving_nw.mnc[2]);
    mnc_len = 3;
  } else {
    mnc_len = 2;
  }
  serving_nw_proto->set_mcc(mcc, 3);
  serving_nw_proto->set_mnc(mnc, mnc_len);
}

void sgw_create_session_message_to_proto(
    const itti_s11_create_session_request_t* session_request,
    magma::lte::oai::CreateSessionMessage* proto) {
  proto->Clear();

  if (session_request->trxn != nullptr) {
    proto->set_trxn(reinterpret_cast<char*>(session_request->trxn));
  }

  proto->set_teid(session_request->teid);
  proto->set_imsi(reinterpret_cast<const char*>(session_request->imsi.digit));
  char msisdn[MSISDN_LENGTH + 1] = {};
  uint32_t msisdn_len = get_msisdn_from_session_req(session_request, msisdn);
  proto->set_msisdn(msisdn);

  char imeisv[IMEISV_DIGITS_MAX + 1] = {};
  if (get_imeisv_from_session_req(session_request, imeisv)) {
    convert_imeisv_to_string(imeisv);
    proto->set_mei(imeisv, IMEISV_DIGITS_MAX);
    OAILOG_DEBUG(LOG_SPGW_APP, "imeisv:%s \n", imeisv);
  }

  char uli[14] = {};
  bool uli_exists = get_uli_from_session_req(session_request, uli);
  if (uli_exists) {
    proto->set_uli(uli, 14);
  }

  const auto cc = session_request->charging_characteristics;
  if (cc.length > 0) {
    proto->set_charging_characteristics(cc.value, cc.length);
  }

  char mcc_mnc[7] = {};
  get_plmn_from_session_req(session_request, mcc_mnc);
  convert_serving_network_to_proto(session_request->serving_network,
                                   proto->mutable_serving_network());

  proto->set_rat_type(session_request->rat_type);
  proto->set_pdn_type(session_request->pdn_type);
  proto->mutable_ambr()->set_br_ul(session_request->ambr.br_ul);
  proto->mutable_ambr()->set_br_dl(session_request->ambr.br_dl);

  proto->set_apn(session_request->apn, strlen(session_request->apn));
  proto->mutable_ue_ip_paa()->set_pdn_type(session_request->paa.pdn_type);

  if (session_request->paa.pdn_type == IPv4) {
    char ip_str[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &(session_request->paa.ipv4_address.s_addr), ip_str,
              INET_ADDRSTRLEN);
    proto->mutable_ue_ip_paa()->set_ipv4_addr(ip_str);
  } else if (session_request->paa.pdn_type == IPv6) {
    char ip6_str[INET6_ADDRSTRLEN];
    inet_ntop(AF_INET6, &(session_request->paa.ipv6_address), ip6_str,
              INET6_ADDRSTRLEN);
    proto->mutable_ue_ip_paa()->set_ipv6_addr(ip6_str);
    proto->mutable_ue_ip_paa()->set_ipv6_prefix_length(
        session_request->paa.ipv6_prefix_length);
    proto->mutable_ue_ip_paa()->set_vlan(session_request->paa.vlan);
  } else if (session_request->paa.pdn_type == IPv4_AND_v6) {
    char ip4_str[INET_ADDRSTRLEN] = {};
    inet_ntop(AF_INET, &(session_request->paa.ipv4_address.s_addr), ip4_str,
              INET_ADDRSTRLEN);
    char ip6_str[INET6_ADDRSTRLEN] = {};
    inet_ntop(AF_INET6, &(session_request->paa.ipv6_address), ip6_str,
              INET6_ADDRSTRLEN);
    proto->mutable_ue_ip_paa()->set_ipv4_addr(ip4_str);
    proto->mutable_ue_ip_paa()->set_ipv6_addr(ip6_str);
    proto->mutable_ue_ip_paa()->set_ipv6_prefix_length(
        session_request->paa.ipv6_prefix_length);
    proto->mutable_ue_ip_paa()->set_vlan(session_request->paa.vlan);
  }

  proto->set_peer_ip(session_request->edns_peer_ip.addr_v4.sin_addr.s_addr);

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

  for (int i = 0; i < session_request->pco.num_protocol_or_container_id; i++) {
    auto* pco_protocol = &session_request->pco.protocol_or_container_ids[i];
    auto* pco_protocol_proto = proto->mutable_pco()->add_pco_protocol();
    if (pco_protocol->contents) {
      pco_protocol_proto->set_id(pco_protocol->id);
      pco_protocol_proto->set_length(pco_protocol->length);
      BSTRING_TO_STRING(pco_protocol->contents,
                        pco_protocol_proto->mutable_contents());
    }
  }
  for (int i = 0;
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
