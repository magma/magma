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

/**
 * BaseApplication is used to set up the generic table 0 flow, that's not
 * specific to an application. The generic flow passes on all unhandled traffic
 * to the next table. Needs to be registered first, because it clears flows
 * on connection
 */
class BaseApplication : public Application {
 public:
  BaseApplication(bool persist_state);

 private:
  bool persist_state_                = false;
  static const uint32_t LOW_PRIORITY = 0;
  static const uint16_t NEXT_TABLE   = 1;
  /**
   * Main callback event required by inherited Application class. Whenever
   * the controller gets an event like packet in or switch up, it will pass
   * it to the application here
   *
   * @param ev (in) - some subclass of ControllerEvent that occurred
   */
  virtual void event_callback(
      const ControllerEvent& ev, const OpenflowMessenger& messenger);

  /**
   * Creates the default table 0 flow, which resubmits to table 1
   */
  void install_default_flow(
      fluid_base::OFConnection* ofconn, const OpenflowMessenger& messenger);

  void remove_all_flows(
      fluid_base::OFConnection* ofconn, const OpenflowMessenger& messenger);

  /**
   * Log all error messages sent to the controller. These can be rejected flows
   * or bad requests.
   *
   * @param ev (in) - Error event containing type and code of OF error
   */
  void handle_error_message(const ErrorEvent& ev);
};

}  // namespace openflow
