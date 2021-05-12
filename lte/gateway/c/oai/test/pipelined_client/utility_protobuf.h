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

#include <thread>
#include <gtest/gtest.h>
#include <arpa/inet.h>
#include <stdio.h>

#include "PipelinedServiceClient.h"
#include "lte/protos/pipelined.grpc.pb.h"
#include "lte/protos/pipelined.pb.h"
#include "lte/protos/mobilityd.pb.h"
#include "proto_converters.h"

namespace magma {
namespace lte {

class UpdateRequestV4 {
 private:
  std::string enb_v4;
  std::string ue_v4;
  uint32_t in_teid;
  uint32_t out_teid;
  uint32_t vlan;
  uint32_t ue_state;

 public:
  UESessionSet update_request;

  UpdateRequestV4(
      uint32_t set_ue_state, const std::string enb_str = "192.168.60.141",
      const std::string ue_str = "192.168.128.11", uint32_t incoming_teid = 100,
      uint32_t outgoing_teid = 200, uint32_t out_vlan = 0)
      : enb_v4(enb_str),
        ue_v4(ue_str),
        in_teid(incoming_teid),
        out_teid(outgoing_teid),
        vlan(out_vlan),
        ue_state(set_ue_state) {}

  void get_enb_v4_addr(struct in_addr* enb_ipv4_addr) {
    inet_pton(AF_INET, enb_v4.c_str(), enb_ipv4_addr);
  }

  void get_ue_v4_addr(struct in_addr* ue_ipv4_addr) {
    inet_pton(AF_INET, ue_v4.c_str(), ue_ipv4_addr);
  }

  uint32_t get_in_teid() { return in_teid; }

  uint32_t get_out_teid() { return out_teid; }

  uint32_t get_vlan() { return vlan; }

  void set_update_request_ipv4() {
    struct in_addr enb_ipv4_addr;
    struct in_addr ue_ipv4_addr;

    if (!enb_v4.empty()) {
      inet_pton(AF_INET, enb_v4.c_str(), &enb_ipv4_addr);
    }

    if (!ue_v4.empty()) {
      inet_pton(AF_INET, ue_v4.c_str(), &ue_ipv4_addr);
    }

    update_request = create_update_request_ipv4(
        enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, vlan, ue_state);
  }

  UESessionSet get_update_request_ipv4() { return update_request; }

  u_int32_t get_config_ue_state(UESessionState sess_state) {
    if (sess_state.ue_config_state() == UESessionState::ACTIVE) {
      return UE_SESSION_ACTIVE_STATE;
    } else if (sess_state.ue_config_state() == UESessionState::UNREGISTERED) {
      return UE_SESSION_UNREGISTERED_STATE;
    } else if (sess_state.ue_config_state() == UESessionState::INSTALL_IDLE) {
      return UE_SESSION_INSTALL_IDLE_STATE;
    } else if (sess_state.ue_config_state() == UESessionState::UNINSTALL_IDLE) {
      return UE_SESSION_UNINSTALL_IDLE_STATE;
    } else if (sess_state.ue_config_state() == UESessionState::SUSPENDED_DATA) {
      return UE_SESSION_SUSPENDED_DATA_STATE;
    } else if (sess_state.ue_config_state() == UESessionState::RESUME_DATA) {
      return UE_SESSION_RESUME_DATA_STATE;
    }

    return 0;
  }

  bool validate_req_msg(UESessionSet request) {
    char req_enb_ip_str[INET_ADDRSTRLEN];
    char req_ue_ip_str[INET_ADDRSTRLEN];

    inet_ntop(
        AF_INET, request.enb_ip_address().address().c_str(), req_enb_ip_str,
        INET_ADDRSTRLEN);

    if (std::strcmp(req_enb_ip_str, enb_v4.c_str())) {
      return false;
    }

    inet_ntop(
        AF_INET, (request.ue_ipv4_address().address().c_str()), req_ue_ip_str,
        INET_ADDRSTRLEN);

    if (std::strcmp(req_ue_ip_str, ue_v4.c_str())) {
      return false;
    }

    if (in_teid != request.in_teid()) {
      return false;
    }

    if (out_teid != request.out_teid()) {
      return false;
    }

    if (ue_state != get_config_ue_state(request.ue_session_state())) {
      return false;
    }

    if (vlan != request.vlan()) {
      return false;
    }

    return true;
  }
};

class FlowDLOps {
 public:
  // Constructor definition
  FlowDLOps()
      : set_params_((DST_IPV4 | SRC_IPV4)),
        tcp_dst_port_(5002),
        tcp_src_port_(60),
        udp_dst_port_(0),
        udp_src_port_(0),
        ip_proto_(6),
        dst_v4_("192.168.60.141"),
        src_v4_("192.168.128.11") {}

  struct ip_flow_dl get_flow_dl() {
    return flow_dl_;
  }

  void set_flow_dl() {
    flow_dl_.set_params   = set_params_;
    flow_dl_.tcp_dst_port = tcp_dst_port_;
    flow_dl_.tcp_src_port = tcp_src_port_;
    flow_dl_.udp_dst_port = udp_dst_port_;
    flow_dl_.udp_src_port = udp_src_port_;
    flow_dl_.ip_proto     = ip_proto_;
    inet_pton(AF_INET, dst_v4_.c_str(), &(flow_dl_.dst_ip));
    inet_pton(AF_INET, src_v4_.c_str(), &(flow_dl_.src_ip));
  }

  bool validate_flow_dl(IPFlowDL req_flow_dl) {
    IPAddress* dst_ip_addr;
    char flow_dl_dst_addr[INET_ADDRSTRLEN];
    IPAddress* src_ip_addr;
    char flow_dl_src_addr[INET_ADDRSTRLEN];

    dst_ip_addr = req_flow_dl.mutable_dest_ip();
    inet_ntop(
        AF_INET, (dst_ip_addr->address().c_str()), flow_dl_dst_addr,
        INET_ADDRSTRLEN);

    if (std::strcmp(flow_dl_dst_addr, dst_v4_.c_str())) {
      return false;
    }

    src_ip_addr = req_flow_dl.mutable_src_ip();
    inet_ntop(
        AF_INET, (src_ip_addr->address().c_str()), flow_dl_src_addr,
        INET_ADDRSTRLEN);

    if (std::strcmp(flow_dl_src_addr, src_v4_.c_str())) {
      return false;
    }

    if (req_flow_dl.set_params() != set_params_) {
      return false;
    }

    if (req_flow_dl.tcp_dst_port() != tcp_dst_port_) {
      return false;
    }

    if (req_flow_dl.tcp_src_port() != tcp_src_port_) {
      return false;
    }

    if (req_flow_dl.udp_dst_port() != udp_dst_port_) {
      return false;
    }

    if (req_flow_dl.udp_src_port() != udp_src_port_) {
      return false;
    }

    if (req_flow_dl.ip_proto() != ip_proto_) {
      return false;
    }

    return true;
  }

 private:
  uint32_t set_params_;
  uint16_t tcp_dst_port_;
  uint16_t tcp_src_port_;
  uint16_t udp_dst_port_;
  uint16_t udp_src_port_;
  uint8_t ip_proto_;
  std::string dst_v4_;
  std::string src_v4_;
  struct ip_flow_dl flow_dl_;
};

}  // namespace lte
}  // namespace magma
