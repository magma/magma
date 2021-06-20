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

Source      emm_recv.h

Version     0.1

Date        2013/01/30

Product     NAS stack

Subsystem   EPS Mobility Management

Author      Frederic Maurel

Description Defines functions executed at the EMMAS Service Access
        Point upon receiving EPS Mobility Management messages
        from the Access Stratum sublayer.

*****************************************************************************/
#ifndef FILE_EMM_RECV_SEEN
#define FILE_EMM_RECV_SEEN

#include <stdbool.h>

#include "EmmStatus.h"
#include "DetachRequest.h"
#include "DetachAccept.h"
#include "AttachRequest.h"
#include "AttachComplete.h"
#include "TrackingAreaUpdateRequest.h"
#include "TrackingAreaUpdateComplete.h"
#include "ServiceRequest.h"
#include "ExtendedServiceRequest.h"
#include "GutiReallocationComplete.h"
#include "AuthenticationResponse.h"
#include "AuthenticationFailure.h"
#include "IdentityResponse.h"
#include "NASSecurityModeComplete.h"
#include "SecurityModeReject.h"
#include "UplinkNasTransport.h"
#include "nas_message.h"
#include "3gpp_23.003.h"
#include "3gpp_36.401.h"
#include "TrackingAreaIdentity.h"

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
 * Functions executed by the MME upon receiving EMM message from the UE
 * --------------------------------------------------------------------------
 */
int emm_recv_status(
    mme_ue_s1ap_id_t ueid, emm_status_msg* msg, int* emm_cause,
    const nas_message_decode_status_t* const status);

int emm_recv_attach_request(
    const mme_ue_s1ap_id_t ue_id, const tai_t* const originating_tai,
    const ecgi_t* const originating_ecgi, attach_request_msg* const msg,
    const bool is_initial, const bool ctx_is_new, int* const emm_cause,
    const nas_message_decode_status_t* decode_status);

int emm_recv_attach_complete(
    const mme_ue_s1ap_id_t ueid, const attach_complete_msg* msg,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_detach_request(
    mme_ue_s1ap_id_t ueid, const detach_request_msg* msg, const bool is_initial,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_tracking_area_update_request(
    const mme_ue_s1ap_id_t ueid, tracking_area_update_request_msg* const msg,
    const bool is_initial, tac_t tac, int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_service_request(
    mme_ue_s1ap_id_t ueid, const service_request_msg* msg,
    const bool is_initial, int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_initial_ext_service_request(
    mme_ue_s1ap_id_t ue_id, const extended_service_request_msg* msg,
    int* emm_cause, const nas_message_decode_status_t* decode_status);

int emm_recv_ext_service_request(
    mme_ue_s1ap_id_t ue_id, const extended_service_request_msg* msg,
    int* emm_cause, const nas_message_decode_status_t* decode_status);

int emm_recv_identity_response(
    const mme_ue_s1ap_id_t ueid, identity_response_msg* msg,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_authentication_response(
    const mme_ue_s1ap_id_t ueid, authentication_response_msg* msg,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_authentication_failure(
    const mme_ue_s1ap_id_t ueid, authentication_failure_msg* msg,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_security_mode_complete(
    const mme_ue_s1ap_id_t ueid, security_mode_complete_msg* msg,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_security_mode_reject(
    const mme_ue_s1ap_id_t ueid, security_mode_reject_msg* msg,
    int* const emm_cause,
    const nas_message_decode_status_t* const decode_status);

int emm_recv_detach_accept(mme_ue_s1ap_id_t ueid, int* emm_cause);

int emm_recv_tau_complete(
    mme_ue_s1ap_id_t ue_id, const tracking_area_update_complete_msg* msg);

int emm_recv_uplink_nas_transport(
    mme_ue_s1ap_id_t ueid, uplink_nas_transport_msg* msg, int* emm_cause,
    const nas_message_decode_status_t* status);
#endif /* FILE_EMM_RECV_SEEN*/
