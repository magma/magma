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
/*! \file gtpv1u.h
 * \brief
 * \author Sebastien ROUX, Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#ifndef FILE_GTPV1_U_SEEN
#define FILE_GTPV1_U_SEEN

#include <arpa/inet.h>
#include <net/if.h>
#include "sgw_ie_defs.h"

#define GTPU_HEADER_OVERHEAD_MAX 64

/*
 * downlink flow description for a dedicated bearer
 */

#define SRC_IPV4 0x1
#define DST_IPV4 0x2
#define TCP_SRC_PORT 0x4
#define TCP_DST_PORT 0x8
#define UDP_SRC_PORT 0x10
#define UDP_DST_PORT 0x20
#define IP_PROTO 0x40
#define SRC_IPV6 0x80
#define DST_IPV6 0x100

// This is the default precedence value for flow rules.
// A flow rule with precedence value 0 takes precedence over
// all other flows with higher value. Flow rules that use
// the default precedence have the lowest priority.
#define DEFAULT_PRECEDENCE 65535
// This is the maximum priority an OVS rule can take.
// A high priority value takes precedence over the lower value.
// For equal priority values, the behavior is undefined, but
// it is not atypical to see a previously installed rule to take
// precedence over latter rules.
#define MAX_PRIORITY 65535

struct ip_flow_dl {
  uint32_t set_params;
  uint16_t tcp_dst_port;
  uint16_t tcp_src_port;
  uint16_t udp_dst_port;
  uint16_t udp_src_port;
  uint8_t ip_proto;
  union {
    struct {
      struct in_addr dst_ip;
      struct in_addr src_ip;
    };
    struct {
      struct in6_addr dst_ip6;
      struct in6_addr src_ip6;
    };
  };
};

/*
 * This structure defines the management hooks for GTP tunnels and paging
 * support. The following hooks can be defined; unless noted otherwise, they are
 * optional and can be filled with a null pointer.
 *
 * int (*init)(struct in_addr *ue_net, uint32_t mask,
 *             int mtu, int *fd0, int *fd1u);
 *     This function is called when initializing GTP network device. How to use
 *     these input parameters are defined by the actual function
 * implementations.
 *         @ue_net: subnet assigned to UEs
 *         @mask: network mask for the UE subnet
 *         @mtu: MTU for the GTP network device.
 *         @fd0: socket file descriptor for GTPv0.
 *         @fd1u: socket file descriptor for GTPv1u.
 *
 * int (*uninit)(void);
 *     This function is called to destroy GTP network device.
 *
 * int (*reset)(void);
 *     This function is called to reset the GTP network device to clean state.
 *
 * int (*add_tunnel)(struct in_addr ue, struct in_addr enb, uint32_t i_tei,
 * uint32_t o_tei, Imsi_t imsi); Add a gtp tunnel.
 *         @ue: UE IP address
 *         @enb: eNB IP address
 *         @i_tei: RX GTP Tunnel ID
 *         @o_tei: TX GTP Tunnel ID.
 *         @imsi: UE IMSI
 *         @flow_dl: Downlink flow rule
 *         @flow_precedence_dl: Downlink flow rule precedence
 *
 * int (*del_tunnel)(uint32_t i_tei, uint32_t o_tei);
 *     Delete a gtp tunnel.
 *         @ue: UE IP address
 *         @i_tei: RX GTP Tunnel ID
 *         @o_tei: TX GTP Tunnel ID.
 *
 * int (*discard_data_on_tunnel)(struct in_addr ue, uint32_t i_tei);
 *         @ue: UE IP address
 *         @i_tei: RX GTP Tunnel ID
 *
 * int (*forward_data_on_tunnel)(struct in_addr ue, uint32_t i_tei);
 *         @ue: UE IP address
 *         @i_tei: RX GTP Tunnel ID
 *         @flow_dl: Downlink flow rule
 *         @flow_precedence_dl: Downlink flow rule precedence
 *
 * int (*add_paging_rule)(struct in_addr ue);
 *        @ue: UE IP address
 *
 * int (*send_end_marker) (struct in_addr enb, uint32_t i_tei);
 *        @enb: eNB IP address
 *        @i_tei: RX GTP Tunnel ID
 */
struct gtp_tunnel_ops {
  int (*init)(
      struct in_addr* ue_net, uint32_t mask, int mtu, int* fd0, int* fd1u,
      bool persist_state);
  int (*uninit)(void);
  int (*reset)(void);
  int (*add_tunnel)(
      struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
      uint32_t i_tei, uint32_t o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
      uint32_t flow_precedence_dl, char* apn);
  int (*del_tunnel)(
      struct in_addr enb, struct in_addr ue, struct in6_addr* ue_ipv6,
      uint32_t i_tei, uint32_t o_tei, struct ip_flow_dl* flow_dl);
  int (*add_s8_tunnel)(
      struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
      struct in_addr pgw, uint32_t i_tei, uint32_t o_tei, uint32_t pgw_i_tei,
      uint32_t pgw_o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
      uint32_t flow_precedence_dl);
  int (*del_s8_tunnel)(
      struct in_addr enb, struct in_addr pgw, struct in_addr ue,
      struct in6_addr* ue_ipv6, uint32_t i_tei, uint32_t o_tei,
      struct ip_flow_dl* flow_dl);
  int (*discard_data_on_tunnel)(
      struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
      struct ip_flow_dl* flow_dl);
  int (*forward_data_on_tunnel)(
      struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
      struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl);
  int (*add_paging_rule)(struct in_addr ue);
  int (*delete_paging_rule)(struct in_addr ue);
  int (*send_end_marker)(struct in_addr enbode, uint32_t i_tei);
  const char* (*get_dev_name)(void);
};

#if ENABLE_OPENFLOW
const struct gtp_tunnel_ops* gtp_tunnel_ops_init_openflow(void);
#else
const struct gtp_tunnel_ops* gtp_tunnel_ops_init_libgtpnl(void);
#endif

int gtpv1u_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl, char* apn);

int gtpv1u_add_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in_addr pgw, uint32_t i_tei, uint32_t o_tei, uint32_t pgw_i_tei,
    uint32_t pgw_o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl);

int gtpv1u_del_s8_tunnel(
    struct in_addr enb, struct in_addr pgw, struct in_addr ue,
    struct in6_addr* ue_ipv6, uint32_t i_tei, uint32_t o_tei,
    struct ip_flow_dl* flow_dl);
#endif /* FILE_GTPV1_U_SEEN */
