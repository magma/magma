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

#include <netinet/ip.h>
#include <netinet/ip6.h>
#include <arpa/inet.h>
#include "lte/gateway/c/core/oai/lib/openflow/controller/OpenflowController.hpp"
#include "lte/gateway/c/core/oai/lib/openflow/controller/PagingApplication.hpp"
#include "lte/gateway/c/core/oai/lib/mobility_client/MobilityClientAPI.hpp"
#include "lte/gateway/c/core/oai/tasks/sgw/sgw_paging.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
}

using namespace fluid_msg;

namespace openflow {

uint32_t prefix2mask(int prefix) {
  if (prefix) {
    return htonl(~((1 << (32 - prefix)) - 1));
  } else {
    return htonl(0);
  }
}

void PagingApplication::event_callback(const ControllerEvent& ev,
                                       const OpenflowMessenger& messenger) {
  if (ev.get_type() == EVENT_PACKET_IN) {
    const PacketInEvent& pi = static_cast<const PacketInEvent&>(ev);
    of13::PacketIn ofpi;
    ofpi.unpack(const_cast<uint8_t*>(pi.get_data()));
    OAILOG_DEBUG(LOG_GTPV1U,
                 "Handling packet-in message in paging app: tbl: %d\n",
                 ofpi.table_id());
    if (ofpi.table_id() == SPGW_OVS_TABLE_ID) {
      handle_paging_message(ev.get_connection(),
                            static_cast<uint8_t*>(ofpi.data()), messenger);
    }
  } else if (ev.get_type() == EVENT_ADD_PAGING_RULE) {
    auto add_paging_rule_event = static_cast<const AddPagingRuleEvent&>(ev);
    // Add paging rule for ipv4 and ipv6
    add_paging_flow(add_paging_rule_event, messenger);
    add_paging_flow_ipv6(add_paging_rule_event, messenger);
  } else if (ev.get_type() == EVENT_DELETE_PAGING_RULE) {
    auto delete_paging_rule_event =
        static_cast<const DeletePagingRuleEvent&>(ev);
    // Delete paging rule for ipv4 and ipv6
    delete_paging_flow_ipv6(delete_paging_rule_event, messenger);
    delete_paging_flow(delete_paging_rule_event, messenger);
  }
}

void PagingApplication::handle_paging_message(
    fluid_base::OFConnection* ofconn, uint8_t* data,
    const OpenflowMessenger& messenger) {
  // send paging request to MME
  struct ip ip_header;
  memcpy(&ip_header, data + ETH_HEADER_LENGTH, sizeof(ip_header));

  if (ip_header.ip_v == 6) {
    handle_paging_ipv6_message(ofconn, data, messenger);
  } else {
    struct in_addr* dest_ipv4 = &ip_header.ip_dst;
    char* dest_ip_str = inet_ntoa(*dest_ipv4);

    OAILOG_DEBUG(LOG_GTPV1U, "Initiating paging procedure for IP %s\n",
                 dest_ip_str);
    sgw_send_paging_request(dest_ipv4, NULL);

    /*
     * Clamp on this ip for configured amount of time
     * Priority is above default paging flow, but below gtp flow. This way when
     * paging succeeds, this flow will be ignored.
     * The clamping time is necessary to prevent packets from continually
     * hitting userspace, and as a retry time if paging fails
     */
    of13::FlowMod fm = messenger.create_default_flow_mod(
        SPGW_OVS_TABLE_ID, of13::OFPFC_ADD, MID_PRIORITY + 1);
    fm.hard_timeout(CLAMPING_TIMEOUT);
    of13::EthType type_match(IP_ETH_TYPE);
    fm.add_oxm_field(type_match);

    of13::IPv4Dst ip_match(dest_ipv4->s_addr);
    fm.add_oxm_field(ip_match);

    // No actions mean packet is dropped
    messenger.send_of_msg(fm, ofconn);
    return;
  }
}

static void mask_ipv6_address(uint8_t* dst, const uint8_t* src,
                              const uint8_t* mask) {
  for (int i = 0; i < sizeof(struct in6_addr); i++) {
    dst[i] = src[i] & mask[i];
  }
}

void PagingApplication::handle_paging_ipv6_message(
    fluid_base::OFConnection* ofconn, uint8_t* data,
    const OpenflowMessenger& messenger) {
  if (!data) {
    OAILOG_ERROR(LOG_GTPV1U, "IPv6 header is NULL\n");
    return;
  }
  // send paging request to MME
  struct ip6_hdr ipv6_header;
  memcpy(&ipv6_header, data + ETH_HEADER_LENGTH, sizeof(ipv6_header));

  char ip6_str[INET6_ADDRSTRLEN];
  struct in6_addr* dest_ipv6 = &ipv6_header.ip6_dst;
  inet_ntop(AF_INET6, dest_ipv6, ip6_str, INET6_ADDRSTRLEN);

  OAILOG_DEBUG(LOG_GTPV1U, "Initiating paging procedure for IPv6 %s\n",
               ip6_str);

  sgw_send_paging_request(NULL, dest_ipv6);

  /*
   * Clamp on this ip for configured amount of time
   * Priority is above default paging flow, but below gtp flow. This way when
   * paging succeeds, this flow will be ignored.
   * The clamping time is necessary to prevent packets from continually hitting
   * userspace, and as a retry time if paging fails
   */
  of13::FlowMod fm = messenger.create_default_flow_mod(
      SPGW_OVS_TABLE_ID, of13::OFPFC_ADD, MID_PRIORITY + 1);
  fm.hard_timeout(CLAMPING_TIMEOUT);
  of13::EthType ip6_type(0x86DD);
  fm.add_oxm_field(ip6_type);

  static IPAddress mask("ffff:ffff:ffff:ffff::");

  // Match UE IP destination
  of13::IPv6Dst ipv6_match(IPAddress((const uint8_t*)dest_ipv6), mask);
  fm.add_oxm_field(ipv6_match);

  // No actions mean packet is dropped
  messenger.send_of_msg(fm, ofconn);
  return;
}

void PagingApplication::add_paging_flow(const AddPagingRuleEvent& ev,
                                        const OpenflowMessenger& messenger) {
  of13::FlowMod fm = messenger.create_default_flow_mod(
      SPGW_OVS_TABLE_ID, of13::OFPFC_ADD, MID_PRIORITY);
  // IP eth type
  of13::EthType type_match(IP_ETH_TYPE);
  fm.add_oxm_field(type_match);

  // Match on UE IP addr
  UeNetworkInfo ue_info_ = ev.get_ue_info();
  if (!(ue_info_.is_ue_ipv4_addr_valid())) {
    OAILOG_DEBUG(LOG_GTPV1U, "Not an IPv4 UE\n");
  } else {
    const struct in_addr& ue_ip = ev.get_ue_ip();
    of13::IPv4Dst ip_match(ue_ip.s_addr);
    fm.add_oxm_field(ip_match);

    // Output to controller
    of13::OutputAction act(of13::OFPP_CONTROLLER, of13::OFPCML_NO_BUFFER);
    of13::ApplyActions inst;
    inst.add_action(act);
    fm.add_instruction(inst);

    messenger.send_of_msg(fm, ev.get_connection());
    // Convert to string for logging
    char ip_str[INET_ADDRSTRLEN];
    inet_ntop(AF_INET, &(ue_ip.s_addr), ip_str, INET_ADDRSTRLEN);
    OAILOG_INFO(LOG_GTPV1U, "Added paging flow rule for UE IPv4 %s\n", ip_str);
  }
}

void PagingApplication::add_paging_flow_ipv6(
    const AddPagingRuleEvent& ev, const OpenflowMessenger& messenger) {
  of13::FlowMod fm = messenger.create_default_flow_mod(
      SPGW_OVS_TABLE_ID, of13::OFPFC_ADD, MID_PRIORITY);
  // IP eth type
  of13::EthType ip6_type(0x86DD);
  fm.add_oxm_field(ip6_type);

  // Match on UE IP addr, compare get_ue_ipv6 to uin6_addr_any
  UeNetworkInfo ue_info_ = ev.get_ue_info();
  if (!(ue_info_.is_ue_ipv6_addr_valid())) {
    OAILOG_ERROR(LOG_GTPV1U, "Not an IPv6 UE\n");
  } else {
    const struct in6_addr& ue_ipv6 = ev.get_ue_ipv6();
    static IPAddress mask("ffff:ffff:ffff:ffff::");
    // Match UE IP destination
    struct in6_addr ue_ip6_masked;
    mask_ipv6_address((uint8_t*)&ue_ip6_masked, (const uint8_t*)&ue_ipv6,
                      mask.getIPv6());

    of13::IPv6Dst ipv6_match(IPAddress(ue_ip6_masked), mask);
    fm.add_oxm_field(ipv6_match);

    // Output to controller
    of13::OutputAction act(of13::OFPP_CONTROLLER, of13::OFPCML_NO_BUFFER);
    of13::ApplyActions inst;
    inst.add_action(act);
    fm.add_instruction(inst);

    messenger.send_of_msg(fm, ev.get_connection());
    // Convert to string for logging
    char ip_str[INET6_ADDRSTRLEN] = {};
    inet_ntop(AF_INET6, &(ue_ipv6), ip_str, INET6_ADDRSTRLEN);
    OAILOG_INFO(LOG_GTPV1U, "Added paging flow rule for UE IPv6 %s\n", ip_str);
  }
}

void PagingApplication::delete_paging_flow(const DeletePagingRuleEvent& ev,
                                           const OpenflowMessenger& messenger) {
  of13::FlowMod fm = messenger.create_default_flow_mod(SPGW_OVS_TABLE_ID,
                                                       of13::OFPFC_DELETE, 0);

  // IP eth type
  of13::EthType type_match(IP_ETH_TYPE);
  fm.add_oxm_field(type_match);

  // match all ports and groups
  fm.out_port(of13::OFPP_ANY);
  fm.out_group(of13::OFPG_ANY);

  // Match on UE IP addr
  const struct in_addr& ue_ip = ev.get_ue_ip();
  of13::IPv4Dst ip_match(ue_ip.s_addr);
  fm.add_oxm_field(ip_match);

  // Output to controller
  // (This has actually no effect on deletion, but included
  // for symmetry purposes wrt add_paging_flow)
  of13::OutputAction act(of13::OFPP_CONTROLLER, of13::OFPCML_NO_BUFFER);
  of13::ApplyActions inst;
  inst.add_action(act);
  fm.add_instruction(inst);

  messenger.send_of_msg(fm, ev.get_connection());
  // Convert to string for logging
  char* ip_str = inet_ntoa(ue_ip);
  OAILOG_INFO(LOG_GTPV1U, "Deleted paging flow rule for UE IP %s\n", ip_str);
}

void PagingApplication::delete_paging_flow_ipv6(
    const DeletePagingRuleEvent& ev, const OpenflowMessenger& messenger) {
  of13::FlowMod fm = messenger.create_default_flow_mod(SPGW_OVS_TABLE_ID,
                                                       of13::OFPFC_DELETE, 0);

  // IP eth type
  of13::EthType ip6_type(0x86DD);
  fm.add_oxm_field(ip6_type);

  // match all ports and groups
  fm.out_port(of13::OFPP_ANY);
  fm.out_group(of13::OFPG_ANY);

  // Match on UE IP addr
  UeNetworkInfo ue_info_ = ev.get_ue_info();
  if (!(ue_info_.is_ue_ipv6_addr_valid())) {
    OAILOG_ERROR(LOG_GTPV1U, "Not an IPv6 UE\n");
  } else {
    const struct in6_addr& ue_ipv6 = ev.get_ue_ipv6();
    static IPAddress mask("ffff:ffff:ffff:ffff::");
    struct in6_addr ue_ip6_masked;
    mask_ipv6_address((uint8_t*)&ue_ip6_masked, (const uint8_t*)&ue_ipv6,
                      mask.getIPv6());

    of13::IPv6Dst ipv6_match(IPAddress(ue_ip6_masked), mask);
    fm.add_oxm_field(ipv6_match);

    // Output to controller
    // (This has actually no effect on deletion, but included
    // for symmetry purposes wrt add_paging_flow  // )
    of13::OutputAction act(of13::OFPP_CONTROLLER, of13::OFPCML_NO_BUFFER);
    of13::ApplyActions inst;
    inst.add_action(act);
    fm.add_instruction(inst);

    messenger.send_of_msg(fm, ev.get_connection());
    // Convert to string for logging
    char ip_str[INET6_ADDRSTRLEN];
    inet_ntop(AF_INET6, &(ue_ipv6), ip_str, INET6_ADDRSTRLEN);
    OAILOG_INFO(LOG_GTPV1U, "Deleted paging flow rule for UE IPv6 %s\n",
                ip_str);
  }
}

}  // namespace openflow
