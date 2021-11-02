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
Source    emm_msg.h

Version   0.1

Date    2012/09/27

Product   NAS stack

Subsystem EPS Mobility Management

Author    Frederic Maurel

Description Defines EPS Mobility Management messages and functions used
    to encode and decode

*****************************************************************************/
#ifndef FILE_EMM_MSG_SEEN
#define FILE_EMM_MSG_SEEN

#include <stdint.h>

#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/emm_msgDef.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateResult.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/Cli.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/CsfbResponse.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/DetachType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EmmCause.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsAttachResult.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsAttachType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsBearerContextStatus.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsMobileIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsNetworkFeatureSupport.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsUpdateResult.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsUpdateType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EsmMessageContainer.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/GutiType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/KsiAndSequenceNumber.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/LcsClientIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/LcsIndicator.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/MessageType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasKeySetIdentifier.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasMessageContainer.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasSecurityAlgorithms.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/Nonce.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/PagingIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/SecurityHeaderType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/ServiceType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/ShortMac.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/SsCode.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeNetworkCapability.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeRadioCapabilityInformationUpdateNeeded.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeSecurityCapability.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachComplete.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationFailure.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationResponse.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/CsServiceNotification.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DetachAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DetachRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DownlinkNasTransport.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmInformation.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmStatus.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ExtendedServiceRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/GutiReallocationCommand.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/GutiReallocationComplete.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/IdentityRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/IdentityResponse.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeCommand.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeComplete.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/SecurityModeReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ServiceReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ServiceRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateAccept.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateComplete.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateReject.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/UplinkNasTransport.h"

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
int emm_msg_decode_header(
    emm_msg_header_t* header, const uint8_t* buffer, uint32_t len);

int emm_msg_decode(EMM_msg* msg, uint8_t* buffer, uint32_t len);

int emm_msg_encode(EMM_msg* msg, uint8_t* buffer, uint32_t len);

int emm_msg_encode_header(
    const emm_msg_header_t* header, uint8_t* buffer, uint32_t len);

#endif /* FILE_EMM_MSG_SEEN */
