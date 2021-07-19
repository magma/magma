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
#include <map>

namespace magma {

/*Shards represent groups of UEs placed into buckets of
a certain size, to make polling more manageable*/
class ShardTracker {
  ShardTracker();

  /**
   * add UE to shards based on availability
   * @return index(shard id) where UE was placed
   */
  // TODO: Store IMSI as well for easier subscriber reallocation
  int add_ue();

  /**
   * remove UE from shard
   * @param shard_id location of UE to be removed
   */
  void remove_ue(int shard_id);

 /* shards: a vector of quantities, where the indices represent
             the shard id and the values represent the number of
             UEs held in each shard
     max_shard_size: represents the largest number of UEs that can
                     fill a shard
     number_of_shards: the number of shards we currently have
  */
 private:
  std::vector<uint16_t> shards_;
  uint16_t max_shard_size_;
  uint16_t number_of_shards_;
};

}  // namespace magma
