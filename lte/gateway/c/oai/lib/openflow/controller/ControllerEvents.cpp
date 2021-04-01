/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <netinet/in.h>
#include <string.h>
#include "ControllerEvents.h"

using namespace fluid_msg;

namespace openflow {

ControllerEvent::ControllerEvent(
    fluid_base::OFConnection* ofconn, const ControllerEventType type)
    : ofconn_(ofconn), type_(type) {}

const ControllerEventType ControllerEvent::get_type() const {
  return type_;
}

fluid_base::OFConnection* ControllerEvent::get_connection() const {
  return ofconn_;
}

DataEvent::DataEvent(
    fluid_base::OFConnection* ofconn, fluid_base::OFHandler& ofhandler,
    const void* data, const size_t len, const ControllerEventType type)
    : ControllerEvent(ofconn, type),
      ofhandler_(ofhandler),
      data_(static_cast<const uint8_t*>(data)),
      len_(len) {}

DataEvent::~DataEvent() {
  ofhandler_.free_data(const_cast<uint8_t*>(data_));
}

const uint8_t* DataEvent::get_data() const {
  return data_;
}

const size_t DataEvent::get_length() const {
  return len_;
}

PacketInEvent::PacketInEvent(
    fluid_base::OFConnection* ofconn, fluid_base::OFHandler& ofhandler,
    const void* data, const size_t len)
    : DataEvent(ofconn, ofhandler, data, len, EVENT_PACKET_IN) {}

SwitchUpEvent::SwitchUpEvent(
    fluid_base::OFConnection* ofconn, fluid_base::OFHandler& ofhandler,
    const void* data, const size_t len)
    : DataEvent(ofconn, ofhandler, data, len, EVENT_SWITCH_UP) {}

SwitchDownEvent::SwitchDownEvent(fluid_base::OFConnection* ofconn)
    : ControllerEvent(ofconn, EVENT_SWITCH_DOWN) {}

ErrorEvent::ErrorEvent(
    fluid_base::OFConnection* ofconn, const struct ofp_error_msg* error_msg)
    : error_type_(ntohs(error_msg->type)),
      error_code_(ntohs(error_msg->code)),
      ControllerEvent(ofconn, EVENT_ERROR) {}

const uint16_t ErrorEvent::get_error_type() const {
  return error_type_;
}

const uint16_t ErrorEvent::get_error_code() const {
  return error_code_;
}

ExternalEvent::ExternalEvent(const ControllerEventType type)
    : ControllerEvent(NULL, type) {}

void ExternalEvent::set_of_connection(fluid_base::OFConnection* ofconn) {
  ofconn_ = ofconn;
}

UeNetworkInfo::UeNetworkInfo(const struct in_addr ue_ip)
    : ue_ip_(ue_ip), vlan_(0), ue_ipv6_(in6addr_any) {}

UeNetworkInfo::UeNetworkInfo(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6)
    : ue_ip_(ue_ip), vlan_(0) {
  if (ue_ipv6) {
    ue_ipv6_ = *ue_ipv6;
  } else {
    ue_ipv6_ = in6addr_any;
  }
}

UeNetworkInfo::UeNetworkInfo(const struct in_addr ue_ip, int vlan)
    : ue_ip_(ue_ip), vlan_(vlan), ue_ipv6_(in6addr_any) {}

UeNetworkInfo::UeNetworkInfo(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, int vlan)
    : ue_ip_(ue_ip), vlan_(vlan) {
  if (ue_ipv6) {
    ue_ipv6_ = *ue_ipv6;
  } else {
    ue_ipv6_ = in6addr_any;
  }
}

const bool UeNetworkInfo::is_ue_ipv6_addr_valid() const {
  return !!memcmp(&ue_ipv6_, &in6addr_any, sizeof(ue_ipv6_));
}

const bool UeNetworkInfo::is_ue_ipv4_addr_valid() const {
  return ue_ip_.s_addr != INADDR_ANY;
}

const struct in_addr& UeNetworkInfo::get_ip() const {
  return ue_ip_;
}

const struct in6_addr& UeNetworkInfo::get_ipv6() const {
  return ue_ipv6_;
}

const int UeNetworkInfo::get_vlan() const {
  return vlan_;
}

AddGTPTunnelEvent::AddGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, int vlan,
    const struct in_addr enb_ip, const uint32_t in_tei, const uint32_t out_tei,
    const char* imsi, uint32_t enb_gtp_port)
    : ue_info_(ue_ip, ue_ipv6, vlan),
      enb_ip_(enb_ip),
      pgw_ip_(INADDR_ZERO),
      in_tei_(in_tei),
      out_tei_(out_tei),
      imsi_(imsi),
      dl_flow_valid_(false),
      dl_flow_(),
      dl_flow_precedence_(DEFAULT_PRECEDENCE),
      ExternalEvent(EVENT_ADD_GTP_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(0) {}

AddGTPTunnelEvent::AddGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, int vlan,
    const struct in_addr enb_ip, const uint32_t in_tei, const uint32_t out_tei,
    const char* imsi, const struct ip_flow_dl* dl_flow,
    const uint32_t dl_flow_precedence, uint32_t enb_gtp_port)
    : ue_info_(ue_ip, vlan),
      enb_ip_(enb_ip),
      pgw_ip_(INADDR_ZERO),
      in_tei_(in_tei),
      out_tei_(out_tei),
      imsi_(imsi),
      dl_flow_valid_(true),
      dl_flow_(*dl_flow),
      dl_flow_precedence_(dl_flow_precedence),
      ExternalEvent(EVENT_ADD_GTP_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(0) {}

AddGTPTunnelEvent::AddGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, int vlan,
    const struct in_addr enb_ip, const struct in_addr pgw_ip,
    const uint32_t in_tei, const uint32_t out_tei, const char* imsi,
    const struct ip_flow_dl* dl_flow, const uint32_t dl_flow_precedence,
    uint32_t enb_gtp_port, uint32_t pgw_gtp_port)
    : ue_info_(ue_ip, ue_ipv6, vlan),
      enb_ip_(enb_ip),
      pgw_ip_(pgw_ip),
      in_tei_(in_tei),
      out_tei_(out_tei),
      imsi_(imsi),
      dl_flow_valid_(false),
      dl_flow_(),
      dl_flow_precedence_(DEFAULT_PRECEDENCE),
      ExternalEvent(EVENT_ADD_GTP_S8_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(pgw_gtp_port) {}

AddGTPTunnelEvent::AddGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, int vlan,
    const struct in_addr enb_ip, const struct in_addr pgw_ip,
    const uint32_t in_tei, const uint32_t out_tei, const char* imsi,
    uint32_t enb_gtp_port, uint32_t pgw_gtp_port)
    : ue_info_(ue_ip, vlan),
      enb_ip_(enb_ip),
      pgw_ip_(pgw_ip),
      in_tei_(in_tei),
      out_tei_(out_tei),
      imsi_(imsi),
      dl_flow_valid_(false),
      dl_flow_(),
      dl_flow_precedence_(DEFAULT_PRECEDENCE),
      ExternalEvent(EVENT_ADD_GTP_S8_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(pgw_gtp_port) {}

const struct in_addr& AddGTPTunnelEvent::get_ue_ip() const {
  return ue_info_.get_ip();
}

const struct UeNetworkInfo& AddGTPTunnelEvent::get_ue_info() const {
  return ue_info_;
}

const struct in_addr& AddGTPTunnelEvent::get_enb_ip() const {
  return enb_ip_;
}

const struct in_addr& AddGTPTunnelEvent::get_pgw_ip() const {
  return pgw_ip_;
}

const uint32_t AddGTPTunnelEvent::get_in_tei() const {
  return in_tei_;
}

const uint32_t AddGTPTunnelEvent::get_out_tei() const {
  return out_tei_;
}

const std::string& AddGTPTunnelEvent::get_imsi() const {
  return imsi_;
}

const bool AddGTPTunnelEvent::is_dl_flow_valid() const {
  return dl_flow_valid_;
}

const struct ip_flow_dl& AddGTPTunnelEvent::get_dl_flow() const {
  return dl_flow_;
}

const uint32_t AddGTPTunnelEvent::get_dl_flow_precedence() const {
  return dl_flow_precedence_;
}

const uint32_t AddGTPTunnelEvent::get_enb_gtp_portno() const {
  return enb_gtp_port_;
}

const uint32_t AddGTPTunnelEvent::get_pgw_gtp_portno() const {
  return pgw_gtp_port_;
}

DeleteGTPTunnelEvent::DeleteGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, const uint32_t in_tei,
    const struct ip_flow_dl* dl_flow, uint32_t enb_gtp_port)
    : ue_info_(ue_ip, ue_ipv6),
      in_tei_(in_tei),
      dl_flow_valid_(true),
      dl_flow_(*dl_flow),
      ExternalEvent(EVENT_DELETE_GTP_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(0) {}

DeleteGTPTunnelEvent::DeleteGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, const uint32_t in_tei,
    uint32_t enb_gtp_port)
    : ue_info_(ue_ip, ue_ipv6),
      in_tei_(in_tei),
      dl_flow_valid_(false),
      dl_flow_(),
      ExternalEvent(EVENT_DELETE_GTP_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(0) {}

DeleteGTPTunnelEvent::DeleteGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, const uint32_t in_tei,
    const struct ip_flow_dl* dl_flow, uint32_t enb_gtp_port,
    uint32_t pgw_gtp_port)
    : ue_info_(ue_ip, ue_ipv6),
      in_tei_(in_tei),
      dl_flow_valid_(true),
      dl_flow_(*dl_flow),
      ExternalEvent(EVENT_DELETE_GTP_S8_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(pgw_gtp_port) {}

DeleteGTPTunnelEvent::DeleteGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, const uint32_t in_tei,
    uint32_t enb_gtp_port, uint32_t pgw_gtp_port)
    : ue_info_(ue_ip, ue_ipv6),
      in_tei_(in_tei),
      dl_flow_valid_(false),
      dl_flow_(),
      ExternalEvent(EVENT_DELETE_GTP_S8_TUNNEL),
      enb_gtp_port_(enb_gtp_port),
      pgw_gtp_port_(pgw_gtp_port) {}

