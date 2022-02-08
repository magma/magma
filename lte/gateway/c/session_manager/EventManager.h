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

#include "event2/event.h"

#include <functional>
#include <iostream>
#include <queue>
#include <thread>

namespace magma {

// Libevent documentation helpful:
// http://www.wangafu.net/~nickm/libevent-book/Ref3_eventloop.html

// EventManager will keep a FIFO of events to be processed
// by arbitrary functions attached to the events.
//
// EventManager contains a thread that will wake when one or more
// tasks exists to be accomplished (signaled by new_event_) and
// then process the full FIFO of event actions' functions, then
// sleep again on new_event_;
//
// For now we will assume that clean-up / halt is not necessary.
//
// We will guarantee that events are processed in the order they
// are placed into the events_ container.
class EventManager {
 public:
  EventManager(struct event_base* base);
  void AddEvent(void (*action)(evutil_socket_t sig, short events,
                               void* user_data),
                void* callback_arg);
  void AddEventWithDelay(std::function<void> action, uint16_t delay_seconds);
  void Terminate();

 private:
  struct event_base* base_;
  std::unique_ptr<std::thread> processing_thread_;

  void Dispatcher();
};

EventManager::EventManager(struct event_base* base) {
  base_ = base;

  // Create thread to act as context for libevent dispatch.
  processing_thread_ =
      std::make_unique<std::thread>(&EventManager::Dispatcher, this);
}

// Blocking call waits for processing_thread_ join.
void EventManager::Terminate() {
  std::cout << "In Terminate! dumping the current event base entries"
            << std::endl;
  event_base_dump_events(base_, stdout);
  std::cout << "going to event_base_loopexit" << std::endl;
  assert(!event_base_loopexit(base_, nullptr));
  std::cout << "going to join processing_thread_" << std::endl;
  processing_thread_->join();
  std::cout << "done with Terminate!" << std::endl;
}

void EventManager::AddEvent(void (*action)(evutil_socket_t sig, short events,
                                           void* user_data),
                            void* callback_arg) {
  std::cout << "In AddEvent! going to make a new event: event_new" << std::endl;
  // TODO: I think this event needs to be somewhere
  struct event* new_event = event_new(base_, -1, 0, action, callback_arg);

  const timeval zero_tv{};
  std::cout << "adding the new event!" << std::endl;
  event_add(new_event, &zero_tv);
  event_base_dump_events(base_, stdout);
}

// Blocking function call that continuously executes libevent
// callbacks in this context until / unless libevent exit is
// triggered.
void EventManager::Dispatcher() {
  std::cout << "In Dispatcher! going to call a blocking func: event_base_loop!"
            << std::endl;
  int retval = event_base_loop(base_, EVLOOP_NO_EXIT_ON_EMPTY);
  if (retval != 0) {
    std::cout << "LOG ANGRY MESSAGE: " << retval << std::endl;
  }
  // log SO BAD WHY DID THIS HAPPEN?!@? WHY
}

}  // namespace magma
