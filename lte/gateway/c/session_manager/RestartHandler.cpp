/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#include "RestartHandler.h"
#include "magma_logging.h"

namespace magma {
namespace sessiond {

const uint RestartHandler::max_cleanup_retries_ = 3;

RestartHandler::RestartHandler(
  std::shared_ptr<AsyncDirectorydClient> directoryd_client,
  std::shared_ptr<LocalEnforcer> enforcer,
  SessionReporter* reporter):
  directoryd_client_(directoryd_client),
  enforcer_(enforcer),
  reporter_(reporter)
{
}

void RestartHandler::cleanup_previous_sessions() {
  auto res = directoryd_client_->get_all_directoryd_records(
    [this](Status status, AllDirectoryRecords response) {
      if (!status.ok()) {
        MLOG(MERROR) << "DirectoryD call failed. "
          "Not terminating previous sessions";
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
        uint current_try = 0;
        while (!sessions_to_terminate_.empty() &&
          current_try < max_cleanup_retries_) {
            for (const auto& iter: sessions_to_terminate_) {
              terminate_previous_session(iter.first, iter.second);
            }
            current_try++;
        }
        if (sessions_to_terminate_.empty()) {
          MLOG(MINFO) << "Successfully terminated all old sessions";
        } else {
          MLOG(MERROR) << "Terminating old sessions failed";
        }
      }
    });
}

void RestartHandler::terminate_previous_session(
  const std::string& sid,
  const std::string& session_id)
{
  SessionTerminateRequest term_req;
  term_req.set_sid(sid);
  term_req.set_session_id(session_id);
  (*reporter_).report_terminate_session(
    term_req,
    [this, sid, session_id] (Status status, SessionTerminateResponse response) {
      if (!status.ok()) {
        MLOG(MERROR) << "CCR-T cleanup for subscriber " << sid
          << ", session id: " << session_id << " failed";
        return;
      }
      // Don't delete subscriber from directoryD if IMSI is known
      if (enforcer_->is_imsi_duplicate(response.sid())) {
        MLOG(MINFO) << "Not cleaning up previous session after restart "
          << "for subscriber " << response.sid() << ", session id: "
          << response.session_id() << ": subscriber session exists.";
        return;
      }
      DeleteRecordRequest del_request;
      del_request.set_id(response.sid());
      directoryd_client_->delete_directoryd_record(
        del_request,
        [this, &del_request, sid] (Status status, Void) {
          if (!status.ok()) {
            MLOG(MERROR) << "DirectoryD DeleteRecord failed to remove "
            << "subscriber " << sid << " from DirectoryD";
            return;
          }
          sessions_to_terminate_.erase(del_request.id());
          MLOG(MDEBUG) << "Successfully terminated previous session for "
            << "subscriber " << sid;
        });
    });
}
} // namespace sessiond
} // namespace magma
