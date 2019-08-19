/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#define UPSTREAM_SOCK "unix:///tmp/sctpd_upstream.sock"
#define DOWNSTREAM_SOCK "unix:///tmp/sctpd_downstream.sock"

#define SCTP_OUT_STREAMS (16)
#define SCTP_IN_STREAMS (16)
#define SCTP_MAX_ATTEMPTS (2)
#define SCTP_TIMEOUT (5)
#define SCTP_RECV_BUFFER_SIZE (4096)
