// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <MetricsSingleton.h>

#include <folly/GLog.h>
#include <folly/json.h>

#include <devmand/Application.h>
#include <devmand/magma/Service.h>
#include <orc8r/protos/service303.grpc.pb.h>
#include <orc8r/protos/service303.pb.h>

namespace devmand {
namespace magma {

Service::Service(Application& application)
    : ::devmand::Service(application),
      magmaService(app.getName(), app.getVersion()) {
  magmaService.SetServiceInfoCallback([this]() {
    auto uv = app.getUnifiedView();

    LOG(INFO) << "publishing :=\n";
    for (auto& kv : uv) {
      LOG(INFO) << "\t\"" << kv.first << "\" : \"" << kv.second << "\"\n";
    }
    return uv;
  });

  magmaService.SetOperationalStatesCallback([this]() {
    auto uv = app.getUnifiedView();
    std::list<std::map<std::string, std::string>> states;

    folly::dynamic devices = uv.find("devmand") != uv.end()
        ? folly::parseJson(uv["devmand"])
        : folly::dynamic::object;

    LOG(INFO) << "Publishing op-state :=\n";
    for (auto& device : devices.items()) {
      std::string device_id = device.first.asString();
      folly::dynamic device_state = folly::dynamic::object;
      device_state["raw_state"] = folly::toJson(device.second);

      std::map<std::string, std::string> state = {
          {"type", orc8rDeviceType},
          {"device_id", device.first.asString()},
          {"value", folly::toJson(device_state)}};
      LOG(INFO) << device_id << " : " << folly::toJson(device_state) << "\n";
      states.push_back(state);
    }
    return states;
  });
}

void Service::setGauge(const std::string& key, double value) {
  setGaugeVA(key, value, 0);
}

void Service::setGauge(
    const std::string& key,
    double value,
    const std::string& label_name,
    const std::string& label_value) {
  if (label_name.length() == 0 || label_value.length() == 0) {
    setGaugeVA(key, value, 0);
  } else {
    setGaugeVA(key, value, 1, label_name.c_str(), label_value.c_str());
  }
}

void Service::setGaugeVA(
    const std::string& key,
    double value,
    size_t label_count,
    ...) {
  va_list labels;
  va_start(labels, label_count);
  ::magma::service303::MetricsSingleton::Instance().SetGauge(
      key.c_str(), value, label_count, labels);
  va_end(labels);
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
