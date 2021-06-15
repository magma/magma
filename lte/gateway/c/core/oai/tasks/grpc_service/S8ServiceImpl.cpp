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

extern "C" {
#include "intertask_interface.h"
#include "log.h"
#include "common_types.h"
#include "common_defs.h"
extern task_zmq_ctx_t grpc_service_task_zmq_ctx;
}
#include "S8ServiceImpl.h"
#include "s8_itti_proto_conversion.h"
#include "spgw_state_converter.h"

namespace grpc {
class ServerContext;
}  // namespace grpc


namespace magma {
namespace feg {
class CreateBearerRequestPgw;
}  // namespace feg
}  // namespace magma
using grpc::ServerContext;

namespace magma {

S8ServiceImpl::S8ServiceImpl() {}

static void convert_proto_msg_to_itti_create_bearer_req(
    const CreateBearerRequestPgw* request,
    s8_create_bearer_request_t* itti_msg) {
  itti_msg->context_teid         = request->c_agw_teid();
  itti_msg->linked_eps_bearer_id = request->linked_bearer_id();
  get_pco_from_proto_msg(
      request->protocol_configuration_options(), &itti_msg->pco);
  s8_bearer_context_t* s8_bc = &(itti_msg->bearer_context[0]);
  s8_bc->eps_bearer_id       = request->bearer_context().id();
  s8_bc->charging_id         = request->bearer_context().charging_id();

  if (request->bearer_context().has_qos()) {
    get_qos_from_proto_msg(request->bearer_context().qos(), &s8_bc->qos);
  }
  if (request->bearer_context().has_tft()) {
    magma::lte::SpgwStateConverter::proto_to_traffic_flow_template(
        request->bearer_context().tft(), &s8_bc->tft);
  }
  get_fteid_from_proto_msg(
      request->bearer_context().user_plane_fteid(), &s8_bc->pgw_s8_up);
  return;
}

grpc::Status CreateBearer(
    ServerContext* context, const CreateBearerRequestPgw* request,
    CreateBearerResponsePgw* response) {
  s8_create_bearer_request_t* cb_req = NULL;
  OAILOG_INFO(
      LOG_SGW_S8,
      " Received Create Bearer Request from roaming network's PGW"
      " for context teid: " TEID_FMT "\n",
      request->c_agw_teid());

  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, S8_CREATE_BEARER_REQ);
  if (!message_p) {
    OAILOG_ERROR(
        LOG_SGW_S8,
        "Failed to allocate memory for S8_CREATE_BEARER_REQ for "
        "context_teid" TEID_FMT "\n",
        request->c_agw_teid());
    return grpc::Status::OK;
  }

  cb_req = &message_p->ittiMsg.s8_create_bearer_req;
  convert_proto_msg_to_itti_create_bearer_req(request, cb_req);
  if ((send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SGW_S8, message_p)) !=
      RETURNok) {
    OAILOG_ERROR(
        LOG_SGW_S8,
        "Failed to send Create bearer request to sgw_s8 task for"
        "context_teid " TEID_FMT "\n",
        request->c_agw_teid());
  }

  return grpc::Status::OK;
}

}  // namespace magma
