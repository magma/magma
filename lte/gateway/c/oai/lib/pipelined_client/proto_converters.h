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

#include "lte/protos/pipelined.grpc.pb.h"
#include "lte/protos/pipelined.pb.h"
#include "gtpv1u.h"

namespace magma {
namespace lte {

/**
 * @brief
 *
 * @param enb_ipv4_addr
 * @param ue_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param flow_dl
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet make_update_request_ipv4(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const struct ip_flow_dl& flow_dl,
    uint32_t ue_state);

/**
 * @brief
 *
 * @param flow_dl
 * @return IPFlowDL
 */
IPFlowDL to_proto_ip_flow_dl(struct ip_flow_dl flow_dl);

/**
 * @brief Set the ue ipv4 addr object
 *
 * @param ue_ipv4_addr
 * @param request
 */
void set_ue_ipv4_addr(
    const struct in_addr& ue_ipv4_addr, UESessionSet& request);

/**
 * @brief Set the ue ipv6 addr object
 *
 * @param ue_ipv6_addr
 * @param request
 */
void set_ue_ipv6_addr(
    const struct in6_addr& ue_ipv6_addr, UESessionSet& request);

/**
 * @brief Set the gnb ipv4 addr object
 *
 * @param gnb_ipv4_addr
 * @param request
 */
void set_gnb_ipv4_addr(
    const struct in_addr& gnb_ipv4_addr, UESessionSet& request);

/**
 * @brief
 *
 * @param ue_state
 * @param request
 */
void config_ue_session_state(uint32_t& ue_state, UESessionSet& request);

}  // namespace lte
}  // namespace magma
