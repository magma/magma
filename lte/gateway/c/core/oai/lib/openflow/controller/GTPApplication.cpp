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
#include <string>

#include "GTPApplication.h"
#include "IMSIEncoder.h"
#include "gtpv1u.h"

extern "C" {
#include "log.h"
#include "bstrlib.h"
}

using namespace fluid_msg;

namespace openflow {

const std::string GTPApplication::GTP_PORT_MAC = "02:00:00:00:00:01";
const std::uint16_t OFPVID_PRESENT             = 0x1000;

GTPApplication::GTPApplication(
    const std::string& uplink_mac, uint32_t gtp_port_num, uint32_t mtr_port_num,
    uint32_t internal_sampling_port_num, uint32_t internal_sampling_fwd_tbl_num,
    uint32_t uplink_port_num)
    : uplink_mac_(uplink_mac),
      gtp0_port_num_(gtp_port_num),
      mtr_port_num_(mtr_port_num),
      internal_sampling_port_num_(internal_sampling_port_num),
      internal_sampling_fwd_tbl_num_(internal_sampling_fwd_tbl_num),
      uplink_port_num_(uplink_port_num) {}

void GTPApplication::event_callback(
    const ControllerEvent& ev, const OpenflowMessenger& messenger) {
  if (ev.get_type() == EVENT_ADD_GTP_TUNNEL) {
    auto add_tunnel_event = static_cast<const AddGTPTunnelEvent&>(ev);
    add_uplink_tunnel_flow(add_tunnel_event, messenger);
    add_downlink_tunnel_flow(
        add_tunnel_event, messenger, uplink_port_num_, false, false);
    add_downlink_tunnel_flow(
        add_tunnel_event, messenger, mtr_port_num_, false, false);
    add_downlink_arp_flow(add_tunnel_event, messenger, uplink_port_num_);
    add_downlink_arp_flow(add_tunnel_event, messenger, mtr_port_num_);
  } else if (ev.get_type() == EVENT_DELETE_GTP_TUNNEL) {
    auto del_tunnel_event = static_cast<const DeleteGTPTunnelEvent&>(ev);
    delete_uplink_tunnel_flow(del_tunnel_event, messenger);
    delete_downlink_tunnel_flow(del_tunnel_event, messenger, uplink_port_num_);
    delete_downlink_tunnel_flow(del_tunnel_event, messenger, mtr_port_num_);
    delete_downlink_arp_flow(del_tunnel_event, messenger, uplink_port_num_);
    delete_downlink_arp_flow(del_tunnel_event, messenger, mtr_port_num_);
  } else if (ev.get_type() == EVENT_ADD_GTP_S8_TUNNEL) {
    auto add_tunnel_event = static_cast<const AddGTPTunnelEvent&>(ev);
    auto imsi = IMSIEncoder::compact_imsi(add_tunnel_event.get_imsi());

    OAILOG_DEBUG_UE(
        LOG_GTPV1U, imsi,
        "s8: add: TEID: s1-in %u s1-out %u s8-in %u s8-out %u\n",
        add_tunnel_event.get_in_tei(), add_tunnel_event.get_out_tei(),
        add_tunnel_event.get_pgw_in_tei(), add_tunnel_event.get_pgw_out_tei());

    add_uplink_s8_tunnel_flow(add_tunnel_event, messenger);
    int pgw_port = add_tunnel_event.get_pgw_gtp_portno();
    if (pgw_port == 0) {
      pgw_port = GTPApplication::gtp0_port_num_;
    }
    add_downlink_tunnel_flow(add_tunnel_event, messenger, pgw_port, true, true);
    add_downlink_tunnel_flow(
        add_tunnel_event, messenger, mtr_port_num_, true, true);
    add_downlink_arp_flow(add_tunnel_event, messenger, mtr_port_num_);
  } else if (ev.get_type() == EVENT_DELETE_GTP_S8_TUNNEL) {
    auto del_tunnel_event = static_cast<const DeleteGTPTunnelEvent&>(ev);
    // Same delete can be used, since uplink flow match is same for S8 tunnel.
    delete_uplink_tunnel_flow(del_tunnel_event, messenger);
    int pgw_port = del_tunnel_event.get_pgw_gtp_portno();
    if (pgw_port == 0) {
      pgw_port = GTPApplication::gtp0_port_num_;
    }
    delete_downlink_tunnel_flow(del_tunnel_event, messenger, pgw_port);

    delete_downlink_tunnel_flow(del_tunnel_event, messenger, mtr_port_num_);
    delete_downlink_arp_flow(del_tunnel_event, messenger, mtr_port_num_);
  } else if (ev.get_type() == EVENT_DISCARD_DATA_ON_GTP_TUNNEL) {
    auto discard_tunnel_flow =
        static_cast<const HandleDataOnGTPTunnelEvent&>(ev);
    discard_uplink_tunnel_flow(discard_tunnel_flow, messenger);
    discard_downlink_tunnel_flow(
        discard_tunnel_flow, messenger, uplink_port_num_);
    discard_downlink_tunnel_flow(discard_tunnel_flow, messenger, mtr_port_num_);
  } else if (ev.get_type() == EVENT_FORWARD_DATA_ON_GTP_TUNNEL) {
    auto forward_tunnel_flow =
        static_cast<const HandleDataOnGTPTunnelEvent&>(ev);
    forward_uplink_tunnel_flow(forward_tunnel_flow, messenger);
    forward_downlink_tunnel_flow(
        forward_tunnel_flow, messenger, uplink_port_num_);
    forward_downlink_tunnel_flow(forward_tunnel_flow, messenger, mtr_port_num_);
  } else if (ev.get_type() == EVENT_SWITCH_UP) {
    install_internal_pkt_fwd_flow(
        ev.get_connection(), messenger, internal_sampling_port_num_,
        internal_sampling_fwd_tbl_num_);
  }
}

void GTPApplication::install_internal_pkt_fwd_flow(
    fluid_base::OFConnection* ofconn, const OpenflowMessenger& messenger,
    uint32_t port, uint32_t next_table) {
  of13::FlowMod fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, DEFAULT_PRIORITY);

