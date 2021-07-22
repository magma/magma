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
#pragma once
#include <string>
#include <vector>

namespace magma {

/* Shards represent groups of UEs placed into buckets of
a certain size, to make polling more manageable*/
class ShardTracker {
  ShardTracker();

  /**
   * Add UE to shards based on availability
   * @return index(shard id) where UE was placed
   * TODO(veshkemburu): Store IMSI as well for easier subscriber reallocation
   * (GH8167)
   */
  int add_ue();

  /**
   * Remove UE from shard
   * @param shard_id location of UE to be removed
   * @return true for successful removal, false for failed removal
   */
  bool remove_ue(int shard_id);

 private:
  /*
   * a vector of quantities, where the indices represent
   * the shard id and the values represent the number of
   * UEs held in each shard
   */
  std::vector<uint16_t> ue_count_per_shard_;
  /*
   * largest number of UEs that can fill a shard
   */
  const uint16_t max_shard_size_ = 100;
};

}  // namespace magma
