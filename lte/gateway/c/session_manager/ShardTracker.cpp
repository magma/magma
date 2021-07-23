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
#include <set>

#include "ShardTracker.h"

namespace magma {

uint16_t ShardTracker::add_ue(const std::string imsi) {
  if (imsis_per_shard_.empty()) {
    imsis_per_shard_.push_back({imsi});
  }
  for (size_t shard_id = 0; shard_id < imsis_per_shard_.size(); shard_id++) {
    // If the UE is already in the shard, return the shard id. This check
    // is meant to avoid multiple sessions for a UE being assigned duplicate
    // entries.
    std::set<std::string>& imsis_at_id = imsis_per_shard_[shard_id];
    if (imsis_at_id.find(imsi) != imsis_at_id.end()) {
      return shard_id;
    }
    // If the shard has space, insert the UE
    if (imsis_at_id.size() < max_shard_size_) {
      imsis_at_id.insert(imsi);
      return shard_id;
    }
  }
  // If all shards are filled, add a new shard and insert the UE,
  // return it's index
  imsis_per_shard_.push_back({imsi});
  return imsis_per_shard_.size() - 1;
}

bool ShardTracker::remove_ue(const std::string imsi, uint16_t shard_id) {
  // Check if there are any shards, any UEs at a particular shard,
  // and whether the UE is actually part of the shard, before removal
  if (imsis_per_shard_.empty() || imsis_per_shard_[shard_id].empty()) {
    return false;
  }
  imsis_per_shard_[shard_id].erase(imsi);
  return true;
}

}  // namespace magma
