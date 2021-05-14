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

// Utility : IPv4 tunnel create generic request
UESessionSet create_update_request_ipv4(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, int vlan, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the enb IPv4 address
  set_gnb_ipv4_addr(enb_ipv4_addr, request);

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);
  request.set_out_teid(out_teid);

  if (vlan) {
    request.set_vlan(vlan);
  }

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}

//---------- TUNNEL IPV4 ONLY ADD FUNCTION ----------

// IPv4 tunnel add request
UESessionSet create_add_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, int vlan, struct in_addr& enb_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const std::string& imsi,
    uint32_t flow_precedence, const std::string& apn, uint32_t ue_state) {
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, vlan, ue_state);

  // Set the Subscriber ID
  SubscriberID* sid = request.mutable_subscriber_id();
  sid->set_id("IMSI" + std::string(imsi));
  sid->set_type(SubscriberID::IMSI);

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the APN
  request.set_apn(apn);

  return request;
}

// IPv4 tunnel add request with flow_dl
UESessionSet create_add_update_request_ipv4_flow_dl(
    const struct in_addr& ue_ipv4_addr, int vlan, struct in_addr& enb_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const std::string& imsi,
    const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
    const std::string& apn, uint32_t ue_state) {
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, vlan, ue_state);

  // Set the Subscriber ID
  SubscriberID* sid = request.mutable_subscriber_id();
  sid->set_id("IMSI" + std::string(imsi));
  sid->set_type(SubscriberID::IMSI);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the APN
  request.set_apn(apn);

  return request;
}

//---------- TUNNEL IPV4 ONLY DEL FUNCTION ----------

// IPv4 tunnel del request
UESessionSet create_del_update_request_ipv4(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, uint32_t ue_state) {
  // For deletion vlan=0
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, 0, ue_state);

  return request;
}

// IPv4 tunnel del request with flow_dl
UESessionSet create_del_update_request_ipv4_flow_dl(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const struct ip_flow_dl& flow_dl,
    uint32_t ue_state) {
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, 0, ue_state);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  return request;
}

//---------- TUNNEL IPv4-v6 ADD FUNCTION ----------

// IPv4-v6 tunnel add request
UESessionSet create_add_update_request_ipv4v6(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr, int vlan,
    struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
    const std::string& imsi, uint32_t flow_precedence, const std::string& apn,
    uint32_t ue_state) {
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, vlan, ue_state);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  // Set the Subscriber ID
  SubscriberID* sid = request.mutable_subscriber_id();
  sid->set_id("IMSI" + std::string(imsi));
  sid->set_type(SubscriberID::IMSI);

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the APN
  request.set_apn(apn);

  return request;
}

// IPv4-v6 tunnel add request with flow_dl
UESessionSet create_add_update_request_ipv4v6_flow_dl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr, int vlan,
    struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
    const std::string& imsi, const struct ip_flow_dl& flow_dl,
    uint32_t flow_precedence, const std::string& apn, uint32_t ue_state) {
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, vlan, ue_state);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  // Set the Subscriber ID
  SubscriberID* sid = request.mutable_subscriber_id();
  sid->set_id("IMSI" + std::string(imsi));
  sid->set_type(SubscriberID::IMSI);

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the APN
  request.set_apn(apn);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  return request;
}

//---------- TUNNEL IPv4-v6 DEL FUNCTION ----------

// IPv4-v6 tunnel del request
UESessionSet create_del_update_request_ipv4v6(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
    uint32_t ue_state) {
  // For deletion vlan=0
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, 0, ue_state);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  return request;
}

// IPv4-v6 tunnel del request with flow dl
UESessionSet create_del_update_request_ipv4v6_flow_dl(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
    const struct ip_flow_dl& flow_dl, uint32_t ue_state) {
  // For deletion vlan=0
  UESessionSet request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, 0, ue_state);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  return request;
}

//---------- TUNNEL DISCARD FUNCTION ----------

// IPv4 request for discarding data on tunnel
UESessionSet create_discard_data_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}

// IPv4 request for discarding data on tunnel with flow_dl
UESessionSet create_discard_data_update_request_ipv4_flow_dl(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    const struct ip_flow_dl& flow_dl, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  return request;
}

// IPv4v6 request for discarding data on tunnel
UESessionSet create_discard_data_update_request_ipv4v6(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  return request;
}

// IPv4v6 request for discarding data on tunnel with flow_dl
UESessionSet create_discard_data_update_request_ipv4v6_flow_dl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, const struct ip_flow_dl& flow_dl, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  return request;
}

//---------- TUNNEL FORWARD FUNCTION ----------

// IPv4 request for forwarding data on tunnel
UESessionSet create_forwarding_data_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    uint32_t flow_precedence, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}

// IPv4 request for forwarding data on tunnel with flow_dl
UESessionSet create_forwarding_data_update_request_ipv4_flow_dl(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
    uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}

// IPv4v6 request for forwarding data on tunnel
UESessionSet create_forwarding_data_update_request_ipv4v6(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, uint32_t flow_precedence, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}  // namespace lte

// IPv4v6 request for forwarding data on tunnel with flow_dl
UESessionSet create_forwarding_data_update_request_ipv4v6_flow_dl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, const struct ip_flow_dl& flow_dl,
    uint32_t flow_precedence, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the UE IPv6 address
  set_ue_ipv6_addr(ue_ipv6_addr, request);

  // Set the incoming and outgoing teid
  request.set_in_teid(in_teid);

  // Flow dl
  request.mutable_ip_flow_dl()->CopyFrom(to_proto_ip_flow_dl(flow_dl));

  // Set the precedence
  request.set_precedence(flow_precedence);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}

//---------- UE PAGING FUNCTION ----------

// Create paging request
UESessionSet create_paging_update_request_ipv4(
    const struct in_addr& ue_ipv4_addr, uint32_t ue_state) {
  UESessionSet request = UESessionSet();

  // Set the UE IPv4 address
  set_ue_ipv4_addr(ue_ipv4_addr, request);

  // Set the ue_state
  config_ue_session_state(ue_state, request);

  return request;
}

}  // namespace lte
}  // namespace magma
