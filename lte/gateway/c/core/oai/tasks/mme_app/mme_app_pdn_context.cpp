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

/*! \file mme_app_pdn_context.cpp
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_pdn_context.hpp"

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_apn_selection.hpp"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.hpp"

static void mme_app_pdn_context_init(ue_mm_context_t* const ue_context,
                                     pdn_context_t* const pdn_context);

//------------------------------------------------------------------------------
void mme_app_free_pdn_context(pdn_context_t** const pdn_context,
                              imsi64_t imsi64) {
  bdestroy_wrapper(&(*pdn_context)->apn_in_use);
  bdestroy_wrapper(&(*pdn_context)->apn_subscribed);
  bdestroy_wrapper(&(*pdn_context)->apn_oi_replacement);
  free_protocol_configuration_options(&(*pdn_context)->pco);
  free_wrapper((void**)pdn_context);
}
//------------------------------------------------------------------------------
static void mme_app_pdn_context_init(ue_mm_context_t* const ue_context,
                                     pdn_context_t* const pdn_context) {
  if ((pdn_context) && (ue_context)) {
    memset(pdn_context, 0, sizeof(*pdn_context));
    pdn_context->is_active = false;
    for (int i = 0; i < BEARERS_PER_UE; i++) {
      // contains bearer indexes in ue_mm_context_t.bearer_contexts[]
      pdn_context->bearer_contexts[i] = INVALID_BEARER_INDEX;
    }
  }
}
//------------------------------------------------------------------------------
pdn_context_t* mme_app_create_pdn_context(
    ue_mm_context_t* const ue_mm_context, const pdn_cid_t pdn_cid,
    const context_identifier_t context_identifier) {
  OAILOG_FUNC_IN(LOG_MME_APP);
  if (!ue_mm_context->pdn_contexts[pdn_cid]) {
    pdn_context_t* pdn_context =
        reinterpret_cast<pdn_context_t*>(calloc(1, sizeof(*pdn_context)));

    if (pdn_context) {
      struct apn_configuration_s* apn_configuration =
          mme_app_get_apn_config(ue_mm_context, context_identifier);
      if (apn_configuration) {
        mme_app_pdn_context_init(ue_mm_context, pdn_context);

        pdn_context->context_identifier = apn_configuration->context_identifier;
        pdn_context->pdn_type = apn_configuration->pdn_type;
        if (apn_configuration->nb_ip_address) {
          pdn_context->paa.pdn_type =
              apn_configuration->ip_address[0]
                  .pdn_type;  // TODO check this later...
          pdn_context->paa.ipv4_address =
              apn_configuration->ip_address[0].address.ipv4_address;
          if (2 == apn_configuration->nb_ip_address) {
            pdn_context->paa.ipv6_address =
                apn_configuration->ip_address[1].address.ipv6_address;
            pdn_context->paa.ipv6_prefix_length = 64;
          }
        }
        // pdn_context->apn_oi_replacement
        memcpy(&pdn_context->default_bearer_eps_subscribed_qos_profile,
               &apn_configuration->subscribed_qos,
               sizeof(eps_subscribed_qos_profile_t));
        memcpy(&pdn_context->subscribed_apn_ambr, &apn_configuration->ambr,
               sizeof(ambr_t));
        // pdn_context->pgw_apn_ambr = ;
        pdn_context->apn_subscribed =
            blk2bstr(apn_configuration->service_selection,
                     apn_configuration->service_selection_length);
        pdn_context->default_ebi = EPS_BEARER_IDENTITY_UNASSIGNED;
        // TODO pdn_context->apn_in_use     =

        ue_mm_context->pdn_contexts[pdn_cid] = pdn_context;

        OAILOG_FUNC_RETURN(LOG_MME_APP, pdn_context);
      } else {
        free_wrapper((void**)&pdn_context);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_MME_APP, NULL);
}
