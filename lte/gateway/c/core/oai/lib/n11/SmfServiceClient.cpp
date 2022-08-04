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
#include <google/protobuf/util/time_util.h>
#include <cassert>
#include <grpcpp/impl/codegen/client_context.h>
#include <grpcpp/impl/codegen/status.h>
#include <cstring>
#include <iostream>
#include <memory>
#include <string>
#include <thread>
#include <lte/protos/session_manager.grpc.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <arpa/inet.h>
#include <utility>

#include "orc8r/gateway/c/common/service_registry/ServiceRegistrySingleton.hpp"
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/include/amf_service_handler.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

using grpc::Status;
using magma::AsyncLocalResponse;
using magma::ServiceRegistrySingleton;

void handle_session_context_response(grpc::Status status,
                                     magma::lte::SmContextVoid response) {
  if (!status.ok()) {
    OAILOG_ERROR(LOG_AMF_APP,
                 "AsyncSetAmfSessionContext fails with"
                 "error code [%d] error message [%s] ",
                 status.error_code(), status.error_message().c_str());
  }
}

int send_n11_create_pdu_session_failure_itti(
    itti_n11_create_pdu_session_failure_t* itti_msg) {
  OAILOG_DEBUG(LOG_UTIL,
               "Sending itti_n11_create_pdu_session_failure to AMF \n");
  MessageDef* message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, N11_CREATE_PDU_SESSION_FAILURE);
  if (message_p == NULL) {
    OAILOG_ERROR(
        LOG_UTIL,
        "Failed to allocate memory for N11_CREATE_PDU_SESSION_FAILURE\n");
    return RETURNerror;
  }
  message_p->ittiMsg.n11_create_pdu_session_failure = *itti_msg;
  return send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

void handle_session_context_response(grpc::Status status,
                                     const SetSMSessionContext& request,
                                     magma::lte::SmContextVoid response) {
  if (status.ok()) {
    return;
  }

  OAILOG_ERROR(LOG_AMF_APP,
               "AsyncSetAmfSessionContext fails with"
               "error code [%d] error message [%s] ",
               status.error_code(), status.error_message().c_str());

  auto& req_common = request.common_context();
  itti_n11_create_pdu_session_failure_t itti_msg = {};

  if (!req_common.sid().id().empty()) {
    strncpy(itti_msg.imsi, req_common.sid().id().c_str() + 4,
            IMSI_BCD_DIGITS_MAX);
  }
  auto& req_rat_specific =
      request.rat_specific_context().m5gsm_session_context();
  itti_msg.pdu_session_id = req_rat_specific.pdu_session_id();
  itti_msg.error_code = static_cast<uint8_t>(status.error_code());

  send_n11_create_pdu_session_failure_itti(&itti_msg);
}
using namespace magma::lte;
namespace magma5g {

SetSMSessionContext create_sm_pdu_session(
    std::string& imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, std::string& ip4, std::string& ip6,
    const ambr_t& state_ambr, uint32_t version,
    const eps_subscribed_qos_profile_t& qos_profile) {
  magma::lte::SetSMSessionContext req;
  M5GQosInformationRequest qos_info;

  auto* req_common = req.mutable_common_context();

  // Encode IMSI
  req_common->mutable_sid()->set_id("IMSI" + imsi);

  // Encode TYPE IMSI
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);

  // Encode APU, storing apn value
  req_common->set_apn((char*)apn);