const struct UeNetworkInfo& DeleteGTPTunnelEvent::get_ue_info() const {
  return ue_info_;
}

const struct in_addr& DeleteGTPTunnelEvent::get_ue_ip() const {
  return ue_info_.get_ip();
}

const uint32_t DeleteGTPTunnelEvent::get_in_tei() const {
  return in_tei_;
}

const bool DeleteGTPTunnelEvent::is_dl_flow_valid() const {
  return dl_flow_valid_;
}

const struct ip_flow_dl& DeleteGTPTunnelEvent::get_dl_flow() const {
  return dl_flow_;
}

const uint32_t DeleteGTPTunnelEvent::get_enb_gtp_portno() const {
  return enb_gtp_port_;
}

const uint32_t DeleteGTPTunnelEvent::get_pgw_gtp_portno() const {
  return pgw_gtp_port_;
}

HandleDataOnGTPTunnelEvent::HandleDataOnGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, const uint32_t in_tei,
    const ControllerEventType event_type, const struct ip_flow_dl* dl_flow,
    const uint32_t dl_flow_precedence)
    : ue_info_(ue_ip),
      in_tei_(in_tei),
      dl_flow_valid_(true),
      dl_flow_(*dl_flow),
      dl_flow_precedence_(dl_flow_precedence),
      ExternalEvent(event_type) {}

