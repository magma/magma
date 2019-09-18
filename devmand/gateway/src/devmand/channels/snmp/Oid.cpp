// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <assert.h>
#include <cstring>
#include <sstream>
#include <stdexcept>

#include <devmand/channels/snmp/Oid.h>

namespace devmand {
namespace channels {
namespace snmp {

const Oid Oid::error{""};

Oid::Oid(oid* buf, size_t len) : length(len) {
  std::memcpy(buffer, buf, length * sizeof(oid));
}

Oid::Oid(const std::string& repr) {
  if (repr.empty()) {
    length = 0;
    return;
  }

  if (read_objid(repr.c_str(), buffer, &length) == 0) {
    throw std::runtime_error("read_objid error");
  }
}

std::string Oid::toString() const {
  std::stringstream str;
  for (unsigned int i = 0; i < length; ++i) {
    str << "." << folly::to<std::string>(buffer[i]);
  }
  return str.str();
}

const oid* const Oid::get() const {
  return buffer;
}

size_t Oid::getLength() const {
  return length;
}

} // namespace snmp
} // namespace channels
} // namespace devmand
