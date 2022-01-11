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

#include "orc8r/gateway/c/common/service_registry/includes/ServiceRegistrySingleton.h"
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.h"
using grpc::Status;
using magma::AsyncLocalResponse;
using magma::ServiceRegistrySingleton;

void handle_session_context_response(
    grpc::Status status, magma::lte::SmContextVoid response) {
  if (!status.ok()) {
    std::cout << "AsyncSetAmfSessionContext fails with code "
              << status.error_code() << ", msg: " << status.error_message()
              << std::endl;
  }
}

using namespace magma::lte;
namespace magma5g {

SetSMSessionContext create_sm_pdu_session_v4(
    std::string& imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version,
    const ambr_t& state_ambr) {
  magma::lte::SetSMSessionContext req;

  auto* req_common = req.mutable_common_context();

  // Encode IMSI
  req_common->mutable_sid()->set_id("IMSI" + imsi);

  // Encode TYPE IMSI
  req_common->mutable_sid()->set_type(
      magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI);

  // Encode APU, storing apn value
  req_common->set_apn((char*) apn);

  // UE IPv4 address set
  req_common->set_ue_ipv4((char*) ipv4_addr);
  // Encode RAT TYPE
  req_common->set_rat_type(magma::lte::RATType::TGPP_NR);

  // Put in CREATING STATE
  req_common->set_sm_session_state(magma::lte::SMSessionFSMState::CREATING_0);

  // Create with Default Version
  req_common->set_sm_session_version(version);

  auto* req_rat_specific =
      req.mutable_rat_specific_context()->mutable_m5gsm_session_context();

  // Set the Session ID
  req_rat_specific->set_pdu_session_id(pdu_session_id);

  // Set the Type of Request
  req_rat_specific->set_request_type(magma::lte::RequestType::INITIAL_REQUEST);

  // Type is IPv4
  req_rat_specific->set_pdu_session_type(magma::lte::PduSessionType::IPV4);

  // TEID of GNB
  req_rat_specific->mutable_gnode_endpoint()->set_teid(gnb_gtp_teid);

  // IP Address of GNB

  char ipv4_str[INET_ADDRSTRLEN] = {0};
  inet_ntop(AF_INET, gnb_gtp_teid_ip_addr, ipv4_str, INET_ADDRSTRLEN);
  req_rat_specific->mutable_gnode_endpoint()->set_end_ipv4_addr(ipv4_str);

  // Set the PTI
  req_rat_specific->set_procedure_trans_identity((const char*) (&(pti)));

  // Set the default QoS values
  req_rat_specific->mutable_default_ambr()->set_max_bandwidth_ul(
      state_ambr.br_ul);
  req_rat_specific->mutable_default_ambr()->set_max_bandwidth_dl(
      state_ambr.br_dl);

  // For UT
  // state_ambr.br_unit = 1;

  req_rat_specific->mutable_default_ambr()->set_br_unit(
      static_cast<magma::lte::AggregatedMaximumBitrate::BitrateUnitsAMBR>(
          state_ambr.br_unit));

  return (req);
}

int AsyncSmfServiceClient::amf_smf_create_pdu_session_ipv4(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr, uint32_t version,
    const ambr_t& state_ambr) {
  auto imsi_str = std::string(imsi);

  magma::lte::SetSMSessionContext req = create_sm_pdu_session_v4(
      imsi_str, apn, pdu_session_id, pdu_session_type, gnb_gtp_teid, pti,
      gnb_gtp_teid_ip_addr, ipv4_addr, version, state_ambr);

  AsyncSmfServiceClient::getInstance().set_smf_session(req);
  return 0;
}

bool AsyncSmfServiceClient::set_smf_session(SetSMSessionContext& request) {
  SetSMFSessionRPC(
      request, [](const Status& status, const SmContextVoid& response) {
        handle_session_context_response(status, response);
      });

  return true;
}

void AsyncSmfServiceClient::SetSMFSessionRPC(
    SetSMSessionContext& request,
    const std::function<void(Status, SmContextVoid)>& callback) {
  auto localResp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);

  localResp->set_response_reader(std::move(stub_->AsyncSetAmfSessionContext(
      localResp->get_context(), request, &queue_)));
}

bool AsyncSmfServiceClient::set_smf_notification(
    SetSmNotificationContext& notify) {
  SetSMFNotificationRPC(
      notify, [](const Status& status, const SmContextVoid& response) {
        handle_session_context_response(status, response);
      });

  return true;
}

void AsyncSmfServiceClient::SetSMFNotificationRPC(
    SetSmNotificationContext& notify,
    const std::function<void(Status, SmContextVoid)>& callback) {
  auto localResp = new AsyncLocalResponse<SmContextVoid>(
      std::move(callback), RESPONSE_TIMEOUT);

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
