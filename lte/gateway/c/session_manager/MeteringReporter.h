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
#pragma once

#include <string>

#include "StoredState.h"
#include "SessionCredit.h"

namespace magma {
namespace lte {

class MeteringReporter {
 public:
  MeteringReporter();

  /**
   * Report all unreported traffic usage for a session.
   * All charging and monitoring keys are aggregated.
   */
  void report_usage(
      const std::string& imsi, const std::string& session_id,
      SessionStateUpdateCriteria& update_criteria);

  /**
   * Reports the usage as described in TotalCreditUsage
   * This function is intended to be used on service restart to offset the
   * counter value. TotalCreditUsage contains the cumulative usage since the
   * start, not a delta value.
   */
  void initialize_usage(
      const std::string& imsi, const std::string& session_id,
      TotalCreditUsage usage);

 private:
  /**
   * Report traffic usage for a session
   */
  void report_traffic(
      const std::string& imsi, const std::string& session_id,
      const std::string& traffic_direction, double unreported_usage_bytes);
};

}  // namespace lte
}  // namespace magma
