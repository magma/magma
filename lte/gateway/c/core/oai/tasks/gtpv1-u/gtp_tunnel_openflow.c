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

#include <assert.h>
#include <errno.h>
#include <stdint.h>
#include <netinet/in.h>
#include <stdlib.h>

#include "assertions.h"
#include "bstrlib.h"
#include "log.h"
#include "gtpv1u.h"
#include "ControllerMain.h"
#include "3gpp_23.003.h"
#include "spgw_config.h"

extern struct gtp_tunnel_ops gtp_tunnel_ops;

// Tunnel port related functionality
static const char* ovs_gtp_type;

#define MAX_GTP_PORT_NAME_LENGTH 15

#define INIT_GTP_TABLE_SIZE 64
#define MAX_GTP_TABLE_SIZE 1024

struct gtp_portno {
  // GTP port name is limited at 14 char
  char name[MAX_GTP_PORT_NAME_LENGTH];
  // zero port number is unknown port
  uint32_t portno;
};

struct gtp_portno_record {
  struct gtp_portno* arr;
  int allocated;
  int size;
};

static struct gtp_portno_record gtp_portno_rec;

/**
 * Generate GTP port name from eNodeB IP address
 */
static void ip_addr_to_gtp_port_name(struct in_addr enb_addr, char name[]) {
  int rc;
  rc = snprintf(
      name, MAX_GTP_PORT_NAME_LENGTH, "g_%x", (uint32_t) enb_addr.s_addr);
  assert(rc > 0);
}

static uint32_t search_portno_records(char* name) {
  int i;
  for (i = 0; i < gtp_portno_rec.allocated; i++) {
    if (!strncmp(name, gtp_portno_rec.arr[i].name, MAX_GTP_PORT_NAME_LENGTH)) {
      return gtp_portno_rec.arr[i].portno;
    }
  }
  return 0;
}

/**
 * Implements basic port number map with linear seach.
 */
static void add_portno_rec(char* port_name, uint32_t portno) {
  if (gtp_portno_rec.allocated == gtp_portno_rec.size) {
    int new_size = gtp_portno_rec.size * 2;
    if (new_size > MAX_GTP_TABLE_SIZE) {
      // in case of unpected increase in table size, flush
      // all records.
      new_size = INIT_GTP_TABLE_SIZE;
    }
    struct gtp_portno* new_arr = calloc(new_size, sizeof(struct gtp_portno));
    if (!new_arr) {
      return;
    }
    if (new_size > gtp_portno_rec.size) {
      // In case of table size increase copy existing records.
      memcpy(
          new_arr, gtp_portno_rec.arr,
          sizeof(gtp_portno_rec.arr[0]) * gtp_portno_rec.size);
    } else {
      // Flush all records if size resets to init-size.
      gtp_portno_rec.allocated = 0;
    }
    free(gtp_portno_rec.arr);
    gtp_portno_rec.arr  = new_arr;
    gtp_portno_rec.size = new_size;
  }
  // Now we shld have space to add new port.
  int i;

  for (i = 0; i < gtp_portno_rec.size; i++) {
    struct gtp_portno* rec = &gtp_portno_rec.arr[i];

    if (!rec->portno) {
      strncpy(rec->name, port_name, MAX_GTP_PORT_NAME_LENGTH);
      rec->portno = portno;
      gtp_portno_rec.allocated++;
      return;
    }
  }
  // port number caching is best effort.
}

/**
 * Read GTP tunnel port number from OVSDB.
 */
static uint32_t get_gtp_port_no(char port_name[]) {
  FILE* fp;
  uint32_t port_no = 0;
  char ovsdb_dump[256];

  /* Open the command for reading. */
  fp = popen("sudo ovsdb-client dump Interface name ofport", "r");
  if (fp == NULL) {
    OAILOG_ERROR(LOG_GTPV1U, "could not read ovsdb");
    return 0;
  }

  /* Read the output a line at a time - output it. */
  while (fgets(ovsdb_dump, sizeof(ovsdb_dump), fp) != NULL) {
    OAILOG_DEBUG(LOG_GTPV1U, "ovsdb: %s\n", ovsdb_dump);
    char* p = strstr(ovsdb_dump, port_name);
    if (p) {
      int len = strlen(port_name);
      // ovsdb dump has port number after portname separated by whitespaces.
      port_no = atoi(&p[len + 1]);
      break;
    }
  }

  pclose(fp);
  return port_no;
}

