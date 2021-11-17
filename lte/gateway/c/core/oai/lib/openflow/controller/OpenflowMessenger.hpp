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

#include <fluid/of13msg.hh>
#include <fluid/OFServer.hh>

namespace openflow {
/**
 * Abstract helper class with libfluid message utilities
 */
class OpenflowMessenger {
 public:
  /**
   * Create a standard flow mod, where every parameter is the default except the
   * table id, command, and priority.
   *
   * @param table_id - table number to add flow to
   * @param command - type of flow mod, e.g. add/delete
   * @param priority - int priority of the flow
   * @return a filled-in FlowMod object
   */
  virtual fluid_msg::of13::FlowMod create_default_flow_mod(
      uint8_t table_id, fluid_msg::of13::ofp_flow_mod_command command,
      uint16_t priority) const = 0;

  /**
   * Sends a completed flow modification to OVS
   *
   * @param flow_mod - a flow modification (add/delete) to make
   * @param ofconn - the connection to send the flow mod to
   */
  virtual void send_of_msg(
      fluid_msg::OFMsg& of_msg, fluid_base::OFConnection* ofconn) const {}
};

/**
 * Implemented messenger class
 */
class DefaultMessenger : public OpenflowMessenger {
 public:
  fluid_msg::of13::FlowMod create_default_flow_mod(
      uint8_t table_id, fluid_msg::of13::ofp_flow_mod_command command,
      uint16_t priority) const;

  void send_of_msg(
      fluid_msg::OFMsg& of_msg, fluid_base::OFConnection* ofconn) const;
};

}  // namespace openflow
