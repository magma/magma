// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/dynamic.h>

#include <devmand/utils/YangUtils.h>

namespace devmand {
namespace models {
namespace wifi {

class Model {
 public:
  Model() = delete;
  ~Model() = delete;
  Model(const Model&) = delete;
  Model& operator=(const Model&) = delete;
  Model(Model&&) = delete;
  Model& operator=(Model&&) = delete;

 public:
  static void init(folly::dynamic& state);
  static void updateRadio(
      folly::dynamic& state,
      int index,
      const YangPath& path,
      const folly::dynamic& value);

  static void updateSsid(
      folly::dynamic& state,
      int index,
      const YangPath& path,
      const folly::dynamic& value);

  static void updateSsidBssid(
      folly::dynamic& state,
      int indexSsid,
      int indexBssid,
      const YangPath& path,
      const folly::dynamic& value);
};

} // namespace wifi
} // namespace models
} // namespace devmand
