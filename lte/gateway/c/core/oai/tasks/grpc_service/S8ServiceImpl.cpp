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
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
extern task_zmq_ctx_t grpc_service_task_zmq_ctx;
}
#include "lte/gateway/c/core/oai/tasks/grpc_service/S8ServiceImpl.hpp"
#include "lte/gateway/c/core/oai/lib/s8_proxy/s8_itti_proto_conversion.h"
#include "lte/gateway/c/core/oai/tasks/sgw/spgw_state_converter.hpp"

namespace grpc {
class ServerContext;
}  // namespace grpc

namespace magma {
namespace feg {
class CreateBearerRequestPgw;
}  // namespace feg
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma
using grpc::ServerContext;

namespace magma {

S8ServiceImpl::S8ServiceImpl() {}

static void convert_proto_msg_to_itti_create_bearer_req(
    const CreateBearerRequestPgw* request,
    s8_create_bearer_request_t* itti_msg) {
  auto ip = request->pgwaddrs();
  itti_msg->pgw_cp_address =
      reinterpret_cast<char*>(calloc(1, (ip.size() + 1)));
  snprintf(itti_msg->pgw_cp_address, (ip.size() + 1), "%s", ip.c_str());
  itti_msg->context_teid = request->c_agw_teid();
  itti_msg->sequence_number = request->sequence_number();
  itti_msg->linked_eps_bearer_id = request->linked_bearer_id();
  get_pco_from_proto_msg(request->protocol_configuration_options(),
                         &itti_msg->pco);
  s8_bearer_context_t* s8_bc = &(itti_msg->bearer_context[0]);
  s8_bc->eps_bearer_id = request->bearer_context().id();
  s8_bc->charging_id = request->bearer_context().charging_id();

  if (request->bearer_context().has_qos()) {
    get_qos_from_proto_msg(request->bearer_context().qos(), &s8_bc->qos);
  }
  if (request->bearer_context().has_tft()) {
    magma::lte::SpgwStateConverter::proto_to_traffic_flow_template(
        request->bearer_context().tft(), &s8_bc->tft);
  }
  get_fteid_from_proto_msg(request->bearer_context().user_plane_fteid(),
                           &s8_bc->pgw_s8_up);
  itti_msg->sgw_s8_up_teid = request->u_agw_teid();

  return;
}

grpc::Status S8ServiceImpl::CreateBearer(ServerContext* context,
                                         const CreateBearerRequestPgw* request,
                                         orc8r::Void* response) {
  OAILOG_INFO(LOG_SGW_S8,
              " Received Create Bearer Request from roaming network's PGW"
              " for context teid: " TEID_FMT "\n",
              request->c_agw_teid());

  MessageDef* message_p = NULL;
  s8_create_bearer_request_t* cb_req = NULL;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, S8_CREATE_BEARER_REQ);
  if (!message_p) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Failed to allocate memory for S8_CREATE_BEARER_REQ for "
                 "context_teid" TEID_FMT "\n",
                 request->c_agw_teid());
    return grpc::Status::CANCELLED;
  }

  cb_req = &message_p->ittiMsg.s8_create_bearer_req;
  convert_proto_msg_to_itti_create_bearer_req(request, cb_req);
  if ((send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SGW_S8, message_p)) !=
      RETURNok) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Failed to send Create bearer request to sgw_s8 task for"
                 "context_teid " TEID_FMT "\n",
                 request->c_agw_teid());
  }
  return grpc::Status::OK;
}

static void convert_proto_msg_to_itti_delete_bearer_req(
    const DeleteBearerRequestPgw* request,
    s8_delete_bearer_request_t* itti_msg_db_req) {
  itti_msg_db_req->sequence_number = request->sequence_number();
  itti_msg_db_req->context_teid = request->c_agw_teid();
  itti_msg_db_req->linked_eps_bearer_id = request->linked_bearer_id();
  for (int i = 0; i < request->eps_bearer_id_size(); i++) {
    itti_msg_db_req->eps_bearer_id[i] = request->eps_bearer_id(i);
    itti_msg_db_req->num_eps_bearer_id++;
  }
  get_pco_from_proto_msg(request->protocol_configuration_options(),
                         &itti_msg_db_req->pco);
}

grpc::Status S8ServiceImpl::DeleteBearerRequest(
    ServerContext* context, const DeleteBearerRequestPgw* request,
    orc8r::Void* response) {
  OAILOG_INFO(LOG_SGW_S8,
              " Received Delete Bearer Request from roaming network's PGW"
              " for context teid: " TEID_FMT "\n",
              request->c_agw_teid());

  MessageDef* message_p = NULL;
  s8_delete_bearer_request_t* db_req = NULL;
  message_p = itti_alloc_new_message(TASK_GRPC_SERVICE, S8_DELETE_BEARER_REQ);
  if (!message_p) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Failed to allocate memory for S8_DELETE_BEARER_REQ for "
                 "context_teid" TEID_FMT "\n",
                 request->c_agw_teid());
    return grpc::Status::CANCELLED;
  }

  db_req = &message_p->ittiMsg.s8_delete_bearer_req;
  convert_proto_msg_to_itti_delete_bearer_req(request, db_req);
  if ((send_msg_to_task(&grpc_service_task_zmq_ctx, TASK_SGW_S8, message_p)) !=
      RETURNok) {
    OAILOG_ERROR(LOG_SGW_S8,
                 "Failed to send Delete bearer request to sgw_s8 task for"
                 "context_teid " TEID_FMT "\n",
                 request->c_agw_teid());
  }
  return grpc::Status::OK;
}
}  // namespace magma
