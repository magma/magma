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

#include <stdbool.h>

#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u.hpp"

#define CONTROLLER_ADDR "127.0.0.1"
#define CONTROLLER_PORT 6654
#define NUM_WORKERS 2

int start_of_controller(bool persist_state);

int stop_of_controller(void);

int openflow_controller_add_gtp_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in6_addr* enb_ipv6, uint32_t i_tei, uint32_t o_tei, const char* imsi,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl,
    uint32_t gtp_portno);

int openflow_controller_del_gtp_tunnel(struct in_addr ue,
                                       struct in6_addr* ue_ipv6, uint32_t i_tei,
                                       struct ip_flow_dl* flow_dl,
                                       uint32_t gtp_portno);

int openflow_controller_discard_data_on_tunnel(struct in_addr ue,
                                               struct in6_addr* ue_ipv6,
                                               uint32_t i_tei,
                                               struct ip_flow_dl* flow_dl);

int openflow_controller_forward_data_on_tunnel(struct in_addr ue,
                                               struct in6_addr* ue_ipv6,
                                               uint32_t i_tei,
                                               struct ip_flow_dl* flow_dl,
                                               uint32_t flow_precedence_dl);

int openflow_controller_add_paging_rule(const char* imsi, struct in_addr ue_ip,
                                        struct in6_addr* ue_ipv6);

int openflow_controller_delete_paging_rule(struct in_addr ue_ip,
                                           struct in6_addr* ue_ipv6);

int openflow_controller_add_gtp_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in6_addr* enb_ipv6, struct in_addr pgw, struct in6_addr* pgw_ipv6,
    uint32_t i_tei, uint32_t o_tei, uint32_t pgw_in_tei, uint32_t pgw_o_tei,
    const char* imsi, uint32_t enb_gtp_port, uint32_t pgw_gtp_port);

int openflow_controller_del_gtp_s8_tunnel(struct in_addr ue,
                                          struct in6_addr* ue_ipv6,
                                          uint32_t i_tei, uint32_t pgw_o_tei,
                                          uint32_t enb_gtp_port,
                                          uint32_t pgw_gtp_port);
