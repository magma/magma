/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include <lte/protos/policydb.pb.h>
#include <lte/protos/session_manager.pb.h>

#include <functional>
#include <ostream>

namespace magma {
using namespace lte;

struct CreditKey {
  uint32_t rating_group;
  uint32_t service_identifier;
  CreditKey() : rating_group(0), service_identifier(0) {}
  CreditKey(uint32_t rg) : rating_group(rg), service_identifier(0) {}
  CreditKey(uint32_t rg, uint32_t si)
      : rating_group(rg), service_identifier(si) {}
  CreditKey(const PolicyRule* rule) { set(rule); }
  CreditKey(const PolicyRule& rule) { set(&rule); }
  CreditKey(const CreditUsage* usage) { set(usage); }
  CreditKey(const CreditUsage& usage) { set(&usage); }
  CreditKey(const CreditUpdateResponse* update) { set(update); }
  CreditKey(const CreditUpdateResponse& update) { set(&update); }
  CreditKey(const ChargingReAuthRequest* reauth) { set(reauth); }
  CreditKey(const ChargingReAuthRequest& reauth) { set(&reauth); }

  CreditKey* set(const PolicyRule* rule) {
    if (rule != nullptr) {
      rating_group       = rule->rating_group();
      service_identifier = rule->has_service_identifier() ?
                               rule->service_identifier().value() :
                               0;
    }
    return this;
  }
  void set_rule(PolicyRule* rule) const {
    if (rule != nullptr) {
      rule->set_rating_group(rating_group);
      if (service_identifier) {
        rule->mutable_service_identifier()->set_value(service_identifier);
      } else {
        rule->release_service_identifier();
      }
    }
  }
  CreditKey* set(const CreditUsage* usage) {
    if (usage != nullptr) {
      rating_group       = usage->charging_key();
      service_identifier = usage->has_service_identifier() ?
                               usage->service_identifier().value() :
                               0;
    }
    return this;
  }
  void set_credit_usage(CreditUsage* usage) const {
    if (usage != nullptr) {
      usage->set_charging_key(rating_group);
      if (service_identifier) {
        usage->mutable_service_identifier()->set_value(service_identifier);
      } else {
        usage->release_service_identifier();
      }
    }
  }
  CreditKey* set(const CreditUpdateResponse* update) {
    if (update != nullptr) {
      rating_group       = update->charging_key();
      service_identifier = update->has_service_identifier() ?
                               update->service_identifier().value() :
                               0;
    }
    return this;
  }
  CreditKey* set(const ChargingReAuthRequest* reath) {
    if (reath != nullptr) {
      rating_group       = reath->charging_key();
      service_identifier = reath->has_service_identifier() ?
                               reath->service_identifier().value() :
                               0;
    }
    return this;
  }
};

inline std::ostream& operator<<(std::ostream& s, const CreditKey& k) {
  s << "RG: " << k.rating_group;
  if (k.service_identifier) {
    s << ", SI: " << k.service_identifier;
  }
  return s;
}

inline size_t ccHash(const CreditKey& k) {
  static const int ccHashShift = sizeof(size_t) > sizeof(uint32_t) ? 32 : 1;
  size_t res = std::hash<uint32_t>()(k.rating_group) << ccHashShift;
  if (k.service_identifier) {
    res += std::hash<uint32_t>()(k.service_identifier);
  }
  return res;
};

inline bool ccEqual(const CreditKey& l, const CreditKey& r) {
  return (l.rating_group == r.rating_group) &&
         ((!l.service_identifier) ||
          (l.service_identifier == r.service_identifier));
};

}  // namespace magma
