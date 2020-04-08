/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "CreditKey.h"
#include "StoredState.h"

namespace magma {

SessionStateUpdateCriteria get_default_update_criteria()
{
  SessionStateUpdateCriteria uc {};
  uc.is_config_updated = false;
  uc.charging_credit_to_install = std::unordered_map<
    CreditKey,
    StoredSessionCredit,
    decltype(&ccHash),
    decltype(&ccEqual)>(4, &ccHash, &ccEqual);
  uc.charging_credit_map = std::unordered_map<
    CreditKey,
    SessionCreditUpdateCriteria,
    decltype(&ccHash),
    decltype(&ccEqual)>(4, &ccHash, &ccEqual);
  return uc;
}

}; // namespace magma
