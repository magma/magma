// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <folly/futures/Future.h>

#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

class Channel;

class Request final {
 public:
  Request(Channel* channel_);
  Request() = delete;
  ~Request() = default;
  Request(const Request&) = delete;
  Request& operator=(const Request&) = delete;
  Request(Request&&) = default;
  Request& operator=(Request&&) = delete;

 public:
  Channel* channel{nullptr};
  folly::Promise<Response> responsePromise{};
};

} // namespace snmp
} // namespace channels
} // namespace devmand
