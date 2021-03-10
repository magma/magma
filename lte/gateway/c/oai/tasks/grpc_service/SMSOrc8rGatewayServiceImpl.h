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
#include <grpcpp/impl/codegen/status.h>

#include "lte/protos/sms_orc8r.grpc.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace lte {
class SMODownlinkUnitdata;
}  // namespace lte
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma

using grpc::ServerContext;

namespace magma {
using namespace lte;
using namespace orc8r;

class SMSOrc8rGatewayServiceImpl final
    : public SMSOrc8rGatewayService::Service {
 public:
  SMSOrc8rGatewayServiceImpl();

  /*
   * Sent from the sms_orc8r service to the MME to transparently relay a NAS
   * message.
   *
   * @param context: the grpc Server context
   * @param request: SMODownlinkUnitdata, contains IMSI, NAS message container
   * @param response (out): Void defined in common.proto
   * @return grpc Status instance
   */
  grpc::Status SMODownlink(
      ServerContext* context, const SMODownlinkUnitdata* request,
      Void* response) override;
};

}  // namespace magma