HandleDataOnGTPTunnelEvent::HandleDataOnGTPTunnelEvent(
    const struct in_addr ue_ip, struct in6_addr* ue_ipv6, const uint32_t in_tei,
    const ControllerEventType event_type)
    : ue_info_(ue_ip),
      in_tei_(in_tei),
      dl_flow_valid_(false),
      dl_flow_(),
      dl_flow_precedence_(DEFAULT_PRECEDENCE),
      ExternalEvent(event_type) {}

const struct in_addr& HandleDataOnGTPTunnelEvent::get_ue_ip() const {
  return ue_info_.get_ip();
}

const uint32_t HandleDataOnGTPTunnelEvent::get_in_tei() const {
  return in_tei_;
}

const bool HandleDataOnGTPTunnelEvent::is_dl_flow_valid() const {
  return dl_flow_valid_;
}

const struct ip_flow_dl& HandleDataOnGTPTunnelEvent::get_dl_flow() const {
  return dl_flow_;
}

const uint32_t HandleDataOnGTPTunnelEvent::get_dl_flow_precedence() const {
  return dl_flow_precedence_;
}

AddPagingRuleEvent::AddPagingRuleEvent(const struct in_addr ue_ip)
    : ue_info_(ue_ip), ExternalEvent(EVENT_ADD_PAGING_RULE) {}

const struct in_addr& AddPagingRuleEvent::get_ue_ip() const {
  return ue_info_.get_ip();
}

DeletePagingRuleEvent::DeletePagingRuleEvent(const struct in_addr ue_ip)
    : ue_info_(ue_ip), ExternalEvent(EVENT_DELETE_PAGING_RULE) {}

const struct in_addr& DeletePagingRuleEvent::get_ue_ip() const {
  return ue_info_.get_ip();
}

}  // namespace openflow
