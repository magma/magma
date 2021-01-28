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
#ifndef AMF_ASDEFS_SEEN
#define AMF_ASDEFS_SEEN

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
using namespace std;
#define AUTH_KNAS_INT_SIZE 16 /* NAS integrity key     */
#define AUTH_KNAS_ENC_SIZE 16 /* NAS cyphering key     */
typedef location_area_identification_t LAI_t;
namespace magma5g {
// Data used to setup 5g CN NAS security /
typedef struct amf_as_security_data_s {
  amf_as_security_data_s() {}
  ~amf_as_security_data_s() {}
  bool is_new;                           // New security data indicator      /
  ksi_t ksi;                             // NAS key set identifier       /
  uint8_t sqn;                           // Sequence number          /
  uint32_t count;                        // NAS counter              /
  uint8_t knas_enc[AUTH_KNAS_ENC_SIZE];  // NAS cyphering key               /
  uint8_t knas_int[AUTH_KNAS_INT_SIZE];  // NAS integrity key               /
  bool is_knas_enc_present;
  bool is_knas_int_present;
} amf_as_security_data_t;

////////////////////typecast require///////////////////////////////
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * AMFAS-SAP primitives
 */
typedef enum amf_as_primitive_s {
  _AMFAS_START = 200,
  //_AMFREG_REGISTRATION_REJ,
  _AMFAS_SECURITY_REQ,   /* AMF->AS: Security request          */
  _AMFAS_SECURITY_IND,   /* AS->AMF: Security indication       */
  _AMFAS_SECURITY_RES,   /* AMF->AS: Security response         */
  _AMFAS_SECURITY_REJ,   /* AMF->AS: Security reject           */
  _AMFAS_ESTABLISH_REQ,  /* AMF->AS: Connection establish request  */
  _AMFAS_ESTABLISH_CNF,  /* AS->AMF: Connection establish confirm  */
  _AMFAS_ESTABLISH_REJ,  /* AS->AMF: Connection establish reject   */
  _AMFAS_RELEASE_REQ,    /* AMF->AS: Connection release request    */
  _AMFAS_RELEASE_IND,    /* AS->AMF: Connection release indication */
  _AMFAS_ERAB_SETUP_REQ, /* AMF->AS: ERAB setup request  */
  _AMFAS_ERAB_SETUP_CNF, /* AS->AMF  */
  _AMFAS_ERAB_SETUP_REJ, /* AS->AMF  */
  _AMFAS_DATA_REQ,       /* AMF->AS: Data transfer request     */
  _AMFAS_DATA_IND,       /* AS->AMF: Data transfer indication      */
  _AMFAS_PAGE_IND,       /* AS->AMF: Paging data indication        */
  _AMFAS_STATUS_IND,     /* AS->AMF: Status indication         */
  _AMFAS_ERAB_REL_CMD,   /* AMF->AS: ERAB Release Cmd  */
  _AMFAS_ERAB_REL_RSP,   /* AMF->AS: ERAB Release Rsp  */
  _AMFAS_END
} amf_as_primitive_t;
/*
 * AMFS primitive for data transfer
 * ---------------------------------
 */
static const char* amf_cn_primitive_str[] = {
    "AMF_CN_AUTHENTICATION_PARAM_RES",
    "AMF_CN_AUTHENTICATION_PARAM_FAIL",
    "AMF_CN_ULA_SUCCESS",
    "AMFCN_NW_INITIATED_DEREGISTRATION_UE",

};

typedef enum amfcn_primitive_s {
  _AMFCN_START = 400,
  AMFCN_AUTHENTICATION_PARAM_RES,
  AMFCN_AUTHENTICATION_PARAM_FAIL,
  AMFCN_NW_INITIATED_DEREGISTRATION_UE,
  AMFCN_DEACTIVATE_PDUSESSION_REQ,
  AMFCN_IDENTITY_PARAM_RES,
  AMFCN_SMC_PARAM_RES,
  _AMFCN_END
} amf_cn_primitive_t;

class amf_as_data_t {
 public:
  amf_as_data_t() {}
  ~amf_as_data_t() {}
  amf_ue_ngap_id_t ue_id; /* UE lower layer identifier        */
  // amf_as_M5GS_identity_t m5gs_id; /* UE's M5G mobile identity *///TODO
  guti_m5_t* guti; /* GUTI temporary mobile identity   */
  // guti_t* new_guti;       /* New GUTI, if re-allocated        */
  amf_as_security_data_t sctx; /* M5G NAS security context         */
  // uint8_t encryption : 4;       /* Ciphering algorithm              */
  // uint8_t integrity : 4;        /* Integrity protection algorithm   */
  // plmn_t* plmn_id;        /* Identifier of the selected PLMN  */
  plmn_t* plmn_id; /* Identifier of the selected PLMN  */
  ecgi_t ecgi;     /* NR CGI This information element is used to globally
                     identify a cell */
  // tai_t* tai; /* Code of the first tracking area identity the UE is
  // registered to          */ tai_list_t tai_list;                  /* Valid
  // field if num tai > 0 */

