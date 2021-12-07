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

#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtp_tunnel_mock.h"
#include "lte/gateway/c/core/oai/include/spgw_config.h"

static const struct gtp_tunnel_ops mock_tunnel_ops = {
    .init                   = mock_tunnel_init,
    .uninit                 = mock_tunnel_uninit,
    .reset                  = mock_tunnel_reset,
    .add_tunnel             = mock_add_tunnel,
    .del_tunnel             = mock_del_tunnel,
    .discard_data_on_tunnel = mock_discard_data_on_tunnel,
    .forward_data_on_tunnel = mock_forward_data_on_tunnel,
    .add_paging_rule        = mock_add_paging_rule,
    .delete_paging_rule     = mock_delete_paging_rule,
    .send_end_marker        = mock_send_end_marker,
    .get_dev_name           = mock_get_dev_name,
};

const struct gtp_tunnel_ops* mock_gtp_tunnel_ops_init(void) {
  return &mock_tunnel_ops;
}

int mock_tunnel_init(
    struct in_addr* ue_net, uint32_t mask, int mtu, int* fd0, int* fd1u,
    bool persist_state) {
  return 0;
}

int mock_tunnel_uninit(void) {
  return 0;
}

int mock_tunnel_reset(void) {
  return 0;
}

int mock_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in6_addr* enb_ipv6, uint32_t i_tei, uint32_t o_tei, Imsi_t imsi,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl, char* apn) {
  return 0;
}

int mock_del_tunnel(
    struct in_addr enb, struct in6_addr* enb_ipv6, struct in_addr ue,
    struct in6_addr* ue_ipv6, uint32_t i_tei, uint32_t o_tei,
    struct ip_flow_dl* flow_dl) {
  return 0;
}

int mock_add_s8_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    struct in6_addr* enb_ipv6, struct in_addr pgw, struct in6_addr* pgw_ipv6,
    uint32_t i_tei, uint32_t o_tei, uint32_t pgw_in_tei, uint32_t pgw_o_tei,
    Imsi_t imsi) {
  return 0;
}

int mock_del_s8_tunnel(
    struct in_addr enb, struct in6_addr* enb_ipv6, struct in_addr pgw,
    struct in6_addr* pgw_ipv6, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t pgw_in_tei) {
  return 0;
}

int mock_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl) {
  return 0;
}

int mock_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl) {
  return 0;
}

int mock_add_paging_rule(Imsi_t imsi, struct in_addr ue) {
  return 0;
}

int mock_delete_paging_rule(struct in_addr ue) {
  return 0;
}

int mock_send_end_marker(struct in_addr enb, uint32_t tei) {
  return 0;
}

const char* mock_get_dev_name(void) {
  return bdata(spgw_config.sgw_config.ovs_config.bridge_name);
}
