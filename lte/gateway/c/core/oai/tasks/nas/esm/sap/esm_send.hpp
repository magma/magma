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
Source      esm_send.hpp

Version     0.1

Date        2013/02/11

Product     NAS stack

Subsystem   EPS Session Management

Author      Frederic Maurel

Description Defines functions executed at the ESM Service Access
        Point to send EPS Session Management messages to the
        EPS Mobility Management sublayer.

*****************************************************************************/
#pragma once

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/include/EpsQualityOfService.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDedicatedEpsBearerContextRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ActivateDefaultEpsBearerContextRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/BearerResourceAllocationReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/BearerResourceModificationReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/DeactivateEpsBearerContextRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/EsmInformationRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/EsmStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/ModifyEpsBearerContextRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/PdnConnectivityReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/PdnDisconnectReject.hpp"

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

status_code_e esm_send_status(pti_t pti, ebi_t ebi, esm_status_msg* msg,
                              int esm_cause);

/*
 * Transaction related messages
 * ----------------------------
 */
status_code_e esm_send_esm_information_request(
    pti_t pti, ebi_t ebi, esm_information_request_msg* msg);

status_code_e esm_send_pdn_connectivity_reject(pti_t pti,
                                               pdn_connectivity_reject_msg* msg,
                                               int esm_cause);

status_code_e esm_send_pdn_disconnect_reject(pti_t pti,
                                             pdn_disconnect_reject_msg* msg,
                                             int esm_cause);

/*
 * Messages related to EPS bearer contexts
 * ---------------------------------------
 */
status_code_e esm_send_activate_default_eps_bearer_context_request(
    pti_t pti, ebi_t ebi, activate_default_eps_bearer_context_request_msg* msg,
    pdn_context_t* pdn_context_p, const protocol_configuration_options_t* pco,
    int pdn_type, bstring pdn_addr, const EpsQualityOfService* qos,
    int esm_cause);

status_code_e esm_send_activate_dedicated_eps_bearer_context_request(
    pti_t pti, ebi_t ebi,
    activate_dedicated_eps_bearer_context_request_msg* msg, ebi_t linked_ebi,
    const EpsQualityOfService* qos, traffic_flow_template_t* tft,
    protocol_configuration_options_t* pco);

status_code_e esm_send_deactivate_eps_bearer_context_request(
    pti_t pti, ebi_t ebi, deactivate_eps_bearer_context_request_msg* msg,
    int esm_cause);
