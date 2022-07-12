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

/*! \file s1ap_mme.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#pragma once
#include "lte/gateway/c/core/oai/include/s1ap_state.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"
#ifdef __cplusplus
extern "C" {
#endif
#if MME_CLIENT_TEST == 0
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#endif

#include "lte/gateway/c/core/common/common_defs.h"

#define S1AP_ZMQ_LATENCY_TH \
  s1ap_zmq_th  // absolute threshold to be used for initial UE messages

extern bool hss_associated;

/** \brief Allocate and add to the list a new eNB descriptor
 * @returns Reference to the new eNB element in list
 **/
enb_description_t* s1ap_new_enb(void);

/** \brief Allocate and add to the right eNB list a new UE descriptor
 * \param sctp_assoc_id association ID over SCTP
 * \param enb_ue_s1ap_id ue ID over S1AP
 * @returns Reference to the new UE element in list
 **/
ue_description_t* s1ap_new_ue(s1ap_state_t* state,
                              sctp_assoc_id_t sctp_assoc_id,
                              enb_ue_s1ap_id_t enb_ue_s1ap_id);

/** \brief Remove target UE from the list
 * \param ue_ref UE structure reference to remove
 **/
void s1ap_remove_ue(s1ap_state_t* state, ue_description_t* ue_ref);

/** \brief Remove target eNB from the list and remove any UE associated
 * \param enb_ref eNB structure reference to remove
 **/
void s1ap_remove_enb(s1ap_state_t* state, enb_description_t* enb_ref);
#ifdef __cplusplus
}
#endif
