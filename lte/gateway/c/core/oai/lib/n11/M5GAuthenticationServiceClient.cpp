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

#include "lte/gateway/c/core/oai/lib/n11/M5GAuthenticationServiceClient.hpp"

#include <cstring>
#include <iostream>
#include <memory>
#include <string>
#include <thread>
#include <cassert>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.413.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include <google/protobuf/util/time_util.h>
#include <grpcpp/impl/codegen/client_context.h>
#include <grpcpp/impl/codegen/status.h>

#include "lte/protos/subscriberauth.grpc.pb.h"
#include "lte/protos/subscriberauth.pb.h"
#include "orc8r/gateway/c/common/service_registry/ServiceRegistrySingleton.hpp"
#include "lte/gateway/c/core/oai/lib/n11/amf_client_proto_msg_to_itti_msg.hpp"

using grpc::Channel;
using grpc::ClientContext;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::AsyncLocalResponse;
using magma::ServiceRegistrySingleton;
using magma::lte::M5GAuthenticationInformationAnswer;
using magma::lte::M5GAuthenticationInformationRequest;

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

static void handle_subs_authentication_info_ans(
    grpc::Status status,
    magma::lte::M5GAuthenticationInformationAnswer response,
    const std::string& imsi, uint8_t imsi_length, amf_ue_ngap_id_t ue_id) {
  MessageDef* message_p;
  message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, AMF_APP_SUBS_AUTH_INFO_RESP);

  itti_amf_subs_auth_info_ans_t* amf_app_subs_auth_info_resp_p;
  amf_app_subs_auth_info_resp_p =
      &message_p->ittiMsg.amf_app_subs_auth_info_resp;
  memset(amf_app_subs_auth_info_resp_p, 0,
         sizeof(itti_amf_subs_auth_info_ans_t));

  amf_app_subs_auth_info_resp_p->result = response.error_code();

  if (!status.ok()) {
    OAILOG_ERROR(LOG_AMF_APP,
                 "get_subs_auth_info fails with "
                 "code=%d and message=%s\n",
                 status.error_code(), status.error_message().c_str());
  }
  strncpy(amf_app_subs_auth_info_resp_p->imsi, imsi.c_str(), imsi_length);
  amf_app_subs_auth_info_resp_p->imsi_length = imsi_length;
  amf_app_subs_auth_info_resp_p->ue_id = ue_id;

  magma5g::convert_proto_msg_to_itti_m5g_auth_info_ans(
      response, amf_app_subs_auth_info_resp_p);

  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

namespace magma5g {

M5GAuthenticationInformationRequest create_subs_auth_request(
    const std::string& imsi, const std::string& snni) {
  M5GAuthenticationInformationRequest request;

  request.Clear();
  request.set_user_name(imsi);
  request.set_serving_network_name(snni);

  return request;
}

M5GAuthenticationInformationRequest create_subs_auth_request(
    const std::string& imsi, const std::string& snni, const void* resync_info,
    uint8_t resync_info_len) {
  M5GAuthenticationInformationRequest request;

  request.Clear();
  request.set_user_name(imsi);
  request.set_serving_network_name(snni);
  request.set_resync_info(resync_info, resync_info_len);

  return (request);
}

bool AsyncM5GAuthenticationServiceClient::get_subs_auth_info(
    const std::string& imsi, uint8_t imsi_length, const char* snni,
    amf_ue_ngap_id_t ue_id) {
  M5GAuthenticationInformationRequest request =
      create_subs_auth_request(imsi, snni);

  GetSubscriberAuthInfoRPC(
      request, [imsi, imsi_length, ue_id](
                   const Status& status,
                   const M5GAuthenticationInformationAnswer& response) {
        handle_subs_authentication_info_ans(status, response, imsi, imsi_length,
                                            ue_id);
      });
  return true;
}

bool AsyncM5GAuthenticationServiceClient::get_subs_auth_info_resync(
    const std::string& imsi, uint8_t imsi_length, const char* snni,
    const void* resync_info, uint8_t resync_info_len, amf_ue_ngap_id_t ue_id) {
  M5GAuthenticationInformationRequest request =
      create_subs_auth_request(imsi, snni, resync_info, resync_info_len);

  GetSubscriberAuthInfoRPC(
      request, [imsi, imsi_length, ue_id](
                   const Status& status,
                   const M5GAuthenticationInformationAnswer& response) {
        handle_subs_authentication_info_ans(status, response, imsi, imsi_length,
                                            ue_id);
      });
  return true;
}

void AsyncM5GAuthenticationServiceClient::GetSubscriberAuthInfoRPC(
    M5GAuthenticationInformationRequest& request,
    const std::function<void(Status, M5GAuthenticationInformationAnswer)>&
        callback) {
  auto localResp = new AsyncLocalResponse<M5GAuthenticationInformationAnswer>(
      std::move(callback), RESPONSE_TIMEOUT);
  localResp->set_response_reader(
      std::move(stub_->AsyncM5GAuthenticationInformation(
          localResp->get_context(), request, &queue_)));
}

AsyncM5GAuthenticationServiceClient::AsyncM5GAuthenticationServiceClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "subscriberdb", ServiceRegistrySingleton::LOCAL);
  stub_ = M5GSubscriberAuthentication::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

AsyncM5GAuthenticationServiceClient&
AsyncM5GAuthenticationServiceClient::getInstance() {
  static AsyncM5GAuthenticationServiceClient instance;
  return instance;
}

}  // namespace magma5g
