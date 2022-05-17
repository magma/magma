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

#include <lte/protos/apn.pb.h>
#include <lte/protos/spgw_service.grpc.pb.h>
#include <lte/protos/subscriberdb.pb.h>
#include <stdint.h>
#include <functional>
#include <memory>
#include <string>
#include <vector>

#include "orc8r/gateway/c/common/async_grpc/GRPCReceiver.hpp"

namespace grpc {
class Channel;
}
namespace grpc {
class Status;
}
namespace grpc {
class Status;
}
namespace magma {
namespace lte {
class CreateBearerRequest;
}
}  // namespace magma
namespace magma {
namespace lte {
class CreateBearerResult;
}
}  // namespace magma
namespace magma {
namespace lte {
class DeleteBearerRequest;
}
}  // namespace magma
namespace magma {
namespace lte {
class DeleteBearerResult;
}
}  // namespace magma

using grpc::Status;

namespace magma {
using namespace lte;

/**
 * SpgwServiceClient is the base class for sending dedicated bearer
 * create/delete to PGW
 */
class SpgwServiceClient {
 public:
  virtual ~SpgwServiceClient() = default;

  /**
   * Delete a default bearer (all session bearers)
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @return true if the operation was successful
   */
  virtual bool delete_default_bearer(const std::string& imsi,
                                     const std::string& apn_ip_addr,
                                     const uint32_t linked_bearer_id) = 0;

  /**
   * Delete a dedicated bearer
   * @param DeleteBearerRequest
   * @return always returns true
   */
  virtual bool delete_dedicated_bearer(
      const magma::DeleteBearerRequest& request) = 0;

  /**
   * Create a dedicated bearer
   * @param CreateBearerRequest
   * @return always returns true
   */
  virtual bool create_dedicated_bearer(
      const magma::CreateBearerRequest& request) = 0;
};

/**
 * AsyncSpgwServiceClient implements SpgwServiceClient but sends calls
 * asynchronously to PGW.
 */
class AsyncSpgwServiceClient : public GRPCReceiver, public SpgwServiceClient {
 public:
  AsyncSpgwServiceClient();

  explicit AsyncSpgwServiceClient(std::shared_ptr<grpc::Channel> pgw_channel);
  /**
   * Delete a default bearer (all session bearers)
   * @param imsi - msi to identify a UE
   * @param apn_ip_addr - imsi and apn_ip_addrs identify a default bearer
   * @param linked_bearer_id - identifier for link bearer
   * @return true if the operation was successful
   */
  bool delete_default_bearer(const std::string& imsi,
                             const std::string& apn_ip_addr,
                             const uint32_t linked_bearer_id);

  /**
   * Delete a dedicated bearer
   * @param DeleteBearerRequest
   * @return always returns true
   */
  bool delete_dedicated_bearer(const magma::DeleteBearerRequest& request);

  /**
   * Create a dedicated bearer
   * @param CreateBearerRequest
   * @return always returns true
   */
  bool create_dedicated_bearer(const magma::CreateBearerRequest& request);

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<SpgwService::Stub> stub_;

 private:
  bool delete_bearer(const std::string& imsi, const std::string& apn_ip_addr,
                     const uint32_t linked_bearer_id,
                     const std::vector<uint32_t>& eps_bearer_ids);

  void delete_bearer_rpc(
      const DeleteBearerRequest& request,
      std::function<void(Status, DeleteBearerResult)> callback);

  void create_dedicated_bearer_rpc(
      const CreateBearerRequest& request,
      std::function<void(Status, CreateBearerResult)> callback);
};

}  // namespace magma
