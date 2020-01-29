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

#include "OpenflowController.h"
#include "PagingApplication.h"
#include "BaseApplication.h"
#include "ControllerMain.h"
#include "GTPApplication.h"
extern "C" {
#include "log.h"
#include "spgw_config.h"
}

namespace {
openflow::OpenflowController
  ctrl(CONTROLLER_ADDR, CONTROLLER_PORT, NUM_WORKERS, false);
}

int start_of_controller(bool persist_state)
{
  static openflow::PagingApplication paging_app;
  static openflow::BaseApplication base_app(persist_state);
  static openflow::GTPApplication gtp_app(
    std::string(bdata(spgw_config.sgw_config.ovs_config.uplink_mac)),
    spgw_config.sgw_config.ovs_config.gtp_port_num);
  // Base app registers first, because it deletes/creates default flow
  ctrl.register_for_event(&base_app, openflow::EVENT_SWITCH_UP);
  ctrl.register_for_event(&base_app, openflow::EVENT_ERROR);
  ctrl.register_for_event(&paging_app, openflow::EVENT_PACKET_IN);
  ctrl.register_for_event(&paging_app, openflow::EVENT_SWITCH_UP);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_ADD_GTP_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_DELETE_GTP_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_DISCARD_DATA_ON_GTP_TUNNEL);
  ctrl.register_for_event(&gtp_app, openflow::EVENT_FORWARD_DATA_ON_GTP_TUNNEL);
  ctrl.start();
  OAILOG_INFO(LOG_GTPV1U, "Started openflow controller\n");
  return 0;
}

int stop_of_controller(void)
{
  ctrl.stop();
  OAILOG_INFO(LOG_GTPV1U, "Stopped openflow controller\n");
  return 0;
}

/**
 * This callback is called from the event loop itself to dispatch an external
 * event to all registered applications
 */
static void *external_event_callback(std::shared_ptr<void> data)
{
  auto external_event = std::static_pointer_cast<openflow::ExternalEvent>(data);
  ctrl.dispatch_event(*external_event);
}

int openflow_controller_add_gtp_tunnel(
  struct in_addr ue,
  struct in_addr enb,
  uint32_t i_tei,
  uint32_t o_tei,
  const char *imsi,
  struct ipv4flow_dl *flow_dl,
  uint32_t flow_precedence_dl)
{
  if (flow_dl) {
    auto add_tunnel = std::make_shared<openflow::AddGTPTunnelEvent>(
      ue, enb, i_tei, o_tei, imsi, flow_dl, flow_precedence_dl);
    ctrl.inject_external_event(add_tunnel, external_event_callback);
  } else {
    auto add_tunnel = std::make_shared<openflow::AddGTPTunnelEvent>(
      ue, enb, i_tei, o_tei, imsi);
    ctrl.inject_external_event(add_tunnel, external_event_callback);
  }
  return 0;
}

int openflow_controller_del_gtp_tunnel(struct in_addr ue, uint32_t i_tei,
    struct ipv4flow_dl *flow_dl)
{
  if (flow_dl) {
    auto del_tunnel =
      std::make_shared<openflow::DeleteGTPTunnelEvent>(ue, i_tei, flow_dl);
    ctrl.inject_external_event(del_tunnel, external_event_callback);
  } else {
    auto del_tunnel =
      std::make_shared<openflow::DeleteGTPTunnelEvent>(ue, i_tei);
    ctrl.inject_external_event(del_tunnel, external_event_callback);
  }
  return 0;
}

int openflow_controller_discard_data_on_tunnel(
  struct in_addr ue,
  uint32_t i_tei,
  struct ipv4flow_dl *flow_dl)
{
  if (flow_dl) {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
        ue, i_tei, openflow::EVENT_DISCARD_DATA_ON_GTP_TUNNEL, flow_dl, false);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  } else {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
        ue, i_tei, openflow::EVENT_DISCARD_DATA_ON_GTP_TUNNEL);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  }
  return 0;
}

int openflow_controller_forward_data_on_tunnel(
  struct in_addr ue,
  uint32_t i_tei,
  struct ipv4flow_dl *flow_dl,
  uint32_t flow_precedence_dl)
{
  if (flow_dl) {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
      ue,
      i_tei,
      openflow::EVENT_FORWARD_DATA_ON_GTP_TUNNEL,
      flow_dl,
      flow_precedence_dl);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  } else {
    auto gtp_tunnel = std::make_shared<openflow::HandleDataOnGTPTunnelEvent>(
      ue, i_tei, openflow::EVENT_FORWARD_DATA_ON_GTP_TUNNEL);
    ctrl.inject_external_event(gtp_tunnel, external_event_callback);
  }
  return 0;
}
