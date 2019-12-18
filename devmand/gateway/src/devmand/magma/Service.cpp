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
  magmaService.SetServiceInfoCallback([this]() { return getServiceInfo(); });

  magmaService.SetOperationalStatesCallback(
      [this]() { return getOperationalStates(); });
}

std::list<std::map<std::string, std::string>> Service::getOperationalStates() {
  auto unifiedView = app.getUnifiedView();
  std::list<std::map<std::string, std::string>> states;

  for (auto& device : unifiedView) {
    std::string deviceId = device.first;
    folly::dynamic deviceState = folly::dynamic::object;
    deviceState["raw_state"] = folly::toJson(device.second);

    std::map<std::string, std::string> state = {
        {"type", orc8rDeviceType},
        {"device_id", deviceId},
        {"value", folly::toJson(deviceState)}};
    states.push_back(state);
  }
  return states;
}

std::map<std::string, std::string> Service::getServiceInfo() {
  auto unifiedView = app.getUnifiedView();

  folly::dynamic devices = folly::dynamic::object;
  for (auto& device : unifiedView) {
    devices[device.first] = device.second;
  }

  return std::map<std::string, std::string>{
      {"devmand", folly::toJson(devices)}};
}

void Service::setGauge(
    const std::string& key,
    double value,
    const std::string& labelName,
    const std::string& labelValue) {
  if (labelName.empty() or labelValue.empty()) {
    setGaugeVA(key, value, 0);
  } else {
    setGaugeVA(key, value, 1, labelName.c_str(), labelValue.c_str());
  }
}

void Service::setGaugeVA(
    const std::string& key,
    double value,
    size_t labelCount,
    ...) {
  va_list labels;
  va_start(labels, labelCount);
  ::magma::service303::MetricsSingleton::Instance().SetGauge(
      key.c_str(), value, labelCount, labels);
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
