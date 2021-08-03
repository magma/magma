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

#pragma once
#include <sstream>
#include "securityDef.h"
#include "common_types.h"
#include "3gpp_24.008.h"
#include "TrackingAreaIdentity.h"
#include "amf_config.h"
#include "SmfMessage.h"

#define NAS_MESSAGE_SECURITY_HEADER_SIZE 7
typedef uint8_t amf_cause_t;
namespace magma5g {
#define OFFSET_OF(TyPe, MeMBeR) ((size_t) & ((TyPe*) 0)->MeMBeR)
#define PARENT_STRUCT(cOnTaiNeD, TyPe, MeMBeR)                                 \
  ({                                                                           \
    const typeof(((TyPe*) 0)->MeMBeR)* __MemBeR_ptr = (cOnTaiNeD);             \
    (TyPe*) ((char*) __MemBeR_ptr - OFFSET_OF(TyPe, MeMBeR));                  \
  })
//------------------------------------------------------------------------------
// Causes related to invalid messages
//------------------------------------------------------------------------------
#define AMF_CAUSE_INVALID_MANDATORY_INFO 96
#define AMF_CAUSE_IE_NOT_IMPLEMENTED 99
#define AMF_CAUSE_PROTOCOL_ERROR 111
#define PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t 0
#define PROCEDURE_TRANSACTION_IDENTITY_LAST_t 254

//------------------------------------------------------------------------------
// Causes related to nature of request
//------------------------------------------------------------------------------
#define SMF_CAUSE_UNKNOWN_PDN_TYPE 28
#define SMF_CAUSE_INVALID_PTI_VALUE 81
#define SMF_CAUSE_SUCCESS 0

// PDN Session Type
#define PDN_TYPE_IPV4 0
#define PDN_TYPE_IPV6 1
#define PDN_TYPE_IPV4V6 2

// PDN connection type
#define NET_PDN_TYPE_IPV4 (0 + 1)
#define NET_PDN_TYPE_IPV6 (1 + 1)
#define NET_PDN_TYPE_IPV4V6 (2 + 1)

/*
 * AMFSMF gRPC/proto related elements for connection Establish
 * -----------------------------------------
 */
typedef struct amf_smf_establish_s {
  uint32_t pdu_session_id;    // Session Identity
  uint8_t pti;                // Procedure Tranction Identity
  uint32_t pdu_session_type;  // Session type
  uint8_t gnb_gtp_teid[5];
  uint8_t gnb_gtp_teid_ip_addr[16];
  uint8_t cause_value;  // M5GSMCause
} amf_smf_establish_t;

/*
 * AMFSMF primitive for connection Establishment and Release
 * --------------------------------------------------------
 */
typedef struct amf_smf_release_s {
  uint32_t pdu_session_id;  // Session Identity
  uint8_t pti;              // Procedure Tranction Identity
  uint8_t cause_value;      // M5GSMCause
} amf_smf_release_t;

typedef struct smf_primitive_s {
  amf_smf_establish_t establish;
  amf_smf_release_t release;
} smf_primitive_t;

typedef struct amf_smf_s {
  smf_primitive_t u;
} amf_smf_t;

// Routines for communication from AMF to SMF on PDU sessions
int amf_smf_handle_pdu_establishment_request(
    SmfMsg* msg, amf_smf_t* amf_smf_msg);
int amf_smf_handle_pdu_release_request(SmfMsg* msg, amf_smf_t* amf_smf_msg);
int amf_smf_create_pdu_session(
    amf_smf_establish_t* message, char* imsi, uint32_t version);

int amf_smf_create_ipv4_session_grpc_req(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint8_t* gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ipv4_addr);

int create_session_grpc_req_on_gnb_setup_rsp(
    amf_smf_establish_t* message, char* imsi, uint32_t version);
int create_session_grpc_req(amf_smf_establish_t* message, char* imsi);
int release_session_gprc_req(amf_smf_release_t* message, char* imsi);
}  // namespace magma5g
