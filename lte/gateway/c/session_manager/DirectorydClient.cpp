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

#include <memory>
#include <utility>

#include "DirectorydClient.h"
#include "magma_logging.h"
#include "includes/ServiceRegistrySingleton.h"

using grpc::Status;

namespace magma {

AsyncDirectorydClient::AsyncDirectorydClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(GatewayDirectoryService::NewStub(channel)) {}

AsyncDirectorydClient::AsyncDirectorydClient()
    : AsyncDirectorydClient(
          ServiceRegistrySingleton::Instance()->GetGrpcChannel(
              "directoryd", ServiceRegistrySingleton::LOCAL)) {}

void AsyncDirectorydClient::update_directoryd_record(
    const UpdateRecordRequest& request,
    std::function<void(Status status, Void)> callback) {
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  local_response->set_response_reader(std::move(stub_->AsyncUpdateRecord(
      local_response->get_context(), request, &queue_)));
}

void AsyncDirectorydClient::delete_directoryd_record(
    const DeleteRecordRequest& request,
    std::function<void(Status status, Void)> callback) {
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  local_response->set_response_reader(std::move(stub_->AsyncDeleteRecord(
      local_response->get_context(), request, &queue_)));
}

void AsyncDirectorydClient::get_all_directoryd_records(
    std::function<void(Status status, AllDirectoryRecords)> callback) {
  magma::Void request;
  auto local_resp = new AsyncLocalResponse<AllDirectoryRecords>(
      std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->AsyncGetAllDirectoryRecords(
      local_resp->get_context(), request, &queue_)));
}
}  // namespace magma
