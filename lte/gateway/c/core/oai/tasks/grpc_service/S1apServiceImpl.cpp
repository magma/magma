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
#include "lte/gateway/c/core/oai/tasks/grpc_service/S1apServiceImpl.h"

#include <memory>
#include <string>

#include "lte/protos/oai/s1ap_state.pb.h"
#include "orc8r/protos/common.pb.h"
#include "redis_utils/redis_client.h"

using grpc::ServerContext;
using grpc::Status;
using magma::lte::EnbStateResult;
using magma::lte::S1apService;
using magma::lte::oai::S1apState;
using magma::orc8r::Void;

namespace magma {
namespace lte {

S1apServiceImpl::S1apServiceImpl() : client_(nullptr) {}

void S1apServiceImpl::init(std::shared_ptr<RedisClient> client) {
  client_ = client;
}

Status S1apServiceImpl::GetENBState(
    ServerContext* context, const Void* request, EnbStateResult* response) {
  OAILOG_DEBUG(LOG_UTIL, "Received GetENBState GRPC request\n");

  S1apState s1ap_state_msg = S1apState();
  auto rc                  = client_->read_proto("s1ap_state", s1ap_state_msg);
  if (rc != RETURNok) {
    OAILOG_DEBUG(
        LOG_UTIL,
        "Failure while reading s1ap_state from redis. Cancelling GetENBState "
        "request");
    return Status::CANCELLED;
  }

  // Construct EnbStateResult
  for (const auto& enb_kv : s1ap_state_msg.enbs()) {
    (*response->mutable_enb_state_map())[enb_kv.second.enb_id()] =
        enb_kv.second.nb_ue_associated();
  }

  return Status::OK;
}

}  // namespace lte
}  // namespace magma
