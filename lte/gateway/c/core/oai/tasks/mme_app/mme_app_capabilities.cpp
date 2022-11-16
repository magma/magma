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

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_app_desc.hpp"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.hpp"

status_code_e mme_app_handle_s1ap_ue_capabilities_ind(
    const itti_s1ap_ue_cap_ind_t* const s1ap_ue_cap_ind_pP) {
  ue_mm_context_t* ue_context_p = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);
  if (s1ap_ue_cap_ind_pP == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid S1AP UE Capability Indication ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  ue_context_p =
      mme_ue_context_exists_mme_ue_s1ap_id(s1ap_ue_cap_ind_pP->mme_ue_s1ap_id);
  if (!ue_context_p) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "UE context doesn't exist for enb_ue_s1ap_ue_id " ENB_UE_S1AP_ID_FMT
        " mme_ue_s1ap_id " MME_UE_S1AP_ID_FMT "\n",
        s1ap_ue_cap_ind_pP->enb_ue_s1ap_id, s1ap_ue_cap_ind_pP->mme_ue_s1ap_id);

    free_wrapper((void**)&s1ap_ue_cap_ind_pP->radio_capabilities);
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  if (ue_context_p->ue_radio_capability) {
    bdestroy_wrapper(&ue_context_p->ue_radio_capability);
  }

  // Allocate the radio capabilities memory. Note that this takes care of the
  // length = 0 case for us quite nicely.
  ue_context_p->ue_radio_capability =
      blk2bstr(s1ap_ue_cap_ind_pP->radio_capabilities,
               s1ap_ue_cap_ind_pP->radio_capabilities_length);

  OAILOG_DEBUG(LOG_MME_APP,
               "UE radio capabilities of length %d found and cached\n",
               blength(ue_context_p->ue_radio_capability));

  OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
}
