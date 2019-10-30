// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <vector>

#include <folly/dynamic.h>

#include <devmand/channels/snmp/Oid.h>

namespace devmand {
namespace channels {
namespace snmp {

class Response {
 public:
  Response(const Oid& _oid, const folly::dynamic& _value);
  Response() = delete;
  ~Response() = default;
  Response(const Response&) = default;
  Response& operator=(const Response&) = default;
  Response(Response&&) = default;
  Response& operator=(Response&&) = default;

  bool isError() const;

  friend bool operator==(const Response& lhs, const Response& rhs) {
    return lhs.oid == rhs.oid and lhs.value == lhs.value;
  }

 public:
  Oid oid;
  folly::dynamic value;
};

class ErrorResponse final : public Response {
 public:
  ErrorResponse(const folly::dynamic& _value);
  ErrorResponse() = delete;
  ~ErrorResponse() = default;
  ErrorResponse(const ErrorResponse&) = default;
  ErrorResponse& operator=(const ErrorResponse&) = default;
  ErrorResponse(ErrorResponse&&) = default;
  ErrorResponse& operator=(ErrorResponse&&) = default;
};

using Responses = std::vector<Response>;

} // namespace snmp
} // namespace channels
} // namespace devmand
