// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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
  void setGauge(
      const std::string& key,
      double value,
      const std::string& labelName,
      const std::string& labelValue) override;

 private:
  void setGaugeVA(const std::string& key, double value, size_t labelCount, ...);

  std::list<std::map<std::string, std::string>> getOperationalStates();
  std::map<std::string, std::string> getServiceInfo();

  // The key that will tell orc8r how to store these states
  static constexpr auto orc8rDeviceType = "symphony_device";

 private:
  ::magma::service303::MagmaService magmaService;
};

} // namespace magma
} // namespace devmand
