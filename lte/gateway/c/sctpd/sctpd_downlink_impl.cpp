/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sctpd_downlink_impl.h"

#include <arpa/inet.h>
#include <assert.h>
#include <netinet/sctp.h>
#include <unistd.h>
#include "sctpd.h"
#include "util.h"

namespace magma {
namespace sctpd {

SctpdDownlinkImpl::SctpdDownlinkImpl(SctpEventHandler &uplink_handler):
  _uplink_handler(uplink_handler),
  _sctp_connection(nullptr)
{
}

Status SctpdDownlinkImpl::Init(
  ServerContext *context,
  const InitReq *req,
  InitRes *res)
{
  MLOG(MDEBUG) << "SctpdDownlinkImpl::Init starting";

  if (_sctp_connection != nullptr && !req->force_restart()) {
    MLOG(MINFO) << "SctpdDownlinkImpl::Init reusing existing connection";
    res->set_result(InitRes::INIT_OK);
    return Status::OK;
  }

  if (_sctp_connection != nullptr) {
    MLOG(MDEBUG)
      << "SctpdDownlinkImpl::Init cleaning up sctp_desc and listener";

    auto conn = std::move(_sctp_connection);
    conn->Close();
  }

  MLOG(MDEBUG) << "SctpdDownlinkImpl::Init creating new socket and listener";

  try {
    _sctp_connection = std::make_unique<SctpConnection>(*req, _uplink_handler);
  } catch (...) {
    res->set_result(InitRes::INIT_FAIL);
    return Status::OK;
  }

  _sctp_connection->Start();

  res->set_result(InitRes::INIT_OK);
  return Status::OK;
}

Status SctpdDownlinkImpl::SendDl(
  ServerContext *context,
  const SendDlReq *req,
  SendDlRes *res)
{
  MLOG(MDEBUG) << "SctpdDownlinkImpl::SendDl starting";

  try {
    _sctp_connection->Send(req->assoc_id(), req->stream(), req->payload());
  } catch (...) {
    res->set_result(SendDlRes::SEND_DL_FAIL);
    return Status::OK;
  }

  res->set_result(SendDlRes::SEND_DL_OK);
  return Status::OK;
}

void SctpdDownlinkImpl::stop()
{
  if (_sctp_connection != nullptr) {
    _sctp_connection->Close();
  }
}

} // namespace sctpd
} // namespace magma
