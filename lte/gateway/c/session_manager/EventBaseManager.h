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
/**
 * This is a helper that wraps C style API objects that need to be deleted with
 * a smart pointer.
 * Based off of EnvoyProxy's implementation
 */
template <class T, void (*deleter)(T*)>
class CSmartPtr : public std::unique_ptr<T, void (*)(T*)> {
 public:
  CSmartPtr() : std::unique_ptr<T, void (*)(T*)>(nullptr, deleter) {}
  CSmartPtr(T* object) : std::unique_ptr<T, void (*)(T*)>(object, deleter) {}
};
using EventBasePtr = CSmartPtr<event_base, event_base_free>;

// Libevent documentation helpful:
// http://www.wangafu.net/~nickm/libevent-book/Ref3_eventloop.html

// EventBaseManager will keep a FIFO of events to be processed
// by arbitrary functions attached to the events.

// EventBaseManager contains a thread that will wake when one or more
// tasks exists to be accomplished (signaled by new_event_) and
// then process the full FIFO of event actions' functions, then
// sleep again on new_event_;

// For now we will assume that clean-up / halt is not necessary.

// We will guarantee that events are processed in the order they
// are placed into the events_ container.
class EventBaseManager {
 public:
  EventBaseManager(EventBasePtr& base);
  void Terminate();

 private:
  EventBasePtr& base_;
  std::unique_ptr<std::thread> processing_thread_;

  void Dispatcher();
};

EventBaseManager::EventBaseManager(EventBasePtr& base) : base_(base) {
  // Create thread to act as context for libevent dispatch.
  processing_thread_ =
      std::make_unique<std::thread>(&EventBaseManager::Dispatcher, this);
}

// Blocking call waits for processing_thread_ join.
void EventBaseManager::Terminate() {
  std::cout << "In Terminate!" << std::endl;
  assert(!event_base_loopexit(base_.get(), nullptr));
  processing_thread_->join();
  std::cout << "done with Terminate! Dispatcher thread has been joined!"
            << std::endl;
}

// Blocking function call that continuously executes libevent
// callbacks in this context until / unless libevent exit is
// triggered.
void EventBaseManager::Dispatcher() {
  std::cout << "In Dispatcher! going to call a blocking func: event_base_loop!"
            << std::endl;
  int retval = event_base_loop(base_.get(), EVLOOP_NO_EXIT_ON_EMPTY);
  if (retval != 0) {
    std::cout << "LOG ANGRY MESSAGE: " << retval << std::endl;
  }
  // log SO BAD WHY DID THIS HAPPEN?!@? WHY
}

}  // namespace magma
