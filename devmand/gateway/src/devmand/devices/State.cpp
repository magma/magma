// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/State.h>

#include <devmand/MetricSink.h>
#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {

std::shared_ptr<State> State::make(MetricSink& sink_, const Id& device_) {
  return std::shared_ptr<State>(new State(sink_, device_));
}

State::State(MetricSink& sink_, const Id& device_)
    : sink(sink_), device(device_), state(folly::dynamic::object) {}

folly::Future<folly::dynamic> State::collect() {
  return folly::collect(std::move(requests))
      .thenValue([s = shared_from_this()](auto reqs) {
        s->setErrors();
        s->setStatus(true);
        for (auto& f : s->finals) {
          f();
        }
        s->clear();
        reqs.clear();
        return std::move(*(s->state.wlock()));
      });
}

void State::clear() {
  finals.clear();
  requests.clear();
}

void State::update(std::function<void(folly::dynamic&)> func) {
  state.withWLock([&func](auto& unlockedState){
    func(unlockedState);
  });
}

void State::addFinally(std::function<void()>&& f) {
  finals.push_back(f);
}

void State::addError(std::string&& error) {
  errorQueue.add(std::forward<std::string>(error));
}

void State::addRequest(folly::Future<folly::Unit> future) {
  auto lifetimeFuture = std::move(future).thenValue(
      [s = shared_from_this()](auto v) { return v; });
  requests.push_back(
      std::move(lifetimeFuture)
          .thenError(
              folly::tag_t<std::exception>{},
              [s = shared_from_this()](std::exception const& e) {
                LOG(ERROR) << "Caught exception from future: " << e.what();
                s->errorQueue.add(folly::sformat(
                    "Caught exception from future: {}", e.what()));
              })
          .thenError([s = shared_from_this()](folly::exception_wrapper) {
            LOG(ERROR) << "Caught unknown exception from future";
            s->errorQueue.add("Caught unknown exception from future");
          }));
}

folly::dynamic& State::getFbcPlatformDevice(
    const std::string& key,
    folly::dynamic& unlockedState) {
  auto k = folly::sformat("fbc-symphony-device:{}", key);
  folly::dynamic* system = unlockedState.get_ptr(k);
  if (system == nullptr) {
    system = &(unlockedState[k] = folly::dynamic::object);
  }
  return *system;
}

void State::setStatus(bool systemIsUp) {
  state.withWLock([this, systemIsUp](auto& unlockedState) {
    folly::dynamic& system =
        this->getFbcPlatformDevice("system", unlockedState);
    system["status"] = systemIsUp ? "UP" : "DOWN";
  });
  setGauge("device.status", systemIsUp ? 1 : 0);
}

void State::setErrors() {
  folly::dynamic errors = errorQueue.get();
  if (not errors.empty()) {
    state.withWLock([&errors](auto& lockedState) {
      lockedState["fbc-symphony-device:errors"] = std::move(errors);
    });
  }
}

void State::setGauge(const std::string& key, double value) {
  sink.setGauge(
      key,
      value,
      // adds the label deviceID = {deviceID}
      "deviceID",
      device);
}

} // namespace devices
} // namespace devmand
