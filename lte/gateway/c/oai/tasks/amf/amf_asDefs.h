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
#ifdef __cplusplus
extern "C" {
#endif
#include "3gpp_24.008.h"
#ifdef __cplusplus
};
#endif
#include "TrackingAreaIdentity.h"
#include "common_types.h"
#include "amf_config.h"
#include "amf_securityDef.h"

#define AUTH_KNAS_INT_SIZE 16 /* NAS integrity key     */
#define AUTH_KNAS_ENC_SIZE 16 /* NAS cyphering key     */

namespace magma5g {
// Data used to setup 5g CN NAS security /
typedef struct amf_as_security_data_s {
  bool is_new;                           // New security data indicator
  ksi_t ksi;                             // NAS key set identifier
  uint8_t sqn;                           // Sequence number
  uint32_t count;                        // NAS counter
  uint8_t knas_enc[AUTH_KNAS_ENC_SIZE];  // NAS cyphering key
  uint8_t knas_int[AUTH_KNAS_INT_SIZE];  // NAS integrity key
  bool is_knas_enc_present;
  bool is_knas_int_present;
} amf_as_security_data_t;

// AMFAS-SAP primitives
typedef enum amf_as_primitive_s {
  _AMFAS_START = 200,
  _AMFAS_SECURITY_REQ,   // AMF->AS: Security request
  _AMFAS_SECURITY_IND,   // AS->AMF: Security indication
  _AMFAS_SECURITY_RES,   // AMF->AS: Security response
  _AMFAS_SECURITY_REJ,   // AMF->AS: Security reject
  _AMFAS_ESTABLISH_REQ,  // AMF->AS: Connection establish request
  _AMFAS_ESTABLISH_CNF,  // AS->AMF: Connection establish confirm
  _AMFAS_ESTABLISH_REJ,  // AS->AMF: Connection establish reject
  _AMFAS_RELEASE_REQ,    // AMF->AS: Connection release request
  _AMFAS_RELEASE_IND,    // AS->AMF: Connection release indication
  _AMFAS_DATA_REQ,       // AMF->AS: Data transfer request
  _AMFAS_DATA_IND,       // AS->AMF: Data transfer indication
  _AMFAS_PAGE_IND,       // AS->AMF: Paging data indication
  _AMFAS_STATUS_IND,     // AS->AMF: Status indication
  _AMFAS_END
} amf_as_primitive_t;

typedef enum amf_cn_primitive_s {
  _AMFCN_START = 400,
  _AMFCN_AUTHENTICATION_PARAM_RES,
  _AMFCN_AUTHENTICATION_PARAM_FAIL,
  _AMFCN_NW_INITIATED_DEREGISTRATION_UE,
  _AMFCN_DEACTIVATE_PDUSESSION_REQ,
  _AMFCN_IDENTITY_PARAM_RES,
  _AMFCN_SMC_PARAM_RES,
  _AMFCN_END
} amf_cn_primitive_t;

// AMF to access related information
class amf_as_data_t {
 public:
  amf_as_data_t() {}
  ~amf_as_data_t() {}
  amf_ue_ngap_id_t ue_id;       // UE lower layer identifier
  guti_m5_t* guti;              // GUTI temporary mobile identity
  amf_as_security_data_t sctx;  // M5G NAS security context
#define AMF_AS_NAS_DATA_REGISTRATION_ACCEPT 0x04    // REGISTRATION Accept
#define AMF_AS_NAS_AMF_INFORMATION 0x05             // Amf information
#define AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT 0x06  // DEREGISTRATION Accept
#define AMF_AS_NAS_DL_NAS_TRANSPORT 0x09            // Downlink Nas Transport
  uint8_t nas_info;     // Type of NAS information to transfer
  std::string nas_msg;  // NAS message to be transferred
  void amf_as_set_security_data(
      amf_as_security_data_t* data, const void* context, bool is_new,
      bool is_ciphered);
};

typedef struct amf_as_pdusession_identity_s {
  const guti_m5_t* guti;  // The GUTI, if valid
} amf_as_pdusession_identity_t;

// Structure to handle UL/DL NAS message in AMF
typedef struct amf_as_establish_s {
  amf_ue_ngap_id_t ue_id;               // UE lower layer identifier
  uint64_t puid;                        // linked to procedure UID
  amf_as_pdusession_identity_t pds_id;  // UE's 5g cn mobile identity
  amf_as_security_data_t sctx;          // 5g cn NAS security context
  bool is_initial;                      // true if contained in initial message
  bool is_amf_ctx_new;
  uint8_t amf_cause;  // amf failure cause code
  tai_t tai;          // The first tracking area the UE is registered
  ecgi_t ecgi;  // E-UTRAN CGI This information element is used to globally
                // identify a cell
#define AMF_AS_NAS_INFO_REGISTERD 0x01  // REGISTERD request
#define AMF_AS_NAS_INFO_TAU 0x03        // Tracking Area Update request
#define AMF_AS_NAS_INFO_SR 0x04         // Service Request
#define AMF_AS_NAS_INFO_NONE 0xFF       // No Nas Message
  uint8_t nas_info;  // Type of initial NAS information to transfer
  bstring nas_msg;   // NAS message to be transferred within initial NAS
                     // information message
  guti_m5_t guti;    // TAU GUTI
  uint32_t t3502;    // TAU GPRS T3502 timer
  uint8_t
      presencemask;  // Indicates the presence of some params like service type
  uint8_t service_type;  // Extended service request initiated for which service
                         // type
} amf_as_establish_t;

/*
 * AMFAS primitive for security
 * ----------------------------
 */
typedef struct amf_as_security_s {
  amf_ue_ngap_id_t ue_id;       // UE lower layer identifier
  guti_m5_t guti;               // GUTI temporary mobile identity
  amf_as_security_data_t sctx;  // 5G CN NAS security context
  uint64_t puid;                // linked to procedure UID
  /*
   * Identity request/response
   */
  uint8_t ident_type;  // Type of requested UE's identity
  uint8_t imeisv_request_enabled;
  /*
   * Authentication request/response
   */
  ksi_t ksi;                     // NAS key set identifier
  uint8_t rand[AUTH_RAND_SIZE];  // Random challenge number
  uint8_t autn[AUTH_AUTN_SIZE];  // Authentication token
  /*
   * Security Mode Command
   */
  uint8_t eea;  // Replayed EPS encryption algorithms
  uint8_t eia;  // Replayed EPS integrity algorithms
  uint8_t ucs2;
  bool imeisv_request;

  uint8_t selected_eea;  // Selected EPS encryption algorithms
  uint8_t selected_eia;  // Selected EPS integrity algorithms

#define AMF_AS_MSG_TYPE_IDENT 0x01  // Identification message
#define AMF_AS_MSG_TYPE_AUTH 0x02   // Authentication message
#define AMF_AS_MSG_TYPE_SMC 0x03    // Security Mode Command
  uint8_t msg_type;                 // Type of NAS security message to transfer
} amf_as_security_t;

// AMF AS/CN primitives
typedef struct as_primitive_s {
  amf_as_security_t security;
  amf_as_establish_t establish;
  amf_as_data_t data;
} as_primitive_t;

typedef struct amf_as_s {
  amf_as_primitive_t primitive;
  as_primitive_t u;
} amf_as_t;

}  // namespace magma5g
