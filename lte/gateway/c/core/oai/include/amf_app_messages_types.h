/*
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

#include <stdint.h>
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.413.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/nas/securityDef.h"
#include "lte/gateway/c/core/oai/include/amf_as_message.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/include/nas/as_message.h"

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
#define AMF_APP_AUTH_RESPONSE_DATA(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.amf_app_subs_auth_info_resp
#define AMF_APP_DECRYPT_IMSI_RESPONSE_DATA(mSGpTR)                             \
  (mSGpTR)->ittiMsg.amf_app_decrypt_info_resp

typedef struct itti_amf_app_connection_establishment_cnf_s {
  Ngap_initial_context_setup_request_t contextSetupRequest;
} itti_amf_app_connection_establishment_cnf_t;

typedef struct itti_amf_app_initial_context_setup_rsp_s {
  amf_ue_ngap_id_t ue_id;
  Ngap_PDUSession_Resource_Setup_Response_List_t
      PDU_Session_Resource_Setup_Response_Transfer;  // Containing the PDU
                                                     // Session Resource
                                                     // Setup Request
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
  amf_ue_ngap_id_t ue_id;      /* UE lower layer identifier        */
  nas5g_error_code_t err_code; /* Transaction status*/
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
  tai_t tai;
  /* Indicating the cell from which the UE has sent the NAS message  */
  ecgi_t cgi;
} itti_amf_app_ul_data_ind_t;

typedef struct itti_amf_subs_auth_info_ans_s {
  /* IMSI of the subscriber */
  char imsi[IMSI_BCD_DIGITS_MAX + 1];

  /* Length of the Imsi. Mostly 15 */
  uint8_t imsi_length;

  /* Authentication is success or failure with code */
  int result;

  /* UE identifier */
  amf_ue_ngap_id_t ue_id;

  /* Authentication info containing the vector(s) */
  m5g_authentication_info_t auth_info;

} itti_amf_subs_auth_info_ans_t;

typedef struct itti_amf_decrypted_imsi_info_ans_s {
  /* IMSI of the subscriber */
  char imsi[IMSI_BCD_DIGITS_MAX + 1];

  /* Length of the Imsi. Mostly 15 */
  uint8_t imsi_length;

  /* Authentication is success or failure with code */
  int result;

  /* UE identifier */
  amf_ue_ngap_id_t ue_id;
} itti_amf_decrypted_imsi_info_ans_t;

typedef struct itti_amf_ip_allocation_response_s {
  /* IMSI of the subscriber */
  char imsi[IMSI_BCD_DIGITS_MAX + 1];

  /* Length of the Imsi. Mostly 15 */
  uint8_t imsi_length;

  uint8_t apn[101];

  /* PDU Session Identity */
  uint32_t pdu_session_id;

  /* Procedure Tranction Identity */
  uint8_t pti;

  /* PDU Session Type */
  uint32_t pdu_session_type;

  /* GNB GTP Tunnel Identifier */
  uint32_t gnb_gtp_teid;

  /* GNB End Point IP address */
  uint8_t gnb_gtp_teid_ip_addr[16];

  /* UE IP Address Allocated by MobilityD */
  paa_t paa;

  /* Result Type of IP Allocation */
  int result;

  ambr_t default_ambr;
} itti_amf_ip_allocation_response_t;
