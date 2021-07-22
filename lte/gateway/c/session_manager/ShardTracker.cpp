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

int ShardTracker::add_ue(std::string imsi) {
  if (shards_.empty()) {
    std::set<std::string> new_shard = {imsi};
    shards_.push_back(new_shard);
  }
  for (size_t shard_id = 0; shard_id < shards_.size(); shard_id++) {
    /**
     * If the UE is already in the shard, return the shard id. This check
     * is meant to avoid multiple sessions for a UE being assigned duplicate
     * entries.
     * */
    if (shards_[shard_id].find(imsi) != shards_[shard_id].end()) {
      return shard_id;
    }
    /**
     * If the shard has space, insert the UE
     * */
    if (shards_[shard_id].size() < max_shard_size_) {
      shards_[shard_id].insert(imsi);
      return shard_id;
    }
  }
  /**
   * If all shards are filled, add a new shard and insert the UE,
   * return it's index
   * */
  std::set<std::string> new_shard = {imsi};
  shards_.push_back(new_shard);
  return shards_.size() - 1;
}

bool ShardTracker::remove_ue(std::string imsi, int shard_id) {
  /**
   * Check if there are any shards, any UEs at a particular shard,
   * and whether the UE is actually part of the shard, before removal
   * */
  if (shards_.empty() || shards_[shard_id].empty() ||
      (shards_[shard_id].find(imsi) == shards_[shard_id].end())) {
    return false;
  }
  shards_[shard_id].erase(imsi);
  return true;
}

}  // namespace magma
