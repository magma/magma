/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#include <grpc++/grpc++.h>
#include <stdint.h>
#include <functional>
#include <memory>
#include <string>

#include "orc8r/protos/directoryd.grpc.pb.h"
#include "includes/GRPCReceiver.h"
#include "orc8r/protos/directoryd.pb.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc
namespace magma {
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma

using grpc::Channel;
using grpc::ClientContext;
using grpc::Status;

namespace magma {
using namespace orc8r;
/*
 * gRPC client for DirectoryService
 */
class GatewayDirectoryServiceClient : public GRPCReceiver {
 public:
  static bool UpdateRecord(
      const std::string& id, const std::string& location,
      std::function<void(Status, Void)> callback);

  static bool UpdateRecordField(
      const std::string& id, const std::string& field_key,
      const std::string& field_value,
      std::function<void(Status, Void)> callback);

  static bool DeleteRecord(
      const std::string& id, std::function<void(Status, Void)> callback);

 private:
  static bool updateRecordImpl(
      UpdateRecordRequest& request, std::function<void(Status, Void)> callback);

 public:
  GatewayDirectoryServiceClient(GatewayDirectoryServiceClient const&) = delete;
  void operator=(GatewayDirectoryServiceClient const&) = delete;

 private:
  GatewayDirectoryServiceClient();
  static GatewayDirectoryServiceClient& get_instance();
  std::shared_ptr<GatewayDirectoryService::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 30;  // seconds
};

}  // namespace magma
