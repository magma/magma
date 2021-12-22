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

#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include "lte/gateway/c/core/oai/common/common_types.h"

/**
 * Helper function to initiate AsyncEventdClient in its own thread
 */
void event_client_init(void);

/**
 * Fire event for having received an AttachRequest message
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @return response code
 */
int attach_request_event(
    imsi64_t imsi64, const guti_t guti, const char* imei, const char* mme_id,
    const char* enb_id, const char* enb_ip, const char* apn);

/**
 * Fire event for having sent an AttachAccept message
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn

 * @return response code
 */
int attach_accept_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn);

/**
 * Fire event for having sent an AttachReject message
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @param cause
 * @return response code
 */
int attach_reject_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause);

/**
 * Fire event for having received an AttachComplete message
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @return response code
 */
int attach_complete_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn);

/**
 * Fire event for successfully processing an AttachComplete message.
 * Only fires if a new UE is successfully attached after processing the
 * AttachComplete.
 * If a duplicate AttachComplete was processed for example, and no new UE
 * was attached, then this event should not be fired.
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @return response code
 */
int attach_success_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn);

/**
 * Fire event for when an Attach failure occurs.
 * This should be called, for example, when an AttachRequest was received,
 * but neither an AttachAccept or AttachReject was sent back to the UE.
 * Or, if an AttachComplete was received, but the MME was unable to finish
 * the Attach process.
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @param cause
 * @return response code
 */
int attach_failure_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause);

/**
 * Fire event to record a DetachRequest message being sent/received
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @param source
 * @return response code
 */
int detach_request_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* source);

/**
 * Fire event to record a DetachAccept message being sent/received
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @param source
 * @return response code
 */
int detach_accept_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* source);

/**
 * Fire event for when an implicit Detach has occurred
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @return response code
 */
int detach_implicit_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn);

/**
 * Fire event for a successful Detach
 *
 * @param imsi
 * @param action Indicates whether explicit detach accept action was sent to UE
 * @return response code
 */
int detach_success_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* action);

/**
 * Fire event for a failed Detach
 *
 * @param imsi
 * @param guti
 * @param mme_id
 * @param enb_id
 * @param enb_ip
 * @param apn
 * @param cause
 * @return response code
 */
int detach_failure_event(
    imsi64_t imsi64, const guti_t guti, const char* mme_id, const char* enb_id,
    const char* enb_ip, const char* apn, const char* cause);

/**
 * Logs s1 setup success event
 * @param enb_name name assigned to eNodeb
 * @param enb_id unique identifier of eNodeb
 * @return response code
 */
int s1_setup_success_event(const char* enb_name, uint32_t enb_id);

#ifdef __cplusplus
}
#endif
