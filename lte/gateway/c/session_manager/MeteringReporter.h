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

#include "StoredState.h"

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

 private:
  /**
   * Report upload traffic usage for a session
   */
  void report_upload(
      const std::string& imsi, const std::string& session_id,
      double unreported_usage_bytes);

  /**
   * Report download traffic usage for a session
   */
  void report_download(
      const std::string& imsi, const std::string& session_id,
      double unreported_usage_bytes);

  /**
   * Report traffic usage for a session
   */
  void report_traffic(
      const std::string& imsi, const std::string& session_id,
      const std::string& traffic_direction, double unreported_usage_bytes);
};

}  // namespace lte
}  // namespace magma
