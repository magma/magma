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
Source    esm_msg.h

Version   0.1

Date    2012/09/27

Product   NAS stack

Subsystem EPS Session Management

Author    Frederic Maurel

Description Defines EPS Session Management messages and functions used
    to encode and decode

*****************************************************************************/
#ifndef __ESM_MSG_H__
#define __ESM_MSG_H__

#include <stdint.h>

#include "esm_msgDef.h"
#include "EsmInformationTransferFlag.h"
#include "NasRequestType.h"
#include "PdnType.h"
#include "TrafficFlowAggregateDescription.h"
#include "ActivateDedicatedEpsBearerContextRequest.h"
#include "ActivateDedicatedEpsBearerContextAccept.h"
#include "ActivateDedicatedEpsBearerContextReject.h"
#include "ActivateDefaultEpsBearerContextRequest.h"
#include "ActivateDefaultEpsBearerContextAccept.h"
#include "ActivateDefaultEpsBearerContextReject.h"
#include "ModifyEpsBearerContextRequest.h"
#include "ModifyEpsBearerContextAccept.h"
#include "ModifyEpsBearerContextReject.h"
#include "DeactivateEpsBearerContextRequest.h"
#include "DeactivateEpsBearerContextAccept.h"
#include "PdnDisconnectRequest.h"
#include "PdnDisconnectReject.h"
#include "PdnConnectivityRequest.h"
#include "PdnConnectivityReject.h"
#include "BearerResourceAllocationRequest.h"
#include "BearerResourceAllocationReject.h"
#include "BearerResourceModificationRequest.h"
#include "BearerResourceModificationReject.h"
#include "EsmInformationRequest.h"
#include "EsmInformationResponse.h"
#include "EsmStatus.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Structure of ESM plain NAS message
 * ----------------------------------
 */
typedef union {
  esm_msg_header_t header;
  activate_default_eps_bearer_context_request_msg
      activate_default_eps_bearer_context_request;
  activate_default_eps_bearer_context_accept_msg
      activate_default_eps_bearer_context_accept;
  activate_default_eps_bearer_context_reject_msg
      activate_default_eps_bearer_context_reject;
  activate_dedicated_eps_bearer_context_request_msg
      activate_dedicated_eps_bearer_context_request;
  activate_dedicated_eps_bearer_context_accept_msg
      activate_dedicated_eps_bearer_context_accept;
  activate_dedicated_eps_bearer_context_reject_msg
      activate_dedicated_eps_bearer_context_reject;
  modify_eps_bearer_context_request_msg modify_eps_bearer_context_request;
  modify_eps_bearer_context_accept_msg modify_eps_bearer_context_accept;
  modify_eps_bearer_context_reject_msg modify_eps_bearer_context_reject;
  deactivate_eps_bearer_context_request_msg
      deactivate_eps_bearer_context_request;
  deactivate_eps_bearer_context_accept_msg deactivate_eps_bearer_context_accept;
  pdn_connectivity_request_msg pdn_connectivity_request;
  pdn_connectivity_reject_msg pdn_connectivity_reject;
  pdn_disconnect_request_msg pdn_disconnect_request;
  pdn_disconnect_reject_msg pdn_disconnect_reject;
  bearer_resource_allocation_request_msg bearer_resource_allocation_request;
  bearer_resource_allocation_reject_msg bearer_resource_allocation_reject;
  bearer_resource_modification_request_msg bearer_resource_modification_request;
  bearer_resource_modification_reject_msg bearer_resource_modification_reject;
  esm_information_request_msg esm_information_request;
  esm_information_response_msg esm_information_response;
  esm_status_msg esm_status;
} ESM_msg;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
int esm_msg_decode_header(
    esm_msg_header_t* header, const uint8_t* buffer, uint32_t len);

int esm_msg_decode(ESM_msg* msg, uint8_t* buffer, uint32_t len);

int esm_msg_encode(ESM_msg* msg, uint8_t* buffer, uint32_t len);

#endif /* __ESM_MSG_H__ */
