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

#include "lte/protos/ha_service.grpc.pb.h"

extern "C" {
#include "log.h"
#include "hashtable.h"
}

namespace magma {
namespace lte {

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

}  // namespace lte
}  // namespace magma
