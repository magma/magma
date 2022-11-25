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
Source      esm_ebr.hpp

Version     0.1

Date        2013/01/29

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines functions used to handle state of EPS bearer contexts
        and manage ESM messages re-transmission.

*****************************************************************************/
#pragma once

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.hpp"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/* Unassigned EPS bearer identity value */
#define ESM_EBI_UNASSIGNED (EPS_BEARER_IDENTITY_UNASSIGNED)
#define ERAB_SETUP_RSP_COUNTER_MAX 5
// TODO - Make it configurable
#define ERAB_SETUP_RSP_TMR 5  // In secs

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

const char* esm_ebr_state2string(esm_ebr_state esm_ebr_state);

void esm_ebr_initialize(void);
ebi_t esm_ebr_assign(emm_context_t* emm_context);
status_code_e esm_ebr_release(emm_context_t* emm_context, ebi_t ebi);

bool esm_ebr_is_reserved(ebi_t ebi);
status_code_e esm_ebr_start_timer(emm_context_t* emm_context, ebi_t ebi,
                                  CLONE_REF const_bstring msg, uint32_t msec,
                                  time_out_t cb);

bool esm_ebr_is_not_in_use(emm_context_t* emm_context, ebi_t ebi);

status_code_e esm_ebr_stop_timer(emm_context_t* emm_context, ebi_t ebi);
ebi_t esm_ebr_get_pending_ebi(emm_context_t* emm_context, esm_ebr_state status);

status_code_e esm_ebr_set_status(emm_context_t* emm_context, ebi_t ebi,
                                 esm_ebr_state status, bool ue_requested);
esm_ebr_state esm_ebr_get_status(emm_context_t* emm_context, ebi_t ebi);
