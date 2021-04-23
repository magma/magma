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
#pragma once

#include <mutex>
#include <unordered_map>

#include <grpc++/grpc++.h>
#include <orc8r/protos/common.pb.h>
#include <orc8r/protos/directoryd.pb.h>
#include <orc8r/protos/directoryd.grpc.pb.h>

#include "GRPCReceiver.h"

namespace magma {
using namespace orc8r;
using grpc::Status;

class DirectorydClient {
 public:
  virtual ~DirectorydClient() = default;

  /**
   * Gets the directoryd imsi's 'xid' field
   * @param imsi - UE to query
   */
  virtual void get_directoryd_xid_field(
      const std::string& ip,
      std::function<void(Status status, DirectoryField)> callback) = 0;
};
/**
 * AsyncDirectorydClient sends asynchronous calls to directoryd to retrieve
 * UE information.
 */
class AsyncDirectorydClient : public GRPCReceiver, public DirectorydClient {
 public:
  AsyncDirectorydClient();

  explicit AsyncDirectorydClient(
      std::shared_ptr<grpc::Channel> directoryd_channel);

  /**
   * Gets the directoryd imsi's 'xid' field
   * @param imsi - UE to query
   */
  void get_directoryd_xid_field(
      const std::string& ip,
      std::function<void(Status status, DirectoryField)> callback);

 private:
  static const uint32_t RESPONSE_TIMEOUT_SECONDS = 6;
  std::unique_ptr<GatewayDirectoryService::Stub> stub_;

 private:
  void get_directoryd_xid_field_rpc(
      const GetDirectoryFieldRequest& request,
      std::function<void(Status, DirectoryField)> callback);
};

}  // namespace magma
