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
#include <arpa/inet.h>
#include "OpenflowController.h"
#include "PagingApplication.h"
#include "MobilityClientAPI.h"

extern "C" {
#include "log.h"
#include "sgw_paging.h"
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

void PagingApplication::event_callback(
    const ControllerEvent& ev, const OpenflowMessenger& messenger) {
  if (ev.get_type() == EVENT_PACKET_IN) {
    OAILOG_DEBUG(LOG_GTPV1U, "Handling packet-in message in paging app\n");
    const PacketInEvent& pi = static_cast<const PacketInEvent&>(ev);
    of13::PacketIn ofpi;
    ofpi.unpack(const_cast<uint8_t*>(pi.get_data()));
    handle_paging_message(
        ev.get_connection(), static_cast<uint8_t*>(ofpi.data()), messenger);
  } else if (ev.get_type() == EVENT_ADD_PAGING_RULE) {
    auto add_paging_rule_event = static_cast<const AddPagingRuleEvent&>(ev);
    add_paging_flow(add_paging_rule_event, messenger);
  } else if (ev.get_type() == EVENT_DELETE_PAGING_RULE) {
    auto delete_paging_rule_event =
        static_cast<const DeletePagingRuleEvent&>(ev);
    delete_paging_flow(delete_paging_rule_event, messenger);
  }
}

void PagingApplication::handle_paging_message(
    fluid_base::OFConnection* ofconn, uint8_t* data,
    const OpenflowMessenger& messenger) {
  // send paging request to MME
  struct ip* ip_header = (struct ip*) (data + ETH_HEADER_LENGTH);
  struct in_addr dest_ip;
  memcpy(&dest_ip, &ip_header->ip_dst, sizeof(struct in_addr));
  char* dest_ip_str = inet_ntoa(dest_ip);
  OAILOG_DEBUG(
      LOG_GTPV1U, "Initiating paging procedure for IP %s\n", dest_ip_str);
  sgw_send_paging_request(&dest_ip);

  /*
   * Clamp on this ip for configured amount of time
   * Priority is above default paging flow, but below gtp flow. This way when
   * paging succeeds, this flow will be ignored.
   * The clamping time is necessary to prevent packets from continually hitting
   * userspace, and as a retry time if paging fails
   */
  of13::FlowMod fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, MID_PRIORITY + 1);
  fm.hard_timeout(CLAMPING_TIMEOUT);
  of13::EthType type_match(IP_ETH_TYPE);
  fm.add_oxm_field(type_match);

  of13::IPv4Dst ip_match(dest_ip.s_addr);
  fm.add_oxm_field(ip_match);

  // No actions mean packet is dropped
  messenger.send_of_msg(fm, ofconn);
  return;
}

void PagingApplication::add_paging_flow(
    const AddPagingRuleEvent& ev, const OpenflowMessenger& messenger) {
  of13::FlowMod fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, MID_PRIORITY);
  // IP eth type
  of13::EthType type_match(IP_ETH_TYPE);
  fm.add_oxm_field(type_match);

  // Match on UE IP addr
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
  OAILOG_INFO(LOG_GTPV1U, "Added paging flow rule for UE IP %s\n", ip_str);
}

void PagingApplication::delete_paging_flow(
    const DeletePagingRuleEvent& ev, const OpenflowMessenger& messenger) {
  of13::FlowMod fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_DELETE, 0);

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

}  // namespace openflow
