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

#include <devmand/MetricSink.h>
#include <devmand/devices/Id.h>
#include <devmand/error/ErrorHandler.h>
#include <devmand/error/ErrorQueue.h>
#include <devmand/utils/LifetimeTracker.h>
#include <devmand/utils/Time.h>

namespace devmand {

namespace devices {

class Datastore;

struct Request {
  std::shared_ptr<Datastore> datastore{nullptr};
  utils::TimePoint start{};
  utils::TimePoint end{};
  bool isError{false};
};

class Datastore final : public std::enable_shared_from_this<Datastore>,
                        public utils::LifetimeTracker<Datastore> {
 private:
  Datastore(MetricSink& sink_, const Id& device_);

 public:
  Datastore() = delete;
  virtual ~Datastore() = default;
  Datastore(const Datastore&) = delete;
  Datastore& operator=(const Datastore&) = delete;
  Datastore(Datastore&&) = delete;
  Datastore& operator=(Datastore&&) = delete;

  static std::shared_ptr<Datastore> make(MetricSink& sink, const Id& device_);

 public:
  void update(std::function<void(folly::dynamic&)> func);
  void addRequest(folly::Future<folly::Unit> future);
  void setStatus(bool systemIsUp);
  void setErrors();
  void addError(std::string&& error);

  template <class T>
  void setGauge(const std::string& key, T value) {
    sink.setGauge(
        key,
        folly::to<double>(value),
        // adds the label deviceID = {deviceID}
        "deviceID",
        device);
  }

  void setGauge(const std::string& key, long unsigned int value);

  // Adds a callback to be executed on collect.
  void addFinally(std::function<void()>&& f);

  // NOTE a datastore object that is never collected will be a leak.
  folly::Future<folly::dynamic> collect();

  // clears requests and finalies to break circle.
  void clear();

 private:
  folly::dynamic& getFbcPlatformDevice(
      const std::string& key,
      folly::dynamic& unlockedDatastore);

  static std::chrono::microseconds getAverageRequestDuration(
      std::vector<std::shared_ptr<Request>> reqs);

 private:
  // A link to the sink.
  MetricSink& sink;

  // The id of the device which created this state.
  Id device;

  // The state of an object formated according to the yang models supported.
  folly::Synchronized<folly::dynamic> datastore;

  // This is a queue of errors occuring on this system.
  ErrorQueue errorQueue;

  // This is a list of functions to execute after collect as returned. Most
  // likely these will operate on the state object by capturing it.
  std::list<std::function<void()>> finals;

  // This is a vector of futures which will be collected for the final
  // coalescing of state.
  std::vector<folly::Future<std::shared_ptr<Request>>> requests;
};

} // namespace devices
} // namespace devmand
