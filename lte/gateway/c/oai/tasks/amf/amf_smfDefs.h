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

  Source      amf_asDefs.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include "securityDef.h"
#include "common_types.h"
#include "3gpp_24.008.h"
#include "TrackingAreaIdentity.h"
#include "amf_config.h"
#include "SmfMessage.h"
using namespace std;
#pragma once
namespace magma5g {
//------------------------------------------------------------------------------
// Causes related to invalid messages
//------------------------------------------------------------------------------
#define EMM_CAUSE_SEMANTICALLY_INCORRECT 95
#define EMM_CAUSE_INVALID_MANDATORY_INFO 96
#define AMF_CAUSE_INVALID_MANDATORY_INFO 96
#define EMM_CAUSE_MESSAGE_TYPE_NOT_IMPLEMENTED 97
#define EMM_CAUSE_MESSAGE_TYPE_NOT_COMPATIBLE 98
#define EMM_CAUSE_IE_NOT_IMPLEMENTED 99
#define AMF_CAUSE_IE_NOT_IMPLEMENTED 99
#define EMM_CAUSE_CONDITIONAL_IE_ERROR 100
#define EMM_CAUSE_MESSAGE_NOT_COMPATIBLE 101
#define EMM_CAUSE_PROTOCOL_ERROR 111
#define AMF_CAUSE_PROTOCOL_ERROR 111

#define PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t 0
#define PROCEDURE_TRANSACTION_IDENTITY_FIRST_t 1
#define PROCEDURE_TRANSACTION_IDENTITY_LAST_t 254
#define PROCEDURE_TRANSACTION_IDENTITY_RESERVED_t 255

//------------------------------------------------------------------------------
// Causes related to nature of request
//------------------------------------------------------------------------------
#define SMF_CAUSE_OPERATOR_DETERMINED_BARRING 8
#define SMF_CAUSE_INSUFFICIENT_RESOURCES 26
#define SMF_CAUSE_UNKNOWN_ACCESS_POINT_NAME 27
#define SMF_CAUSE_UNKNOWN_PDN_TYPE 28
#define SMF_CAUSE_USER_AUTHENTICATION_FAILED 29
#define SMF_CAUSE_REQUEST_REJECTED_BY_GW 30
#define SMF_CAUSE_REQUEST_REJECTED_UNSPECIFIED 31
#define SMF_CAUSE_SERVICE_OPTION_NOT_SUPPORTED 32
#define SMF_CAUSE_REQUESTED_SERVICE_OPTION_NOT_SUBSCRIBED 33
#define SMF_CAUSE_SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER 34
#define SMF_CAUSE_PTI_ALREADY_IN_USE 35
#define SMF_CAUSE_REGULAR_DEACTIVATION 36
#define SMF_CAUSE_EPS_QOS_NOT_ACCEPTED 37
#define SMF_CAUSE_NETWORK_FAILURE 38
#define SMF_CAUSE_REACTIVATION_REQUESTED 39
#define SMF_CAUSE_SEMANTIC_ERROR_IN_THE_TFT_OPERATION 41
#define SMF_CAUSE_SYNTACTICAL_ERROR_IN_THE_TFT_OPERATION 42
#define SMF_CAUSE_INVALID_EPS_BEARER_IDENTITY 43
#define SMF_CAUSE_SEMANTIC_ERRORS_IN_PACKET_FILTER 44
#define SMF_CAUSE_SYNTACTICAL_ERROR_IN_PACKET_FILTER 45
#define SMF_CAUSE_PTI_MISMATCH 47
#define SMF_CAUSE_LAST_PDN_DISCONNECTION_NOT_ALLOWED 49
#define SMF_CAUSE_PDN_TYPE_IPV4_ONLY_ALLOWED 50
#define SMF_CAUSE_PDN_TYPE_IPV6_ONLY_ALLOWED 51
#define SMF_CAUSE_SINGLE_ADDRESS_BEARERS_ONLY_ALLOWED 52
#define SMF_CAUSE_SMF_INFORMATION_NOT_RECEIVED 53
#define SMF_CAUSE_PDN_CONNECTION_DOES_NOT_EXIST 54
#define SMF_CAUSE_MULTIPLE_PDN_CONNECTIONS_NOT_ALLOWED 55
#define SMF_CAUSE_COLLISION_WITH_NETWORK_INITIATED_REQUEST 56
#define SMF_CAUSE_UNSUPPORTED_QCI_VALUE 59
#define SMF_CAUSE_BEARER_HANDLING_NOT_SUPPORTED 60
#define SMF_CAUSE_INVALID_PTI_VALUE 81
#define SMF_CAUSE_APN_RESTRICTION_VALUE_NOT_COMPATIBLE 112
#define SMF_CAUSE_REQUESTED_APN_NOT_SUPPORTED_IN_CURRENT_RAT 66
#define SMF_CAUSE_SUCCESS 0

// PDN Session Type
#define PDN_TYPE_IPV4 0
#define PDN_TYPE_IPV6 1
#define PDN_TYPE_IPV4V6 2

// PDN connection type
#define NET_PDN_TYPE_IPV4 (0 + 1)
#define NET_PDN_TYPE_IPV6 (1 + 1)
#define NET_PDN_TYPE_IPV4V6 (2 + 1)

typedef struct pkt_filter_s {
  uint8_t pkt_filter_dir;     // Packet filter direction
  uint8_t pkt_filter_id;      // Pkt filter id
  std::string contents[255];  // pkt filter contents
} pkt_filter_t;

// QOSRule
typedef struct qos_rule_s {
  uint8_t qos_rule_id;          // QOS rule id
  uint8_t rule_opercode;        // Rule operation code
  bool dqr_bit;                 // Default QOS Rule
  uint8_t no_of_pktfilters;     // Number of pkt filters
  pkt_filter_t pkt_filter[16];  // Pkt filter
  uint8_t qos_rule_precedence;  // Rule Precedence
  bool segregation;             // Segregation request reject
  uint8_t qfi;                  // QoS flow identifier
} qos_rule_t;

/*
 * AMFSMF gRPC/proto related elements for connection Establish
 * -----------------------------------------
 */
typedef struct amf_smf_establish_s {
  uint8_t pdu_session_id;    // Session Identity
  uint8_t pti;               // Procedure Tranction Identity
  uint8_t max_uplink;        // Integrity protection maximum data rate
  uint8_t max_downlink;      // Integrity protection minimum data rate
  uint8_t pdu_session_type;  // Session type
  uint8_t ssc_mode;          // SSC mode selection
  uint8_t gnb_gtp_teid[5];
  uint8_t gnb_gtp_teid_ip_addr[16];
  qos_rule_t
      qos_rules[4];  // QOS rules TODO verify index, 32 in nas5g, 3 in pcap
  uint8_t dl_unit;   // Session amber downlink unit
  uint16_t dl_session_ambr;  // Session amber downlink
  uint8_t ul_unit;           // Session amber uplink unit
  uint16_t ul_session_ambr;  // Session amber uplink
  uint8_t cause_value;       // M5GSMCause
  uint8_t address_info[12];  // PDU Address Info
} amf_smf_establish_t;

typedef struct amf_smf_modif_s {
  uint8_t pti;             // Procedure Tranction Identity
  uint8_t gnb_gtp_teid[5];
  uint8_t gnb_gtp_teid_ip_addr[16];
  qos_rule_t
      qos_rules[4];  // QOS rules TODO verify index, 32 in nas5g, 3 in pcap
 uint8_t cause_value;     // M5GSMCause
} amf_smf_modif_t;

/*
 * AMFSMF primitive for connection Release
 * ---------------------------------------
 */
typedef struct amf_smf_release_s {
  uint8_t pdu_session_id;  // Session Identity
  uint8_t pti;             // Procedure Tranction Identity
  uint8_t cause_value;     // M5GSMCause
} amf_smf_release_t;

typedef enum amf_smf_primitive_s {
  _AMFSMF_START = 400,

  _AMFSMF_END
} amf_smf_primitive_t;

typedef struct smf_primitive_s {
  smf_primitive_s(){};
  ~smf_primitive_s(){};
  amf_smf_establish_t establish;
  amf_smf_release_t release;
  amf_smf_modif_t modif;
} smf_primitive_t;

typedef struct amf_smf_s {
  amf_smf_s(){};
  ~amf_smf_s(){};
  amf_smf_primitive_t primitive;
  smf_primitive_t u;
} amf_smf_t;

class amf_smf_procedure_handler {
 public:
  int amf_smf_handle_pdu_establishment_request(
      SmfMsg* msg, amf_smf_t* amf_smf_msg);
  int amf_smf_handle_pdu_release_request(SmfMsg* msg, amf_smf_t* amf_smf_msg);
  int amf_smf_handle_pdu_modif_request(SmfMsg* msg, amf_smf_t* amf_smf_msg);
  int amf_smf_handle_pdu_modif_complete(SmfMsg* msg, amf_smf_t* amf_smf_msg);
  int amf_smf_handle_pdu_modif_cmd_reject(SmfMsg* msg, amf_smf_t* amf_smf_msg);
};

int create_session_grpc_req_on_gnb_setup_rsp(amf_smf_establish_t* message,
	       	char* imsi, uint32_t version);
int create_session_grpc_req(amf_smf_establish_t* message, char* imsi);
int release_session_gprc_req(amf_smf_release_t* message, char* imsi);
int mod_sessioncomp_grpc_req(amf_smf_modif_t* message, char* imsi);
int mod_sessionreq_grpc_req(amf_smf_modif_t* message, char* imsi);
int mod_sessioncmd_reject_grpc_req(amf_smf_modif_t* message, char* imsi);
}  // namespace magma5g
