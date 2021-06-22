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

/*****************************************************************************
Source      emm_proc.h

Version     0.1

Date        2012/10/16

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines the EPS Mobility Management procedures executed at
        the EMM Service Access Points.

*****************************************************************************/
#ifndef FILE_EMM_PROC_SEEN
#define FILE_EMM_PROC_SEEN

#include "nas/commonDef.h"
#include "bstrlib.h"

#include "emm_data.h"

#include "nas_message.h"  //nas_message_decode_status_t
#include "mme_app_ue_context.h"
#include "MobileStationClassmark2.h"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/* Type of network attachment */
typedef enum {
  EMM_ATTACH_TYPE_EPS = 0,
  EMM_ATTACH_TYPE_COMBINED_EPS_IMSI,
  EMM_ATTACH_TYPE_EMERGENCY,
  EMM_ATTACH_TYPE_RESERVED,
} emm_proc_attach_type_t;

/* Type of network detach */
typedef enum {
  EMM_DETACH_TYPE_EPS = 0,
  EMM_DETACH_TYPE_IMSI,
  EMM_DETACH_TYPE_EPS_IMSI,
  EMM_DETACH_TYPE_REATTACH,
  EMM_DETACH_TYPE_NOT_REATTACH,
  EMM_DETACH_TYPE_RESERVED,
} emm_proc_detach_type_t;

/* Type of SGS detach */
typedef enum {
  EMM_SGS_UE_INITIATED_EPS_DETACH = 0,
  EMM_SGS_UE_INITIATED_EXPLICIT_NONEPS_DETACH,
  EMM_SGS_UE_INITIATED_COMBINED_DETACH,
  EMM_SGS_NW_INITIATED_EPS_DETACH,
  EMM_SGS_NW_INITIATED_IMPLICIT_NONEPS_DETACH,
  EMM_SGS_DETACH_TYPE_RESERVED,
} emm_proc_sgs_detach_type_t;

/* Type of requested identity */
typedef enum {
  EMM_IDENT_TYPE_NOT_AVAILABLE = 0,
  EMM_IDENT_TYPE_IMSI,
  EMM_IDENT_TYPE_IMEI,
  EMM_IDENT_TYPE_IMEISV,
  EMM_IDENT_TYPE_TMSI
} emm_proc_identity_type_t;

typedef enum {
  EMM_IMEISV_NOT_REQUESTED = 0,
  EMM_IMEISV_REQUESTED     = 1
} emm_imeisv_req_type_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

typedef struct emm_attach_request_ies_s {
  bool is_initial;
  emm_proc_attach_type_t type;
  additional_update_type_t additional_update_type;
  bool is_native_sc;
  ksi_t ksi;
  bool is_native_guti;
  guti_t* guti;
  imsi_t* imsi;
  imei_t* imei;
  tai_t* last_visited_registered_tai;
  tai_t* originating_tai;
  ecgi_t* originating_ecgi;
  ue_network_capability_t ue_network_capability;
  ms_network_capability_t* ms_network_capability;
  drx_parameter_t* drx_parameter;
  bstring esm_msg;
  nas_message_decode_status_t decode_status;
  MobileStationClassmark2*
      mob_st_clsMark2; /* Mobile station classmark2 provided by the UE */
  voice_domain_preference_and_ue_usage_setting_t*
      voicedomainpreferenceandueusagesetting;
  ue_additional_security_capability_t* ueadditionalsecuritycapability;
} emm_attach_request_ies_t;

typedef struct emm_detach_request_ies_s {
  emm_proc_detach_type_t type;
  bool switch_off;
  bool is_native_sc;
  ksi_t ksi;
  guti_t* guti;
  imsi_t* imsi;
  imei_t* imei;
  nas_message_decode_status_t decode_status;
} emm_detach_request_ies_t;

typedef struct emm_tau_request_ies_s {
  bool is_initial;
  emm_proc_attach_type_t type;
  EpsUpdateType eps_update_type;
  bool is_native_sc;
  ksi_t ksi;
  guti_t old_guti;

  /* Optional fields */
  bool is_native_non_current_sc;
  ksi_t non_current_ksi;
  guti_t* additional_guti;
  ue_network_capability_t* ue_network_capability;
  tai_t* last_visited_registered_tai;
  drx_parameter_t* drx_parameter;
  bool is_ue_radio_capability_information_update_needed;
  eps_bearer_context_status_t* eps_bearer_context_status;
  ms_network_capability_t* ms_network_capability;
  tmsi_status_t* tmsi_status;
  mobile_station_classmark2_t* mobile_station_classmark2;
  mobile_station_classmark3_t* mobile_station_classmark3;
  supported_codec_list_t* supported_codecs;
  additional_update_type_t* additional_updatetype;
  guti_type_t* old_guti_type;
  voice_domain_preference_and_ue_usage_setting_t*
      voicedomainpreferenceandueusagesetting;
  nas_message_decode_status_t decode_status;
} emm_tau_request_ies_t;
/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
 *---------------------------------------------------------------------------
 *              EMM status procedure
 *---------------------------------------------------------------------------
 */
int emm_proc_status_ind(mme_ue_s1ap_id_t ue_id, emm_cause_t emm_cause);
int emm_proc_status(mme_ue_s1ap_id_t ue_id, emm_cause_t emm_cause);

/*
 *---------------------------------------------------------------------------
 *              Lower layer procedure
 *---------------------------------------------------------------------------
 */

