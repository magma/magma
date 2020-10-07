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

#include <stdint.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sys/types.h>
#include <arpa/inet.h>
#include <net/if.h>

#include <libgtpnl/gtp.h>
#include <libgtpnl/gtpnl.h>
#include <libmnl/libmnl.h>
#include <errno.h>

#include "log.h"
#include "common_defs.h"
#include "gtpv1u.h"
#include "gtpv1u_sgw_defs.h"

extern struct gtp_tunnel_ops gtp_tunnel_ops;

static struct {
  int genl_id;
  struct mnl_socket* nl;
  bool is_enabled;
} gtp_nl;

#define GTP_DEVNAME "gtp0"

int libgtpnl_init(
    struct in_addr* ue_net, uint32_t mask, int mtu, int* fd0, int* fd1u,
    bool persist_state) {
  // we don't need GTP v0, but interface with kernel requires 2 file descriptors
  *fd0                            = socket(AF_INET, SOCK_DGRAM, 0);
  *fd1u                           = socket(AF_INET, SOCK_DGRAM, 0);
  struct sockaddr_in sockaddr_fd0 = {
      .sin_family = AF_INET,
      .sin_port   = htons(3386),
      .sin_addr =
          {
              .s_addr = INADDR_ANY,
          },
  };
  struct sockaddr_in sockaddr_fd1 = {
      .sin_family = AF_INET,
      .sin_port   = htons(GTPV1U_UDP_PORT),
      .sin_addr =
          {
              .s_addr = INADDR_ANY,
          },
  };

  if (bind(*fd0, (struct sockaddr*) &sockaddr_fd0, sizeof(sockaddr_fd0)) < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "bind GTPv0 port");
    return RETURNerror;
  }
  if (bind(*fd1u, (struct sockaddr*) &sockaddr_fd1, sizeof(sockaddr_fd1)) < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "bind S1U port");
    return RETURNerror;
  }

  if (gtp_dev_create(-1, GTP_DEVNAME, *fd0, *fd1u) < 0) {
    OAILOG_ERROR(
        LOG_GTPV1U, "Cannot create GTP tunnel device: %s\n", strerror(errno));
    return RETURNerror;
  }
  gtp_nl.is_enabled = true;

  gtp_nl.nl = genl_socket_open();
  if (gtp_nl.nl == NULL) {
    OAILOG_ERROR(LOG_GTPV1U, "Cannot create genetlink socket\n");
    return RETURNerror;
  }
  gtp_nl.genl_id = genl_lookup_family(gtp_nl.nl, "gtp");
  if (gtp_nl.genl_id < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "Cannot lookup GTP genetlink ID\n");
    return RETURNerror;
  }
  OAILOG_NOTICE(
      LOG_GTPV1U, "Using the GTP kernel mode (genl ID is %d)\n",
      gtp_nl.genl_id);

  bstring system_cmd = bformat("ip link set dev %s mtu %u", GTP_DEVNAME, mtu);
  int ret            = system((const char*) system_cmd->data);
  if (ret) {
    OAILOG_ERROR(
        LOG_GTPV1U, "ERROR in system command %s: %d at %s:%u\n",
        bdata(system_cmd), ret, __FILE__, __LINE__);
    bdestroy(system_cmd);
    return RETURNerror;
  }
  bdestroy(system_cmd);

  struct in_addr ue_gw;
  ue_gw.s_addr = ue_net->s_addr | htonl(1);
  system_cmd =
      bformat("ip addr add %s/%u dev %s", inet_ntoa(ue_gw), mask, GTP_DEVNAME);
  ret = system((const char*) system_cmd->data);
  if (ret) {
    OAILOG_ERROR(
        LOG_GTPV1U, "ERROR in system command %s: %d at %s:%u\n",
        bdata(system_cmd), ret, __FILE__, __LINE__);
    bdestroy(system_cmd);
    return RETURNerror;
  }
  bdestroy(system_cmd);

  OAILOG_DEBUG(
      LOG_GTPV1U, "Setting route to reach UE net %s via %s\n",
      inet_ntoa(*ue_net), GTP_DEVNAME);

  if (gtp_dev_config(GTP_DEVNAME, ue_net, mask) < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "Cannot add route to reach network\n");
    return RETURNerror;
  }

  OAILOG_NOTICE(LOG_GTPV1U, "GTP kernel configured\n");

  return RETURNok;
}

