// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/State.h>

#include <devmand/Application.h>
#include <devmand/devices/Device.h>

namespace devmand {
namespace devices {

std::shared_ptr<State> State::make(Application& application, Device& device_) {
  return std::shared_ptr<State>(new State(application, device_));
}

State::State(Application& application, Device& device_)
    : app(application), device(device_), state(folly::dynamic::object) {}

folly::Future<folly::dynamic> State::collect() {
  // TODO how to capture errors collect?
  return folly::collect(std::move(requests))
      .thenValue([s = shared_from_this()](auto) {
        s->setErrors();
        s->setStatus(true);
        for (auto& f : s->finals) {
          f();
        }
        return s->state;
      });
}

folly::dynamic& State::update() {
  return state;
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

folly::dynamic& State::getFbcPlatformDevice(const std::string& key) {
  auto k = folly::sformat("fbc-symphony-device:{}", key);
  folly::dynamic* system = state.get_ptr(k);
  if (system == nullptr) {
    system = &(state[k] = folly::dynamic::object);
  }
  return *system;
}

void State::setStatus(bool systemIsUp) {
  folly::dynamic& system = getFbcPlatformDevice("system");
  system["status"] = systemIsUp ? "UP" : "DOWN";
  app.setGauge(folly::sformat("{}.status", device.getId()), systemIsUp ? 1 : 0);
}

void State::setErrors() {
  folly::dynamic errors = errorQueue.get();
  if (not errors.empty()) {
    state["fbc-symphony-device:errors"] = std::move(errors);
  }
}

} // namespace devices
} // namespace devmand
