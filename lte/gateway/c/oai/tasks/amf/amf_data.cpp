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
/*****************************************************************************

  Source      amf_data.cpp

  Version     0.1

  Date        2020/07/28

  Product     AMF

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#ifdef __cplusplus
}
#endif
#include "amf_fsm.h"
#include "amf_asDefs.h"
#include "amf_as.h"
#include "amf_sap.h"
#include "amf_app_ue_context_and_proc.h"
#include "conversions.h"
#include "secu_defs.h"

using namespace std;
typedef uint32_t amf_ue_ngap_id_t;
#define QUADLET 4
#define AMF_GET_BYTE_ALIGNED_LENGTH(LENGTH)                                    \
  LENGTH += QUADLET - (LENGTH % QUADLET)

namespace magma5g {
nas_amf_smc_proc_t smc_data;
// void amf_ctx_set_attribute_valid(
//   amf_context_t* ctxt, const int attribute_bit_pos) {
// ctxt->member_present_mask |= attribute_bit_pos; //TODO -  NEED-RECHECK
// ctxt->member_valid_mask |= attribute_bit_pos;
//}
/* Set IMSI, mark it as valid */
void amf_ctx_set_valid_imsi(
    amf_context_t* ctxt, imsi_t* imsi, const imsi64_t imsi64) {
  ctxt->_imsi                    = *imsi;
  ctxt->_imsi64                  = imsi64;
  ctxt->is_initial_identity_imsi = true;
// amf_ctx_set_attribute_valid(ctxt, AMF_CTXT_MEMBER_IMSI);
#if DEBUG_IS_ON
  char imsi_str[IMSI_BCD_DIGITS_MAX + 1] = {0};
  IMSI64_TO_STRING(ctxt->_imsi64, imsi_str, ctxt->_imsi.length);
#if 0
    OAILOG_DEBUG(LOG_NAS_AMF, "ue_id=" AMF_UE_NGAP_ID_FMT " set IMSI %s (valid)\n",
        (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))
            ->amf_ue_ngap_id,
        imsi_str);
#endif
#endif
  // TODO
  // amf_api_notify_imsi((PARENT_STRUCT(ctxt, ue_m5gmm_context_s,
  // amf_context))->amf_ue_ngap_id, imsi64);
}
//------------------------------------------------------------------------------
void nas_amf_smc_proc_t::amf_ctx_set_security_eksi(
    amf_context_t* ctxt, ksi_t eksi) {
  ctxt->_security.eksi = eksi;
  OAILOG_TRACE(
      LOG_NAS_AMF,
      "ue_id=" AMF_UE_NGAP_ID_FMT " set security context eksi %d\n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id,
      eksi);
}

//------------------------------------------------------------------------------
void nas_amf_smc_proc_t::amf_ctx_set_security_type(
    amf_context_t* ctxt, amf_sc_type_t sc_type) {
  ctxt->_security.sc_type = sc_type;
  OAILOG_TRACE(
      LOG_NAS_AMF,
      "ue_id=" AMF_UE_NGAP_ID_FMT " set security context security type %d\n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id,
      sc_type);
}

inline void amf_ctx_clear_attribute_present(  // TODO and define new variables
    amf_context_t* ctxt, const int attribute_bit_pos) {
  // ctxt->member_present_mask &= ~attribute_bit_pos;
  // ctxt->member_valid_mask &= ~attribute_bit_pos;
}
//------------------------------------------------------------------------------
/* Clear security  */

void nas_amf_smc_proc_t::amf_ctx_clear_security(amf_context_t* ctxt) {
  memset(&ctxt->_security, 0, sizeof(ctxt->_security));
  smc_data.amf_ctx_set_security_type(ctxt, SECURITY_CTX_TYPE_NOT_AVAILABLE);
  smc_data.amf_ctx_set_security_eksi(ctxt, KSI_NO_KEY_AVAILABLE);
  ctxt->_security.selected_algorithms.encryption =
      0;  // NAS_SECURITY_ALGORITHMS_EEA0;
  ctxt->_security.selected_algorithms.integrity =
      0;  // NAS_SECURITY_ALGORITHMS_EIA0;
  amf_ctx_clear_attribute_present(ctxt, AMF_CTXT_MEMBER_SECURITY);
  ctxt->_security.direction_decode = SECU_DIRECTION_UPLINK;
  ctxt->_security.direction_encode = SECU_DIRECTION_DOWNLINK;
  OAILOG_DEBUG(
      LOG_NAS_AMF, "ue_id=" AMF_UE_NGAP_ID_FMT " cleared security context \n",
      (PARENT_STRUCT(ctxt, ue_m5gmm_context_s, amf_context))->amf_ue_ngap_id);
}

}  // namespace magma5g
