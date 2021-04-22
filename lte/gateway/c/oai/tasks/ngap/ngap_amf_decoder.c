/**
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
/****************************************************************************
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages
*****************************************************************************/

#include <stdbool.h>
#include <string.h>

#include "ngap_amf_decoder.h"
#include "bstrlib.h"
#include "log.h"
#include "assertions.h"
#include "common_defs.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_InitiatingMessage.h"
#include "Ngap_ProcedureCode.h"
#include "Ngap_SuccessfulOutcome.h"
#include "Ngap_UnsuccessfulOutcome.h"
#include "asn_codecs.h"
#include "constr_TYPE.h"
#include "per_decoder.h"

int ngap_amf_decode_pdu(Ngap_NGAP_PDU_t* pdu, const_bstring const raw) {
  asn_dec_rval_t dec_ret;
  DevAssert(pdu != NULL);
  DevAssert(blength(raw) != 0);
  dec_ret = aper_decode(
      NULL, &asn_DEF_Ngap_NGAP_PDU, (void**) &pdu, bdata(raw), blength(raw), 0,
      0);

  if (dec_ret.code != RC_OK) {
    OAILOG_ERROR(LOG_NGAP, "Failed to decode PDU\n");
    return -1;
  }
  return 0;
}