  // Set match on the internal pkt sampling port
  of13::InPort port_match(port);
  fm.add_oxm_field(port_match);

  // Output to next table
  of13::GoToTable inst(next_table);
  fm.add_instruction(inst);
  messenger.send_of_msg(fm, ofconn);
  OAILOG_DEBUG(LOG_GTPV1U, "Session tracker forward flow added\n");
}

/*
 * Helper method to add matching for adding/deleting the uplink flow
 */
void GTPApplication::add_tunnel_match(
    of13::FlowMod& uplink_fm, uint32_t gtp_port, uint32_t i_tei) {
  if (gtp_port == 0) {
    gtp_port = GTPApplication::gtp0_port_num_;
  }
  // Match on tunnel id and gtp in port
  of13::InPort gtp_port_match(gtp_port);
  uplink_fm.add_oxm_field(gtp_port_match);
  of13::TUNNELId in_tunnel_id(i_tei);
  uplink_fm.add_oxm_field(in_tunnel_id);
}

/*
 * Helper method to add imsi as metadata to the packet
 */
void add_imsi_metadata(of13::ApplyActions& apply_actions, uint64_t imsi) {
  auto metadata_field = new of13::Metadata(imsi);
  of13::SetFieldAction set_metadata(metadata_field);
  apply_actions.add_action(set_metadata);
}

