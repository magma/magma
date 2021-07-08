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

/* File : gtp_tunnel_upf.h
 */

#pragma once

#include "gtpv1u.h"
#include "gtp_tunnel_openflow.h"

const struct gtp_tunnel_ops* upf_gtp_tunnel_ops_init_openflow(void);

int upf_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl, char* apn);

int upf_del_tunnel(
    struct in_addr enb, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t o_tei, struct ip_flow_dl* flow_dl);

int upf_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl);

int upf_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl);

int upf_add_paging_rule(struct in_addr ue);

int upf_delete_paging_rule(struct in_addr ue);
