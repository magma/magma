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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/common_types.h"

#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/state_converter.hpp"
#include "lte/protos/oai/std_3gpp_types.pb.h"
#include "lte/protos/oai/spgw_state.pb.h"
#include "lte/gateway/c/core/oai/include/spgw_types.hpp"
#include "lte/gateway/c/core/oai/include/pgw_types.h"
#include "lte/gateway/c/core/oai/tasks/sgw/pgw_procedures.hpp"
#include "lte/gateway/c/core/oai/include/spgw_state.hpp"

void eps_bearer_qos_to_proto(
    const bearer_qos_t* eps_bearer_qos_state,
    magma::lte::oai::SgwBearerQos* eps_bearer_qos_proto);

void traffic_flow_template_to_proto(
    const traffic_flow_template_t* tft_state,
    magma::lte::oai::TrafficFlowTemplate* tft_proto);

void sgw_create_session_message_to_proto(
    const itti_s11_create_session_request_t* session_request,
    magma::lte::oai::CreateSessionMessage* proto);

void proto_to_packet_filter(
    const magma::lte::oai::PacketFilter& packet_filter_proto,
    packet_filter_t* packet_filter);

void proto_to_traffic_flow_template(
    const magma::lte::oai::TrafficFlowTemplate& tft_proto,
    traffic_flow_template_t* tft_state);
