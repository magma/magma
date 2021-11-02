/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <string>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/include/amf_service_handler.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/tasks/grpc_service/AmfServiceImpl.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"

namespace grpc {

class ServerContext;
}  // namespace grpc

using grpc::ServerContext;
using grpc::Status;
using magma::lte::SetSmNotificationContext;
using magma::lte::SetSMSessionContextAccess;
using magma::lte::SmContextVoid;
using magma::lte::SmfPduSessionSmContext;

namespace magma {
using namespace lte;
#define TEID_SIZE 4
#define UPF_IPV4_ADDR_SIZE 4

AmfServiceImpl::AmfServiceImpl() {}

Status AmfServiceImpl::SetAmfNotification(
    ServerContext* context, const SetSmNotificationContext* notif,
    SmContextVoid* response) {
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSmNotificationContext request\n");
  // ToDo processing ITTI,ZMQ

  itti_n11_received_notification_t itti_msg;
  auto& notify_common = notif->common_context();
  auto& req_m5g       = notif->rat_specific_notification();

  // CommonSessionContext
  strcpy(itti_msg.imsi, notify_common.sid().id().c_str());
  itti_msg.sm_session_fsm_state =
      (SMSessionFSMState_response) notify_common.sm_session_state();
  itti_msg.sm_session_version = notify_common.sm_session_version();

  // RatSpecificContextAccess
  itti_msg.pdu_session_id   = req_m5g.pdu_session_id();
  itti_msg.request_type     = (RequestType_received) req_m5g.request_type();
  itti_msg.pdu_session_type = (pdu_session_type_t) req_m5g.pdu_session_type();
  itti_msg.m5g_sm_capability.reflective_qos =
      req_m5g.m5g_sm_capability().reflective_qos();
  itti_msg.m5g_sm_capability.multi_homed_ipv6_pdu_session =
      req_m5g.m5g_sm_capability().multi_homed_ipv6_pdu_session();
  itti_msg.m5gsm_cause = (m5g_sm_cause_t) req_m5g.m5gsm_cause();

  // pdu_change
  itti_msg.notify_ue_evnt = (notify_ue_event) req_m5g.notify_ue_event();

  send_n11_notification_received_itti(&itti_msg);
  return Status::OK;
}
// Set message from SessionD received
Status AmfServiceImpl::SetSmfSessionContext(
    ServerContext* context, const SetSMSessionContextAccess* request,
    SmContextVoid* response) {
  struct in_addr ip_addr       = {0};
  char ip_str[INET_ADDRSTRLEN] = {0};
  uint32_t ip_int              = 0;
  OAILOG_INFO(
      LOG_UTIL, "Received GRPC SetSmfSessionContext request from SMF\n");

  itti_n11_create_pdu_session_response_t itti_msg;
  auto& req_common = request->common_context();
  auto& req_m5g    = request->rat_specific_context().m5g_session_context_rsp();

  // CommonSessionContext
  strcpy(itti_msg.imsi, req_common.sid().id().c_str());
  itti_msg.sm_session_fsm_state =
      (sm_session_fsm_state_t) req_common.sm_session_state();
  itti_msg.sm_session_version = req_common.sm_session_version();

  // RatSpecificContextAccess
  itti_msg.pdu_session_id    = req_m5g.pdu_session_id();
  itti_msg.pdu_session_type  = (pdu_session_type_t) req_m5g.pdu_session_type();
  itti_msg.selected_ssc_mode = (ssc_mode_t) req_m5g.selected_ssc_mode();
  itti_msg.m5gsm_cause       = (m5g_sm_cause_t) req_m5g.m5gsm_cause();

  itti_msg.session_ambr.uplink_unit_type = req_m5g.session_ambr().br_unit();
  itti_msg.session_ambr.uplink_units =
      (uint32_t) req_m5g.session_ambr().max_bandwidth_ul();

  itti_msg.session_ambr.downlink_unit_type = req_m5g.session_ambr().br_unit();
  itti_msg.session_ambr.downlink_units =
      (uint32_t) req_m5g.session_ambr().max_bandwidth_dl();

  itti_msg.qos_list.qos_flow_req_item.qos_flow_identifier = req_m5g.qos().qci();

  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .qos_characteristic.non_dynamic_5QI_desc.fiveQI =
      req_m5g.qos().qci();  // enum
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.priority_level =
      req_m5g.qos().arp().priority_level();  // uint32
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_cap =
      (pre_emption_capability) req_m5g.qos().arp().pre_capability();  // enum
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_vul =
      (pre_emption_vulnerability) req_m5g.qos()
          .arp()
          .pre_vulnerability();  // enum

  // get the 4 byte UPF TEID and UPF IP message
  uint32_t nteid                = req_m5g.upf_endpoint().teid();
  itti_msg.upf_endpoint.teid[0] = (nteid >> 24) & 0xFF;
  itti_msg.upf_endpoint.teid[1] = (nteid >> 16) & 0xFF;
  itti_msg.upf_endpoint.teid[2] = (nteid >> 8) & 0xFF;
  itti_msg.upf_endpoint.teid[3] = nteid & 0xFF;

  if (req_m5g.upf_endpoint().end_ipv4_addr().size() > 0) {
    inet_pton(
        AF_INET, req_m5g.upf_endpoint().end_ipv4_addr().c_str(),
        itti_msg.upf_endpoint.end_ipv4_addr);
  }

  strcpy(
      (char*) itti_msg.procedure_trans_identity,
      req_m5g.procedure_trans_identity().c_str());  // pdu_change
  itti_msg.always_on_pdu_session_indication =
      req_m5g.always_on_pdu_session_indication();
  itti_msg.allowed_ssc_mode = (ssc_mode_t) req_m5g.allowed_ssc_mode();
  itti_msg.m5gsm_congetion_re_attempt_indicator =
      req_m5g.m5g_sm_congestion_reattempt_indicator();
  itti_msg.pdu_address.redirect_address_type =
      (redirect_address_type_t) req_m5g.pdu_address().redirect_address_type();
  // PDU IP address coming from SMF in human-readable format has to be packed
  // into 4 raw bytes in hex for NAS5G layer
  strcpy(ip_str, req_m5g.pdu_address().redirect_server_address().c_str());
  inet_pton(AF_INET, ip_str, &(ip_addr.s_addr));
  ip_int = ntohl(ip_addr.s_addr);
  INT32_TO_BUFFER(ip_int, itti_msg.pdu_address.redirect_server_address);
  send_n11_create_pdu_session_resp_itti(&itti_msg);
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSMSessionContextAccess request \n");
  return Status::OK;
}

}  // namespace magma