  // Encode RAT TYPE
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);

  // Put in CREATING STATE
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);

  // Create with Default Version
  req_common->set_sm_session_version(version);

  if (!ip4.empty()) {
    // Set the IPv4 PDU Address
    req_common->set_ue_ipv4(ip4);
  }

  if (!ip6.empty()) {
    // Set the IPv6 PDU Address
    req_common->set_ue_ipv6(ip6);
  }

  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  // Set the Session ID
  req_rat_specific->set_pdu_session_id(pdu_session_id);

  // Set the Type of Request
  req_rat_specific->set_request_type(magma::lte::RequestType::INITIAL_REQUEST);

  // TEID of GNB
  req_rat_specific->mutable_gnode_endpoint()->set_teid(gnb_gtp_teid);

  // IP Address of GNB

  char ipv4_str[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, gnb_gtp_teid_ip_addr, ipv4_str, INET_ADDRSTRLEN);
  req_rat_specific->mutable_gnode_endpoint()->set_end_ipv4_addr(ipv4_str);

  // Set the PTI
  req_rat_specific->set_procedure_trans_identity(pti);

  // qos_info
  qos_info.set_qos_class_id(static_cast<magma::lte::QCI>(qos_profile.qci));
  qos_info.set_priority_level(static_cast<magma::lte::prem_capab>(
      qos_profile.allocation_retention_priority.priority_level));
  qos_info.set_preemption_capability(static_cast<magma::lte::prem_capab>(
      qos_profile.allocation_retention_priority.pre_emp_capability));
  qos_info.set_preemption_vulnerability(static_cast<magma::lte::prem_vuner>(
      qos_profile.allocation_retention_priority.pre_emp_vulnerability));
  qos_info.set_apn_ambr_ul(state_ambr.br_ul);
  qos_info.set_apn_ambr_dl(state_ambr.br_dl);
  qos_info.set_br_unit(
      static_cast<magma::lte::M5GQosInformationRequest::BitrateUnitsAMBR>(
          state_ambr.br_unit));

  req_rat_specific->mutable_subscribed_qos()->CopyFrom(qos_info);

  return (req);
}

int AsyncSmfServiceClient::amf_smf_create_pdu_session(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ue_ipv4_addr, char* ue_ipv6_addr,
    const ambr_t& state_ambr, uint32_t version,
    const eps_subscribed_qos_profile_t& qos_profile) {
  std::string ip4_str, ip6_str;

  if (ue_ipv4_addr) {
    ip4_str = ue_ipv4_addr;
  }
  if (ue_ipv6_addr) {
    ip6_str = ue_ipv6_addr;
  }

  auto imsi_str = std::string(imsi);

  magma::lte::SetSMSessionContext req = create_sm_pdu_session(
      imsi_str, apn, pdu_session_id, pdu_session_type, gnb_gtp_teid, pti,
      gnb_gtp_teid_ip_addr, ip4_str, ip6_str, state_ambr, version, qos_profile);

  AsyncSmfServiceClient::getInstance().set_smf_session(req);
  return 0;
}

bool AsyncSmfServiceClient::set_smf_session(SetSMSessionContext& request) {
  SetSMFSessionRPC(
      request, [request](const Status& status, const SmContextVoid& response) {
        handle_session_context_response(status, request, response);
      });

  return true;
}

void AsyncSmfServiceClient::SetSMFSessionRPC(
    SetSMSessionContext& request,
    const std::function<void(Status, SmContextVoid)>& callback) {
  auto localResp = new AsyncLocalResponse<SmContextVoid>(std::move(callback),
                                                         RESPONSE_TIMEOUT);

  localResp->set_response_reader(std::move(stub_->AsyncSetAmfSessionContext(
      localResp->get_context(), request, &queue_)));
}

bool AsyncSmfServiceClient::set_smf_notification(
    const SetSmNotificationContext& notify) {
  SetSMFNotificationRPC(
      notify, [](const Status& status, const SmContextVoid& response) {
        handle_session_context_response(status, response);
      });

  return true;
}

void AsyncSmfServiceClient::SetSMFNotificationRPC(
    const SetSmNotificationContext& notify,
    const std::function<void(Status, SmContextVoid)>& callback) {
  auto localResp = new AsyncLocalResponse<SmContextVoid>(std::move(callback),
                                                         RESPONSE_TIMEOUT);

  localResp->set_response_reader(std::move(stub_->AsyncSetSmfNotification(
      localResp->get_context(), notify, &queue_)));
}

AsyncSmfServiceClient::AsyncSmfServiceClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "sessiond", ServiceRegistrySingleton::LOCAL);
  stub_ = AmfPduSessionSmContext::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

AsyncSmfServiceClient& AsyncSmfServiceClient::getInstance() {
  static AsyncSmfServiceClient instance;
  return instance;
}

bool AsyncSmfServiceClient::n11_update_location_req(
    const s6a_update_location_req_t* const ulr_p) {
  return s6a_update_location_req(ulr_p);
}

}  // namespace magma5g
