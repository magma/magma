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
#include <stdint.h>
#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <time.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "common_defs.h"
#include "3gpp_24.008.h"
#include "3gpp_24.301.h"
#include "emm_msgDef.h"
#include "emm_proc.h"
#include "mme_config.h"
#include "emm_send.h"
#include "emm_data.h"
#include "mme_app_ue_context.h"
#include "conversions.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "AdditionalUpdateResult.h"
#include "AdditionalUpdateType.h"
#include "Cli.h"
#include "DetachType.h"
#include "EmmCause.h"
#include "EpsAttachResult.h"
#include "EpsBearerContextStatus.h"
#include "EpsMobileIdentity.h"
#include "EpsNetworkFeatureSupport.h"
#include "EpsUpdateResult.h"
#include "EsmMessageContainer.h"
#include "MobileIdentity.h"
#include "NasKeySetIdentifier.h"
#include "NasMessageContainer.h"
#include "NasSecurityAlgorithms.h"
#include "PagingIdentity.h"
#include "TrackingAreaIdentityList.h"
#include "esm_data.h"
#include "mme_api.h"
#include "nas/securityDef.h"

#define MAX_MINUTE_DIGITS 3 /* Maximum digits required to hold minute value */

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Functions executed by both the UE and the MME to send EMM messages
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_status()                                         **
 **                                                                        **
 ** Description: Builds EMM status message                                 **
 **                                                                        **
 **      The EMM status message is sent by the UE or the network   **
 **      at any time to report certain error conditions.           **
 **                                                                        **
 ** Inputs:  emm_cause: EMM cause code                             **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_status(const emm_as_status_t* msg, emm_status_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_WARNING(
      LOG_NAS_EMM, "EMMAS-SAP - Send EMM Status message (cause=%d)\n",
      msg->emm_cause);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = EMM_STATUS;
  /*
   * Mandatory - EMM cause
   */
  size += EMM_CAUSE_MAXIMUM_LENGTH;
  emm_msg->emmcause = msg->emm_cause;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_detach_accept()                                  **
 **                                                                        **
 ** Description: Builds Detach Accept message                              **
 **                                                                        **
 **      The Detach Accept message is sent by the UE or the net-   **
 **      work to indicate that the detach procedure has been com-  **
 **      pleted.                                                   **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_detach_accept(
    const emm_as_data_t* msg, detach_accept_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send Detach Accept message\n");
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = DETACH_ACCEPT;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}
/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_nw_detach_request()                                  **
 **                                                                        **
 ** Description: Builds Detach Accept message                              **
 **                                                                        **
 **      The Detach Accept message is sent by the n/w                      **
 **      to initiate detach procedure.                                     **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                     **
 **      Return:    The size of the EMM message                            **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int emm_send_nw_detach_request(
    const emm_as_data_t* msg, nw_detach_request_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send Detach Request message\n");
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = NW_DETACH_REQUEST;
  /*
   * Mandatory Detach type
   */
  size += DETACH_TYPE_MAXIMUM_LENGTH;

  emm_msg->nw_detachtype = msg->type;
  /*
   * Optional EMM cause. Not present
   */
  emm_msg->presenceMask = 0;

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/*
   --------------------------------------------------------------------------
   Functions executed by the MME to send EMM messages to the UE
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_attach_accept()                                  **
 **                                                                        **
 ** Description: Builds Attach Accept message                              **
 **                                                                        **
 **      The Attach Accept message is sent by the network to the   **
 **      UE to indicate that the corresponding attach request has  **
 **      been accepted.                                            **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_attach_accept(
    const emm_as_establish_t* msg, attach_accept_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  // Get the UE context
  emm_context_t* emm_ctx = emm_context_get(&_emm_data, msg->ue_id);
  DevAssert(emm_ctx);
  ue_mm_context_t* ue_mm_context_p =
      PARENT_STRUCT(emm_ctx, struct ue_mm_context_s, emm_context);
  mme_ue_s1ap_id_t ue_id = ue_mm_context_p->mme_ue_s1ap_id;
  DevAssert(msg->ue_id == ue_id);

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send Attach Accept message, ue_id = " MME_UE_S1AP_ID_FMT
      "\n",
      msg->ue_id);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - size = EMM_HEADER_MAXIMUM_LENGTH(%d)\n", size);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = ATTACH_ACCEPT;
  /*
   * Mandatory - EPS attach result
   */
  size += EPS_ATTACH_RESULT_MAXIMUM_LENGTH;
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += EPS_ATTACH_RESULT_MAXIMUM_LENGTH(%d)  (%d)\n",
      EPS_ATTACH_RESULT_MAXIMUM_LENGTH, size);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - EMM Context Attach Type (%d) for (ue_id = %u)\n",
      emm_ctx->attach_type, ue_id);
  switch (emm_ctx->attach_type) {
    case EMM_ATTACH_TYPE_COMBINED_EPS_IMSI:
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "EMMAS-SAP - Combined EPS/IMSI attach for (ue_id = %u)\n", ue_id);
      /* It is observed that UE/Handest (with usage setting = voice centric and
       * voice domain preference = CS voice only) sends detach after successful
       * attach .UEs with such settings sends attach type = combined EPS and
       * IMSI attach as attach_type in attach request message. At present, EPC
       * does not support interface with MSC /CS domain and it supports only LTE
       * data service, hence it is supposed to send attach_result as EPS-only
       * and add emm_cause = "CS domain not available" for such cases. Ideally
       * in data service only n/w ,UE's usage setting should be set to data
       * centric mode and should send attach type as EPS attach only. However UE
       * settings may not be in our control. To take care of this as a
       * workaround in this patch we modified MME implementation to set EPS
       * result to Combined attach if attach type is combined attach to prevent
       * such UEs from sending detach so that such UEs can remain attached in
       * the n/w and should be able to get data service from the n/w.
       */

      /* Added check for CSFB. If Location Update procedure towards MSC/VLR
       * fails or Network Access mode received as PACKET_ONLY from HSS in ULS
       * message send epsattachresult to EPS_ATTACH_RESULT_EPS. If is it
       * successful set epsattachresult to EPS_ATTACH_RESULT_EPS_IMSI
       */
      if (((_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED) ||
           (_esm_data.conf.features & MME_API_SMS_SUPPORTED)) &&
          ((emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) ||
           is_mme_ue_context_network_access_mode_packet_only(
               ue_mm_context_p))) {
        emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS;
      } else {
        emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS_IMSI;
      }
      break;
    case EMM_ATTACH_TYPE_RESERVED:
    default:
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "EMMAS-SAP - Unused attach type defaults to EPS attach for (ue_id = "
          "%u)\n",
          ue_id);
    case EMM_ATTACH_TYPE_EPS:
      emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS;
      OAILOG_DEBUG(
          LOG_NAS_EMM, "EMMAS-SAP - EPS attach for (ue_id = %u)\n", ue_id);
      break;
    case EMM_ATTACH_TYPE_EMERGENCY:  // We should not reach here
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMMAS-SAP - EPS emergency attach, currently unsupported for (ue_id "
          "= " MME_UE_S1AP_ID_FMT ")\n",
          ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);  // TODO: fix once supported
      break;
  }
  /*
   * Mandatory - T3412 value
   */
  size += GPRS_TIMER_IE_MAX_LENGTH;
  // Check whether Periodic TAU timer is disabled
  if (mme_config.nas_config.t3412_min == 0) {
    emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_0S;
    emm_msg->t3412value.timervalue = mme_config.nas_config.t3412_min;
  } else if (mme_config.nas_config.t3412_min <= 31) {
    emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_60S;
    emm_msg->t3412value.timervalue = mme_config.nas_config.t3412_min;
  } else {
    emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_360S;
    emm_msg->t3412value.timervalue = mme_config.nas_config.t3412_min / 6;
  }
  // emm_msg->t3412value.unit = GPRS_TIMER_UNIT_0S;
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += GPRS_TIMER_IE_MAX_LENGTH(%d)  (%d) for (ue_id = "
      "%u)\n",
      GPRS_TIMER_IE_MAX_LENGTH, size, ue_id);
  /*
   * Mandatory - Tracking area identity list
   */
  size +=
      TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH * msg->tai_list.numberoflists;
  memcpy(&emm_msg->tailist, &msg->tai_list, sizeof(msg->tai_list));
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += "
      "TRACKING_AREA_IDENTITY_LIST_LENGTH(%d*%d)  (%d) for (ue_id = %u)\n",
      TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH,
      emm_msg->tailist.numberoflists, size, ue_id);

  for (int p = 0; p < emm_msg->tailist.numberoflists; p++) {
    if (TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS ==
        emm_msg->tailist.partial_tai_list[p].typeoflist) {
      size = size + (2 * emm_msg->tailist.partial_tai_list[p].numberofelements);
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "EMMAS-SAP - size += "
          "TRACKING AREA CODE LENGTH(%d*%d)  (%d) for (ue_id = %u)\n",
          2, emm_msg->tailist.partial_tai_list[p].numberofelements, size,
          ue_id);
    } else if (
        TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS ==
        emm_msg->tailist.partial_tai_list[p].typeoflist) {
      size = size + (5 * emm_msg->tailist.partial_tai_list[p].numberofelements);
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "EMMAS-SAP - size += "
          "TRACKING AREA CODE LENGTH(%d*%d)  (%d) for (ue_id = %u)\n",
          5, emm_msg->tailist.partial_tai_list[p].numberofelements, size,
          ue_id);
    }
  }
  /*
   * Mandatory - ESM message container
   */
  size += ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH + blength(msg->nas_msg);
  emm_msg->esmmessagecontainer = bstrcpy(msg->nas_msg);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += "
      "ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH(%d)  (%d) for (ue_id = %u)\n",
      ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH, size, ue_id);

  /*
   * Optional - GUTI
   */
  if (msg->new_guti) {
    size += EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_GUTI_PRESENT;
    emm_msg->guti.guti.typeofidentity = EPS_MOBILE_IDENTITY_GUTI;
    emm_msg->guti.guti.oddeven        = EPS_MOBILE_IDENTITY_EVEN;
    emm_msg->guti.guti.mme_group_id   = msg->new_guti->gummei.mme_gid;
    emm_msg->guti.guti.mme_code       = msg->new_guti->gummei.mme_code;
    emm_msg->guti.guti.m_tmsi         = msg->new_guti->m_tmsi;
    emm_msg->guti.guti.mcc_digit1     = msg->new_guti->gummei.plmn.mcc_digit1;
    emm_msg->guti.guti.mcc_digit2     = msg->new_guti->gummei.plmn.mcc_digit2;
    emm_msg->guti.guti.mcc_digit3     = msg->new_guti->gummei.plmn.mcc_digit3;
    emm_msg->guti.guti.mnc_digit1     = msg->new_guti->gummei.plmn.mnc_digit1;
    emm_msg->guti.guti.mnc_digit2     = msg->new_guti->gummei.plmn.mnc_digit2;
    emm_msg->guti.guti.mnc_digit3     = msg->new_guti->gummei.plmn.mnc_digit3;
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH(%d)  (%d) for (ue_id = %u)\n",
        EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH, size, ue_id);
  }

  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - Out of GUTI for (ue_id = %u)\n", ue_id);
  /*
   * Optional - LAI
   */
  if (msg->location_area_identification) {
    size += LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT;
    emm_msg->locationareaidentification.mccdigit2 =
        msg->location_area_identification->mccdigit2;
    emm_msg->locationareaidentification.mccdigit1 =
        msg->location_area_identification->mccdigit1;
    emm_msg->locationareaidentification.mncdigit3 =
        msg->location_area_identification->mncdigit3;
    emm_msg->locationareaidentification.mccdigit3 =
        msg->location_area_identification->mccdigit3;
    emm_msg->locationareaidentification.mncdigit2 =
        msg->location_area_identification->mncdigit2;
    emm_msg->locationareaidentification.mncdigit1 =
        msg->location_area_identification->mncdigit1;
    emm_msg->locationareaidentification.lac =
        msg->location_area_identification->lac;
  }

  OAILOG_DEBUG(LOG_NAS_EMM, "EMMAS-SAP - Out of LAI for (ue_id = %u)\n", ue_id);
  /*
   * Optional - Mobile Identity
   */
  if (msg->ms_identity) {
    size += MOBILE_IDENTITY_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_MS_IDENTITY_PRESENT;
    if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
      memcpy(
          &emm_msg->msidentity.imsi, &msg->ms_identity->imsi,
          sizeof(emm_msg->msidentity.imsi));
    } else if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
      memcpy(
          &emm_msg->msidentity.tmsi, &msg->ms_identity->tmsi,
          sizeof(emm_msg->msidentity.tmsi));
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "TMSI  digit1 %d\n"
          "TMSI  digit2 %d\n"
          "TMSI  digit3 %d\n"
          "TMSI  digit4 %d\n",
          emm_msg->msidentity.tmsi.tmsi[0], emm_msg->msidentity.tmsi.tmsi[1],
          emm_msg->msidentity.tmsi.tmsi[2], emm_msg->msidentity.tmsi.tmsi[3]);
    }
  }
  /*
   * CSFB -Optional - Send failure cause
   */

  if (emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) {
    emm_msg->emmcause = emm_ctx->emm_cause;
  }

  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - Out of Mobile Identity for (ue_id = %u)\n",
      ue_id);
  /*
   * Optional - Additional Update Result
   */
  if (msg->additional_update_result) {
    size += ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
    emm_msg->additionalupdateresult = *msg->additional_update_result;
  }

  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - Out of Additional Update Result for (ue_id = %u)\n", ue_id);
  /*
   * CSFB -Optional - Send failure cause
   */

  if ((emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) ||
      (is_mme_ue_context_network_access_mode_packet_only(ue_mm_context_p))) {
    size += EMM_CAUSE_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_EMM_CAUSE_PRESENT;
    emm_msg->emmcause = emm_ctx->emm_cause;
  }
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - Out of Send failure cause for (ue_id = %u)\n",
      ue_id);

  /*
   * Optional - T3402
   */
  if (msg->t3402) {
    size += GPRS_TIMER_IE_MAX_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_T3402_VALUE_PRESENT;
    if (mme_config.nas_config.t3402_min <= 31) {
      emm_msg->t3402value.unit       = GPRS_TIMER_UNIT_60S;
      emm_msg->t3402value.timervalue = mme_config.nas_config.t3402_min;
    } else {
      emm_msg->t3402value.unit       = GPRS_TIMER_UNIT_360S;
      emm_msg->t3402value.timervalue = mme_config.nas_config.t3402_min / 6;
    }
  }

  /*
   * Optional - Network feature support
   */
  if (msg->eps_network_feature_support) {
    size += EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT;
    emm_msg->epsnetworkfeaturesupport = *msg->eps_network_feature_support;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}
