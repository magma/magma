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
#include "feg/protos/csfb.pb.h"
#include "sgs_messages_types.h"

extern "C" {
#include "intertask_interface.h"
}

namespace magma {
using namespace feg;

AlertAck convert_itti_sgsap_alert_ack_to_proto_msg(
    const itti_sgsap_alert_ack_t* msg);

AlertReject convert_itti_sgsap_alert_reject_to_proto_msg(
    const itti_sgsap_alert_reject_t* msg);

LocationUpdateRequest convert_itti_sgsap_location_update_req_to_proto_msg(
    const itti_sgsap_location_update_req_t* msg);

TMSIReallocationComplete convert_itti_sgsap_tmsi_reallocation_comp_to_proto_msg(
    const itti_sgsap_tmsi_reallocation_comp_t* msg);

EPSDetachIndication convert_itti_sgsap_eps_detach_ind_to_proto_msg(
    const itti_sgsap_eps_detach_ind_t* msg);

IMSIDetachIndication convert_itti_sgsap_imsi_detach_ind_to_proto_msg(
    const itti_sgsap_imsi_detach_ind_t* msg);

PagingReject convert_itti_sgsap_paging_reject_to_proto_msg(
    const itti_sgsap_paging_reject_t* msg);

ServiceRequest convert_itti_sgsap_service_request_to_proto_msg(
    const itti_sgsap_service_request_t* msg);

UEActivityIndication convert_itti_sgsap_ue_activity_indication_to_proto_msg(
    const itti_sgsap_ue_activity_ind_t* msg);

UEUnreachable convert_itti_sgsap_ue_unreachable_to_proto_msg(
    const itti_sgsap_ue_unreachable_t* msg);

UplinkUnitdata convert_itti_sgsap_uplink_unitdata_to_proto_msg(
    const itti_sgsap_uplink_unitdata_t* msg);

}  // namespace magma
