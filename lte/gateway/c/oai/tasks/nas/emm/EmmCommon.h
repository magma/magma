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
#ifndef FILE_EMM_COMMON_SEEN
#define FILE_EMM_COMMON_SEEN
#include <pthread.h>

#include "common_types.h"
#include "tree.h"
#include "3gpp_36.401.h"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Type of EMM procedure callback functions
 * ----------------------------------------
 * EMM procedure to be executed under certain conditions, when an EMM common
 * procedure has been initiated by the ongoing EMM procedure.
 * - The EMM common procedure successfully completed
 * - The EMM common procedure failed or is rejected
 * - Lower layer failure occured before the EMM common procedure completion
 */
typedef int (*emm_common_success_callback_t)(void*);
typedef int (*emm_common_reject_callback_t)(void*);
typedef int (*emm_common_failure_callback_t)(void*);
typedef int (*emm_common_ll_failure_callback_t)(void*);
typedef int (*emm_common_non_delivered_callback_t)(void*);
/* EMM common procedure to be executed when the ongoing EMM procedure is
 * aborted.
 */
typedef int (*emm_common_abort_callback_t)(void*);

/* Ongoing EMM procedure callback functions */
typedef struct emm_common_data_s {
  mme_ue_s1ap_id_t ue_id;
  int ref_count;

  emm_common_success_callback_t success;
  emm_common_reject_callback_t reject;
  emm_common_failure_callback_t failure;

  emm_common_ll_failure_callback_t ll_failure;
  emm_common_non_delivered_callback_t non_delivered;
  emm_common_abort_callback_t abort;

  void* args;
  RB_ENTRY(emm_common_data_s) entries;
} emm_common_data_t;

typedef struct emm_common_data_head_s {
  pthread_mutex_t mutex;
  RB_HEAD(emm_common_data_map, emm_common_data_s) emm_common_data_root;
} emm_common_data_head_t;

extern emm_common_data_head_t emm_common_data_head;
/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

int emm_proc_common_initialize(
    mme_ue_s1ap_id_t ue_id, emm_common_success_callback_t success,
    emm_common_reject_callback_t reject, emm_common_failure_callback_t failure,
    emm_common_ll_failure_callback_t ll_failure,
    emm_common_non_delivered_callback_t non_delivered,
    emm_common_abort_callback_t abort, void* args);

int emm_proc_common_success(emm_common_data_t* emm_common_data_ctx);
int emm_proc_common_reject(emm_common_data_t* emm_common_data_ctx);
int emm_proc_common_failure(emm_common_data_t* emm_common_data_ctx);
int emm_proc_common_ll_failure(emm_common_data_t* emm_common_data_ctx);
int emm_proc_common_non_delivered(emm_common_data_t* emm_common_data_ctx);
int emm_proc_common_abort(emm_common_data_t* emm_common_data_ctx);

void* emm_proc_common_get_args(mme_ue_s1ap_id_t ue_id);
// Free args and set it to NULL
void emm_proc_common_clear_args(mme_ue_s1ap_id_t ue_id);
void emm_common_cleanup(emm_common_data_t* emm_common_data_ctx);
void emm_common_cleanup_by_ueid(mme_ue_s1ap_id_t ue_id);

struct emm_common_data_s* emm_common_data_context_get(
    struct emm_common_data_head_s* root, mme_ue_s1ap_id_t _ueid);

#endif /* FILE_EMM_COMMON_SEEN*/
