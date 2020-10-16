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

#include <gmp.h>

#include "feg/protos/csfb.grpc.pb.h"
#include "lte/protos/sms_orc8r.grpc.pb.h"
#include "sgs_messages_types.h"

extern "C" {
#include "intertask_interface.h"

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
class ResetIndication;
class Status;
}  // namespace feg
namespace lte {
class SMODownlinkUnitdata;
}  // namespace lte
}  // namespace magma
}

namespace magma {
using namespace lte;

void convert_proto_msg_to_itti_sgsap_downlink_unitdata(
    const SMODownlinkUnitdata* msg, itti_sgsap_downlink_unitdata_t* itti_msg);

}  // namespace magma

namespace magma {
using namespace feg;

void convert_proto_msg_to_itti_sgsap_eps_detach_ack(
    const EPSDetachAck* msg, itti_sgsap_eps_detach_ack_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_imsi_detach_ack(
    const IMSIDetachAck* msg, itti_sgsap_imsi_detach_ack_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_location_update_accept(
    const LocationUpdateAccept* msg,
    itti_sgsap_location_update_acc_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_location_update_reject(
    const LocationUpdateReject* msg,
    itti_sgsap_location_update_rej_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_paging_request(
    const PagingRequest* msg, itti_sgsap_paging_request_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_status_t(
    const Status* msg, itti_sgsap_status_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_vlr_reset_indication(
    const ResetIndication* msg, itti_sgsap_vlr_reset_indication_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_downlink_unitdata(
    const DownlinkUnitdata* msg, itti_sgsap_downlink_unitdata_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_release_req(
    const ReleaseRequest* msg, itti_sgsap_release_req_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_alert_request(
    const AlertRequest* msg, itti_sgsap_alert_request_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_service_abort_req(
    const ServiceAbortRequest* msg, itti_sgsap_service_abort_req_t* itti_msg);

void convert_proto_msg_to_itti_sgsap_mm_information_req(
    const MMInformationRequest* msg, itti_sgsap_mm_information_req_t* itti_msg);
}  // namespace magma
