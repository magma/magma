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
#include <gtest/gtest-message.h>    // for Message
#include <gtest/gtest-test-part.h>  // for TestPartResult
#include <gtest/gtest.h>            // for TestInfo (ptr only), InitGoog...
#include <stddef.h>                 // for NULL
#include <stdint.h>                 // for uint8_t
#include <fluid/base/EventLoop.hh>  // for fluid_base
#include <fluid/util/ethaddr.hh>    // for fluid_msg
#include <memory>                   // for unique_ptr
// TODO: Once #5146 is resolved this can be re-ordered above <memory>
#include <fluid/OFConnection.hh>        // for OFConnection, OFConnection::E...
#include <stdexcept>                    // for runtime_error
#include "ControllerEvents.h"           // for EVENT_PACKET_IN, EVENT_SWITCH...
#include "OpenflowController.h"         // for OpenflowController, OFPT_PACK...
#include "gmock/gmock-matchers.h"       // for AnythingMatcher, _
#include "gmock/gmock-spec-builders.h"  // for EXPECT_CALL, TypedExpectation
#include "gmock/gmock.h"                // for InitGoogleMock
#include "openflow_mocks.h"             // for MockApplication

using ::testing::_;
using ::testing::InSequence;
using ::testing::Test;
using namespace openflow;
using namespace fluid_msg;
using namespace fluid_base;

namespace {

/**
 * Test fixture that instantiates an openflow controller for testing.
 */
class ControllerTest : public ::testing::Test {
 public:
  void default_message_callback(uint8_t type) {
    controller->message_callback(
        NULL,  // no connection
        type,
        NULL,  // no data
        0);    // no length
  }

  void default_connection_callback(OFConnection::Event type) {
    controller->connection_callback(NULL, type);
  }

 protected:
  virtual void SetUp() {
    controller = std::unique_ptr<OpenflowController>(
        new OpenflowController("127.0.0.1", 6666, 2, false));
  }

  virtual void TearDown() { controller = NULL; }

  std::unique_ptr<OpenflowController> controller;
};

// Test simple registration and see if application received event
TEST_F(ControllerTest, TestRegistration) {
  MockApplication app;
  EXPECT_CALL(app, event_callback(_, _)).Times(1);

  controller->register_for_event((Application*) &app, EVENT_PACKET_IN);
  default_message_callback(OFPT_PACKET_IN_TYPE);
}

// Ensure controller can only handle events if it's running
TEST_F(ControllerTest, TestRunningAssertion) {
  controller->stop();
  EXPECT_THROW(
      controller->message_callback(NULL, OFPT_PACKET_IN_TYPE, NULL, 0),
      std::runtime_error);
}

// Test that with multiple apps registered, the correct apps receive the
// correct events in order.
TEST_F(ControllerTest, TestMultipleApplications) {
  MockApplication app1, app2, app3;
  controller->register_for_event(&app1, openflow::EVENT_PACKET_IN);
  controller->register_for_event(&app2, openflow::EVENT_SWITCH_UP);
  controller->register_for_event(&app3, openflow::EVENT_SWITCH_DOWN);
  {
    InSequence dummy;
    // Initial flow
    EXPECT_CALL(app2, event_callback(_, _)).Times(1);
    EXPECT_CALL(app1, event_callback(_, _)).Times(1);
    EXPECT_CALL(app3, event_callback(_, _)).Times(1);
  }
  default_message_callback(OFPT_FEATURES_REPLY_TYPE);
  default_message_callback(OFPT_PACKET_IN_TYPE);
  default_connection_callback(OFConnection::EVENT_CLOSED);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  ::testing::InitGoogleMock(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace
