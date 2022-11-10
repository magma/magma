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

#include "lte/protos/oai/s1ap_state.pb.h"

#ifdef __cplusplus
extern "C" {
#endif
#if MME_CLIENT_TEST == 0
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#endif

#include "lte/gateway/c/core/common/common_defs.h"

#define S1AP_ZMQ_LATENCY_TH \
  s1ap_zmq_th  // absolute threshold to be used for initial UE messages

#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/s1ap_state.hpp"
#include "lte/gateway/c/core/oai/include/s1ap_types.hpp"

extern bool hss_associated;

namespace magma {
namespace lte {

/** \Initialize an object of EnbDescription, which is passed as
 * argument
 **/
void s1ap_new_enb(oai::EnbDescription* enb_ref);

/** \brief Allocate and add to the right eNB list a new UE descriptor
 * \param sctp_assoc_id association ID over SCTP
 * \param enb_ue_s1ap_id ue ID over S1AP
 * @returns Reference to the new UE element in list
 **/
oai::UeDescription* s1ap_new_ue(oai::EnbDescription* enb_ref,
                                sctp_assoc_id_t sctp_assoc_id,
                                enb_ue_s1ap_id_t enb_ue_s1ap_id);

/** \brief Remove target UE from the list
 * \param ue_ref UE structure reference to remove
 **/
void s1ap_remove_ue(oai::S1apState* state, oai::UeDescription* ue_ref);

/** \brief Remove target eNB from the list and remove any UE associated
 * \param enb_ref eNB structure reference to remove
 **/
void s1ap_remove_enb(oai::S1apState* state, oai::EnbDescription* enb_ref);

void free_enb_description(void** ptr);

void free_ue_description(void** ptr);

}  // namespace lte
}  // namespace magma
