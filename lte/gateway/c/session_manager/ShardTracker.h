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

// Shards represent groups of UEs placed into buckets of
// a certain size, to make polling more manageable
class ShardTracker {
  ShardTracker();

  // add UE to shards based on availability
  int add_ue();

  // remove UE from shard
  void remove_ue(int shard_id);

 private:
  std::vector<int> shards;
  int max_shard_size;
  int number_of_shards;
};

}  // namespace magma
