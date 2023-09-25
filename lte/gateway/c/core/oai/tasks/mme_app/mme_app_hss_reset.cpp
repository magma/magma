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

/*! \file mme_app_hss_reset.cpp
   \brief
   \author Sebastien ROUX, Lionel GAUTHIER
   \version 1.0
   \company Eurecom
   \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <pthread.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/mme_app_state.hpp"
#include "lte/gateway/c/core/oai/include/mme_app_desc.hpp"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.hpp"

status_code_e mme_app_handle_s6a_reset_req(
    const s6a_reset_req_t* const rsr_pP) {
  status_code_e rc = RETURNok;
  struct ue_mm_context_s* ue_context_p = nullptr;
  unsigned int i = 0;
  unsigned int num_elements = 0;

  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_DEBUG(LOG_MME_APP, "S6a Reset Request received\n");

  if (rsr_pP == nullptr) {
    OAILOG_ERROR(LOG_MME_APP,
                 "Invalid S6a Reset Request ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  proto_map_uint32_ue_context_t* mme_app_state_ue_map = get_mme_ue_state();
  if (!mme_app_state_ue_map) {
    OAILOG_ERROR(LOG_MME_APP, "mme_app_state_ue_map doesn't exist");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  for (auto itr = mme_app_state_ue_map->map->begin();
       itr != mme_app_state_ue_map->map->end(); itr++) {
    mme_app_state_ue_map->get(itr->first, &ue_context_p);
    if (ue_context_p != nullptr) {
      if (ue_context_p->mm_state == UE_REGISTERED) {
        /*
         * set the flag: location_info_confirmed_in_hss to indicate that,
         * hss has restarted and MME shall send ULR to hss
         */
        ue_context_p->location_info_confirmed_in_hss = true;
        /*
         * set the sgs context flag: neaf to indicate that,
         * hss has restarted and MME shall send SGS Ue Activity Indication to
         * MSC/VLR to indicate that activity from a UE has been detected
         */
        if (ue_context_p->sgs_context != nullptr) {
          ue_context_p->sgs_context->neaf = true;
        }

        if (ue_context_p->ecm_state == ECM_CONNECTED) {
          /*
           * hss has restarted and MME shall send ULR to hss for connected Ue
           */
          rc = mme_app_send_s6a_update_location_req(ue_context_p);
        }
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}
