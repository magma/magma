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

//#include "lte/gateway/c/core/common/assertions.h"
//#include "lte/gateway/c/core/oai/common/log.h"
//#include "lte/gateway/c/core/oai/include/pgw_config.h"
//#include "lte/gateway/c/core/oai/include/spgw_config.h"
//#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
//#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
//#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtp_tunnel_upf.h"
#include "lte/gateway/c/session_manager/gtp/gtpv1u.h"
// #include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u_sgw_defs.h"
//#include "lte/gateway/c/core/oai/tasks/sgw/pgw_ue_ip_address_alloc.h"
#include "lte/gateway/c/session_manager/bstr/bstrlib.h"
#include "lte/gateway/c/session_manager/gtp/gtpv1u_task.h"
const struct gtp_tunnel_ops* gtp_tunnel_ops;
static struct in_addr current_ue_net;
static int current_ue_net_mask;


//------------------------------------------------------------------------------
void add_route_for_ue_block(struct in_addr ue_net, uint32_t mask) {
  if (ue_net.s_addr == htonl(INADDR_ANY) || mask == 0) {
    return;
  }
  // Use replace to avoid error related to existing routes.
  bstring system_cmd =
      bformat("ip route replace %s/%u dev %s", inet_ntoa(ue_net), mask,
              gtp_tunnel_ops->get_dev_name());
  int ret = system((const char*)system_cmd->data);
  if (ret) {
    //bdestroy(system_cmd);
    return;
  }

  //OAILOG_DEBUG(LOG_GTPV1U, "route updated: %s\n", bdata(system_cmd));
  //bdestroy(system_cmd);
  // cache updated route.
  current_ue_net = ue_net;
  current_ue_net_mask = mask;
}

static void del_route_for_ue_block(struct in_addr ue_net, uint32_t mask) {
  if (ue_net.s_addr == htonl(INADDR_ANY) || mask == 0) {
    return;
  }
  bstring system_cmd = bformat("ip route del %s/%u", inet_ntoa(ue_net), mask);
  int ret = system((const char*)system_cmd->data);
  if (ret) {
    //bdestroy(system_cmd);
    return;
  }

  //OAILOG_DEBUG(LOG_GTPV1U, "Deleted route%s\n", bdata(system_cmd));
  //bdestroy(system_cmd);
  current_ue_net_mask = 0;
}

/**
 * Check if _addr is in given subnet (_net/mask)
 */
static bool ue_ip_is_in_subnet(struct in_addr _net, int mask,
                               struct in_addr _addr) {
  if (mask == 0) {
    // This is first time checking for subnect.
    return false;
  }
  uint32_t net = ntohl(_net.s_addr);
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
int gtpv1u_init_1() {
  int rv = 0;
  struct in_addr netaddr;
  uint32_t netmask = 0;
  int fd = 1;

  //OAILOG_DEBUG(LOG_GTPV1U, "Initializing GTPV1U interface\n");

  // Init gtp_tunnel_ops
  // If pipeline config is enabled initialize userplane ops

  gtp_tunnel_ops = gtp_tunnel_ops_init_openflow();

  if (gtp_tunnel_ops == NULL) {
    return -1;
  }

  // Reset GTP tunnel states
  rv = gtp_tunnel_ops->reset();
  if (rv != 0) {
    return -1;
  }

  // === HARDCODE gtp_tunnel_ops
  gtp_tunnel_ops->init(&netaddr, netmask, 1500, &fd, &fd, 0);

  return 0;
}

//=== HARDCODE
int gtpv1u_add_tunnel(struct in_addr ue, struct in6_addr* ue_ipv6, int vlan,
                      struct in_addr enb, struct in6_addr* enb_ipv6,
                      uint32_t i_tei, uint32_t o_tei, char* imsi,
                      struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl,
                      char* apn) {
  return gtp_tunnel_ops->add_tunnel(ue, ue_ipv6, vlan, enb, enb_ipv6, i_tei,
                                    o_tei, imsi, flow_dl, flow_precedence_dl,
                                    apn);
}

//------------------------------------------------------------------------------
void gtpv1u_exit(void) { gtp_tunnel_ops->uninit(); }
