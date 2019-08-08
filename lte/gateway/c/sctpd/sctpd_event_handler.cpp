/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sctpd_event_handler.h"

#include <lte/protos/sctpd.grpc.pb.h>

namespace magma {
namespace sctpd {

SctpdEventHandler::SctpdEventHandler(SctpdUplinkClient &client): _client(client)
{
}

void SctpdEventHandler::HandleNewAssoc(
  uint32_t assoc_id,
  uint32_t instreams,
  uint32_t outstreams)
{
  NewAssocReq req;
  NewAssocRes res;

  req.set_assoc_id(assoc_id);
  req.set_instreams(instreams);
  req.set_outstreams(outstreams);

  _client.newAssoc(req, &res);
}

void SctpdEventHandler::HandleCloseAssoc(uint32_t assoc_id, bool reset)
{
  CloseAssocReq req;
  CloseAssocRes res;

  req.set_assoc_id(assoc_id);
  req.set_is_reset(reset);

  _client.closeAssoc(req, &res);
}

void SctpdEventHandler::HandleRecv(
  uint32_t assoc_id,
  uint32_t stream,
  const std::string &payload)
{
  SendUlReq req;
  SendUlRes res;

  req.set_assoc_id(assoc_id);
  req.set_stream(stream);
  req.set_payload(payload);

  _client.sendUl(req, &res);
}

} // namespace sctpd
} // namespace magma
