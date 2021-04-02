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
Source      esm_send.h

Version     0.1

Date        2013/02/11

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines functions executed at the ESM Service Access
        Point to send EPS Session Management messages to the
        EPS Mobility Management sublayer.

*****************************************************************************/
#ifndef __ESM_SEND_H__
#define __ESM_SEND_H__

#include "common_defs.h"
#include "EsmStatus.h"
#include "PdnConnectivityReject.h"
#include "PdnDisconnectReject.h"
#include "BearerResourceAllocationReject.h"
#include "BearerResourceModificationReject.h"
#include "ActivateDefaultEpsBearerContextRequest.h"
#include "ActivateDedicatedEpsBearerContextRequest.h"
#include "ModifyEpsBearerContextRequest.h"
#include "DeactivateEpsBearerContextRequest.h"
#include "EsmInformationRequest.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "EpsQualityOfService.h"
#include "bstrlib.h"
#include "mme_app_ue_context.h"

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
 * Functions executed by the MME to send ESM message to the UE
 * --------------------------------------------------------------------------
 */
int esm_send_esm_information_request(
    pti_t pti, ebi_t ebi, esm_information_request_msg* msg);

int esm_send_status(pti_t pti, ebi_t ebi, esm_status_msg* msg, int esm_cause);

/*
 * Transaction related messages
 * ----------------------------
 */
int esm_send_pdn_connectivity_reject(
    pti_t pti, pdn_connectivity_reject_msg* msg, int esm_cause);

int esm_send_pdn_disconnect_reject(
    pti_t pti, pdn_disconnect_reject_msg* msg, int esm_cause);

/*
 * Messages related to EPS bearer contexts
 * ---------------------------------------
 */
int esm_send_activate_default_eps_bearer_context_request(
    pti_t pti, ebi_t ebi, activate_default_eps_bearer_context_request_msg* msg,
    pdn_context_t* pdn_context_p, const protocol_configuration_options_t* pco,
    int pdn_type, bstring pdn_addr, const EpsQualityOfService* qos,
    int esm_cause);

int esm_send_activate_dedicated_eps_bearer_context_request(
    pti_t pti, ebi_t ebi,
    activate_dedicated_eps_bearer_context_request_msg* msg, ebi_t linked_ebi,
    const EpsQualityOfService* qos, traffic_flow_template_t* tft,
    protocol_configuration_options_t* pco);

int esm_send_deactivate_eps_bearer_context_request(
    pti_t pti, ebi_t ebi, deactivate_eps_bearer_context_request_msg* msg,
    int esm_cause);

#endif /* __ESM_SEND_H__*/