/**
 * Create GTP tunnel using OVS tool
 */
static uint32_t create_gtp_port(struct in_addr enb_addr, char port_name[]) {
  char gtp_port_create[512];
  int rc;
  rc = snprintf(
      gtp_port_create, sizeof(gtp_port_create),
      "sudo ovs-vsctl --may-exist add-port gtp_br0 %s -- set Interface %s "
      "type=%s "
      "options:remote_ip=%s options:key=flow",
      port_name, port_name, ovs_gtp_type, inet_ntoa(enb_addr));
  if (rc < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "gtp-port create: format error %d", rc);
    return rc;
  }
  rc = system(gtp_port_create);
  if (rc != 0) {
    // ignore failures. we can always fallback to gtp0 for GTP tunnel traffic.
    OAILOG_ERROR(
        LOG_GTPV1U, "gtp port create: [%s] failed: %d", gtp_port_create, rc);
  } else {
    OAILOG_DEBUG(
        LOG_GTPV1U, "gtp port create done: for ENB: %s ", inet_ntoa(enb_addr));
  }

  return get_gtp_port_no(port_name);
}

/**
 * seach port in cached table. otherwise create tunnel and
 * retrieve port number from OVSDB.
 */
static uint32_t find_gtp_port_no(struct in_addr enb_addr) {
  if (!spgw_config.sgw_config.ovs_config.multi_tunnel) {
    return 0;
  }
  if ((uint32_t) enb_addr.s_addr == 0) {
    OAILOG_WARNING(LOG_GTPV1U, "zero enb IP address not supported");
    return 0;
  }
  char port_name[MAX_GTP_PORT_NAME_LENGTH];
  ip_addr_to_gtp_port_name(enb_addr, port_name);

  uint32_t portno = search_portno_records(port_name);
  if (portno) {
    return portno;
  }

  portno = create_gtp_port(enb_addr, port_name);
  add_portno_rec(port_name, portno);
  return portno;
}

/**
 * Initialize GTP port table for caching GTP tunnel port numbers.
 */
static void openflow_multi_tunnel_init(void) {
  char* probe_gtp_type = "sudo ovs-vsctl list Open_vSwitch | grep gtpu";

  // OVS GTP tunnel type has changed upstream, for better compatibility
  // detect it on initilization.
  int rc = system(probe_gtp_type);
  if (rc != 0) {
    ovs_gtp_type = strdup("gtp");
  } else {
    ovs_gtp_type = strdup("gtpu");
  }
  OAILOG_INFO(LOG_GTPV1U, "Using GTP type: %s", ovs_gtp_type);

  gtp_portno_rec.arr = calloc(INIT_GTP_TABLE_SIZE, sizeof(struct gtp_portno));
  assert(gtp_portno_rec.arr != NULL);
  gtp_portno_rec.size = INIT_GTP_TABLE_SIZE;
}

// tunnel flows
int openflow_uninit(void) {
  int ret;
  if ((ret = stop_of_controller()) < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "Could not stop openflow controller on uninit\n");
  }
  return ret;
}

int openflow_init(
    struct in_addr* ue_net, uint32_t mask, int mtu, int* fd0, int* fd1u,
    bool persist_state) {
  AssertFatal(
      start_of_controller(persist_state) >= 0,
      "Could not start openflow controller\n");
  return 0;
}

int openflow_reset(void) {
  int rv = 0;
  return rv;
}

int openflow_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl, char* apn) {
  uint32_t gtp_portno = find_gtp_port_no(enb);

  return openflow_controller_add_gtp_tunnel(
      ue, ue_ipv6, vlan, enb, i_tei, o_tei, (const char*) imsi.digit, flow_dl,
      flow_precedence_dl, gtp_portno);
}

int openflow_del_tunnel(
    struct in_addr enb, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t o_tei, struct ip_flow_dl* flow_dl) {
  uint32_t gtp_portno = find_gtp_port_no(enb);

  return openflow_controller_del_gtp_tunnel(
      ue, ue_ipv6, i_tei, flow_dl, gtp_portno);
}

