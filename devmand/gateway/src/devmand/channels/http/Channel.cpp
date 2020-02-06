// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <folly/executors/Async.h>
#include <iostream>

#include <devmand/channels/http/Channel.h>
#include <devmand/error/ErrorHandler.h>

namespace devmand {
namespace channels {
namespace http {

Channel::Channel(const std::string& controllerHost_, const int controllerPort)
    : controllerHost(controllerHost_),
      controller(
          controllerHost,
          controllerPort
          // TODO convert the above to something configurable
      ) {
  controller.set_timeout_sec(20000); /* timeout in seconds */
}

folly::Future<Response> Channel::asyncMsg(
    const std::string& endpoint,
    std::function<std::shared_ptr<httplib::Response>()> send) {
  auto requestId = requestGuid++;
  auto result = outstandingRequests.emplace(
      std::piecewise_construct,
      std::forward_as_tuple(requestId),
      std::forward_as_tuple(folly::Promise<Response>{}));
  if (result.second) {
    std::thread t([this, send, requestId, endpoint]() {
      std::thread::id id = std::this_thread::get_id();
      std::cerr << "started " << id << std::endl;
      ErrorHandler::executeWithCatch(
          [this, &send, &id, requestId, &endpoint]() {
            auto res = send();
            auto pResult = outstandingRequests.find(requestId);
            if (pResult == outstandingRequests.end()) {
              throw std::runtime_error("unknown outstanding request");
            } else if (
                res != nullptr and res->status >= 200 and res->status <= 299) {
              std::cerr << id << ": response is " << res->body << std::endl;
              pResult->second.setValue(Response(res->body));
            } else if (res != nullptr) {
              std::cerr << id
                        << folly::sformat(
                               ": http error {} on {}", res->status, endpoint)
                        << std::endl;
              pResult->second.setValue(ErrorResponse(folly::sformat(
                  "http error {} on {}", res->status, endpoint)));
            } else {
              std::cerr << id << ": http error on " << endpoint << std::endl;
              pResult->second.setValue(
                  ErrorResponse(folly::sformat("http error on {}", endpoint)));
            }
            outstandingRequests.erase(pResult);
          });
      std::cerr << "ended " << id << std::endl;
    });
    t.detach(); // TODO Yeah this isn't nice but its a hack for now. Will
                // switch to a better library later.
    return result.first->second.getFuture();
  } else {
    throw std::runtime_error("emplace failed");
  }
}

folly::Future<Response> Channel::asyncPut(
    const httplib::Headers& headers,
    const std::string& endpoint,
    const std::string& body,
    const std::string& contentType) {
  return asyncMsg(endpoint, [this, endpoint, headers, body, contentType]() {
    return controller.Put(
        endpoint.c_str(), headers, body.c_str(), contentType.c_str());
  });
}

folly::Future<Response> Channel::asyncGet(
    const httplib::Headers& headers,
    const std::string& endpoint) {
  return asyncMsg(endpoint, [this, endpoint, headers]() {
    return controller.Get(endpoint.c_str(), headers);
  });
}

} // namespace http
} // namespace channels
} // namespace devmand
