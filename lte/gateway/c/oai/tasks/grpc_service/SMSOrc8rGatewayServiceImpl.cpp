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
#include <iostream>

#include "SMSOrc8rGatewayServiceImpl.h"
#include "proto_msg_to_itti_msg.h"
#include "common_ies.h"
#include "sgs_messages_types.h"

extern "C" {
#include "sms_orc8r_service_handler.h"
#include "log.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace lte {
class SMODownlinkUnitdata;
}  // namespace lte
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma
}

using grpc::ServerContext;

namespace magma {

SMSOrc8rGatewayServiceImpl::SMSOrc8rGatewayServiceImpl() {}

grpc::Status SMSOrc8rGatewayServiceImpl::SMODownlink(
    ServerContext* context, const SMODownlinkUnitdata* request,
    Void* response) {
  itti_sgsap_downlink_unitdata_t itti_msg;
  convert_proto_msg_to_itti_sgsap_downlink_unitdata(request, &itti_msg);
  OAILOG_DEBUG(
      LOG_MME_APP,
      "Received SMS_ORC8R_DOWNLINK_UNITDATA message from orc8r with IMSI %s\n",
      itti_msg.imsi);
  handle_sms_orc8r_downlink_unitdata(&itti_msg);
  return grpc::Status::OK;
}

}  // namespace magma
