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

#include <string>

#include "includes/MagmaService.h"
#include "MeteringReporter.h"
#include "includes/MetricsHelpers.h"

using magma::service303::increment_counter;

namespace magma {
namespace lte {

const char* COUNTER_NAME     = "ue_traffic";
const char* LABEL_IMSI       = "IMSI";
const char* LABEL_SESSION_ID = "session_id";
const char* LABEL_DIRECTION  = "direction";
const char* DIRECTION_UP     = "up";
const char* DIRECTION_DOWN   = "down";

MeteringReporter::MeteringReporter() {}

void MeteringReporter::report_usage(
    const std::string& imsi, const std::string& session_id,
    SessionStateUpdateCriteria& update_criteria) {
  double total_tx = 0;
  double total_rx = 0;

  // Charging credit
  for (const auto& it : update_criteria.charging_credit_map) {
    auto credit_update = it.second;
    total_tx += (double) credit_update.bucket_deltas[USED_TX];
    total_rx += (double) credit_update.bucket_deltas[USED_RX];
  }

  // Monitoring credit
  for (const auto& it : update_criteria.monitor_credit_map) {
    auto credit_update = it.second;
    total_tx += (double) credit_update.bucket_deltas[USED_TX];
    total_rx += (double) credit_update.bucket_deltas[USED_RX];
  }

  report_traffic(imsi, session_id, DIRECTION_UP, total_tx);
  report_traffic(imsi, session_id, DIRECTION_DOWN, total_rx);
}

void MeteringReporter::initialize_usage(
    const std::string& imsi, const std::string& session_id,
    TotalCreditUsage usage) {
  auto tx = usage.monitoring_tx + usage.charging_tx;
  auto rx = usage.monitoring_rx + usage.charging_rx;
  report_traffic(imsi, session_id, DIRECTION_UP, tx);
  report_traffic(imsi, session_id, DIRECTION_DOWN, rx);
}

void MeteringReporter::report_traffic(
    const std::string& imsi, const std::string& session_id,
    const std::string& traffic_direction, double unreported_usage_bytes) {
  increment_counter(
      COUNTER_NAME, unreported_usage_bytes, size_t(3), LABEL_IMSI, imsi.c_str(),
      LABEL_SESSION_ID, session_id.c_str(), LABEL_DIRECTION,
      traffic_direction.c_str());
}

}  // namespace lte
}  // namespace magma
