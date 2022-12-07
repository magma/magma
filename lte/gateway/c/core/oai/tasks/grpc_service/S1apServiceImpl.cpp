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

#include "lte/gateway/c/core/oai/tasks/grpc_service/S1apServiceImpl.hpp"
#include <string>
#include "lte/gateway/c/core/oai/include/s1ap_state.hpp"

extern "C" {
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/hashtable/hashtable.h"
}

using grpc::ServerContext;
using grpc::Status;
using magma::EnbStateResult;
using magma::S1apService;

namespace magma {
using namespace lte;
using namespace orc8r;

S1apServiceImpl::S1apServiceImpl() {}

Status S1apServiceImpl::GetENBState(ServerContext* context, const Void* request,
                                    EnbStateResult* response) {
  OAILOG_DEBUG(LOG_UTIL, "Received GetENBState GRPC request\n");

  // Get state from S1APStateManager
  // TODO: Get state through ITTI message from S1AP task, as it's read only
  // it will not affect ownership
  oai::S1apState* s1ap_state = get_s1ap_state(false);
  if (s1ap_state != nullptr) {
    if (!(s1ap_state->enbs_size())) {
      return Status::OK;
    }
    proto_map_uint32_enb_description_t enb_map;
    enb_map.map = s1ap_state->mutable_enbs();
    for (auto itr = enb_map.map->begin(); itr != enb_map.map->end(); itr++) {
      oai::EnbDescription enb_ref = itr->second;
      if (enb_ref.sctp_assoc_id()) {
        (*response->mutable_enb_state_map())[enb_ref.enb_id()] =
            enb_ref.nb_ue_associated();
      }
    }
  }
  return Status::OK;
}

}  // namespace magma
