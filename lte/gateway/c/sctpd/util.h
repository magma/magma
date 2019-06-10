/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include <arpa/inet.h>
#include <errno.h>
#include <stdint.h>

#include <string>

#include <lte/protos/sctpd.grpc.pb.h>

#include <magma_logging.h>

namespace magma {
namespace sctpd {

#define MLOG_perror(fname)                                                     \
  do {                                                                         \
    MLOG(MERROR) << fname << " error (" << std::to_string(errno)               \
                 << "): " << strerror(errno);                                  \
  } while (0)
#define MLOG_grpcerr(status)                                                   \
  do {                                                                         \
    MLOG(MERROR) << "grpc error (" << std::to_string(status.error_code())      \
                 << "): " << status.error_message();                           \
  } while (0)

int create_sctp_sock(const InitReq &req);

} // namespace sctpd
} // namespace magma
