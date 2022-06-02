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
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.hpp"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"

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

void amf_ctx_set_valid_imsi(amf_context_t* ctxt, imsi_t* imsi,
                            const imsi64_t imsi64) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ctxt->imsi = *imsi;
  ctxt->imsi64 = imsi64;
  ctxt->is_initial_identity_imsi = true;

  amf_ue_ngap_id_t ue_id =
      PARENT_STRUCT(ctxt, struct ue_m5gmm_context_s, amf_context)
          ->amf_ue_ngap_id;

  amf_api_notify_imsi(ue_id, imsi64);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_security_eksi()                                   **
**                                                                        **
** Description: sets security context eksi                                **
**                                                                        **
**                                                                        **
***************************************************************************/
void nas_amf_smc_proc_t::amf_ctx_set_security_eksi(amf_context_t* ctxt,
                                                   ksi_t eksi) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ctxt->_security.eksi = eksi;
  ctxt->ksi = eksi;

  OAILOG_TRACE(
      LOG_AMF_APP,
      "ue_id= " AMF_UE_NGAP_ID_FMT " set security context eksi %d\n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id,
      eksi);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_security_type()                                   **
**                                                                        **
** Description: sets security context type                                **
**                                                                        **
**                                                                        **
***************************************************************************/
void nas_amf_smc_proc_t::amf_ctx_set_security_type(amf_context_t* ctxt,
                                                   amf_sc_type_t sc_type) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ctxt->_security.sc_type = sc_type;
  OAILOG_TRACE(
      LOG_AMF_APP,
      "ue_id= " AMF_UE_NGAP_ID_FMT " set security context security type %d\n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id,
      sc_type);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
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
  OAILOG_FUNC_IN(LOG_AMF_APP);
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
      LOG_AMF_APP, "ue_id= " AMF_UE_NGAP_ID_FMT " cleared security context \n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

}  // namespace magma5g
