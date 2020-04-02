// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/devices/cli/schema/Path.h>
#include <folly/dynamic.h>

namespace devmand::channels::cli::datastore {
using devmand::devices::cli::Path;
using folly::dynamic;
using std::string;
enum DatastoreDiffType { create, update, deleted };

struct DatastoreDiff {
  const dynamic before = nullptr;
  const dynamic after = nullptr;
  const DatastoreDiffType type;
  const Path keyedPath;
  const Path path;

  string toText() {
    return string("change on keyedpath: ") + keyedPath.str();
  }
  DatastoreDiff(
      const dynamic& _before,
      const dynamic& _after,
      const DatastoreDiffType _type,
      const Path _path)
      : before(_before),
        after(_after),
        type(_type),
        keyedPath(_path),
        path(_path.unkeyed()) {}

  DatastoreDiff(const DatastoreDiff& diff)
      : before(diff.before),
        after(diff.after),
        type(diff.type),
        keyedPath(diff.keyedPath),
        path(diff.path) {}
};
} // namespace devmand::channels::cli::datastore
