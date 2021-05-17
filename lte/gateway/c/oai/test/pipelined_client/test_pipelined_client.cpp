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

#include "utility_protobuf.h"

using grpc::Status;
using ::testing::Test;

namespace magma {
namespace lte {

std::string convert_to_ipv4_host_format(IPAddress req_ipv4_addr) {
  char req_ipv4_addr_str[INET_ADDRSTRLEN];

  inet_ntop(
      AF_INET, (req_ipv4_addr.address().c_str()), req_ipv4_addr_str,
      INET_ADDRSTRLEN);

  return req_ipv4_addr_str;
}

std::string convert_to_ipv6_host_format(IPAddress req_ipv6_addr) {
  char req_ipv6_addr_str[INET6_ADDRSTRLEN];

  inet_ntop(
      AF_INET6, (req_ipv6_addr.address().c_str()), req_ipv6_addr_str,
      INET6_ADDRSTRLEN);

  return req_ipv6_addr_str;
}

// utility testcase-1
TEST(test_classifier_rpc, test_utility_update_request_ipv4) {
  std::string enb_ipv4_addr_str = "192.168.60.141";
  struct in_addr enb_ipv4_addr;

  std::string ue_ipv4_addr_str = "192.168.128.11";
  struct in_addr ue_ipv4_addr;

  uint32_t in_teid  = 100;
  uint32_t out_teid = 200;

  uint32_t vlan = 100;

  uint32_t ue_state = UE_SESSION_ACTIVE_STATE;

  UESessionSet request;

  inet_pton(AF_INET, enb_ipv4_addr_str.c_str(), &enb_ipv4_addr);
  inet_pton(AF_INET, ue_ipv4_addr_str.c_str(), &ue_ipv4_addr);

  request = create_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, vlan, ue_state);

  // Validate GNB IP Address
  EXPECT_TRUE(
      convert_to_ipv4_host_format(request.enb_ip_address()) ==
      enb_ipv4_addr_str);

  // Validate UE IP Address
  EXPECT_TRUE(
      convert_to_ipv4_host_format(request.ue_ipv4_address()) ==
      ue_ipv4_addr_str);

  // Validate TEID
  EXPECT_TRUE(in_teid == request.in_teid());
  EXPECT_TRUE(out_teid == request.out_teid());

