/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <string>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/include/amf_service_handler.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/grpc_service/AmfServiceImpl.hpp"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"

namespace grpc {

class ServerContext;
}  // namespace grpc

using grpc::ServerContext;
using grpc::Status;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContextAccess;
using magma::lte::SmContextVoid;
using magma::lte::SmfPduSessionSmContext;

namespace magma {
using namespace lte;
#define TEID_SIZE 4
#define UPF_IPV4_ADDR_SIZE 4

AmfServiceImpl::AmfServiceImpl() {}

// Remove the Leading IMSI string if present
inline void get_subscriber_id(const std::string& subscriber_id, char* imsi) {
  // No parameter check as these should always be filled up
  uint8_t imsi_len = 0;

  // Check if the subscriber information received contains IMSI
  if (subscriber_id.compare(0, 4, "IMSI") == 0) {
    // If yes then remove the same
    imsi_len = strlen("IMSI");
  }

  strcpy(imsi, subscriber_id.c_str() + imsi_len);
}

Status AmfServiceImpl::SetAmfNotification(ServerContext* context,
                                          const SetSmNotificationContext* notif,
                                          SmContextVoid* response) {
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSmNotificationContext request\n");
  // ToDo processing ITTI,ZMQ

  itti_n11_received_notification_t itti_msg;
  memset(&itti_msg, 0, sizeof(itti_n11_received_notification_t));

  SetAmfNotification_itti(notif, &itti_msg);

  send_n11_notification_received_itti(&itti_msg);
  return Status::OK;
}

Status AmfServiceImpl::SetAmfNotification_itti(
    const SetSmNotificationContext* notif,
    itti_n11_received_notification_t* itti_msg) {
  auto& notify_common = notif->common_context();
  auto& req_m5g = notif->rat_specific_notification();

  // CommonSessionContext
  get_subscriber_id(notify_common.sid().id(), itti_msg->imsi);

  itti_msg->sm_session_fsm_state =
      (SMSessionFSMState_response)notify_common.sm_session_state();
  itti_msg->sm_session_version = notify_common.sm_session_version();

  // RatSpecificContextAccessS
  itti_msg->pdu_session_id = req_m5g.pdu_session_id();
  itti_msg->request_type = (RequestType_received)req_m5g.request_type();
  itti_msg->pdu_session_type = (pdu_session_type_t)req_m5g.pdu_session_type();
  itti_msg->m5g_sm_capability.reflective_qos =
      req_m5g.m5g_sm_capability().reflective_qos();
  itti_msg->m5g_sm_capability.multi_homed_ipv6_pdu_session =
      req_m5g.m5g_sm_capability().multi_homed_ipv6_pdu_session();
  itti_msg->m5gsm_cause = (m5g_sm_cause_t)req_m5g.m5gsm_cause();

  // pdu_change
  itti_msg->notify_ue_evnt = (notify_ue_event)req_m5g.notify_ue_event();
  return Status::OK;
}
// Set message from SessionD received
Status AmfServiceImpl::SetSmfSessionContext(
    ServerContext* context, const SetSMSessionContextAccess* request,
    SmContextVoid* response) {
  itti_n11_create_pdu_session_response_t itti_msg;
  memset(&itti_msg, 0, sizeof(itti_n11_create_pdu_session_response_t));

  if (SetSmfSessionContext_itti(request, &itti_msg) == false) {
    return Status::CANCELLED;
  }

  send_n11_create_pdu_session_resp_itti(&itti_msg);
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSMSessionContextAccess request \n");
  return Status::OK;
}

bool AmfServiceImpl::SetSmfSessionContext_itti(
    const SetSMSessionContextAccess* request,
    itti_n11_create_pdu_session_response_t* itti_msg_p) {
  uint32_t i = 0;
  traffic_flow_template_t* ul_tft = NULL;
  int ul_count_packetfilters = 0;
  OAILOG_INFO(LOG_UTIL,
              "Received GRPC SetSmfSessionContext request from SMF\n");

  auto& req_common = request->common_context();
  auto& req_m5g = request->rat_specific_context().m5g_session_context_rsp();

  // CommonSessionContext
  get_subscriber_id(req_common.sid().id(), itti_msg_p->imsi);

  itti_msg_p->sm_session_fsm_state =
      (sm_session_fsm_state_t)req_common.sm_session_state();
  itti_msg_p->sm_session_version = req_common.sm_session_version();

  // RatSpecificContextAccess
  itti_msg_p->pdu_session_id = req_m5g.pdu_session_id();
  itti_msg_p->pdu_session_type = (pdu_session_type_t)req_m5g.pdu_session_type();
  itti_msg_p->selected_ssc_mode = (ssc_mode_t)req_m5g.selected_ssc_mode();
  itti_msg_p->m5gsm_cause = (m5g_sm_cause_t)req_m5g.m5gsm_cause();

  if (!(req_m5g.qos_policy_size()) && req_m5g.has_subscribed_qos()) {
    itti_msg_p->session_ambr.uplink_unit_type =
        req_m5g.subscribed_qos().br_unit();
    itti_msg_p->session_ambr.uplink_units =
        req_m5g.subscribed_qos().apn_ambr_ul();

    itti_msg_p->session_ambr.downlink_unit_type =
        req_m5g.subscribed_qos().br_unit();
    itti_msg_p->session_ambr.downlink_units =
        req_m5g.subscribed_qos().apn_ambr_dl();

    // authorized qos profile
    itti_msg_p->qos_flow_list.item[i].qos_flow_req_item.qos_flow_identifier =
        req_m5g.subscribed_qos().qos_class_id();

    // default flow descriptors
    if (req_m5g.subscribed_qos().qos_class_id()) {
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier =
          req_m5g.subscribed_qos().qos_class_id();

      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_descriptor.fiveQi =
          req_m5g.subscribed_qos().qos_class_id();
    }

    itti_msg_p->qos_flow_list.item[i]
        .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
        .non_dynamic_5QI_desc.fiveQI =
        req_m5g.subscribed_qos().qos_class_id();  // enum
    itti_msg_p->qos_flow_list.item[i]
        .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
        .priority_level = req_m5g.subscribed_qos().priority_level();  // uint32
    itti_msg_p->qos_flow_list.item[i]
        .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
        .pre_emption_cap = (pre_emption_capability)req_m5g.subscribed_qos()
                               .preemption_capability();  // enum
    itti_msg_p->qos_flow_list.item[i]
        .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
        .pre_emption_vul = (pre_emption_vulnerability)req_m5g.subscribed_qos()
                               .preemption_vulnerability();  // enum
    i++;
  }

  for (; i < req_m5g.qos_policy_size(); i++) {
    ul_tft = &itti_msg_p->qos_flow_list.item[i].qos_flow_req_item.ul_tft;
    memset(ul_tft, 0, sizeof(traffic_flow_template_t));
    auto& qos_rule = req_m5g.qos_policy(i);
    // Session ambr is policy ambr if policy attached
    itti_msg_p->session_ambr.uplink_units =
        qos_rule.qos().qos().max_req_bw_ul();

    itti_msg_p->session_ambr.downlink_units =
        qos_rule.qos().qos().max_req_bw_dl();
    itti_msg_p->qos_flow_list.item[i].qos_flow_req_item.qos_flow_identifier =
        qos_rule.qos().qos().qci();

    itti_msg_p->qos_flow_list.item[i]
        .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
        .non_dynamic_5QI_desc.fiveQI = qos_rule.qos().qos().qci();  // enum
    if (qos_rule.qos().qos().arp().priority_level()) {
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
          .priority_level =
          qos_rule.qos().qos().arp().priority_level();  // uint32
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
          .pre_emption_cap = (pre_emption_capability)qos_rule.qos()
                                 .qos()
                                 .arp()
                                 .pre_capability();  // enum
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
          .pre_emption_vul = (pre_emption_vulnerability)qos_rule.qos()
                                 .qos()
                                 .arp()
                                 .pre_vulnerability();  // enum
    } else {
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
          .priority_level =
          req_m5g.subscribed_qos().priority_level();  // uint32
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
          .pre_emption_cap = (pre_emption_capability)req_m5g.subscribed_qos()
                                 .preemption_capability();  // enum
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
          .pre_emption_vul = (pre_emption_vulnerability)req_m5g.subscribed_qos()
                                 .preemption_vulnerability();  // enum
    }
    itti_msg_p->qos_flow_list.item[i].qos_flow_req_item.qos_flow_action =
        (qos_flow_action_t)qos_rule.policy_action();  // enum
    itti_msg_p->qos_flow_list.item[i].qos_flow_req_item.qos_flow_version =
        qos_rule.version();  // uint32
    strncpy(reinterpret_cast<char*>(
                itti_msg_p->qos_flow_list.item[i].qos_flow_req_item.rule_id),
            qos_rule.qos().id().c_str(), strlen(qos_rule.qos().id().c_str()));

    // flow descriptor
    if (qos_rule.qos().has_qos()) {
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier =
          qos_rule.qos().qos().qci();
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_descriptor.fiveQi =
          qos_rule.qos().qos().qci();
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_descriptor.mbr_dl =
          qos_rule.qos().qos().max_req_bw_dl();
      itti_msg_p->qos_flow_list.item[i]
          .qos_flow_req_item.qos_flow_descriptor.mbr_ul =
          qos_rule.qos().qos().max_req_bw_ul();
      if (qos_rule.qos().qos().gbr_dl()) {
        itti_msg_p->qos_flow_list.item[i]
            .qos_flow_req_item.qos_flow_descriptor.gbr_dl =
            qos_rule.qos().qos().gbr_dl();
        itti_msg_p->qos_flow_list.item[i]
            .qos_flow_req_item.qos_flow_descriptor.gbr_ul =
            qos_rule.qos().qos().gbr_ul();
      } else {
        itti_msg_p->qos_flow_list.item[i]
            .qos_flow_req_item.qos_flow_descriptor.gbr_dl =
            qos_rule.qos().qos().max_req_bw_dl();
        itti_msg_p->qos_flow_list.item[i]
            .qos_flow_req_item.qos_flow_descriptor.gbr_ul =
            qos_rule.qos().qos().max_req_bw_ul();
      }
    }
    for (const auto& flow : qos_rule.qos().flow_list()) {
      if (flow.action() == FlowDescription::DENY) {
        continue;
      }

      if ((flow.match().direction() == FlowMatch::UPLINK) &&
          (ul_count_packetfilters <
           TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX)) {
        if (qos_rule.policy_action() == QosPolicy::ADD) {
          ul_tft->tftoperationcode =
              TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
          ul_tft->packetfilterlist.createnewtft[ul_count_packetfilters]
              .direction = TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
          ul_tft->packetfilterlist.createnewtft[ul_count_packetfilters]
              .eval_precedence = qos_rule.qos().priority();
          if (!fillUpPacketFilterContents(
                  &ul_tft->packetfilterlist.createnewtft[ul_count_packetfilters]
                       .packetfiltercontents,
                  &flow.match())) {
            OAILOG_ERROR(
                LOG_UTIL,
                "The uplink packet filter contents are not formatted correctly."
                "Canceling qos flow creation request. \n");
            return false;
          }
          ++ul_count_packetfilters;
          ul_tft->numberofpacketfilters++;
        }
      } else if (qos_rule.policy_action() == QosPolicy::DEL) {
        ul_tft->tftoperationcode =
            TRAFFIC_FLOW_TEMPLATE_OPCODE_DELETE_EXISTING_TFT;
      }
    }
  }
  itti_msg_p->qos_flow_list.maxNumOfQosFlows = i;

  // get the 4 byte UPF TEID and UPF IP message
  uint32_t nteid = req_m5g.upf_endpoint().teid();
  itti_msg_p->upf_endpoint.teid[0] = (nteid >> 24) & 0xFF;
  itti_msg_p->upf_endpoint.teid[1] = (nteid >> 16) & 0xFF;
  itti_msg_p->upf_endpoint.teid[2] = (nteid >> 8) & 0xFF;
  itti_msg_p->upf_endpoint.teid[3] = nteid & 0xFF;

  if (req_m5g.upf_endpoint().end_ipv4_addr().size() > 0) {
    inet_pton(AF_INET, req_m5g.upf_endpoint().end_ipv4_addr().c_str(),
              itti_msg_p->upf_endpoint.end_ipv4_addr);
  }

  itti_msg_p->procedure_trans_identity = req_m5g.procedure_trans_identity();
  itti_msg_p->always_on_pdu_session_indication =
      req_m5g.always_on_pdu_session_indication();
  itti_msg_p->allowed_ssc_mode = (ssc_mode_t)req_m5g.allowed_ssc_mode();
  itti_msg_p->m5gsm_congetion_re_attempt_indicator =
      req_m5g.m5g_sm_congestion_reattempt_indicator();

  // PDU IP address coming from SMF in human-readable format has to be packed
  // into 4 raw bytes in hex for NAS5G layer

  if (req_common.ue_ipv4().size() > 0) {
    inet_pton(AF_INET, req_common.ue_ipv4().c_str(),
              &(itti_msg_p->pdu_address.ipv4_address));
    itti_msg_p->pdu_address.pdn_type = IPv4;
  }

  if (req_common.ue_ipv6().size() > 0) {
    inet_pton(AF_INET6, req_common.ue_ipv6().c_str(),
              &(itti_msg_p->pdu_address.ipv6_address));

    if (req_common.ue_ipv4().size() == 0) {
      itti_msg_p->pdu_address.pdn_type = IPv6;
    } else {
      itti_msg_p->pdu_address.pdn_type = IPv4_AND_v6;
    }

    itti_msg_p->pdu_address.ipv6_prefix_length = IPV6_PREFIX_LEN;
  }

  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSMSessionContextAccess request \n");
  return true;
}

bool AmfServiceImpl::fillUpPacketFilterContents(
    packet_filter_contents_t* pf_content, const FlowMatch* flow_match_rule) {
  uint16_t flags = 0;
  pf_content->protocolidentifier_nextheader = flow_match_rule->ip_proto();
  if (pf_content->protocolidentifier_nextheader) {
    flags |= TRAFFIC_FLOW_TEMPLATE_PROTOCOL_NEXT_HEADER_FLAG;
  }
  // If flow match rule is for UL, remote server is TCP destination
  // Else, remote server is TCP source
  // GRPC interface does not support a third option (e.g., bidirectional)
  if (flow_match_rule->direction() == FlowMatch::UPLINK) {
    if (!flow_match_rule->ip_dst().address().empty()) {
      if (flow_match_rule->ip_dst().version() ==
          flow_match_rule->ip_dst().IPV4) {
        flags |= TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
        if (!fillIpv4(pf_content, flow_match_rule->ip_dst().address())) {
          return false;
        }
      }
      if (flow_match_rule->ip_dst().version() ==
          flow_match_rule->ip_dst().IPV6) {
        flags |= TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG;
        if (!fillIpv6(pf_content, flow_match_rule->ip_dst().address())) {
          return false;
        }
      }
    }
    if (flow_match_rule->tcp_src() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
      pf_content->singlelocalport = flow_match_rule->tcp_src();
    } else if (flow_match_rule->udp_src() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
      pf_content->singlelocalport = flow_match_rule->udp_src();
    }
    if (flow_match_rule->tcp_dst() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
      pf_content->singleremoteport = flow_match_rule->tcp_dst();
    } else if (flow_match_rule->udp_dst() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
      pf_content->singleremoteport = flow_match_rule->udp_dst();
    }
  } else if (flow_match_rule->direction() == FlowMatch::DOWNLINK) {
    if (!flow_match_rule->ip_src().address().empty()) {
      if (flow_match_rule->ip_src().version() ==
          flow_match_rule->ip_src().IPV4) {
        flags |= TRAFFIC_FLOW_TEMPLATE_IPV4_REMOTE_ADDR_FLAG;
        if (!fillIpv4(pf_content, flow_match_rule->ip_src().address())) {
          return false;
        }
      }
      if (flow_match_rule->ip_src().version() ==
          flow_match_rule->ip_src().IPV6) {
        flags |= TRAFFIC_FLOW_TEMPLATE_IPV6_REMOTE_ADDR_FLAG;
        if (!fillIpv6(pf_content, flow_match_rule->ip_src().address())) {
          return false;
        }
      }
    }
    if (flow_match_rule->tcp_dst() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
      pf_content->singlelocalport = flow_match_rule->tcp_dst();
    } else if (flow_match_rule->udp_dst() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_LOCAL_PORT_FLAG;
      pf_content->singlelocalport = flow_match_rule->udp_dst();
    }
    if (flow_match_rule->tcp_src() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
      pf_content->singleremoteport = flow_match_rule->tcp_src();
    } else if (flow_match_rule->udp_src() != 0) {
      flags |= TRAFFIC_FLOW_TEMPLATE_SINGLE_REMOTE_PORT_FLAG;
      pf_content->singleremoteport = flow_match_rule->udp_src();
    }
  }

  pf_content->flags = flags;
  return true;
}

// Extract and validate IP address and subnet mask
// IPv4 network format ex.: 192.176.128.10/24
ipv4_networks_t AmfServiceImpl::parseIpv4Network(
    const std::string& ipv4network_str) {
  ipv4_networks_t result;
  const int slash_pos = ipv4network_str.find("/");
  std::string ipv4addr = (slash_pos != std::string::npos)
                             ? ipv4network_str.substr(0, slash_pos)
                             : ipv4network_str;
  in_addr addr;
  if (inet_pton(AF_INET, ipv4addr.c_str(), &addr) != 1) {
    OAILOG_ERROR(LOG_UTIL, "Invalid address string %s \n",
                 ipv4network_str.c_str());
    result.success = false;
    return result;
  }
  // Host Byte Order
  result.addr_hbo = ntohl(addr.s_addr);
  constexpr char default_mask_len_str[] = "32";
  std::string mask_len_str = (slash_pos != std::string::npos)
                                 ? ipv4network_str.substr(slash_pos + 1)
                                 : default_mask_len_str;
  int mask_len;
  try {
    mask_len = std::stoi(mask_len_str);
  } catch (...) {
    OAILOG_ERROR(LOG_UTIL, "Invalid address string %s \n",
                 ipv4network_str.c_str());
    result.success = false;
    return result;
  }
  if (mask_len > 32 || mask_len < 0) {
    OAILOG_ERROR(LOG_UTIL, "Invalid address string %s \n",
                 ipv4network_str.c_str());
    result.success = false;
    return result;
  }
  result.mask_len = mask_len;
  result.success = true;
  return result;
}

// IPv4 address format ex.: 192.176.128.10/24
// FEG can provide an empty string which indicates
// ANY and it is equivalent to 0.0.0.0/0
// But this function is called only for non-empty ipv4 string

bool AmfServiceImpl::fillIpv4(packet_filter_contents_t* pf_content,
                              const std::string& ipv4network_str) {
  ipv4_networks_t ipv4network = parseIpv4Network(ipv4network_str);
  if (!ipv4network.success) {
    return false;
  }

  uint32_t mask = UINT32_MAX;  // all ones
  mask =
      (mask << (32 -
                ipv4network.mask_len));  // first mask_len bits are 1s, rest 0s
  uint32_t ipv4addrHBO = ipv4network.addr_hbo;

  for (int i = (TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE - 1); i >= 0; --i) {
    pf_content->ipv4remoteaddr[i].mask = (unsigned char)mask & 0xFF;
    pf_content->ipv4remoteaddr[i].addr =
        (unsigned char)ipv4addrHBO & pf_content->ipv4remoteaddr[i].mask;
    ipv4addrHBO = ipv4addrHBO >> 8;
    mask = mask >> 8;
  }

  OAILOG_DEBUG(
      LOG_UTIL,
      "Network Address: %d.%d.%d.%d "
      "Network Mask: %d.%d.%d.%d \n",
      pf_content->ipv4remoteaddr[0].addr, pf_content->ipv4remoteaddr[1].addr,
      pf_content->ipv4remoteaddr[2].addr, pf_content->ipv4remoteaddr[3].addr,
      pf_content->ipv4remoteaddr[0].mask, pf_content->ipv4remoteaddr[1].mask,
      pf_content->ipv4remoteaddr[2].mask, pf_content->ipv4remoteaddr[3].mask);
  return true;
}

bool AmfServiceImpl::fillIpv6(packet_filter_contents_t* pf_content,
                              const std::string ipv6network_str) {
  struct in6_addr in6addr;
  if (inet_pton(AF_INET6, ipv6network_str.c_str(), &in6addr) != 1) {
    OAILOG_ERROR(LOG_UTIL, "Invalid address string %s \n",
                 ipv6network_str.c_str());
    return false;
  }
  for (int i = 0; i < TRAFFIC_FLOW_TEMPLATE_IPV6_ADDR_SIZE; i++) {
    pf_content->ipv6remoteaddr[i].addr = in6addr.s6_addr[i];
  }

  OAILOG_DEBUG(
      LOG_UTIL,
      "Network Address: %x:%x:%x:%x:%x:%x:%x:%x:%x:%x:%x:%x:%x:%x:%x:%x\n",
      pf_content->ipv6remoteaddr[0].addr, pf_content->ipv6remoteaddr[1].addr,
      pf_content->ipv6remoteaddr[2].addr, pf_content->ipv6remoteaddr[3].addr,
      pf_content->ipv6remoteaddr[4].addr, pf_content->ipv6remoteaddr[5].addr,
      pf_content->ipv6remoteaddr[6].addr, pf_content->ipv6remoteaddr[7].addr,
      pf_content->ipv6remoteaddr[8].addr, pf_content->ipv6remoteaddr[9].addr,
      pf_content->ipv6remoteaddr[10].addr, pf_content->ipv6remoteaddr[11].addr,
      pf_content->ipv6remoteaddr[12].addr, pf_content->ipv6remoteaddr[13].addr,
      pf_content->ipv6remoteaddr[14].addr, pf_content->ipv6remoteaddr[15].addr);
  return true;
}
}  // namespace magma
