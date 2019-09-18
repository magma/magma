// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <MetricsSingleton.h>

#include <devmand/Application.h>
#include <devmand/magma/Service.h>

namespace devmand {
namespace magma {

Service::Service(Application& application)
    : ::devmand::Service(application),
      magmaService(app.getName(), app.getVersion()) {
  magmaService.SetServiceInfoCallback([this]() {
    auto uv = app.getUnifiedView();

    std::cerr << "publishing :=\n";
    for (auto& kv : uv) {
      std::cerr << "\t\"" << kv.first << "\" : \"" << kv.second << "\"\n";
    }
    return uv;
  });
}

void Service::setGauge(const std::string& key, double value) {
  va_list ap;
  ::magma::service303::MetricsSingleton::Instance().SetGauge(
      key.c_str(), value, 0, ap);
}

void Service::start() {
  magmaService.Start();
}

void Service::wait() {
  magmaService.WaitForShutdown();
}

void Service::stop() {
  magmaService.Stop();
}

} // namespace magma
} // namespace devmand
