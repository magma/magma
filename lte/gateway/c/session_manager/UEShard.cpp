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

int UEShard::add_ue() {
  int shard_id;
  if (shards.size() == 0) {
    shard_id = 0;
    shards.push_back(1);
    number_of_shards++;
  } else {
    for (auto i = 0; i < number_of_shards; i++) {
      if (shards[i] < MAX_SHARD_SIZE) {
        shards[i]++;
        return i;
      }
    }
    shard_id = number_of_shards++;
    shards.push_back(1);
  }
  return shard_id;
}

void UEShard::remove_ue(int shard_id) {
  shards[shard_id]--;
}

}  // namespace magma
