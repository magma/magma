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

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_as.h"
#include "conversions.h"
#include "secu_defs.h"

typedef uint32_t amf_ue_ngap_id_t;
#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH)                                    \
  LENGTH += QUADLET - (LENGTH % QUADLET)

namespace magma5g {
nas_amf_smc_proc_t smc_data;

/***************************************************************************
**                                                                        **
** Name:  amf_ctx_set_valid_imsi()                                        **
**                                                                        **
** Description: Set IMSI, mark it as valid                                **
**                                                                        **
**                                                                        **
***************************************************************************/

void amf_ctx_set_valid_imsi(
    amf_context_t* ctxt, imsi_t* imsi, const imsi64_t imsi64) {
  ctxt->imsi                     = *imsi;
  ctxt->imsi64                   = imsi64;
  ctxt->is_initial_identity_imsi = true;
#if DEBUG_IS_ON
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
  IMSI64_TO_STRING(ctxt->imsi64, imsi_str, ctxt->imsi.length);
  OAILOG_DEBUG(LOG_AMF_APP, "imsi : %s", imsi_str);
#endif
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_security_eksi()                                   **
**                                                                        **
** Description: sets security context eksi                                **
**                                                                        **
**                                                                        **
***************************************************************************/
void nas_amf_smc_proc_t::amf_ctx_set_security_eksi(
    amf_context_t* ctxt, ksi_t eksi) {
  ctxt->_security.eksi = eksi;
  OAILOG_TRACE(
      LOG_NAS_AMF,
      "ue_id= " AMF_UE_NGAP_ID_FMT " set security context eksi %d\n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id,
      eksi);
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_security_type()                                   **
**                                                                        **
** Description: sets security context type                                **
**                                                                        **
**                                                                        **
***************************************************************************/
void nas_amf_smc_proc_t::amf_ctx_set_security_type(
    amf_context_t* ctxt, amf_sc_type_t sc_type) {
  ctxt->_security.sc_type = sc_type;
  OAILOG_TRACE(
      LOG_NAS_AMF,
      "ue_id= " AMF_UE_NGAP_ID_FMT " set security context security type %d\n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id,
      sc_type);
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_clear_security                                        **
**                                                                        **
** Description: clears amf security context                               **
**                                                                        **
**                                                                        **
***************************************************************************/
void nas_amf_smc_proc_t::amf_ctx_clear_security(amf_context_t* ctxt) {
  memset(&ctxt->_security, 0, sizeof(ctxt->_security));
  smc_data.amf_ctx_set_security_type(ctxt, SECURITY_CTX_TYPE_NOT_AVAILABLE);
  smc_data.amf_ctx_set_security_eksi(ctxt, KSI_NO_KEY_AVAILABLE);
  ctxt->_security.selected_algorithms.encryption =
      0;  // NAS_SECURITY_ALGORITHMS_EEA0;
  ctxt->_security.selected_algorithms.integrity =
      0;  // NAS_SECURITY_ALGORITHMS_EIA0;
  ctxt->_security.direction_decode = SECU_DIRECTION_UPLINK;
  ctxt->_security.direction_encode = SECU_DIRECTION_DOWNLINK;
  OAILOG_DEBUG(
      LOG_NAS_AMF, "ue_id= " AMF_UE_NGAP_ID_FMT " cleared security context \n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id);
}

}  // namespace magma5g
