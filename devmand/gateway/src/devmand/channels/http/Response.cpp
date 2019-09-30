// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/http/Response.h>

namespace devmand {
namespace channels {
namespace http {

Response::Response(const std::string& msg_) : Response(msg_, false) {}
std::string Response::get() const {
  return msg;
}

bool Response::isError() const {
  return isErr;
}

Response::Response(const std::string& msg_, bool isError_)
    : msg(msg_), isErr(isError_) {}

ErrorResponse::ErrorResponse(const std::string& msg_) : Response(msg_, true) {}

} // namespace http
} // namespace channels
} // namespace devmand
