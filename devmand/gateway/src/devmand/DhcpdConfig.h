// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <set>
#include <string>
#include <tuple>

#include <folly/IPAddress.h>
#include <folly/MacAddress.h>

namespace devmand {

struct Host {
  const std::string name;
  const folly::MacAddress mac;
  const folly::IPAddress ip;

  friend bool operator<(const Host& lhs, const Host& rhs) {
    return std::tie(lhs.name, lhs.mac, lhs.ip) <
        std::tie(rhs.name, rhs.mac, rhs.ip);
  }

  friend bool operator==(const Host& lhs, const Host& rhs) {
    return std::tie(lhs.name, lhs.mac, lhs.ip) ==
        std::tie(rhs.name, rhs.mac, rhs.ip);
  }

  friend bool operator!=(const Host& lhs, const Host& rhs) {
    return std::tie(lhs.name, lhs.mac, lhs.ip) !=
        std::tie(rhs.name, rhs.mac, rhs.ip);
  }
};

class DhcpdConfig final {
 public:
  DhcpdConfig();
  ~DhcpdConfig() = default;
  DhcpdConfig(const DhcpdConfig&) = delete;
  DhcpdConfig& operator=(const DhcpdConfig&) = delete;
  DhcpdConfig(DhcpdConfig&&) = delete;
  DhcpdConfig& operator=(DhcpdConfig&&) = delete;

 public:
  void add(Host& host);

  void remove(Host& host);

 private:
  void rewrite();

  static std::string getHostSection(const Host& host);

 private:
  std::set<Host> hosts;
};

} // namespace devmand
