// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <string>

#include <devmand/channels/Engine.h>

namespace devmand {
namespace channels {
namespace snmp {

class Engine final : public channels::Engine {
 public:
  Engine(const std::string& appName);

  Engine() = delete;
  ~Engine() override = default;
  Engine(const Engine&) = delete;
  Engine& operator=(const Engine&) = delete;
  Engine(Engine&&) = delete;
  Engine& operator=(Engine&&) = delete;

 public:
  void run();
  void stopEventually();

 private:
  bool stopping{false};
};

} // namespace snmp
} // namespace channels
} // namespace devmand
