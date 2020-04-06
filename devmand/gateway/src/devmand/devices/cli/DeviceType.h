// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <devmand/cartography/DeviceConfig.h>
#include <magma_logging.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

static const char* const ANY_VERSION = "*";

class DeviceType {
 private:
  string device;
  string version;

 public:
  DeviceType(const string& device, const string& version);

  static DeviceType getDefaultInstance() {
    return DeviceType("default", ANY_VERSION); // TODO
  }

  static DeviceType create(
      const devmand::cartography::DeviceConfig& deviceConfig);

  friend ostream& operator<<(ostream& os, const DeviceType& type);

  string str() const;

  bool operator==(const DeviceType& rhs) const;

  bool operator!=(const DeviceType& rhs) const;

  bool operator<(const DeviceType& rhs) const;

  bool operator>(const DeviceType& rhs) const;

  bool operator<=(const DeviceType& rhs) const;

  bool operator>=(const DeviceType& rhs) const;
};

} // namespace cli
} // namespace devices
} // namespace devmand
