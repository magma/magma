// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/error/ErrorQueue.h>

namespace devmand {

void ErrorQueue::add(std::string&& error) {
  errors.withWLock([this, &error](auto& queue) {
    queue.emplace_back(std::forward<std::string>(error));
    // on max size, discard oldest error
    if (queue.size() > maxSize) {
      queue.pop_front();
    }
  });
}

folly::dynamic ErrorQueue::get() {
  auto ret = folly::dynamic::array();
  // NOTE: this is a shared read lock but if you modify "get()" to clear
  // the errors, you'll need a write lock.
  errors.withRLock([this, &ret](auto& queue) {
    for (auto& error : queue) {
      ret.push_back(error);
    }
  });
  return ret;
}

} // namespace devmand