  bool switch_off; /* true if the UE is switched off   */
  uint8_t type;    /* Network deregister type          */
#define AMF_AS_DATA_DELIVERED_LOWER_LAYER_FAILURE 0
#define AMF_AS_DATA_DELIVERED_TRUE 1
#define AMF_AS_DATA_DELIVERED_LOWER_LAYER_NON_DELIVERY_INDICATION_DUE_TO_HO 2
  uint8_t delivered;                      /* Data message delivery indicator  */
#define AMF_AS_NAS_DATA_REGISTRATION 0x01 /* REGISTRATION Complete          */
#define AMF_AS_NAS_DATA_DEREGISTRATION_REQ                                     \
  0x02                           /* DEREGISTRATION Request           */
#define AMF_AS_NAS_DATA_TAU 0x03 /* TAU    REGISTRATION            */
#define AMF_AS_NAS_DATA_REGISTRATION_ACCEPT                                    \
  0x04                                  /* REGISTRATION Accept            */
#define AMF_AS_NAS_AMF_INFORMATION 0x05 /* Emm information          */
#define AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT                                  \
  0x06 /* DEREGISTRATION Accept            */
#define AMF_AS_NAS_DATA_CS_SERVICE_NOTIFICATION                                \
  0x07                                   /* CS Service Notification  */
#define AMF_AS_NAS_DATA_INFO_SR 0x08     /* Service Reject in DL NAS */
#define AMF_AS_NAS_DL_NAS_TRANSPORT 0x09 /* Downlink Nas Transport */
  uint8_t nas_info;    /* Type of NAS information to transfer  */
  std::string nas_msg; /* NAS message to be transferred     */
  std::string full_network_name;
  std::string short_network_name;
  std::uint8_t daylight_saving_time;
  // LAI_t* location_area_identification; /* Location area identification */
  mobile_identity_t*
      ms_identity; /* MS identity This IE may be included to assign or unassign
       a new TMSI to a UE during a combined TA/LA update. */
  std::uint8_t* additional_update_result; /* TAU Additional update result   */
  std::uint32_t* amf_cause;               /* EMM failure cause code        */
  std::string cli;                        /* Calling Line Identification  */

