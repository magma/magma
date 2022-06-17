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
 *------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */
/*! \file gtpv1u_task.c
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdint.h>
#include <string.h>

//------------------------------------------------------------------------------
void add_route_for_ue_block(struct in_addr ue_net, uint32_t mask);

static void del_route_for_ue_block(struct in_addr ue_net, uint32_t mask);

static bool ue_ip_is_in_subnet(struct in_addr _net, int mask,
                               struct in_addr _addr);

//------------------------------------------------------------------------------
int gtpv1u_init_1();
//=== HARDCODE
int gtpv1u_add_tunnel(struct in_addr ue, struct in6_addr* ue_ipv6, int vlan,
                      struct in_addr enb, struct in6_addr* enb_ipv6,
                      uint32_t i_tei, uint32_t o_tei, char* imsi,
                      struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl,
                      char* apn);

//------------------------------------------------------------------------------
void gtpv1u_exit(void);
