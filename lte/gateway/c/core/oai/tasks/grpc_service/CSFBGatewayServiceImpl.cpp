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

#include "CSFBGatewayServiceImpl.h"
#include "proto_msg_to_itti_msg.h"
#include "common_ies.h"
#include "sgs_messages_types.h"

extern "C" {
#include "sgs_service_handler.h"

namespace grpc {
class ServerContext;
}  // namespace grpc
namespace magma {
namespace feg {
class AlertRequest;
class DownlinkUnitdata;
class EPSDetachAck;
class IMSIDetachAck;
class LocationUpdateAccept;
class LocationUpdateReject;
class MMInformationRequest;
class PagingRequest;
class ReleaseRequest;
class ResetAck;
class ResetIndication;
class ServiceAbortRequest;
class Status;
}  // namespace feg
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma
}

using grpc::ServerContext;

namespace magma {

CSFBGatewayServiceImpl::CSFBGatewayServiceImpl() {}

grpc::Status CSFBGatewayServiceImpl::AlertReq(
    ServerContext* context, const AlertRequest* request, Void* response) {
  itti_sgsap_alert_request_t itti_msg;
  convert_proto_msg_to_itti_sgsap_alert_request(request, &itti_msg);
  std::cout << "[DEBUG] Received SGSAP_ALERT_REQUEST message from"
               "FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgsap_alert_request(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::Downlink(
    ServerContext* context, const DownlinkUnitdata* request, Void* response) {
  itti_sgsap_downlink_unitdata_t itti_msg;
  convert_proto_msg_to_itti_sgsap_downlink_unitdata(request, &itti_msg);
  std::cout << "[DEBUG] "
            << "Received SGSAP_DOWNLINK_UNITDATA message from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_downlink_unitdata(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::EPSDetachAc(
    ServerContext* context, const EPSDetachAck* request, Void* response) {
  itti_sgsap_eps_detach_ack_t itti_msg;
  convert_proto_msg_to_itti_sgsap_eps_detach_ack(request, &itti_msg);
  std::cout << "[DEBUG] Received SGSAP_EPS_DETACH_ACK message"
               "from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_eps_detach_ack(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::IMSIDetachAc(
    ServerContext* context, const IMSIDetachAck* request, Void* response) {
  itti_sgsap_imsi_detach_ack_t itti_msg;
  convert_proto_msg_to_itti_sgsap_imsi_detach_ack(request, &itti_msg);
  std::cout << "[DEBUG] Received SGSAP_IMSI_DETACH_ACK message"
               "from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_imsi_detach_ack(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::LocationUpdateAcc(
    ServerContext* context, const LocationUpdateAccept* request,
    Void* response) {
  itti_sgsap_location_update_acc_t itti_msg;
  convert_proto_msg_to_itti_sgsap_location_update_accept(request, &itti_msg);
  std::cout
      << "[DEBUG] "
      << "Received SGSAP_LOCATION_UPDATE_ACCEPT message from FeG with IMSI: "
      << itti_msg.imsi << std::endl;
  if (itti_msg.presencemask & SGSAP_MOBILE_IDENTITY) {
    std::cout << "[DEBUG] Mobile Identity presented." << std::endl;
    if (itti_msg.mobileid.typeofidentity == MOBILE_IDENTITY_IMSI) {
      std::cout << "[DEBUG] New IMSI included" << std::endl;
    } else {
      std::cout << "[DEBUG] New TMSI included" << std::endl;
    }
  }
  handle_sgs_location_update_accept(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::LocationUpdateRej(
    ServerContext* context, const LocationUpdateReject* request,
    Void* response) {
  itti_sgsap_location_update_rej_t itti_msg;
  convert_proto_msg_to_itti_sgsap_location_update_reject(request, &itti_msg);
  std::cout
      << "[DEBUG] "
      << "Received SGSAP_LOCATION_UPDATE_REJECT message from FeG with IMSI: "
      << itti_msg.imsi << std::endl;
  handle_sgs_location_update_reject(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::MMInformationReq(
    ServerContext* context, const MMInformationRequest* request,
    Void* response) {
  itti_sgsap_mm_information_req_t itti_msg;
  convert_proto_msg_to_itti_sgsap_mm_information_req(request, &itti_msg);
  std::cout << "[DEBUG] "
            << "Received SGSAP_MM_INFORMATION_REQ message from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_mm_information_request(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::PagingReq(
    ServerContext* context, const PagingRequest* request, Void* response) {
  itti_sgsap_paging_request_t itti_msg;
  convert_proto_msg_to_itti_sgsap_paging_request(request, &itti_msg);
  std::cout << "[DEBUG] Received SGSAP_PAGING_REQUEST message"
               "from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_paging_request(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::ReleaseReq(
    ServerContext* context, const ReleaseRequest* request, Void* response) {
  itti_sgsap_release_req_t itti_msg;
  convert_proto_msg_to_itti_sgsap_release_req(request, &itti_msg);
  std::cout << "[DEBUG] Received SGSAP_RELEASE_REQ message from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_release_req(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::ServiceAbort(
    ServerContext* context, const ServiceAbortRequest* request,
    Void* response) {
  itti_sgsap_service_abort_req_t itti_msg;
  convert_proto_msg_to_itti_sgsap_service_abort_req(request, &itti_msg);
  std::cout << "[DEBUG]"
               "Received SGSAP_SERVICE_ABORT_REQ message from FeG with IMSI: "
            << itti_msg.imsi << std::endl;
  handle_sgs_service_abort_req(&itti_msg);
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::VLRResetAck(
    ServerContext* context, const ResetAck* request, Void* response) {
  std::cout << "[DEBUG] Received SGSAP_RESET_ACK message from FeG" << std::endl;
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::VLRResetIndication(
    ServerContext* context, const ResetIndication* request, Void* response) {
  itti_sgsap_vlr_reset_indication_t itti_msg;
  convert_proto_msg_to_itti_sgsap_vlr_reset_indication(request, &itti_msg);
  handle_sgs_vlr_reset_indication(&itti_msg);
  std::cout << "[DEBUG] Received SGSAP_VLR_RESET_INDICATION message from FeG"
            << std::endl;
  return grpc::Status::OK;
}

grpc::Status CSFBGatewayServiceImpl::VLRStatus(
    ServerContext* context, const Status* request, Void* response) {
  itti_sgsap_status_t itti_msg;
  convert_proto_msg_to_itti_sgsap_status_t(request, &itti_msg);
  std::cout << "[DEBUG] Received SGSAP_STATUS message from FeG" << std::endl;
  return grpc::Status::OK;
}

}  // namespace magma
