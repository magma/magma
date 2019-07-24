/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sctp_assoc.h"

#include <iostream>

#include "util.h"

namespace magma {
namespace sctpd {

SctpAssoc::SctpAssoc():
  sd(0),
  ppid(0),
  instreams(0),
  outstreams(0),
  assoc_id(0),
  messages_recv(0),
  messages_sent(0)
{
}

void SctpAssoc::dump() const
{
  MLOG(MDEBUG) << "SctpAssoc<id: " << std::to_string(this->assoc_id) << ">"
               << std::endl;
}

} // namespace sctpd
} // namespace magma
