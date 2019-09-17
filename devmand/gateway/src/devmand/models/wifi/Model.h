// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <folly/dynamic.h>

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
};

} // namespace wifi
} // namespace models
} // namespace devmand
