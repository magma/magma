/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <memory>
#include <string>
#include <utility>

#include <folly/Format.h>
#include <folly/dynamic.h>
#include <folly/json.h>
#include <orc8r/protos/eventd.pb.h>

#include "EventdClient.h"
#include "SessionState.h"
#include "magma_logging.h"

namespace magma {
/**
 * Session Events are sent to the eventd service for logging.
 */
namespace session_events {

void session_created(
    AsyncEventdClient& client,
    const std::string& imsi,
    const std::string& session_id);

void session_terminated(
    AsyncEventdClient& client,
    const std::unique_ptr<SessionState>& session);

}  // namespace session_events
}  // namespace magma
