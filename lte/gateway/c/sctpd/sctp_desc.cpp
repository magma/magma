/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "sctp_desc.h"

#include "assert.h"

namespace magma {
namespace sctpd {

SctpDesc::SctpDesc(int sd): _sd(sd)
{
  assert(sd >= 0);
}

void SctpDesc::addAssoc(const SctpAssoc &assoc)
{
  _assocs[assoc.assoc_id] = assoc;
}

SctpAssoc &SctpDesc::getAssoc(uint32_t assoc_id)
{
  return _assocs.at(assoc_id); // throws std::out_of_range
}

int SctpDesc::delAssoc(uint32_t assoc_id)
{
  auto num_removed = _assocs.erase(assoc_id);
  return num_removed == 1 ? 0 : -1;
}

AssocMap::const_iterator SctpDesc::begin() const
{
  return _assocs.cbegin();
}

AssocMap::const_iterator SctpDesc::end() const
{
  return _assocs.cend();
}

int SctpDesc::sd() const
{
  return _sd;
}

void SctpDesc::dump() const
{
  for (auto const &kv : _assocs) {
    auto assoc = kv.second;
    assoc.dump();
  }
}

} // namespace sctpd
} // namespace magma
