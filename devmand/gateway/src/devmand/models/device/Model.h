// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/dynamic.h>

#include <devmand/channels/ping/Engine.h>

namespace devmand {
namespace models {
namespace device {

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

  static void addLatency(
      folly::dynamic& state,
      const std::string& type,
      const std::string& src,
      const std::string& dst,
      channels::ping::Rtt rtt);
};

} // namespace device
} // namespace models
} // namespace devmand
