/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sctpd_uplink_client.h"

#include "util.h"

namespace magma {
namespace sctpd {

using grpc::ClientContext;

SctpdUplinkClient::SctpdUplinkClient(std::shared_ptr<Channel> channel)
{
  _stub = SctpdUplink::NewStub(channel);
}

int SctpdUplinkClient::sendUl(const SendUlReq &req, SendUlRes *res)
{
  assert(res != nullptr);

  ClientContext context;

  auto status = _stub->SendUl(&context, req, res);

  if (!status.ok()) {
    MLOG(MERROR) << "sctpul.sendul error";
    MLOG_grpcerr(status);
  }

  return status.ok() ? 0 : -1;
}

int SctpdUplinkClient::newAssoc(const NewAssocReq &req, NewAssocRes *res)
{
  assert(res != nullptr);

  ClientContext context;

  auto status = _stub->NewAssoc(&context, req, res);

  if (!status.ok()) {
    MLOG(MERROR) << "sctpul.newassoc error";
    MLOG_grpcerr(status);
  }

  return status.ok() ? 0 : -1;
}

int SctpdUplinkClient::closeAssoc(const CloseAssocReq &req, CloseAssocRes *res)
{
  assert(res != nullptr);

  ClientContext context;

  auto status = _stub->CloseAssoc(&context, req, res);

  if (!status.ok()) {
    MLOG(MERROR) << "sctpul.closeassoc error";
    MLOG_grpcerr(status);
  }

  return status.ok() ? 0 : -1;
}

} // namespace sctpd
} // namespace magma
