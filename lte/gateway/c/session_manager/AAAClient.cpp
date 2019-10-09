/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "AAAClient.h"
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
