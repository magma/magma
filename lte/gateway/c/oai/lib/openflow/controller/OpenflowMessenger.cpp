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

#include "OpenflowMessenger.h"

namespace openflow {

fluid_msg::of13::FlowMod DefaultMessenger::create_default_flow_mod(
    uint8_t table_id, fluid_msg::of13::ofp_flow_mod_command command,
    uint16_t priority) const {
  fluid_msg::of13::FlowMod fm;
  // Defaults
  fm.xid(1);                           // Transaction id, can be anything
  fm.cookie(0);                        // Not used
  fm.cookie_mask(0xffffffffffffffff);  // Not used
  fm.buffer_id(OFP_NO_BUFFER);         // Not used
  fm.out_port(0);                      // Default to not going out of a port
  fm.out_group(0);                     // Groups not used
  fm.flags(0);                         // Not used

  // Custom
  fm.table_id(table_id);
  fm.command(command);
  fm.priority(priority);
  return fm;
}

void DefaultMessenger::send_of_msg(
    fluid_msg::OFMsg& of_msg, fluid_base::OFConnection* ofconn) const {
  uint8_t* buffer;
  buffer = of_msg.pack();
  ofconn->send(buffer, of_msg.length());
  // TODO OF_ERROR_HANDLING - check if OF message successfully installed
  fluid_msg::OFMsg::free_buffer(buffer);
}

}  // namespace openflow
