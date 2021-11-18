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

#include "feg/protos/s6a_proxy.grpc.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace feg {
class CancelLocationAnswer;
class CancelLocationRequest;
class ResetAnswer;
class ResetRequest;
}  // namespace feg
}  // namespace magma

using grpc::ServerContext;

namespace magma {
using namespace feg;
class S6aGatewayImpl final : public S6aGatewayService::Service {
 public:
  S6aGatewayImpl();

  /*
   * Cancel Location Request
   * S6a Command Code: 317
   *
   * @param context: the grpc Server context
   * @param request: CancelLocationRequest
   * @param response (out): CancelLocationAnswer
   * @return grpc Status instance
   */
  grpc::Status CancelLocation(
      ServerContext* context, const CancelLocationRequest* request,
      CancelLocationAnswer* response) override;
  /*
   * Reset Request
   * S6a Command Code: 322
   *
   * @param context: the grpc Server context
   * @param request: ResetRequest
   * @param response (out): ResetAnswer
   * @return grpc Status instance
   */
  grpc::Status Reset(
      ServerContext* context, const ResetRequest* request,
      ResetAnswer* response) override;
};

}  // namespace magma
