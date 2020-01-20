// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/Service.h>

namespace devmand {
namespace fscache {

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
};

} // namespace fscache
} // namespace devmand
