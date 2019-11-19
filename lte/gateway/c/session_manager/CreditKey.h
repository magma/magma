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
#include <ostream>

#include <lte/protos/policydb.pb.h>
#include <lte/protos/session_manager.pb.h>

namespace magma {
using namespace lte;

struct CreditKey {
  uint32_t rating_group;
  uint32_t service_identifier;
  bool use_sid;
  CreditKey() : rating_group(0), service_identifier(0), use_sid(false) {}
  CreditKey(uint32_t rg) : rating_group(rg), service_identifier(0), use_sid(false) {}
  CreditKey(uint32_t rg, uint32_t si) : rating_group(rg), service_identifier(si), use_sid(true) {}
  CreditKey(const PolicyRule *rule)   { set(rule); }
  CreditKey(const PolicyRule &rule)   { set(&rule); }
  CreditKey(const CreditUsage *usage) { set(usage); }
  CreditKey(const CreditUsage& usage) { set(&usage); }
  CreditKey(const CreditUpdateResponse *update) { set(update); }
  CreditKey(const CreditUpdateResponse& update) { set(&update); }
  CreditKey(const ChargingReAuthRequest *reauth) { set(reauth); }
  CreditKey(const ChargingReAuthRequest& reauth) { set(&reauth); }

  CreditKey *set(const PolicyRule *rule) {
    if (rule != nullptr) {
      rating_group = rule->rating_group();
      use_sid = rule->has_service_identifier();
      service_identifier = use_sid ? rule->service_identifier().value() : 0;
    }
    return this;
  }
  void set_rule(PolicyRule *rule) const {
    if (rule != nullptr) {
      rule->set_rating_group(rating_group);
      if (use_sid) {
        rule->mutable_service_identifier()->set_value(service_identifier);
      } else {
        rule->release_service_identifier();
      }
    }
  }
  CreditKey *set(const CreditUsage *usage) {
    if (usage != nullptr) {
      rating_group = usage->charging_key();
      use_sid = usage->has_service_identifier();
      service_identifier = use_sid ? usage->service_identifier().value() : 0;
    }
    return this;
  }
  void set_credit_usage(CreditUsage *usage) const {
    if (usage != nullptr) {
      usage->set_charging_key(rating_group);
      if (use_sid) {
        usage->mutable_service_identifier()->set_value(service_identifier);
      } else {
        usage->release_service_identifier();
      }
    }
  }
  CreditKey *set(const CreditUpdateResponse *update) {
    if (update != nullptr) {
      rating_group = update->charging_key();
      use_sid = update->has_service_identifier();
      service_identifier = use_sid ? update->service_identifier().value() : 0;
    }
    return this;
  }
  CreditKey *set(const ChargingReAuthRequest *reath) {
    if (reath != nullptr) {
      rating_group = reath->charging_key();
      use_sid = reath->has_service_identifier();
      service_identifier = use_sid ? reath->service_identifier().value() : 0;
    }
    return this;
  }
};

inline std::ostream &operator<<(std::ostream &s, const CreditKey &k) {
  s << "RG: " << k.rating_group;
  if (k.use_sid) {
    s << ", SI: " << k.service_identifier;
  }
  return s;
}

inline size_t ccHash(const CreditKey& k) {
  static const int ccHashShift = sizeof(size_t) > sizeof(uint32_t) ? 32 : 1;
  size_t res = std::hash<uint32_t>()(k.rating_group) << ccHashShift;
  if (k.use_sid) {
    res += std::hash<uint32_t>()(k.service_identifier);
  }
  return res;
};

inline bool ccEqual(const CreditKey& l, const CreditKey& r){
  return (l.rating_group == r.rating_group) &&
         (l.use_sid == r.use_sid) &&
         ((!l.use_sid) || (l.service_identifier == r.service_identifier));
};

} // namespace magma
