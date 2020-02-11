// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/dynamic.h>

namespace devmand::channels::cli::datastore {
using folly::dynamic;
using std::string;

enum DatastoreDiffType { create, update, deleted };

struct DatastoreDiff {
  const dynamic before = nullptr;
  const dynamic after = nullptr;
  const DatastoreDiffType type;
  DatastoreDiff(
      const dynamic& _before,
      const dynamic& _after,
      const DatastoreDiffType _type)
      : before(_before), after(_after), type(_type) {}
};
} // namespace devmand::channels::cli::datastore
