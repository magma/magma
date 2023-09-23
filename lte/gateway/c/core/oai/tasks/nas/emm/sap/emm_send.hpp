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

Source      emm_send.hpp

Version     0.1

Date        2013/01/30

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines functions executed at the EMMAS Service Access
        Point to send EPS Mobility Management messages to the
        Access Stratum sublayer.

*****************************************************************************/
#pragma once

#include <stdint.h>

#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AttachReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/AuthenticationRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/CsServiceNotification.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DetachAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DetachRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/DownlinkNasTransport.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmInformation.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/EmmStatus.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/GutiReallocationCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/IdentityRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/NASSecurityModeCommand.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ServiceReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateAccept.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateReject.hpp"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_asDef.hpp"

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
 * Functions executed by the MME to send EMM messages to the UE
 * --------------------------------------------------------------------------
 */
int emm_send_status(const emm_as_status_t*, emm_status_msg*);

int emm_send_detach_accept(const emm_as_data_t*, detach_accept_msg*);

int emm_send_attach_accept(const emm_as_establish_t*, attach_accept_msg*);
int emm_send_attach_accept_dl_nas(const emm_as_data_t* msg, attach_accept_msg*);
int emm_send_attach_reject(const emm_as_establish_t*, attach_reject_msg*);

int emm_send_tracking_area_update_reject(
    const emm_as_establish_t* msg, tracking_area_update_reject_msg* emm_msg);
int emm_send_tracking_area_update_accept(
    const emm_as_establish_t* msg, tracking_area_update_accept_msg* emm_msg);

int emm_send_tracking_area_update_accept_dl_nas(
    const emm_as_data_t* msg, tracking_area_update_accept_msg* emm_msg);

int emm_send_service_reject(const uint8_t emm_cause,
                            service_reject_msg* emm_msg);

int emm_send_identity_request(const emm_as_security_t*, identity_request_msg*);
int emm_send_authentication_request(const emm_as_security_t*,
                                    authentication_request_msg*);
void emm_free_send_authentication_request(authentication_request_msg*);
int emm_send_authentication_reject(authentication_reject_msg*);
int emm_send_security_mode_command(const emm_as_security_t*,
                                   security_mode_command_msg*);
int emm_send_emm_information(const emm_as_data_t* msg,
                             emm_information_msg* emm_msg);
void emm_free_send_emm_information(emm_information_msg* emm_msg);

int emm_send_nw_detach_request(const emm_as_data_t*, nw_detach_request_msg*);

int emm_send_dl_nas_transport(const emm_as_data_t*,
                              downlink_nas_transport_msg*);
void emm_free_send_dl_nas_transport(downlink_nas_transport_msg* emm_msg);

int emm_send_cs_service_notification(const emm_as_data_t* msg,
                                     cs_service_notification_msg* emm_msg);
void emm_free_send_cs_service_notification(
    cs_service_notification_msg* emm_msg);
