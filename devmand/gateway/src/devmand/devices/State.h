// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <list>
#include <memory>
#include <vector>

#include <folly/dynamic.h>
#include <folly/futures/Future.h>

#include <devmand/ErrorHandler.h>
#include <devmand/ErrorQueue.h>
#include <devmand/devices/Id.h>
#include <devmand/utils/LifetimeTracker.h>

namespace devmand {

class MetricSink;

namespace devices {

class State final : public std::enable_shared_from_this<State>,
                    public utils::LifetimeTracker<State> {
 private:
  State(MetricSink& sink_, const Id& device_);

 public:
  State() = delete;
  virtual ~State() = default;
  State(const State&) = delete;
  State& operator=(const State&) = delete;
  State(State&&) = delete;
  State& operator=(State&&) = delete;

  static std::shared_ptr<State> make(MetricSink& sink, const Id& device_);

 public:
  void update(std::function<void(folly::dynamic&)> func);
  void addRequest(folly::Future<folly::Unit> future);
  void setStatus(bool systemIsUp);
  void setErrors();
  void addError(std::string&& error);
  void setGauge(const std::string& key, double value);

  // Adds a callback to be executed on collect.
  void addFinally(std::function<void()>&& f);

  // NOTE a state object that is never collected will be a leak.
  folly::Future<folly::dynamic> collect();

  // clears requests and finalies to break circle.
  void clear();

 private:
  folly::dynamic& getFbcPlatformDevice(
      const std::string& key,
      folly::dynamic& unlockedState);

 private:
  // A link to the sink.
  MetricSink& sink;

  // The id of the device which created this state.
  Id device;

  // The state of an object formated according to the yang models supported.
  folly::Synchronized<folly::dynamic> state;

  // This is a queue of errors occuring on this system.
  ErrorQueue errorQueue;

  // This is a list of functions to execute after collect as returned. Most
  // likely these will operate on the state object by capturing it.
  std::list<std::function<void()>> finals;

  // This is a vector of futures which will be collected for the final
  // coalescing of state.
  std::vector<folly::Future<folly::Unit>> requests;
};

} // namespace devices
} // namespace devmand
