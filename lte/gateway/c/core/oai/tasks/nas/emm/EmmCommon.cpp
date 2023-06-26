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

/*****************************************************************************
  Source      EmmCommon.cpp

  Version     0.1

  Date        2013/04/19

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author      Frederic Maurel

  Description Defines callback functions executed within EMM common procedures
        by the Non-Access Stratum running at the network side.

        Following EMM common procedures can always be initiated by the
        network whilst a NAS signalling connection exists:

        GUTI reallocation
        authentication
        security mode control
        identification
        EMM information

*****************************************************************************/
#include "lte/gateway/c/core/oai/tasks/nas/emm/EmmCommon.hpp"

#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <pthread.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/lib/gtpv2-c/nwgtpv2c-0.11/include/tree.h"
/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

emm_common_data_head_t emm_common_data_head = {PTHREAD_MUTEX_INITIALIZER,
                                               RB_INITIALIZER()};

static inline int emm_common_data_compare_ueid(struct emm_common_data_s* p1,
                                               struct emm_common_data_s* p2);

RB_HEAD(emm_common_data_map, emm_common_data_s) emm_common_data_root;

RB_PROTOTYPE(emm_common_data_map, emm_common_data_s, entries,
             emm_common_data_compare_ueid);

/* Generate functions used for the MAP */
RB_GENERATE(emm_common_data_map, emm_common_data_s, entries,
            emm_common_data_compare_ueid);

static inline int emm_common_data_compare_ueid(struct emm_common_data_s* p1,
                                               struct emm_common_data_s* p2) {
  if (p1->ue_id > p2->ue_id) {
    return 1;
  }

  if (p1->ue_id < p2->ue_id) {
    return -1;
  }

  /*
   * Matching reference -> return 0
   */
  return 0;
}

