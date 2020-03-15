// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/DeviceType.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

bool DeviceType::operator==(const DeviceType& rhs) const {
  return device == rhs.device && version == rhs.version;
}

bool DeviceType::operator!=(const DeviceType& rhs) const {
  return !(rhs == *this);
}

bool DeviceType::operator<(const DeviceType& rhs) const {
  if (device < rhs.device)
    return true;
  if (rhs.device < device)
    return false;
  return version < rhs.version;
}

bool DeviceType::operator>(const DeviceType& rhs) const {
  return rhs < *this;
}

bool DeviceType::operator<=(const DeviceType& rhs) const {
  return !(rhs < *this);
}

bool DeviceType::operator>=(const DeviceType& rhs) const {
  return !(*this < rhs);
}

ostream& operator<<(ostream& os, const DeviceType& type) {
  os << "{" << type.device << ": " << type.version << "}";
  return os;
}

string DeviceType::str() const {
  stringstream strStream = stringstream();
  strStream << *this;
  return strStream.str();
}

DeviceType::DeviceType(const string& _device, const string& _version)
    : device(_device), version(_version) {}

DeviceType DeviceType::create(
    const devmand::cartography::DeviceConfig& deviceConfig) {
  const std::map<std::string, std::string>& config =
      deviceConfig.channelConfigs.at("cli").kvPairs;

  if (config.find("flavour") != config.end()) {
    string device = config.at("flavour");
    string version = "*";
    if (config.find("flavourVersion") != config.end()) {
      version = config.at("flavourVersion");
    }
    return DeviceType{device, version};
  } else {
    return DeviceType::getDefaultInstance();
  }
}

} // namespace cli
} // namespace devices
} // namespace devmand
