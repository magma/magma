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
static int mme_api_pdn_id = 0;

static tmsi_t generate_random_TMSI(void);

static int copy_plmn_from_config(
    const served_tai_t* served_tai, int index, plmn_t* plmn);

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
    mme_api_emm_config_t* config, const struct mme_config_s* mme_config_p) {
  OAILOG_FUNC_IN(LOG_NAS);
  if (mme_config_p->served_tai.nb_tai < 1) {
    OAILOG_ERROR(LOG_NAS, "No TAI configured\n");
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }
  OAILOG_INFO(
      LOG_NAS, "Number of GUMMEIs supported = %d\n", mme_config_p->gummei.nb);

  config->tai_list.numberoflists = 0;
  // We can store 16 TAIs per list and have 16 partial lists maximum
  // As per TS 124.301 V15.4.0 Section 9.9.3.33, we will be sending at most
  // 16 TAIs for which UE data is set during
  switch (mme_config_p->served_tai.list_type) {
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS: {
      int tai_list_i = 0, tac_i = 0;
      for (int i = 0; i < mme_config_p->served_tai.nb_tai; ++i) {
        if (tac_i == 16) {  // cannot have more than 16 TACs per partial list
          // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2
          // elements...), see 3GPP TS 24.301, section 9.9.3.33.1
          config->tai_list.partial_tai_list[tai_list_i].numberofelements =
              tac_i - 1;
          ++tai_list_i;
          tac_i = 0;
        }
        config->tai_list.partial_tai_list[tai_list_i].typeoflist =
            mme_config_p->served_tai.list_type;
        int rc = copy_plmn_from_config(
            &(mme_config_p->served_tai), i,
            &(config->tai_list.partial_tai_list[tai_list_i]
                  .u.tai_one_plmn_consecutive_tacs.plmn));
        if (rc != RETURNok) {
          OAILOG_ERROR(LOG_NAS, "BAD MNC length in SERVED TAI\n");
          OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
        }
        config->tai_list.partial_tai_list[tai_list_i]
            .u.tai_one_plmn_consecutive_tacs.tac =
            mme_config_p->served_tai.tac[16 * tai_list_i];
        ++tac_i;
      }
      if (tai_list_i >= TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI) {
        OAILOG_ERROR(LOG_NAS, "Too many TAI partial list in TAI list\n");
        OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
      }
      // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2
      // elements...), see 3GPP TS 24.301, section 9.9.3.33.1
      config->tai_list.partial_tai_list[tai_list_i].numberofelements =
          tac_i - 1;
      config->tai_list.numberoflists = tai_list_i + 1;

      break;
    }
    case TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS: {
      // no need to optimize here as we do not really support multi PLMN
      // use case in Magma yet
      int tai_list_i = 0, tac_i = 0;
      uint16_t last_mcc = mme_config_p->served_tai.plmn_mcc[0];
      uint16_t last_mnc = mme_config_p->served_tai.plmn_mnc[0];
      for (int i = 0; i < mme_config_p->served_tai.nb_tai; i++) {
        bool is_plmn_changed = false;
        if ((mme_config_p->served_tai.plmn_mcc[i] != last_mcc) ||
            (mme_config_p->served_tai.plmn_mnc[i] != last_mnc)) {
          last_mcc        = mme_config_p->served_tai.plmn_mcc[i];
          last_mnc        = mme_config_p->served_tai.plmn_mnc[i];
          is_plmn_changed = true;
        }
        if ((tac_i == 16) || is_plmn_changed) {  // cannot have more than 16
                                                 // TACs per partial list
          // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2
          // elements...), see 3GPP TS 24.301, section 9.9.3.33.1
          config->tai_list.partial_tai_list[tai_list_i].numberofelements =
              tac_i - 1;
          ++tai_list_i;
          tac_i = 0;
        }
        config->tai_list.partial_tai_list[tai_list_i].typeoflist =
            mme_config_p->served_tai.list_type;
        int rc = copy_plmn_from_config(
            &(mme_config_p->served_tai), i,
            &(config->tai_list.partial_tai_list[tai_list_i]
                  .u.tai_many_plmn[tac_i]
                  .plmn));
        if (rc != RETURNok) {
          OAILOG_ERROR(LOG_NAS, "BAD MNC length in SERVED TAI\n");
          OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
        }
        config->tai_list.partial_tai_list[tai_list_i].u.tai_many_plmn[i].tac =
            mme_config_p->served_tai.tac[i];
        ++tac_i;
      }
      if (tai_list_i >= TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI) {
        OAILOG_ERROR(LOG_NAS, "Too many TAI partial list in TAI list\n");
        OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
      }
      // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2
      // elements...), see 3GPP TS 24.301, section 9.9.3.33.1
      config->tai_list.partial_tai_list[tai_list_i].numberofelements =
          tac_i - 1;
      config->tai_list.numberoflists = tai_list_i + 1;
      break;
    }
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS: {
      int tai_list_i = 0, tac_i = 0;
      for (int i = 0; i < mme_config_p->served_tai.nb_tai; ++i) {
        if (tac_i == 16) {  // cannot have more than 16 TACs per partial list
          // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2
          // elements...), see 3GPP TS 24.301, section 9.9.3.33.1
          config->tai_list.partial_tai_list[tai_list_i].numberofelements =
              tac_i - 1;
          ++tai_list_i;
          tac_i = 0;
        }
        config->tai_list.partial_tai_list[tai_list_i].typeoflist =
            mme_config_p->served_tai.list_type;
        int rc = copy_plmn_from_config(
            &(mme_config_p->served_tai), i,
            &(config->tai_list.partial_tai_list[tai_list_i]
                  .u.tai_one_plmn_non_consecutive_tacs.plmn));
        if (rc != RETURNok) {
          OAILOG_ERROR(LOG_NAS, "BAD MNC length in SERVED TAI\n");
          OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
        }
        config->tai_list.partial_tai_list[tai_list_i]
            .u.tai_one_plmn_non_consecutive_tacs.tac[tac_i] =
            mme_config_p->served_tai.tac[i];
        ++tac_i;
      }
      if (tai_list_i >= TRACKING_AREA_IDENTITY_LIST_MAXIMUM_NUM_TAI) {
        OAILOG_ERROR(LOG_NAS, "Too many TAI partial list in TAI list\n");
        OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
      }
      // LW: number of elements is coded as N-1 (0 -> 1 element, 1 -> 2
      // elements...), see 3GPP TS 24.301, section 9.9.3.33.1
      config->tai_list.partial_tai_list[tai_list_i].numberofelements =
          tac_i - 1;
      config->tai_list.numberoflists = tai_list_i + 1;
      break;
    }
    default:
      OAILOG_ERROR(
          LOG_NAS, "BAD TAI list configuration, unknown TAI list type %u\n",
          mme_config_p->served_tai.list_type);
      OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }
  // Read GUMMEI List
  config->gummei.num_gummei = mme_config_p->gummei.nb;
  for (uint8_t num_gummei = 0; num_gummei < mme_config_p->gummei.nb;
       num_gummei++) {
    config->gummei.gummei[num_gummei] = mme_config_p->gummei.gummei[num_gummei];
  }

  // hardcoded
  config->eps_network_feature_support[0] =
      EPS_NETWORK_FEATURE_SUPPORT_CS_LCS_LOCATION_SERVICES_VIA_CS_DOMAIN_NOT_SUPPORTED;
  if (mme_config_p->eps_network_feature_support
          .emergency_bearer_services_in_s1_mode != 0) {
    config->eps_network_feature_support[0] |=
        EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE_SUPPORTED;
  }
  if (mme_config_p->eps_network_feature_support
          .ims_voice_over_ps_session_in_s1 != 0) {
    config->eps_network_feature_support[0] |=
        EPS_NETWORK_FEATURE_SUPPORT_IMS_VOICE_OVER_PS_SESSION_IN_S1_SUPPORTED;
  }
  if (mme_config_p->eps_network_feature_support.location_services_via_epc !=
      0) {
    config->eps_network_feature_support[0] |=
        EPS_NETWORK_FEATURE_SUPPORT_LOCATION_SERVICES_VIA_EPC_SUPPORTED;
  }
  if (mme_config_p->eps_network_feature_support.extended_service_request != 0) {
    config->eps_network_feature_support[0] |=
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

  config->full_network_name    = bstrcpy(mme_config_p->full_network_name);
  config->short_network_name   = bstrcpy(mme_config_p->short_network_name);
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
int mme_api_get_esm_config(mme_api_esm_config_t* config) {
  OAILOG_FUNC_IN(LOG_NAS);
  if (strcmp((const char*) mme_config.non_eps_service_control->data, "SMS") ==
      0) {
    config->features = config->features | MME_API_SMS_SUPPORTED;
  } else if (
      strcmp(
          (const char*) mme_config.non_eps_service_control->data, "CSFB_SMS") ==
      0) {
    config->features = config->features | MME_API_CSFB_SMS_SUPPORTED;
  } else if (
      strcmp(
          (const char*) mme_config.non_eps_service_control->data,
          "SMS_ORC8R") == 0) {
    config->features = config->features | MME_API_SMS_ORC8R_SUPPORTED;
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
int mme_api_notify_imsi(const mme_ue_s1ap_id_t id, imsi64_t imsi64) {
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context_t* ue_mm_context = NULL;

  OAILOG_FUNC_IN(LOG_NAS);
  ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(id);

  if (ue_mm_context) {
    mme_ue_context_update_coll_keys(
        &mme_app_desc_p->mme_ue_contexts, ue_mm_context,
        ue_mm_context->enb_s1ap_id_key, id, imsi64, ue_mm_context->mme_teid_s11,
        &ue_mm_context->emm_context._guti);
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
int mme_api_notify_new_guti(const mme_ue_s1ap_id_t id, guti_t* const guti) {
  ue_mm_context_t* ue_mm_context = NULL;
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  OAILOG_FUNC_IN(LOG_NAS);
  ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(id);

  if (ue_mm_context) {
    mme_ue_context_update_coll_keys(
        &mme_app_desc_p->mme_ue_contexts, ue_mm_context,
        ue_mm_context->enb_s1ap_id_key, id, ue_mm_context->emm_context._imsi64,
        ue_mm_context->mme_teid_s11, guti);
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
}

/************************************************************************
 **                                                                    **
 ** Name:    mme_api_new_guti()                                        **
 **                                                                    **
 ** Description: Requests the MME to assign a new GUTI to the UE       **
 **      identified by the given IMSI.                                 **
 **                                                                    **
 ** Description: Requests the MME to assign a new GUTI to the UE       **
 **      identified by the given IMSI and returns the list of          **
 **      consecutive tracking areas the UE is registered to.           **
 **                                                                    **
 ** Inputs:  imsi:      International Mobile Subscriber Identity       **
 **      Others:    None                                               **
 **                                                                    **
 ** Outputs:     guti:      The new assigned GUTI                      **
 **      tai_list:       TAIs belonging to the PLMN                    **
 **      Return:    RETURNok, RETURNerror                              **
 **      Others:    None                                               **
 **                                                                    **
 ***********************************************************************/
int mme_api_new_guti(
    const imsi_t* const imsi, const guti_t* const old_guti, guti_t* const guti,
    const tai_t* const originating_tai, tai_list_t* const tai_list) {
  OAILOG_FUNC_IN(LOG_NAS);
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context_t* ue_context    = NULL;
  imsi64_t imsi64                = imsi_to_imsi64(imsi);
  bool is_plmn_equal             = false;

  ue_context =
      mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);

  if (ue_context) {
    for (uint8_t nb_gummei = 0; nb_gummei < _emm_data.conf.gummei.num_gummei;
         nb_gummei++) {
      /* comparing UE serving cell plmn with the gummei list in
       * mme configuration. */
      if (IS_PLMN_EQUAL(
              ue_context->emm_context.originating_tai.plmn,
              mme_config.gummei.gummei[nb_gummei].plmn)) {
        is_plmn_equal = true;
        /* Copies the GUMMEI value from configuration to the emm context */
        COPY_GUMMEI(guti, _emm_data.conf.gummei.gummei[nb_gummei]);
        break;
      }
    }
    if (!is_plmn_equal) {
      OAILOG_ERROR(LOG_NAS, "Serving PLMN not matching with GUMMEI List!\n");
      OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
    }

    // TODO Find another way to generate m_tmsi
    guti->m_tmsi = generate_random_TMSI();
    if (guti->m_tmsi == INVALID_M_TMSI) {
      OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
    }
    mme_api_notify_new_guti(ue_context->mme_ue_s1ap_id, guti);
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }

  int j = 0;  // keeps track of number of partial lists with matching PLMN
  for (int i = 0; i < _emm_data.conf.tai_list.numberoflists; i++) {
    /* Comparing mme configuration plmn of TAI_LIST with plmn of GUMMEI_LIST
     * if PLMN matches, _emm_data.conf.tai_list value gets updated with TAI_LIST
     * values from mme configuration file
     */
    switch (_emm_data.conf.tai_list.partial_tai_list[i].typeoflist) {
      case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS:
        if (IS_PLMN_EQUAL(
                _emm_data.conf.tai_list.partial_tai_list[i]
                    .u.tai_one_plmn_non_consecutive_tacs.plmn,
                guti->gummei.plmn)) {
          tai_list->partial_tai_list[j].numberofelements =
              _emm_data.conf.tai_list.partial_tai_list[i].numberofelements;
          tai_list->partial_tai_list[j].typeoflist =
              _emm_data.conf.tai_list.partial_tai_list[i].typeoflist;
          COPY_PLMN(
              tai_list->partial_tai_list[j]
                  .u.tai_one_plmn_non_consecutive_tacs.plmn,
              guti->gummei.plmn);

          // _emm_data.conf.tai_list is sorted
          for (int t = 0;
               t < (tai_list->partial_tai_list[j].numberofelements + 1); t++) {
            tai_list->partial_tai_list[j]
                .u.tai_one_plmn_non_consecutive_tacs.tac[t] =
                _emm_data.conf.tai_list.partial_tai_list[i]
                    .u.tai_one_plmn_non_consecutive_tacs.tac[t];
          }
          j += 1;
        } else {
          OAILOG_ERROR(
              LOG_NAS,
              "GUTI PLMN does not match with mme configuration tai list\n");
        }
        break;
      case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS:
        if (IS_PLMN_EQUAL(
                _emm_data.conf.tai_list.partial_tai_list[i]
                    .u.tai_one_plmn_consecutive_tacs.plmn,
                guti->gummei.plmn)) {
          tai_list->partial_tai_list[j].numberofelements =
              _emm_data.conf.tai_list.partial_tai_list[i].numberofelements;
          tai_list->partial_tai_list[j].typeoflist =
              _emm_data.conf.tai_list.partial_tai_list[i].typeoflist;

          COPY_PLMN(
              tai_list->partial_tai_list[j]
                  .u.tai_one_plmn_consecutive_tacs.plmn,
              guti->gummei.plmn);

          // _emm_data.conf.tai_list is sorted
          tai_list->partial_tai_list[j].u.tai_one_plmn_consecutive_tacs.tac =
              _emm_data.conf.tai_list.partial_tai_list[i]
                  .u.tai_one_plmn_consecutive_tacs.tac;
          j += 1;
        } else {
          OAILOG_ERROR(
              LOG_NAS,
              "GUTI PLMN does not match with mme configuration tai list\n");
        }
        break;
      case TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS:
        // Each partial list has the same plmn
        if (IS_PLMN_EQUAL(
                _emm_data.conf.tai_list.partial_tai_list[i]
                    .u.tai_many_plmn[0]
                    .plmn,
                guti->gummei.plmn)) {
          tai_list->partial_tai_list[j].numberofelements =
              _emm_data.conf.tai_list.partial_tai_list[i].numberofelements;
          tai_list->partial_tai_list[j].typeoflist =
              _emm_data.conf.tai_list.partial_tai_list[i].typeoflist;

          for (int t = 0;
               t < (tai_list->partial_tai_list[j].numberofelements + 1); t++) {
            COPY_PLMN(
                tai_list->partial_tai_list[j].u.tai_many_plmn[t].plmn,
                _emm_data.conf.tai_list.partial_tai_list[i]
                    .u.tai_many_plmn[t]
                    .plmn);

            // _emm_data.conf.tai_list is sorted
            tai_list->partial_tai_list[j].u.tai_many_plmn[t].tac =
                _emm_data.conf.tai_list.partial_tai_list[i]
                    .u.tai_many_plmn[t]
                    .tac;
          }
          j += 1;
        } else {
          OAILOG_ERROR(
              LOG_NAS,
              "GUTI PLMN does not match with mme configuration tai list\n");
        }
        break;
      default:
        Fatal(
            "BAD TAI list configuration, unknown TAI list type %u",
            _emm_data.conf.tai_list.partial_tai_list[i].typeoflist);
    }

    // TS 124.301 V15.4.0 Section 9.9.3.33:
    // "The Tracking area identity list is a type 4 information element,
    // with a minimum length of 8 octets and a maximum length of 98 octets.
    // The list can contain a maximum of 16 different tracking area identities."
    // We will limit the number to 1 partial list which can have maximum of 16
    // TAIs.
    if (j == 1) {
      break;  // for loop
    }
  }
  tai_list->numberoflists = j;
  OAILOG_INFO(
      LOG_NAS,
      "UE " MME_UE_S1AP_ID_FMT "  Got GUTI " GUTI_FMT
      ". The number of TAI partial lists: %d",
      ue_context->mme_ue_s1ap_id, GUTI_ARG(guti), tai_list->numberoflists);
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
    bstring* apn, mme_api_ip_version_t mme_pdn_index, bstring* pdn_addr,
    int is_emergency, mme_api_qos_t* qos) {
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
int mme_api_unsubscribe(bstring apn) {
  OAILOG_FUNC_IN(LOG_NAS);
  int rc = RETURNok;

  /*
   * Decrement the total number of PDN connections
   */
  mme_api_pdn_id -= 1;
  OAILOG_FUNC_RETURN(LOG_NAS, rc);
}

static tmsi_t generate_random_TMSI() {
  // note srand with seed is initialized at main
  return (tmsi_t) rand();
}

/****************************************************************************
 **                                                                        **
 ** Name:        copy_plmn_from_config()                                   **
 **                                                                        **
 ** Description: Copies the tai list configuration to partial tai list.    **
 **                                                                        **
 ** Inputs:  served_tai:        Served tai constructed from MME config.    **
 **          index:             Index to used on served_tai                **
 **                  Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **          plmn:   PLMN filled from served_tai                           **
 **                  Return:    RETURNok, RETURNerror                      **
 **                  Others:                                               **
 **                                                                        **
 ***************************************************************************/
static int copy_plmn_from_config(
    const served_tai_t* served_tai, int index, plmn_t* plmn) {
  plmn->mcc_digit1 = (served_tai->plmn_mcc[index] / 100) % 10;
  plmn->mcc_digit2 = (served_tai->plmn_mcc[index] / 10) % 10;
  plmn->mcc_digit3 = served_tai->plmn_mcc[index] % 10;
  if (served_tai->plmn_mnc_len[index] == 2) {
    plmn->mnc_digit1 = (served_tai->plmn_mnc[index] / 10) % 10;
    plmn->mnc_digit2 = served_tai->plmn_mnc[index] % 10;
    plmn->mnc_digit3 = 0xf;
  } else if (served_tai->plmn_mnc_len[index] == 3) {
    plmn->mnc_digit1 = (served_tai->plmn_mnc[index] / 100) % 10;
    plmn->mnc_digit2 = (served_tai->plmn_mnc[index] / 10) % 10;
    plmn->mnc_digit3 = served_tai->plmn_mnc[index] % 10;
  } else {
    return RETURNerror;
  }
  return RETURNok;
}
