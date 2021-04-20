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
Source      esm_recv.h

Version     0.1

Date        2013/02/06

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines functions executed at the ESM Service Access
        Point upon receiving EPS Session Management messages
        from the EPS Mobility Management sublayer.

*****************************************************************************/
#ifndef __ESM_RECV_H__
#define __ESM_RECV_H__

#include "EsmStatus.h"
#include "emm_data.h"
#include "PdnConnectivityRequest.h"
#include "PdnDisconnectRequest.h"
#include "BearerResourceAllocationRequest.h"
#include "BearerResourceModificationRequest.h"
#include "ActivateDefaultEpsBearerContextAccept.h"
#include "ActivateDefaultEpsBearerContextReject.h"
#include "ActivateDedicatedEpsBearerContextAccept.h"
#include "ActivateDedicatedEpsBearerContextReject.h"
#include "ModifyEpsBearerContextAccept.h"
#include "ModifyEpsBearerContextReject.h"
#include "DeactivateEpsBearerContextAccept.h"
#include "EsmInformationResponse.h"
#include "3gpp_24.007.h"
#include "EsmCause.h"
#include "common_types.h"

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

/*
 * --------------------------------------------------------------------------
 * Functions executed by both the UE and the MME upon receiving ESM messages
 * --------------------------------------------------------------------------
 */

esm_cause_t esm_recv_status(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const esm_status_msg* msg);

/*
 * --------------------------------------------------------------------------
 * Functions executed by the MME upon receiving ESM message from the UE
 * --------------------------------------------------------------------------
 */
/*
 * Transaction related messages
 * ----------------------------
 */
esm_cause_t esm_recv_pdn_connectivity_request(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const pdn_connectivity_request_msg* msg, ebi_t* new_ebi,
    bool is_standalone);

esm_cause_t esm_recv_pdn_disconnect_request(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const pdn_disconnect_request_msg* msg);

esm_cause_t esm_recv_information_response(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const esm_information_response_msg* msg);

/*
 * Messages related to EPS bearer contexts
 * ---------------------------------------
 */
esm_cause_t esm_recv_activate_default_eps_bearer_context_accept(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_default_eps_bearer_context_accept_msg* msg);

esm_cause_t esm_recv_activate_default_eps_bearer_context_reject(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_default_eps_bearer_context_reject_msg* msg);

esm_cause_t esm_recv_activate_dedicated_eps_bearer_context_accept(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_dedicated_eps_bearer_context_accept_msg* msg);

esm_cause_t esm_recv_activate_dedicated_eps_bearer_context_reject(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_dedicated_eps_bearer_context_reject_msg* msg);

esm_cause_t esm_recv_deactivate_eps_bearer_context_accept(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const deactivate_eps_bearer_context_accept_msg* msg);

#endif /* __ESM_RECV_H__*/
