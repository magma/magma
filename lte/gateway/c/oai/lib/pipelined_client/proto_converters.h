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
 * @param vlan
 * @param ue_state
 * @return UESessionSet
 */

UESessionSet create_update_request_ipv4(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, int vlan, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param vlan
 * @param enb_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param imsi
 * @param flow_precedence
 * @param apn
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_add_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, int vlan, struct in_addr& enb_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const std::string& imsi,
    uint32_t flow_precedence, const std::string& apn, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param vlan
 * @param enb_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param imsi
 * @param flow_dl
 * @param flow_precedence
 * @param apn
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_add_update_request_ipv4_flow_dl(
    const struct in_addr& ue_ipv4_addr, int vlan, struct in_addr& enb_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const std::string& imsi,
    const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
    const std::string& apn, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param enb_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_del_update_request_ipv4(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param enb_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param ue_state
 * @param flow_dl
 * @return UESessionSet
 */
UESessionSet create_del_update_request_ipv4_flow_dl(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const struct ip_flow_dl& flow_dl,
    uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param vlan
 * @param ue_ipv6_addr
 * @param enb_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param imsi
 * @param flow_precedence
 * @param apn
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_add_update_request_ipv4v6(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr, int vlan,
    struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
    const std::string& imsi, uint32_t flow_precedence, const std::string& apn,
    uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param vlan
 * @param ue_ipv6_addr
 * @param enb_ipv4_addr
 * @param in_teid
 * @param out_teid
 * @param imsi
 * @param flow_dl
 * @param flow_precedence
 * @param apn
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_add_update_request_ipv4v6_flow_dl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr, int vlan,
    struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
    const std::string& imsi, const struct ip_flow_dl& flow_dl,
    uint32_t flow_precedence, const std::string& apn, uint32_t ue_state);

/**
 * @brief
 *
 * @param enb_ipv4_addr
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param out_teid
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_del_update_request_ipv4v6(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
    uint32_t ue_state);

/**
 * @brief
 *
 * @param enb_ipv4_addr
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param out_teid
 * @param ue_state
 * @param flow_dl
 * @return UESessionSet
 */
UESessionSet create_del_update_request_ipv4v6_flow_dl(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
    const struct ip_flow_dl& flow_dl, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param in_teid
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_discard_data_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param in_teid
 * @param flow_dl
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_discard_data_update_request_ipv4_flow_dl(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    const struct ip_flow_dl& flow_dl, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_discard_data_update_request_ipv4v6(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param ue_state
 * @param flow_dl
 * @return UESessionSet
 */
UESessionSet create_discard_data_update_request_ipv4v6_flow_dl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, const struct ip_flow_dl& flow_dl, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_forwarding_data_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    uint32_t flow_precedence, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param ue_state
 * @param flow_dl
 * @return UESessionSet
 */
UESessionSet create_forwarding_data_update_request_ipv4_flow_dl(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
    uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_forwarding_data_update_request_ipv4v6(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, uint32_t flow_precedence, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_ipv6_addr
 * @param in_teid
 * @param ue_state
 * @param flow_dl
 * @return UESessionSet
 */
UESessionSet create_forwarding_data_update_request_ipv4v6_flow_dl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, const struct ip_flow_dl& flow_dl,
    uint32_t flow_precedence, uint32_t ue_state);

/**
 * @brief
 *
 * @param ue_ipv4_addr
 * @param ue_state
 * @return UESessionSet
 */
UESessionSet create_paging_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, uint32_t ue_state);

/**
 * @brief
 *
 * @param flow_dl
 * @return IPFlowDL
 */
IPFlowDL to_proto_ip_flow_dl(struct ip_flow_dl flow_dl);

}  // namespace lte
}  // namespace magma
