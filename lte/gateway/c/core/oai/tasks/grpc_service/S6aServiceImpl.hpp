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

#include "lte/protos/s6a_service.grpc.pb.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace lte {
class DeleteSubscriberRequest;
class DeleteSubscriberResponse;
}  // namespace lte
}  // namespace magma

using grpc::ServerContext;
using magma::lte::DeleteSubscriberRequest;
using magma::lte::DeleteSubscriberResponse;
using magma::lte::S6aService;

namespace magma {
using namespace lte;

class S6aServiceImpl final : public S6aService::Service {
 public:
  S6aServiceImpl();

  /*
       * Deletes the subscribers in the DeleteSubscriberRequest.
       *
       * @param context: the grpc Server context
       * @param request: deleteSubscriberRequest, contains a list of IMSI of
                        subscriber to delete.
       * @param response (out): the DeleteSubscriberResponse that contains
                                err message.
       * @return grpc Status instance
       */
  grpc::Status DeleteSubscriber(
      ServerContext* context, const DeleteSubscriberRequest* request,
      DeleteSubscriberResponse* response) override;
};

}  // namespace magma
