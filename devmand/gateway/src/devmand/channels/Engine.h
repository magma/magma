// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

namespace devmand {
namespace channels {

/*
 * This is the common interface all channel engines must implement. A channel
 * engine lives for the duration of the process and maintains state the is
 * independent of individual channels. Not all channels require engines: just
 * ones that maintain per process state.
 */
class Engine {
 public:
  Engine() = default;
  virtual ~Engine() = default;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;
};

} // namespace channels
} // namespace devmand
