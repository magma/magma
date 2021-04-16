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

#include "DirectorydClient.h"

#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

using grpc::Status;

namespace {  // anonymous

magma::GetDirectoryFieldRequest create_directory_field_req(
    const std::string& imsi) {
  magma::GetDirectoryFieldRequest req;
  req.set_id(imsi);
  req.set_field_key("xid");
  return req;
}

}  // namespace

namespace magma {

AsyncDirectorydClient::AsyncDirectorydClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(GatewayDirectoryService::NewStub(channel)) {}

AsyncDirectorydClient::AsyncDirectorydClient()
    : AsyncDirectorydClient(
          ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "directoryd", ServiceRegistrySingleton::LOCAL)) {}

void AsyncDirectorydClient::get_directoryd_xid_field(
    const std::string& ip,
    std::function<void(Status status, DirectoryField)> callback) {
  GetDirectoryFieldRequest req = create_directory_field_req(ip);
  get_directoryd_xid_field_rpc(req, callback);
}

void AsyncDirectorydClient::get_directoryd_xid_field_rpc(
    const GetDirectoryFieldRequest& request,
    std::function<void(Status, DirectoryField)> callback) {
  auto local_resp = new AsyncLocalResponse<DirectoryField>(
      std::move(callback), RESPONSE_TIMEOUT_SECONDS);
  local_resp->set_response_reader(std::move(stub_->AsyncGetDirectoryField(
      local_resp->get_context(), request, &queue_)));
}
}  // namespace magma