void GTPApplication::add_uplink_tunnel_flow(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger) {
  auto imsi = IMSIEncoder::compact_imsi(ev.get_imsi());
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod uplink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, flow_priority);
  add_tunnel_match(uplink_fm, ev.get_enb_gtp_portno(), ev.get_in_tei());

  // Set eth src and dst
  of13::ApplyActions apply_ul_inst;
  EthAddress gtp_port(GTP_PORT_MAC);
  // libfluid handles memory freeing of fields
  of13::SetFieldAction set_eth_src(new of13::EthSrc(gtp_port));
  apply_ul_inst.add_action(set_eth_src);

  EthAddress uplink_port(uplink_mac_);
  of13::SetFieldAction set_eth_dst(new of13::EthDst(uplink_port));
  apply_ul_inst.add_action(set_eth_dst);

  of13::SetFieldAction set_tunnel_id(
      new of13::NXMRegX(TUNNEL_ID_REG, ev.get_in_tei()));
  apply_ul_inst.add_action(set_tunnel_id);

  int vlan_id = ev.get_ue_info().get_vlan();
  if (vlan_id > 0) {
    of13::PushVLANAction push_vlan(0x8100);
    apply_ul_inst.add_action(push_vlan);

    uint16_t vid = OFPVID_PRESENT | vlan_id;
    of13::SetFieldAction set_vlan(new of13::VLANVid(vid));
    apply_ul_inst.add_action(set_vlan);
  }
  // add imsi to packet metadata to pass to other tables
  add_imsi_metadata(apply_ul_inst, imsi);

  uplink_fm.add_instruction(apply_ul_inst);

  // Output to inout table
  of13::GoToTable goto_inst(NEXT_TABLE);
  uplink_fm.add_instruction(goto_inst);

  // Finally, send flow mod
  messenger.send_of_msg(uplink_fm, ev.get_connection());

  OAILOG_DEBUG_UE(LOG_GTPV1U, imsi, "Uplink flow added\n");
}

void GTPApplication::add_uplink_s8_tunnel_flow(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod uplink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, flow_priority);
  add_tunnel_match(uplink_fm, ev.get_enb_gtp_portno(), ev.get_in_tei());

  add_tunnel_flow_action(
      ev.get_pgw_out_tei(), ev.get_in_tei(), ev.get_imsi(), ev.get_pgw_ip(),
      ev.get_pgw_gtp_portno(), ev.get_connection(), messenger, uplink_fm,
      "S8 Uplink", true);
}

void GTPApplication::delete_uplink_tunnel_flow(
    const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger) {
  of13::FlowMod uplink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_DELETE, 0);
  // match all ports and groups
  uplink_fm.out_port(of13::OFPP_ANY);
  uplink_fm.out_group(of13::OFPG_ANY);

  add_tunnel_match(uplink_fm, ev.get_enb_gtp_portno(), ev.get_in_tei());

  messenger.send_of_msg(uplink_fm, ev.get_connection());
}

/*
 * Helper method to add matching for adding/deleting the downlink flow
 */
static void add_downlink_arp_match(
    of13::FlowMod& downlink_fm, const struct in_addr& ue_ip, uint32_t port)

{
  // Set match on uplink port and IP eth type
  of13::InPort uplink_port_match(port);
  downlink_fm.add_oxm_field(uplink_port_match);
  of13::EthType ip_type(0x0806);
  downlink_fm.add_oxm_field(ip_type);

  // Match UE IP destination
  of13::ARPTPA arptpa(ue_ip.s_addr);
  downlink_fm.add_oxm_field(arptpa);
}

static void add_downlink_match(
    of13::FlowMod& downlink_fm, const struct in_addr& ue_ip, uint32_t port) {
  // Set match on uplink port and IP eth type
  of13::InPort uplink_port_match(port);
  downlink_fm.add_oxm_field(uplink_port_match);
  of13::EthType ip_type(0x0800);
  downlink_fm.add_oxm_field(ip_type);

  // Match UE IP destination
  of13::IPv4Dst ip_match(ue_ip.s_addr);
  downlink_fm.add_oxm_field(ip_match);
}

static void mask_ipv6_address(
    uint8_t* dst, const uint8_t* src, const uint8_t* mask) {
  for (int i = 0; i < INET6_ADDRSTRLEN; i++) {
    dst[i] = src[i] & mask[i];
  }
}

