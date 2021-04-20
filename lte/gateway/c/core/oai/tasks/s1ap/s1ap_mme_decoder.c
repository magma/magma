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

/*! \file s1ap_mme_decoder.c
   \brief s1ap decode procedures for MME
   \author Sebastien ROUX <sebastien.roux@eurecom.fr>
   \date 2012
   \version 0.1
*/

#include <stdbool.h>
#include <string.h>

#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "common_defs.h"
#include "s1ap_mme_decoder.h"
#include "S1ap_S1AP-PDU.h"
#include "S1ap_InitiatingMessage.h"
#include "S1ap_ProcedureCode.h"
#include "S1ap_SuccessfulOutcome.h"
#include "S1ap_UnsuccessfulOutcome.h"
#include "asn_codecs.h"
#include "constr_TYPE.h"
#include "per_decoder.h"

//-----------------------------------------------------------------------------
int s1ap_mme_decode_pdu(S1ap_S1AP_PDU_t* pdu, const_bstring const raw) {
  if ((pdu) && (raw)) {
    if (blength(raw) == 0) {
      OAILOG_DEBUG(LOG_S1AP, "Buffer length is Zero \n");
    }
    asn_dec_rval_t dec_ret = aper_decode(
        NULL, &asn_DEF_S1ap_S1AP_PDU, (void**) &pdu, bdata(raw), blength(raw),
        0, 0);

    if (dec_ret.code != RC_OK) {
      OAILOG_ERROR(LOG_S1AP, "Failed to decode PDU\n");
      return RETURNerror;
    }
    return RETURNok;
  } else {
    OAILOG_DEBUG(LOG_S1AP, "PDU is NULL \n");
    return RETURNerror;
  }
}
