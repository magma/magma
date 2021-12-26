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

#include <sys/types.h>
#include <future>
#include <memory>
#include <mutex>
#include <string>
#include <unordered_map>

#include "AAAClient.h"
#include "DirectorydClient.h"
#include "LocalEnforcer.h"
#include "SessionReporter.h"

namespace aaa {
class AsyncAAAClient;
}  // namespace aaa

namespace magma {
class DirectorydClient;
class LocalEnforcer;
class SessionReporter;
namespace lte {
class SessionStore;
}  // namespace lte

namespace sessiond {

/**
 * Restart handler cleans up previous sessions after a SessionD service restart
 */
class RestartHandler {
 public:
  RestartHandler(std::shared_ptr<DirectorydClient> directoryd_client,
                 std::shared_ptr<aaa::AsyncAAAClient> aaa_client,
                 std::shared_ptr<LocalEnforcer> enforcer,
                 SessionReporter* reporter, SessionStore& session_store);

  /**
   * Cleanup previous sessions stored in directoryD
   */
  void cleanup_previous_sessions();

  /**
   * Re-create AAA sessions stored in sessiond
   */
  void setup_aaa_sessions();

 private:
  void terminate_previous_session(const std::string& sid,
                                  const std::string& session_id);
  bool populate_sessions_to_terminate_with_retries();
  bool launch_threads_to_terminate_with_retries();

 private:
  std::shared_ptr<DirectorydClient> directoryd_client_;
  std::shared_ptr<aaa::AsyncAAAClient> aaa_client_;
  std::shared_ptr<LocalEnforcer> enforcer_;
  SessionReporter* reporter_;
  SessionStore& session_store_;
  std::mutex sessions_to_terminate_lock_;  // mutex to guard add/remove access
                                           // to sessions_to_terminate
  std::unordered_map<std::string, std::string> sessions_to_terminate_;
  static const uint max_cleanup_retries_;
  static const uint rpc_retry_interval_s_;
};
}  // namespace sessiond
}  // namespace magma