static void add_downlink_match_ipv6(
    of13::FlowMod& downlink_fm, const struct in6_addr& ue_ip6, uint32_t port) {
  // Set match on uplink port and IP eth type
  static IPAddress mask("ffff:ffff:ffff:ffff::");
  of13::InPort uplink_port_match(port);
  downlink_fm.add_oxm_field(uplink_port_match);
  /* TODO-Made this local fix as it was not yet available in master
   * without this fix ovs rules are not getting created
   */
  of13::EthType ip6_type(0x86DD);
  downlink_fm.add_oxm_field(ip6_type);

  // Match UE IP destination
  struct in6_addr ue_ip6_masked;
  mask_ipv6_address(
      (uint8_t*) &ue_ip6_masked, (const uint8_t*) &ue_ip6, mask.getIPv6());

  of13::IPv6Dst ipv6_dst(IPAddress(ue_ip6_masked), mask);
  downlink_fm.add_oxm_field(ipv6_dst);
}

/**
 * Helper function to add dedicated brr flow match
 */
static void add_ded_brr_dl_match(
    of13::FlowMod& downlink_fm, const struct ip_flow_dl& flow, uint32_t port) {
  // Set match on uplink port and IP eth type
  of13::InPort uplink_port_match(port);
  downlink_fm.add_oxm_field(uplink_port_match);

  // Match UE IP destination
  if ((flow.set_params & DST_IPV4) || (flow.set_params & SRC_IPV4)) {
    of13::EthType ip_type(0x0800);
    downlink_fm.add_oxm_field(ip_type);

    if (flow.set_params & DST_IPV4) {
      of13::IPv4Dst ipv4_dst(flow.dst_ip.s_addr);
      downlink_fm.add_oxm_field(ipv4_dst);
    }

    // Match IP source
    if (flow.set_params & SRC_IPV4) {
      of13::IPv4Src ipv4_src(flow.src_ip.s_addr);
      downlink_fm.add_oxm_field(ipv4_src);
    }
  } else {
    of13::EthType ip_type(0x86DD);
    downlink_fm.add_oxm_field(ip_type);

    if (flow.set_params & DST_IPV6) {
      // Match on the prefix portion -- this is UE IPv6 address
      static IPAddress mask("ffff:ffff:ffff:ffff::");
      struct in6_addr dst_ip6_masked;
      mask_ipv6_address(
          (uint8_t*) &dst_ip6_masked, (const uint8_t*) &flow.dst_ip6,
          mask.getIPv6());

      of13::IPv6Dst ipv6_dst(IPAddress(dst_ip6_masked), mask);
      downlink_fm.add_oxm_field(ipv6_dst);
    }

    // Match IP source
    if (flow.set_params & SRC_IPV6) {
      of13::IPv6Src ipv6_src(IPAddress(flow.src_ip6));
      downlink_fm.add_oxm_field(ipv6_src);
    }
  }

  if (flow.set_params & IP_PROTO) {
    of13::IPProto proto(flow.ip_proto);
    downlink_fm.add_oxm_field(proto);
  }

  if (flow.set_params & TCP_SRC_PORT) {
    of13::TCPSrc tcp_src_port(flow.tcp_src_port);
    downlink_fm.add_oxm_field(tcp_src_port);
  }

  if (flow.set_params & TCP_DST_PORT) {
    of13::TCPDst tcp_dst_port(flow.tcp_dst_port);
    downlink_fm.add_oxm_field(tcp_dst_port);
  }

  if (flow.set_params & UDP_SRC_PORT) {
    of13::UDPSrc udp_src_port(flow.udp_src_port);
    downlink_fm.add_oxm_field(udp_src_port);
  }

  if (flow.set_params & UDP_DST_PORT) {
    of13::UDPDst udp_dst_port(flow.udp_dst_port);
    downlink_fm.add_oxm_field(udp_dst_port);
  }
}

/**
 * Helper function to add downlink flow action.
 */
