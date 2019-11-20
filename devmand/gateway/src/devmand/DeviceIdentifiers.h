// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <map>

#include <devmand/devices/Id.h>

namespace devmand {

// This class provides a mapping between
class DeviceIdentifiers final {
 public:
  DeviceIdentifiers() = default;
  ~DeviceIdentifiers() = default;
  DeviceIdentifiers(const DeviceIdentifiers&) = delete;
  DeviceIdentifiers& operator=(const DeviceIdentifiers&) = delete;
  DeviceIdentifiers(DeviceIdentifiers&&) = delete;
  DeviceIdentifiers& operator=(DeviceIdentifiers&&) = delete;

 public:
  void addIdentifier(const std::string& identifer, const devices::Id& id);
  void removeIdentifier(const std::string& identifer, const devices::Id& id);

  devices::Id lookup(const std::string& identifer) const;

 private:
  std::multimap<std::string, devices::Id> identifiers;
};

} // namespace devmand
