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
#include <stdlib.h>

#include "bstrlib.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "3gpp_36.401.h"
#include "common_defs.h"
#include "mme_app_ue_context.h"
#include "esm_proc.h"
#include "log.h"
#include "esm_data.h"
#include "esm_pt.h"
#include "emm_sap.h"
#include "3gpp_24.301.h"
#include "EsmCause.h"
#include "emm_data.h"
#include "emm_esmDef.h"
#include "mme_api.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

extern int pdn_connectivity_delete(emm_context_t* emm_context, pdn_cid_t pid);

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
    Internal data handled by the PDN disconnect procedure in the MME
   --------------------------------------------------------------------------
*/
/*
   PDN disconnection handlers
*/
static int pdn_disconnect_get_pid(emm_context_t* emm_context, proc_tid_t pti);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
          PDN disconnect procedure executed by the MME
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_pdn_disconnect_request()                         **
 **                                                                        **
 ** Description: Performs PDN disconnect procedure requested by the UE.    **
 **                                                                        **
 **              3GPP TS 24.301, section 6.5.2.3                           **
 **      Upon receipt of the PDN DISCONNECT REQUEST message, if it **
 **      is accepted by the network, the MME shall initiate the    **
 **      bearer context deactivation procedure.                    **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      pti:       Identifies the PDN disconnect procedure    **
 **             requested by the UE                        **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    The identifier of the PDN connection to be **
 **             released, if it exists;                    **
 **             RETURNerror otherwise.                     **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_pdn_disconnect_request(
    emm_context_t* emm_context, proc_tid_t pti, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  pdn_cid_t pid = RETURNerror;
  ue_mm_context_t* ue_context_p =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
  mme_ue_s1ap_id_t ue_id     = ue_context_p->mme_ue_s1ap_id;
  int nb_active_pdn_contexts = ue_context_p->nb_active_pdn_contexts;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - PDN disconnect requested by the UE "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d)\n",
      ue_id, pti);

  /*
   * Get UE's ESM context
   */
  if (nb_active_pdn_contexts > 1) {
    /*
     * Get the identifier of the PDN connection entry assigned to the
     * * * * procedure transaction identity
     */
    pid = pdn_disconnect_get_pid(emm_context, pti);

    if (pid >= MAX_APN_PER_UE) {
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "ESM-PROC  - No PDN connection found (pti=%d) for ue "
          "id " MME_UE_S1AP_ID_FMT "\n",
          pti, ue_id);
      *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
    }
  } else {
    /*
     * Attempt to disconnect from the last PDN disconnection
     * * * * is not allowed
     */
    *esm_cause = ESM_CAUSE_LAST_PDN_DISCONNECTION_NOT_ALLOWED;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, pid);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_pdn_disconnect_accept()                          **
 **                                                                        **
 ** Description: Performs PDN disconnect procedure accepted by the UE.     **
 **                                                                        **
 **              3GPP TS 24.301, section 6.5.2.3                           **
 **      On reception of DEACTIVATE EPS BEARER CONTEXT ACCEPT mes- **
 **      sage from the UE, the MME releases all the resources re-  **
 **      served for the PDN in the network.                        **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      pid:       Identifier of the PDN connection to be     **
 **             released                                   **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     esm_cause: Cause code returned upon ESM procedure     **
 **             failure                                    **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_pdn_disconnect_accept(
    emm_context_t* emm_context, pdn_cid_t pid, esm_cause_t* esm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;
  OAILOG_INFO(
      LOG_NAS_ESM,
      "ESM-PROC  - PDN disconnect accepted by the UE "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", pid=%d)\n",
      ue_id, pid);
  /*
   * Release the connectivity with the requested PDN
   */
  int rc = mme_api_unsubscribe(NULL);

  if (rc != RETURNerror) {
    /*
     * Delete the PDN connection entry
     */
    proc_tid_t pti = pdn_connectivity_delete(emm_context, pid);

    if (pti != ESM_PT_UNASSIGNED) {
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
    }
  }

  *esm_cause = ESM_CAUSE_PROTOCOL_ERROR;
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_proc_pdn_disconnect_reject()                          **
 **                                                                        **
 ** Description: Performs PDN disconnect procedure not accepted by the     **
 **      network.                                                  **
 **                                                                        **
 **              3GPP TS 24.301, section 6.5.2.4                           **
 **      Upon receipt of the PDN DISCONNECT REQUEST message, if it **
 **      is not accepted by the network, the MME shall send a PDN  **
 **      DISCONNECT REJECT message to the UE.                      **
 **                                                                        **
 ** Inputs:  is_standalone: Not used - Always true                     **
 **      ue_id:      UE lower layer identifier                  **
 **      ebi:       Not used                                   **
 **      msg:       Encoded PDN disconnect reject message to   **
 **             be sent                                    **
 **      ue_triggered:  Not used                                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int esm_proc_pdn_disconnect_reject(
    const bool is_standalone, emm_context_t* emm_context, ebi_t ebi,
    STOLEN_REF bstring* msg, const bool ue_triggered) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc;
  emm_sap_t emm_sap = {0};
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_WARNING(
      LOG_NAS_ESM,
      "ESM-PROC  - PDN disconnect not accepted by the network "
      "(ue_id=" MME_UE_S1AP_ID_FMT ")\n",
      ue_id);
  /*
   * Notity EMM that ESM PDU has to be forwarded to lower layers
   */
  emm_sap.primitive            = EMMESM_UNITDATA_REQ;
  emm_sap.u.emm_esm.ue_id      = ue_id;
  emm_sap.u.emm_esm.ctx        = emm_context;
  emm_sap.u.emm_esm.u.data.msg = *msg;
  rc                           = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
                Timer handlers
   --------------------------------------------------------------------------
*/

/*
  ---------------------------------------------------------------------------
                PDN disconnection handlers
  ---------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    _pdn_disconnect_get_pid()                                 **
 **                                                                        **
 ** Description: Returns the identifier of the PDN connection to which the **
 **      given procedure transaction identity has been assigned    **
 **      to establish connectivity to the specified UE             **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       The procedure transaction identity         **
 **      Others:    _esm_data                                  **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    The identifier of the PDN connection if    **
 **             found in the list; -1 otherwise.           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static pdn_cid_t pdn_disconnect_get_pid(
    emm_context_t* emm_context, proc_tid_t pti) {
  pdn_cid_t i = MAX_APN_PER_UE;

  if (emm_context) {
    ue_mm_context_t* ue_mm_context =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
    for (i = 0; i < MAX_APN_PER_UE; i++) {
      if (ue_mm_context->pdn_contexts[i]) {
        if (ue_mm_context->pdn_contexts[i]->esm_data.pti == pti) {
          return (i);
        }
      }
    }
  }

  return RETURNerror;
}
