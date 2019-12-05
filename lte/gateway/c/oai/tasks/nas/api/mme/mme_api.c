/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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
  Source      mme_api.c

  Version     0.1

  Date        2013/02/28

  Product     NAS stack

  Subsystem   Application Programming Interface

  Author      Frederic Maurel, Lionel GAUTHIER

  Description Implements the API used by the NAS layer running in the MME
        to interact with a Mobility Management Entity

*****************************************************************************/
#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>
#include <string.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "conversions.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "common_types.h"
#include "common_defs.h"
#include "mme_api.h"
#include "mme_app_ue_context.h"
#include "mme_config.h"
#include "emm_data.h"
#include "EpsNetworkFeatureSupport.h"
#include "mme_app_state.h"

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* Maximum number of PDN connections the MME may simultaneously support */
#define MME_API_PDN_MAX 10

/* Subscribed QCI */
#define MME_API_QCI 3

/* Data bit rates */
#define MME_API_BIT_RATE_64K 0x40
#define MME_API_BIT_RATE_128K 0x48
#define MME_API_BIT_RATE_512K 0x78
#define MME_API_BIT_RATE_1024K 0x87

/* Total number of PDN connections (should not exceed MME_API_PDN_MAX) */
static int _mme_api_pdn_id = 0;

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    mme_api_get_emm_config()                                  **
 **                                                                        **
 ** Description: Retreives MME configuration data related to EPS mobility  **
 **      management                                                **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int mme_api_get_emm_config(
  mme_api_emm_config_t *config,
  struct mme_config_s *mme_config_p)
{
  OAILOG_FUNC_IN(LOG_NAS);
  AssertFatal(mme_config_p->served_tai.nb_tai >= 1, "No TAI configured");
  AssertFatal(mme_config_p->gummei.nb >= 1, "No GUMMEI configured");

