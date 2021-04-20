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

/*! \file s11_ie_formatter.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#include "sgw_ie_defs.h"

#ifndef FILE_S11_IE_FORMATTER_SEEN
#define FILE_S11_IE_FORMATTER_SEEN

//-----------------
typedef struct ebis_to_be_deleted_s {
#define MSG_DELETE_BEARER_REQUEST_MAX_EBIS_TO_BE_DELETED 11
  uint8_t num_ebis;
  uint8_t eps_bearer_id
      [MSG_DELETE_BEARER_REQUEST_MAX_EBIS_TO_BE_DELETED];  ///< EPS bearer ID
} ebis_to_be_deleted_t;

nw_rc_t gtpv2c_msisdn_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

/* Node Type Information Element
 * 3GPP TS 29.274 #8.34
 * Node type:
 *      * 0 = MME
 *      * 1 = SGSN
 */
nw_rc_t gtpv2c_node_type_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_node_type_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const node_type_t* node_type);

/* PDN Type Information Element
 * 3GPP TS 29.274 #8.34
 * PDN type:
 *      * 1 = IPv4
 *      * 2 = IPv6
 *      * 3 = IPv4v6
 */
nw_rc_t gtpv2c_pdn_type_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_pdn_type_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const pdn_type_t* pdn_type);

/* RAT type Information Element
 * WARNING: the RAT type used in MME and S/P-GW is not the same as the one
 * for S11 interface defined in 3GPP TS 29.274 #8.17.
 */
nw_rc_t gtpv2c_rat_type_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_rat_type_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const rat_type_t* rat_type);

nw_rc_t gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_bearer_context_to_be_created_within_create_bearer_request_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_to_be_created_t* bearer_context);

nw_rc_t gtpv2c_failed_bearer_contexts_within_delete_bearer_request_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_failed_bearer_context_within_delete_bearer_request_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_to_be_removed_t* bearer_context);

int gtpv2c_bearer_context_within_create_bearer_response_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_within_create_bearer_response_t* bearer_context);
nw_rc_t gtpv2c_bearer_context_within_create_bearer_response_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_ebis_within_delete_bearer_request_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const ebis_to_be_deleted_t* ebis_tbd);
nw_rc_t gtpv2c_ebis_to_be_deleted_within_delete_bearer_request_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_bearer_context_within_delete_bearer_response_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_within_delete_bearer_response_t* bearer_context);
nw_rc_t gtpv2c_bearer_context_within_delete_bearer_response_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_bearer_context_ebi_only_ie_set(
    nw_gtpv2c_msg_handle_t* const msg, const ebi_t ebi);

int gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_to_be_modified_t* bearer_context);
nw_rc_t
gtpv2c_bearer_context_to_be_modified_within_modify_bearer_request_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_bearer_context_to_be_removed_within_modify_bearer_request_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_to_be_removed_t* bearer_context);
nw_rc_t gtpv2c_bearer_context_to_be_removed_within_modify_bearer_request_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

nw_rc_t gtpv2c_pti_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const pti_t pti, const uint8_t ieInstance);
nw_rc_t gtpv2c_pti_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

nw_rc_t gtpv2c_ebi_ie_get_list(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

/* Bearer Context Created grouped Information Element */
nw_rc_t gtpv2c_bearer_context_modified_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

nw_rc_t gtpv2c_bearer_context_marked_for_removal_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_bearer_context_marked_for_removal_ie_set(
    nw_gtpv2c_msg_handle_t* msg,
    const bearer_context_marked_for_removal_t* const bearer);

/* Serving Network Information Element
 * 3GPP TS 29.274 #8.18
 */
nw_rc_t gtpv2c_serving_network_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_serving_network_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const ServingNetwork_t* serving_network);

/* Protocol Configuration Options Information Element */
nw_rc_t gtpv2c_pco_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_pco_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const protocol_configuration_options_t* pco);

/* Access Point Name Information Element
 * 3GPP TS 29.274 #8.6
 * NOTE: The APN field is not encoded as a dotted string as commonly used in
 * documentation.
 * The encoding of the APN field follows 3GPP TS 23.003 subclause 9.1
 */
nw_rc_t gtpv2c_apn_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

int gtpv2c_apn_ie_set(nw_gtpv2c_msg_handle_t* msg, const char* apn);

int gtpv2c_apn_plmn_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const char* apn,
    const ServingNetwork_t* serving_network);

int gtpv2c_uli_ie_set(nw_gtpv2c_msg_handle_t* msg, const Uli_t* uli);

nw_rc_t gtpv2c_mei_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

nw_rc_t gtpv2c_uli_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

/* APN restriction Information Element
 * 3GPP TS 29.274 #8.57
 */
nw_rc_t gtpv2c_apn_restriction_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_apn_restriction_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const uint8_t apn_restriction);

/* IP address Information Element
 * 3GPP TS 29.274 #8.9
 */
nw_rc_t gtpv2c_ip_address_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_ip_address_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const gtp_ip_address_t* ip_address);

/* Delay Value Information Element
 * 3GPP TS 29.274 #8.27
 */
nw_rc_t gtpv2c_delay_value_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_delay_value_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const DelayValue_t* delay_value);

/* UE Time Zone Information Element
 * 3GPP TS 29.274 #8.44
 */
nw_rc_t gtpv2c_ue_time_zone_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_ue_time_zone_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const UETimeZone_t* ue_time_zone);

/* Bearer Flags Information Element
 * 3GPP TS 29.274 #8.32
 */
nw_rc_t gtpv2c_bearer_flags_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_bearer_flags_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const bearer_flags_t* bearer_flags);

/* Indication Element
 * 3GPP TS 29.274 #8.12
 */
nw_rc_t gtpv2c_indication_flags_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);
int gtpv2c_indication_flags_ie_set(
    nw_gtpv2c_msg_handle_t* msg, const indication_flags_t* indication_flags);

/* FQ-CSID Information Element
 * 3GPP TS 29.274 #8.62
 */

nw_rc_t gtpv2c_fqcsid_ie_get(
    uint8_t ieType, uint16_t ieLength, uint8_t ieInstance, uint8_t* ieValue,
    void* arg);

#endif /* FILE_S11_IE_FORMATTER_SEEN */
