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

#include <glog/logging.h>                // for COMPACT_GOOGLE_LOG...
#include <grpcpp/channel.h>              // for Channel
#include <grpcpp/impl/codegen/status.h>  // for Status
#include <cstdint>                       // for uint32_t
#include <ostream>                       // for operator<<, basic_...
#include <utility>                       // for move

#include "SpgwServiceClient.h"
#include "includes/ServiceRegistrySingleton.h"  // for ServiceRegistrySin...
#include "lte/protos/policydb.pb.h"             // for RepeatedField, Rep...
#include "lte/protos/spgw_service.grpc.pb.h"    // for SpgwService::Stub
#include "lte/protos/spgw_service.pb.h"         // for DeleteBearerRequest
#include "lte/protos/subscriberdb.pb.h"
#include "magma_logging.h"  // for MLOG, MERROR, MINFO

using grpc::Status;

namespace {  // anonymous

magma::DeleteBearerRequest create_delete_bearer_req(
    const std::string& imsi, const std::string& apn_ip_addr,
    const uint32_t linked_bearer_id,
    const std::vector<uint32_t>& eps_bearer_ids) {
  magma::DeleteBearerRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(apn_ip_addr);
  req.set_link_bearer_id(linked_bearer_id);

  auto ebis = req.mutable_eps_bearer_ids();
  for (const auto& eps_bearer_id : eps_bearer_ids) {
    ebis->Add(eps_bearer_id);
  }

  return req;
}

}  // namespace

namespace magma {

AsyncSpgwServiceClient::AsyncSpgwServiceClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(SpgwService::NewStub(channel)) {}

AsyncSpgwServiceClient::AsyncSpgwServiceClient()
    : AsyncSpgwServiceClient(
          ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "spgw_service", ServiceRegistrySingleton::LOCAL)) {}

bool AsyncSpgwServiceClient::delete_default_bearer(
    const std::string& imsi, const std::string& apn_ip_addr,
    const uint32_t linked_bearer_id) {
  MLOG(MINFO) << "Deleting default bearer and corresponding PDN session for"
              << " IMSI: " << imsi << " APN IP addr " << apn_ip_addr
              << " Bearer ID " << linked_bearer_id;
  std::vector<uint32_t> eps_bearer_ids = {linked_bearer_id};
  return delete_bearer(imsi, apn_ip_addr, linked_bearer_id, eps_bearer_ids);
}

bool AsyncSpgwServiceClient::delete_dedicated_bearer(
    const magma::DeleteBearerRequest& request) {
  std::string bearer_ids = "{ ";
  for (const auto& bearer_id : request.eps_bearer_ids()) {
    bearer_ids += std::to_string(bearer_id) + " ";
  }
  bearer_ids += "}";
  MLOG(MINFO) << "Deleting dedicated bearers for " << request.sid().id() << ", "
              << request.ip_addr()
              << ", default bearer id=" << request.link_bearer_id()
              << ", dedicated bearer ids=" << bearer_ids;
  delete_bearer_rpc(request, [request](Status status, DeleteBearerResult resp) {
    if (!status.ok()) {
      // only log error for now
      MLOG(MERROR) << "Could not delete bearer" << request.sid().id() << ", "
                   << request.ip_addr() << ": " << status.error_message();
    }
  });
  return true;
}

bool AsyncSpgwServiceClient::create_dedicated_bearer(
    const magma::CreateBearerRequest& request) {
  std::string rule_ids = "{ ";
  for (const auto& rule : request.policy_rules()) {
    rule_ids += rule.id() + " ";
  }
  rule_ids += "}";
  MLOG(MINFO) << "Creating dedicated bearers for " << request.sid().id() << ", "
              << request.ip_addr() << ", rules=" << rule_ids;
  create_dedicated_bearer_rpc(
      request, [request](Status status, CreateBearerResult resp) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not create dedicated bearer"
                       << request.sid().id() << ", " << request.ip_addr()
                       << ": " << status.error_message();
        }
      });
  return true;
}

// delete_bearer creates the DeleteBearerRequest and logs the error
bool AsyncSpgwServiceClient::delete_bearer(
    const std::string& imsi, const std::string& apn_ip_addr,
    const uint32_t linked_bearer_id,
    const std::vector<uint32_t>& eps_bearer_ids) {
  auto req = create_delete_bearer_req(imsi, apn_ip_addr, linked_bearer_id,
                                      eps_bearer_ids);
  delete_bearer_rpc(
      req, [imsi, apn_ip_addr](Status status, DeleteBearerResult resp) {
        if (!status.ok()) {
          // only log error for now
          MLOG(MERROR) << "Could not delete bearer" << imsi << apn_ip_addr
                       << ": " << status.error_message();
        }
      });
  return true;
}

void AsyncSpgwServiceClient::delete_bearer_rpc(
    const DeleteBearerRequest& request,
    std::function<void(Status, DeleteBearerResult)> callback) {
  auto local_resp = new AsyncLocalResponse<DeleteBearerResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(
      stub_->AsyncDeleteBearer(local_resp->get_context(), request, &queue_));
}

void AsyncSpgwServiceClient::create_dedicated_bearer_rpc(
    const CreateBearerRequest& request,
    std::function<void(Status, CreateBearerResult)> callback) {
  auto local_resp = new AsyncLocalResponse<CreateBearerResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(
      stub_->AsyncCreateBearer(local_resp->get_context(), request, &queue_));
}

}  // namespace magma
