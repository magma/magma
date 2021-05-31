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

#include "includes/MagmaService.h"
#include "MeteringReporter.h"
#include "includes/MetricsSingleton.h"

using magma::orc8r::MetricsContainer;
using ::testing::Test;

namespace magma {

class MeteringReporterTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    reporter = std::make_shared<MeteringReporter>();
    magma_service =
        std::make_shared<service303::MagmaService>("test_service", "1.0");
  }
  bool is_equal(
      io::prometheus::client::LabelPair label_pair, const char*& name,
      const char*& value) {
    return label_pair.name().compare(name) == 0 &&
           label_pair.value().compare(value) == 0;
  }

 protected:
  std::shared_ptr<service303::MagmaService> magma_service;
  std::shared_ptr<MeteringReporter> reporter;
};

TEST_F(MeteringReporterTest, test_reporting) {
  auto IMSI_LABEL       = "IMSI";
  auto SESSION_ID_LABEL = "session_id";
  auto DIRECTION_LABEL  = "direction";

  auto IMSI           = "imsi";
  auto SESSION_ID     = "session_1";
  auto MONITORING_KEY = "mk1";
  auto DIRECTION_UP   = "up";
  auto DIRECTION_DOWN = "down";

  auto UPLOADED_BYTES   = 5;
  auto DOWNLOADED_BYTES = 7;

  auto uc = get_default_update_criteria();
  SessionCreditUpdateCriteria credit_uc{};
  credit_uc.bucket_deltas[USED_TX]      = UPLOADED_BYTES;
  credit_uc.bucket_deltas[USED_RX]      = DOWNLOADED_BYTES;
  uc.monitor_credit_map[MONITORING_KEY] = credit_uc;

  reporter->report_usage(IMSI, SESSION_ID, uc);

  // verify if UE traffic metrics are recorded properly
  MetricsContainer resp;
  magma_service->GetMetrics(nullptr, nullptr, &resp);
  for (auto const& fam : resp.family()) {
    if (fam.name().compare("ue_traffic") == 0) {
      for (auto const& m : fam.metric()) {
        for (auto const& l : m.label()) {
          EXPECT_TRUE(
              is_equal(l, IMSI_LABEL, IMSI) ||
              is_equal(l, SESSION_ID_LABEL, SESSION_ID) ||
              l.name().compare(DIRECTION_LABEL) == 0);

          if (is_equal(l, DIRECTION_LABEL, DIRECTION_UP)) {
            EXPECT_EQ(m.counter().value(), UPLOADED_BYTES);
          } else if (is_equal(l, DIRECTION_LABEL, DIRECTION_DOWN)) {
            EXPECT_EQ(m.counter().value(), DOWNLOADED_BYTES);
          }
        }
      }
      break;
    }
  }
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma
