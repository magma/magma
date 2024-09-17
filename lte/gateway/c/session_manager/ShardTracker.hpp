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
#include <stdint.h>
#include <set>
#include <string>
#include <vector>

namespace magma {

/* Shards represent groups of UEs placed into buckets of
a certain size, to make polling more manageable*/
class ShardTracker {
 public:
  ShardTracker();
  /**
   * Add UE to shards based on availability, if the UE already has an
   * existing shard, return the existing shard id and don't perform
   * an addition
   * @return index(shard id) where UE was placed
   */
  uint16_t add_ue(const std::string imsi);

  /**
   * Remove UE from shard
   * @param shard_id location of UE to be removed
   * @return true for successful removal, false for failed removal
   */
  bool remove_ue(const std::string imsi, const uint16_t shard_id);

 private:
  /*
   * a vector of quantities, where the indices represent
   * the shard id and the values are the UEs(IMSIs) assigned
   * to that shard id
   */
  std::vector<std::set<std::string>> imsis_per_shard_;
  /*
   * largest number of UEs that can fill a shard
   */
  const uint16_t max_shard_size_ = 100;
};

}  // namespace magma
