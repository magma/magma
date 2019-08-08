/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#pragma once

#include <map>
#include <memory>
#include <stdint.h>

#include "sctp_assoc.h"

namespace magma {
namespace sctpd {

using AssocMap = std::map<uint32_t, SctpAssoc>;

// Models the state of an SCTP connection and its assocations
class SctpDesc {
 public:
  // Construct a SCTP assocation on socket, sd
  SctpDesc(int sd);

  // Add assocation, assoc, to the list of assocations - keyed by assoc_id
  void addAssoc(const SctpAssoc &assoc);
  // Get association keyed by assoc_id, throw std::out_of_range otherwise
  SctpAssoc &getAssoc(uint32_t assoc_id);
  // Remove assoc keyed by assoc_id from assoc list, returns 0/-1 on ok/fail
  int delAssoc(uint32_t assoc_id);

  // Return the starting const_iterator of associations in the SCTP connection
  AssocMap::const_iterator begin() const;
  // Return the ending const_iterator of associations in the SCTP connection
  AssocMap::const_iterator end() const;

  // Return the socket descriptor for the SCTP connection
  int sd() const;

  // Dump debug information about the SCTP connection to the log
  void dump() const;

 private:
  // List (map) of assocations for the SCTP connection
  AssocMap _assocs;
  // Socket descriptor for the SCTP connection
  int _sd;
};

} // namespace sctpd
} // namespace magma
