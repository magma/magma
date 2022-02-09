/**
 * Copyright 2022 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limit
 */
#pragma once

#include <event2/event.h>
#include "event2/event_struct.h"
#include "lte/gateway/c/session_manager/EventBaseManager.h"

#include <functional>
#include <iostream>
#include <queue>
#include <thread>

// Modified from EnvoyProxy/Envoy's wrapper around libevent

namespace magma {
class SchedulableCallback {
 public:
  static SchedulableCallback* MakeSchedulableCallback(EventBasePtr& base,
                                                      std::function<void()> cb);
  void scheduleCallbackSoon();
  void scheduleCallbackWithDelay(uint16_t delay_seconds);

 private:
  SchedulableCallback(EventBasePtr& base, std::function<void()> cb);
  std::function<void()> cb_;
  event* raw_event_;
};

SchedulableCallback* SchedulableCallback::MakeSchedulableCallback(
    EventBasePtr& base, std::function<void()> cb) {
  return new SchedulableCallback(base, cb);
}

SchedulableCallback::SchedulableCallback(EventBasePtr& base,
                                         std::function<void()> cb)
    : cb_(cb) {
  std::cout << "In SchedulableCallback constructor!" << std::endl;
  raw_event_ = event_new(
      base.get(), -1, 0,
      [](evutil_socket_t, short, void* arg) -> void {
        SchedulableCallback* cb = static_cast<SchedulableCallback*>(arg);
        cb->cb_();

        // Clean up everything
        event_del(cb->raw_event_);
        event_free(cb->raw_event_);
        free(cb);
      },
      this);
  std::cout << "Exiting SchedulableCallback constructor!" << std::endl;
}

void SchedulableCallback::scheduleCallbackSoon() {
  std::cout << "Calling scheduleCallbackSoon!" << std::endl;
  // libevent computes the list of timers to move to the work list after polling
  // for fd events, but iteration through the work list starts. Zero delay
  // timers added while iterating through the work list execute on the next
  // iteration of the event loop.
  const timeval zero_tv{};
  event_add(raw_event_, &zero_tv);
  std::cout << "Exiting scheduleCallbackSoon!" << std::endl;
}

void SchedulableCallback::scheduleCallbackWithDelay(uint16_t delay_seconds) {
  std::cout << "Calling scheduleCallbackWithDelay!" << std::endl;
  // libevent computes the list of timers to move to the work list after polling
  // for fd events, but iteration through the work list starts. Zero delay
  // timers added while iterating through the work list execute on the next
  // iteration of the event loop.
  const timeval tv{
      .tv_sec = delay_seconds,
  };
  event_add(raw_event_, &tv);
  std::cout << "Exiting scheduleCallbackWithDelay!" << std::endl;
}

}  // namespace magma
