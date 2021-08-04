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
#include <memory>
#include <vector>
#include "StatsPoller.h"
#define OVS_COOKIE_MATCH_ALL 0xffffffff

namespace magma {

void StatsPoller::start_loop(
    std::shared_ptr<magma::LocalEnforcer> local_enforcer,
    std::shared_ptr<magma::ShardTracker> shard_tracker) {
  // when no shard ids are available poll all UEs
  // to poll by cookie pass in a vector of shard ids, where
  // each element is a shard id that maps directly to a cookie
  while (true) {
    // TODO(veshkemburu): add smart polling function to determine which
    // shard ids to poll
    std::vector<int> active_shard_ids = shard_tracker->get_active_shards();
    uint32_t interval;
    if (active_shard_ids.size() == 0) {
      interval = 10;
    } else {
      interval = 10 / active_shard_ids.size();
    }

    if (active_shard_ids.size() == 0) {
      local_enforcer->poll_stats_enforcer(0, 0);
      std::this_thread::sleep_for(std::chrono::seconds(interval));
    } else {
      for (size_t shard_index = 0; shard_index < active_shard_ids.size();
           shard_index++) {
        local_enforcer->poll_stats_enforcer(
            active_shard_ids[shard_index], OVS_COOKIE_MATCH_ALL);
        std::this_thread::sleep_for(std::chrono::seconds(interval));
      }
    }
  }
}

}  // namespace magma
