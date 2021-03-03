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
Source      esm_proc.h

Version     0.1

Date        2013/01/02

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines the EPS Session Management procedures executed at
        the ESM Service Access Points.

*****************************************************************************/
#ifndef __ESM_PROC_H__
#define __ESM_PROC_H__

#include "3gpp_29.274.h"
#include "common_defs.h"
#include "nas/networkDef.h"
#include "EsmCause.h"
#include "emm_data.h"
#include "NasRequestType.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/* Type of PDN address */
typedef enum {
  ESM_PDN_TYPE_IPV4   = NET_PDN_TYPE_IPV4,
  ESM_PDN_TYPE_IPV6   = NET_PDN_TYPE_IPV6,
  ESM_PDN_TYPE_IPV4V6 = NET_PDN_TYPE_IPV4V6
} esm_proc_pdn_type_t;

/* Type of PDN request */
typedef enum {
  ESM_PDN_REQUEST_INITIAL = 1,
  ESM_PDN_REQUEST_HANDOVER,
  ESM_PDN_REQUEST_EMERGENCY
} esm_proc_pdn_request_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

struct emm_context_s;
/*
 * Type of the ESM procedure callback executed when requested by the UE
 * or initiated by the network
 */
typedef int (*esm_proc_procedure_t)(
    const bool, struct emm_context_s* const, const ebi_t, bstring*, const bool);

/* PDN connection and EPS bearer context data */
typedef struct esm_proc_data_s {
  proc_tid_t pti;
  request_type_t request_type;
  bstring apn;
  pdn_cid_t pdn_cid;
  esm_proc_pdn_type_t pdn_type;
  bstring pdn_addr;
  bearer_qos_t bearer_qos;
  protocol_configuration_options_t pco;
} esm_proc_data_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
 * --------------------------------------------------------------------------
 *              ESM status procedure
 * --------------------------------------------------------------------------
 */

int esm_proc_status_ind(
    emm_context_t* emm_context, const proc_tid_t pti, ebi_t ebi,
    esm_cause_t* esm_cause);
int esm_proc_status(
    const bool is_standalone, emm_context_t* const emm_context, const ebi_t ebi,
    bstring* msg, const bool sent_by_ue);

/*
 * --------------------------------------------------------------------------
 *          PDN connectivity procedure
 * --------------------------------------------------------------------------
 */

int esm_proc_pdn_connectivity_request(
    emm_context_t* emm_context, const proc_tid_t pti, const pdn_cid_t pdn_cid,
    const context_identifier_t context_identifier,
    const esm_proc_pdn_request_t request_type, const_bstring const apn,
    pdn_type_t pdn_type, const_bstring const pdn_addr,
    bearer_qos_t* default_qos, protocol_configuration_options_t* const pco,
    esm_cause_t* esm_cause);

int esm_proc_pdn_connectivity_reject(
    bool is_standalone, emm_context_t* emm_context, ebi_t ebi, bstring* msg,
    bool ue_triggered);
int esm_proc_pdn_connectivity_failure(
    emm_context_t* emm_context, pdn_cid_t pid);

/*
 * --------------------------------------------------------------------------
 *              PDN disconnect procedure
 * --------------------------------------------------------------------------
 */

int esm_proc_pdn_disconnect_request(
    emm_context_t* emm_context, const proc_tid_t pti, esm_cause_t* esm_cause);

int esm_proc_pdn_disconnect_accept(
    emm_context_t* emm_context, pdn_cid_t pid, esm_cause_t* esm_cause);
int esm_proc_pdn_disconnect_reject(
    const bool is_standalone, emm_context_t* emm_context, ebi_t ebi,
    bstring* msg, const bool ue_triggered);

/*
 * --------------------------------------------------------------------------
 *              ESM information procedure
 * --------------------------------------------------------------------------
 */

int esm_proc_esm_information_request(
    emm_context_t* const emm_context_p, const pti_t pti);

int esm_proc_esm_information_response(
    emm_context_t* emm_context_p, pti_t pti, const_bstring const apn,
    const protocol_configuration_options_t* const pco,
    esm_cause_t* const esm_cause);

/*
 * --------------------------------------------------------------------------
 *      Default EPS bearer context activation procedure
 * --------------------------------------------------------------------------
 */
int esm_proc_default_eps_bearer_context(
    emm_context_t* emm_context, const proc_tid_t pti, pdn_cid_t pid, ebi_t* ebi,
    const qci_t qci, esm_cause_t* esm_cause);
int esm_proc_default_eps_bearer_context_request(
    bool is_standalone, emm_context_t* const emm_context, const ebi_t ebi,
    STOLEN_REF bstring* msg, const bool ue_triggered);
int esm_proc_default_eps_bearer_context_failure(
    emm_context_t* emm_context, pdn_cid_t* const pid);

int esm_proc_default_eps_bearer_context_accept(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause);
int esm_proc_default_eps_bearer_context_reject(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause);

/*
 * --------------------------------------------------------------------------
 *      Dedicated EPS bearer context activation procedure
 * --------------------------------------------------------------------------
 */
int esm_proc_dedicated_eps_bearer_context(
    emm_context_t* emm_context, const proc_tid_t pti, pdn_cid_t pid, ebi_t* ebi,
    ebi_t* default_ebi, const qci_t qci, const bitrate_t gbr_dl,
    const bitrate_t gbr_ul, const bitrate_t mbr_dl, const bitrate_t mbr_ul,
    traffic_flow_template_t* tft, protocol_configuration_options_t* pco,
    fteid_t* sgw_fteid, esm_cause_t* esm_cause);

int esm_proc_dedicated_eps_bearer_context_request(
    const bool is_standalone, emm_context_t* const emm_context, const ebi_t ebi,
    STOLEN_REF bstring* msg, const bool ue_triggered);

int esm_proc_dedicated_eps_bearer_context_accept(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause);
int esm_proc_dedicated_eps_bearer_context_reject(
    emm_context_t* emm_context, ebi_t ebi);

/*
 * --------------------------------------------------------------------------
 *      EPS bearer context deactivation procedure
 * --------------------------------------------------------------------------
 */
int esm_proc_eps_bearer_context_deactivate(
    emm_context_t* const emm_context_p, const bool is_local, const ebi_t ebi,
    pdn_cid_t* pid, int* const bidx, esm_cause_t* const esm_cause);
int esm_proc_eps_bearer_context_deactivate_request(
    const bool is_standalone, emm_context_t* const emm_context, const ebi_t ebi,
    STOLEN_REF bstring* msg, const bool ue_triggered);
pdn_cid_t esm_proc_eps_bearer_context_deactivate_accept(
    emm_context_t* emm_context, ebi_t ebi, esm_cause_t* esm_cause);

#endif /* __ESM_PROC_H__*/
