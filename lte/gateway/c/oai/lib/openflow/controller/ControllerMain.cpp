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

#include "OpenflowController.h"
#include "PagingApplication.h"
#include "BaseApplication.h"
#include "ControllerMain.h"
#include "GTPApplication.h"
extern "C" {
#include "log.h"
#include "spgw_config.h"
#include "common_defs.h"
}

static const int OFP_LOCAL   = 65534;
static const int OF13P_LOCAL = 0xfffffffe;

namespace {
openflow::OpenflowController ctrl(
    CONTROLLER_ADDR, CONTROLLER_PORT, NUM_WORKERS, false);
}

int start_of_controller(bool persist_state) {
  static openflow::PagingApplication paging_app;
  static openflow::BaseApplication base_app(persist_state);
  int uplink_port_num_ = OF13P_LOCAL;  // default is LOCAL port.

  if (spgw_config.pgw_config.enable_nat == false) {
    // For Non NAT config we need to know uplink bridge port.
    uplink_port_num_ = spgw_config.sgw_config.ovs_config.uplink_port_num;
  }
  if (uplink_port_num_ == OFP_LOCAL) {  // convert it to OF 1.3 LOCAL port no.
    uplink_port_num_ = OF13P_LOCAL;
  }

  static openflow::GTPApplication gtp_app(
      std::string(bdata(spgw_config.sgw_config.ovs_config.uplink_mac)),
      spgw_config.sgw_config.ovs_config.gtp_port_num,
      spgw_config.sgw_config.ovs_config.mtr_port_num,
      spgw_config.sgw_config.ovs_config.internal_sampling_port_num,
      spgw_config.sgw_config.ovs_config.internal_sampling_fwd_tbl_num,
      uplink_port_num_);
  // Base app registers first, because it deletes/creates default flow
  ctrl.register_for_event(&base_app, openflow::EVENT_SWITCH_UP);
  ctrl.register_for_event(&base_app, openflow::EVENT_ERROR);
  ctrl.register_for_event(&paging_app, openflow::EVENT_PACKET_IN);
  ctrl.register_for_event(&paging_app, openflow::EVENT_SWITCH_UP);
  ctrl.register_for_event(&paging_app, openflow::EVENT_ADD_PAGING_RULE);
  ctrl.register_for_event(&paging_app, openflow::EVENT_DELETE_PAGING_RULE);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_SWITCH_UP);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_ADD_GTP_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_DELETE_GTP_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_ADD_GTP_S8_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_DELETE_GTP_S8_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_DISCARD_DATA_ON_GTP_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_FORWARD_DATA_ON_GTP_TUNNEL);
  ctrl.start();
  OAILOG_INFO(LOG_GTPV1U, "Started openflow controller\n");
#define CONNECTION_WAIT_TIME 300
  OAILOG_FUNC_RETURN(
      LOG_GTPV1U, ctrl.is_controller_connected_to_switch(CONNECTION_WAIT_TIME));
}

int stop_of_controller(void) {
  ctrl.stop();
  OAILOG_INFO(LOG_GTPV1U, "Stopped openflow controller\n");
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

/**
 * This callback is called from the event loop itself to dispatch an external
 * event to all registered applications
 */
static void* external_event_callback(std::shared_ptr<void> data) {
  auto external_event = std::static_pointer_cast<openflow::ExternalEvent>(data);
  ctrl.dispatch_event(*external_event);
  return NULL;
}

int openflow_controller_add_gtp_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, const char* imsi,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl,
    uint32_t gtp_portno) {
  if (flow_dl) {
    auto add_tunnel = std::make_shared<openflow::AddGTPTunnelEvent>(
        ue, ue_ipv6, vlan, enb, i_tei, o_tei, imsi, flow_dl, flow_precedence_dl,
        gtp_portno);
    ctrl.inject_external_event(add_tunnel, external_event_callback);
  } else {
    auto add_tunnel = std::make_shared<openflow::AddGTPTunnelEvent>(
        ue, ue_ipv6, vlan, enb, i_tei, o_tei, imsi, gtp_portno);
    ctrl.inject_external_event(add_tunnel, external_event_callback);
  }
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_del_gtp_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t gtp_portno) {
  if (flow_dl) {
    auto del_tunnel = std::make_shared<openflow::DeleteGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, flow_dl, gtp_portno);
    ctrl.inject_external_event(del_tunnel, external_event_callback);
  } else {
    auto del_tunnel = std::make_shared<openflow::DeleteGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, gtp_portno);
    ctrl.inject_external_event(del_tunnel, external_event_callback);
  }
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_add_gtp_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in_addr pgw, uint32_t i_tei, uint32_t o_tei, uint32_t pgw_i_tei,
    uint32_t pgw_o_tei, const char* imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl, uint32_t enb_gtp_port, uint32_t pgw_gtp_port) {
  if (flow_dl) {
    auto add_tunnel = std::make_shared<openflow::AddGTPTunnelEvent>(
        ue, ue_ipv6, vlan, enb, pgw, i_tei, o_tei, pgw_i_tei, pgw_o_tei, imsi,
        flow_dl, flow_precedence_dl, enb_gtp_port, pgw_gtp_port);
    ctrl.inject_external_event(add_tunnel, external_event_callback);
  } else {
    auto add_tunnel = std::make_shared<openflow::AddGTPTunnelEvent>(
        ue, ue_ipv6, vlan, enb, pgw, i_tei, o_tei, pgw_i_tei, pgw_o_tei, imsi,
        enb_gtp_port, pgw_gtp_port);
    ctrl.inject_external_event(add_tunnel, external_event_callback);
  }
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_del_gtp_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t enb_gtp_port, uint32_t pgw_gtp_port) {
  if (flow_dl) {
    auto del_tunnel = std::make_shared<openflow::DeleteGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, flow_dl, enb_gtp_port, pgw_gtp_port);
    ctrl.inject_external_event(del_tunnel, external_event_callback);
  } else {
    auto del_tunnel = std::make_shared<openflow::DeleteGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, enb_gtp_port, pgw_gtp_port);
    ctrl.inject_external_event(del_tunnel, external_event_callback);
  }
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl) {
  if (flow_dl) {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, openflow::EVENT_DISCARD_DATA_ON_GTP_TUNNEL, flow_dl,
        false);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  } else {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, openflow::EVENT_DISCARD_DATA_ON_GTP_TUNNEL);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  }
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl) {
  if (flow_dl) {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, openflow::EVENT_FORWARD_DATA_ON_GTP_TUNNEL, flow_dl,
        flow_precedence_dl);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  } else {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
        ue, ue_ipv6, i_tei, openflow::EVENT_FORWARD_DATA_ON_GTP_TUNNEL);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  }
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_add_paging_rule(struct in_addr ue_ip) {
  auto paging_event = std::make_shared<openflow::AddPagingRuleEvent>(ue_ip);
  ctrl.inject_external_event(paging_event, external_event_callback);
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}

int openflow_controller_delete_paging_rule(struct in_addr ue_ip) {
  auto paging_event = std::make_shared<openflow::DeletePagingRuleEvent>(ue_ip);
  ctrl.inject_external_event(paging_event, external_event_callback);
  OAILOG_FUNC_RETURN(LOG_GTPV1U, RETURNok);
}
