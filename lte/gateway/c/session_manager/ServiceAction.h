/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <memory>

namespace magma {
using namespace lte;

enum ServiceActionType {
  CONTINUE_SERVICE = 0,
  TERMINATE_SERVICE = 1,
  ACTIVATE_SERVICE = 2
};

/**
 * ServiceAction is the base class for any action that needs to be taken on
 * a subscriber's service. This could be terminate, redirect data, or just
 * continue.
 */
class ServiceAction {
public:
  ServiceAction(ServiceActionType action_type) : action_type_(action_type) {}

  ServiceActionType get_type() const {
    return action_type_;
  }

  ServiceAction& set_imsi(const std::string& imsi) {
    imsi_ = std::make_unique<std::string>(imsi);
    return *this;
  }

  ServiceAction& set_ip_addr(const std::string& ip_addr) {
    ip_addr_ = std::make_unique<std::string>(ip_addr);
    return *this;
  }

  /**
   * get_imsi returns the associated IMSI for the action, or throws a nullptr
   * exception if there is none stored
   */
  const std::string& get_imsi() const {
    return *imsi_;
  }

  /**
   * get_ip_addr returns the associated subscriber's ip_addr for the action,
   * or throws a nullptr exception if there is none stored
   */
  const std::string& get_ip_addr() const {
    return *ip_addr_;
  }


  const std::vector<std::string>& get_rule_ids() const {
    return rule_ids_;
  }

  const std::vector<PolicyRule>& get_rule_definitions() const {
    return rule_definitions_;
  }

  std::vector<std::string>* get_mutable_rule_ids() {
    return &rule_ids_;
  }

  std::vector<PolicyRule>* get_mutable_rule_definitions() {
    return &rule_definitions_;
  }

private:
  ServiceActionType action_type_;
  std::unique_ptr<std::string> imsi_;
  std::unique_ptr<std::string> ip_addr_;
  std::vector<std::string> rule_ids_;
  std::vector<PolicyRule> rule_definitions_;
};

}
