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

#include <stdbool.h>
#include <string.h>

/*****************************************************************************
  Source      NasTransportHdl.c

  Version     0.1

  Date        2018/06/11

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author

  Description Defines the Nas Transport EMM procedure executed by the
        Non-Access Stratum.

        The purpose of the nas transport procedure is to transfer
        the nas message from ue to msc/vlr and vice versa

*****************************************************************************/
#include "emm_proc.h"
#include "log.h"
#include "emm_data.h"
#include "mme_app_itti_messaging.h"
#include "conversions.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "DetachRequest.h"
#include "MobileStationClassmark2.h"
#include "bstrlib.h"
#include "common_defs.h"
#include "common_types.h"
#include "esm_data.h"
#include "mme_api.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/
/*
   --------------------------------------------------------------------------
    Internal data handled by the service request procedure in the UE
   --------------------------------------------------------------------------
*/

/*
   --------------------------------------------------------------------------
    Internal data handled by the service request procedure in the MME
   --------------------------------------------------------------------------
*/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_uplink_nas_transport()                               **
 **                                                                        **
 ** Description: Send the uplink nas transport procedure upon receiving    **
 **      Uplink Nas Transport message from the UE.                         **
 **                                                                        **
 **              3GPP TS 24.301, section 5.6.3.2                           **
 **      Upon receiving an UPLINK NAS TRANSPORT message, the MME shall  **
 **      send the available imsi,imeisv,ue time zone,                   **
 **      mobilestationclassmark2,tai,ecgi and recieved nas message      **
 **      container(SMS) to MSC/VLR.                                        **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      nas_msg_pP:   uplink nas message container                 **
 **      Others:    _emm_data                                  **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    _emm_data                                  **
 **                                                                        **
 ***************************************************************************/
int emm_proc_uplink_nas_transport(mme_ue_s1ap_id_t ue_id, bstring nas_msg_pP) {
  int rc                                     = RETURNok;
  emm_context_t* emm_ctxt_p                  = NULL;
  imeisv_t* p_imeisv                         = NULL;
  MobileStationClassmark2* p_mob_st_clsMark2 = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * Get the UE's EMM context if it exists
   */

  emm_ctxt_p = emm_context_get(&_emm_data, ue_id);

  if (emm_ctxt_p != NULL) {
    ue_mm_context_t* ue_mm_context_p =
        PARENT_STRUCT(emm_ctxt_p, struct ue_mm_context_s, emm_context);
    /* check if the non EPS service control is enable and combined attach. If
     * in SMS_ORC8R, we still want to send the uplink message, but we should
     * disable vlr checks since SGs is not present.*/
    if (((_esm_data.conf.features & MME_API_SMS_SUPPORTED) ||
         (_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED) ||
         (_esm_data.conf.features & MME_API_SMS_ORC8R_SUPPORTED)) &&
        (emm_ctxt_p->attach_type == EMM_ATTACH_TYPE_COMBINED_EPS_IMSI)) {
      // check if vlr reliable flag is true for sgs association.
      if (mme_ue_context_get_ue_sgs_vlr_reliable(ue_id) ||
          (_esm_data.conf.features & MME_API_SMS_ORC8R_SUPPORTED)) {
        if (IS_EMM_CTXT_PRESENT_IMEISV(emm_ctxt_p)) {
          p_imeisv = &emm_ctxt_p->_imeisv;
        }
        if (IS_EMM_CTXT_PRESENT_MOB_STATION_CLSMARK2(emm_ctxt_p)) {
          p_mob_st_clsMark2 = &emm_ctxt_p->_mob_st_clsMark2;
        }
        // Send SGS Uplink unitdata message towards SGS task.
        char imsi_str[IMSI_BCD_DIGITS_MAX + 1];

        IMSI_TO_STRING(&emm_ctxt_p->_imsi, imsi_str, IMSI_BCD_DIGITS_MAX + 1);

        nas_itti_sgsap_uplink_unitdata(
            imsi_str, strlen(imsi_str), nas_msg_pP, p_imeisv, p_mob_st_clsMark2,
            &emm_ctxt_p->originating_tai, &ue_mm_context_p->e_utran_cgi,
            _esm_data.conf.features & MME_API_SMS_ORC8R_SUPPORTED);
      } else {
        if (emm_ctxt_p->is_imsi_only_detach == true) {
          OAILOG_DEBUG(
              LOG_NAS_EMM,
              "Already triggred Detach Request for the UE "
              "(ue_id=" MME_UE_S1AP_ID_FMT ") \n",
              ue_id);
        } else {
          // NAS trigger UE to re-attach for non-EPS services.
          emm_proc_nw_initiated_detach_request(
              ue_id, NW_DETACH_TYPE_IMSI_DETACH);
        }
      }
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "No EMM context exists for the UE (ue_id=" MME_UE_S1AP_ID_FMT ") \n",
        ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
