// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <string>

namespace devmand {
namespace channels {
namespace cnmaestro {

class Response {
 public:
  Response(const std::string& msg_);
  Response() = delete;
  ~Response() = default;
  Response(const Response&) = delete;
  Response& operator=(const Response&) = delete;
  Response(Response&&) = default;
  Response& operator=(Response&&) = delete;

 public:
  std::string get() const;

  bool isError() const;

 protected:
  Response(const std::string& msg_, bool isError_);

 private:
  std::string msg;
  const bool isErr{false};
};

class ErrorResponse : public Response {
 public:
  ErrorResponse(const std::string& msg);
  ErrorResponse() = delete;
  ~ErrorResponse() = default;
  ErrorResponse(const ErrorResponse&) = delete;
  ErrorResponse& operator=(const ErrorResponse&) = delete;
  ErrorResponse(ErrorResponse&&) = default;
  ErrorResponse& operator=(ErrorResponse&&) = delete;
};

} // namespace cnmaestro
} // namespace channels
} // namespace devmand