/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_attach_accept_dl_nas()                               **
 **                                                                        **
 ** Description: Builds Attach Accept message to be sent is S1AP:DL NAS Tx **
 **                                                                        **
 **      The Attach Accept message is sent by the network to the           **
 **      UE to indicate that the corresponding attach request has          **
 **      been accepted.                                                    **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process                 **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                     **
 **      Return:    The size of the EMM message                            **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int emm_send_attach_accept_dl_nas(
    const emm_as_data_t* msg, attach_accept_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  // Get the UE context
  emm_context_t* emm_ctx = emm_context_get(&_emm_data, msg->ue_id);
  DevAssert(emm_ctx);
  ue_mm_context_t* ue_mm_context_p =
      PARENT_STRUCT(emm_ctx, struct ue_mm_context_s, emm_context);
  mme_ue_s1ap_id_t ue_id = ue_mm_context_p->mme_ue_s1ap_id;
  DevAssert(msg->ue_id == ue_id);

  OAILOG_DEBUG(LOG_NAS_EMM, "EMMAS-SAP - Send Attach Accept message\n");
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - size = EMM_HEADER_MAXIMUM_LENGTH(%d)\n", size);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = ATTACH_ACCEPT;
  /*
   * Mandatory - EPS attach result
   */
  size += EPS_ATTACH_RESULT_MAXIMUM_LENGTH;
  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += EPS_ATTACH_RESULT_MAXIMUM_LENGTH(%d)  (%d)\n",
      EPS_ATTACH_RESULT_MAXIMUM_LENGTH, size);
  switch (emm_ctx->attach_type) {
    case EMM_ATTACH_TYPE_COMBINED_EPS_IMSI:
      OAILOG_DEBUG(LOG_NAS_EMM, "EMMAS-SAP - Combined EPS/IMSI attach\n");
      /* It is observed that UE/Handest (with usage setting = voice centric and
       * voice domain preference = CS voice only) sends detach after successful
       * attach. UEs with such settings sends attach type = combined EPS and
       * IMSI attach as attach_type in attach request message. At present, EPC
       * does not support interface with MSC /CS domain and it supports only LTE
       * data service, hence it is supposed to send attach_result as EPS-only
       * and add emm_cause = "CS domain not available" for such cases. Ideally
       * in data service only n/w ,UE's usage setting should be set to data
       * centric mode and should send attach type as EPS attach only. However UE
       * settings may not be in our control. To take care of this as a
       * workaround in this patch we modified MME implementation to set EPS
       * result to Combined attach if attach type is combined attach to prevent
       * such UEs from sending detach so that such UEs can remain attached in
       * the n/w and should be able to get data service from the n/w.
       */

      /* Added check for CSFB. If Location Update procedure towards MSC/VLR
       * fails or Network Access mode received as PACKET_ONLY from HSS in ULS
       * message send epsattachresult to EPS_ATTACH_RESULT_EPS. If is it
       * successful set epsattachresult to EPS_ATTACH_RESULT_EPS_IMSI
       */
      if (((_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED) ||
           (_esm_data.conf.features & MME_API_SMS_SUPPORTED)) &&
          ((emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) ||
           is_mme_ue_context_network_access_mode_packet_only(
               ue_mm_context_p))) {
        emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS;
      } else {
        emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS_IMSI;
      }
      emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS_IMSI;
      break;
    case EMM_ATTACH_TYPE_RESERVED:
    default:
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "EMMAS-SAP - Unused attach type defaults to EPS attach\n");
    case EMM_ATTACH_TYPE_EPS:
      emm_msg->epsattachresult = EPS_ATTACH_RESULT_EPS;
      OAILOG_DEBUG(LOG_NAS_EMM, "EMMAS-SAP - EPS attach\n");
      break;
    case EMM_ATTACH_TYPE_EMERGENCY:  // We should not reach here
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMMAS-SAP - EPS emergency attach, currently unsupported\n");
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);  // TODO: fix once supported
      break;
  }
  /*
   * Mandatory - T3412 value
   */
  size += GPRS_TIMER_IE_MAX_LENGTH;
  // Check whether Periodic TAU timer is disabled
  if (mme_config.nas_config.t3412_min == 0) {
    emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_0S;
    emm_msg->t3412value.timervalue = mme_config.nas_config.t3412_min;
  } else if (mme_config.nas_config.t3412_min <= 31) {
    emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_60S;
    emm_msg->t3412value.timervalue = mme_config.nas_config.t3412_min;
  } else {
    emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_360S;
    emm_msg->t3412value.timervalue = mme_config.nas_config.t3412_min / 6;
  }
  // emm_msg->t3412value.unit = GPRS_TIMER_UNIT_0S;
  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - size += GPRS_TIMER_MAXIMUM_LENGTH(%d)  (%d)\n",
      GPRS_TIMER_IE_MAX_LENGTH, size);
  /*
   * Mandatory - Tracking area identity list
   */
  size +=
      TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH * msg->tai_list.numberoflists;
  memcpy(&emm_msg->tailist, &msg->tai_list, sizeof(msg->tai_list));
  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += "
      "TRACKING_AREA_IDENTITY_LIST_LENGTH(%d*%d)  (%d)\n",
      TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH,
      emm_msg->tailist.numberoflists, size);
  AssertFatal(
      emm_msg->tailist.numberoflists <= 16, "Too many TAIs in TAI list");
  for (int p = 0; p < emm_msg->tailist.numberoflists; p++) {
    if (TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS ==
        emm_msg->tailist.partial_tai_list[p].typeoflist) {
      size = size + (2 * emm_msg->tailist.partial_tai_list[p].numberofelements);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - size += "
          "TRACKING AREA CODE LENGTH(%d*%d)  (%d)\n",
          2, emm_msg->tailist.partial_tai_list[p].numberofelements, size);
    } else if (
        TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS ==
        emm_msg->tailist.partial_tai_list[p].typeoflist) {
      size = size + (5 * emm_msg->tailist.partial_tai_list[p].numberofelements);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - size += "
          "TRACKING AREA CODE LENGTH(%d*%d)  (%d)\n",
          5, emm_msg->tailist.partial_tai_list[p].numberofelements, size);
    }
  }
  /*
   * Mandatory - ESM message container
   */
  size += ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH + blength(msg->nas_msg);
  emm_msg->esmmessagecontainer = bstrcpy(msg->nas_msg);
  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += "
      "ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH(%d)  (%d)\n",
      ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH, size);

  /*
   * Optional - GUTI
   */
  if (msg->new_guti) {
    size += EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_GUTI_PRESENT;
    emm_msg->guti.guti.typeofidentity = EPS_MOBILE_IDENTITY_GUTI;
    emm_msg->guti.guti.oddeven        = EPS_MOBILE_IDENTITY_EVEN;
    emm_msg->guti.guti.mme_group_id   = msg->new_guti->gummei.mme_gid;
    emm_msg->guti.guti.mme_code       = msg->new_guti->gummei.mme_code;
    emm_msg->guti.guti.m_tmsi         = msg->new_guti->m_tmsi;
    emm_msg->guti.guti.mcc_digit1     = msg->new_guti->gummei.plmn.mcc_digit1;
    emm_msg->guti.guti.mcc_digit2     = msg->new_guti->gummei.plmn.mcc_digit2;
    emm_msg->guti.guti.mcc_digit3     = msg->new_guti->gummei.plmn.mcc_digit3;
    emm_msg->guti.guti.mnc_digit1     = msg->new_guti->gummei.plmn.mnc_digit1;
    emm_msg->guti.guti.mnc_digit2     = msg->new_guti->gummei.plmn.mnc_digit2;
    emm_msg->guti.guti.mnc_digit3     = msg->new_guti->gummei.plmn.mnc_digit3;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH(%d)  (%d)\n",
        EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH, size);
  }

  /*
   * Optional - LAI
   */
  if (msg->location_area_identification) {
    size += LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT;
    emm_msg->locationareaidentification.mccdigit2 =
        msg->location_area_identification->mccdigit2;
    emm_msg->locationareaidentification.mccdigit1 =
        msg->location_area_identification->mccdigit1;
    emm_msg->locationareaidentification.mncdigit3 =
        msg->location_area_identification->mncdigit3;
    emm_msg->locationareaidentification.mccdigit3 =
        msg->location_area_identification->mccdigit3;
    emm_msg->locationareaidentification.mncdigit2 =
        msg->location_area_identification->mncdigit2;
    emm_msg->locationareaidentification.mncdigit1 =
        msg->location_area_identification->mncdigit1;
    emm_msg->locationareaidentification.lac =
        msg->location_area_identification->lac;
  }

  /*
   * Optional - Mobile Identity
   */
  if (msg->ms_identity) {
    size += MOBILE_IDENTITY_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_MS_IDENTITY_PRESENT;
    if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
      memcpy(
          &emm_msg->msidentity.imsi, &msg->ms_identity->imsi,
          sizeof(emm_msg->msidentity.imsi));
    } else if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
      memcpy(
          &emm_msg->msidentity.tmsi, &msg->ms_identity->tmsi,
          sizeof(emm_msg->msidentity.tmsi));
    }
  }

  /*
   * Optional - Additional Update Result
   */
  if (msg->additional_update_result) {
    size += ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
    emm_msg->additionalupdateresult = SMS_ONLY;
  }
  /*
   * CSFB -Optional - Send failure cause
   */

  if (((emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) ||
       is_mme_ue_context_network_access_mode_packet_only(ue_mm_context_p)) &&
      (msg->emm_cause)) {
    size += EMM_CAUSE_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_EMM_CAUSE_PRESENT;
    emm_msg->emmcause = *msg->emm_cause;
  }
  /*
   * Optional - Network feature support
   */
  if (msg->eps_network_feature_support) {
    size += EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT;
    emm_msg->epsnetworkfeaturesupport = *msg->eps_network_feature_support;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_attach_reject()                                  **
 **                                                                        **
 ** Description: Builds Attach Reject message                              **
 **                                                                        **
 **      The Attach Reject message is sent by the network to the   **
 **      UE to indicate that the corresponding attach request has  **
 **      been rejected.                                            **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_attach_reject(
    const emm_as_establish_t* msg, attach_reject_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - Send Attach Reject message (cause=%d)\n",
      msg->emm_cause);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = ATTACH_REJECT;
  /*
   * Mandatory - EMM cause
   */
  size += EMM_CAUSE_MAXIMUM_LENGTH;
  emm_msg->emmcause = msg->emm_cause;

  /*
   * Optional - ESM message container
   */
  if (msg->nas_msg) {
    size +=
        ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH + blength(msg->nas_msg) +
        1;  // Adding 1 byte since ESM Container is optional IE in Attach Reject
    emm_msg->presencemask |= ATTACH_REJECT_ESM_MESSAGE_CONTAINER_PRESENT;
    emm_msg->esmmessagecontainer = msg->nas_msg;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_tracking_area_update_accept()                    **
 **                                                                        **
 ** Description: Builds Tracking Area Update Accept message                **
 **                                                                        **
 **              The Tracking Area Update Accept message is sent by the    **
 **              network to the UE to indicate that the corresponding      **
 **              tracking area update has been accepted.                   **
 **              This function is used to send TAU Accept message together **
 **              with Initial context setup request message to establish   **
 **              radio bearers as well.                                    **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_tracking_area_update_accept(
    const emm_as_establish_t* msg, tracking_area_update_accept_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send Tracking Area Update Accept message (cause=%d)\n",
      msg->emm_cause);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = TRACKING_AREA_UPDATE_ACCEPT;
  /*
   * Mandatory - EMM cause
   */
  size += EPS_UPDATE_RESULT_MAXIMUM_LENGTH;
  emm_msg->epsupdateresult = msg->eps_update_result;
  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - epsupdateresult (%d)\n",
      emm_msg->epsupdateresult);
  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - size += EPS_UPDATE_RESULT_MAXIMUM_LENGTH(%d)  (%d)\n",
      EPS_UPDATE_RESULT_MAXIMUM_LENGTH, size);

  // Optional - GPRS Timer T3412
  if (msg->t3412) {
    size += GPRS_TIMER_IE_MAX_LENGTH;
    emm_msg->presencemask |= TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_PRESENT;
    if (*msg->t3412 <= 31) {
      emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_60S;
      emm_msg->t3412value.timervalue = *msg->t3412;
    } else {
      emm_msg->t3412value.unit       = GPRS_TIMER_UNIT_360S;
      emm_msg->t3412value.timervalue = *msg->t3412 / 6;
    }
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "GPRS_TIMER_IE_MAX_LENGTH(%d)  (%d)\n",
        GPRS_TIMER_IE_MAX_LENGTH, size);
  }
  // Optional - GUTI
  if (msg->new_guti) {
    size += EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH;
    emm_msg->presencemask |= ATTACH_ACCEPT_GUTI_PRESENT;
    emm_msg->guti.guti.typeofidentity = EPS_MOBILE_IDENTITY_GUTI;
    emm_msg->guti.guti.oddeven        = EPS_MOBILE_IDENTITY_EVEN;
    emm_msg->guti.guti.mme_group_id   = msg->guti->gummei.mme_gid;
    emm_msg->guti.guti.mme_code       = msg->guti->gummei.mme_code;
    emm_msg->guti.guti.m_tmsi         = msg->guti->m_tmsi;
    emm_msg->guti.guti.mcc_digit1     = msg->guti->gummei.plmn.mcc_digit1;
    emm_msg->guti.guti.mcc_digit2     = msg->guti->gummei.plmn.mcc_digit2;
    emm_msg->guti.guti.mcc_digit3     = msg->guti->gummei.plmn.mcc_digit3;
    emm_msg->guti.guti.mnc_digit1     = msg->guti->gummei.plmn.mnc_digit1;
    emm_msg->guti.guti.mnc_digit2     = msg->guti->gummei.plmn.mnc_digit2;
    emm_msg->guti.guti.mnc_digit3     = msg->guti->gummei.plmn.mnc_digit3;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH(%d)  (%d)\n",
        EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH, size);
  }
  /* Optional - TAI list
   * This IE may be included to assign a TAI list to a UE.
   */
  if (msg->tai_list.numberoflists > 0) {
    emm_msg->presencemask |= TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_PRESENT;

    size += TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH *
            msg->tai_list.numberoflists;
    memcpy(&emm_msg->tailist, &msg->tai_list, sizeof(msg->tai_list));
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "TRACKING_AREA_IDENTITY_LIST_LENGTH(%d*%d)  (%d)\n",
        TRACKING_AREA_IDENTITY_LIST_MINIMUM_LENGTH,
        emm_msg->tailist.numberoflists, size);
    AssertFatal(
        emm_msg->tailist.numberoflists <= 16, "Too many TAIs in TAI list");
    for (int p = 0; p < emm_msg->tailist.numberoflists; p++) {
      if (TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS ==
          emm_msg->tailist.partial_tai_list[p].typeoflist) {
        size =
            size + (2 * emm_msg->tailist.partial_tai_list[p].numberofelements);
        OAILOG_INFO(
            LOG_NAS_EMM,
            "EMMAS-SAP - size += "
            "TRACKING AREA CODE LENGTH(%d*%d)  (%d)\n",
            2, emm_msg->tailist.partial_tai_list[p].numberofelements, size);
      } else if (
          TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS ==
          emm_msg->tailist.partial_tai_list[p].typeoflist) {
        size =
            size + (5 * emm_msg->tailist.partial_tai_list[p].numberofelements);
        OAILOG_INFO(
            LOG_NAS_EMM,
            "EMMAS-SAP - size += "
            "TRACKING AREA CODE LENGTH(%d*%d)  (%d)\n",
            5, emm_msg->tailist.partial_tai_list[p].numberofelements, size);
      }
    }
  }
  // Optional - EPS Bearer context status
  if (msg->eps_bearer_context_status) {
    size += EPS_BEARER_CONTEXT_STATUS_MAXIMUM_LENGTH;
    emm_msg->presencemask |=
        TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_PRESENT;
    emm_msg->epsbearercontextstatus = *(msg->eps_bearer_context_status);
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "EPS_BEARER_CONTEXT_STATUS_MAXIMUM_LENGTH(%d)  (%d)\n",
        EPS_BEARER_CONTEXT_STATUS_MAXIMUM_LENGTH, size);
  }
  // Optional - Location Area Identification, should be sent if CSFB feature is
  // enabled
  if (msg->location_area_identification) {
    size += LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH;
    emm_msg->presencemask |=
        TRACKING_AREA_UPDATE_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT;
    emm_msg->locationareaidentification.mccdigit2 =
        msg->location_area_identification->mccdigit2;
    emm_msg->locationareaidentification.mccdigit1 =
        msg->location_area_identification->mccdigit1;
    emm_msg->locationareaidentification.mncdigit3 =
        msg->location_area_identification->mncdigit3;
    emm_msg->locationareaidentification.mccdigit3 =
        msg->location_area_identification->mccdigit3;
    emm_msg->locationareaidentification.mncdigit2 =
        msg->location_area_identification->mncdigit2;
    emm_msg->locationareaidentification.mncdigit1 =
        msg->location_area_identification->mncdigit1;
    emm_msg->locationareaidentification.lac =
        msg->location_area_identification->lac;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH(%d)  (%d)",
        LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH, size);
  }
  // Optional - Mobile Identity,should be sent if CSFB feature is enabled
  if (msg->ms_identity) {
    size += EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH;
    emm_msg->presencemask |= TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_PRESENT;
    if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
      memcpy(
          &emm_msg->msidentity.imsi, &msg->ms_identity->imsi,
          sizeof(msg->ms_identity->imsi));
    } else if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_TMSI) {
      memcpy(
          &emm_msg->msidentity.tmsi, &msg->ms_identity->tmsi,
          sizeof(msg->ms_identity->tmsi));
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - size += "
          "MOBILE_IDENTITY_MINIMUM_LENGTH(%d)  (%d)",
          MOBILE_IDENTITY_MINIMUM_LENGTH, size);
    }
  }
  // Optional - EMM cause
  if (msg->emm_cause) {
    size += EMM_CAUSE_MAXIMUM_LENGTH;
    emm_msg->presencemask |= TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_PRESENT;
    emm_msg->emmcause = msg->emm_cause;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "MOBILE_IDENTITY_MINIMUM_LENGTH(%d)  (%d)\n",
        EMM_CAUSE_MAXIMUM_LENGTH, size);
  }
  // Optional - GPRS Timer T3402
  if (msg->t3402) {
    size += GPRS_TIMER_IE_MAX_LENGTH;
    emm_msg->presencemask |= TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_PRESENT;
    if (*msg->t3402 <= 31) {
      emm_msg->t3402value.unit       = GPRS_TIMER_UNIT_60S;
      emm_msg->t3402value.timervalue = *msg->t3402;
    } else {
      emm_msg->t3402value.unit       = GPRS_TIMER_UNIT_360S;
      emm_msg->t3402value.timervalue = *msg->t3402 / 6;
    }
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "GPRS_TIMER_IE_MAX_LENGTH(%d)  (%d)\n",
        GPRS_TIMER_IE_MAX_LENGTH, size);
  }
  // Optional - GPRS Timer T3423
  if (msg->t3423) {
    size += GPRS_TIMER_IE_MAX_LENGTH;
    emm_msg->presencemask |= TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_PRESENT;
    if (*msg->t3423 <= 31) {
      emm_msg->t3423value.unit       = GPRS_TIMER_UNIT_60S;
      emm_msg->t3423value.timervalue = *msg->t3423;
    } else {
      emm_msg->t3423value.unit       = GPRS_TIMER_UNIT_360S;
      emm_msg->t3423value.timervalue = *msg->t3423 / 6;
    }
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "GPRS_TIMER_IE_MAX_LENGTH(%d)  (%d)\n",
        GPRS_TIMER_IE_MAX_LENGTH, size);
  }
  // Useless actually, Optional - Equivalent PLMNs
  /*if (msg->equivalent_plmns) {
    size += PLMN_LIST_MINIMUM_LENGTH;
    emm_msg->presencemask |=
  TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_PRESENT;
    emm_msg->equivalentplmns.       = msg->;
    OAILOG_INFO (LOG_NAS_EMM, "EMMAS-SAP - size += "
  "PLMN_LIST_MINIMUM_LENGTH(%d)  (%d)", PLMN_LIST_MINIMUM_LENGTH, size);
  }*/
  /* Useless actually, Optional - Emergency number list
  if (msg->emergency_number_list) {
    size += EMERGENCY_NUMBER_LIST_MINIMUM_LENGTH;
    emm_msg->presencemask |=
  TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_PRESENT;
    emm_msg->emergencynumberlist.       = msg->;
    OAILOG_INFO (LOG_NAS_EMM, "EMMAS-SAP - size += "
  "EMERGENCY_NUMBER_LIST_MINIMUM_LENGTH(%d)  (%d)",
  EMERGENCY_NUMBER_LIST_MINIMUM_LENGTH, size);
  }*/
  // Optional - EPS network feature support
  if (msg->eps_network_feature_support) {
    size += EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH;
    emm_msg->presencemask |=
        TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT;
    emm_msg->epsnetworkfeaturesupport = *msg->eps_network_feature_support;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH(%d)  (%d)\n",
        EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH, size);
  }
  // Useless actually, Optional - Additional update result to be sent in case of
  // CSFB SMS
  if (msg->additional_update_result) {
    size += ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH;
    emm_msg->presencemask |=
        TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
    emm_msg->additionalupdateresult = *msg->additional_update_result;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH(%d)  (%d)",
        ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH, size);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}
