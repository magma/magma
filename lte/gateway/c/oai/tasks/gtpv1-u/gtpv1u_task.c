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
/*! \file gtpv1u_task.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>

#include "log.h"
#include "assertions.h"
#include "intertask_interface.h"
#include "gtpv1u.h"
#include "gtpv1u_sgw_defs.h"
#include "gtp_tunnel_upf.h"
#include "pgw_ue_ip_address_alloc.h"
#include "intertask_interface_types.h"
#include "pgw_config.h"
#include "spgw_config.h"

const struct gtp_tunnel_ops* gtp_tunnel_ops;
static struct in_addr current_ue_net;
static int current_ue_net_mask;

//------------------------------------------------------------------------------
static void add_route_for_ue_block(struct in_addr ue_net, uint32_t mask) {
  if (ue_net.s_addr == htonl(INADDR_ANY) || mask == 0) {
    return;
  }
  // Use replace to avoid error related to existing routes.
  bstring system_cmd = bformat(
      "ip route replace %s/%u dev %s", inet_ntoa(ue_net), mask,
      gtp_tunnel_ops->get_dev_name());
  int ret = system((const char*) system_cmd->data);
  if (ret) {
    OAILOG_ERROR(
        LOG_GTPV1U, "ERROR in system command %s: %d at %s:%u\n",
        bdata(system_cmd), ret, __FILE__, __LINE__);
    bdestroy(system_cmd);
    return;
  }

  OAILOG_DEBUG(LOG_GTPV1U, "route updated: %s\n", bdata(system_cmd));
  bdestroy(system_cmd);
  // cache updated route.
  current_ue_net      = ue_net;
  current_ue_net_mask = mask;
}

static void del_route_for_ue_block(struct in_addr ue_net, uint32_t mask) {
  if (ue_net.s_addr == htonl(INADDR_ANY) || mask == 0) {
    return;
  }
  bstring system_cmd = bformat("ip route del %s/%u", inet_ntoa(ue_net), mask);
  int ret            = system((const char*) system_cmd->data);
  if (ret) {
    OAILOG_ERROR(
        LOG_GTPV1U, "ERROR in system command %s: %d at %s:%u\n",
        bdata(system_cmd), ret, __FILE__, __LINE__);
    bdestroy(system_cmd);
    return;
  }

  OAILOG_DEBUG(LOG_GTPV1U, "Deleted route%s\n", bdata(system_cmd));
  bdestroy(system_cmd);
  current_ue_net_mask = 0;
}

/**
 * Check if _addr is in given subnet (_net/mask)
 */
static bool ue_ip_is_in_subnet(
    struct in_addr _net, int mask, struct in_addr _addr) {
  if (mask == 0) {
    // This is first time checking for subnect.
    return false;
  }
  uint32_t net  = ntohl(_net.s_addr);
  uint32_t addr = ntohl(_addr.s_addr);
  if (addr < net) {
    return false;
  }
  uint32_t no_of_ips = 1 << (32 - mask);
  if (net + no_of_ips < addr) {
    return false;
  }

  return true;
}

