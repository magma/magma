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

/*! \file mme_app_bearer_context.cpp
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_bearer_context.hpp"

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>

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
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_29.274.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.hpp"

static void mme_app_bearer_context_init(bearer_context_t* const bearer_context);

//------------------------------------------------------------------------------
static void mme_app_bearer_context_init(
    bearer_context_t* const bearer_context) {
  if (bearer_context) {
    memset(bearer_context, 0, sizeof(*bearer_context));

    esm_bearer_context_init(&bearer_context->esm_ebr_context);
  }
}
//------------------------------------------------------------------------------
bearer_context_t* mme_app_create_bearer_context(
    ue_mm_context_t* const ue_mm_context, const pdn_cid_t pdn_cid,
    const ebi_t ebi, const bool is_default) {
  ebi_t lebi = ebi;
  if ((EPS_BEARER_IDENTITY_FIRST > ebi) || (EPS_BEARER_IDENTITY_LAST < ebi)) {
    OAILOG_ERROR(LOG_NAS_ESM, "Received invalid ebi :%u \n", ebi);
    return NULL;
  }
  lebi = mme_app_get_free_bearer_id(ue_mm_context);
  if (EPS_BEARER_IDENTITY_UNASSIGNED == lebi) {
    return NULL;
  }

  bearer_context_t* bearer_context =
      reinterpret_cast<bearer_context_t*>(malloc(sizeof(*bearer_context)));

  if (bearer_context) {
    mme_app_bearer_context_init(bearer_context);
    bearer_context->ebi = lebi;
    mme_app_add_bearer_context(ue_mm_context, bearer_context, pdn_cid,
                               is_default);
  }
  return bearer_context;
}

//------------------------------------------------------------------------------
void mme_app_free_bearer_context(bearer_context_t** const bearer_context) {
  free_esm_bearer_context(&(*bearer_context)->esm_ebr_context);
  free_wrapper((void**)bearer_context);
}

//------------------------------------------------------------------------------
bearer_context_t* mme_app_get_bearer_context(ue_mm_context_t* const ue_context,
                                             const ebi_t ebi) {
  if ((ue_context) && (EPS_BEARER_IDENTITY_LAST >= ebi) &&
      (EPS_BEARER_IDENTITY_FIRST <= ebi)) {
    return ue_context->bearer_contexts[EBI_TO_INDEX(ebi)];
  }
  return NULL;
}

//------------------------------------------------------------------------------
void mme_app_add_bearer_context(ue_mm_context_t* const ue_context,
                                bearer_context_t* const bc,
                                const pdn_cid_t pdn_cid,
                                const bool is_default) {
  if (bc->ebi > EPS_BEARER_IDENTITY_LAST ||
      bc->ebi < EPS_BEARER_IDENTITY_FIRST) {
    OAILOG_ERROR(LOG_MME_APP,
                 "Invalid EBI (%u) received in bearer context "
                 "for MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
                 bc->ebi, ue_context->mme_ue_s1ap_id);
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }
  int index = EBI_TO_INDEX(bc->ebi);
  if (!ue_context->bearer_contexts[index]) {
    if (ue_context->pdn_contexts[pdn_cid]) {
      bc->pdn_cx_id = pdn_cid;
      ue_context->pdn_contexts[pdn_cid]->bearer_contexts[index] = index;
      ue_context->bearer_contexts[index] = bc;

      bc->preemption_capability =
          ue_context->pdn_contexts[pdn_cid]
              ->default_bearer_eps_subscribed_qos_profile
              .allocation_retention_priority.pre_emp_capability;
      bc->preemption_vulnerability =
          ue_context->pdn_contexts[pdn_cid]
              ->default_bearer_eps_subscribed_qos_profile
              .allocation_retention_priority.pre_emp_vulnerability;
      bc->priority_level = ue_context->pdn_contexts[pdn_cid]
                               ->default_bearer_eps_subscribed_qos_profile
                               .allocation_retention_priority.priority_level;
      return;
    }
    OAILOG_WARNING(LOG_MME_APP,
                   "No PDN id %u exist for ue id " MME_UE_S1AP_ID_FMT, pdn_cid,
                   ue_context->mme_ue_s1ap_id);
    return;
  }
  OAILOG_WARNING(
      LOG_MME_APP,
      "Bearer ebi %u PDN id %u already exist for ue id " MME_UE_S1AP_ID_FMT,
      bc->ebi, pdn_cid, ue_context->mme_ue_s1ap_id);
}

//------------------------------------------------------------------------------
ebi_t mme_app_get_free_bearer_id(ue_mm_context_t* const ue_context) {
  for (int i = 0; i < BEARERS_PER_UE; i++) {
    if (!ue_context->bearer_contexts[i]) {
      return INDEX_TO_EBI(i);
    }
  }
  return EPS_BEARER_IDENTITY_UNASSIGNED;
}

//------------------------------------------------------------------------------
void mme_app_bearer_context_s1_release_enb_informations(
    bearer_context_t* const bc) {
  if (bc) {
    memset(&bc->enb_fteid_s1u, 0, sizeof(bc->enb_fteid_s1u));
    bc->enb_fteid_s1u.teid = INVALID_TEID;
  }
}
