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
#include <utility>
#include <vector>
#include <algorithm>
#include <map>

#include "ShardTracker.h"
#define MAX_SHARD_SIZE 100

namespace magma {

ShardTracker::ShardTracker() : number_of_shards_(0) {}

int ShardTracker::add_ue() {
  int shard_id;
  // If we have no shards, we add an entry for the shard representing
  // of 1 representing a single UE at the new shard*/
  if (shards_.size() == 0) {
    shard_id = 0;
    shards_.push_back(1);
    number_of_shards_++;
  } else {
    // We iterate through all existing shards, if all are full
    // we create a new shard with quantity 1, otherwise we increment
    // the quantity of UEs in the latest shard
    for (auto i = 0; i < number_of_shards_; i++) {
      if (shards_[i] < MAX_SHARD_SIZE) {
        shards_[i]++;
        return i;
      }
    }
    shard_id = number_of_shards_++;
    shards_.push_back(1);
  }
  return shard_id;
}

void ShardTracker::remove_ue(int shard_id) {
  // Since we only keep global state of all UEs, we just
  // need to decrement the number of UEs at a particular id
  shards_[shard_id]--;
}

}  // namespace magma
