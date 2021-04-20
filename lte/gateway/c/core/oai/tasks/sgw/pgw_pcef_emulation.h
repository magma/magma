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

/*! \file pgw_pcef_emulation.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */

#ifndef FILE_PGW_PCEF_EMULATION_SEEN
#define FILE_PGW_PCEF_EMULATION_SEEN

#include <stdbool.h>
#include <stdint.h>

#include "queue.h"
#include "bstrlib.h"
#include "pgw_config.h"
#include "pgw_types.h"
#include "spgw_state.h"

int pgw_pcef_emulation_init(
    spgw_state_t* state_p, const pgw_config_t* pgw_config_p);
void pgw_pcef_emulation_apply_rule(
    spgw_state_t* state_p, sdf_id_t sdf_id, const pgw_config_t* pgw_config_p);
void pgw_pcef_emulation_apply_sdf_filter(
    sdf_filter_t* sdf_f, sdf_id_t sdf_id, const pgw_config_t* pgw_config_p);
bstring pgw_pcef_emulation_packet_filter_2_iptable_string(
    packet_filter_contents_t* packetfiltercontents, uint8_t direction);
int pgw_pcef_get_sdf_parameters(
    spgw_state_t* state, sdf_id_t sdf_id, bearer_qos_t* bearer_qos,
    packet_filter_t* packet_filter, uint8_t* num_pf);

#endif /* FILE_PGW_PCEF_EMULATION_SEEN */
