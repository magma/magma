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

class UEShard {
 private:
  std::map<int, std::vector<std::string>> shards;
  int number_of_shards;
  int max_shard_size;

  UEShard();

  // returns shard_id associated with UE added
  int add_ue(std::string imsi);

  std::pair<int, int> find_ue_shard(std::string imsi);

  void remove_ue(std::string imsi);

  // change shard for ue
  int move_ue(std::string imsi, int shard_id);

  int total_ues_for_shard(int shard_id);
};

}  // namespace magma