void GTPApplication::add_tunnel_flow_action(
    uint32_t out_tei, uint32_t in_tei, std::string ue_imsi,
    struct in_addr remote_ip, uint32_t egress_gtp_port,
    fluid_base::OFConnection* connection, const OpenflowMessenger& messenger,
    of13::FlowMod downlink_fm, const std::string& flow_type, bool passthrough) {
  of13::ApplyActions apply_dl_inst;
  auto imsi = IMSIEncoder::compact_imsi(ue_imsi);

  // Set outgoing tunnel id and tunnel destination ip
  of13::SetFieldAction set_out_tunnel(new of13::TUNNELId(out_tei));
  apply_dl_inst.add_action(set_out_tunnel);
  of13::SetFieldAction set_tunnel_dst(
      new of13::TunnelIPv4Dst(remote_ip.s_addr));
  apply_dl_inst.add_action(set_tunnel_dst);

  int gtp_port = egress_gtp_port;
  if (gtp_port == 0) {
    gtp_port = GTPApplication::gtp0_port_num_;
  }
  of13::SetFieldAction set_tunnel_port(
      new of13::NXMRegX(TUNNEL_PORT_REG, gtp_port));
  apply_dl_inst.add_action(set_tunnel_port);

  of13::SetFieldAction set_tunnel_id(new of13::NXMRegX(TUNNEL_ID_REG, in_tei));
  apply_dl_inst.add_action(set_tunnel_id);

  EthAddress uplink_port(uplink_mac_);
  of13::SetFieldAction set_eth_dst(new of13::EthDst(uplink_port));
  apply_dl_inst.add_action(set_eth_dst);

  if (passthrough) {
    // Set register6, this is pipelineD internal details.
    // Once GTP app is moved to pipelineD we can remove this hack.
    of13::SetFieldAction set_passthrough(new of13::NXMRegX(6, 1));
    apply_dl_inst.add_action(set_passthrough);
  }
  // add imsi to packet metadata to pass to other tables
  add_imsi_metadata(apply_dl_inst, imsi);

  // Output to inout table
  of13::GoToTable goto_inst(NEXT_TABLE);

  downlink_fm.add_instruction(apply_dl_inst);
  downlink_fm.add_instruction(goto_inst);

  // Finally, send flow mod
  messenger.send_of_msg(downlink_fm, connection);
  OAILOG_DEBUG_UE(LOG_GTPV1U, imsi, "%s flow added\n", flow_type.c_str());
}

void GTPApplication::add_downlink_tunnel_flow_action(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    of13::FlowMod downlink_fm, bool passthrough, bool from_pgw) {
  uint32_t in_teid;
  if (from_pgw) {
    in_teid = ev.get_pgw_in_tei();
  } else {
    in_teid = ev.get_in_tei();
  }

  add_tunnel_flow_action(
      ev.get_out_tei(), in_teid, ev.get_imsi(), ev.get_enb_ip(),
      ev.get_enb_gtp_portno(), ev.get_connection(), messenger, downlink_fm,
      "S1 Downlink", passthrough);
}

void GTPApplication::add_downlink_tunnel_flow_ipv4(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port, bool passthrough, bool from_pgw) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, flow_priority);

  add_downlink_match(downlink_fm, ev.get_ue_ip(), ingress_port);
  add_downlink_tunnel_flow_action(
      ev, messenger, downlink_fm, passthrough, from_pgw);
}

void GTPApplication::add_downlink_tunnel_flow_ipv6(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port, bool passthrough, bool from_pgw) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, flow_priority);

  add_downlink_match_ipv6(
      downlink_fm, ev.get_ue_info().get_ipv6(), ingress_port);
  add_downlink_tunnel_flow_action(
      ev, messenger, downlink_fm, passthrough, from_pgw);
}

void GTPApplication::add_downlink_tunnel_flow_ded_brr(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port, bool passthrough, bool from_pgw) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, flow_priority);

  add_ded_brr_dl_match(downlink_fm, ev.get_dl_flow(), ingress_port);
  add_downlink_tunnel_flow_action(
      ev, messenger, downlink_fm, passthrough, from_pgw);
}