  config->tai_list.numberoflists = 0;
  // TODO actually we support only one partial TAI list.
  // reminder mme_config_p->served_tai is sorted in ascending order of TAIs
  switch (mme_config_p->served_tai.list_type) {
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS:
      config->tai_list.numberoflists = 1;
      config->tai_list.partial_tai_list[0].typeoflist =
        mme_config_p->served_tai.list_type;
      // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2 elements...), see 3GPP TS 24.301, section 9.9.3.33.1
      config->tai_list.partial_tai_list[0].numberofelements =
        mme_config_p->served_tai.nb_tai - 1;
      config->tai_list.partial_tai_list[0]
        .u.tai_one_plmn_consecutive_tacs.mcc_digit1 =
        (mme_config_p->served_tai.plmn_mcc[0] / 100) % 10;
      config->tai_list.partial_tai_list[0]
        .u.tai_one_plmn_consecutive_tacs.mcc_digit2 =
        (mme_config_p->served_tai.plmn_mcc[0] / 10) % 10;
      config->tai_list.partial_tai_list[0]
        .u.tai_one_plmn_consecutive_tacs.mcc_digit3 =
        mme_config_p->served_tai.plmn_mcc[0] % 10;
      if (mme_config_p->served_tai.plmn_mnc_len[0] == 2) {
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_consecutive_tacs.mnc_digit1 =
          (mme_config_p->served_tai.plmn_mnc[0] / 10) % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_consecutive_tacs.mnc_digit2 =
          mme_config_p->served_tai.plmn_mnc[0] % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_consecutive_tacs.mnc_digit3 = 0xf;
      } else if (mme_config_p->served_tai.plmn_mnc_len[0] == 3) {
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_consecutive_tacs.mnc_digit1 =
          (mme_config_p->served_tai.plmn_mnc[0] / 100) % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_consecutive_tacs.mnc_digit2 =
          (mme_config_p->served_tai.plmn_mnc[0] / 10) % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_consecutive_tacs.mnc_digit3 =
          mme_config_p->served_tai.plmn_mnc[0] % 10;
      } else {
        AssertFatal(
          (mme_config_p->served_tai.plmn_mnc_len[0] >= 2) &&
            (mme_config_p->served_tai.plmn_mnc_len[0] <= 3),
          "BAD MNC length for GUMMEI");
      }
      config->tai_list.partial_tai_list[0].u.tai_one_plmn_consecutive_tacs.tac =
        mme_config_p->served_tai.tac[0];
      break;

    case TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS:
      config->tai_list.numberoflists = 1;
      config->tai_list.partial_tai_list[0].typeoflist =
        mme_config_p->served_tai.list_type;
      config->tai_list.partial_tai_list[0].numberofelements =
        mme_config_p->served_tai.nb_tai - 1;
      for (int i = 0; i < mme_config_p->served_tai.nb_tai; i++) {
        config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mcc_digit1 =
          (mme_config_p->served_tai.plmn_mcc[i] / 100) % 10;
        config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mcc_digit2 =
          (mme_config_p->served_tai.plmn_mcc[i] / 10) % 10;
        config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mcc_digit3 =
          mme_config_p->served_tai.plmn_mcc[i] % 10;
        if (mme_config_p->served_tai.plmn_mnc_len[0] == 2) {
          config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mnc_digit1 =
            (mme_config_p->served_tai.plmn_mnc[i] / 10) % 10;
          config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mnc_digit2 =
            mme_config_p->served_tai.plmn_mnc[i] % 10;
          config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mnc_digit3 =
            0xf;
        } else if (mme_config_p->served_tai.plmn_mnc_len[0] == 3) {
          config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mnc_digit1 =
            (mme_config_p->served_tai.plmn_mnc[i] / 100) % 10;
          config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mnc_digit2 =
            (mme_config_p->served_tai.plmn_mnc[i] / 10) % 10;
          config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].mnc_digit3 =
            mme_config_p->served_tai.plmn_mnc[i] % 10;
        } else {
          AssertFatal(
            (mme_config_p->served_tai.plmn_mnc_len[0] >= 2) &&
              (mme_config_p->served_tai.plmn_mnc_len[i] <= 3),
            "BAD MNC length for GUMMEI");
        }
        config->tai_list.partial_tai_list[0].u.tai_many_plmn[i].tac =
          mme_config_p->served_tai.tac[i];
        // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2 elements...), see 3GPP TS 24.301, section 9.9.3.33.1
      }
      break;

    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS:
      config->tai_list.numberoflists = 1;
      config->tai_list.partial_tai_list[0].typeoflist =
        mme_config_p->served_tai.list_type;
      config->tai_list.partial_tai_list[0]
        .u.tai_one_plmn_non_consecutive_tacs.mcc_digit1 =
        (mme_config_p->served_tai.plmn_mcc[0] / 100) % 10;
      config->tai_list.partial_tai_list[0]
        .u.tai_one_plmn_non_consecutive_tacs.mcc_digit2 =
        (mme_config_p->served_tai.plmn_mcc[0] / 10) % 10;
      config->tai_list.partial_tai_list[0]
        .u.tai_one_plmn_non_consecutive_tacs.mcc_digit3 =
        mme_config_p->served_tai.plmn_mcc[0] % 10;
      if (mme_config_p->served_tai.plmn_mnc_len[0] == 2) {
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.mnc_digit1 =
          (mme_config_p->served_tai.plmn_mnc[0] / 10) % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.mnc_digit2 =
          mme_config_p->served_tai.plmn_mnc[0] % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.mnc_digit3 = 0xf;
      } else if (mme_config_p->served_tai.plmn_mnc_len[0] == 3) {
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.mnc_digit1 =
          (mme_config_p->served_tai.plmn_mnc[0] / 100) % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.mnc_digit2 =
          (mme_config_p->served_tai.plmn_mnc[0] / 10) % 10;
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.mnc_digit3 =
          mme_config_p->served_tai.plmn_mnc[0] % 10;
      } else {
        AssertFatal(
          (mme_config_p->served_tai.plmn_mnc_len[0] >= 2) &&
            (mme_config_p->served_tai.plmn_mnc_len[0] <= 3),
          "BAD MNC length for GUMMEI");
      }
      for (int i = 0; i < mme_config_p->served_tai.nb_tai; i++) {
        config->tai_list.partial_tai_list[0]
          .u.tai_one_plmn_non_consecutive_tacs.tac[i] =
          mme_config_p->served_tai.tac[i];
      }
      // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2 elements...), see 3GPP TS 24.301, section 9.9.3.33.1
      config->tai_list.partial_tai_list[0].numberofelements =
        mme_config_p->served_tai.nb_tai - 1;
      break;
    default:
      AssertFatal(
        0,
        "BAD TAI list configuration, unknown TAI list type %u",
        mme_config_p->served_tai.list_type);
  }

  config->gummei = mme_config_p->gummei.gummei[0];

  // hardcoded
  config->eps_network_feature_support =
    EPS_NETWORK_FEATURE_SUPPORT_CS_LCS_LOCATION_SERVICES_VIA_CS_DOMAIN_NOT_SUPPORTED;
  if (
    mme_config_p->eps_network_feature_support
      .emergency_bearer_services_in_s1_mode != 0) {
    config->eps_network_feature_support |=
      EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE_SUPPORTED;
  }
  if (
    mme_config_p->eps_network_feature_support.ims_voice_over_ps_session_in_s1 !=
    0) {
    config->eps_network_feature_support |=
      EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1_SUPPORTED;
  }
  if (
    mme_config_p->eps_network_feature_support.location_services_via_epc != 0) {
    config->eps_network_feature_support |=
      EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC_SUPPORTED;
  }
  if (mme_config_p->eps_network_feature_support.extended_service_request != 0) {
    config->eps_network_feature_support |=
      EPS_NETWORK_FEATURE_SUPPORT_EXTENDED_SERVICE_REQUEST_SUPPORTED;
  }

  if (mme_config_p->unauthenticated_imsi_supported != 0) {
    config->features |= MME_API_UNAUTHENTICATED_IMSI;
  }

  for (int i = 0; i < 8; i++) {
    config->prefered_integrity_algorithm[i] =
      mme_config_p->nas_config.prefered_integrity_algorithm[i];
    config->prefered_ciphering_algorithm[i] =
      mme_config_p->nas_config.prefered_ciphering_algorithm[i];
  }

  config->full_network_name = bstrcpy(mme_config_p->full_network_name);
  config->short_network_name = bstrcpy(mme_config_p->short_network_name);
  config->daylight_saving_time = mme_config_p->daylight_saving_time;
  OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    mme_api_get_config()                                      **
 **                                                                        **
 ** Description: Retreives MME configuration data related to EPS session   **
 **      management                                                **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int mme_api_get_esm_config(mme_api_esm_config_t *config)
{
  OAILOG_FUNC_IN(LOG_NAS);

  if (strcmp((const char *) mme_config.ip_capability->data, "IPV4") == 0) {
    config->features = MME_API_IPV4;
  } else if (
    strcmp((const char *) mme_config.ip_capability->data, "IPV6") == 0) {
    config->features = MME_API_IPV6;
  } else if (
    strcmp((const char *) mme_config.ip_capability->data, "IPV4V6") == 0) {
    config->features = MME_API_IPV4 | MME_API_IPV6;
  } else {
    config->features = 0;
  }
  if (
    strcmp((const char *) mme_config.non_eps_service_control->data, "SMS") ==
    0) {
    config->features = config->features | MME_API_SMS_SUPPORTED;
  } else if (
    strcmp(
      (const char *) mme_config.non_eps_service_control->data, "CSFB_SMS") ==
    0) {
    config->features = config->features | MME_API_CSFB_SMS_SUPPORTED;
  }

  OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
}

