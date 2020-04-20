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

#ifdef __cplusplus
extern "C" {
#endif

#include "intertask_interface.h"
#include "sgs_messages_types.h"

void send_alert_ack(const itti_sgsap_alert_ack_t* msg);

void send_alert_reject(const itti_sgsap_alert_reject_t* msg);

void send_location_update_request(const itti_sgsap_location_update_req_t* msg);

void send_tmsi_reallocation_complete(
    const itti_sgsap_tmsi_reallocation_comp_t* msg);

void send_eps_detach_indication(const itti_sgsap_eps_detach_ind_t* msg);

void send_imsi_detach_indication(const itti_sgsap_imsi_detach_ind_t* msg);

void send_reset_ack(const itti_sgsap_vlr_reset_ack_t* msg);

void send_paging_reject(const itti_sgsap_paging_reject_t* msg);

void send_service_request(const itti_sgsap_service_request_t* msg);

void send_ue_activity_indication(const itti_sgsap_ue_activity_ind_t* msg);

void send_ue_unreachable(const itti_sgsap_ue_unreachable_t* msg);

void send_uplink_unitdata(const itti_sgsap_uplink_unitdata_t* msg);

#ifdef __cplusplus
}
#endif
