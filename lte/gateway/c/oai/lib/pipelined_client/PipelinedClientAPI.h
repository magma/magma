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
#ifndef PIPELINE_RPC_CLIENT_H
#define PIPELINE_RPC_CLIENT_H

#include "stdint.h"

#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
typedef enum {
  UE_SESSION_ACTIVE_STATE,
  UE_SESSION_UNREGISTERED_STATE,
  UE_SESSION_INSTALL_IDLE_STATE,
  UE_SESSION_UNINSTALL_IDLE_STATE,
  UE_SESSION_SUSPENDED_DATA_STATE,
  UE_SESSION_RESUME_DATA_STATE,
} ue_session_states;

#include "3gpp_23.003.h"
#include "gtpv1u.h"

int upf_classifier_add_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, int vlan, struct in_addr enb,
    uint32_t i_tei, uint32_t o_tei, const char* imsi,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl, const char* apn);

int upf_classifier_del_tunnel(
    struct in_addr enb, struct in_addr ue, struct in6_addr* ue_ipv6,
    uint32_t i_tei, uint32_t o_tei, struct ip_flow_dl* flow_dl);

int upf_classifier_discard_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl);

int upf_classifier_forward_data_on_tunnel(
    struct in_addr ue, struct in6_addr* ue_ipv6, uint32_t i_tei,
    struct ip_flow_dl* flow_dl, uint32_t flow_precedence_dl);

int upf_classifier_add_paging_rule(struct in_addr ue);

int upf_classifier_delete_paging_rule(struct in_addr ue);

#ifdef __cplusplus
}
#endif

#endif  // PIPELINE_RPC_CLIENT_H
