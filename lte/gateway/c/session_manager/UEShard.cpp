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

int UEShard::add_ue(const std::string& imsi) {
  int shard_id;
  if (shards.size() == 0) {
    shard_id = 0;
    shards.push_back({imsi});
    number_of_shards++;
  } else {
    for (auto i = 0; i < number_of_shards; i++) {
      if (total_ues_for_shard(i) < max_shard_size) {
        shards[i].push_back(imsi);
        return i;
      }
    }
    shard_id = number_of_shards++;
    shards.push_back({imsi});
  }
  return shard_id;
}

std::pair<int, int> UEShard::find_ue_shard(const std::string& imsi) {
  std::pair<int, int> location;
  for (size_t i = 0; i < shards.size(); i++) {
    if (std::find(shards[i].begin(), shards[i].end(), imsi) !=
        shards[i].end()) {
      std::vector<std::string>::iterator itr =
          std::find(shards[i].begin(), shards[i].end(), imsi);
      int indexOfImsi = std::distance(shards[i].begin(), itr);
      location        = std::make_pair(i, indexOfImsi);
      return location;
    }
  }
  location = std::make_pair(0, 0);
  return location;
}

void UEShard::remove_ue(const std::string& imsi) {
  std::pair<int, int> location = find_ue_shard(imsi);
  int shard_id                 = location.first;
  int indexOfImsi              = location.second;
  shards[shard_id].erase(shards[shard_id].begin() + indexOfImsi);
}

int UEShard::total_ues_for_shard(int shard_id) {
  if (shard_id <= 0 || shard_id > number_of_shards) {
    return 0;
  }
  return shards[shard_id].size();
}

}  // namespace magma
