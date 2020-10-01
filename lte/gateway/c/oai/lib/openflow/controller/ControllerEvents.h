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

#pragma once

#include <arpa/inet.h>
#include <fluid/OFServer.hh>
#include <fluid/ofcommon/openflow-common.hh>
#include "gtpv1u.h"

using namespace fluid_msg;

namespace openflow {

enum ControllerEventType {
  EVENT_PACKET_IN,
  EVENT_SWITCH_DOWN,
  EVENT_SWITCH_UP,
  EVENT_ERROR,
  EVENT_ADD_GTP_TUNNEL,
  EVENT_DELETE_GTP_TUNNEL,
  EVENT_DISCARD_DATA_ON_GTP_TUNNEL,
  EVENT_FORWARD_DATA_ON_GTP_TUNNEL,
  EVENT_ADD_PAGING_RULE,
  EVENT_DELETE_PAGING_RULE,
};

/**
 * Superclass for all controller events. These classes are used to pass info
 * like the connection from the controller or any external source to the
 * application
 */
class ControllerEvent {
 public:
  ControllerEvent(
      fluid_base::OFConnection* ofconn, const ControllerEventType type);

  virtual ~ControllerEvent() {}

  fluid_base::OFConnection* get_connection() const;

  const ControllerEventType get_type() const;

 private:
  const ControllerEventType type_;

 protected:
  fluid_base::OFConnection* ofconn_;
};

/**
 * Superclass for any event that gets data passed in through the event
 */
class DataEvent : public ControllerEvent {
 public:
  DataEvent(
      fluid_base::OFConnection* ofconn, fluid_base::OFHandler& ofhandler,
      const void* data, const size_t len, const ControllerEventType type);

  ~DataEvent();

  const uint8_t* get_data() const;
  const size_t get_length() const;

 private:
  fluid_base::OFHandler& ofhandler_;
  const uint8_t* data_;
  const size_t len_;
};

/**
 * Event triggered when a packet gets pushed to user space
 */
class PacketInEvent : public DataEvent {
 public:
  PacketInEvent(
      fluid_base::OFConnection* ofconn, fluid_base::OFHandler& ofhandler,
      const void* data, const size_t len);
};

/**
 * Event triggered when the controller connects with the switch
 */
class SwitchUpEvent : public DataEvent {
 public:
  SwitchUpEvent(
      fluid_base::OFConnection* ofconn, fluid_base::OFHandler& ofhandler,
      const void* data, const size_t len);
};

/**
 * Event triggered when the controller loses connection with the switch
 */
class SwitchDownEvent : public ControllerEvent {
 public:
  SwitchDownEvent(fluid_base::OFConnection* ofconn);
};

/**
 * Event triggered when there is an openflow error reported from the switch
 */
class ErrorEvent : public ControllerEvent {
 public:
  ErrorEvent(
      fluid_base::OFConnection* ofconn, const struct ofp_error_msg* error_msg);

  const uint16_t get_error_type() const;
  const uint16_t get_error_code() const;

 private:
  const uint16_t error_type_;
  const uint16_t error_code_;
};

/*
 * Event triggered externally, so it allows for delayed assignment of the
 * openflow connection. This way, the controller can set the latest known
 * connection, instead of an external file
 */
class ExternalEvent : public ControllerEvent {
 public:
  ExternalEvent(const ControllerEventType type);

  void set_of_connection(fluid_base::OFConnection* ofconn);
};

/*
 * This object contains info about UE IP and vlan.
 * Together this uniquely identifies a UE.
 */

class UeNetworkInfo {
 public:
  UeNetworkInfo(const struct in_addr ue_ip);
  UeNetworkInfo(const struct in_addr ue_ip, int vlan);

  const struct in_addr& get_ip() const;
  const int get_vlan() const;

 private:
  const struct in_addr ue_ip_;
  const int vlan_;
};

/*
 * Event triggered by SPGW to add a GTP tunnel for a UE
 */
class AddGTPTunnelEvent : public ExternalEvent {
 public:
  AddGTPTunnelEvent(
      const struct in_addr ue_ip, int vlan,  const struct in_addr enb_ip,
      const uint32_t in_tei, const uint32_t out_tei, const char* imsi,
      const struct ipv4flow_dl* dl_flow, const uint32_t dl_flow_precedence,
      uint32_t gtp_port_no);

