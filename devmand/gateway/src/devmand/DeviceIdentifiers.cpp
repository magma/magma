// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/DeviceIdentifiers.h>

namespace devmand {

void DeviceIdentifiers::addIdentifier(
    const std::string& identifer,
    const devices::Id& id) {
  identifiers.emplace(identifer, id);
}

void DeviceIdentifiers::removeIdentifier(
    const std::string& identifer,
    const devices::Id& id) {
  auto range = identifiers.equal_range(identifer);
  for (auto it = range.first; it != range.second; ++it) {
    if (it->second == id) {
      identifiers.erase(it);
      break;
    }
  }
}

devices::Id DeviceIdentifiers::lookup(const std::string& identifer) const {
  auto range = identifiers.equal_range(identifer);
  for (auto it = range.first; it != range.second; ++it) {
    return it->second;
  }
  return "";
}

} // namespace devmand
