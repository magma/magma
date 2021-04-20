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

#pragma once

#include <gmp.h>  // gross but necessary to link spgw_config.h

#include "OpenflowController.h"
#include "gtpv1u.h"

namespace openflow {

/**
 * GTPApplication handles external callbacks to add/delete tunnel flows for a
 * UE when it connects
 */
class GTPApplication : public Application {
 public:
  GTPApplication(
      const std::string& uplink_mac, uint32_t gtp_port_num,
      uint32_t mtr_port_num, uint32_t internal_sampling_port_num,
      uint32_t internal_sampling_fwd_tbl_num, uint32_t uplink_port_num);

 private:
  /**
   * Main callback event required by inherited Application class. Whenever
   * the controller gets an event like packet in or switch up, it will pass
   * it to the application here
   *
   * @param ev (in) - pointer to some subclass of ControllerEvent that occurred
   */
  virtual void event_callback(
      const ControllerEvent& ev, const OpenflowMessenger& messenger);

  void install_internal_pkt_fwd_flow(
      fluid_base::OFConnection* ofconn, const OpenflowMessenger& messenger,
      uint32_t port, uint32_t next_table);

  /*
   * Add uplink flow from UE to internet
   * @param ev - AddGTPTunnelEvent containing ue ip, enb ip, and tunnel id's
   */
  void add_uplink_tunnel_flow(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger);

  /*
   * Add downlink flow from internet to UE
   * @param ev - AddGTPTunnelEvent containing ue ip, enb ip, and tunnel id's
   */
  void add_downlink_tunnel_flow(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number, bool passthrough, bool from_pgw);

  /*
   * Add downlink tunnel flow for S8
   */
  void add_uplink_s8_tunnel_flow(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger);

  /*
   * Remove uplink tunnel flow on disconnect
   * @param ev - DeleteGTPTunnelEvent containing ue ip, and inbound tei
   */
  void delete_uplink_tunnel_flow(
      const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger);

  /*
   * Remove downlink tunnel flow on disconnect
   * @param ev - DeleteGTPTunnelEvent containing ue ip, and inbound tei
   */
  void delete_downlink_tunnel_flow(
      const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);
  /*
   * Discard downlink data received for UE IP during UE suspended state
   * @param ev - HandleDataOnGTPTunnelEvent containing ue ip, and inbound tei
   */
  void discard_downlink_tunnel_flow(
      const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);
  /*
   * Discard uplink data received for sgw-S1U-teid during UE suspended state
   * @param ev - HandleDataOnGTPTunnelEvent containing ue ip, and inbound tei
   */
  void discard_uplink_tunnel_flow(
      const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger);
  /*
   * Remove the rule inserted to discard data for UE IP  UE suspended state
   * And Forward data existing rule
   * @param ev - HandleDataOnGTPTunnelEvent containing ue ip, and inbound tei
   */
  void forward_downlink_tunnel_flow(
      const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);
  /*
   * Remove the rule inserted to discard data for UE in suspended state
   * And Forward data existing rule
   * @param ev - HandleDataOnGTPTunnelEvent containing ue ip, and inbound tei
   */
  void forward_uplink_tunnel_flow(
      const HandleDataOnGTPTunnelEvent& ev, const OpenflowMessenger& messenger);
  /*
   * Convert flow rule precedence to OF flow priority
   * @param precedence - can be between 0 and DEFAULT_PRECEDENCE
   *
   * @return uint32_t flow priority (minimum value set to DEFAULT_PRIORITY)
   */
  uint32_t convert_precedence_to_priority(const uint32_t precedence);

  /**
   * Add Arp flow for UE IP address
   * @param ev - AddGTPTunnelEvent containing ue ip
   */
  void add_downlink_arp_flow(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);

  /**
   * Delete arp flow of the UE.
   * @param ev - AddGTPTunnelEvent containing ue ip
   */
  void delete_downlink_arp_flow(
      const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);

  /**
   * Add uplink port match to UL flows
   * @param uplink_fm OF flow mod msg
   * @param gtp_port GTP port from event
   * @param i_tei tunnel id.
   */
  void add_tunnel_match(
      of13::FlowMod& uplink_fm, uint32_t gtp_port, uint32_t i_tei);

 private:
  static const uint32_t DEFAULT_PRIORITY = 10;
  static const std::string GTP_PORT_MAC;
  static const uint16_t NEXT_TABLE      = 1;
  static const uint32_t LOW_PRIORITY    = 0;
  static const uint32_t TUNNEL_PORT_REG = 8;
  static const uint32_t TUNNEL_ID_REG   = 9;

  const std::string uplink_mac_;
  const uint32_t gtp0_port_num_;
  // Internal port number for monitoring service
  const uint32_t mtr_port_num_;
  // Internal port for sampling internal ipfix packets
  const uint32_t internal_sampling_port_num_;
  const uint32_t internal_sampling_fwd_tbl_num_;
  /* cookie is added to identify the rules enforced for the flow controller
   * Initialising with 1
   */
  const uint64_t cookie = 1;

  const uint32_t uplink_port_num_;

  void add_downlink_arp_flow_action(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      of13::FlowMod downlink_fm);

  void add_tunnel_flow_action(
      uint32_t out_tei, uint32_t in_tei, std::string ue_imsi,
      struct in_addr remote_ip, uint32_t egress_gtp_port,
      fluid_base::OFConnection* connection, const OpenflowMessenger& messenger,
      of13::FlowMod downlink_fm, const std::string& flow_type,
      bool passthrough);

  void add_downlink_tunnel_flow_action(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      of13::FlowMod downlink_fm, bool passthrough, bool from_pgw);

  void add_downlink_tunnel_flow_ipv4(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number, bool passthrough, bool from_pgw);
  void add_downlink_tunnel_flow_ipv6(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number, bool passthrough, bool from_pgw);
  void add_downlink_tunnel_flow_ded_brr(
      const AddGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number, bool passthrough, bool from_pgw);

  void delete_downlink_tunnel_flow_ipv4(
      const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);
  void delete_downlink_tunnel_flow_ipv6(
      const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);
  void delete_downlink_tunnel_flow_ded_brr(
      const DeleteGTPTunnelEvent& ev, const OpenflowMessenger& messenger,
      uint32_t port_number);
};

}  // namespace openflow
