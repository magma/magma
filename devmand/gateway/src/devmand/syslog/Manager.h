// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <map>

#include <devmand/devices/Id.h>
#include <devmand/utils/ConfigGenerator.h>

namespace devmand {
namespace syslog {

// This class provides a mapping between hostnames and ip addresses
class Manager final {
 public:
  Manager();
  ~Manager() = default;
  Manager(const Manager&) = delete;
  Manager& operator=(const Manager&) = delete;
  Manager(Manager&&) = delete;
  Manager& operator=(Manager&&) = delete;

 public:
  void addIdentifier(const std::string& identifer, const devices::Id& id);
  void removeIdentifier(const std::string& identifer, const devices::Id& id);

  devices::Id lookup(const std::string& identifer) const;

  void restartTdAgentBitAsync() const;

 private:
  // This needs to keep track of hostnames and ip addresses.
  std::multimap<std::string, devices::Id> identifiers;
  utils::ConfigGenerator configGenerator;

  static constexpr const char* configFile =
      "/etc/magma/td-agent-bit-devmand.conf";

  static constexpr const char* configTemplate = R"template([FILTER]
    Name modify
    Match *
    Condition Key_Value_Equals host {}
    Set device {}
)template";
};

} // namespace syslog
} // namespace devmand
