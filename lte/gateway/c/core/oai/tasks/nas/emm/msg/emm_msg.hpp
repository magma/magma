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
Source    emm_msg.hpp

Version   0.1

Date    2012/09/27

Product   NAS stack

Subsystem EPS Mobility Management

Author    Frederic Maurel

Description Defines EPS Mobility Management messages and functions used
    to encode and decode

*****************************************************************************/

#pragma once

#include <stdint.h>

#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/emm_msgDef.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateResult.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/Cli.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/CsfbResponse.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/DetachType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EmmCause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsAttachResult.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsAttachType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsBearerContextStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsMobileIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsNetworkFeatureSupport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsUpdateResult.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsUpdateType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EsmMessageContainer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/GutiType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/KsiAndSequenceNumber.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/LcsClientIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/LcsIndicator.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/MessageType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasKeySetIdentifier.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasMessageContainer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasSecurityAlgorithms.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/Nonce.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/PagingIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/SecurityHeaderType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/ServiceType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/ShortMac.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/SsCode.hpp"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.hpp"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeNetworkCapability.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeRadioCapabilityInformationUpdateNeeded.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeSecurityCapability.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachComplete.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationFailure.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationResponse.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/CsServiceNotification.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DetachAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DetachRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DownlinkNasTransport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmInformation.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ExtendedServiceRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/GutiReallocationCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/GutiReallocationComplete.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/IdentityRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/IdentityResponse.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeComplete.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/SecurityModeReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ServiceReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ServiceRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateComplete.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/UplinkNasTransport.hpp"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * Structure of EMM plain NAS message
 * ----------------------------------
 */
typedef union {
  emm_msg_header_t header;
  attach_request_msg attach_request;
  attach_accept_msg attach_accept;
  attach_complete_msg attach_complete;
  attach_reject_msg attach_reject;
  detach_request_msg detach_request;
  nw_detach_request_msg nw_detach_request;
  detach_accept_msg detach_accept;
  tracking_area_update_request_msg tracking_area_update_request;
  tracking_area_update_accept_msg tracking_area_update_accept;
  tracking_area_update_complete_msg tracking_area_update_complete;
  tracking_area_update_reject_msg tracking_area_update_reject;
  extended_service_request_msg extended_service_request;
  service_request_msg service_request;
  service_reject_msg service_reject;
  guti_reallocation_command_msg guti_reallocation_command;
  guti_reallocation_complete_msg guti_reallocation_complete;
  authentication_request_msg authentication_request;
  authentication_response_msg authentication_response;
  authentication_reject_msg authentication_reject;
  authentication_failure_msg authentication_failure;
  identity_request_msg identity_request;
  identity_response_msg identity_response;
  security_mode_command_msg security_mode_command;
  security_mode_complete_msg security_mode_complete;
  security_mode_reject_msg security_mode_reject;
  emm_status_msg emm_status;
  emm_information_msg emm_information;
  downlink_nas_transport_msg downlink_nas_transport;
  uplink_nas_transport_msg uplink_nas_transport;
  cs_service_notification_msg cs_service_notification;
} EMM_msg;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
int emm_msg_decode_header(emm_msg_header_t* header, const uint8_t* buffer,
                          uint32_t len);
int emm_msg_decode(EMM_msg* msg, uint8_t* buffer, uint32_t len);

int emm_msg_encode(EMM_msg* msg, uint8_t* buffer, uint32_t len);

int emm_msg_encode_header(const emm_msg_header_t* header, uint8_t* buffer,
                          uint32_t len);
