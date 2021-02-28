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

#include <assert.h>
#include <errno.h>
#include <stdint.h>
#include <netinet/in.h>
#include <stdlib.h>

#include "assertions.h"
#include "bstrlib.h"
#include "log.h"
#include "3gpp_23.003.h"
#include "spgw_config.h"
#include "gtp_tunnel_upf.h"
#include "PipelinedClientAPI.h"

int upf_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, Imsi_t imsi, struct ip_flow_dl* flow_dl,
    uint32_t flow_precedence_dl, char* apn) {
  upf_classifier_add_tunnel(
      ue, ue_ipv6, vlan, enb, i_tei, o_tei, (const char*) imsi.digit, flow_dl,
      flow_precedence_dl, apn);
  return 0;
}

int upf_del_tunnel(
    struct in_addr enb, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t o_tei, struct ip_flow_dl* flow_dl) {
  upf_classifier_del_tunnel(enb, ue, ue_ipv6, i_tei, o_tei, flow_dl);

  return 0;
}

int upf_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl) {
  upf_classifier_discard_data_on_tunnel(ue, ue_ipv6, i_tei, flow_dl);

  return 0;
}

int upf_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl) {
  upf_classifier_forward_data_on_tunnel(
      ue, ue_ipv6, i_tei, flow_dl, flow_precedence_dl);

  return 0;
}

int upf_add_paging_rule(struct in_addr ue) {
  upf_classifier_add_paging_rule(ue);

  return 0;
}

int upf_delete_paging_rule(struct in_addr ue) {
  upf_classifier_delete_paging_rule(ue);

  return 0;
}
