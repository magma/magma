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

#include "OpenflowController.h"

namespace openflow {
#define ETH_HEADER_LENGTH 14

class PagingApplication : public Application {
 private:
  static const int MID_PRIORITY = 5;
  // TODO: move to config file
  static const int CLAMPING_TIMEOUT = 30;  // seconds
  /**
   * Main callback event required by inherited Application class. Whenever
   * the controller gets an event like packet in or switch up, it will pass
   * it to the application here
   *
   * @param ev (in) - pointer to some subclass of ControllerEvent that occurred
   */
  virtual void event_callback(
      const ControllerEvent& ev, const OpenflowMessenger& messenger);

  /**
   * Handles downlink data intended for a UE in idle mode, then forwards the
   * paging request to SPGW. After initiating the paging process, it also clamps
   * on the destination IP, to prevent multiple packet-in messages
   *
   * @param ofconn (in) - given connection to OVS switch
   * @param data (in) - the ethernet packet received by the switch
   */
  void handle_paging_message(
      fluid_base::OFConnection* ofconn, uint8_t* data,
      const OpenflowMessenger& messenger);

  /**
   * Creates exact paging flow, which sends a packet intended for an
   * idle UE to this application
   */
  void add_paging_flow(
      const AddPagingRuleEvent& ev, const OpenflowMessenger& messenger);

  /**
   * Removes exact paging flow rule to stop paging UE
   */
  void delete_paging_flow(
      const DeletePagingRuleEvent& ev, const OpenflowMessenger& messenger);
};

}  // namespace openflow
