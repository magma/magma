/**
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
#if 0
#include <string>

#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"

extern "C" {
#include "amf_service_handler.h"
#include "log.h"
}
#include "AmfServiceImpl.h"
#endif

#include <string>

extern "C" {
#include "amf_service_handler.h"
#include "log.h"
}
#include "AmfServiceImpl.h"
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

AmfServiceImpl::AmfServiceImpl() {}

Status AmfServiceImpl::SetAmfNotification(
    ServerContext* context, const SetSmNotificationContext* notif,
    SmContextVoid* response) {
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSmNotificationContext request\n");
  // ToDo processing ITTI,ZMQ

  return Status::OK;
}
// Set message from SessionD received
Status AmfServiceImpl::SetSmfSessionContext(
    ServerContext* context, const SetSMSessionContextAccess* request,
    SmContextVoid* response) {
  OAILOG_INFO(LOG_UTIL, "Received  GRPC SetSMSessionContextAccess request\n");
  // ToDo processing ITTI,ZMQ

  itti_n11_create_pdu_session_response_t itti_msg;
  auto& req_common = request->common_context();
  auto& req_m5g    = request->rat_specific_context().m5g_session_context_rsp();

  // CommonSessionContext
  strcpy(itti_msg.imsi, req_common.sid().id().c_str());
  itti_msg.sm_session_fsm_state =
      (SMSessionFSMState_response) req_common.sm_session_state();
  itti_msg.sm_session_version = req_common.sm_session_version();

  // RatSpecificContextAccess
  // strcpy((char*)itti_msg.pdu_session_id, req_m5g.pdu_session_id().c_str());
  // itti_msg.pdu_session_id =  req_m5g.pdu_session_id(); cobraranu
  // itti_msg.pdu_session_id =  (uint8_t)1;
  // itti_msg.pdu_session_id = req_m5g.pdu_session_id();
  itti_msg.pdu_session_type =
      (PduSessionType_response) req_m5g.pdu_session_type();

  // pdu_change
  itti_msg.selected_ssc_mode = (SscMode_response) req_m5g.selected_ssc_mode();
  itti_msg.M5gsm_cause       = (M5GSMCause_response) req_m5g.m5gsm_cause();
  for (int i = 0, n = req_m5g.authorized_qos_rules_size(); (i < n) && (i < 4);
       i++) {  // TODO 32 in NAS5g,3 in pcap, 4 in zmq struct, revisit later
    itti_msg.authorized_qos_rules[i].qos_rule_identifier =
        (uint32_t) req_m5g.authorized_qos_rules(i).qos_rule_identifier();
    itti_msg.authorized_qos_rules[i].dqr =
        req_m5g.authorized_qos_rules(i).dqr();
    itti_msg.authorized_qos_rules[i].number_of_packet_filters =
        (uint32_t) req_m5g.authorized_qos_rules(i).number_of_packet_filters();
    for (int j = 0;
         (j < req_m5g.authorized_qos_rules(i).packet_filter_identifier_size() ||
          j < 16);
         j++) {
      itti_msg.authorized_qos_rules[i].packet_filter_identifier[j] =
          (uint32_t) req_m5g.authorized_qos_rules(i).packet_filter_identifier(
              j);
    }
    itti_msg.authorized_qos_rules[i].qos_rule_precedence =
        (uint32_t) req_m5g.authorized_qos_rules(i).qos_rule_precedence();
    itti_msg.authorized_qos_rules[i].segregation =
        req_m5g.authorized_qos_rules(i).segregation();
    itti_msg.authorized_qos_rules[i].qos_flow_identifier =
        (uint32_t) req_m5g.authorized_qos_rules(i).qos_flow_identifier();
  }
  itti_msg.session_ambr.uplink_unit_type =
      (AmbrUnit_response) req_m5g.session_ambr().uplink_unit_type();
  itti_msg.session_ambr.uplink_units =
      (uint32_t) req_m5g.session_ambr().uplink_units();
  itti_msg.session_ambr.downlink_unit_type =
      (AmbrUnit_response) req_m5g.session_ambr().downlink_unit_type();
  itti_msg.session_ambr.downlink_units =
      (uint32_t) req_m5g.session_ambr().downlink_units();

  itti_msg.qos_list.qos_flow_req_item.qos_flow_identifier =
      (uint32_t) req_m5g.qos_list().flow().qos_flow_ident();
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .qos_characteristic.non_dynamic_5QI_desc.fiveQI =
      req_m5g.qos_list().flow().param().qos_chars().fiveqi();
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.priority_level =
      req_m5g.qos_list().flow().param().alloc_reten_prio().prio_level();
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_cap =
      (pre_emption_capability) req_m5g.qos_list()
          .flow()
          .param()
          .alloc_reten_prio()
          .pre_emtion_cap();
  itti_msg.qos_list.qos_flow_req_item.qos_flow_level_qos_param
      .alloc_reten_priority.pre_emption_vul =
      (pre_emption_vulnerability) req_m5g.qos_list()
          .flow()
          .param()
          .alloc_reten_prio()
          .pre_emtion_vul();

  std::string byte4_value;
  // Lets get the TEID 4 byte value from grpc message
  byte4_value.assign(req_m5g.upf_endpoint().teid());
  int len = byte4_value.length();
  for (int i = 0; i < len; i++) {
    itti_msg.upf_endpoint.teid[i] = byte4_value[i];
  }  // No need of null termination

  // Lets get the ip address  byte value from grpc message
  byte4_value.assign(req_m5g.upf_endpoint().end_ipv4_addr());
  len = byte4_value.length();
  for (int i = 0; i < len; i++) {
    itti_msg.upf_endpoint.end_ipv4_addr[i] = byte4_value[i];
  }
  OAILOG_INFO(
      LOG_AMF_APP, "#######TIED: %02x %02x %02x %02x \n",
      itti_msg.upf_endpoint.teid[0], itti_msg.upf_endpoint.teid[1],
      itti_msg.upf_endpoint.teid[2], itti_msg.upf_endpoint.teid[3]);

  OAILOG_INFO(
      LOG_AMF_APP, "#######IP: %02x %02x %02x %02x \n",
      itti_msg.upf_endpoint.end_ipv4_addr[0],
      itti_msg.upf_endpoint.end_ipv4_addr[1],
      itti_msg.upf_endpoint.end_ipv4_addr[2],
      itti_msg.upf_endpoint.end_ipv4_addr[3]);
#if 0
  //strcpy((char*)itti_msg.upf_endpoint.teid, req_m5g.upf_endpoint().teid().c_str());
  memcpy(itti_msg.upf_endpoint.teid, req_m5g.upf_endpoint().teid(), 4);
  //strcpy(
  //    (char*)itti_msg.upf_endpoint.end_ipv4_addr,
  //    req_m5g.upf_endpoint().end_ipv4_addr().c_str());
  memcpy(
      itti_msg.upf_endpoint.end_ipv4_addr,
      req_m5g.upf_endpoint().end_ipv4_addr(), 4);
#endif
  strcpy(
      (char*) itti_msg.procedure_trans_identity,
      req_m5g.procedure_trans_identity().c_str());  // pdu_change

  itti_msg.always_on_pdu_session_indication =
      req_m5g.always_on_pdu_session_indication();
  itti_msg.allowed_ssc_mode = (SscMode_response) req_m5g.allowed_ssc_mode();
  itti_msg.M5gsm_congetion_re_attempt_indicator =
      req_m5g.m5gsm_congetion_re_attempt_indicator();
  itti_msg.pdu_address.redirect_address_type =
      (RedirectAddressType_response) req_m5g.pdu_address()
          .redirect_address_type();
  strcpy(
      (char*) itti_msg.pdu_address.redirect_server_address,
      req_m5g.pdu_address().redirect_server_address().c_str());
  send_n11_create_pdu_session_resp_itti(&itti_msg);
  return Status::OK;
}

}  // namespace magma