struct emm_common_data_s* emm_common_data_context_get(
    struct emm_common_data_head_s* root, mme_ue_s1ap_id_t _ueid) {
  struct emm_common_data_s reference;
  struct emm_common_data_s* reference_p = NULL;

  DevAssert(root);
  DevCheck(_ueid > 0, _ueid, 0, 0);
  memset(&reference, 0, sizeof(struct emm_common_data_s));
  reference.ue_id = _ueid;
  pthread_mutex_lock(&root->mutex);
  reference_p =
      RB_FIND(emm_common_data_map,
              (emm_common_data_map*)&root->emm_common_data_root, &reference);
  pthread_mutex_unlock(&root->mutex);
  return reference_p;
}

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_initialize()                              **
 **                                                                        **
 ** Description: Initialize EMM procedure callback functions executed for  **
 **      the UE with the given identifier                          **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      success:   EMM procedure executed upon successful EMM **
 **             common procedure completion                **
 **      reject:    EMM procedure executed if the EMM common   **
 **             procedure failed or is rejected            **
 **      failure:   EMM procedure executed upon transmission   **
 **             failure reported by lower layer            **
 **      abort:     EMM common procedure executed when the on- **
 **             going EMM procedure is aborted             **
 **      args:      EMM common procedure argument parameters   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    _emm_common_data                           **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_common_initialize(
    mme_ue_s1ap_id_t ue_id, emm_common_success_callback_t _success,
    emm_common_reject_callback_t _reject,
    emm_common_failure_callback_t _failure,
    emm_common_ll_failure_callback_t _ll_failure,
    emm_common_non_delivered_callback_t _non_delivered,
    emm_common_abort_callback_t _abort, void* args) {
  struct emm_common_data_s* emm_common_data_ctx = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  assert(ue_id > 0);
  emm_common_data_ctx =
      emm_common_data_context_get(&emm_common_data_head, ue_id);

  if (emm_common_data_ctx == NULL) {
    emm_common_data_ctx = reinterpret_cast<emm_common_data_t*>(
        calloc(1, sizeof(emm_common_data_t)));
    emm_common_data_ctx->ue_id = ue_id;
    pthread_mutex_lock(&emm_common_data_head.mutex);
    RB_INSERT(emm_common_data_map,
              (emm_common_data_map*)&emm_common_data_head.emm_common_data_root,
              emm_common_data_ctx);
    pthread_mutex_unlock(&emm_common_data_head.mutex);

    if (emm_common_data_ctx) {
      emm_common_data_ctx->ref_count = 0;
    }
  }

  if (emm_common_data_ctx) {
    __sync_fetch_and_add(&emm_common_data_ctx->ref_count, 1);
    emm_common_data_ctx->success = _success;
    emm_common_data_ctx->reject = _reject;
    emm_common_data_ctx->failure = _failure;
    emm_common_data_ctx->ll_failure = _ll_failure;
    emm_common_data_ctx->non_delivered = _non_delivered;
    emm_common_data_ctx->abort = _abort;
    emm_common_data_ctx->args = args;
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_success()                                 **
 **                                                                        **
 ** Description: The EMM common procedure initiated between the UE with    **
 **      the specified identifier and the MME completed success-   **
 **      fully. The network performs required actions related to   **
 **      the ongoing EMM procedure.                                **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    _emm_common_data, _emm_data                **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_common_success(emm_common_data_t* emm_common_data_ctx) {
  emm_common_success_callback_t emm_callback = {0};
  status_code_e rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->success;

    if (emm_callback) {
      struct emm_context_s* ctx =
          emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc = (*emm_callback)(ctx);
    }

    emm_common_cleanup(emm_common_data_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_reject()                                  **
 **                                                                        **
 ** Description: The EMM common procedure initiated between the UE with    **
 **      the specified identifier and the MME failed or has been   **
 **      rejected. The network performs required actions related   **
 **      to the ongoing EMM procedure.                             **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    _emm_common_data, _emm_data                **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_common_reject(emm_common_data_t* emm_common_data_ctx) {
  status_code_e rc = RETURNerror;
  emm_common_reject_callback_t emm_callback;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->reject;

    if (emm_callback) {
      struct emm_context_s* ctx =
          emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc = (*emm_callback)(ctx);
    }

    emm_common_cleanup(emm_common_data_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

status_code_e emm_proc_common_failure(emm_common_data_t* emm_common_data_ctx) {
  status_code_e rc = RETURNerror;
  emm_common_reject_callback_t emm_callback;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->failure;

    if (emm_callback) {
      struct emm_context_s* ctx =
          emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc = (*emm_callback)(ctx);
    }

    emm_common_cleanup(emm_common_data_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_ll_failure()                                 **
 **                                                                        **
 ** Description: The EMM common procedure has been initiated between the   **
 **      UE with the specified identifier and the MME, and a lower **
 **      layer failure occurred before the EMM common procedure    **
 **      being completed. The network performs required actions    **
 **      related to the ongoing EMM procedure.                     **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    _emm_common_data, _emm_data                **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_common_ll_failure(
    emm_common_data_t* emm_common_data_ctx) {
  emm_common_ll_failure_callback_t emm_callback;
  status_code_e rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->ll_failure;

    if (emm_callback) {
      struct emm_context_s* ctx = NULL;

      ctx = emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc = (*emm_callback)(ctx);
    }

    emm_common_cleanup(emm_common_data_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_non_delivered()                                 **
 **                                                                        **
 ** Description: The EMM common procedure has been initiated between the   **
 **      UE with the specified identifier and the MME, and a report **
 **      of lower layers stated that the EMM common procedure message   **
 **      could not be delivered. The network performs required actions    **
 **      related to the ongoing EMM procedure.                     **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    _emm_common_data, _emm_data                **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_common_non_delivered(
    emm_common_data_t* emm_common_data_ctx) {
  emm_common_non_delivered_callback_t emm_callback;
  status_code_e rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->non_delivered;

    if (emm_callback) {
      struct emm_context_s* ctx = NULL;

      ctx = emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc = (*emm_callback)(ctx);
    }

    // emm_common_cleanup (emm_common_data_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_abort()                                   **
 **                                                                        **
 ** Description: The ongoing EMM procedure has been aborted. The network   **
 **      performs required actions related to the EMM common pro-  **
 **      cedure previously initiated between the UE with the spe-  **
 **      cified identifier and the MME.                            **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    _emm_common_data                           **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
status_code_e emm_proc_common_abort(emm_common_data_t* emm_common_data_ctx) {
  emm_common_abort_callback_t emm_callback;
  status_code_e rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->abort;

    if (emm_callback) {
      struct emm_context_s* ctx = NULL;

      ctx = emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc = (*emm_callback)(ctx);
    }

    emm_common_cleanup(emm_common_data_ctx);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_get_args()                                **
 **                                                                        **
 ** Description: Returns pointer to the EMM common procedure argument pa-  **
 **      rameters allocated for the UE with the given identifier.  **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    _emm_common_data                           **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    pointer to the EMM common procedure argu-  **
 **             ment parameters                            **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
void* emm_proc_common_get_args(mme_ue_s1ap_id_t ue_id) {
  emm_common_data_t* emm_common_data_ctx = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_common_data_ctx =
      emm_common_data_context_get(&emm_common_data_head, ue_id);
  if (emm_common_data_ctx) {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, emm_common_data_ctx->args);
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, NULL);
  }
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_common_cleanup()                                 **
 **                                                                        **
 ** Description: Cleans EMM procedure callback functions upon completion   **
 **      of an EMM common procedure previously initiated within an **
 **      EMM procedure currently in progress between the network   **
 **      and the UE with the specified identifier.                 **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    _emm_common_data                           **
 **                                                                        **
 ***************************************************************************/
void emm_common_cleanup(emm_common_data_t* emm_common_data_ctx) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  if (emm_common_data_ctx) {
    __sync_fetch_and_sub(&emm_common_data_ctx->ref_count, 1);

    if (emm_common_data_ctx->ref_count == 0) {
      /*
       * Release the callback functions
       */
      pthread_mutex_lock(&emm_common_data_head.mutex);
      RB_REMOVE(
          emm_common_data_map,
          (emm_common_data_map*)&emm_common_data_head.emm_common_data_root,
          emm_common_data_ctx);
      free_wrapper(&emm_common_data_ctx->args);
      free_wrapper((void**)&emm_common_data_ctx);
      pthread_mutex_unlock(&emm_common_data_head.mutex);
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void emm_common_cleanup_by_ueid(mme_ue_s1ap_id_t ue_id) {
  emm_common_data_t* emm_common_data_ctx = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_common_data_ctx =
      emm_common_data_context_get(&emm_common_data_head, ue_id);

  if (emm_common_data_ctx) {
    __sync_fetch_and_sub(&emm_common_data_ctx->ref_count, 1);
    pthread_mutex_lock(&emm_common_data_head.mutex);
    RB_REMOVE(emm_common_data_map,
              (emm_common_data_map*)&emm_common_data_head.emm_common_data_root,
              emm_common_data_ctx);
    if (emm_common_data_ctx->args) {
      free_wrapper(&emm_common_data_ctx->args);
    }
    free_wrapper((void**)&emm_common_data_ctx);
    pthread_mutex_unlock(&emm_common_data_head.mutex);
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void emm_proc_common_clear_args(mme_ue_s1ap_id_t ue_id) {
  emm_common_data_t* emm_common_data_ctx = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_common_data_ctx =
      emm_common_data_context_get(&emm_common_data_head, ue_id);
  if (emm_common_data_ctx && emm_common_data_ctx->args) {
    free_wrapper(&emm_common_data_ctx->args);
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void create_new_attach_info(emm_context_t* emm_context_p,
                            mme_ue_s1ap_id_t mme_ue_s1ap_id,
                            STOLEN_REF struct emm_attach_request_ies_s* ies,
                            bool is_mm_ctx_new) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_context_p->new_attach_info = reinterpret_cast<new_attach_info_t*>(
      calloc(1, sizeof(new_attach_info_t)));
  emm_context_p->new_attach_info->mme_ue_s1ap_id = mme_ue_s1ap_id;
  emm_context_p->new_attach_info->ies = ies;
  emm_context_p->new_attach_info->is_mm_ctx_new = is_mm_ctx_new;
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:        emm_verify_orig_tai()                                     **
 **                                                                        **
 ** Description: Verifies if the TAI received in s1ap message              **
 **              is configured                                             **
 **                                                                        **
 ** Inputs:      orig_tai: TAI received in the s1ap message                **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    matching partial list, NULL                            **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
partial_list_t* emm_verify_orig_tai(const tai_t orig_tai) {
  partial_list_t* par_list = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  if (!mme_config.partial_list) {
    OAILOG_ERROR(LOG_NAS_EMM, "partial_list in mme_config is NULL\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, par_list);
  }

  for (uint8_t list_i = 0; list_i < mme_config.num_par_lists; list_i++) {
    for (uint8_t elem_i = 0; elem_i < mme_config.partial_list[list_i].nb_elem;
         elem_i++) {
      if (((mme_config.partial_list[list_i].plmn) &&
           (IS_PLMN_EQUAL(orig_tai.plmn,
                          mme_config.partial_list[list_i].plmn[elem_i]))) &&
          (mme_config.partial_list[list_i].tac &&
           (orig_tai.tac == mme_config.partial_list[list_i].tac[elem_i]))) {
        par_list = &mme_config.partial_list[list_i];
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, par_list);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, par_list);
}

/****************************************************************************
 **                                                                        **
 ** Name:        update_tai_list_to_emm_context                            **
 **                                                                        **
 ** Description: Updates the new TAI list to emm context                   **
 **                                                                        **
 ** Inputs:      imsi64, guti                                              **
 **              par_tai_list: pointer to the matching partial_list_t      **
 **                            in mme_config                               **
 **              tai_list: pointer to tai_list_t in emm context            **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
status_code_e update_tai_list_to_emm_context(
    uint64_t imsi64, guti_t guti, const partial_list_t* const par_tai_list,
    tai_list_t* tai_list) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  if (!par_tai_list->plmn) {
    OAILOG_ERROR_UE(LOG_NAS, imsi64, "config PLMN is NULL\n");
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }
  if (!par_tai_list->tac) {
    OAILOG_ERROR_UE(LOG_NAS, imsi64, "config TAC is NULL\n");
    OAILOG_FUNC_RETURN(LOG_NAS, RETURNerror);
  }

  OAILOG_INFO_UE(
      LOG_NAS, imsi64,
      "Matching partial list for originating TAI found! typeOfList=%d\n",
      par_tai_list->list_type);
  int itr = 0;
  /* Comparing PLMN of mme configuration with PLMN of GUMMEI_LIST.
   * If PLMN matches, TAI_LIST in emm_context gets updated with TAI_LIST
   * values from mme configuration file
   */
  switch (par_tai_list->list_type) {
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_NON_CONSECUTIVE_TACS:
      if (IS_PLMN_EQUAL(par_tai_list->plmn[0], guti.gummei.plmn)) {
        /* As per 3gpp spec 24.301 sec-9.9.3.33, numberofelements=0
         * corresponds to 1 element,
         * numberofelements=1 corresponds to 2 elements ...
         * So set numberofelements = nb_elem - 1
         */
        tai_list->partial_tai_list[0].numberofelements =
            par_tai_list->nb_elem - 1;
        tai_list->partial_tai_list[0].typeoflist = par_tai_list->list_type;
        COPY_PLMN(tai_list->partial_tai_list[0]
                      .u.tai_one_plmn_non_consecutive_tacs.plmn,
                  guti.gummei.plmn);

        // par_tai_list is sorted
        for (itr = 0; itr < (par_tai_list->nb_elem); itr++) {
          tai_list->partial_tai_list[0]
              .u.tai_one_plmn_non_consecutive_tacs.tac[itr] =
              par_tai_list->tac[itr];
        }
      } else {
        OAILOG_ERROR_UE(
            LOG_NAS, imsi64,
            "GUTI PLMN does not match with mme configuration tai list\n");
      }
      break;
    case TRACKING_AREA_IDENTITY_LIST_ONE_PLMN_CONSECUTIVE_TACS:
      if (IS_PLMN_EQUAL(par_tai_list->plmn[0], guti.gummei.plmn)) {
        tai_list->partial_tai_list[0].numberofelements =
            par_tai_list->nb_elem - 1;
        tai_list->partial_tai_list[0].typeoflist = par_tai_list->list_type;

        COPY_PLMN(
            tai_list->partial_tai_list[0].u.tai_one_plmn_consecutive_tacs.plmn,
            guti.gummei.plmn);

        tai_list->partial_tai_list[0].u.tai_one_plmn_consecutive_tacs.tac =
            par_tai_list->tac[0];
      } else {
        OAILOG_ERROR_UE(
            LOG_NAS, imsi64,
            "GUTI PLMN does not match with mme configuration tai list\n");
      }
      break;
    case TRACKING_AREA_IDENTITY_LIST_MANY_PLMNS:
      /* Include all the TAIs as we do not support equivalent PLMN list.
       * Once equivalent PLMN list is supported,check if the TAI PLMNs are
       * present in equivalent PLMN list
       */
      tai_list->partial_tai_list[0].numberofelements =
          par_tai_list->nb_elem - 1;
      tai_list->partial_tai_list[0].typeoflist = par_tai_list->list_type;

      for (itr = 0; itr < (par_tai_list->nb_elem); itr++) {
        COPY_PLMN(tai_list->partial_tai_list[0].u.tai_many_plmn[itr].plmn,
                  par_tai_list->plmn[itr]);

        // partial tai_list is sorted
        tai_list->partial_tai_list[0].u.tai_many_plmn[itr].tac =
            par_tai_list->tac[itr];
      }
      break;
    default:
      OAILOG_ERROR_UE(LOG_NAS, imsi64,
                      "BAD TAI list configuration, unknown TAI list type %u",
                      par_tai_list->list_type);
  }

  /* TS 124.301 V15.4.0 Section 9.9.3.33:
   * "The Tracking area identity list is a type 4 information element,
   * with a minimum length of 8 octets and a maximum length of 98 octets.
   * The list can contain a maximum of 16 different tracking area identities."
   * We will limit the number to 1 partial list which can have maximum of 16
   * TAIs.
   */
  tai_list->numberoflists = 1;
  OAILOG_INFO_UE(LOG_NAS, imsi64,
                 "  Got GUTI " GUTI_FMT ". The number of TAI partial lists: %d",
                 GUTI_ARG(&guti), tai_list->numberoflists);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:        verify_tau_tai                                            **
 **                                                                        **
 ** Description: Verifies if the TAI received during TAU                   **
 **              is configured                                             **
 **                                                                        **
 ** Inputs:      imsi,guti,                                                **
 **              tai: TAI received in TAU                                  **
 **              emm_ctx_tai: pointer to tai_list_t in emm context         **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
status_code_e verify_tau_tai(uint64_t imsi64, guti_t guti, tai_t tai,
                             tai_list_t* emm_ctx_tai) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /* Check if the TAI matches with the TAI in emm context.
   * If it does not match, check if it matches with one of the partial
   * TAI lists stored in mme_config and update the new TAI to emm context.
   * Note that there is only one partial list stored in emm context.
   */
  switch (emm_ctx_tai->partial_tai_list[0].typeoflist) {
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_CONSECUTIVE_TACS:
      if ((IS_PLMN_EQUAL(tai.plmn,
                         emm_ctx_tai->partial_tai_list[0]
                             .u.tai_one_plmn_consecutive_tacs.plmn)) &&
          ((tai.tac >= emm_ctx_tai->partial_tai_list[0]
                           .u.tai_one_plmn_consecutive_tacs.tac) &&
           (tai.tac <= (emm_ctx_tai->partial_tai_list[0]
                            .u.tai_one_plmn_consecutive_tacs.tac +
                        emm_ctx_tai->partial_tai_list[0].numberofelements)))) {
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      }
      break;
    case TRACKING_AREA_IDENTITY_LIST_TYPE_ONE_PLMN_NON_CONSECUTIVE_TACS:
      if (IS_PLMN_EQUAL(tai.plmn,
                        emm_ctx_tai->partial_tai_list[0]
                            .u.tai_one_plmn_non_consecutive_tacs.plmn)) {
        for (uint8_t idx = 0;
             idx < emm_ctx_tai->partial_tai_list[0].numberofelements; idx++) {
          if (tai.tac == emm_ctx_tai->partial_tai_list[0]
                             .u.tai_one_plmn_non_consecutive_tacs.tac[idx]) {
            OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
          }
        }
      }
      break;
    case TRACKING_AREA_IDENTITY_LIST_TYPE_MANY_PLMNS:
      for (uint8_t idx = 0;
           idx < emm_ctx_tai->partial_tai_list[0].numberofelements; idx++) {
        if ((IS_PLMN_EQUAL(
                tai.plmn,
                emm_ctx_tai->partial_tai_list[0].u.tai_many_plmn[idx].plmn)) &&
            (tai.tac ==
             emm_ctx_tai->partial_tai_list[0].u.tai_many_plmn[idx].tac)) {
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        }
      }
      break;
    default:
      OAILOG_ERROR_UE(LOG_NAS_EMM, imsi64,
                      "Unknown TAI list type in verify_tai"
                      "%u\n",
                      emm_ctx_tai->partial_tai_list[0].typeoflist);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  // Check if the TAI matches with the partial lists in mme_config
  partial_list_t* par_list = emm_verify_orig_tai(tai);
  if (par_list) {
    /* Update the new partial list to emm_context. For now, emm context
     * contains only one partial list
     */
    if (update_tai_list_to_emm_context(imsi64, guti, par_list, emm_ctx_tai) ==
        RETURNok) {
      OAILOG_DEBUG_UE(LOG_NAS_EMM, imsi64,
                      "New TAI list updated to emm context\n");
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
    }
  }
  OAILOG_ERROR_UE(LOG_NAS_EMM, imsi64,
                  " Verification of TAI received in TAU failed\n");
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
}