void GTPApplication::add_downlink_tunnel_flow(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port, bool passthrough, bool from_pgw) {
  if (ev.is_dl_flow_valid()) {
    add_downlink_tunnel_flow_ded_brr(
        ev, messenger, ingress_port, passthrough, from_pgw);
    return;
  }
  UeNetworkInfo ue_info = ev.get_ue_info();

  if (ue_info.is_ue_ipv4_addr_valid()) {
    add_downlink_tunnel_flow_ipv4(
        ev, messenger, ingress_port, passthrough, from_pgw);
  }
  if (ue_info.is_ue_ipv6_addr_valid()) {
    add_downlink_tunnel_flow_ipv6(
        ev, messenger, ingress_port, passthrough, from_pgw);
  }
}

/**
 * Helper function to add ARP actions
 */
void GTPApplication::add_downlink_arp_flow_action(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    of13::FlowMod downlink_fm) {
  auto imsi = IMSIEncoder::compact_imsi(ev.get_imsi());
  of13::ApplyActions apply_dl_inst;

  // add imsi to packet metadata to pass to other tables
  add_imsi_metadata(apply_dl_inst, imsi);

  // Output to inout table
  of13::GoToTable goto_inst(NEXT_TABLE);

  downlink_fm.add_instruction(apply_dl_inst);
  downlink_fm.add_instruction(goto_inst);

  // Finally, send flow mod
  messenger.send_of_msg(downlink_fm, ev.get_connection());
  OAILOG_DEBUG_UE(LOG_GTPV1U, imsi, "ARP flow added\n");
}

void GTPApplication::add_downlink_arp_flow(
    const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_ADD, flow_priority);

  add_downlink_arp_match(downlink_fm, ev.get_ue_ip(), ingress_port);

  add_downlink_arp_flow_action(ev, messenger, downlink_fm);
}

void GTPApplication::delete_downlink_tunnel_flow_ipv4(
    const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_DELETE, 0);
  // match all ports and groups
  downlink_fm.out_port(of13::OFPP_ANY);
  downlink_fm.out_group(of13::OFPG_ANY);

  add_downlink_match(downlink_fm, ev.get_ue_ip(), ingress_port);
  messenger.send_of_msg(downlink_fm, ev.get_connection());
}

void GTPApplication::delete_downlink_tunnel_flow_ded_brr(
    const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_DELETE, 0);
  // match all ports and groups
  downlink_fm.out_port(of13::OFPP_ANY);
  downlink_fm.out_group(of13::OFPG_ANY);

  add_ded_brr_dl_match(downlink_fm, ev.get_dl_flow(), ingress_port);
  messenger.send_of_msg(downlink_fm, ev.get_connection());
}

void GTPApplication::delete_downlink_tunnel_flow_ipv6(
    const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_DELETE, 0);
  // match all ports and groups
  downlink_fm.out_port(of13::OFPP_ANY);
  downlink_fm.out_group(of13::OFPG_ANY);

  add_downlink_match_ipv6(
      downlink_fm, ev.get_ue_info().get_ipv6(), ingress_port);
  messenger.send_of_msg(downlink_fm, ev.get_connection());
}

void GTPApplication::delete_downlink_tunnel_flow(
    const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  if (ev.is_dl_flow_valid()) {
    delete_downlink_tunnel_flow_ded_brr(ev, messenger, ingress_port);
    return;
  }
  UeNetworkInfo ue_info = ev.get_ue_info();

  if (ue_info.is_ue_ipv4_addr_valid()) {
    delete_downlink_tunnel_flow_ipv4(ev, messenger, ingress_port);
  }
  if (ue_info.is_ue_ipv6_addr_valid()) {
    delete_downlink_tunnel_flow_ipv6(ev, messenger, ingress_port);
  }
}

void GTPApplication::delete_downlink_arp_flow(
    const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  of13::FlowMod downlink_fm =
      messenger.create_default_flow_mod(0, of13::OFPFC_DELETE, 0);
  // match all ports and groups
  downlink_fm.out_port(of13::OFPP_ANY);
  downlink_fm.out_group(of13::OFPG_ANY);

  add_downlink_arp_match(downlink_fm, ev.get_ue_ip(), ingress_port);

  messenger.send_of_msg(downlink_fm, ev.get_connection());
}

