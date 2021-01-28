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

#include "emm_msgDef.h"
#include "AdditionalUpdateResult.h"
#include "AdditionalUpdateType.h"
#include "Cli.h"
#include "CsfbResponse.h"
#include "DetachType.h"
#include "EmmCause.h"
#include "EpsAttachResult.h"
#include "EpsAttachType.h"
#include "EpsBearerContextStatus.h"
#include "EpsMobileIdentity.h"
#include "EpsNetworkFeatureSupport.h"
#include "EpsUpdateResult.h"
#include "EpsUpdateType.h"
#include "EsmMessageContainer.h"
#include "GutiType.h"
#include "KsiAndSequenceNumber.h"
#include "LcsClientIdentity.h"
#include "LcsIndicator.h"
#include "MessageType.h"
#include "NasKeySetIdentifier.h"
#include "NasMessageContainer.h"
#include "NasSecurityAlgorithms.h"
#include "Nonce.h"
#include "PagingIdentity.h"
#include "SecurityHeaderType.h"
#include "ServiceType.h"
#include "ShortMac.h"
#include "SsCode.h"
#include "TrackingAreaIdentity.h"
#include "TrackingAreaIdentityList.h"
#include "3gpp_24.301.h"
#include "UeNetworkCapability.h"
#include "UeRadioCapabilityInformationUpdateNeeded.h"
#include "UeSecurityCapability.h"
#include "AttachAccept.h"
#include "AttachComplete.h"
#include "AttachReject.h"
#include "AttachRequest.h"
#include "AuthenticationFailure.h"
#include "AuthenticationReject.h"
#include "AuthenticationRequest.h"
#include "AuthenticationResponse.h"
#include "CsServiceNotification.h"
#include "DetachAccept.h"
#include "DetachRequest.h"
#include "DownlinkNasTransport.h"
#include "EmmInformation.h"
#include "EmmStatus.h"
#include "ExtendedServiceRequest.h"
#include "GutiReallocationCommand.h"
#include "GutiReallocationComplete.h"
#include "IdentityRequest.h"
#include "IdentityResponse.h"
#include "NASSecurityModeCommand.h"
#include "NASSecurityModeComplete.h"
#include "SecurityModeReject.h"
#include "ServiceReject.h"
#include "ServiceRequest.h"
#include "TrackingAreaUpdateAccept.h"
#include "TrackingAreaUpdateComplete.h"
#include "TrackingAreaUpdateReject.h"
#include "TrackingAreaUpdateRequest.h"
#include "UplinkNasTransport.h"

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
