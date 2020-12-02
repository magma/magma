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

/*! \file pgw_ue_ip_address_alloc.h
 * \brief
 * \author
 * \company
 * \email:
 */

#ifndef PGW_UE_IP_ADDRESS_ALLOC_SEEN
#define PGW_UE_IP_ADDRESS_ALLOC_SEEN

#include <arpa/inet.h>
#include <stdint.h>

#include "spgw_state.h"
#include "ip_forward_messages_types.h"

int release_ue_ipv4_address(
    const char* imsi, const char* apn, struct in_addr* addr);

int get_ip_block(struct in_addr* netaddr, uint32_t* netmask);

int release_ue_ipv6_address(
    const char* imsi, const char* apn, struct in6_addr* addr);

#endif /*PGW_UE_IP_ADDRESS_ALLOC_SEEN */
