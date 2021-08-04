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
#include <iostream>

#include "ShardTracker.h"
#include "magma_logging.h"

namespace magma {

ShardTracker::ShardTracker() {
  // initialize with at least one empty vector to avoid checking for empty
  imsis_per_shard_.push_back({});
}

uint16_t ShardTracker::add_ue(const std::string imsi) {
  MLOG(MINFO) << "The max size for a shard is: " << max_shard_size_;
  for (size_t shard_id = 0; shard_id < get_shard_list_size(); shard_id++) {
    // If the UE is already in the shard, return the shard id. This check
    // is meant to avoid multiple sessions for a UE being assigned duplicate
    // entries.
    if (imsis_per_shard_[shard_id].find(imsi) !=
        imsis_per_shard_[shard_id].end()) {
      return shard_id;
    }
    std::set<std::string>& imsis_at_id = imsis_per_shard_[shard_id];
    // If the shard has space, insert the UE
    if (imsis_at_id.size() < max_shard_size_) {
      imsis_per_shard_[shard_id].insert(imsi);
      std::cout << "\nThe size of this shard entry is: \n"
                << imsis_per_shard_[shard_id].size() << "\n";
      return shard_id;
    }
  }
  // If all shards are filled, add a new shard and insert the UE,
  // return it's index
  imsis_per_shard_.push_back({imsi});
  std::cout << "\nThe size of this shard entry is: "
            << imsis_per_shard_[imsis_per_shard_.size() - 1].size() << "\n";
  return imsis_per_shard_.size() - 1;
}

bool ShardTracker::remove_ue(const std::string imsi, const uint16_t shard_id) {
  // Check if the shard id exists(shard ids are index based),
  // and whether the UE is actually part of the shard, before removal
  if (shard_id >= imsis_per_shard_.size() ||
      imsis_per_shard_[shard_id].empty()) {
    return false;
  }
  imsis_per_shard_[shard_id].erase(imsi);
  return true;
}

uint16_t ShardTracker::get_shard_list_size() {
  return imsis_per_shard_.size();
}

std::vector<int> ShardTracker::get_active_shards() {
  std::vector<int> active_shard_ids;
  for (size_t i = 0; i < imsis_per_shard_.size(); i++) {
    if (!imsis_per_shard_[i].empty()) {
      active_shard_ids.push_back(i);
    }
  }
  return active_shard_ids;
}

}  // namespace magma
