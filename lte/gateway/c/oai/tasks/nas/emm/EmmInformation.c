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
  Source      EmmInformation.c

  Version     0.1

  Date        2017/09/29

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author

  Description Defines sending of EMM Information  from the Network

        The purpose of sending the EMM INFORMATION message is to allow
        the network to provide information to the UE.

        The EMM information procedure may be invoked by the network
        at any time during an established EMM context.

*****************************************************************************/
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "log.h"
#include "common_defs.h"
#include "emm_data.h"
#include "emm_sap.h"
#include "mme_app_ue_context.h"
#include "bstrlib.h"
#include "emm_asDef.h"
#include "mme_api.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

static void emm_information_pack_gsm_7Bit(bstring str, unsigned char* result);

int emm_proc_emm_informtion(ue_mm_context_t* ue_emm_ctx) {
  int rc                    = RETURNerror;
  unsigned char result[256] = {0};
  emm_sap_t emm_sap         = {0};
  emm_as_data_t* emm_as     = &emm_sap.u.emm_as.u.data;
  emm_context_t* emm_ctx    = &(ue_emm_ctx->emm_context);
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  /*
   * Setup NAS information message to transfer
   */
  emm_as->nas_info = EMM_AS_NAS_EMM_INFORMATION;
  emm_as->nas_msg  = NULL;  // No ESM container
  /*
   * Set the UE identifier
   */
  emm_as->ue_id = ue_emm_ctx->mme_ue_s1ap_id;

  emm_as->daylight_saving_time = _emm_data.conf.daylight_saving_time;

  /*
   * Encode full_network_name with gsm 7 bit encoding
   * The encoding is done referring to 3gpp 24.008
   * (section: 10.5.3.5a)and 23.038
   */
  emm_information_pack_gsm_7Bit(_emm_data.conf.full_network_name, result);
  emm_as->full_network_name = bfromcstr((const char*) result);
  /*
   * Encode short_network_name with gsm 7 bit encoding
   */
  memset(result, 0, sizeof(result));
  emm_information_pack_gsm_7Bit(_emm_data.conf.short_network_name, result);
  emm_as->short_network_name = bfromcstr((const char*) result);

  /*
   * Setup EPS NAS security data
   */
  emm_as_set_security_data(&emm_as->sctx, &emm_ctx->_security, false, true);
  /*
   * Notify EMM-AS SAP that TAU Accept message has to be sent to the network
   */
  emm_sap.primitive = EMMAS_DATA_REQ;
  rc                = emm_sap_send(&emm_sap);
  bdestroy(emm_as->full_network_name);
  bdestroy(emm_as->short_network_name);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

static void emm_information_pack_gsm_7Bit(bstring str, unsigned char* result) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int encIdx = 0;
  int len, i = 0, j = 0;
  len = blength(str);
  len -= 1;

  while (i <= len) {
    if (i < len) {
      result[encIdx++] =
          ((bchar(str, i) >> j) | ((bchar(str, i + 1) << (7 - j)) & 0xFF));
    } else {
      result[encIdx++] = ((bchar(str, i) >> j) & 0x7f);
    }

    i++;
    j = (j + 1) % 7;

    if (j == 0) {
      i++;
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}
