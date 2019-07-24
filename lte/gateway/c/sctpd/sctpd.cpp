/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sctpd.h"

#include <memory>

#include <grpcpp/grpcpp.h>

#include "sctpd_downlink_impl.h"
#include "sctpd_event_handler.h"
#include "sctpd_uplink_client.h"
#include "util.h"

using grpc::ServerBuilder;
using magma::sctpd::SctpdDownlinkImpl;
using magma::sctpd::SctpdEventHandler;
using magma::sctpd::SctpdUplinkClient;

int main()
{
  magma::init_logging("sctpd");
  magma::set_verbosity(MDEBUG);

  auto channel =
    grpc::CreateChannel(UPSTREAM_SOCK, grpc::InsecureChannelCredentials());

  SctpdUplinkClient client(channel);
  SctpdEventHandler handler(client);
  SctpdDownlinkImpl service(handler);

  ServerBuilder builder;
  builder.AddListeningPort(DOWNSTREAM_SOCK, grpc::InsecureServerCredentials());
  builder.RegisterService(&service);

  auto server = builder.BuildAndStart();

  MLOG(MINFO) << "sctp downlink server started, waiting for init";

  server->Wait();

  return 0;
}
