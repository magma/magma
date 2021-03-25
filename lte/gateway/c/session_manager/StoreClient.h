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

#include <memory>

#include <lte/protos/session_manager.grpc.pb.h>

#include "SessionState.h"

namespace magma {
namespace lte {

typedef std::vector<std::unique_ptr<SessionState>> SessionVector;
typedef std::unordered_map<std::string, SessionVector> SessionMap;

/**
 * StoreClient is responsible for reading/writing sessions to/from storage.
 */
class StoreClient {
 public:
  virtual ~StoreClient() = default;

  /**
   * @brief Return a boolean to indicate whether the client is ready to accept
   * requests
   */
  virtual bool is_ready() = 0;

  /**
   * Directly read the subscriber's sessions from storage
   *
   * If one or more of the subscribers have no sessions, empty entries will be
   * returned.
   * @param subscriber_ids typically in IMSI
   * @return All sessions for the subscribers
   */
  virtual SessionMap read_sessions(std::set<std::string> subscriber_ids) = 0;

  /**
   * Directly read all subscriber sessions from storage
   *
   * If one or more of the subscribers have no sessions, empty entries will be
   * returned.
   * @return All sessions for the subscribers
   */
  virtual SessionMap read_all_sessions() = 0;

  /**
   * Directly write the subscriber sessions into storage, overwriting previous
   * values.
   *
   * @param sessions Sessions to write into storage
   * @return True if writes have completed successfully for all sessions.
   */
  virtual bool write_sessions(SessionMap sessions) = 0;
};

}  // namespace lte
}  // namespace magma
