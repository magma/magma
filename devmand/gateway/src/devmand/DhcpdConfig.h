// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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
