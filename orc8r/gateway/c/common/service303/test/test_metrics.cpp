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
#include "includes/MetricsRegistry.h"
#include <prometheus/registry.h>

using io::prometheus::client::MetricFamily;
using magma::service303::MetricsRegistry;
using prometheus::BuildCounter;
using prometheus::Registry;
using prometheus::detail::CounterBuilder;
using ::testing::Test;

namespace magma {

// Tests the MetricsRegistry properly initializes and retrieves metrics
TEST_F(Test, test_metrics_registry) {
  auto prometheus_registry = std::make_shared<Registry>();
  auto registry = MetricsRegistry<prometheus::Counter, CounterBuilder (&)()>(
      prometheus_registry, BuildCounter);
  EXPECT_EQ(registry.SizeFamilies(), 0);
  EXPECT_EQ(registry.SizeMetrics(), 0);

  // Create two new timeseries that will construct two families and metrics
  registry.Get("test", {});
  registry.Get("another", {{"key", "value"}});
  EXPECT_EQ(registry.SizeFamilies(), 2);
  EXPECT_EQ(registry.SizeMetrics(), 2);

  // This should retrieve the previously constructed family
  registry.Get("test", {});
  EXPECT_EQ(registry.SizeFamilies(), 2);
  EXPECT_EQ(registry.SizeMetrics(), 2);

  // Add new unique timeseries to an existing family
  registry.Get("test", {{"key", "value1"}});
  registry.Get("test", {{"key", "value2"}});
  EXPECT_EQ(registry.SizeFamilies(), 2);
  EXPECT_EQ(registry.SizeMetrics(), 4);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}
}  // namespace magma
