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

#include "lte/gateway/c/core/oai/lib/n11/M5GSUCIRegistrationServiceClient.h"

#include <cstring>
#include <iostream>
#include <memory>
#include <string>
#include <thread>
#include <cassert>

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.413.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.501.h"
#include "lte/gateway/c/core/oai/common/log.h"

#include <google/protobuf/util/time_util.h>
#include <grpcpp/impl/codegen/client_context.h>
#include <grpcpp/impl/codegen/status.h>

#include "lte/protos/subscriberdb.grpc.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "orc8r/gateway/c/common/service_registry/includes/ServiceRegistrySingleton.h"
#include "lte/gateway/c/core/oai/lib/n11/amf_client_proto_msg_to_itti_msg.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"

using grpc::Channel;
using grpc::ClientContext;
using grpc::InsecureChannelCredentials;
using grpc::Status;
using magma::AsyncLocalResponse;
using magma::ServiceRegistrySingleton;
using magma::lte::M5GSUCIRegistrationAnswer;
using magma::lte::M5GSUCIRegistrationRequest;
using magma5g::amf_proc_registration_reject;

extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

static void handle_decrypted_imsi_info_ans(
    grpc::Status status, magma::lte::M5GSUCIRegistrationAnswer response,
    amf_ue_ngap_id_t ue_id) {
  MessageDef* message_p;

  if (!status.ok() || (response.ue_msin_recv().length() == 0)) {
    std::cout << "get_decrypt_imsi_info fails with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
    MLOG(MERROR)
        << "   Error : Deconcealing IMSI Failed, sending Registration Reject";
    int amf_cause = AMF_UE_ILLEGAL;
    amf_proc_registration_reject(ue_id, amf_cause);
    return;
  }

  message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, AMF_APP_DECRYPT_IMSI_INFO_RESP);

  itti_amf_decrypted_imsi_info_ans_t* amf_app_decrypted_imsi_info_resp;
  amf_app_decrypted_imsi_info_resp =
      &message_p->ittiMsg.amf_app_decrypt_info_resp;
  memset(
      amf_app_decrypted_imsi_info_resp, 0,
      sizeof(itti_amf_decrypted_imsi_info_ans_t));

  magma5g::convert_proto_msg_to_itti_amf_decrypted_imsi_info_ans(
      response, amf_app_decrypted_imsi_info_resp);

  amf_app_decrypted_imsi_info_resp->ue_id = ue_id;

  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_AMF_APP, message_p);
}

namespace magma5g {

AsyncM5GSUCIRegistrationServiceClient::AsyncM5GSUCIRegistrationServiceClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "subscriberdb", ServiceRegistrySingleton::LOCAL);
  stub_ = M5GSUCIRegistration::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

AsyncM5GSUCIRegistrationServiceClient&
AsyncM5GSUCIRegistrationServiceClient::getInstance() {
  static AsyncM5GSUCIRegistrationServiceClient instance;
  return instance;
}

M5GSUCIRegistrationRequest create_decrypt_imsi_request(
    const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
    const std::string& ciphertext, const std::string& mac_tag) {
  M5GSUCIRegistrationRequest request;

  request.Clear();
  request.set_ue_pubkey_identifier(ue_pubkey_identifier);
  request.set_ue_pubkey(ue_pubkey);
  request.set_ue_ciphertext(ciphertext);
  request.set_ue_encrypted_mac(mac_tag);

  return request;
}

bool AsyncM5GSUCIRegistrationServiceClient::get_decrypt_imsi_info(
    const uint8_t ue_pubkey_identifier, const std::string& ue_pubkey,
    const std::string& ciphertext, const std::string& mac_tag,
    amf_ue_ngap_id_t ue_id) {
  M5GSUCIRegistrationRequest request = create_decrypt_imsi_request(
      ue_pubkey_identifier, ue_pubkey, ciphertext, mac_tag);

  GetSuciInfoRPC(
      request,
      [ue_id](const Status& status, const M5GSUCIRegistrationAnswer& response) {
        handle_decrypted_imsi_info_ans(status, response, ue_id);
      });
  return true;
}

void AsyncM5GSUCIRegistrationServiceClient::GetSuciInfoRPC(
    const M5GSUCIRegistrationRequest& request,
    const std::function<void(Status, M5GSUCIRegistrationAnswer)>& callback) {
  auto localResp = new AsyncLocalResponse<M5GSUCIRegistrationAnswer>(
      std::move(callback), RESPONSE_TIMEOUT);
  localResp->set_response_reader(
      std::move(stub_->AsyncM5GDecryptImsiSUCIRegistration(
          localResp->get_context(), request, &queue_)));
}

}  // namespace magma5g
