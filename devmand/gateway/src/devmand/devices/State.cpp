// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <numeric>

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
        for (auto& f : s->finals) {
          f();
        }

        s->state.withRLock([&s](auto& unlockedState) {
          auto status = YangUtils::lookup(
              unlockedState, "fbc-symphony-device:system/status");
          if (status != nullptr and status.isString()) {
            s->setGauge("device.status", status.asString() == "UP" ? 1 : 0);
          } else {
            s->setGauge("device.status", 0);
          }
        });

        auto averageRequestDuration = getAverageRequestDuration(reqs).count();
        s->setGauge("device.request.duration.avg", averageRequestDuration);

        LOG(INFO) << s->device << " average request duration was "
                  << averageRequestDuration << " usec";
        s->clear();
        return std::move(*(s->state.wlock()));
      });
}

void State::clear() {
  finals.clear();
  requests.clear();
}

void State::update(std::function<void(folly::dynamic&)> func) {
  state.withWLock([&func](auto& unlockedState) { func(unlockedState); });
}

void State::addFinally(std::function<void()>&& f) {
  finals.push_back(f);
}

void State::addError(std::string&& error) {
  errorQueue.add(std::forward<std::string>(error));
}

void State::addRequest(folly::Future<folly::Unit> future) {
  auto request = std::make_shared<Request>();

  // In theory this end could already have occured but that's fine this isn't
  // very precise anyways.
  request->start = utils::Time::now();

  // Here we capture the state object as a shared ref and store it in the
  // request. This way if the user doesn't hold the ref. we have ensured someone
  // did.
  request->state = shared_from_this();

  requests.push_back(
      std::move(future)
          .thenValue([request](auto) {
            request->end = utils::Time::now();
            return request;
          })
          .thenError(
              folly::tag_t<std::exception>{},
              [request](std::exception const& e) {
                request->end = utils::Time::now();
                request->isError = true;

                LOG(ERROR) << "Caught exception from future: " << e.what();
                request->state->errorQueue.add(folly::sformat(
                    "Caught exception from future: {}", e.what()));
                return request;
              })
          .thenError([request](folly::exception_wrapper) {
            request->end = utils::Time::now();
            request->isError = true;

            LOG(ERROR) << "Caught unknown exception from future";
            request->state->errorQueue.add(
                "Caught unknown exception from future");
            return request;
          }));
}

std::chrono::microseconds State::getAverageRequestDuration(
    std::vector<std::shared_ptr<Request>> reqs) {
  return reqs.empty()
      ? std::chrono::microseconds(0)
      : std::accumulate(
            reqs.begin(),
            reqs.end(),
            std::chrono::microseconds(0),
            [](const std::chrono::microseconds& a,
               const std::shared_ptr<Request>& b) {
              return std::chrono::duration_cast<std::chrono::microseconds>(
                  a + (b->end - b->start));
            }) /
          static_cast<long int>(reqs.size());
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
}

void State::setErrors() {
  folly::dynamic errors = errorQueue.get();
  if (not errors.empty()) {
    state.withWLock([&errors](auto& lockedState) {
      lockedState["fbc-symphony-device:errors"] = std::move(errors);
    });
  }
}

} // namespace devices
} // namespace devmand
