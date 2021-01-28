

/*! \fopyright 2020 The Magma Authors.
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

  Source      amf_app_messages_types.h

  Version     0.1

  Date        2020/09/07

  Product     NGAP

  Subsystem   NG Application Protocol IEs

  Author      Sandeep Kumar Mall

 Description Defines NG Application Protocol Messages

*****************************************************************************/

#ifndef FILE_AMF_APP_MESSAGES_TYPES_SEEN
#define FILE_AMF_APP_MESSAGES_TYPES_SEEN

#include <stdint.h>

#include "bstrlib.h"
#include "3gpp_38.413.h"
#include "3gpp_24.007.h"
#include "3gpp_38.401.h"
#include "common_types.h"
#include "nas/securityDef.h"
//#include "nas/as_message.h"
#include "amf_as_message.h"
#include "TrackingAreaIdentity.h"

#define AMF_APP_CONNECTION_ESTABLISHMENT_CNF(mSGpTR)                           \
  (mSGpTR)->ittiMsg.amf_app_connection_establishment_cnf
#define AMF_APP_INITIAL_CONTEXT_SETUP_RSP(mSGpTR)                              \
  (mSGpTR)->ittiMsg.amf_app_initial_context_setup_rsp
#define AMF_APP_INITIAL_CONTEXT_SETUP_FAILURE(mSGpTR)                          \
  (mSGpTR)->ittiMsg.amf_app_initial_context_setup_failure
#define AMF_APP_NGAP_AMF_UE_ID_NOTIFICATION(mSGpTR)                            \
  (mSGpTR)->ittiMsg.amf_app_ngap_amf_ue_id_notification
#define AMF_APP_UL_DATA_IND(mSGpTR) (mSGpTR)->ittiMsg.amf_app_ul_data_ind
#define AMF_APP_DL_DATA_CNF(mSGpTR) (mSGpTR)->ittiMsg.amf_app_dl_data_cnf
#define AMF_APP_DL_DATA_REJ(mSGpTR) (mSGpTR)->ittiMsg.amf_app_dl_data_rej

typedef struct itti_amf_app_connection_establishment_cnf_s {
#if 0
  amf_ue_ngap_id_t ue_id;
  ambr_t ue_ambr;
  // E-RAB to Be Setup List
  uint8_t no_of_e_rabs;  // spec says max 256, actually stay with BEARERS_PER_UE
  //     >>E-RAB ID
  ebi_t e_rab_id[BEARERS_PER_UE];
  //     >>E-RAB Level QoS Parameters
  qci_t e_rab_level_qos_qci[BEARERS_PER_UE];
  //       >>>Allocation and Retention Priority
  priority_level_t e_rab_level_qos_priority_level[BEARERS_PER_UE];
  //       >>>Pre-emption Capability
  pre_emption_capability_t
      e_rab_level_qos_preemption_capability[BEARERS_PER_UE];
  //       >>>Pre-emption Vulnerability
  pre_emption_vulnerability_t
      e_rab_level_qos_preemption_vulnerability[BEARERS_PER_UE];
  //     >>Transport Layer Address
  bstring transport_layer_address[BEARERS_PER_UE];
  //     >>GTP-TEID
  teid_t gtp_teid[BEARERS_PER_UE];
  //     >>NAS-PDU (optional)
  bstring nas_pdu[BEARERS_PER_UE];
  //     >>Correlation ID TODO? later...

  // UE Security Capabilities
  uint16_t ue_security_capabilities_encryption_algorithms;
  uint16_t ue_security_capabilities_integrity_algorithms;

  // Security key
  uint8_t kgnb[AUTH_KGNB_SIZE];  // TODO -  NEED-RECHECK

  bstring ue_radio_capability;

  uint8_t presencemask;
#define NGAP_CSFB_INDICATOR_PRESENT (1 << 0)
  // ngap_csfb_indicator_t cs_fallback_indicator; cobraranu commented in header
  // Trace Activation (optional)
  // Handover Restriction List (optional)
  // UE Radio Capability (optional)
  // Subscriber Profile ID for RAT/Frequency priority (optional)
  // CS Fallback Indicator (optional)
  // SRVCC Operation Possible (optional)
  // CSG Membership Status (optional)
  // Registered LAI (optional)
  // GUAMFI ID (optional)
  // AMF UE NGAP ID 2  (optional)
  // Management Based MDT Allowed (optional)
#endif
  Ngap_initial_context_setup_request_t contextSetupRequest;
} itti_amf_app_connection_establishment_cnf_t;

typedef struct itti_amf_app_initial_context_setup_rsp_s {
  //  amf_ue_ngap_id_t ue_id;
  //  uint8_t no_of_e_rabs;
  //  ebi_t e_rab_id[BEARERS_PER_UE];
  //  bstring transport_layer_address[BEARERS_PER_UE];
  //  ngu_teid_t gtp_teid[BEARERS_PER_UE];
  amf_ue_ngap_id_t ue_id;
  Ngap_PDUSession_Resource_Setup_Request_List_t pdusesssion_setup_list;
  // Optional
  // e_rab_list_t e_rab_failed_to_setup_list;
} itti_amf_app_initial_context_setup_rsp_t;

typedef struct itti_amf_app_initial_context_setup_failure_s {
  uint32_t amf_ue_ngap_id;
} itti_amf_app_initial_context_setup_failure_t;

typedef struct itti_amf_app_delete_session_rsp_s {
  /* UE identifier */
  amf_ue_ngap_id_t ue_id;
} itti_amf_app_delete_session_rsp_t;

typedef struct itti_amf_app_ngap_amf_ue_id_notification_s {
  gnb_ue_ngap_id_t gnb_ue_ngap_id;
  amf_ue_ngap_id_t amf_ue_ngap_id;
  sctp_assoc_id_t sctp_assoc_id;
} itti_amf_app_ngap_amf_ue_id_notification_t;

typedef struct itti_amf_app_dl_data_cnf_s {
  amf_ue_ngap_id_t ue_id; /* UE lower layer identifier        */
  nas5g_error_code_t err_code; /* Transaction status*/  // TODO -  NEED-RECHECK
} itti_amf_app_dl_data_cnf_t;

typedef struct itti_amf_app_dl_data_rej_s {
  amf_ue_ngap_id_t ue_id; /* UE lower layer identifier   */
  bstring nas_msg;        /* Uplink NAS message           */
  int err_code;
} itti_amf_app_dl_data_rej_t;

typedef struct itti_amf_app_ul_data_ind_s {
  amf_ue_ngap_id_t ue_id; /* UE lower layer identifier    */
  bstring nas_msg;        /* Uplink NAS message           */
  /* Indicating the Tracking Area from which the UE has sent the NAS message */
  tai_t tai;  // TODO -  NEED-RECHECK
  /* Indicating the cell from which the UE has sent the NAS message  */
  ecgi_t cgi;
} itti_amf_app_ul_data_ind_t;

#endif /* FILE_AMF_APP_MESSAGES_TYPES_SEEN */
