// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/snmp/Snmp.h>

namespace devmand {
namespace channels {
namespace snmp {

/* A simple implementation of snmp pdu for managing lifetime. For now this
 * doesn't support much more than a single mib as that is the only current
 * need.
 */
class Pdu final {
 public:
  Pdu(int type, const Oid& _oid);
  Pdu() = delete;
  ~Pdu();
  Pdu(const Pdu&) = delete;
  Pdu& operator=(const Pdu&) = delete;
  Pdu(Pdu&&) = delete;
  Pdu& operator=(Pdu&&) = delete;

 public:
  snmp_pdu* get();

  // Once the pdu ownership has been transfered to the library release it
  // TODO dig into this more...
  void release();

 private:
  snmp_pdu* pdu{nullptr};
};

} // namespace snmp
} // namespace channels
} // namespace devmand
