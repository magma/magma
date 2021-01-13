/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************

  Source      n11_messages_types.h 

  Date        2020/09/07

  Subsystem   NG Application Protocol IEs

  Description Defines NG Application Protocol Messages

*****************************************************************************/

#pragma once

#include "common_types.h"
//-----------------------------------------------------------------------------
/** @struct itti_n11_create_pdu_session_response_t
 *  @brief Create PDU Session Response */

typedef enum {
  SHALL_NOT_TRIGGER_PRE_EMPTION,
  MAY_TRIGGER_PRE_EMPTION,
} pre_emption_capability;

typedef enum {
  NOT_PREEMPTABLE,
  PRE_EMPTABLE,
} pre_emption_vulnerability;

typedef struct m5g_allocation_and_retention_priority_s {
  int priority_level;
  pre_emption_capability pre_emption_cap;
  pre_emption_vulnerability pre_emption_vul;
} m5g_allocation_and_retention_priority;

typedef struct non_dynamic_5QI_descriptor_s {
  int fiveQI;
} non_dynamic_5QI_descriptor;
// Dynamic_5QI not cosidered

typedef struct qos_characteristics_s {
  non_dynamic_5QI_descriptor non_dynamic_5QI_desc;
} qos_characteristics;

typedef struct qos_flow_level_qos_parameters_s {
  qos_characteristics qos_characteristic;
  m5g_allocation_and_retention_priority alloc_reten_priority;

} qos_flow_level_qos_parameters;

typedef struct qos_flow_setup_request_item_s {
  uint32_t qos_flow_identifier;
  qos_flow_level_qos_parameters qos_flow_level_qos_param;
  // E-RAB ID is optional spec-38413 - 9.3.4.1
} qos_flow_setup_request_item;

typedef struct qos_flow_request_list_s {
  qos_flow_setup_request_item qos_flow_req_item;

} qos_flow_request_list;

typedef struct amf_pdn_type_value_s {
  pdn_type_value_t pdn_type;
} amf_pdn_type_value_t;

typedef struct gtp_tunnel_s {
  bstring endpoint_ip_address;  // Transport_Layer_Information
  uint8_t gtp_tied[4];
} gtp_tunnel;

typedef struct up_transport_layer_information_s {
  gtp_tunnel gtp_tnl;
} up_transport_layer_information_t;

typedef struct amf_ue_aggregate_maximum_bit_rate_s {
  uint64_t dl;
  uint64_t ul;
} amf_ue_aggregate_maximum_bit_rate_t;

typedef struct pdu_session_resource_setup_request_transfer_s {
  amf_ue_aggregate_maximum_bit_rate_t pdu_aggregate_max_bit_rate;
  up_transport_layer_information_t up_transport_layer_info;
  amf_pdn_type_value_t pdu_ip_type;
  qos_flow_request_list qos_flow_setup_request_list;
} pdu_session_resource_setup_request_transfer_t;