void GTPApplication::discard_uplink_tunnel_flow(
    const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger) {
  of13::FlowMod uplink_fm = messenger.create_default_flow_mod(
      0, of13::OFPFC_ADD, DEFAULT_PRIORITY + 1);
  // match all ports and groups
  uplink_fm.out_port(of13::OFPP_ANY);
  uplink_fm.out_group(of13::OFPG_ANY);
  uplink_fm.cookie(cookie);
  uplink_fm.cookie_mask(cookie);

  add_tunnel_match(uplink_fm, gtp0_port_num_, ev.get_in_tei());

  messenger.send_of_msg(uplink_fm, ev.get_connection());
}

void GTPApplication::discard_downlink_tunnel_flow(
    const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  of13::FlowMod downlink_fm = messenger.create_default_flow_mod(
      0, of13::OFPFC_ADD, DEFAULT_PRIORITY + 1);
  // match all ports and groups
  downlink_fm.out_port(of13::OFPP_ANY);
  downlink_fm.out_group(of13::OFPG_ANY);
  downlink_fm.cookie(cookie + 1);
  downlink_fm.cookie_mask(cookie + 1);

  if (ev.is_dl_flow_valid()) {
    add_ded_brr_dl_match(downlink_fm, ev.get_dl_flow(), ingress_port);
  } else {
    add_downlink_match(downlink_fm, ev.get_ue_ip(), ingress_port);
  }

  messenger.send_of_msg(downlink_fm, ev.get_connection());
}

void GTPApplication::forward_uplink_tunnel_flow(
    const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod uplink_fm = messenger.create_default_flow_mod(
      0, of13::OFPFC_DELETE, flow_priority + 1);
  // match all ports and groups
  uplink_fm.out_port(of13::OFPP_ANY);
  uplink_fm.out_group(of13::OFPG_ANY);
  uplink_fm.cookie(cookie);
  uplink_fm.cookie_mask(cookie);

  add_tunnel_match(uplink_fm, gtp0_port_num_, ev.get_in_tei());

  messenger.send_of_msg(uplink_fm, ev.get_connection());
}

void GTPApplication::forward_downlink_tunnel_flow(
    const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
    uint32_t ingress_port) {
  uint32_t flow_priority =
      convert_precedence_to_priority(ev.get_dl_flow_precedence());
  of13::FlowMod downlink_fm = messenger.create_default_flow_mod(
      0, of13::OFPFC_DELETE, flow_priority + 1);
  // match all ports and groups
  downlink_fm.out_port(of13::OFPP_ANY);
  downlink_fm.out_group(of13::OFPG_ANY);
  downlink_fm.cookie(cookie + 1);
  downlink_fm.cookie_mask(cookie + 1);

  if (ev.is_dl_flow_valid()) {
    add_ded_brr_dl_match(downlink_fm, ev.get_dl_flow(), ingress_port);
  } else {
    add_downlink_match(downlink_fm, ev.get_ue_ip(), ingress_port);
  }

  messenger.send_of_msg(downlink_fm, ev.get_connection());
}

// Precedence in TFT and flow rule priority in OVS are inversely
// related. Rules with a low precedence value takes precedence,
// where 0 has the highest precedence. In OVS rules with high
// priority value takes precedence with the maximum value of
// 65535. Typical range of precedence is in [0,255] in line
// with the 8-bit TFT field for precedence in the current code.
// This implementation:
// - Allows 32-bit unsigned value for precedence velue, but truncates
//   precedence values higher than 65535 (i.e., 16 bits) to 65535.
// - Maps precendence values to priority values in [10, 65535].
// - Sets the minimum priority value to 10 in order to give GTP App
//   a sufficient margin to take priority over CP and management
//   related rules.
// - DEFAULT_PRECEDENCE always maps to a priority value of 10.
uint32_t GTPApplication::convert_precedence_to_priority(
    const uint32_t precedence) {
  uint32_t priority =
      (precedence < MAX_PRIORITY) ? (MAX_PRIORITY - precedence) : 0;
  if (priority < DEFAULT_PRIORITY) {
    priority = DEFAULT_PRIORITY;
  }
  return priority;
}

}  // namespace openflow
