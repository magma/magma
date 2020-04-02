/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include <chrono>
#include <future>

#include "RestartHandler.h"
#include "magma_logging.h"

namespace magma {
namespace sessiond {

const uint RestartHandler::max_cleanup_retries_ = 3;
const uint RestartHandler::rpc_retry_interval_s_ = 5;

RestartHandler::RestartHandler(
  std::shared_ptr<AsyncDirectorydClient> directoryd_client,
  std::shared_ptr<LocalEnforcer> enforcer,
  SessionReporter* reporter,
  SessionMap& session_map):
  directoryd_client_(directoryd_client),
  enforcer_(enforcer),
  reporter_(reporter),
  session_map_(session_map)
{
}

void RestartHandler::cleanup_previous_sessions()
{
  uint rpc_try = 0;
  bool finished = false;
  while (!finished && rpc_try < max_cleanup_retries_) {
    std::promise<bool> directoryd_res;
    std::future<bool> directoryd_future = directoryd_res.get_future();
    auto res = directoryd_client_->get_all_directoryd_records(
      [this, &directoryd_res](Status status, AllDirectoryRecords response) {
        if (!status.ok()) {
          directoryd_res.set_value(false);
          return;
        }
        auto records = response.records();
        for (auto& record : records) {
          auto session_iter = record.fields().find("session_id");
          if (session_iter == record.fields().end()) {
            continue;
          }
          auto session_id = session_iter->second;
          sessions_to_terminate_.insert({record.id(), session_id});
        }
        directoryd_res.set_value(true);
      });
    finished = directoryd_future.get();
    rpc_try++;
    directoryd_res = std::promise<bool>();
    if (!finished) {
      std::this_thread::sleep_for(std::chrono::seconds(rpc_retry_interval_s_));
    }
  }
  if (!finished) {
    MLOG(MERROR) << "DirectoryD call to fetch all records failed after "
                 << max_cleanup_retries_
                 << " retries. Not cleaning up previous sessions";
    return;
  }
  uint termination_try = 0;
  while (!sessions_to_terminate_.empty() &&
         termination_try < max_cleanup_retries_) {
    std::vector<std::future<void>> termination_futures;
    for (const auto& iter : sessions_to_terminate_) {
      termination_futures.push_back(std::async(
        std::launch::async,
        &magma::sessiond::RestartHandler::terminate_previous_session,
        this,
        iter.first,
        iter.second));
    }
    termination_try++;
    for (auto&& fut : termination_futures) {
      fut.get();
    }
    termination_futures.clear();
    std::this_thread::sleep_for(std::chrono::seconds(rpc_retry_interval_s_));
  }
  if (sessions_to_terminate_.empty()) {
    MLOG(MINFO) << "Successfully terminated all old sessions";
  } else {
    MLOG(MERROR) << "Terminating old sessions failed";
  }
}

void RestartHandler::terminate_previous_session(
  const std::string& sid,
  const std::string& session_id)
{
  SessionTerminateRequest term_req;
  term_req.set_sid(sid);
  term_req.set_session_id(session_id);
  std::promise<bool> termination_res;
  std::future<bool> termination_future = termination_res.get_future();
  (*reporter_)
    .report_terminate_session(
      term_req,
      [this, &termination_res, sid, session_id](
        Status status, SessionTerminateResponse response) {
        if (!status.ok()) {
          MLOG(MERROR) << "CCR-T cleanup for subscriber " << sid
                       << ", session id: " << session_id << " failed";
          termination_res.set_value(false);
          return;
        }
        // Don't delete subscriber from directoryD if IMSI is known
        if (enforcer_->session_with_imsi_exists(session_map_, response.sid())) {
          MLOG(MINFO) << "Not cleaning up previous session after restart "
                      << "for subscriber " << response.sid()
                      << ", session id: " << response.session_id()
                      << ": subscriber session exists.";
          termination_res.set_value(true);
          return;
        }
        DeleteRecordRequest del_request;
        del_request.set_id(response.sid());
        directoryd_client_->delete_directoryd_record(
          del_request,
          [this, &del_request, &termination_res, sid](Status status, Void) {
            if (!status.ok()) {
              MLOG(MERROR) << "DirectoryD DeleteRecord failed to remove "
                           << "subscriber " << sid << " from DirectoryD";
              termination_res.set_value(false);
              return;
            }
            MLOG(MDEBUG) << "Successfully terminated previous session for "
                         << "subscriber " << sid;
            termination_res.set_value(true);
          });
      });
  bool should_erase = termination_future.get();
  if (should_erase) {
    sessions_to_terminate_.erase(sid);
  }
}
} // namespace sessiond
} // namespace magma
