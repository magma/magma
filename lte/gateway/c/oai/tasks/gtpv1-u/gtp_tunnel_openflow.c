/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

#include <errno.h>
#include <stdint.h>
#include <netinet/in.h>
#include <stdlib.h>

#include "assertions.h"
#include "log.h"
#include "gtpv1u.h"
#include "ControllerMain.h"
#include "3gpp_23.003.h"

extern struct gtp_tunnel_ops gtp_tunnel_ops;

int openflow_uninit(void)
{
  int ret;
  if ((ret = stop_of_controller()) < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "Could not stop openflow controller on uninit\n");
  }
  return ret;
}

int openflow_init(
  struct in_addr *ue_net,
  uint32_t mask,
  int mtu,
  int *fd0,
  int *fd1u,
  bool persist_state)
{
  AssertFatal(
    start_of_controller(persist_state) >= 0,
    "Could not start openflow controller\n");
  return 0;
}

int openflow_reset(void)
{
  int rv = 0;
  return rv;
}

int openflow_add_tunnel(
  struct in_addr ue,
  struct in_addr enb,
  uint32_t i_tei,
  uint32_t o_tei,
  Imsi_t imsi,
  struct ipv4flow_dl* flow_dl,
  uint32_t flow_precedence_dl)
{
  return openflow_controller_add_gtp_tunnel(
    ue,
    enb,
    i_tei,
    o_tei,
    (const char*) imsi.digit,
    flow_dl,
    flow_precedence_dl);
}

int openflow_del_tunnel(struct in_addr ue, uint32_t i_tei,
    uint32_t o_tei, struct ipv4flow_dl *flow_dl)
{
  return openflow_controller_del_gtp_tunnel(ue, i_tei, flow_dl);
}

int openflow_discard_data_on_tunnel(struct in_addr ue, uint32_t i_tei,
    struct ipv4flow_dl *flow_dl)
{
  return openflow_controller_discard_data_on_tunnel(ue, i_tei, flow_dl);
}

int openflow_forward_data_on_tunnel(
  struct in_addr ue,
  uint32_t i_tei,
  struct ipv4flow_dl* flow_dl,
  uint32_t flow_precedence_dl)
{
  return openflow_controller_forward_data_on_tunnel(
    ue, i_tei, flow_dl, flow_precedence_dl);
}

int openflow_add_paging_rule(struct in_addr ue)
{
  return openflow_controller_add_paging_rule(ue);
}

int openflow_delete_paging_rule(struct in_addr ue)
{
  return openflow_controller_delete_paging_rule(ue);
}

/**
 * Send packet marker to enodeB @enb for tunnel @tei.
 */
int openflow_send_end_marker(struct in_addr enb, uint32_t tei)
{
  static bool end_marker_supported = true;
  char end_marker_cmd[2048];
  int rc;

  // End marker needs OVS patch from magma repo, check if it has
  // worked on this host before trying the cmd.
  if (!end_marker_supported) {
    return -ENODEV;
  }

  if (tei == 0) {
    // No need to send end marker for tunnel with tei zero.
    return 0;
  }
  // use a ethernet packet just to make packet out happy.
  rc = snprintf(end_marker_cmd, sizeof(end_marker_cmd),
  "sudo ovs-ofctl packet-out gtp_br0 "
  "'in_port=local packet=50540000000a5054000000008000,"
  "actions=load:%"PRIu32"->tun_id[0..31],"
  "set_field:%s->tun_dst,"
  "set_field:0xfe->tun_gtp_msg_type,set_field:0x30->tun_gtp_flags,output:gtp0'",
  tei, inet_ntoa(enb));

  if (rc < 0) {
    OAILOG_ERROR(LOG_GTPV1U, "end marker cmd: format error %d", rc);
    return rc;
  }
  rc = system(end_marker_cmd);
  if (rc != 0) {
    OAILOG_ERROR(LOG_GTPV1U, "end marker cmd: [%s] failed: %d", end_marker_cmd, rc);
    end_marker_supported = false;
  }
  return rc;
}

static const struct gtp_tunnel_ops openflow_ops = {
  .init = openflow_init,
  .uninit = openflow_uninit,
  .reset = openflow_reset,
  .add_tunnel = openflow_add_tunnel,
  .del_tunnel = openflow_del_tunnel,
  .discard_data_on_tunnel = openflow_discard_data_on_tunnel,
  .forward_data_on_tunnel = openflow_forward_data_on_tunnel,
  .add_paging_rule = openflow_add_paging_rule,
  .delete_paging_rule = openflow_delete_paging_rule,
  .send_end_marker = openflow_send_end_marker,
};

const struct gtp_tunnel_ops *gtp_tunnel_ops_init_openflow(void)
{
  return &openflow_ops;
}
