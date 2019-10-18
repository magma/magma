// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/futures/Future.h>

#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

class Channel;

class Request final {
 public:
  Request(Channel* channel_, Oid oid_);
  Request() = delete;
  ~Request() = default;
  Request(const Request&) = delete;
  Request& operator=(const Request&) = delete;
  Request(Request&&) = default;
  Request& operator=(Request&&) = delete;

 public:
  Channel* channel{nullptr};
  Oid oid;
  folly::Promise<Response> responsePromise{};
};

} // namespace snmp
} // namespace channels
} // namespace devmand
