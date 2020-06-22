/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <functional>
#include <utility>

#include "StoredState.h"
#include "SessionCredit.h"

namespace magma {
struct Monitor {
  SessionCredit credit;
  MonitoringLevel level;

  Monitor() : credit(CreditType::MONITORING) {}

  static StoredMonitor marshal_monitor(std::unique_ptr<Monitor> &monitor) {
    StoredMonitor marshaled{};
    marshaled.credit = monitor->credit.marshal();
    marshaled.level = monitor->level;
    return marshaled;
  }

  static std::unique_ptr<Monitor> unmarshal_monitor(const StoredMonitor &marshaled) {
    Monitor monitor;
    monitor.credit = *SessionCredit::unmarshal(marshaled.credit, MONITORING);
    monitor.level = marshaled.level;
    return std::make_unique<Monitor>(monitor);
  }
};

}  // namespace magma
