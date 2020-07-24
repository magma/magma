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

/*! \file mme_app_hss_reset.c
   \brief
   \author Sebastien ROUX, Lionel GAUTHIER
   \version 1.0
   \company Eurecom
   \email: lionel.gauthier@eurecom.fr
*/

#include <stdio.h>
#include <pthread.h>
#include <stdbool.h>
#include <mme_app_state.h>

#include "common_defs.h"
#include "log.h"
#include "mme_app_ue_context.h"
#include "mme_app_defs.h"
#include "hashtable.h"
#include "mme_api.h"
#include "mme_app_desc.h"
#include "s6a_messages_types.h"

int mme_app_handle_s6a_reset_req(const s6a_reset_req_t* const rsr_pP) {
  int rc                               = RETURNok;
  struct ue_mm_context_s* ue_context_p = NULL;
  hash_node_t* node                    = NULL;
  unsigned int i                       = 0;
  unsigned int num_elements            = 0;
  hash_table_ts_t* hashtblP            = NULL;

  OAILOG_FUNC_IN(LOG_MME_APP);

  OAILOG_DEBUG(LOG_MME_APP, "S6a Reset Request received\n");

  if (rsr_pP == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP, "Invalid S6a Reset Request ITTI message received\n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }

  hashtblP = get_mme_ue_state();
  if (!hashtblP) {
    OAILOG_INFO(LOG_MME_APP, "There is no Ue Context in the MME context \n");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNok);
  }
  while ((num_elements < hashtblP->num_elements) && (i < hashtblP->size)) {
    pthread_mutex_lock(&hashtblP->lock_nodes[i]);
    if (hashtblP->nodes[i] != NULL) {
      node = hashtblP->nodes[i];
    }
    pthread_mutex_unlock(&hashtblP->lock_nodes[i]);
    while (node) {
      num_elements++;
      hashtable_ts_get(
          hashtblP, (const hash_key_t) node->key, (void**) &ue_context_p);
      if (ue_context_p != NULL) {
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
          if (ue_context_p->sgs_context != NULL) {
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
      node = node->next;
    }
    i++;
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, rc);
}
