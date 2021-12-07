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

/* File : gtp_tunnel_upf.c
 * For sending the flow reuqests to pipeplined
 */

#pragma once

#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u.h"

const struct gtp_tunnel_ops* mock_gtp_tunnel_ops_init(void);

int mock_tunnel_init(
    struct in_addr* ue_net, uint32_t mask, int mtu, int* fd0, int* fd1u,
    bool persist_state);

int mock_tunnel_uninit(void);

int mock_tunnel_reset(void);

int mock_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in6_addr* enb_ipv6, uint32_t i_tei, uint32_t o_tei, Imsi_t imsi,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl, char* apn);

int mock_del_tunnel(
    struct in_addr enb, struct in6_addr* enb_ipv6, struct in_addr ue,
    struct in6_addr* ue_ipv6, uint32_t i_tei, uint32_t o_tei,
    struct ip_flow_dl* flow_dl);

int mock_add_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in6_addr* enb_ipv6, struct in_addr pgw, struct in6_addr* pgw_ipv6,
    uint32_t i_tei, uint32_t o_tei, uint32_t pgw_in_tei, uint32_t pgw_o_tei,
    Imsi_t imsi);

int mock_del_s8_tunnel(
    struct in_addr enb, struct in6_addr* enb_ipv6, struct in_addr pgw,
    struct in6_addr* pgw_ipv6, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t pgw_in_tei);

int mock_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl);

int mock_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl);

int mock_add_paging_rule(Imsi_t imsi, struct in_addr ue);

int mock_delete_paging_rule(struct in_addr ue);

int mock_send_end_marker(struct in_addr enb, uint32_t tei);

const char* mock_get_dev_name(void);
