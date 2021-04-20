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
#include <string>

#include "lte/protos/spgw_service.pb.h"
#include <folly/IPAddress.h>

extern "C" {
#include "spgw_service_handler.h"
#include "log.h"
}
#include "SpgwServiceImpl.h"

namespace grpc {
class ServerContext;
}  // namespace grpc

using grpc::ServerContext;
using grpc::Status;
using magma::CreateBearerRequest;
using magma::CreateBearerResult;
using magma::DeleteBearerRequest;
using magma::DeleteBearerResult;
using magma::SpgwService;

namespace magma {
using namespace lte;

SpgwServiceImpl::SpgwServiceImpl() {}
/*
 * CreateBearer is called by North Bound to create dedicated bearers
 */
Status SpgwServiceImpl::CreateBearer(
    ServerContext* context, const CreateBearerRequest* request,
    CreateBearerResult* response) {
  OAILOG_INFO(LOG_UTIL, "Received CreateBearer GRPC request\n");
  itti_gx_nw_init_actv_bearer_request_t itti_msg;
  std::string imsi = request->sid().id();
  // If north bound is sessiond itself, IMSI prefix is used;
  // in S1AP tests, IMSI prefix is not used
  // Strip off any IMSI prefix
  if (imsi.compare(0, 4, "IMSI") == 0) {
    imsi = imsi.substr(4, std::string::npos);
  }
  itti_msg.imsi_length = imsi.size();
  strcpy(itti_msg.imsi, imsi.c_str());
  itti_msg.lbi = request->link_bearer_id();

  memset(&itti_msg.ul_tft, 0, sizeof(traffic_flow_template_t));
  memset(&itti_msg.dl_tft, 0, sizeof(traffic_flow_template_t));

  // NOTE: For each policy rule a separate bearer is created
  // Future improvement:
  // (1) Rather than passing policy rules from sessiond as is
  // it would be better to have a QoS vector to bearer mapping first
  // and then issue create bearer. This will require changes to the RPC
  // request.
  // (2) Refactor this code with functions to copy from
  // policy rules to itti message fields
  bearer_qos_t* qos               = &itti_msg.eps_bearer_qos;
  traffic_flow_template_t* ul_tft = &itti_msg.ul_tft;
  traffic_flow_template_t* dl_tft = &itti_msg.dl_tft;
  for (const auto& policy_rule : request->policy_rules()) {
    // Copy the policy rule name
    std::string policy_rule_name = policy_rule.id();
    // Truncate to maximum allowed in ITTI message
    uint8_t truncated_len = std::min(
        policy_rule_name.size(), (std::size_t) POLICY_RULE_NAME_MAXLEN);
    strncpy(itti_msg.policy_rule_name, policy_rule_name.c_str(), truncated_len);
    itti_msg.policy_rule_name[truncated_len] = '\0';
    itti_msg.policy_rule_name_length         = truncated_len;
    // Copy the QoS vector specified in the policy rule
    qos->pci       = policy_rule.qos().arp().pre_capability();
    qos->pl        = policy_rule.qos().arp().priority_level();
    qos->pvi       = policy_rule.qos().arp().pre_vulnerability();
    qos->qci       = policy_rule.qos().qci();
    qos->gbr.br_ul = policy_rule.qos().gbr_ul();
    qos->gbr.br_dl = policy_rule.qos().gbr_dl();
    qos->mbr.br_ul = policy_rule.qos().max_req_bw_ul();
    qos->mbr.br_dl = policy_rule.qos().max_req_bw_dl();
    // Go through the flow list in the policy rule and
    // populate packet filters for new uplink TFT
    // A new bearer comes with new TFT rules
    ul_tft->tftoperationcode = TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
    dl_tft->tftoperationcode = TRAFFIC_FLOW_TEMPLATE_OPCODE_CREATE_NEW_TFT;
    // Currently we do not have additional parameter list passed by GRPC call
    ul_tft->ebit = TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED;
    dl_tft->ebit = TRAFFIC_FLOW_TEMPLATE_PARAMETER_LIST_IS_NOT_INCLUDED;
    int ul_count_packetfilters = 0;
    int dl_count_packetfilters = 0;
    for (const auto& flow : policy_rule.flow_list()) {
      // Skip to next flow if flow rule is for denying access;
      // TFT is used for mapping admitted flows onto bearers
      if (flow.action() == FlowDescription::DENY) {
        continue;
      }
      // There is flow rule to process, but already maxed out UL or DL PFs
      // handle this as error and cancel RPC call
      if (((ul_count_packetfilters ==
            TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX) &&
           (flow.match().direction() == FlowMatch::UPLINK)) ||
          ((dl_count_packetfilters ==
            TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX) &&
           (flow.match().direction() == FlowMatch::DOWNLINK))) {
        OAILOG_INFO(
            LOG_UTIL,
            "Received More UL or DL Flow Rules in Policy Rule (%d) than the"
            " maximum packet filters allowed (%d), rejecting the request \n",
            policy_rule.flow_list().size(),
            TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX);
        return Status::CANCELLED;
      }

      if ((flow.match().direction() == FlowMatch::UPLINK) &&
          (ul_count_packetfilters <
           TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX)) {
        ul_tft->packetfilterlist.createnewtft[ul_count_packetfilters]
            .direction = TRAFFIC_FLOW_TEMPLATE_UPLINK_ONLY;
        ul_tft->packetfilterlist.createnewtft[ul_count_packetfilters]
            .eval_precedence = policy_rule.priority();
        if (!fillUpPacketFilterContents(
                &ul_tft->packetfilterlist.createnewtft[ul_count_packetfilters]
                     .packetfiltercontents,
                &flow.match())) {
          OAILOG_ERROR(
              LOG_UTIL,
              "The uplink packet filter contents are not formatted correctly."
              "Cancelling dedicated bearer request. \n");
          return Status::CANCELLED;
        }
        ++ul_count_packetfilters;
      } else if (
          (flow.match().direction() == FlowMatch::DOWNLINK) &&
          (dl_count_packetfilters <
           TRAFFIC_FLOW_TEMPLATE_NB_PACKET_FILTERS_MAX)) {
        dl_tft->packetfilterlist.createnewtft[dl_count_packetfilters]
            .direction = TRAFFIC_FLOW_TEMPLATE_DOWNLINK_ONLY;
        dl_tft->packetfilterlist.createnewtft[dl_count_packetfilters]
            .eval_precedence = policy_rule.priority();
        if (!fillUpPacketFilterContents(
                &dl_tft->packetfilterlist.createnewtft[dl_count_packetfilters]
                     .packetfiltercontents,
                &flow.match())) {
          OAILOG_ERROR(
              LOG_UTIL,
              "The downlink packet filter contents are not formatted correctly."
              "Cancelling dedicated bearer request. \n");
          return Status::CANCELLED;
        }
        ++dl_count_packetfilters;
      }

      OAILOG_DEBUG(
          LOG_UTIL,
          " Flow Tuple (0 or empty if field does not exist)"
          " IP Protocol Number: %d"
          " Source IP address: %s"
          " Source TCP port: %d"
          " Source UDP port: %d"
          " Destination IP address: %s"
          " Destination TCP port: %d"
          " Destination UDP port: %d \n",
          flow.match().ip_proto(), flow.match().ip_src().address().c_str(),
          flow.match().tcp_src(), flow.match().udp_src(),
          flow.match().ip_dst().address().c_str(), flow.match().tcp_dst(),
          flow.match().udp_dst());
    }

    ul_tft->numberofpacketfilters = ul_count_packetfilters;
    dl_tft->numberofpacketfilters = dl_count_packetfilters;
    send_activate_bearer_request_itti(&itti_msg);
  }

  return Status::OK;
}  // namespace magma

Status SpgwServiceImpl::DeleteBearer(
    ServerContext* context, const DeleteBearerRequest* request,
    DeleteBearerResult* response) {
  OAILOG_INFO(LOG_UTIL, "Received DeleteBearer GRPC request\n");
  itti_gx_nw_init_deactv_bearer_request_t itti_msg;
  std::string imsi = request->sid().id();
  // If north bound is sessiond itself, IMSI prefix is used;
  // in S1AP tests, IMSI prefix is not used
  // Strip off any IMSI prefix
  if (imsi.compare(0, 4, "IMSI") == 0) {
    imsi = imsi.substr(4, std::string::npos);
  }
  itti_msg.imsi_length = imsi.size();
  strcpy(itti_msg.imsi, imsi.c_str());
  itti_msg.lbi           = request->link_bearer_id();
  itti_msg.no_of_bearers = 1;
  for (int i = 0; i < request->eps_bearer_ids_size() && i < BEARERS_PER_UE;
       i++) {
    itti_msg.ebi[0] = request->eps_bearer_ids(i);
    send_deactivate_bearer_request_itti(&itti_msg);
  }
  return Status::OK;
}

bool SpgwServiceImpl::fillUpPacketFilterContents(
    packet_filter_contents_t* pf_content, const FlowMatch* flow_match_rule) {
  uint16_t flags                            = 0;
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

// IPv4 address format ex.: 192.176.128.10/24
// FEG can provide an empty string which indicates
// ANY and it is equivalent to 0.0.0.0/0
// But this function is called only for non-empty ipv4 string
bool SpgwServiceImpl::fillIpv4(
    packet_filter_contents_t* pf_content, const std::string ipv4addr) {
  const auto cidrNetworkExpect = folly::IPAddress::tryCreateNetwork(ipv4addr);
  if (cidrNetworkExpect.hasError()) {
    OAILOG_ERROR(LOG_UTIL, "Invalid address string %s \n", ipv4addr.c_str());
    return false;
  }
  // Host Byte Order
  uint32_t ipv4addrHBO = cidrNetworkExpect.value().first.asV4().toLongHBO();
  for (int i = (TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE - 1); i >= 0; --i) {
    pf_content->ipv4remoteaddr[i].addr = (unsigned char) ipv4addrHBO & 0xFF;
    ipv4addrHBO                        = ipv4addrHBO >> 8;
  }

  // Get the mask length:
  // folly takes care of absence of mask_len by defaulting to 32
  // i.e., 255.255.255.255.
  int mask_len  = cidrNetworkExpect.value().second;
  uint32_t mask = UINT32_MAX;        // all ones
  mask = (mask << (32 - mask_len));  // first mask_len bits are 1s, rest 0s
  for (int i = (TRAFFIC_FLOW_TEMPLATE_IPV4_ADDR_SIZE - 1); i >= 0; --i) {
    pf_content->ipv4remoteaddr[i].mask = (unsigned char) mask & 0xFF;
    mask                               = mask >> 8;
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

bool SpgwServiceImpl::fillIpv6(
    packet_filter_contents_t* pf_content, const std::string ipv6addr) {
  struct in6_addr in6addr;
  if (inet_pton(AF_INET6, ipv6addr.c_str(), &in6addr) != 1) {
    OAILOG_ERROR(LOG_UTIL, "Invalid address string %s \n", ipv6addr.c_str());
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
