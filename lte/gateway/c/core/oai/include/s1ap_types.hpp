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

#include <stdint.h>

#include "lte/protos/oai/s1ap_state.pb.h"

#include "lte/gateway/c/core/oai/include/proto_map.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.413.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/common_types.h"
#ifdef __cplusplus
}
#endif

#define S1AP_TIMER_INACTIVE_ID (-1)
#define S1AP_UE_CONTEXT_REL_COMP_TIMER 1  // in seconds

// Map- Key: sctp_assoc_id of (uint32_t), Data: EnbDescription
typedef magma::proto_map_s<uint32_t, magma::lte::oai::EnbDescription>
    proto_map_uint32_enb_description_t;

// Map- Key:comp_s1ap_id (uint64_t), Data: pointer to protobuf object,
// UeDescription
typedef magma::proto_map_s<uint64_t, magma::lte::oai::UeDescription*>
    map_uint64_ue_description_t;

/* Maximum no. of Broadcast PLMNs. Value is 6
 * 3gpp spec 36.413 section-9.1.8.4
 */
#define S1AP_MAX_BROADCAST_PLMNS 6
/* Maximum TAI Items configured, can be upto 256 */
#define S1AP_MAX_TAI_ITEMS 16

/* Supported TAI items includes TAC and Broadcast PLMNs */
typedef struct supported_tai_items_s {
  uint16_t tac;             ///< Supported TAC value
  uint8_t bplmnlist_count;  ///< Number of Broadcast PLMNs in the TAI
  plmn_t bplmns[S1AP_MAX_BROADCAST_PLMNS];  ///< List of Broadcast PLMNS
} supported_tai_items_t;

/* Supported TAs by eNB received in S1 Setup request message */
typedef struct supported_ta_list_s {
  uint8_t list_count;  ///< Number of TAIs in the list
  supported_tai_items_t
      supported_tai_items[S1AP_MAX_TAI_ITEMS];  ///< List of TAIs
} supported_ta_list_t;
