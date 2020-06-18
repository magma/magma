/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "AAAClient.h"
#include "SessionState.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

using grpc::Status;

namespace { // anonymous

aaa::terminate_session_request create_deactivate_req(
  const std::string &radius_session_id,
  const std::string &imsi)
{
  aaa::terminate_session_request req;
  req.set_radius_session_id(radius_session_id);
  req.set_imsi(imsi);
  return req;
}

aaa::add_sessions_request create_add_sessions_req(
  magma::lte::SessionMap &session_map)
{
  aaa::add_sessions_request req;
  for (auto it = session_map.begin(); it != session_map.end(); it++) {
    for (const auto &session : it->second) {
      aaa::context ctx;
      if (!session->is_radius_cwf_session()) {
        continue;
      }
      auto config = session->get_config();
      magma::SessionState::SessionInfo session_info;
      session->get_session_info(session_info);
      ctx.set_imsi(session_info.imsi);
      ctx.set_session_id(config.radius_session_id);
      ctx.set_acct_session_id(session->get_session_id());
      ctx.set_mac_addr(config.mac_addr);
      ctx.set_msisdn(config.msisdn);
      ctx.set_apn(config.apn);
      auto mutable_sessions = req.mutable_sessions();
      mutable_sessions->Add()->CopyFrom(ctx);
    }
  }
  return req;
}

} // namespace


namespace aaa {

AsyncAAAClient::AsyncAAAClient(
  std::shared_ptr<grpc::Channel> channel):
  stub_(accounting::NewStub(channel))
{
}

AsyncAAAClient::AsyncAAAClient():
  AsyncAAAClient(magma::ServiceRegistrySingleton::Instance()->GetGrpcChannel(
    "aaa_server",
    magma::ServiceRegistrySingleton::LOCAL))
{
}

bool AsyncAAAClient::terminate_session(
  const std::string &radius_session_id,
  const std::string &imsi)
{
  auto req = create_deactivate_req(radius_session_id, imsi);
  terminate_session_rpc(req, [radius_session_id, imsi](
    Status status, acct_resp resp) {
    if (status.ok()) {
      MLOG(MDEBUG) << "Terminated session for Radius ID:"
                   << radius_session_id << ", IMSI: " << imsi;
    } else {
      MLOG(MERROR) << "Could not add terminate session. Radius ID:"
                   << radius_session_id << ", IMSI: " << imsi
                   << ", Error: " << status.error_message();
    }
  });
  return true;
}

bool AsyncAAAClient::add_sessions(magma::lte::SessionMap &session_map)
{
  auto req = create_add_sessions_req(session_map);
  if (req.sessions().size() == 0) {
    MLOG(MINFO) << "Not sending add_sessions request to AAA server. No AAA "
                << "sessions found";
    return true;
  }
  add_sessions_rpc(req, [this](Status status, acct_resp resp) {
    if (status.ok()) {
      MLOG(MINFO) << "Successfully added all existing sessions to AAA server";
    } else {
      MLOG(MERROR) << "Could not add existing sessions to AAA server,"
                   << " Error: " << status.error_message();
    }
  });
  return true;
}

void AsyncAAAClient::add_sessions_rpc(
  const add_sessions_request &request,
  std::function<void(Status, acct_resp)> callback)
{
  auto local_resp = new magma::AsyncLocalResponse<acct_resp>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->Asyncadd_sessions(
    local_resp->get_context(), request, &queue_)));

}

void AsyncAAAClient::terminate_session_rpc(
  const terminate_session_request &request,
  std::function<void(Status, acct_resp)> callback)
{
  auto local_resp = new magma::AsyncLocalResponse<acct_resp>(
    std::move(callback), RESPONSE_TIMEOUT);
  local_resp->set_response_reader(std::move(stub_->Asyncterminate_session(
    local_resp->get_context(), request, &queue_)));
}

} // namespace aaa
