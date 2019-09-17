// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <devmand/Service.h>

#include <MagmaService.h>

namespace devmand {
namespace magma {

class Service : public ::devmand::Service {
 public:
  Service(Application& application);
  Service() = delete;
  virtual ~Service() = default;
  Service(const Service&) = delete;
  Service& operator=(const Service&) = delete;
  Service(Service&&) = delete;
  Service& operator=(Service&&) = delete;

 public:
  void start() override;
  void wait() override;
  void stop() override;
  void setGauge(const std::string& key, double value) override;

 private:
  ::magma::service303::MagmaService magmaService;
};

} // namespace magma
} // namespace devmand
