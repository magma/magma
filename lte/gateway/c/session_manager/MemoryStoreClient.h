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

#include <lte/protos/session_manager.grpc.pb.h>
#include <memory>
#include <set>
#include <string>
#include <unordered_map>
#include <vector>

#include "StoreClient.h"
#include "StoredState.h"

namespace magma {
class StaticRuleStore;
struct StoredSessionState;

namespace lte {

/**
 * Non-persistent StoreClient used as intermediate stage in development
 */
class MemoryStoreClient final : public StoreClient {
 public:
  MemoryStoreClient(std::shared_ptr<StaticRuleStore> rule_store);
  MemoryStoreClient(MemoryStoreClient const&) = delete;
  MemoryStoreClient(MemoryStoreClient&&) = default;
  ~MemoryStoreClient() = default;

  bool is_ready() { return true; }

  SessionMap read_sessions(std::set<std::string> subscriber_ids);

  SessionMap read_all_sessions();

  bool write_sessions(SessionMap session_map);

 private:
  std::unordered_map<std::string, std::vector<StoredSessionState>> session_map_;
  std::shared_ptr<StaticRuleStore> rule_store_;
};

}  // namespace lte
}  // namespace magma