/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_tracking_area_update_accept_dl_nas()             **
 **                                                                        **
 ** Description: Builds Tracking Area Update Accept message                **
 **                                                                        **
 **              The Tracking Area Update Accept message is sent by the    **
 **              network to the UE to indicate that the corresponding      **
 **              tracking area update has been accepted.                   **
 **              This function is used to send TAU Accept message via      **
 **              S1AP DL NAS Transport message.                            **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/

int emm_send_tracking_area_update_accept_dl_nas(
    const emm_as_data_t* msg, tracking_area_update_accept_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = TRACKING_AREA_UPDATE_ACCEPT;
  /*
   * Mandatory - EMM cause
   */
  size += EPS_UPDATE_RESULT_MAXIMUM_LENGTH;
  emm_msg->epsupdateresult = EPS_UPDATE_RESULT_TA_UPDATED;

  if (msg->eps_bearer_context_status) {
    emm_msg->presencemask |=
        TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_PRESENT;
    emm_msg->epsbearercontextstatus = *msg->eps_bearer_context_status;
    size += EPS_BEARER_CONTEXT_STATUS_MAXIMUM_LENGTH;
  }

  // Optional - EPS network feature support
  if (msg->eps_network_feature_support) {
    size += EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH;
    emm_msg->presencemask |=
        TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT;
    emm_msg->epsnetworkfeaturesupport = *msg->eps_network_feature_support;
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "EMMAS-SAP - size += "
        "EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH(%d)  (%d)\n",
        EPS_NETWORK_FEATURE_SUPPORT_MAXIMUM_LENGTH, size);
  }

  /* If CSFB is enabled send LAI,Mobile Identity
   *  and Additional Update type
   */
  if ((_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED) ||
      (_esm_data.conf.features & MME_API_SMS_SUPPORTED)) {
    if (msg->sgs_loc_updt_status == SUCCESS) {
      // LAI
      if (msg->location_area_identification) {
        emm_msg->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT;
        emm_msg->locationareaidentification.mccdigit2 =
            msg->location_area_identification->mccdigit2;
        emm_msg->locationareaidentification.mccdigit1 =
            msg->location_area_identification->mccdigit1;
        emm_msg->locationareaidentification.mncdigit3 =
            msg->location_area_identification->mncdigit3;
        emm_msg->locationareaidentification.mccdigit3 =
            msg->location_area_identification->mccdigit3;
        emm_msg->locationareaidentification.mncdigit2 =
            msg->location_area_identification->mncdigit2;
        emm_msg->locationareaidentification.mncdigit1 =
            msg->location_area_identification->mncdigit1;
        emm_msg->locationareaidentification.lac =
            msg->location_area_identification->lac;
        size += LOCATION_AREA_IDENTIFICATION_IE_MAX_LENGTH;
      }
      // Mobile Identity
      if (msg->ms_identity) {
        emm_msg->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_PRESENT;
        if (msg->ms_identity->imsi.typeofidentity == MOBILE_IDENTITY_IMSI) {
          memcpy(
              &emm_msg->msidentity.imsi, &msg->ms_identity->imsi,
              sizeof(msg->ms_identity->imsi));
          memcpy(
              &emm_msg->msidentity.tmsi, &msg->ms_identity->tmsi,
              sizeof(msg->ms_identity->tmsi));
        }
        size += MOBILE_IDENTITY_MAXIMUM_LENGTH;
      }
      // Additional Update Type
      if (msg->additional_update_result) {
        emm_msg->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
        emm_msg->additionalupdateresult = *(msg->additional_update_result);
        size += ADDITIONAL_UPDATE_RESULT_MAXIMUM_LENGTH;
      }
      // Overwrite Update type as EPS_UPDATE_RESULT_COMBINED_TA_LA_UPDATED
      emm_msg->epsupdateresult = EPS_UPDATE_RESULT_COMBINED_TA_LA_UPDATED;
    } else {
      emm_msg->emmcause = (emm_cause_t) *msg->sgs_reject_cause;
      size += EMM_CAUSE_MAXIMUM_LENGTH;
    }
  }
  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Sending DL NAS - TAU Accept\n");
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_tracking_area_update_reject()                    **
 **                                                                        **
 ** Description: Builds Tracking Area Update Reject message                **
 **                                                                        **
 **              The Tracking Area Update Reject message is sent by the    **
 **              network to the UE to indicate that the corresponding      **
 **              tracking area update has been rejected.                   **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_tracking_area_update_reject(
    const emm_as_establish_t* msg, tracking_area_update_reject_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send Tracking Area Update Reject message (cause=%d)\n",
      msg->emm_cause);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = TRACKING_AREA_UPDATE_REJECT;
  /*
   * Mandatory - EMM cause
   */
  size += EMM_CAUSE_MAXIMUM_LENGTH;
  emm_msg->emmcause = msg->emm_cause;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_service_reject()                                 **
 **                                                                        **
 ** Description: Builds Service Reject message                             **
 **                                                                        **
 **              The Tracking Area Update Reject message is sent by the    **
 **              network to the UE to indicate that the corresponding      **
 **              tracking area update has been rejected.                   **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_service_reject(
    const uint8_t emm_cause, service_reject_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - Send Service Reject message (cause=%d)\n",
      emm_cause);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = SERVICE_REJECT;
  /*
   * Mandatory - EMM cause
   */
  size += EMM_CAUSE_MAXIMUM_LENGTH;
  emm_msg->emmcause = emm_cause;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_identity_request()                               **
 **                                                                        **
 ** Description: Builds Identity Request message                           **
 **                                                                        **
 **      The Identity Request message is sent by the network to    **
 **      the UE to request the UE to provide the specified identi- **
 **      ty.                                                       **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_identity_request(
    const emm_as_security_t* msg, identity_request_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send Identity Request message for ue_id = (%u)\n",
      msg->ue_id);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = IDENTITY_REQUEST;
  /*
   * Mandatory - Identity type 2
   */
  size += IDENTITY_TYPE_2_IE_MAX_LENGTH;

  if (msg->ident_type == IDENTITY_TYPE_2_IMSI) {
    emm_msg->identitytype = IDENTITY_TYPE_2_IMSI;
  } else if (msg->ident_type == IDENTITY_TYPE_2_TMSI) {
    emm_msg->identitytype = IDENTITY_TYPE_2_TMSI;
  } else if (msg->ident_type == IDENTITY_TYPE_2_IMEI) {
    emm_msg->identitytype = IDENTITY_TYPE_2_IMEI;
  } else if (msg->ident_type == IDENTITY_TYPE_2_IMEISV) {
    emm_msg->identitytype = IDENTITY_TYPE_2_IMEISV;
  } else {
    /*
     * All other values are interpreted as "IMSI"
     */
    emm_msg->identitytype = IDENTITY_TYPE_2_IMSI;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_authentication_request()                         **
 **                                                                        **
 ** Description: Builds Authentication Request message                     **
 **                                                                        **
 **      The Authentication Request message is sent by the network **
 **      to the UE to initiate authentication of the UE identity.  **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_authentication_request(
    const emm_as_security_t* msg, authentication_request_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send Authentication Request message for ue_id "
      "= " MME_UE_S1AP_ID_FMT "\n",
      msg->ue_id);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = AUTHENTICATION_REQUEST;
  /*
   * Mandatory - NAS key set identifier
   */
  size += NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH;
  emm_msg->naskeysetidentifierasme.tsc = NAS_KEY_SET_IDENTIFIER_NATIVE;
  emm_msg->naskeysetidentifierasme.naskeysetidentifier = msg->ksi;
  /*
   * Mandatory - Authentication parameter RAND
   */
  size += AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH;
  emm_msg->authenticationparameterrand =
      blk2bstr((const void*) msg->rand, AUTH_RAND_SIZE);
  if (!emm_msg->authenticationparameterrand) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  /*
   * Mandatory - Authentication parameter AUTN
   */
  size += AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH;
  emm_msg->authenticationparameterautn =
      blk2bstr((const void*) msg->autn, AUTH_AUTN_SIZE);
  if (!emm_msg->authenticationparameterautn) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_free_send_authentication_request()                        **
 **                                                                        **
 ** Description: Frees parameters of previously created authentication     **
 **              request message                                           **
 **      Authentication parameters in the send request message use bstring **
 **      and need to be destroyed using this function                      **
 **                                                                        **
 ** Inputs:  emm_msg:       The EMM message that was sent                  **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:  None                                                         **
 **      Return:    None                                                   **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
void emm_free_send_authentication_request(authentication_request_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - Freeing Send Authentication Request message\n");
  bdestroy(emm_msg->authenticationparameterrand);
  bdestroy(emm_msg->authenticationparameterautn);
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_authentication_reject()                          **
 **                                                                        **
 ** Description: Builds Authentication Reject message                      **
 **                                                                        **
 **      The Authentication Reject message is sent by the network  **
 **      to the UE to indicate that the authentication procedure   **
 **      has failed and that the UE shall abort all activities.    **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_authentication_reject(authentication_reject_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send Authentication Reject message\n");
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = AUTHENTICATION_REJECT;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_send_security_mode_command()                          **
 **                                                                        **
 ** Description: Builds Security Mode Command message                      **
 **                                                                        **
 **      The Security Mode Command message is sent by the network  **
 **      to the UE to establish NAS signalling security.           **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:   The EMM message to be sent                 **
 **      Return:    The size of the EMM message                **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_security_mode_command(
    const emm_as_security_t* msg, security_mode_command_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send Security Mode Command message for ue_id "
      "= " MME_UE_S1AP_ID_FMT "\n",
      msg->ue_id);
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = SECURITY_MODE_COMMAND;
  /*
   * Selected NAS security algorithms
   */
  size += NAS_SECURITY_ALGORITHMS_MAXIMUM_LENGTH;
  emm_msg->selectednassecurityalgorithms.typeofcipheringalgorithm =
      msg->selected_eea;
  emm_msg->selectednassecurityalgorithms.typeofintegrityalgorithm =
      msg->selected_eia;
  /*
   * NAS key set identifier
   */
  size += NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH;
  emm_msg->naskeysetidentifier.tsc = NAS_KEY_SET_IDENTIFIER_NATIVE;
  emm_msg->naskeysetidentifier.naskeysetidentifier = msg->ksi;
  /*
   * Replayed UE Additional security capabilities
   */
  size += UE_ADDITIONAL_SECURITY_CAPABILITY_MAXIMUM_LENGTH;
  emm_msg->replayedueadditionalsecuritycapabilities._5g_ea = msg->nea;
  emm_msg->replayedueadditionalsecuritycapabilities._5g_ia = msg->nia;
  emm_msg->presencemask                                    = 0;
  /*
   * Replayed UE security capabilities
   */
  size += UE_SECURITY_CAPABILITY_MAXIMUM_LENGTH;
  emm_msg->replayeduesecuritycapabilities.eea          = msg->eea;
  emm_msg->replayeduesecuritycapabilities.eia          = msg->eia;
  emm_msg->replayeduesecuritycapabilities.umts_present = msg->umts_present;
  emm_msg->replayeduesecuritycapabilities.gprs_present = msg->gprs_present;
  emm_msg->replayeduesecuritycapabilities.uea          = msg->uea;
  emm_msg->replayeduesecuritycapabilities.uia          = msg->uia;
  emm_msg->replayeduesecuritycapabilities.gea          = msg->gea;
  emm_msg->presencemask                                = 0;

  if (msg->replayed_ue_add_sec_cap_present) {
    emm_msg->presencemask |=
        SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECU_CAPABILITY_PRESENT;
    emm_msg->replayedueadditionalsecuritycapabilities._5g_ea = msg->_5g_ea;
    emm_msg->replayedueadditionalsecuritycapabilities._5g_ia = msg->_5g_ia;
  }

  /*
   *  Setting the IMEISV Request
   */
  if (msg->imeisv_request_enabled) {
    emm_msg->presencemask |= SECURITY_MODE_COMMAND_IMEISV_REQUEST_PRESENT;
    size += IMEISV_REQUEST_IE_MAX_LENGTH;
    emm_msg->imeisvrequest = msg->imeisv_request_enabled;
    OAILOG_DEBUG(LOG_NAS_EMM, "imeisv flag :%d\n", emm_msg->imeisvrequest);
  }
  OAILOG_DEBUG(
      LOG_NAS_EMM, "replayeduesecuritycapabilities.gprs_present %d\n",
      emm_msg->replayeduesecuritycapabilities.gprs_present);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "replayeduesecuritycapabilities.gea          %d\n",
      emm_msg->replayeduesecuritycapabilities.gea);

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_emm_information()                                **
 **                                                                        **
 ** Description: Builds EMM Information message                            **
 **                                                                        **
 **              The EMM Information message is sent by the                **
 **              network to the UE to Allow the network to provide         **
 **              Information to UE, The UE may use the received            **
 **              information.                                              **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_emm_information(
    const emm_as_data_t* msg, emm_information_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  char string[MAX_MINUTE_DIGITS] = {0};
  int size                       = EMM_HEADER_MAXIMUM_LENGTH;
  int addci                      = 0;
  int noOfBits                   = 0;
  time_t t;
  struct tm* tmp;
  struct tm updateTime;
  uint8_t formatted;

  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = EMM_INFORMATION;

  /*
   * optional - full network name
   */
  size += NETWORK_NAME_IE_MAX_LENGTH;
  emm_msg->presencemask |= EMM_INFORMATION_FULL_NAME_FOR_NETWORK_PRESENT;
  emm_msg->fullnamefornetwork.codingscheme = 0x00;
  emm_msg->fullnamefornetwork.addci        = addci;

  noOfBits = blength(msg->full_network_name) * 7;

  /* Check if there are any spare bits */
  if (noOfBits % 8 != 0) {
    emm_msg->fullnamefornetwork.numberofsparebitsinlastoctet =
        8 - (noOfBits % 8);
  } else {
    emm_msg->fullnamefornetwork.numberofsparebitsinlastoctet = 0;
  }
  emm_msg->fullnamefornetwork.textstring = bstrcpy(msg->full_network_name);

  /*
   * optional - short network name
   */
  size += NETWORK_NAME_IE_MAX_LENGTH;
  emm_msg->presencemask |= EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_PRESENT;
  emm_msg->shortnamefornetwork.codingscheme = 0x00;
  emm_msg->shortnamefornetwork.addci        = addci;

  noOfBits = blength(msg->short_network_name) * 7;

  /* Check if there are any spare bits */
  if (noOfBits % 8 != 0) {
    emm_msg->shortnamefornetwork.numberofsparebitsinlastoctet =
        8 - (noOfBits % 8);
  } else {
    emm_msg->shortnamefornetwork.numberofsparebitsinlastoctet = 0;
  }
  emm_msg->shortnamefornetwork.textstring = bstrcpy(msg->short_network_name);

  /*
   * optional - Local Time Zone
   */
  int result = get_time_zone();
  if (result != RETURNerror) {
    emm_msg->localtimezone = result;
    size += TIME_ZONE_IE_MAX_LENGTH;
    emm_msg->presencemask |= EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT;
  }

  /*
   * optional - Universal time and Local Time Zone
   */
  t   = time(NULL);
  tmp = localtime_r(&t, &updateTime);
  if (tmp == NULL) {
    OAILOG_ERROR(LOG_NAS_EMM, "localtime() failed to get local timer info");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  size += TIME_ZONE_AND_TIME_MAX_LENGTH;
  emm_msg->presencemask |=
      EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_PRESENT;

  /*
   * Time SHALL be encoded as specified in 3GPP 23.040 in SM-TL TPDU format.
   */
  /*
   * updateTime.year is "years since 1900"
   * GSM format is the last 2 digits of the year.
   */
  snprintf(string, sizeof(string), "%02d", updateTime.tm_year + 1900 - 2000);
  formatted = (string[1] - '0') << 4;
  formatted |= (string[0] - '0');
  emm_msg->universaltimeandlocaltimezone.year = formatted;
  /*
   * updateTime.tm_mon is "months since January" in [0-11] range.
   * GSM format is months in [1-12] range.
   */
  snprintf(string, sizeof(string), "%02d", updateTime.tm_mon + 1);
  formatted = (string[1] - '0') << 4;
  formatted |= (string[0] - '0');
  emm_msg->universaltimeandlocaltimezone.month = formatted;
  snprintf(string, sizeof(string), "%02d", updateTime.tm_mday);
  formatted = (string[1] - '0') << 4;
  formatted |= (string[0] - '0');
  emm_msg->universaltimeandlocaltimezone.day = formatted;
  snprintf(string, sizeof(string), "%02d", updateTime.tm_hour);
  formatted = (string[1] - '0') << 4;
  formatted |= (string[0] - '0');
  emm_msg->universaltimeandlocaltimezone.hour = formatted;
  snprintf(string, sizeof(string), "%02d", updateTime.tm_min);
  formatted = (string[1] - '0') << 4;
  formatted |= (string[0] - '0');
  emm_msg->universaltimeandlocaltimezone.minute = formatted;
  snprintf(string, sizeof(string), "%02d", updateTime.tm_sec);
  formatted = (string[1] - '0') << 4;
  formatted |= (string[0] - '0');
  emm_msg->universaltimeandlocaltimezone.second = formatted;
  if ((emm_msg->presencemask && EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT) != 0) {
    emm_msg->universaltimeandlocaltimezone.timezone = emm_msg->localtimezone;
  }
  /*
   * optional - Daylight Saving Time
   */
  size += DAYLIGHT_SAVING_TIME_IE_MAX_LENGTH;
  emm_msg->presencemask |= EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_PRESENT;
  emm_msg->networkdaylightsavingtime = msg->daylight_saving_time;

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_dl_nas_transport()                               **
 **                                                                        **
 ** Description: Builds Downlink Nas Transport message                     **
 **                                                                        **
 **              The Downlink Nas Transport message is sent by the         **
 **              network to the UE to transfer the SMS recieved from MSC   **
 **              This function is used to send DL NAS Transport message    **
 **              via S1AP DL NAS Transport message.                        **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/

int emm_send_dl_nas_transport(
    const emm_as_data_t* msg, downlink_nas_transport_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;
  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = DOWNLINK_NAS_TRANSPORT;
  /*
   * Mandatory - Nas message container
   */
  size += NAS_MESSAGE_CONTAINER_MAXIMUM_LENGTH;
  emm_msg->nasmessagecontainer = bstrcpy(msg->nas_msg);
  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Sending DL NAS - DL NAS Transport\n");
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}
/****************************************************************************
 **                                                                        **
 ** Name:    emm_free_send_emm_information()                               **
 **                                                                        **
 ** Description: Frees parameters of previously created emm information    **
 **              message                                                   **
 **      network name text string parameters in the send request message   **
 **      use bstring and need to be destroyed using this function          **
 **                                                                        **
 ** Inputs:  emm_msg:       The EMM message that was sent                  **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:  None                                                         **
 **      Return:    None                                                   **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
void emm_free_send_emm_information(emm_information_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - Freeing send EMM information message\n");
  bdestroy(emm_msg->fullnamefornetwork.textstring);
  bdestroy(emm_msg->shortnamefornetwork.textstring);
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:        get_time_zone()                                             **
 **                                                                        **
 ** Description: Gets the Timezone                                         **
 **                                                                        **
 **                                                                        **
 ** Outputs:            Returns the Timezone numerical format              **
 **                                                                        **
 ***************************************************************************/
int get_time_zone(void) {
  char timestr[20]               = {0};
  char string[MAX_MINUTE_DIGITS] = {0};
  time_t t;
  struct tm* tmp;
  struct tm updateTime;
  int hour = 0, minute = 0, timezone = 0;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  t   = time(NULL);
  tmp = localtime_r(&t, &updateTime);
  if (tmp == NULL) {
    OAILOG_ERROR(LOG_NAS_EMM, "localtime_r() failed to get local timer info ");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  if (strftime(timestr, sizeof(timestr), "%z", &updateTime) == 0) {
    OAILOG_ERROR(LOG_NAS_EMM, "EMMAS-SAP - strftime() Failed to get timezone");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  /* the string shall be in the form of +hhmm (+0530) */
  if (timestr[0] == '-') {
    timezone = 0x08;
  }
  hour   = (10 * (timestr[1] - '0')) | (timestr[2] - '0');
  minute = ((10 * (timestr[3] - '0')) | (timestr[4] - '0')) + (hour * 60);
  minute /= 15;

  snprintf(string, sizeof(string), "%02d", minute);
  timezone |= (string[1] - '0') << 4;
  timezone |= (string[0] - '0');

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, timezone);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_send_cs_service_notification()                        **
 **                                                                        **
 ** Description: Builds CS Service Notification message                    **
 **                                                                        **
 **              The CS Service Notification message is sent by the        **
 **              network to the UE to inform UE about MT call              **
 **              And UE shall switch to CS service network GERAN/UTRAN     **
 **              for Voice service                                         **
 **                                                                        **
 ** Inputs:      msg:           The EMMAS-SAP primitive to process         **
 **              Others:        None                                       **
 **                                                                        **
 ** Outputs:     emm_msg:       The EMM message to be sent                 **
 **              Return:        The size of the EMM message                **
 **              Others:        None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_send_cs_service_notification(
    const emm_as_data_t* msg, cs_service_notification_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = EMM_HEADER_MAXIMUM_LENGTH;

  /*
   * Mandatory - Message type
   */
  emm_msg->messagetype = CS_SERVICE_NOTIFICATION;

  /*
   * Mandatory -PagingIdentity
   */
  size += PAGING_IDENTITY_MAXIMUM_LENGTH;
  emm_msg->pagingidentity = msg->paging_identity;

  /*
   * optional - Calling Line Indentification (CLI)
   */
  if (msg->cli != NULL) {
    size += CLI_MAXIMUM_LENGTH;
    emm_msg->presencemask |= CS_SERVICE_NOTIFICATION_CLI_PRESENT;
    emm_msg->cli = bstrcpy(msg->cli);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, size);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_free_send_dl_nas_transport()                              **
 **                                                                        **
 ** Description: Frees parameters of previously created dl nas transport   **
 **              message                                                   **
 **      nas message container parameter in the send dl nas transport msg  **
 **      use bstring and need to be destroyed using this function          **
 **                                                                        **
 ** Inputs:  emm_msg:       The EMM message that was sent                  **
 ** Outputs:  None                                                         **
 **      Return:    None                                                   **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
void emm_free_send_dl_nas_transport(downlink_nas_transport_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMMAS-SAP - Freeing send DL NAS Transport message\n");
  bdestroy(emm_msg->nasmessagecontainer);
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}
/****************************************************************************
 **                                                                        **
 ** Name:    emm_free_send_cs_service_notification()                       **
 **                                                                        **
 ** Description: Frees parameters of previously created cs service         **
 **              notification                                              **
 **      CLI text string parameters in the send request message            **
 **      use bstring and need to be destroyed using this function          **
 **                                                                        **
 ** Inputs:  emm_msg:       The EMM message: CS Service Notification       **
                            message that was sent                          **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:  None                                                         **
 **      Return:    None                                                   **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
void emm_free_send_cs_service_notification(
    cs_service_notification_msg* emm_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - Freeing CLI sent in CS Service Notification message\n");
  if ((emm_msg->presencemask & CS_SERVICE_NOTIFICATION_CLI_PRESENT) ==
      CS_SERVICE_NOTIFICATION_CLI_PRESENT) {
    if (emm_msg->cli != NULL) {
      bdestroy(emm_msg->cli);
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}
