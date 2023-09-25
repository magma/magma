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

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_asDef.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EmmCause.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_status_ind()                                     **
 **                                                                        **
 ** Description: Processes received EMM status message.                    **
 **                                                                        **
 **      3GPP TS 24.301, section 5.7                               **
 **      On receipt of an EMM STATUS message no state transition   **
 **      and no specific action shall be taken. Local actions are  **
 **      possible and are implementation dependent.                **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **          emm_cause: Received EMM cause code                    **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_status_ind(mme_ue_s1ap_id_t ue_id,
                                  emm_cause_t emm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  status_code_e rc = RETURNok;

  OAILOG_INFO(LOG_NAS_EMM,
              "EMM-PROC  - EMM status procedure requested (cause=%d)",
              emm_cause);
  OAILOG_DEBUG(LOG_NAS_EMM, "EMM-PROC  - To be implemented");
  increment_counter("emm_status_rcvd", 1, NO_LABELS);

  /*
   * TODO
   */
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_status()                                         **
 **                                                                        **
 ** Description: Initiates EMM status procedure.                           **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      emm_cause: EMM cause code to be reported              **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_status(mme_ue_s1ap_id_t ue_id, emm_cause_t emm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  increment_counter("emm_status_sent", 1, NO_LABELS);
  status_code_e rc = RETURNerror;
  emm_sap_t emm_sap = {};
  emm_security_context_t* sctx = NULL;
  struct emm_context_s* ctx = NULL;

  /*
   * Notity EMM that EMM status indication has to be sent to lower layers
   */
  emm_sap.primitive = EMMAS_STATUS_IND;
  emm_sap.u.emm_as.u.status.emm_cause = emm_cause;
  emm_sap.u.emm_as.u.status.ue_id = ue_id;
  emm_sap.u.emm_as.u.status.guti = NULL;
  ue_mm_context_t* ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (ue_mm_context) {
    ctx = &ue_mm_context->emm_context;
    if (ctx) {
      sctx = &ctx->_security;
    }
    OAILOG_INFO_UE(LOG_NAS_EMM, ctx->_imsi64,
                   "EMM-PROC  - EMM status procedure requested\n");
  } else {
    OAILOG_INFO(LOG_NAS_EMM, "EMM-PROC  - EMM status procedure requested\n");
  }

  /*
   * Setup EPS NAS security data
   */
  emm_as_set_security_data(&emm_sap.u.emm_as.u.status.sctx, sctx, false, true);
  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