/*
 *
 *  Name:    mme_api_notify_imsi()
 *
 *  Description: Notify the MME of the IMSI of a UE.
 *
 *  Inputs:
 *         ueid:      nas_ue id
 *         imsi64:    IMSI
 *  Return:    RETURNok, RETURNerror
 *
 */
int mme_api_notify_imsi(
  const mme_ue_s1ap_id_t id,
  const imsi64_t imsi64)
{
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context_t* ue_mm_context = NULL;

  OAILOG_FUNC_IN(LOG_NAS);
  ue_mm_context =
    mme_ue_context_exists_mme_ue_s1ap_id(&mme_app_desc_p->mme_ue_contexts, id);

  if (ue_mm_context) {
    mme_ue_context_update_coll_keys(
      &mme_app_desc_p->mme_ue_contexts,
      ue_mm_context,
      ue_mm_context->enb_s1ap_id_key,
      id,
      imsi64,
      ue_mm_context->mme_teid_s11,
      &ue_mm_context->emm_context._guti);
    unlock_ue_contexts(ue_mm_context);
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
}

/*
 *
 *  Name:    mme_api_notify_new_guti()
 *
 *  Description: Notify the MME of a generated GUTI for a UE(not spec).
 *
 *  Inputs:
 *         ueid:      nas_ue id
 *         guti:      EPS Globally Unique Temporary UE Identity
 *  Return:    RETURNok, RETURNerror
 *
 */
int mme_api_notify_new_guti(const mme_ue_s1ap_id_t id, guti_t* const guti)
{
  ue_mm_context_t* ue_mm_context = NULL;
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  OAILOG_FUNC_IN(LOG_NAS);
  ue_mm_context =
    mme_ue_context_exists_mme_ue_s1ap_id(&mme_app_desc_p->mme_ue_contexts, id);

  if (ue_mm_context) {
    mme_ue_context_update_coll_keys(
      &mme_app_desc_p->mme_ue_contexts,
      ue_mm_context,
      ue_mm_context->enb_s1ap_id_key,
      id,
      ue_mm_context->emm_context._imsi64,
      ue_mm_context->mme_teid_s11,
      guti);
    unlock_ue_contexts(ue_mm_context);
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    mme_api_new_guti()                                        **
 **                                                                        **
 ** Description: Requests the MME to assign a new GUTI to the UE identi-   **
 **      fied by the given IMSI.                                   **
 **                                                                        **
 ** Description: Requests the MME to assign a new GUTI to the UE identi-   **
 **      fied by the given IMSI and returns the list of consecu-   **
 **      tive tracking areas the UE is registered to.              **
 **                                                                        **
 ** Inputs:  imsi:      International Mobile Subscriber Identity   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     guti:      The new assigned GUTI                      **
 **      tai_list:       TAIs belonging to the PLMN                                **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int mme_api_new_guti(
  const imsi_t *const imsi,
  const guti_t *const old_guti,
  guti_t *const guti,
  const tai_t *const originating_tai,
  tai_list_t *const tai_list)
{
  OAILOG_FUNC_IN(LOG_NAS);
  mme_app_desc_t *mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context_t *ue_context = NULL;
  imsi64_t imsi64 = imsi_to_imsi64(imsi);

  ue_context =
    mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);

  if (ue_context) {
    guti->gummei.mme_gid = _emm_data.conf.gummei.mme_gid;
    guti->gummei.mme_code = _emm_data.conf.gummei.mme_code;
    guti->gummei.plmn.mcc_digit1 = _emm_data.conf.gummei.plmn.mcc_digit1;
    guti->gummei.plmn.mcc_digit2 = _emm_data.conf.gummei.plmn.mcc_digit2;
    guti->gummei.plmn.mcc_digit3 = _emm_data.conf.gummei.plmn.mcc_digit3;
    guti->gummei.plmn.mnc_digit1 = _emm_data.conf.gummei.plmn.mnc_digit1;
    guti->gummei.plmn.mnc_digit2 = _emm_data.conf.gummei.plmn.mnc_digit2;
    guti->gummei.plmn.mnc_digit3 = _emm_data.conf.gummei.plmn.mnc_digit3;
    // TODO Find another way to generate m_tmsi
    guti->m_tmsi = (tmsi_t)(uintptr_t) ue_context;
    if (guti->m_tmsi == INVALID_M_TMSI) {
      unlock_ue_contexts(ue_context);
      OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
    }
    mme_api_notify_new_guti(ue_context->mme_ue_s1ap_id, guti);
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }

  int j = 0;
  for (int i = 0; i < _emm_data.conf.tai_list.numberoflists; i++) {
    switch (_emm_data.conf.tai_list.partial_tai_list[i].typeoflist) {
      case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS:
        if (
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mcc_digit1 ==
           guti->gummei.plmn.mcc_digit1) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mcc_digit2 ==
           guti->gummei.plmn.mcc_digit2) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mcc_digit3 ==
           guti->gummei.plmn.mcc_digit3) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mnc_digit1 ==
           guti->gummei.plmn.mnc_digit1) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mnc_digit2 ==
           guti->gummei.plmn.mnc_digit2) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mnc_digit3 ==
           guti->gummei.plmn.mnc_digit3)) {
          tai_list->partial_tai_list[j].numberofelements =
            _emm_data.conf.tai_list.partial_tai_list[i].numberofelements;
          tai_list->partial_tai_list[j].typeoflist =
            _emm_data.conf.tai_list.partial_tai_list[i].typeoflist;

          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_non_consecutive_tacs.mcc_digit1 =
            guti->gummei.plmn.mcc_digit1;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_non_consecutive_tacs.mcc_digit2 =
            guti->gummei.plmn.mcc_digit2;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_non_consecutive_tacs.mcc_digit3 =
            guti->gummei.plmn.mcc_digit3;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_non_consecutive_tacs.mnc_digit1 =
            guti->gummei.plmn.mnc_digit1;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_non_consecutive_tacs.mnc_digit2 =
            guti->gummei.plmn.mnc_digit2;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_non_consecutive_tacs.mnc_digit3 =
            guti->gummei.plmn.mnc_digit3;
          // _emm_data.conf.tai_list is sorted
          for (int t = 0;
               t < (tai_list->partial_tai_list[j].numberofelements + 1);
               t++) {
            tai_list->partial_tai_list[j]
              .u.tai_one_plmn_non_consecutive_tacs.tac[t] =
              _emm_data.conf.tai_list.partial_tai_list[i]
                .u.tai_one_plmn_non_consecutive_tacs.tac[t];
          }
          j += 1;
        }
        break;
      case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS:
        if (
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_consecutive_tacs.mcc_digit1 ==
           guti->gummei.plmn.mcc_digit1) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_consecutive_tacs.mcc_digit2 ==
           guti->gummei.plmn.mcc_digit2) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_consecutive_tacs.mcc_digit3 ==
           guti->gummei.plmn.mcc_digit3) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_consecutive_tacs.mnc_digit1 ==
           guti->gummei.plmn.mnc_digit1) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_consecutive_tacs.mnc_digit2 ==
           guti->gummei.plmn.mnc_digit2) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_consecutive_tacs.mnc_digit3 ==
           guti->gummei.plmn.mnc_digit3)) {
          tai_list->partial_tai_list[j].numberofelements =
            _emm_data.conf.tai_list.partial_tai_list[i].numberofelements;
          tai_list->partial_tai_list[j].typeoflist =
            _emm_data.conf.tai_list.partial_tai_list[i].typeoflist;

          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_consecutive_tacs.mcc_digit1 =
            guti->gummei.plmn.mcc_digit1;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_consecutive_tacs.mcc_digit2 =
            guti->gummei.plmn.mcc_digit2;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_consecutive_tacs.mcc_digit3 =
            guti->gummei.plmn.mcc_digit3;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_consecutive_tacs.mnc_digit1 =
            guti->gummei.plmn.mnc_digit1;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_consecutive_tacs.mnc_digit2 =
            guti->gummei.plmn.mnc_digit2;
          tai_list->partial_tai_list[j]
            .u.tai_one_plmn_consecutive_tacs.mnc_digit3 =
            guti->gummei.plmn.mnc_digit3;
          // _emm_data.conf.tai_list is sorted
          tai_list->partial_tai_list[j].u.tai_one_plmn_consecutive_tacs.tac =
            _emm_data.conf.tai_list.partial_tai_list[i]
              .u.tai_one_plmn_consecutive_tacs.tac;
          j += 1;
        }
        break;
      case TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS:
        if (
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mcc_digit1 ==
           guti->gummei.plmn.mcc_digit1) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mcc_digit2 ==
           guti->gummei.plmn.mcc_digit2) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mcc_digit3 ==
           guti->gummei.plmn.mcc_digit3) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mnc_digit1 ==
           guti->gummei.plmn.mnc_digit1) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mnc_digit2 ==
           guti->gummei.plmn.mnc_digit2) &&
          (_emm_data.conf.tai_list.partial_tai_list[i]
             .u.tai_one_plmn_non_consecutive_tacs.mnc_digit3 ==
           guti->gummei.plmn.mnc_digit3)) {
          tai_list->partial_tai_list[j].numberofelements =
            _emm_data.conf.tai_list.partial_tai_list[i].numberofelements;
          tai_list->partial_tai_list[j].typeoflist =
            _emm_data.conf.tai_list.partial_tai_list[i].typeoflist;

          for (int t = 0;
               t < (tai_list->partial_tai_list[j].numberofelements + 1);
               t++) {
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].mcc_digit1 =
              _emm_data.conf.tai_list.partial_tai_list[i].u.tai_many_plmn[t]
              .mcc_digit1;
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].mcc_digit2 =
              _emm_data.conf.tai_list.partial_tai_list[i].u.tai_many_plmn[t]
              .mcc_digit2;
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].mcc_digit3 =
              _emm_data.conf.tai_list.partial_tai_list[i].u.tai_many_plmn[t]
              .mcc_digit3;
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].mnc_digit1 =
              _emm_data.conf.tai_list.partial_tai_list[i].u.tai_many_plmn[t]
             .mnc_digit1;
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].mnc_digit2 =
              _emm_data.conf.tai_list.partial_tai_list[i].u.tai_many_plmn[t]
              .mnc_digit2;
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].mnc_digit3 =
              _emm_data.conf.tai_list.partial_tai_list[i].u.tai_many_plmn[t]
              .mnc_digit3;
            // _emm_data.conf.tai_list is sorted
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].tac =
              _emm_data.conf.tai_list.partial_tai_list[i]
                .u.tai_many_plmn[t]
                .tac;
          }
          j += 1;
        }
        break;
      default:
        AssertFatal(
          0,
          "BAD TAI list configuration, unknown TAI list type %u",
          _emm_data.conf.tai_list.partial_tai_list[i].typeoflist);
    }
  }
  tai_list->numberoflists = j;
  OAILOG_INFO(
    LOG_NAS,
    "UE " MME_UE_S1AP_ID_FMT "  Got GUTI " GUTI_FMT "\n",
    ue_context->mme_ue_s1ap_id,
    GUTI_ARG(guti));
  unlock_ue_contexts(ue_context);
  OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:        mme_api_subscribe()                                       **
 **                                                                        **
 ** Description: Requests the MME to check whether connectivity with the   **
 **              requested PDN can be established using the specified APN. **
 **              If accepted the MME returns PDN subscription context con- **
 **              taining EPS subscribed QoS profile, the default APN if    **
 **              required and UE's IPv4 address and/or the IPv6 prefix.    **
 **                                                                        **
 ** Inputs:  apn:               If not NULL, Access Point Name of the PDN  **
 **                             to connect to                              **
 **              is_emergency:  true if the PDN connectivity is requested  **
 **                             for emergency bearer services              **
 **                  Others:    None                                       **
 **                                                                        **
 ** Outputs:         apn:       If NULL, default APN or APN configured for **
 **                             emergency bearer services                  **
 **                  pdn_addr:  PDN connection IPv4 address or IPv6 inter- **
 **                             face identifier to be used to build the    **
 **                             IPv6 link local address                    **
 **                  qos:       EPS subscribed QoS profile                 **
 **                  Return:    RETURNok, RETURNerror                      **
 **                  Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int mme_api_subscribe(
  bstring *apn,
  mme_api_ip_version_t mme_pdn_index,
  bstring *pdn_addr,
  int is_emergency,
  mme_api_qos_t *qos)
{
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_NAS);
  OAILOG_FUNC_RETURN(LOG_NAS, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:        mme_api_unsubscribe()                                     **
 **                                                                        **
 ** Description: Requests the MME to release connectivity with the reques- **
 **              ted PDN using the specified APN.                          **
 **                                                                        **
 ** Inputs:  apn:               Access Point Name of the PDN to disconnect **
 **                             from                                       **
 **                  Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **                  Return:    RETURNok, RETURNerror                      **
 **                  Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int mme_api_unsubscribe(bstring apn)
{
  OAILOG_FUNC_IN(LOG_NAS);
  int rc = RETURNok;

  /*
   * Decrement the total number of PDN connections
   */
  _mme_api_pdn_id -= 1;
  OAILOG_FUNC_RETURN(LOG_NAS, rc);
}
