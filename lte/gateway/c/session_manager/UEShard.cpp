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

#include "UEShard.h"
#define MAX_SHARD_SIZE 100

namespace magma {

UEShard::UEShard() {
  number_of_shards = 0;
}

// returns shard_id associated with UE added
int UEShard::add_ue(std::string imsi) {
  int shard_id;
  if (shards.size() == 0) {
    shard_id = 1;
    shards[shard_id].push_back(imsi);
    number_of_shards++;
  } else {
    // iterate through all shards to find open location to put imsi
    for (auto it = shards.begin(); it != shards.end(); it++) {
      if (total_ues_for_shard(it->first) < MAX_SHARD_SIZE) {
        shard_id = it->first;
        shards[shard_id].push_back(imsi);
        return shard_id;
      }
    }
    shard_id = ++number_of_shards;
    shards[shard_id].push_back(imsi);
  }
  return shard_id;
}

std::pair<int, int> UEShard::find_ue_shard(std::string imsi) {
  std::pair<int, int> location;
  for (auto it = shards.begin(); it != shards.end(); it++) {
    if (std::find(it->second.begin(), it->second.end(), imsi) !=
        it->second.end()) {
      std::vector<std::string>::iterator itr =
          std::find(it->second.begin(), it->second.end(), imsi);
      int indexOfImsi = std::distance(it->second.begin(), itr);
      location        = std::make_pair(it->first, indexOfImsi);
      return location;
    }
  }
  location = std::make_pair(0, 0);
  return location;
}

void UEShard::remove_ue(std::string imsi) {
  std::pair<int, int> location = find_ue_shard(imsi);
  int shard_id                 = location.first;
  int indexOfImsi              = location.second;
  shards[shard_id].erase(shards[shard_id].begin() + indexOfImsi);
}

// change shard for ue
int UEShard::move_ue(std::string imsi, int shard_id) {
  if (shard_id <= 0 || shard_id > number_of_shards) {
    return 0;
  }
  if (total_ues_for_shard(shard_id) == MAX_SHARD_SIZE) {
    return 0;
  }
  remove_ue(imsi);
  shards[shard_id].push_back(imsi);
  return shard_id;
}

int UEShard::total_ues_for_shard(int shard_id) {
  if (shard_id <= 0 || shard_id > number_of_shards) {
    return 0;
  }
  return shards[shard_id].size();
}

}  // namespace magma