/*
 *---------------------------------------------------------------------------
 *              UE's Idle mode procedure
 *---------------------------------------------------------------------------
 */

/*
 * --------------------------------------------------------------------------
 *              Attach procedure
 * --------------------------------------------------------------------------
 */

void free_emm_attach_request_ies(emm_attach_request_ies_t** const params);

int emm_proc_attach_request(
    mme_ue_s1ap_id_t ue_id, const bool ctx_is_new,
    emm_attach_request_ies_t* const params);

int _emm_attach_reject(
    emm_context_t* emm_context, struct nas_base_proc_s* nas_base_proc);

int emm_proc_attach_reject(mme_ue_s1ap_id_t ue_id, emm_cause_t emm_cause);

int emm_proc_attach_complete(
    mme_ue_s1ap_id_t ue_id, const_bstring esm_msg_pP, int emm_cause,
    const nas_message_decode_status_t status);

void set_callbacks_for_attach_proc(nas_emm_attach_proc_t* attach_proc);

void free_emm_tau_request_ies(emm_tau_request_ies_t** const ies);

int emm_proc_tracking_area_update_request(
    const mme_ue_s1ap_id_t ue_id, emm_tau_request_ies_t* ies, int* emm_cause,
    tac_t tac);

int emm_proc_tracking_area_update_reject(
    const mme_ue_s1ap_id_t ue_id, const int emm_cause);

int emm_proc_service_reject(
    const mme_ue_s1ap_id_t ue_id, const uint8_t emm_cause);

int emm_proc_extended_service_request(
    const mme_ue_s1ap_id_t ue_id, const extended_service_request_msg* msg);

/*
 * --------------------------------------------------------------------------
 *              Detach procedure
 * --------------------------------------------------------------------------
 */

int emm_proc_detach(mme_ue_s1ap_id_t ue_id, emm_proc_detach_type_t type);
int emm_proc_sgs_detach_request(
    mme_ue_s1ap_id_t ue_id, emm_proc_sgs_detach_type_t type);
int emm_proc_nw_initiated_detach_request(
    mme_ue_s1ap_id_t ue_id, uint8_t detach_type);
void free_emm_detach_request_ies(emm_detach_request_ies_t** const ies);
int emm_proc_detach_request(
    mme_ue_s1ap_id_t ue_id, emm_detach_request_ies_t* params);
int emm_proc_detach_accept(mme_ue_s1ap_id_t ue_id);
/*
 * --------------------------------------------------------------------------
 *              Identification procedure
 * --------------------------------------------------------------------------
 */
struct emm_context_s;

int emm_proc_identification(
    struct emm_context_s* const emm_context, nas_emm_proc_t* const emm_proc,
    const identity_type2_t type, success_cb_t success, failure_cb_t failure);
int emm_proc_identification_complete(
    const mme_ue_s1ap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
    imeisv_t* const imeisv, uint32_t* const tmsi);

/*
 * --------------------------------------------------------------------------
 *              Authentication procedure
 * --------------------------------------------------------------------------
 */

int emm_proc_authentication_ksi(
    struct emm_context_s* emm_context,
    nas_emm_specific_proc_t* const emm_specific_proc, ksi_t ksi,
    const uint8_t* const rand, const uint8_t* const autn, success_cb_t success,
    failure_cb_t failure);

int emm_proc_authentication(
    struct emm_context_s* emm_context,
    nas_emm_specific_proc_t* const emm_specific_proc, success_cb_t success,
    failure_cb_t failure);

int emm_proc_authentication_failure(
    mme_ue_s1ap_id_t ue_id, int emm_cause, const_bstring auts);

int emm_proc_authentication_complete(
    mme_ue_s1ap_id_t ue_id, authentication_response_msg* msg, int emm_cause,
    const_bstring const res);

int emm_attach_security(struct emm_context_s* emm_context);

void set_notif_callbacks_for_auth_proc(nas_emm_auth_proc_t* auth_proc);
void set_callbacks_for_auth_proc(nas_emm_auth_proc_t* auth_proc);
void set_callbacks_for_auth_info_proc(nas_auth_info_proc_t* auth_info_proc);
/*
 * --------------------------------------------------------------------------
 *          Security mode control procedure
 * --------------------------------------------------------------------------
 */

int emm_proc_security_mode_control(
    struct emm_context_s* emm_context,
    nas_emm_specific_proc_t* const emm_specific_proc, ksi_t ksi,
    success_cb_t success, failure_cb_t failure);
int emm_proc_security_mode_complete(
    mme_ue_s1ap_id_t ue_id, const imeisv_mobile_identity_t* const imeisv);
int emm_proc_security_mode_reject(mme_ue_s1ap_id_t ue_id);
int emm_proc_emm_informtion(ue_mm_context_t* emm_ctx);

void _clear_emm_ctxt(emm_context_t* emm_ctx);

int emm_proc_tau_complete(mme_ue_s1ap_id_t ue_id);
int emm_send_service_reject_in_dl_nas(
    const mme_ue_s1ap_id_t ue_id, const uint8_t emm_cause);
int emm_proc_uplink_nas_transport(mme_ue_s1ap_id_t ue_id, bstring nas_msg);

void set_notif_callbacks_for_smc_proc(nas_emm_smc_proc_t* smc_proc);
void set_callbacks_for_smc_proc(nas_emm_smc_proc_t* smc_proc);
/*
 *---------------------------------------------------------------------------
 *             Network indication handlers
 *---------------------------------------------------------------------------
 */

#endif /* FILE_EMM_PROC_SEEN*/