/* S8 tunnel related APIs */
int openflow_add_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in_addr pgw, uint32_t i_tei, uint32_t o_tei, uint32_t pgw_i_tei,
    uint32_t pgw_o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl) {
  uint32_t enb_portno = find_gtp_port_no(enb);
  uint32_t pgw_portno = find_gtp_port_no(pgw);

  return openflow_controller_add_gtp_s8_tunnel(
      ue, ue_ipv6, vlan, enb, pgw, i_tei, o_tei, pgw_i_tei, pgw_o_tei,
      (const char*) imsi.digit, flow_dl, flow_precedence_dl, enb_portno,
      pgw_portno);
}

int openflow_del_s8_tunnel(
    struct in_addr enb, struct in_addr pgw, struct in_addr ue,
    struct in6_addr* ue_ipv6, uint32_t i_tei, uint32_t o_tei,
    struct ip_flow_dl* flow_dl) {
  uint32_t enb_portno = find_gtp_port_no(enb);
  uint32_t pgw_portno = find_gtp_port_no(pgw);

  return openflow_controller_del_gtp_s8_tunnel(
      ue, ue_ipv6, i_tei, flow_dl, enb_portno, pgw_portno);
}

int openflow_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl) {
  return openflow_controller_discard_data_on_tunnel(
      ue, ue_ipv6, i_tei, flow_dl);
}

int openflow_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl) {
  return openflow_controller_forward_data_on_tunnel(
      ue, ue_ipv6, i_tei, flow_dl, flow_precedence_dl);
}

int openflow_add_paging_rule(struct in_addr ue) {
  return openflow_controller_add_paging_rule(ue);
}

int openflow_delete_paging_rule(struct in_addr ue) {
  return openflow_controller_delete_paging_rule(ue);
}

/**
 * Send packet marker to enodeB @enb for tunnel @tei.
 */
int openflow_send_end_marker(struct in_addr enb, uint32_t tei) {
  static bool end_marker_supported = true;
  char end_marker_cmd[2048];
  int rc;

  // End marker needs OVS patch from magma repo, check if it has
  // worked on this host before trying the cmd.
  if (!end_marker_supported) {
    return -ENODEV;
  }

  if (tei == 0 || (uint32_t) enb.s_addr == 0) {
    // No need to send end marker for tunnel with zero tunnel metadata.
    return 0;
  }
  // use a ethernet packet just to make packet out happy.
  rc = snprintf(
      end_marker_cmd, sizeof(end_marker_cmd),
      "sudo ovs-ofctl packet-out gtp_br0 "
      "'in_port=local packet=50540000000a5054000000008000,"
      "actions=load:%" PRIu32
      "->tun_id[0..31],"
      "set_field:%s->tun_dst,"
      "set_field:0xfe->tun_gtpu_msgtype,set_field:0x30->tun_gtpu_flags,output:"
      "gtp0'",
      tei, inet_ntoa(enb));

  if (rc < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "end marker cmd: format error %d", rc);
    return rc;
  }
  rc = system(end_marker_cmd);
  if (rc != 0) {
    OAILOG_ERROR(
        LOG_GTPV1U, "end marker cmd: [%s] failed: %d", end_marker_cmd, rc);
    end_marker_supported = false;
  } else {
    OAILOG_DEBUG(
        LOG_GTPV1U, "End marker sent: tei %" PRIu32 " tun_dst %s", tei,
        inet_ntoa(enb));
  }
  return rc;
}

const char* openflow_get_dev_name(void) {
  return bdata(spgw_config.sgw_config.ovs_config.bridge_name);
}

static const struct gtp_tunnel_ops openflow_ops = {
    .init                   = openflow_init,
    .uninit                 = openflow_uninit,
    .reset                  = openflow_reset,
    .add_tunnel             = openflow_add_tunnel,
    .del_tunnel             = openflow_del_tunnel,
    .add_s8_tunnel          = openflow_add_s8_tunnel,
    .del_s8_tunnel          = openflow_del_s8_tunnel,
    .discard_data_on_tunnel = openflow_discard_data_on_tunnel,
    .forward_data_on_tunnel = openflow_forward_data_on_tunnel,
    .add_paging_rule        = openflow_add_paging_rule,
    .delete_paging_rule     = openflow_delete_paging_rule,
    .send_end_marker        = openflow_send_end_marker,
    .get_dev_name           = openflow_get_dev_name,
};

const struct gtp_tunnel_ops* gtp_tunnel_ops_init_openflow(void) {
  if (spgw_config.sgw_config.ovs_config.multi_tunnel) {
    openflow_multi_tunnel_init();
  }

  return &openflow_ops;
}
