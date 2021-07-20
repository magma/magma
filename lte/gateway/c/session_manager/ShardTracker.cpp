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
  int shard_id;
  // If there are no shards, add an entry for the shard
  // representing a single UE at the new shard
  if (ue_count_per_shard_.size() == 0) {
    ue_count_per_shard_.push_back(1);
    return 0;
  } else {
    // Iterate through all existing shards, if all are full
    // create a new shard with quantity 1, otherwise increment
    // the quantity of UEs in the latest shard
    for (size_t i = 0; i < ue_count_per_shard_.size(); i++) {
      if (ue_count_per_shard_[i] < max_shard_size_) {
        ue_count_per_shard_[i]++;
        return i;
      }
    }
    shard_id = ue_count_per_shard_.size();
    ue_count_per_shard_.push_back(1);
  }
  return shard_id;
}

void ShardTracker::remove_ue(int shard_id) {
  // Since we only keep global state of all UEs, we just
  // need to decrement the number of UEs at a particular id
  ue_count_per_shard_[shard_id]--;
}

}  // namespace magma
