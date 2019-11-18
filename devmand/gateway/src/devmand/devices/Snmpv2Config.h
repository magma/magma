// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <list>
#include <map>
#include <set>
#include <string>

#include <devmand/channels/snmp/Channel.h>

namespace devmand {
namespace devices {

// These enums describe items in the config file.
// To understand what is happening here, refer to:
// https://fburl.com/6w5jyint

enum ResultType { SINGLE, ITERABLE };

enum ActionType { STORE, SET_GAUGE_COUNTER };

enum ActionArgType { TOP_LEVEL, PATH, VALUE };

enum ActionTopLevelEnum { IETF_SYSTEM, OPENCONFIG_INTERFACES };

enum ActionValueEnum {
  RESULT,
  RESULT_TO_STATUS_STR,
  RESULT_IS_ENABLED_BOOLEAN
};

struct ActionArgValue {
  ActionTopLevelEnum topLevel;
  std::string path;
  ActionValueEnum value;
};

using ActionArgMap = std::map<ActionArgType, ActionArgValue>;

struct Action {
  ActionType actionType;
  ActionArgMap args;
};

using ActionList = std::list<Action>;

struct OidConfig {
  channels::snmp::Oid oid;
  ResultType resultType;
  ActionList resultActions;
  ActionList forEachResultActions;
};

using OidConfigList = std::list<OidConfig>;

enum GeneralConfigType {
  INTERFACE_INDICES,
  NUMBER_OF_INTERFACES,
  POLLING_INTERVAL,
  SNMP_TIMEOUT_INTERVAL,
  SNMP_RETRIES
};

using GeneralConfigMap = std::map<GeneralConfigType, std::string>;

class Snmpv2Config final {
 public:
  Snmpv2Config();
  ~Snmpv2Config() = default;
  Snmpv2Config(const Snmpv2Config&) = delete;
  Snmpv2Config& operator=(const Snmpv2Config&) = delete;
  Snmpv2Config(Snmpv2Config&&) = delete;
  Snmpv2Config& operator=(Snmpv2Config&&) = delete;

 public:
  const GeneralConfigMap& getGeneralConfig();
  const OidConfigList& getOidConfigs();

 private:
  GeneralConfigMap generalConfigs;
  OidConfigList oidConfigs;
};

} // namespace devices
} // namespace devmand
