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
#include <cxxabi.h>
#include <glog/logging.h>
#include <google/protobuf/stubs/common.h>
#include <grpcpp/impl/codegen/status.h>
#include <chrono>
#include <future>
#include <memory>
#include <ostream>
#include <string>
#include <system_error>
#include <thread>
#include <unordered_map>
#include <utility>
#include <vector>

#include "AAAClient.h"
#include "DirectorydClient.h"
#include "RestartHandler.h"
#include "SessionReporter.h"
#include "SessionStore.h"
#include "lte/protos/session_manager.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "magma_logging.h"
#include "orc8r/protos/directoryd.pb.h"

namespace magma {
class LocalEnforcer;
namespace orc8r {
class Void;
}  // namespace orc8r

namespace sessiond {

const uint RestartHandler::max_cleanup_retries_ = 3;
const uint RestartHandler::rpc_retry_interval_s_ = 5;

RestartHandler::RestartHandler(
    std::shared_ptr<DirectorydClient> directoryd_client,
    std::shared_ptr<aaa::AsyncAAAClient> aaa_client,
    std::shared_ptr<LocalEnforcer> enforcer, SessionReporter* reporter,
    SessionStore& session_store)
    : directoryd_client_(directoryd_client),
      aaa_client_(aaa_client),
      enforcer_(enforcer),
      reporter_(reporter),
      session_store_(session_store) {}

void RestartHandler::cleanup_previous_sessions() {
  bool populated = populate_sessions_to_terminate_with_retries();
  if (!populated) {
    MLOG(MERROR) << "DirectoryD call to fetch all records failed after "
                 << max_cleanup_retries_
                 << " retries. Not cleaning up previous sessions";
    return;
  }

  bool all_terminated = launch_threads_to_terminate_with_retries();
  if (all_terminated) {
    MLOG(MINFO) << "Successfully terminated all old sessions";
  } else {
    MLOG(MERROR) << "Terminating old sessions failed";
  }
}

void RestartHandler::setup_aaa_sessions() {
  auto session_map = session_store_.read_all_sessions();
  aaa_client_->add_sessions(session_map);
}

// For each session that still need to be terminated, launch a thread to
// run terminate_previous_session. Retry up to max_cleanup_retries_, with
// sleeps in between.
// We must lock access to sessions_to_terminate while iterating and launching
// threads, so that elements are not removed during iteration.
bool RestartHandler::launch_threads_to_terminate_with_retries() {
  if (sessions_to_terminate_.empty()) {
    return true;
  }

  uint termination_try = 0;
  while (!sessions_to_terminate_.empty() &&
         termination_try < max_cleanup_retries_) {
    std::vector<std::future<void>> termination_futures;
    {
      // Modify sessions_to_terminate_
      std::lock_guard<std::mutex> map_guard(sessions_to_terminate_lock_);

      for (const auto& iter : sessions_to_terminate_) {
        termination_futures.push_back(std::async(
            std::launch::async,
            &magma::sessiond::RestartHandler::terminate_previous_session, this,
            iter.first, iter.second));
      }
      // sessions_to_terminate_lock_ released
    }

    termination_try++;
    // Retrieve all thread results
    for (auto&& fut : termination_futures) {
      fut.get();
    }
    termination_futures.clear();

    // Check whether all sessions have successfully terminated before sleeping
    // There are no active terminate threads, so no locking is necessary.
    if (sessions_to_terminate_.empty()) {
      return true;
    }
    std::this_thread::sleep_for(std::chrono::seconds(rpc_retry_interval_s_));
  }
  // At this point, we've re-tried the maximum number of times with no success.
  return false;
}

// This function is executed in the main thread before multiple terminate
// threads are launched. So no locking is necessary here.
bool RestartHandler::populate_sessions_to_terminate_with_retries() {
  uint rpc_try = 0;
  bool finished = false;
  std::promise<bool> directoryd_res;

  while (rpc_try < max_cleanup_retries_) {
    std::future<bool> directoryd_future = directoryd_res.get_future();
    directoryd_client_->get_all_directoryd_records(
        [this, &directoryd_res](Status status,
                                const AllDirectoryRecords& response) {
          if (!status.ok()) {
            directoryd_res.set_value(false);
            return;
          }
          for (auto& record : response.records()) {
            auto session_iter = record.fields().find("session_id");
            if (session_iter == record.fields().end()) {
              continue;
            }
            const std::string& session_id = session_iter->second;
            sessions_to_terminate_.insert({record.id(), session_id});
          }
          directoryd_res.set_value(true);
        });

    // Block until DirectoryD call is complete
    finished = directoryd_future.get();
    if (finished) {
      break;
    }
    // Setup for next iteration
    rpc_try++;
    directoryd_res = std::promise<bool>();
    std::this_thread::sleep_for(std::chrono::seconds(rpc_retry_interval_s_));
  }
  return finished;
}

void RestartHandler::terminate_previous_session(const std::string& sid,
                                                const std::string& session_id) {
  SessionTerminateRequest term_req;
  term_req.mutable_common_context()->mutable_sid()->set_id(sid);
  term_req.set_session_id(session_id);
  std::promise<bool> termination_res;
  std::future<bool> termination_future = termination_res.get_future();
  (*reporter_)
      .report_terminate_session(
          term_req,
          [this, &termination_res, sid, session_id](
              Status status, const SessionTerminateResponse& response) {
            if (!status.ok()) {
              MLOG(MERROR) << "CCR-T cleanup for subscriber " << sid
                           << ", session id: " << session_id << " failed";
              termination_res.set_value(false);
              return;
            }
            DeleteRecordRequest del_request;
            del_request.set_id(response.sid());
            directoryd_client_->delete_directoryd_record(
                del_request,
                [&termination_res, sid](Status status, const Void&) {
                  if (!status.ok()) {
                    MLOG(MERROR) << "DirectoryD DeleteRecord failed to remove "
                                 << "subscriber " << sid << " from DirectoryD";
                    termination_res.set_value(false);
                    return;
                  }
                  MLOG(MDEBUG)
                      << "Successfully terminated previous session for "
                      << "subscriber " << sid;
                  termination_res.set_value(true);
                });
          });
  // Block until Termination call is complete
  bool should_erase = termination_future.get();
  if (should_erase) {
    std::lock_guard<std::mutex> map_guard(sessions_to_terminate_lock_);
    sessions_to_terminate_.erase(sid);
    // sessions_to_terminate_lock_ released
  }
}
}  // namespace sessiond
}  // namespace magma
