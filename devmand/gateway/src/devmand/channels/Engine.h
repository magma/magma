// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
