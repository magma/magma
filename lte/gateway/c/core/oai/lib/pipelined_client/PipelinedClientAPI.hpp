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

#pragma once

#include "stdint.h"

#include "lte/gateway/c/core/oai/tasks/gtpv1-u/gtpv1u.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#ifdef __cplusplus
}
#endif

typedef enum {
  UE_SESSION_ACTIVE_STATE,
  UE_SESSION_UNREGISTERED_STATE,
  UE_SESSION_INSTALL_IDLE_STATE,
  UE_SESSION_UNINSTALL_IDLE_STATE,
  UE_SESSION_SUSPENDED_DATA_STATE,
  UE_SESSION_RESUME_DATA_STATE,
} ue_session_states;

int upf_classifier_add_tunnel(struct in_addr ue, struct in6_addr* ue_ipv6,
                              int vlan, struct in_addr enb, uint32_t i_tei,
                              uint32_t o_tei, const char* imsi,
                              struct ip_flow_dl* flow_dl,
                              uint32_t flow_precedence_dl, const char* apn);

int upf_classifier_del_tunnel(struct in_addr enb, struct in_addr ue,
                              struct in6_addr* ue_ipv6, uint32_t i_tei,
                              uint32_t o_tei, struct ip_flow_dl* flow_dl);

int upf_classifier_discard_data_on_tunnel(struct in_addr ue,
                                          struct in6_addr* ue_ipv6,
                                          uint32_t i_tei,
                                          struct ip_flow_dl* flow_dl);

int upf_classifier_forward_data_on_tunnel(struct in_addr ue,
                                          struct in6_addr* ue_ipv6,
                                          uint32_t i_tei,
                                          struct ip_flow_dl* flow_dl,
                                          uint32_t flow_precedence_dl);

int upf_classifier_add_paging_rule(struct in_addr ue);

int upf_classifier_delete_paging_rule(struct in_addr ue);
