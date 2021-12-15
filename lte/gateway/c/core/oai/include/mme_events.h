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
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EmmCause.h"
/**
 * Helper function to initiate AsyncEventdClient in its own thread
 */
void event_client_init(void);

/**
 * Logs Attach successful event
 * @param imsi
 * @param ue_id
 * @return response code
 */
int attach_success_event(imsi64_t imsi64, mme_ue_s1ap_id_t ue_id);

/**
 * Logs Attach reject event
 * @param ue_id
 * @param cause
 * @param imsi
 * @return response code
 */
int attach_reject_event(mme_ue_s1ap_id_t ue_id, emm_cause_t cause, imsi64_t imsi64);

/**
 * Logs Detach successful event
 * @param imsi
 * @param action Indicates whether explicit detach accept action was sent to UE
 * @return response code
 */
int detach_success_event(imsi64_t imsi64, const char* action);

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
