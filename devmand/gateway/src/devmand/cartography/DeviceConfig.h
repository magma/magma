// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <map>
#include <set>
#include <string>
#include <tuple>

#include <iostream>

namespace devmand {
namespace cartography {

struct ChannelConfig {
  std::map<std::string, std::string> kvPairs;

  friend bool operator<(const ChannelConfig& lhs, const ChannelConfig& rhs) {
    return lhs.kvPairs < rhs.kvPairs;
  }

  friend bool operator==(const ChannelConfig& lhs, const ChannelConfig& rhs) {
    return lhs.kvPairs == rhs.kvPairs;
  }

  friend bool operator!=(const ChannelConfig& lhs, const ChannelConfig& rhs) {
    return lhs.kvPairs != rhs.kvPairs;
  }
};

struct DeviceConfig {
  std::string id;
  std::string platform;
  std::string ip;
  std::string yangConfig;
  bool readonly{false};

  std::map<std::string, ChannelConfig> channelConfigs;

  friend std::ostream& operator<<(std::ostream& out, const DeviceConfig& c) {
    out << "id=" << c.id << ", "
        << "platform=" << c.platform << ", "
        << "ip=" << c.ip << ", "
        << "yangConfig=" << c.yangConfig << ", "
        << "readonly=" << c.readonly << ", channels {";
    for (auto& channel : c.channelConfigs) {
      out << channel.first << ", ";
    }
    out << "}";
    return out;
  }

  friend bool operator<(const DeviceConfig& lhs, const DeviceConfig& rhs) {
    return std::tie(
               lhs.id, lhs.platform, lhs.ip, lhs.readonly, lhs.channelConfigs) <
        std::tie(
               rhs.id, rhs.platform, rhs.ip, rhs.readonly, rhs.channelConfigs);
  }

  friend bool operator==(const DeviceConfig& lhs, const DeviceConfig& rhs) {
    return std::tie(
               lhs.id,
               lhs.platform,
               lhs.ip,
               lhs.yangConfig,
               lhs.readonly,
               lhs.channelConfigs) ==
        std::tie(
               rhs.id,
               rhs.platform,
               rhs.ip,
               rhs.yangConfig,
               rhs.readonly,
               rhs.channelConfigs);
  }

  friend bool operator!=(const DeviceConfig& lhs, const DeviceConfig& rhs) {
    return std::tie(
               lhs.id,
               lhs.platform,
               lhs.ip,
               lhs.yangConfig,
               lhs.readonly,
               lhs.channelConfigs) !=
        std::tie(
               rhs.id,
               rhs.platform,
               rhs.ip,
               rhs.yangConfig,
               rhs.readonly,
               rhs.channelConfigs);
  }
};

using DeviceConfigs = std::set<DeviceConfig>;

} // namespace cartography
} // namespace devmand