  void amf_as_set_security_data(
      amf_as_security_data_t* data, const void* context, bool is_new,
      bool is_ciphered);
};

/*
 * AMFAS primitive for status indication
 * -------------------------------------
 */
typedef struct amf_as_status_s {
  amf_as_status_s() {}
  ~amf_as_status_s() {}
  amf_ue_ngap_id_t ue_id;       // UE lower layer identifier        /
  guti_m5_t guti;               // GUTI temporary mobile identity   */
  amf_as_security_data_t sctx;  // 5G CN NAS security context     /
  int amf_cause;                // AMF failure cause code       /
} amf_as_status_t;

/*
 * AMFAS primitive for connection release
 * --------------------------------------
 */
typedef struct amf_as_release_s {
  amf_ue_ngap_id_t ue_id;  // UE lower layer identifier          /
  guti_m5_t guti;          // GUTI temporary mobile identity     */
#define AMF_AS_CAUSE_AUTHENTICATION 0x01  // Authentication failure /
#define AMF_AS_CAUSE_DEREGISTRATION 0x02  // DeRegistration requested   /
  uint8_t cause;                          // Release cause /
} amf_as_release_t;

typedef struct amf_as_PDUSESSION_identity_s {
  const guti_m5_t* guti; /* The GUTI, if valid               */
  const tai_t* last_tai; /* The last visited registered Tracking
                          * Area Identity, if available          */
  const imsi_t* imsi;    /* IMSI in case of "AttachWithImsi"     */
  const imei_t* imei;    /* UE's IMEI for emergency bearer services  */
} amf_as_PDUSESSION_identity_t;

typedef struct amf_as_establish_s {
  amf_as_establish_s() {}
  ~amf_as_establish_s() {}
  amf_ue_ngap_id_t ue_id;               // UE lower layer identifier         /
  uint64_t puid;                        // linked to procedure UID /
  amf_as_PDUSESSION_identity_t pds_id;  // UE's 5g cn mobile identity      /
  amf_as_security_data_t sctx;          /*5g cn NAS security context      */
  bool switch_off;                      // true if the UE is switched off    /
  bool is_initial;  // true if contained in initial message    /
  bool is_amf_ctx_new;
  uint8_t type;           // Network attach/detach type        /
  uint8_t m5g_rrc_cause;  // Connection establishment cause    /
  uint8_t m5g_rrc_type;   // Associated call type          /
  // plmn_t          plmn_id;                     / Identifier of the
  // selected PLMN   */
  ksi_t ksi;               // NAS key set identifier        /
  uint8_t encryption : 4;  // Ciphering algorithm           /
  uint8_t integrity : 4;   // Integrity protection algorithm    /
  uint8_t amf_cause;       // amf failure cause code        /
  guti_m5_t new_guti;      // New GUTI, if re-allocated         */
  int n_tacs;              // Number of consecutive tracking areas
                           //  the UE is registered to       /
  tai_t tai;               // The first tracking area the UE is registered to */
  tac_t tac;               // Code of the first
  // tracking area the UE is registered to /
  ecgi_t ecgi; /* E-UTRAN CGI This information element is used to globally
                  identify a cell */
#define AMF_AS_NAS_INFO_REGISTERD 0x01    // REGISTERD request        /
#define AMF_AS_NAS_INFO_DEREGISTERD 0x02  // DeREGISTERD request        /
#define AMF_AS_NAS_INFO_TAU 0x03          // Tracking Area Update request  /
#define AMF_AS_NAS_INFO_SR 0x04           // Service Request       /
#define AMF_AS_NAS_INFO_EXTSR 0x05        // Extended Service Request  /
#define AMF_AS_NAS_INFO_NONE 0xFF         // No Nas Message  /
  uint8_t nas_info;  // Type of initial NAS information to transfer   /
  bstring nas_msg;   // NAS message to be transferred within
                     //  initial NAS information message   /

  uint8_t m5gs_update_result;  // TAU EPS update result   /
  uint32_t t3412;              // GPRS T3412 timer   */
  guti_m5_t guti;              // TAU GUTI   */
  // tai_list_t tai_list;                // Valid field if num tai > 0 /
  uint16_t m5gs_pdusession_context_status;  // TAU EPS bearer context status */
  LAI_t location_area_identification;  // TAU Location area identification */
  mobile_identity_t
      ms_identity;              // TAU 8.2.26.7   MS identity This IE may be
                                //   included to assign or unassign a new TMSI
                                //   to a UE during a combined TA/LA update. */
  uint32_t t3402;               // TAU GPRS T3402 timer   */
  uint32_t t3423;               // TAU GPRS T3423 timer   */
  void* equivalent_plmns;       // TAU Equivalent PLMNs   */
  void* emergency_number_list;  // TAU Emergency number list   */
  uint8_t eps_network_feature_support;  // TAU Network feature support   */
  uint8_t additional_update_result;     // TAU Additional update result   */
  uint32_t t3412_extended;              // TAU GPRS timer   */

#define SERVICE_TYPE_PRESENT (1 << 0)
  uint8_t presencemask; /* Indicates the presence of some params like service
                           type */
  uint8_t service_type; /* Extended service request initiated for which service
                           type */
} amf_as_establish_t;

/*
 * AMFAS primitive for security
 * ----------------------------
 */
typedef struct amf_as_security_s {
  amf_as_security_s() {}
  ~amf_as_security_s() {}
  amf_ue_ngap_id_t ue_id;  // UE lower layer identifier        /
  // guti_m5_t guti;        // GUTI temporary mobile identity   */
  guti_m5_t guti;               // GUTI temporary mobile identity   */
  amf_as_security_data_t sctx;  // 5G CN NAS security context     /
  int amf_cause;                // AMF failure cause code       /
  uint64_t puid;                // linked to procedure UID
  /*
   * Identity request/response
   */
  uint8_t ident_type;  // Type of requested UE's identity  /
  imsi_t imsi;         // The requested IMSI of the UE     */
  imei_t imei;         // The requested IMEI of the UE     */
  uint8_t imeisv_request_enabled;
  uint32_t tmsi;  /// The requested TMSI of the UE     /
  /*
   * Authentication request/response
   */
  ksi_t ksi;                     // NAS key set identifier       /
  uint8_t rand[AUTH_RAND_SIZE];  // Random challenge number      /
  uint8_t autn[AUTH_AUTN_SIZE];  // Authentication token         /
  uint8_t res[AUTH_RES_SIZE];    // Authentication response      /
  uint8_t auts[AUTH_AUTS_SIZE];  // Synchronisation failure      /
  /*
   * Security Mode Command
   */
  uint8_t eea;  // Replayed EPS encryption algorithms   /
  uint8_t eia;  // Replayed EPS integrity algorithms    /
  uint8_t uea;  // Replayed UMTS encryption algorithms  /
  uint8_t ucs2;
  bool imeisv_request;

  // Added by LG
  uint8_t selected_eea;  // Selected EPS encryption algorithms   /
  uint8_t selected_eia;  // Selected EPS integrity algorithms    /

#define AMF_AS_MSG_TYPE_IDENT 0x01  // Identification message   /
#define AMF_AS_MSG_TYPE_AUTH 0x02   // Authentication message   /
#define AMF_AS_MSG_TYPE_SMC 0x03    // Security Mode Command    /
  uint8_t msg_type;  // Type of NAS security message to transfer /
} amf_as_security_t;

#if 0
typedef enum amf_as_primitive_u {
  _AMFAS_START = 200,
  _AMFAS_SECURITY_REQ,   // AMF->AS: Security request          /
  _AMFAS_SECURITY_IND,   // AS->AMF: Security indication       /
  _AMFAS_SECURITY_RES,   // AMF->AS: Security response         /
  _AMFAS_SECURITY_REJ,   // AMF->AS: Security reject           /
  _AMFAS_ESTABLISH_REQ,  // AMF->AS: Connection establish request  /
  _AMFAS_ESTABLISH_CNF,  // AS->AMF: Connection establish confirm  /
  _AMFAS_ESTABLISH_REJ,  // AS->AMF: Connection establish reject   /
  _AMFAS_RELEASE_REQ,    // AMF->AS: Connection release request    /
  _AMFAS_RELEASE_IND,    // AS->AMF: Connection release indication /
  _AMFAS_DATA_REQ,       // AMF->AS: Data transfer request     /
  _AMFAS_DATA_IND,       // AS->AMF: Data transfer indication      /
  _AMFAS_PAGE_IND,       // AS->AMF: Paging data indication        /
  _AMFAS_STATUS_IND,     // AS->AMF: Status indication         /
  _AMFAS_END
} amf_as_primitive_t;
#endif

/*
 * AMFAS primitive for cell information
 * ------------------------------------
 */
typedef struct amf_as_cell_info_s {
  uint8_t found;  // Indicates whether a suitable cell is found   /
#define AMF_AS_PLMN_LIST_SIZE 6
  PLMN_LIST_T(AMF_AS_PLMN_LIST_SIZE) plmn_ids;
  // List of identifiers of available PLMNs   /
  uint8_t rat;    // Bitmap of Radio Access Technologies      /
  tac_t tac;      // Tracking Area Code               /
  eci_t cell_id;  // cell identity                /
} amf_as_cell_info_t;

typedef struct as_primitive_s {
  as_primitive_s() {}
  ~as_primitive_s() {}
  amf_as_security_t security;
  amf_as_establish_t establish;
  amf_as_release_t release;
  amf_as_data_t data;
  // amf_as_page_t page;
  amf_as_status_t status;
  amf_as_cell_info_t cell_info;

} as_primitive_t;

typedef struct amf_as_s {
  amf_as_s() {}
  ~amf_as_s() {}
  amf_as_primitive_t primitive;
  as_primitive_t u;
} amf_as_t;

}  // namespace magma5g
#endif
