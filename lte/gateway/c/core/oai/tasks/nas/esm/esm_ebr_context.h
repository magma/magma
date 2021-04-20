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
Source      esm_ebr_context.h

Version     0.1

Date        2013/05/28

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines functions used to handle EPS bearer contexts.

*****************************************************************************/
#ifndef ESM_EBR_CONTEXT_SEEN
#define ESM_EBR_CONTEXT_SEEN

#include <stdbool.h>

#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "common_types.h"
#include "emm_data.h"
#include "esm_data.h"
/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

#define IS_DEFAULT_BEARER_YES true
#define IS_DEFAULT_BEARER_NO false
ebi_t esm_ebr_context_create(
    emm_context_t* emm_context, const proc_tid_t pti, pdn_cid_t pid, ebi_t ebi,
    bool is_default, const qci_t qci, const bitrate_t gbr_dl,
    const bitrate_t gbr_ul, const bitrate_t mbr_dl, const bitrate_t mbr_ul,
    traffic_flow_template_t* tft, protocol_configuration_options_t* pco,
    fteid_t* sgw_fteid);

void esm_ebr_context_init(esm_ebr_context_t* esm_ebr_context);

ebi_t esm_ebr_context_release(
    emm_context_t* emm_context, ebi_t ebi, pdn_cid_t* pid, int* bid);

void free_esm_ebr_context(esm_ebr_context_t* ctx);

void default_eps_bearer_activate_t3485_handler(void* args, imsi64_t* imsi64);

void dedicated_eps_bearer_activate_t3485_handler(void* args, imsi64_t* imsi64);

void eps_bearer_deactivate_t3495_handler(void*, imsi64_t* imsi64);
#endif /* ESM_EBR_CONTEXT_SEEN */
