// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
