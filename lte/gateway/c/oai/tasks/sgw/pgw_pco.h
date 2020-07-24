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

/*! \file pgw_lite_paa.h
 * \brief
 * \author Lionel Gauthier
 * \company Eurecom
 * \email: lionel.gauthier@eurecom.fr
 */
#ifndef FILE_PGW_PCO_SEEN
#define FILE_PGW_PCO_SEEN

#include <stdint.h>

#include "3gpp_24.008.h"

/**
 * protocol_configuration_options_ids_t
 *
 * Container for caching which protocol/container identifiers have been set in
 * the message sent by the UE.
 *
 * ID specifications based on 3GPP #24.008.
 */
typedef struct protocol_configuration_options_ids_s {
  // Protocol identifiers (from configuration protocol options list)
  uint8_t pi_ipcp : 1;

  // Container identifiers (from additional parameters list)
  uint8_t ci_dns_server_ipv4_address_request : 1;
  uint8_t ci_ip_address_allocation_via_nas_signalling : 1;
  uint8_t ci_ipv4_address_allocation_via_dhcpv4 : 1;
  uint8_t ci_ipv4_link_mtu_request : 1;
} protocol_configuration_options_ids_t;

int pgw_pco_push_protocol_or_container_id(
    protocol_configuration_options_t* const pco,
    pco_protocol_or_container_id_t* const poc_id);

int pgw_process_pco_request_ipcp(
    protocol_configuration_options_t* const pco_resp,
    const pco_protocol_or_container_id_t* const poc_id);

int pgw_process_pco_dns_server_request(
    protocol_configuration_options_t* const pco_resp,
    const pco_protocol_or_container_id_t* const poc_id);

int pgw_process_pco_link_mtu_request(
    protocol_configuration_options_t* const pco_resp,
    const pco_protocol_or_container_id_t* const poc_id);

int pgw_process_pco_request(
    const protocol_configuration_options_t* const pco_req,
    protocol_configuration_options_t* pco_resp,
    protocol_configuration_options_ids_t* const pco_ids);

#endif
