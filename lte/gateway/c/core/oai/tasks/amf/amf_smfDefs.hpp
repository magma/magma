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
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/include/nas/securityDef.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/SmfMessage.hpp"

#define NAS_MESSAGE_SECURITY_HEADER_SIZE 7
typedef uint8_t amf_cause_t;
namespace magma5g {
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
#define SMF_CAUSE_FAILURE 1

// PDN Session Type
#define PDN_TYPE_IPV4 0
#define PDN_TYPE_IPV6 1
#define PDN_TYPE_IPV4V6 2

/*
 * AMFSMF gRPC/proto related elements for connection Establish
 * -----------------------------------------
 */
typedef struct amf_smf_establish_s {
  uint32_t pdu_session_id;    // Session Identity
  uint8_t pti;                // Procedure Tranction Identity
  uint32_t pdu_session_type;  // Session type
  uint32_t gnb_gtp_teid;
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
  uint8_t pdu_session_id;
  smf_primitive_t u;
} amf_smf_t;

// Routines for communication from AMF to SMF on PDU sessions
int amf_smf_handle_pdu_establishment_request(SmfMsg* msg,
                                             amf_smf_t* amf_smf_msg);
int amf_smf_handle_pdu_release_request(SmfMsg* msg, amf_smf_t* amf_smf_msg);
status_code_e amf_smf_initiate_pdu_session_creation(
    amf_smf_establish_t* message, char* imsi, uint32_t version);

status_code_e amf_smf_create_session_req(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id,
    uint32_t pdu_session_type, uint32_t gnb_gtp_teid, uint8_t pti,
    uint8_t* gnb_gtp_teid_ip_addr, char* ue_ipv4_addr, char* ue_ipv6_addr,
    const ambr_t& state_ambr, const eps_subscribed_qos_profile_t& qos_profile);

status_code_e create_session_grpc_req_on_gnb_setup_rsp(
    amf_smf_establish_t* message, char* imsi, uint32_t version);
int create_session_grpc_req(amf_smf_establish_t* message, char* imsi);
status_code_e release_session_gprc_req(amf_smf_release_t* message, char* imsi);
}  // namespace magma5g
