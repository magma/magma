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
#include <stdint.h>
#include <experimental/optional>
#include <memory>
#include <set>
#include <string>
#include <unordered_map>

#include "MemoryStoreClient.h"
#include "MeteringReporter.h"
#include "RedisStoreClient.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "StoreClient.h"
#include "StoredState.h"

namespace magma {
class StaticRuleStore;
struct SessionStateUpdateCriteria;

namespace lte {
class MeteringReporter;
class RedisStoreClient;
class UpdateSessionRequest;

using std::experimental::optional;

// Value int represents the request numbers needed for requests to PCRF
using SessionRead = std::set<std::string>;
using SessionUpdate = std::unordered_map<
    std::string, std::unordered_map<std::string, SessionStateUpdateCriteria>>;

enum SessionSearchCriteriaType {
  IMSI_AND_APN = 0,
  IMSI_AND_SESSION_ID = 1,
  IMSI_AND_UE_IPV4 = 2,
  IMSI_AND_UE_IPV4_OR_IPV6 = 3,
  IMSI_AND_BEARER = 4,
  IMSI_AND_TEID = 5,
  IMSI_AND_PDUID = 6,
  IMSI_AND_UE_IPV4_OR_IPV6_OR_UPF_TEID = 7,
};

struct SessionSearchCriteria {
  std::string imsi;
  SessionSearchCriteriaType search_type;
  std::string secondary_key;
  uint32_t secondary_key_unit32;
  uint32_t tertiary_key_unit32;

  SessionSearchCriteria(const std::string p_imsi,
                        SessionSearchCriteriaType p_type,
                        const std::string p_secondary_key)
      : imsi(p_imsi), search_type(p_type), secondary_key(p_secondary_key) {}

  SessionSearchCriteria(const std::string p_imsi,
                        SessionSearchCriteriaType p_type,
                        const uint32_t secondary_key_unit32)
      : imsi(p_imsi),
        search_type(p_type),
        secondary_key_unit32(secondary_key_unit32) {}

  SessionSearchCriteria(const std::string p_imsi,
                        SessionSearchCriteriaType p_type,
                        const std::string p_secondary_key,
                        const uint32_t tertiary_key_unit32)
      : imsi(p_imsi),
        search_type(p_type),
        secondary_key(p_secondary_key),
        tertiary_key_unit32(tertiary_key_unit32) {}
};

/**
 * SessionStore acts as a broker to storage of sessiond state.
 *
 * This allows sessiond to service gRPC requests in a stateless manner.
 * Instead of keeping state in memory, sessiond uses the request parameters and
 * fetches state through SessionStore, handles the request, then writes back
 * to SessionStore, and responds to the gRPC request.
 *
 * SessionStore is intended to be a thread-safe singleton. Each gRPC request
 * should make a single read from SessionStore, and make a single write after
 * the request is serviced. The transactional nature of how requests should be
 * handled is intended to keep sessiond restartable in case of crashes.
 */
class SessionStore {
 public:
  static SessionUpdate get_default_session_update(SessionMap& session_map);

  SessionStore(std::shared_ptr<StaticRuleStore> rule_store,
               std::shared_ptr<magma::MeteringReporter> metering_reporter);

  SessionStore(std::shared_ptr<StaticRuleStore> rule_store,
               std::shared_ptr<magma::MeteringReporter> metering_reporter,
               std::shared_ptr<RedisStoreClient> store_client);

  /**
   * @brief Return a boolean to indicate whether the storage client is ready to
   * accept requests
   */
  bool is_ready() { return store_client_->is_ready(); };

  /**
   * Writes the session map directly to the store. Note that the existing map
   * will be overwriten
   * @param session_map
   * @return
   */
  bool raw_write_sessions(SessionMap session_map);

  /**
   * Read the last written values for the requested sessions through the
   * storage interface.
   * @param req
   * @return Last written values for requested sessions. Returns an empty vector
   * for subscribers that do not have active sessions.
   */
  SessionMap read_sessions(const SessionRead& req);

  /**
   * Read the last written values for all existing sessions through the
   * storage interface.
   * @return Last written values for all sessions. Returns an empty vector
   * for subscribers that do not have active sessions.
   */
  SessionMap read_all_sessions();

  /**
   * Modify the SessionMap in SessionStore to match the current state in
   * the callback.
   * NOTE: Call this method before reporting to other services.
   * NOTE: To avoid race conditions, call this method immediately after
   *       incrementing request numbers and returning control back to the
   *       event loop.
   * @param update_criteria
   */
  void sync_request_numbers(const SessionUpdate& update_criteria);

  /**
   * Goes over all the RG keys and monitoring keys on the UpdateSessionRequest
   * object, and updates is_reporting flab with the value. This function it is
   * used to mark a specific key is currently waiting to get an answer back
   * from the core
   * @param value
   * @param update_session_request
   * @param session_uc
   */
  void set_and_save_reporting_flag(
      bool value, const UpdateSessionRequest& update_session_request,
      SessionUpdate& session_uc);

  /**
   * Read the last written values for the requested sessions through the
   * storage interface. This also modifies the request_numbers stored before
   * returning the SessionMap to the caller, incremented by one for each
   * session.
   * NOTE: It is assumed that the correct number of request_numbers are
   *       reserved on each read_sessions call. If more requests are made to
   *       the OCS/PCRF than are requested, this can cause undefined behavior.
   * NOTE: Here, it is expected that the caller will use one additional
   *       request_number for each session.
   * @param req
   * @return Last written values for requested sessions. Returns an empty
   * vector for subscribers that do not have active sessions.
   */
  SessionMap read_sessions_for_deletion(const SessionRead& req);

  /**
   * Create sessions for a subscriber. Redundant creations will fail.
   * @param subscriber_id
   * @param sessions
   * @return true if successful, otherwise the update to storage is discarded.
   */
  bool create_sessions(const std::string& subscriber_id,
                       SessionVector sessions);

  /**
   * Attempt to update sessions with update criteria. If any update to any of
   * the sessions is invalid, the whole update request is assumed to be invalid,
   * and nothing in storage will be overwritten.
   * NOTE: Will not update request_number. Use sync_request_numbers.
   * @param update_criteria
   * @return true if successful, otherwise the update to storage is discarded.
   */
  bool update_sessions(const SessionUpdate& update_criteria);

  /**
   * @param session_map
   * @param id
   * @return If the session that meets the criteria is found, then it returns an
   * optional of the iterator. Otherwise, it returns an empty value.
   *
   * Usage Example
   * SessionSearchCriteria criteria(IMSI1, IMSI_AND_SESSION_ID,
   * SESSION_ID_1);
   * auto session_it = session_store_.find_session(session_map,
   * id);
   * if (!session_it) { // Log session not found };
   * auto& session = **session_it; // First deference optional, then iterator
   */
  optional<SessionVector::iterator> find_session(
      SessionMap& session_map, SessionSearchCriteria criteria);

  // TODO move this logic outside of this class into MeteringReporter
  /**
   * This function loops through all sessions and propagates the total usage to
   * metering_reporter
   */
  void initialize_metering_counter();

 private:
  std::shared_ptr<StaticRuleStore> rule_store_;
  std::shared_ptr<StoreClient> store_client_;
  std::shared_ptr<MeteringReporter> metering_reporter_;
};

}  // namespace lte
}  // namespace magma
