/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

#pragma once

#include <grpc++/grpc++.h>
#include <grpcpp/impl/codegen/status.h>

#include "feg/protos/s8_proxy.grpc.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace feg {
class CreateBearerRequestPgw;
}  // namespace feg
}  // namespace magma

namespace magma {
using namespace feg;

class S8ServiceImpl final : public S8ProxyResponder::Service {
 public:
  S8ServiceImpl();

  /* Create Bearer Request is sent from roaming network's PGW to initiate
   * dedicated bearer establishment
   */
  grpc::Status CreateBearer(grpc::ServerContext* context,
                            const CreateBearerRequestPgw* request,
                            orc8r::Void* response) override;

  grpc::Status DeleteBearerRequest(grpc::ServerContext* context,
                                   const DeleteBearerRequestPgw* request,
                                   orc8r::Void* response) override;
};

}  // namespace magma
