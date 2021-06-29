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
#include <stdint.h>
#include <atomic>
#include <functional>
#include <vector>
#include "lte/protos/subscriberdb.pb.h"
#include "StatsPoller.h"

namespace magma {

void StatsPoller::start_loop(
    std::shared_ptr<magma::LocalEnforcer> local_enforcer,
    uint32_t loop_interval_seconds) {
  is_running_ = true;
  while (is_running_) {
    local_enforcer->poll_stats_enforcer();
    std::this_thread::sleep_for(std::chrono::seconds(loop_interval_seconds));
  }
}

void StatsPoller::stop() {
  is_running_ = false;
}

}  // namespace magma
