// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <net-snmp/net-snmp-config.h>
#include <net-snmp/net-snmp-includes.h>
#undef FREE // snmp...
#undef READ
#undef WRITE

#include <folly/Conv.h>

namespace devmand {
namespace channels {
namespace snmp {

class Oid final {
 public:
  Oid(oid* buf, size_t len);
  Oid(const std::string& repr);
  Oid() = delete;
  ~Oid() = default;
  Oid(const Oid&) = default;
  Oid& operator=(const Oid&) = delete;
  Oid(Oid&&) = default;
  Oid& operator=(Oid&&) = delete;

 public:
  const oid* const get() const;
  size_t getLength() const;
  std::string toString() const;

 public:
  static const Oid error;

 private:
  size_t length{MAX_OID_LEN};
  oid buffer[MAX_OID_LEN];
};

} // namespace snmp
} // namespace channels
} // namespace devmand
