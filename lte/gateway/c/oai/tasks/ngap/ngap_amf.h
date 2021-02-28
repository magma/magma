/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

#include "intertask_interface.h"
#include "ngap_types.h"
#include "amf_config.h"

extern bool hss_associated;

/** \brief NGAP layer top init
 * @returns -1 in case of failure
 **/
int ngap_amf_init(const amf_config_t* amf_config);

/** \brief NGAP layer top exit
 **/
void ngap_amf_exit(void);

/** \brief Dump gNB related information.
 * Calls dump_ue for each UE in list
 * \param gnb_ref eNB structure reference to dump
 **/
extern void ngap_dump_gnb(const gnb_description_t* const gnb_ref);

/** \brief Allocate and add to the list a new gNB descriptor
 * @returns Reference to the new gNB element in list
 **/
gnb_description_t* ngap_new_gnb(ngap_state_t* state);

/** \brief Allocate and add to the right eNB list a new UE descriptor
 * \param sctp_assoc_id association ID over SCTP
 * \param gnb_ue_ngap_id ue ID over NGAP
 * @returns Reference to the new UE element in list
 **/

m5g_ue_description_t* ngap_new_ue(
    ngap_state_t* state, const sctp_assoc_id_t sctp_assoc_id,
    gnb_ue_ngap_id_t gnb_ue_ngap_id);

/** \brief Remove target UE from the list
 * \param ue_ref UE structure reference to remove
 **/
void ngap_remove_ue(ngap_state_t* state, m5g_ue_description_t* ue_ref);

/** \brief Remove target gNB from the list and remove any UE associated
 * \param gnb_ref gNB structure reference to remove
 **/
void ngap_remove_gnb(ngap_state_t* state, gnb_description_t* gnb_ref);
