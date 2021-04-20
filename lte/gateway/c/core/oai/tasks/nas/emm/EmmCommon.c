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
  Source      EmmCommon.h

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
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <pthread.h>

#include "dynamic_memory_check.h"
#include "assertions.h"
#include "common_defs.h"
#include "log.h"
#include "emm_data.h"
#include "EmmCommon.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

emm_common_data_head_t emm_common_data_head = {PTHREAD_MUTEX_INITIALIZER,
                                               RB_INITIALIZER()};

static inline int emm_common_data_compare_ueid(
    struct emm_common_data_s* p1, struct emm_common_data_s* p2);

RB_PROTOTYPE(
    emm_common_data_map, emm_common_data_s, entries,
    emm_common_data_compare_ueid);

/* Generate functions used for the MAP */
RB_GENERATE(
    emm_common_data_map, emm_common_data_s, entries,
    emm_common_data_compare_ueid);

static inline int emm_common_data_compare_ueid(
    struct emm_common_data_s* p1, struct emm_common_data_s* p2) {
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
      RB_FIND(emm_common_data_map, &root->emm_common_data_root, &reference);
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
int emm_proc_common_initialize(
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
    emm_common_data_ctx =
        (emm_common_data_t*) calloc(1, sizeof(emm_common_data_t));
    emm_common_data_ctx->ue_id = ue_id;
    pthread_mutex_lock(&emm_common_data_head.mutex);
    RB_INSERT(
        emm_common_data_map, &emm_common_data_head.emm_common_data_root,
        emm_common_data_ctx);
    pthread_mutex_unlock(&emm_common_data_head.mutex);

    if (emm_common_data_ctx) {
      emm_common_data_ctx->ref_count = 0;
    }
  }

  if (emm_common_data_ctx) {
    __sync_fetch_and_add(&emm_common_data_ctx->ref_count, 1);
    emm_common_data_ctx->success       = _success;
    emm_common_data_ctx->reject        = _reject;
    emm_common_data_ctx->failure       = _failure;
    emm_common_data_ctx->ll_failure    = _ll_failure;
    emm_common_data_ctx->non_delivered = _non_delivered;
    emm_common_data_ctx->abort         = _abort;
    emm_common_data_ctx->args          = args;
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
int emm_proc_common_success(emm_common_data_t* emm_common_data_ctx) {
  emm_common_success_callback_t emm_callback = {0};
  int rc                                     = RETURNerror;

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
int emm_proc_common_reject(emm_common_data_t* emm_common_data_ctx) {
  int rc = RETURNerror;
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

int emm_proc_common_failure(emm_common_data_t* emm_common_data_ctx) {
  int rc = RETURNerror;
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
int emm_proc_common_ll_failure(emm_common_data_t* emm_common_data_ctx) {
  emm_common_ll_failure_callback_t emm_callback;
  int rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->ll_failure;

    if (emm_callback) {
      struct emm_context_s* ctx = NULL;

      ctx = emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc  = (*emm_callback)(ctx);
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
int emm_proc_common_non_delivered(emm_common_data_t* emm_common_data_ctx) {
  emm_common_non_delivered_callback_t emm_callback;
  int rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->non_delivered;

    if (emm_callback) {
      struct emm_context_s* ctx = NULL;

      ctx = emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc  = (*emm_callback)(ctx);
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
int emm_proc_common_abort(emm_common_data_t* emm_common_data_ctx) {
  emm_common_abort_callback_t emm_callback;
  int rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (emm_common_data_ctx) {
    emm_callback = emm_common_data_ctx->abort;

    if (emm_callback) {
      struct emm_context_s* ctx = NULL;

      ctx = emm_context_get(&_emm_data, emm_common_data_ctx->ue_id);
      rc  = (*emm_callback)(ctx);
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
  if (emm_common_data_ctx) {
    __sync_fetch_and_sub(&emm_common_data_ctx->ref_count, 1);

    if (emm_common_data_ctx->ref_count == 0) {
      /*
       * Release the callback functions
       */
      pthread_mutex_lock(&emm_common_data_head.mutex);
      RB_REMOVE(
          emm_common_data_map, &emm_common_data_head.emm_common_data_root,
          emm_common_data_ctx);
      free_wrapper(&emm_common_data_ctx->args);
      free_wrapper((void**) &emm_common_data_ctx);
      pthread_mutex_unlock(&emm_common_data_head.mutex);
    }
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void emm_common_cleanup_by_ueid(mme_ue_s1ap_id_t ue_id) {
  emm_common_data_t* emm_common_data_ctx = NULL;

  emm_common_data_ctx =
      emm_common_data_context_get(&emm_common_data_head, ue_id);

  if (emm_common_data_ctx) {
    __sync_fetch_and_sub(&emm_common_data_ctx->ref_count, 1);
    pthread_mutex_lock(&emm_common_data_head.mutex);
    RB_REMOVE(
        emm_common_data_map, &emm_common_data_head.emm_common_data_root,
        emm_common_data_ctx);
    if (emm_common_data_ctx->args) {
      free_wrapper(&emm_common_data_ctx->args);
    }
    free_wrapper((void**) &emm_common_data_ctx);
    pthread_mutex_unlock(&emm_common_data_head.mutex);
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

void emm_proc_common_clear_args(mme_ue_s1ap_id_t ue_id) {
  emm_common_data_t* emm_common_data_ctx = NULL;
  emm_common_data_ctx =
      emm_common_data_context_get(&emm_common_data_head, ue_id);
  if (emm_common_data_ctx && emm_common_data_ctx->args) {
    free_wrapper(&emm_common_data_ctx->args);
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}
