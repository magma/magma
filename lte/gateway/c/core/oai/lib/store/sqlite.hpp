/*
 * Copyright 2022 The Magma Authors.
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
#include "lte/protos/subscriberdb.pb.h"

using google::protobuf::Message;

namespace magma {
namespace lte {
class SqliteStore {
 public:
  SqliteStore(std::string db_location, int sid_digits);

  // Initialize data store
  void init_db_connection(std::string db_location, int sid_digits = 2);

  // Add subscriber
  void add_subscriber(const SubscriberData& subscriber_data);

  // Delete subscriber
  void delete_subscriber();  // TODO(vroon2703): add the parameters

 private:
  int _sid_digits;
  int _n_shards;
  std::vector<std::string> _db_locations;
  std::vector<std::string> _create_db_locations(std::string db_location,
                                                int _n_shards);
  void _create_store();
  const char* _get_sid(const SubscriberData& subscriber_data);
  // Map subscriber ID to bucket
  int _sid2bucket(std::string sid);
};
}  // namespace lte
}  // namespace magma
