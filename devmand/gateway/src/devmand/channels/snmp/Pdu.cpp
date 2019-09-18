// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <assert.h>
#include <stdexcept>

#include <devmand/channels/snmp/Pdu.h>

namespace devmand {
namespace channels {
namespace snmp {

Pdu::Pdu(int type, const Oid& _oid) : pdu(snmp_pdu_create(type)) {
  if (pdu == nullptr) {
    throw std::runtime_error("snmp_pdu_create error");
  }

  if (snmp_add_null_var(pdu, _oid.get(), _oid.getLength()) == nullptr) {
    throw std::runtime_error("snmp_add_null_var error");
  }
}

Pdu::~Pdu() {
  if (pdu != nullptr) {
    snmp_free_pdu(pdu);
  }
}

void Pdu::release() {
  assert(pdu != nullptr);
  pdu = nullptr;
}

snmp_pdu* Pdu::get() {
  assert(pdu != nullptr);
  return pdu;
}

} // namespace snmp
} // namespace channels
} // namespace devmand
