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

  Source      3gpp_38413.h

  Date        2020/09/07

  Subsystem   NG Application Protocol IEs

  Description Defines NG Application Protocol Messages

*****************************************************************************/
#pragma once

#include "3gpp_23.003.h"
#include "bstrlib.h"
#include "n11_messages_types.h"

#define MAX_NO_OF_PDUSESSIONS 16

// 9.3.1.58 UE Aggregate Maximum Bit Rate
// all Non-GBR QoS flows per UE which is defined for the downlink and the uplink
// direction and a subscription parameter provided by the AMF to the NG-RAN node
typedef uint64_t bit_rate_t;
typedef struct ngap_ue_aggregate_maximum_bit_rate_s {
  bit_rate_t dl;
  bit_rate_t ul;
} ngap_ue_aggregate_maximum_bit_rate_t;

typedef long Ngap_PDUSessionID_t;

/* Ngap_S-NSSAI */
typedef struct Ngap_S_NSSAI_s {
  bstring sST;
  bstring* sD; /* OPTIONAL */

} Ngap_SNSSAI_t;

typedef struct pdusession_setup_item_s {
  bstring nas_pdu;  // 5GC – UE or UE – 5GC message that is transferred without
                    // interpretation in the NG-RAN node  /*optional*/
  Ngap_PDUSessionID_t
      Pdu_Session_ID;  // PDU Session for a UE. The definition and use of the
                       // PDU Session ID is specified in TS 23.501 [9].
  Ngap_SNSSAI_t Ngap_s_nssai;  // S-NSSAI as defined in TS 23.003 [23].
  pdu_session_resource_setup_request_transfer_t
      PDU_Session_Resource_Setup_Request_Transfer;  // Containing the PDU
                                                    // Session Resource
                                                    // Setup Request
  //  bstring PDU_Session_Resource_Setup_Failed_transfer;  // Containing the PDU
  //  Session
} pdusession_setup_item_t;

typedef struct Ngap_PDUSession_Resource_Setup_Request_List_s {
  uint16_t no_of_items;
  pdusession_setup_item_t item[MAX_NO_OF_PDUSESSIONS];

} Ngap_PDUSession_Resource_Setup_Request_List_t;