  // UE Session State
  EXPECT_TRUE(
      request.ue_session_state().ue_config_state() == UESessionState::ACTIVE);
}

// utility testcase-2
TEST(test_classifier_rpc, test_utility_proto_ip_flow_dl) {
  std::string dst_ip_addr_str = "192.168.60.141";
  std::string src_ip_addr_str = "192.168.128.11";
  struct ip_flow_dl flow_dl;
  IPFlowDL req_flow_dl;

  memset(&flow_dl, 0, sizeof(struct ip_flow_dl));
  flow_dl.set_params   = (DST_IPV4 | SRC_IPV4);
  flow_dl.tcp_dst_port = 5002;
  flow_dl.tcp_src_port = 60;
  flow_dl.ip_proto     = 6;
  inet_pton(AF_INET, dst_ip_addr_str.c_str(), &(flow_dl.dst_ip));
  inet_pton(AF_INET, src_ip_addr_str.c_str(), &(flow_dl.src_ip));

  req_flow_dl = to_proto_ip_flow_dl(flow_dl);

  // Validate flow_dl destination IP Address
  EXPECT_TRUE(
      convert_to_ipv4_host_format(req_flow_dl.dest_ip()) == dst_ip_addr_str);

  // Validate flow_dl source IP Address
  EXPECT_TRUE(
      convert_to_ipv4_host_format(req_flow_dl.src_ip()) == src_ip_addr_str);

  EXPECT_TRUE(req_flow_dl.set_params() == (DST_IPV4 | SRC_IPV4));
  EXPECT_TRUE(req_flow_dl.tcp_dst_port() == flow_dl.tcp_dst_port);
  EXPECT_TRUE(req_flow_dl.tcp_src_port() == flow_dl.tcp_src_port);
  EXPECT_TRUE(req_flow_dl.ip_proto() == flow_dl.ip_proto);
}

TEST(test_classifier_rpc, test_util_class_flow_dl) {
  FlowDLOps flow_dl_ops = FlowDLOps();
  flow_dl_ops.set_flow_dl();
  IPFlowDL req_flow_dl = to_proto_ip_flow_dl(flow_dl_ops.get_flow_dl());

  EXPECT_TRUE(flow_dl_ops.validate_flow_dl(req_flow_dl));
}

TEST(test_classifier_rpc, test_util_update_request_ipv4) {
  UpdateRequestV4 update_v4_msg = UpdateRequestV4(UE_SESSION_ACTIVE_STATE);
  update_v4_msg.set_update_request_ipv4();

  UESessionSet request = update_v4_msg.get_update_request_ipv4();

  EXPECT_TRUE(update_v4_msg.validate_req_msg(request));
}

TEST(test_classifier_add_rpc, test_tunnel_add_ipv4_flow_dl) {
  UESessionSet request;
  std::string subs_id      = "00000001";
  std::string imsi_val     = "IMSI" + subs_id;
  std::string apn          = "magmacore.com";
  uint32_t flow_precedence = 10;
  struct in_addr ue_ipv4_addr;
  struct in_addr enb_ipv4_addr;

  // Initialize base request
  UpdateRequestV4 update_v4 = UpdateRequestV4(UE_SESSION_ACTIVE_STATE);
  update_v4.get_enb_v4_addr(&enb_ipv4_addr);
  update_v4.get_ue_v4_addr(&ue_ipv4_addr);

  // Initialize flow dl ops
  FlowDLOps flow_dl_ops = FlowDLOps();
  flow_dl_ops.set_flow_dl();

  // Create the message
  request = create_add_update_request_ipv4_flow_dl(
      ue_ipv4_addr, update_v4.get_vlan(), enb_ipv4_addr,
      update_v4.get_in_teid(), update_v4.get_out_teid(), subs_id,
      flow_dl_ops.get_flow_dl(), flow_precedence, apn, UE_SESSION_ACTIVE_STATE);

  // Check the base request (without flow_dl)
  EXPECT_TRUE(update_v4.validate_req_msg(request));

  // Check the flow_dl
  EXPECT_TRUE(flow_dl_ops.validate_flow_dl(request.ip_flow_dl()));

  // Check the IMSI
  EXPECT_TRUE(request.subscriber_id().id() == imsi_val);

  // Check the flow_precedence
  EXPECT_TRUE(request.precedence() == flow_precedence);

  // Check the APN
  EXPECT_TRUE(request.apn() == apn);
}

TEST(test_classifier_add_rpc, test_tunnel_add_ipv4v6_flow_dl) {
  UESessionSet request;
  std::string ue_ipv6_str = "2001::1";
  struct in6_addr ue_ipv6_addr;
  std::string subs_id      = "00000001";
  std::string imsi_val     = "IMSI" + subs_id;
  std::string apn          = "magmacore.com";
  uint32_t flow_precedence = 10;
  struct in_addr ue_ipv4_addr;
  struct in_addr enb_ipv4_addr;

  // Initialize base request
  UpdateRequestV4 update_v4 = UpdateRequestV4(UE_SESSION_ACTIVE_STATE);
  update_v4.get_enb_v4_addr(&enb_ipv4_addr);
  update_v4.get_ue_v4_addr(&ue_ipv4_addr);

  inet_pton(AF_INET6, ue_ipv6_str.c_str(), &ue_ipv6_addr);

  // Initialize flow dl ops
  FlowDLOps flow_dl_ops = FlowDLOps();
  flow_dl_ops.set_flow_dl();

  // Create the message
  request = create_add_update_request_ipv4v6_flow_dl(
      ue_ipv4_addr, ue_ipv6_addr, update_v4.get_vlan(), enb_ipv4_addr,
      update_v4.get_in_teid(), update_v4.get_out_teid(), subs_id,
      flow_dl_ops.get_flow_dl(), flow_precedence, apn, UE_SESSION_ACTIVE_STATE);

  // Check the base request (without flow_dl)
  EXPECT_TRUE(update_v4.validate_req_msg(request));

  // Check the flow_dl
  EXPECT_TRUE(flow_dl_ops.validate_flow_dl(request.ip_flow_dl()));

  // Check the UE IPv6 address
  EXPECT_TRUE(
      convert_to_ipv6_host_format(request.ue_ipv6_address()) == ue_ipv6_str);

  // Check the IMSI
  EXPECT_TRUE(request.subscriber_id().id() == imsi_val);

  // Check the flow_precedence
  EXPECT_TRUE(request.precedence() == flow_precedence);

  // Check the APN
  EXPECT_TRUE(request.apn() == apn);
}

TEST(test_classifier_del_rpc, test_tunnel_del_ipv4_flow_dl) {
  UESessionSet request;
  struct in_addr ue_ipv4_addr;
  struct in_addr enb_ipv4_addr;

  // Initialize base request
  UpdateRequestV4 update_v4 = UpdateRequestV4(UE_SESSION_UNREGISTERED_STATE);
  update_v4.get_enb_v4_addr(&enb_ipv4_addr);
  update_v4.get_ue_v4_addr(&ue_ipv4_addr);

  // Initialize flow dl ops
  FlowDLOps flow_dl_ops = FlowDLOps();
  flow_dl_ops.set_flow_dl();

  // Create the message
  request = create_del_update_request_ipv4_flow_dl(
      enb_ipv4_addr, ue_ipv4_addr, update_v4.get_in_teid(),
      update_v4.get_out_teid(), flow_dl_ops.get_flow_dl(),
      UE_SESSION_UNREGISTERED_STATE);

  // Check the base request (without flow_dl)
  EXPECT_TRUE(update_v4.validate_req_msg(request));

  // Check the flow_dl
  EXPECT_TRUE(flow_dl_ops.validate_flow_dl(request.ip_flow_dl()));
}

TEST(test_classifier_rpc, test_tunnel_discard_data_ipv4_flow_dl) {
  UESessionSet discard_request = UESessionSet();
  struct in_addr ue_ipv4_addr;
  std::string ue_v4 = "192.168.1.20";
  uint32_t in_teid  = 455;

  // Initialize base request
  inet_pton(AF_INET, ue_v4.c_str(), &ue_ipv4_addr);

  // Create the message
  discard_request = create_discard_data_update_request_ipv4(
      ue_ipv4_addr, in_teid, UE_SESSION_SUSPENDED_DATA_STATE);

  EXPECT_TRUE(
      convert_to_ipv4_host_format(discard_request.ue_ipv4_address()) == ue_v4);

  EXPECT_TRUE(in_teid == discard_request.in_teid());

  EXPECT_TRUE(
      discard_request.ue_session_state().ue_config_state() ==
      UESessionState::SUSPENDED_DATA);
}

TEST(test_classifier_rpc, test_tunnel_forward_data_ipv4v6_flow_dl) {
  UESessionSet forward_request;

  std::string ue_v4 = "192.168.1.20";
  struct in_addr ue_ipv4_addr;
  std::string ue_ipv6 = "2001::1";
  struct in6_addr ue_ipv6_addr;
  uint32_t in_teid         = 455;
  uint32_t flow_precedence = 10;

  // Initialize flow dl ops
  FlowDLOps flow_dl_ops = FlowDLOps();
  flow_dl_ops.set_flow_dl();

  // Initialize base request
  inet_pton(AF_INET, ue_v4.c_str(), &ue_ipv4_addr);
  inet_pton(AF_INET6, ue_ipv6.c_str(), &ue_ipv6_addr);

  // Create the message
  forward_request = create_forwarding_data_update_request_ipv4v6_flow_dl(
      ue_ipv4_addr, ue_ipv6_addr, in_teid, flow_dl_ops.get_flow_dl(),
      flow_precedence, UE_SESSION_RESUME_DATA_STATE);

  // Check the UE Addresses
  EXPECT_TRUE(
      convert_to_ipv4_host_format(forward_request.ue_ipv4_address()) == ue_v4);

  // Check the UE v6 Addresses
  EXPECT_TRUE(
      convert_to_ipv6_host_format(forward_request.ue_ipv6_address()) ==
      ue_ipv6);
  // Check the in-teid
  EXPECT_TRUE(in_teid == forward_request.in_teid());

  // Check the precedence
  EXPECT_TRUE(forward_request.precedence() == flow_precedence);

  // Check the UE State
  EXPECT_TRUE(
      forward_request.ue_session_state().ue_config_state() ==
      UESessionState::RESUME_DATA);
}

TEST(test_classifier_rpc, test_paging_uev4) {
  UESessionSet paging_request;

  std::string ue_v4 = "192.168.1.20";
  struct in_addr ue_ipv4_addr;

  // Initialize base request
  inet_pton(AF_INET, ue_v4.c_str(), &ue_ipv4_addr);

  // Create the message
  paging_request = create_paging_update_request_ipv4(
      ue_ipv4_addr, UE_SESSION_INSTALL_IDLE_STATE);

  // Check the UE Addresses
  EXPECT_TRUE(
      convert_to_ipv4_host_format(paging_request.ue_ipv4_address()) == ue_v4);

  // Check the UE State
  EXPECT_TRUE(
      paging_request.ue_session_state().ue_config_state() ==
      UESessionState::INSTALL_IDLE);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace lte
}  // namespace magma