int libgtpnl_uninit(void) {
  if (!gtp_nl.is_enabled) return -1;

  return gtp_dev_destroy(GTP_DEVNAME);
}

int libgtpnl_reset(void) {
  int rv = 0;
  rv     = system("rmmod gtp");
  rv     = system("modprobe gtp");
  return rv;
}

int libgtpnl_add_tunnel(
    struct in_addr ue, int vlan, struct in_addr enb, uint32_t i_tei, uint32_t o_tei,
    Imsi_t imsi, struct ipv4flow_dl* flow_dl) {
  struct gtp_tunnel* t;
  int ret;

  if (!gtp_nl.is_enabled) return RETURNok;

  t = gtp_tunnel_alloc();
  if (t == NULL) return RETURNerror;

  gtp_tunnel_set_ifidx(t, if_nametoindex(GTP_DEVNAME));
  gtp_tunnel_set_version(t, 1);
  gtp_tunnel_set_ms_ip4(t, &ue);
  gtp_tunnel_set_sgsn_ip4(t, &enb);
  gtp_tunnel_set_i_tei(t, i_tei);
  gtp_tunnel_set_o_tei(t, o_tei);

  ret = gtp_add_tunnel(gtp_nl.genl_id, gtp_nl.nl, t);
  gtp_tunnel_free(t);

  return ret;
}

int libgtpnl_del_tunnel(__attribute__((unused)) struct in_addr enb,
    __attribute__((unused)) struct in_addr ue, uint32_t i_tei, uint32_t o_tei,
    struct ipv4flow_dl* flow_dl) {
  struct gtp_tunnel* t;
  int ret;

  if (!gtp_nl.is_enabled) return RETURNok;

  t = gtp_tunnel_alloc();
  if (t == NULL) return RETURNerror;

  gtp_tunnel_set_ifidx(t, if_nametoindex(GTP_DEVNAME));
  gtp_tunnel_set_version(t, 1);
  // looking at kernel/drivers/net/gtp.c: not needed gtp_tunnel_set_ms_ip4(t,
  // &ue); looking at kernel/drivers/net/gtp.c: not needed
  // gtp_tunnel_set_sgsn_ip4(t, &enb);
  gtp_tunnel_set_i_tei(t, i_tei);
  gtp_tunnel_set_o_tei(t, o_tei);

  ret = gtp_del_tunnel(gtp_nl.genl_id, gtp_nl.nl, t);
  gtp_tunnel_free(t);

  return ret;
}

/**
 * Send packet marker to enodeB @enb for tunnel @tei.
 */
int libgtpnl_send_end_marker(struct in_addr enb, uint32_t tei) {
  return -EOPNOTSUPP;
}

const char* libgtpnl_get_dev_name() {
  return GTP_DEVNAME;
}

static const struct gtp_tunnel_ops libgtpnl_ops = {
    .init            = libgtpnl_init,
    .uninit          = libgtpnl_uninit,
    .reset           = libgtpnl_reset,
    .add_tunnel      = libgtpnl_add_tunnel,
    .del_tunnel      = libgtpnl_del_tunnel,
    .send_end_marker = libgtpnl_send_end_marker,
    .get_dev_name    = libgtpnl_get_dev_name,
};

const struct gtp_tunnel_ops* gtp_tunnel_ops_init_libgtpnl(void) {
  return &libgtpnl_ops;
}
