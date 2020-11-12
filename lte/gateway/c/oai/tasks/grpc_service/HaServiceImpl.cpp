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
#include <string>

#include "HaServiceImpl.h"
#include "lte/protos/ha_service.pb.h"
#include "mme_app_state_manager.h"
extern "C" {
#include "log.h"
}

// namespace grpc {
// class ServerContext;
// }  // namespace grpc

// using grpc::ServerContext;
// using grpc::Status;
// using magma::HaService;
// using magma::StartAgwOffloadRequest;
// using magma::StartAgwOffloadResponse;

namespace magma {
using namespace lte;

HaServiceImpl::HaServiceImpl() {}
/*
 * StartAgwOffload is called by North Bound to release UEs
 * based on eNB identification.
 */
grpc::Status HaServiceImpl::StartAgwOffload(
    grpc::ServerContext* context, const StartAgwOffloadRequest* request,
    StartAgwOffloadResponse* response) {
  OAILOG_INFO(LOG_UTIL, "Received StartAgwOffloadRequest GRPC request\n");
  hash_table_ts_t* state_imsi_ht =
      MmeNasStateManager::getInstance().get_ue_state_ht();
  hashtable_ts_apply_callback_on_elements(
      state_imsi_ht, process_ue_context, NULL, NULL);
  return grpc::Status::OK;
}  // namespace magma

}  // namespace magma

bool process_ue_context(
    const hash_key_t keyP, void* const elementP, void* parameterP,
    void** resultP) {
  OAILOG_INFO(LOG_UTIL, "Processing UE " IMSI_64_FMT "\n", keyP);
  return false;
}
