/**
 * copyright 2020 The Magma Authors.
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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include <unordered_map>
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.hpp"

extern amf_config_t amf_config;
namespace magma5g {

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_attribute_valid()                                 **
**                                                                        **
** Description: set the amf_context attribute as valid                    **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_ctx_set_attribute_valid(amf_context_t* ctxt,
                                 const uint32_t attribute_bit_pos) {
  ctxt->member_present_mask |= attribute_bit_pos;
  ctxt->member_valid_mask |= attribute_bit_pos;
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_attribute_present()                               **
**                                                                        **
** Description: set the amf_context attribute as present                  **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_ctx_set_attribute_present(amf_context_t* ctxt,
                                   const int attribute_bit_pos) {
  ctxt->member_present_mask |= attribute_bit_pos;
}

/***************************************************************************
**                                                                        **
** Name:    get_5g_nas_common_procedure_identification()                  **
**                                                                        **
** Description:  return identification procedure                          **
**                                                                        **
**                                                                        **
***************************************************************************/
nas_amf_ident_proc_t* get_5g_nas_common_procedure_identification(
    const amf_context_t* ctxt) {
  return (nas_amf_ident_proc_t*)get_nas5g_common_procedure(ctxt,
                                                           AMF_COMM_PROC_IDENT);
}

/***************************************************************************
**                                                                        **
** Name:    amf_ctx_set_valid_imei()                                      **
**                                                                        **
** Description:   set imei and mark it as valid                           **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_ctx_set_valid_imei(amf_context_t* const ctxt, imei_t* imei) {
  ctxt->imei = *imei;
  amf_ctx_set_attribute_valid(ctxt, AMF_CTXT_MEMBER_IMEI);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_proc_identification_complete()                            **
 **                                                                        **
 ** Description: Performs the identification completion procedure executed **
 **      by the network.                                                   **
 **                                                                        **
 **              3GPP TS 24.501, section 5.4.3.4                           **
 **      Upon receiving the IDENTITY RESPONSE message, the MME             **
 **      shall stop timer T3470.                                           **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      imsi:      The IMSI received from the UE                          **
 **      imei:      The IMEI received from the UE                          **
 **      tmsi:      The TMSI received from the UE                          **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    amf_data, T3570                                        **
 **                                                                        **
 ***************************************************************************/
status_code_e amf_proc_identification_complete(
    const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
    imeisv_t* const imeisv, uint32_t* const tmsi, guti_m5_t* amf_ctx_guti) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  status_code_e rc = RETURNok;
  amf_context_t* amf_ctx = NULL;

  OAILOG_DEBUG(LOG_NAS_AMF,
               "Identification procedure complete for "
               "(ue_id= " AMF_UE_NGAP_ID_FMT ")\n",
               ue_id);

  ue_m5gmm_context_s* ue_mm_context =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_mm_context) {
    amf_ctx = &ue_mm_context->amf_context;
    nas_amf_ident_proc_t* ident_proc =
        get_5g_nas_common_procedure_identification(amf_ctx);

    if (ident_proc == NULL) {
      OAILOG_ERROR(LOG_AMF_APP,
                   "Failed to start identity procedure for "
                   "ue_id " AMF_UE_NGAP_ID_FMT "\n",
                   ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }

    OAILOG_DEBUG(LOG_AMF_APP, "Timer: Stopping Identity timer with ID %lu\n",
                 ident_proc->T3570.id);
    amf_app_stop_timer(ident_proc->T3570.id);
    ident_proc->T3570.id = NAS5G_TIMER_INACTIVE_ID;

    if (imsi) {
      /*
       * Update the IMSI
       */
      imsi64_t imsi64 = amf_imsi_to_imsi64(imsi);
      amf_ctx_set_valid_imsi(amf_ctx, imsi, imsi64);
      amf_context_upsert_imsi(amf_ctx);
      amf_ctx->imsi64 = imsi64;
      amf_ctx->imsi.length = 8;
      amf_ctx->m5_guti = *amf_ctx_guti;
    } else {
      OAILOG_ERROR(LOG_AMF_APP,
                   "should not happen because this type of identity is not "
                   "requested by AMF");
    }
    /*
     * Notify AMF that the identification procedure successfully completed
     */

    amf_registration_success_identification_cb(amf_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    amf_app_generate_guti_on_supi()                               **
**                                                                        **
** Description: Generate GUTI based on SUPI/IMSI received                 **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_app_generate_guti_on_supi(amf_guti_m5g_t* amf_guti,
                                   supi_as_imsi_t* supi_imsi) {
  /* Generate GUTI with 5g-tmsi as rand value */
  amf_guti->guamfi.plmn.mcc_digit1 = supi_imsi->plmn.mcc_digit1;
  amf_guti->guamfi.plmn.mcc_digit2 = supi_imsi->plmn.mcc_digit2;
  amf_guti->guamfi.plmn.mcc_digit3 = supi_imsi->plmn.mcc_digit3;
  amf_guti->guamfi.plmn.mnc_digit1 = supi_imsi->plmn.mnc_digit1;
  amf_guti->guamfi.plmn.mnc_digit2 = supi_imsi->plmn.mnc_digit2;
  amf_guti->guamfi.plmn.mnc_digit3 = supi_imsi->plmn.mnc_digit3;

  // tmsi value is 4 octet random value.
  amf_guti->m_tmsi = htonl((uint32_t)rand());

  // Filling data from amf_config file considering only one gNB
  amf_config_read_lock(&amf_config);
  amf_guti->guamfi.amf_regionid = amf_config.guamfi.guamfi[0].amf_regionid;

  // TODO: Temp hardcoded change to remove later
  amf_guti->guamfi.amf_set_id = amf_config.guamfi.guamfi[0].amf_set_id;
  amf_guti->guamfi.amf_pointer = amf_config.guamfi.guamfi[0].amf_pointer;

  OAILOG_DEBUG(LOG_AMF_APP, "amf_region_id %u amf_set_id %u amf_pointer %u",
               amf_config.guamfi.guamfi[0].amf_regionid,
               amf_config.guamfi.guamfi[0].amf_set_id,
               amf_config.guamfi.guamfi[0].amf_pointer);

  amf_config_unlock(&amf_config);
  return;
}

}  // namespace magma5g
