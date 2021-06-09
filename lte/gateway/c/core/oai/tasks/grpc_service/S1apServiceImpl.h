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

#include "orc8r/protos/common.pb.h"

#include "lte/protos/s1ap_service.grpc.pb.h"

extern "C" {
#include "log.h"
}

namespace magma {
using namespace lte;

class S1apServiceImpl final : public magma::S1apService::Service {
 public:
  S1apServiceImpl();

  /**
   * Returns map of S1 connected eNB id as key, with num of UEs connected
   * to each eNB as value
   * @param context grpc ServerContext
   * @param request proto request params
   * @param response proto response EnbStateResult
   * @return status response cod
   */
  grpc::Status GetENBState(
      grpc::ServerContext* context, const magma::orc8r::Void* request,
      magma::lte::EnbStateResult* response) override;
};

}  // namespace magma
