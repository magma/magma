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

/*! \file gtpv2c_ie_formatter.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_GTPV2C_IE_FORMATTER_SEEN
#define FILE_GTPV2C_IE_FORMATTER_SEEN

#include <stdbool.h>

#include "NwTypes.h"
#include "NwGtpv2c.h"
#include "3gpp_23.003.h"
#include "3gpp_29.274.h"
#include "sgw_ie_defs.h"
#include "TrafficFlowAggregateDescription.h"
//#include "mme_ie_defs.h"

/* Imsi Information Element
 * 3GPP TS.29.274 #8.3
 * NOTE: Imsi is TBCD coded
 * octet 5   | Number digit 2 | Number digit 1   |
 * octet n+4 | Number digit m | Number digit m-1 |
 */
nw_rc_t gtpv2c_imsi_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_imsi_ie_set(nw_gtpv2c_msg_handle_t* msg, const imsi_t* imsi);

/* Cause Information Element */
nw_rc_t gtpv2c_cause_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_cause_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const gtpv2c_cause_t* cause);

/* Fully Qualified TEID (F-TEID) Information Element */
int gtpv2c_fteid_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const fteid_t* fteid, const uint8_t instance);
nw_rc_t gtpv2c_fteid_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

/* PDN Address Allocation Information Element */
nw_rc_t gtpv2c_paa_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_paa_ie_set(nw_gtpv2c_msg_handle_t* msg, const paa_t* paa);

/* EPS Bearer Id Information Element
 * 3GPP TS 29.274 #8.8
 */
nw_rc_t gtpv2c_ebi_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_ebi_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const unsigned ebi, const uint8_t instance);

/* traffic flow template */
nw_rc_t gtpv2c_tft_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_tft_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const traffic_flow_template_t* const tft);

/* traffic aggregate description */
nw_rc_t gtpv2c_tad_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_tad_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const traffic_flow_aggregate_description_t* tad);

nw_rc_t gtpv2c_ambr_ie_set(nw_gtpv2c_msg_handle_t* msg, ambr_t* ambr);

nw_rc_t gtpv2c_ambr_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

/* Bearer level Qos Information Element
 * 3GPP TS 29.274 #8.15
 */
nw_rc_t gtpv2c_bearer_qos_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_bearer_qos_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const bearer_qos_t* const bearer_qos);

/* Flow Qos Information Element
 * 3GPP TS 29.274 #8.16
 */
nw_rc_t gtpv2c_flow_qos_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_flow_qos_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const flow_qos_t* flow_qos);

/* Target Identification Information Element
 * 3GPP TS 29.274 #8.51
 */
nw_rc_t gtpv2c_target_identification_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_target_identification_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const target_identification_t* target_identification);

/* Bearer Context Created grouped Information Element */
nw_rc_t gtpv2c_bearer_context_created_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_bearer_context_created_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const bearer_context_created_t* const bearer);

/* Bearer Contexts to Create Information Element
 * 3GPP TS 29.274 Table 7.2.1-2.
 */
nw_rc_t gtpv2c_bearer_context_to_be_created_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_bearer_context_to_be_created_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_to_be_created_t* bearer_context);

/* Selection Mode
 * 3GPP TS 29.274 #8.58
 */

nw_rc_t gtpv2c_selection_mode_ie_set(
    nw_gtpv2c_msg_handle_t* msg, SelectionMode_t* sm);

nw_rc_t gtpv2c_selection_mode_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

#endif /* FILE_GTPV2C_IE_FORMATTER_SEEN */