//------------------------------------------------------------------------------
int gtpv1u_init(
    spgw_state_t* spgw_state_p, spgw_config_t* spgw_config,
    bool persist_state) {
  int rv = 0;
  struct in_addr netaddr;
  uint32_t netmask = 0;

  OAILOG_DEBUG(LOG_GTPV1U, "Initializing GTPV1U interface\n");

  // Init gtp_tunnel_ops
#if ENABLE_OPENFLOW
  // If pipeline config is enabled initialize userplane ops
  if (spgw_config->sgw_config.ovs_config.pipelined_managed_tbl0) {
    OAILOG_INFO(LOG_GTPV1U, "Initializing upf classifier for gtp apps");
    gtp_tunnel_ops = upf_gtp_tunnel_ops_init_openflow();
  } else {
    OAILOG_DEBUG(LOG_GTPV1U, "Initializing gtp_tunnel_ops_openflow\n");
    gtp_tunnel_ops = gtp_tunnel_ops_init_openflow();
  }
#else
  OAILOG_DEBUG(LOG_GTPV1U, "Initializing gtp_tunnel_ops_libgtpnl\n");
  gtp_tunnel_ops = gtp_tunnel_ops_init_libgtpnl();
#endif

  if (gtp_tunnel_ops == NULL) {
    OAILOG_CRITICAL(LOG_GTPV1U, "ERROR in initializing gtp_tunnel_ops\n");
    return -1;
  }

  // Reset GTP tunnel states
  rv = gtp_tunnel_ops->reset();
  if (rv != 0) {
    OAILOG_CRITICAL(LOG_GTPV1U, "ERROR clean existing gtp states.\n");
    return -1;
  }

  if (spgw_config->pgw_config.enable_nat) {
    rv = get_ip_block(&netaddr, &netmask);
    if (rv != 0) {
      OAILOG_CRITICAL(
          LOG_GTPV1U, "ERROR in getting assigned IP block from mobilityd\n");
      return -1;
    }
  } else {
    // Allow All IPs in Non-NAT case.
    netaddr.s_addr = INADDR_ANY;
    netmask        = 0;
  }

  // Init GTP device, using the same MTU as SGi.
  gtp_tunnel_ops->init(
      &netaddr, netmask, spgw_config->pgw_config.ipv4.mtu_SGI,
      &spgw_state_p->gtpv1u_data.fd0, &spgw_state_p->gtpv1u_data.fd1u,
      persist_state);

  // END-GTP quick integration only for evaluation purpose

  // Add route to avoid updating routing during UE attach.
  add_route_for_ue_block(netaddr, netmask);

  OAILOG_DEBUG(LOG_GTPV1U, "Initializing GTPV1U interface: DONE\n");
  return 0;
}

int gtpv1u_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl, char* apn) {
  OAILOG_DEBUG(LOG_GTPV1U, "Add tunnel ue %s", inet_ntoa(ue));

  if (spgw_config.pgw_config.enable_nat) {
    if (!ue_ip_is_in_subnet(current_ue_net, current_ue_net_mask, ue)) {
      struct in_addr netaddr;
      uint32_t netmask = 0;

      // get new block from mobility.
      int rv = get_ip_block(&netaddr, &netmask);
      if (rv != 0) {
        OAILOG_INFO(
            LOG_GTPV1U,
            "ERROR in getting assigned IP block from mobilityd,"
            "could not set the route to UE network\n");
      } else {
        // add the route if needed
        OAILOG_INFO(
            LOG_GTPV1U, "Got new ip-block %s/%d", inet_ntoa(netaddr), netmask);
        if (netaddr.s_addr != current_ue_net.s_addr ||
            current_ue_net_mask != netmask) {
          del_route_for_ue_block(current_ue_net, current_ue_net_mask);
          add_route_for_ue_block(netaddr, netmask);
        }
      }
    }
  }

  return gtp_tunnel_ops->add_tunnel(
      ue, ue_ipv6, vlan, enb, i_tei, o_tei, imsi, flow_dl, flow_precedence_dl,
      apn);
}

int gtpv1u_add_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in_addr pgw, uint32_t i_tei, uint32_t o_tei, uint32_t pgw_i_tei,
    uint32_t pgw_o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl) {
  OAILOG_DEBUG(LOG_GTPV1U, "Add S8 tunnel ue %s", inet_ntoa(ue));
  if (gtp_tunnel_ops->add_s8_tunnel) {
    return gtp_tunnel_ops->add_s8_tunnel(
        ue, ue_ipv6, vlan, enb, pgw, i_tei, o_tei, pgw_i_tei, pgw_o_tei, imsi,
        flow_dl, flow_precedence_dl);
  } else {
    return -EINVAL;
  }
}

int gtpv1u_del_s8_tunnel(
    struct in_addr enb, struct in_addr pgw, struct in_addr ue,
    struct in6_addr* ue_ipv6, uint32_t i_tei, uint32_t o_tei,
    struct ip_flow_dl* flow_dl) {
  OAILOG_DEBUG(LOG_GTPV1U, "Del S8 tunnel ue %s", inet_ntoa(ue));
  if (gtp_tunnel_ops->del_s8_tunnel) {
    return gtp_tunnel_ops->del_s8_tunnel(
        enb, pgw, ue, ue_ipv6, i_tei, o_tei, flow_dl);
  } else {
    return -EINVAL;
  }
}

//------------------------------------------------------------------------------
void gtpv1u_exit(void) {
  gtp_tunnel_ops->uninit();
}
