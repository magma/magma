// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/snmp/Response.h>

namespace devmand {
namespace channels {
namespace snmp {

Response::Response(const Oid& _oid, const folly::dynamic& _value)
    : oid(_oid), value(_value) {}

ErrorResponse::ErrorResponse(const folly::dynamic& _value)
    : Response(Oid::error, _value) {}

bool Response::isError() const {
  return oid.getLength() == 0;
}

} // namespace snmp
} // namespace channels
} // namespace devmand
