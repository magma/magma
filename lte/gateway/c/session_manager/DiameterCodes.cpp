/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "DiameterCodes.h"

namespace magma {
bool DiameterCodeHandler::is_transient_failure(const uint32_t code)
{
  return 4000 <= code && code < 5000;
}

// Diameter code of form 5xxx marks a permanent failure
bool DiameterCodeHandler::is_permanent_failure(const uint32_t code)
{
  return 5000 <= code && code < 6000;
}
} // namespace magma
