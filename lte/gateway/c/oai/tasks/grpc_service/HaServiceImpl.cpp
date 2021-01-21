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

#include <string>

#include "HaServiceImpl.h"
#include "lte/protos/ha_service.pb.h"
extern "C" {
#include "log.h"
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_types.h"
}
extern task_zmq_ctx_t grpc_service_task_zmq_ctx;

namespace magma {
namespace lte {

HaServiceImpl::HaServiceImpl() {}
/*
 * StartAgwOffload is called by North Bound to release UEs
 * based on eNB identification or UE identification.
 */
grpc::Status HaServiceImpl::StartAgwOffload(
    grpc::ServerContext* context, const StartAgwOffloadRequest* request,
    StartAgwOffloadResponse* response) {
  OAILOG_INFO(LOG_UTIL, "Received StartAgwOffloadRequest GRPC request\n");
  MessageDef* message_p =
      itti_alloc_new_message(TASK_GRPC_SERVICE, AGW_OFFLOAD_REQ);

  std::string imsi = request->imsi();
  // if IMSI prefix is used, strip it off
  if (imsi.compare(0, 4, "IMSI") == 0) {
    imsi = imsi.substr(4, std::string::npos);
  }
  AGW_OFFLOAD_REQ(message_p).imsi_length = imsi.size();
  strcpy(AGW_OFFLOAD_REQ(message_p).imsi, imsi.c_str());

  AGW_OFFLOAD_REQ(message_p).eNB_id = request->enb_id();
  AGW_OFFLOAD_REQ(message_p).enb_offload_type =
      (offload_type_t) request->enb_offload_type();

  send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_HA, message_p);
  return grpc::Status::OK;
}

}  // namespace lte
}  // namespace magma
