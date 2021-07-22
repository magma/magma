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
#include <vector>

#include "ShardTracker.h"

namespace magma {

int ShardTracker::add_ue() {
  /*
   * If there are no shards, add an entry for the shard
   * representing a single UE at the new shard
   */
  if (ue_count_per_shard_.size() == 0) {
    ue_count_per_shard_.push_back(1);
    return 0;
  }
  /*
   * Iterate through all existing shards, if all are full
   * create a new shard with quantity 1, otherwise increment
   * the quantity of UEs in the latest shard
   */
  for (size_t shard_id = 0; shard_id < ue_count_per_shard_.size(); shard_id++) {
    if (ue_count_per_shard_[shard_id] < max_shard_size_) {
      ue_count_per_shard_[shard_id]++;
      return shard_id;
    }
  }
  ue_count_per_shard_.push_back(1);
  return ue_count_per_shard_.size() - 1;
}

bool ShardTracker::remove_ue(int shard_id) {
  /*
   * Since we only keep global state of all UEs, we just
   * need to decrement the number of UEs at a particular id
   * if there are no UEs at the shard, removal should fail
   */
  if (ue_count_per_shard_.empty()) {
    return false;
  }
  ue_count_per_shard_[shard_id]--;
  return true;
}

}  // namespace magma
