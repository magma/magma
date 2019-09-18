// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <string>
#include <unordered_map>

#include <folly/futures/Future.h>

// TODO for now I just used http lib which has doesn't intergrate well with an
// epoll and folly dynamics. I want to investigate using proxygen:
//  https://github.com/facebook/proxygen
#include <httplib.h>

#include <devmand/channels/Channel.h>
#include <devmand/channels/http/Response.h>

namespace devmand {
namespace channels {
namespace http {

using OutstandingRequests = std::map<unsigned int, folly::Promise<Response>>;

class Channel final : public channels::Channel {
 public:
  Channel(const std::string& controllerHost, const int controllerPort);
  Channel() = delete;
  ~Channel() override = default;
  Channel(const Channel&) = delete;
  Channel& operator=(const Channel&) = delete;
  Channel(Channel&&) = delete;
  Channel& operator=(Channel&&) = delete;

 public:
  folly::Future<Response> asyncGet(
      const httplib::Headers& headers,
      const std::string& endpoint);
  folly::Future<Response> asyncPut(
      const httplib::Headers& headers,
      const std::string& endpoint,
      const std::string& body,
      const std::string& contentType);

 private:
  folly::Future<Response> asyncMsg(
      const std::string& endpoint,
      std::function<std::shared_ptr<httplib::Response>()> send);

 private:
  std::string controllerHost;
  httplib::Client controller;
  OutstandingRequests outstandingRequests;
  unsigned int requestGuid{0};
};

} // namespace http
} // namespace channels
} // namespace devmand