  AddGTPTunnelEvent(
      const struct in_addr ue_ip, int vlan,  const struct in_addr enb_ip,
      const uint32_t in_tei, const uint32_t out_tei, const char* imsi,
      uint32_t gtp_port_no);

  const struct UeNetworkInfo& get_ue_info() const;
  const struct in_addr& get_ue_ip() const;
  const struct in_addr& get_enb_ip() const;
  const uint32_t get_in_tei() const;
  const uint32_t get_out_tei() const;
  const std::string& get_imsi() const;
  const bool is_dl_flow_valid() const;
  const struct ipv4flow_dl& get_dl_flow() const;
  const uint32_t get_dl_flow_precedence() const;
  const uint32_t get_gtp_portno() const;

 private:
  const UeNetworkInfo ue_info_;
  const struct in_addr enb_ip_;
  const uint32_t in_tei_;
  const uint32_t out_tei_;
  const std::string imsi_;
  const struct ipv4flow_dl dl_flow_;
  const bool dl_flow_valid_;
  const uint32_t dl_flow_precedence_;
  const uint32_t gtp_portno_;
};

/*
 * Event triggered by SPGW to remove a GTP tunnel for a UE on detach
 */
class DeleteGTPTunnelEvent : public ExternalEvent {
 public:
  DeleteGTPTunnelEvent(
      const struct in_addr ue_ip, const uint32_t in_tei,
      const struct ipv4flow_dl* dl_flow,
      uint32_t gtp_port_no);
  DeleteGTPTunnelEvent(const struct in_addr ue_ip, const uint32_t in_tei,
      uint32_t gtp_port_no);

  const struct UeNetworkInfo& get_ue_info() const;
  const struct in_addr& get_ue_ip() const;
  const uint32_t get_in_tei() const;
  const bool is_dl_flow_valid() const;
  const struct ipv4flow_dl& get_dl_flow() const;
  const uint32_t get_gtp_portno() const;

 private:
  const UeNetworkInfo ue_info_;
  const uint32_t in_tei_;
  const struct ipv4flow_dl dl_flow_;
  const bool dl_flow_valid_;
  const uint32_t gtp_portno_;
};

/*
 * Event triggered by SPGW to either Discard/Forward DL data on GTP tunnel
 * identified by sgw-S1u TEID if event_type is set to
 * EVENT_DISCARD_DATA_ON_GTP_TUNNEL; A new rule is set to discard data for the
 * UE if event_type is set to EVENT_FORWARD_DATA_ON_GTP_TUNNEL; Shall delete the
 * previous rule
 */
class HandleDataOnGTPTunnelEvent : public ExternalEvent {
 public:
  HandleDataOnGTPTunnelEvent(
      const struct in_addr ue_ip, const uint32_t in_tei,
      const ControllerEventType event_type, const struct ipv4flow_dl* dl_flow,
      const uint32_t dl_flow_precedence);
  HandleDataOnGTPTunnelEvent(
      const struct in_addr ue_ip, const uint32_t in_tei,
      const ControllerEventType event_type);

  const struct UeNetworkInfo& get_ue_info() const;
  const struct in_addr& get_ue_ip() const;
  const uint32_t get_in_tei() const;
  const bool is_dl_flow_valid() const;
  const struct ipv4flow_dl& get_dl_flow() const;
  const uint32_t get_dl_flow_precedence() const;

 private:
  const UeNetworkInfo ue_info_;
  const uint32_t in_tei_;
  const struct ipv4flow_dl dl_flow_;
  const bool dl_flow_valid_;
  const uint32_t dl_flow_precedence_;
};

/*
 * Event triggered by SPGW to support UE paging when
 * S1 is released (i.e., UE is in IDLE mode)
 */
class AddPagingRuleEvent : public ExternalEvent {
 public:
  AddPagingRuleEvent(const struct in_addr ue_ip);

  const struct UeNetworkInfo& get_ue_info() const;
  const struct in_addr& get_ue_ip() const;

 private:
  const UeNetworkInfo ue_info_;
};

/*
 * Event triggered by SPGW to stop UE paging when
 * UE is detached
 */
class DeletePagingRuleEvent : public ExternalEvent {
 public:
  DeletePagingRuleEvent(const struct in_addr ue_ip);

  const struct UeNetworkInfo& get_ue_info() const;
  const struct in_addr& get_ue_ip() const;

 private:
  const UeNetworkInfo ue_info_;
};

}  // namespace openflow
