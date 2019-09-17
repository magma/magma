// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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

  std::map<std::string, ChannelConfig> channelConfigs;

  friend std::ostream& operator<<(std::ostream& out, const DeviceConfig& c) {
    out << "id=" << c.id << ", "
        << "platform=" << c.platform << ", "
        << "ip=" << c.ip << ", "
        << "yangConfig=" << c.yangConfig << ", channels {";
    for (auto& channel : c.channelConfigs) {
      out << channel.first << ", ";
    }
    out << "}";
    return out;
  }

  friend bool operator<(const DeviceConfig& lhs, const DeviceConfig& rhs) {
    return std::tie(
               lhs.id,
               lhs.platform,
               lhs.ip,
               lhs.yangConfig,
               lhs.channelConfigs) <
        std::tie(
               rhs.id,
               rhs.platform,
               rhs.ip,
               rhs.yangConfig,
               rhs.channelConfigs);
  }

  friend bool operator==(const DeviceConfig& lhs, const DeviceConfig& rhs) {
    return std::tie(
               lhs.id,
               lhs.platform,
               lhs.ip,
               lhs.yangConfig,
               lhs.channelConfigs) ==
        std::tie(
               rhs.id,
               rhs.platform,
               rhs.ip,
               rhs.yangConfig,
               rhs.channelConfigs);
  }

  friend bool operator!=(const DeviceConfig& lhs, const DeviceConfig& rhs) {
    return std::tie(
               lhs.id,
               lhs.platform,
               lhs.ip,
               lhs.yangConfig,
               lhs.channelConfigs) !=
        std::tie(
               rhs.id,
               rhs.platform,
               rhs.ip,
               rhs.yangConfig,
               rhs.channelConfigs);
  }
};

using DeviceConfigs = std::set<DeviceConfig>;

} // namespace cartography
} // namespace devmand
