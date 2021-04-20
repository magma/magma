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
#include <gmock/gmock.h>

#include "OpenflowController.h"

using namespace openflow;
using namespace fluid_base;

/**
 * Mock application that doens't handle any callbacks. It's just used to
 * inspect if an application has received an event.
 */
class MockApplication : public Application {
 public:
  MOCK_METHOD2(
      event_callback, void(const ControllerEvent&, const OpenflowMessenger&));
};

/**
 * Mock Openflow Messenger that doesn't send any messages to OVS
 */
class MockMessenger : public OpenflowMessenger {
 public:
  inline fluid_msg::of13::FlowMod create_default_flow_mod(
      uint8_t table_id, fluid_msg::of13::ofp_flow_mod_command command,
      uint16_t priority) const {
    // Use DefaultMessenger
    DefaultMessenger messenger;
    return messenger.create_default_flow_mod(table_id, command, priority);
  }

  MOCK_CONST_METHOD2(
      send_of_msg,
      void(fluid_msg::OFMsg& of_msg, fluid_base::OFConnection* ofconn));
};
