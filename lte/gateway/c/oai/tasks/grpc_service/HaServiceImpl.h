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

#include "lte/protos/ha_service.grpc.pb.h"

extern "C" {
#include "log.h"
#include "hashtable.h"
}

/*
namespace grpc {
class ServerContext;
}  // namespace grpc

namespace magma {
namespace lte {
class StartAgwOffloadRequest;
class StartAgwOffloadResponse;
}  // namespace lte
}  // namespace magma

using grpc::ServerContext;
using magma::lte::HaService;
using magma::lte::StartAgwOffloadRequest;
using magma::lte::StartAgwOffloadResponse;
*/

bool process_ue_context(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP);

namespace magma {
using namespace lte;

class HaServiceImpl final : public HaService::Service {
 public:
  HaServiceImpl();

  /*
       * StartAgwOffload.
       *
       * @param context: the grpc Server context
       * @param request: StartAgwOffloadRequest
       * @param response (out): the StartAgwOffloadResponse that contains
                                err message.
       * @return grpc Status instance
       */
  grpc::Status StartAgwOffload(
      grpc::ServerContext* context, const StartAgwOffloadRequest* request,
      StartAgwOffloadResponse* response) override;
};

}  // namespace magma
