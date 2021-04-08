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

#include "proto_converters.h"
#include "PipelinedClientAPI.h"

namespace magma {
namespace lte {


// todo call this something better
UESessionSet make_update_request_ipv4(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const struct ip_flow_dl& flow_dl,
    uint32_t ue_state) {
  UESessionSet request;
  // Set the enb IPv4 address
  set_gnb_ipv4_addr(enb_ipv4_addr, request);

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);
  request.set_out_teid(out_teid);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  // Set the ue_state
  config_ue_session_state(ue_state, request);
  return request;
}

IPFlowDL to_proto_ip_flow_dl(struct ip_flow_dl flow_dl) {
  IPFlowDL ue_flow_dl = IPFlowDL();

  ue_flow_dl.set_set_params(flow_dl.set_params);
  ue_flow_dl.set_tcp_dst_port(flow_dl.tcp_dst_port);
  ue_flow_dl.set_tcp_src_port(flow_dl.tcp_src_port);
  ue_flow_dl.set_udp_dst_port(flow_dl.udp_dst_port);
  ue_flow_dl.set_udp_src_port(flow_dl.udp_src_port);
  ue_flow_dl.set_ip_proto(flow_dl.ip_proto);

  if ((flow_dl.set_params & DST_IPV4) || (flow_dl.set_params & SRC_IPV4)) {
    if (flow_dl.set_params & DST_IPV4) {
      IPAddress* dest_ip = ue_flow_dl.mutable_dest_ip();
      dest_ip->set_version(IPAddress::IPV4);
      dest_ip->set_address(&flow_dl.dst_ip, sizeof(struct in_addr));
    }

    if (flow_dl.set_params & SRC_IPV4) {
      IPAddress* src_ip = ue_flow_dl.mutable_src_ip();
      src_ip->set_version(IPAddress::IPV4);
      src_ip->set_address(&flow_dl.src_ip, sizeof(struct in_addr));
    }

  } else {
    if (flow_dl.set_params & DST_IPV6) {
      IPAddress* dest_ip = ue_flow_dl.mutable_dest_ip();
      dest_ip->set_version(IPAddress::IPV6);
      dest_ip->set_address(&flow_dl.dst_ip6, sizeof(struct in6_addr));
    }

    if (flow_dl.set_params & SRC_IPV6) {
      IPAddress* src_ip = ue_flow_dl.mutable_src_ip();
      src_ip->set_version(IPAddress::IPV6);
      src_ip->set_address(&flow_dl.src_ip6, sizeof(struct in6_addr));
    }
  }

  return ue_flow_dl;
}

void set_ue_ipv4_addr(
    const struct in_addr& ue_ipv4_addr, UESessionSet& request) {
  IPAddress* encode_ue_ipv4_addr = request.mutable_ue_ipv4_address();
  encode_ue_ipv4_addr->set_version(IPAddress::IPV4);
  encode_ue_ipv4_addr->set_address(&ue_ipv4_addr, sizeof(struct in_addr));
}

void set_ue_ipv6_addr(
    const struct in6_addr& ue_ipv6_addr, UESessionSet& request) {
  IPAddress* encode_ue_ipv6_addr = request.mutable_ue_ipv6_address();
  encode_ue_ipv6_addr->set_version(IPAddress::IPV6);
  encode_ue_ipv6_addr->set_address(&ue_ipv6_addr, sizeof(struct in6_addr));
}

// Set the GNB IPv4 address
void set_gnb_ipv4_addr(
    const struct in_addr& gnb_ipv4_addr, UESessionSet& request) {
  IPAddress* encode_gnb_ipv4_addr = request.mutable_enb_ip_address();
  encode_gnb_ipv4_addr->set_version(IPAddress::IPV4);
  encode_gnb_ipv4_addr->set_address(&gnb_ipv4_addr, sizeof(struct in_addr));
}

// Set the Session Config State
void config_ue_session_state(uint32_t& ue_state, UESessionSet& request) {
  UESessionState* ue_session_state = request.mutable_ue_session_state();

  if (UE_SESSION_ACTIVE_STATE == ue_state) {
    ue_session_state->set_ue_config_state(UESessionState::ACTIVE);
  } else if (UE_SESSION_UNREGISTERED_STATE == ue_state) {
    ue_session_state->set_ue_config_state(UESessionState::UNREGISTERED);
  } else if (UE_SESSION_INSTALL_IDLE_STATE == ue_state) {
    ue_session_state->set_ue_config_state(UESessionState::INSTALL_IDLE);
  } else if (UE_SESSION_UNINSTALL_IDLE_STATE == ue_state) {
    ue_session_state->set_ue_config_state(UESessionState::UNINSTALL_IDLE);
  } else if (UE_SESSION_SUSPENDED_DATA_STATE == ue_state) {
    ue_session_state->set_ue_config_state(UESessionState::SUSPENDED_DATA);
  } else if (UE_SESSION_RESUME_DATA_STATE == ue_state) {
    ue_session_state->set_ue_config_state(UESessionState::RESUME_DATA);
  }
}

}  // namespace lte
}  // namespace magma