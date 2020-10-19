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

#ifndef FILE_AMF_APP_MESSAGES_TYPES_SEEN
#define FILE_AMF_APP_MESSAGES_TYPES_SEEN
#include "common_types.h"

#define SMF_RESPONSE(mSGpTR) (mSGpTR)->ittiMsg.itti_smf_response

typedef enum SMSessionFSMState_response_s {
  CREATING_0,
  CREATE_1,
  ACTIVE_2,
  INACTIVE_3,
  RELEASED_4
} SMSessionFSMState_response;

typedef enum PduSessionType_response_s {
  IPV4,
  IPV6,
  IPV4IPV6,
  UNSTRUCTURED
} PduSessionType_response;

typedef enum SscMode_response_s {
  SSC_MODE_1,
  SSC_MODE_2,
  SSC_MODE_3
} SscMode_response;

typedef enum M5GSMCause_response_s {
  OPERATOR_DETERMINED_BARRING,
  INSUFFICIENT_RESOURCES,
  MISSING_OR_UNKNOWN_DNN,
  UNKNOWN_PDU_SESSION_TYPE,
  USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED,
  REQUEST_REJECTED_UNSPECIFIED
} M5GSMCause_response;

typedef enum RedirectAddressType_response_s {
  IPV4_1,
  IPV6_1,
  URL,
  SIP_URI
} RedirectAddressType_response;

typedef struct RedirectServer_response_s {
  RedirectAddressType_response redirect_address_type;
  char redirect_server_address[32];
} RedirectServer_response;

typedef struct itti_smf_response_s {
  //common context
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SMSessionFSMState_response sm_session_fsm_state;
  uint32_t sm_session_version;
  //M5GSMSessionContextAccess
  char pdu_session_id[2];
  PduSessionType_response pdu_session_type;
  SscMode_response selected_ssc_mode;
  M5GSMCause_response M5gsm_cause;
  bool always_on_pdu_session_indication;
  SscMode_response allowed_ssc_mode;
  bool M5gsm_congetion_re_attempt_indicator;
  RedirectServer_response pdu_address;
} itti_smf_response_t;

#endif /* FILE_AMF_APP_MESSAGES_TYPES_SEEN */
