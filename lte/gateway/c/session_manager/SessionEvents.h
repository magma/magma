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

#include <folly/Format.h>
#include <folly/dynamic.h>
#include <folly/json.h>
#include <orc8r/protos/eventd.pb.h>
#include <memory>
#include <string>
#include <utility>

#include "SessionState.h"
#include "includes/EventdClient.h"
#include "magma_logging.h"

namespace magma {
class EventdClient;
class SessionState;
struct SessionConfig;
struct UpdateRequests;

namespace lte {

class EventsReporter {
 public:
  virtual ~EventsReporter() = default;

  virtual void session_created(
      const std::string& imsi, const std::string& session_id,
      const SessionConfig& session_context,
      const std::unique_ptr<SessionState>& session) = 0;

  virtual void session_create_failure(const SessionConfig& session_context,
                                      const std::string& failure_reason) = 0;

  virtual void session_updated(const std::string& session_id,
                               const SessionConfig& session_context,
                               const UpdateRequests& update_request) = 0;

  virtual void session_update_failure(const std::string& session_id,
                                      const SessionConfig& session_context,
                                      const UpdateRequests& failed_request,
                                      const std::string& failure_reason) = 0;

  virtual void session_terminated(
      const std::string& imsi,
      const std::unique_ptr<SessionState>& session) = 0;
};

/**
 * Session Events are sent to the eventd service for logging.
 */
class EventsReporterImpl : public EventsReporter {
 public:
  explicit EventsReporterImpl(EventdClient& eventd_client);

  void session_created(const std::string& imsi, const std::string& session_id,
                       const SessionConfig& session_context,
                       const std::unique_ptr<SessionState>& session);

  void session_create_failure(const SessionConfig& session_context,
                              const std::string& failure_reason);

  void session_updated(const std::string& session_id,
                       const SessionConfig& session_context,
                       const UpdateRequests& update_request);

  void session_update_failure(const std::string& session_id,
                              const SessionConfig& session_context,
                              const UpdateRequests& failed_request,
                              const std::string& failure_reason);

  void session_terminated(const std::string& imsi,
                          const std::unique_ptr<SessionState>& session);

 private:
  std::string get_mac_addr(const SessionConfig& config);
  std::string get_imei(const SessionConfig& config);
  std::string get_spgw_ipv4(const SessionConfig& config);
  std::string get_user_location(const SessionConfig& config);
  std::string get_charging_characteristics(const SessionConfig& config);
  folly::dynamic get_update_summary(const UpdateRequests& updates);

 private:
  EventdClient& eventd_client_;
};

}  // namespace lte
}  // namespace magma
