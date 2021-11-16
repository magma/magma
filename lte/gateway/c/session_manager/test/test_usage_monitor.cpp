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
#include <gtest/gtest.h>

#include <chrono>
#include <thread>

#include "ProtobufCreators.h"
#include "SessiondMocks.h"
#include "SessionState.h"
#include "SessionStateTester.h"

using ::testing::Test;

namespace magma {
TEST_F(SessionStateTest, test_insert_monitor) {
  update_criteria = get_default_update_criteria();
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 0);
  insert_rule(1, "m1", "rule1", STATIC, 0, 0);
  EXPECT_EQ(true, session_state->active_monitored_rules_exist());
  EXPECT_TRUE(std::find(update_criteria.static_rules_to_install.begin(),
                        update_criteria.static_rules_to_install.end(),
                        "rule1") !=
              update_criteria.static_rules_to_install.end());

  receive_credit_from_pcrf("m1", 1024, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(session_state->get_monitor("m1", ALLOWED_TOTAL), 1024);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            1024);
}

// Insert a monitor, then remove. Assert that the update criteria reflects the
// deletion
TEST_F(SessionStateTest, test_remove_monitor) {
  EXPECT_EQ(update_criteria.static_rules_to_install.size(), 0);
  insert_rule(0, "m1", "rule1", STATIC, 0, 0);

  receive_credit_from_pcrf("m1", 1000, MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(session_state->get_monitor("m1", ALLOWED_TOTAL), 1000);
  EXPECT_EQ(update_criteria.monitor_credit_to_install["m1"]
                .credit.buckets[ALLOWED_TOTAL],
            1000);

  session_state->add_rule_usage("rule1", 1, 1000, 0, 0, 0, nullptr);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 1000);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 0);

  // UsageMonitorResponse with 0 credit will mark the monitor to be DISABLED
  receive_credit_from_pcrf("m1", 0, MonitoringLevel::PCC_RULE_LEVEL);

  update_criteria = get_default_update_criteria();

  // add usage to trigger the quota exhaustion
  session_state->add_rule_usage("rule1", 1, 1001, 0, 0, 0, &update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 1001);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 0);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].report_last_credit);

  // check last update will be sent
  UpdateSessionRequest update;
  std::vector<std::unique_ptr<ServiceAction>> actions;
  session_state->get_updates(&update, &actions, &update_criteria);
  EXPECT_EQ(actions.size(), 0);
  EXPECT_EQ(update.usage_monitors_size(), 1);
  auto& single_update = update.usage_monitors(0).update();
  EXPECT_EQ(single_update.level(), MonitoringLevel::PCC_RULE_LEVEL);
  EXPECT_EQ(single_update.bytes_tx(), 1001);
  EXPECT_EQ(single_update.bytes_rx(), 0);
  EXPECT_TRUE(update_criteria.monitor_credit_map["m1"].deleted);

  // check the deletion with apply_update_criteria
  session_state->apply_update_criteria(update_criteria);
  EXPECT_EQ(session_state->get_monitor("m1", USED_TX), 0);
  EXPECT_EQ(session_state->get_monitor("m1", USED_RX), 0);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
