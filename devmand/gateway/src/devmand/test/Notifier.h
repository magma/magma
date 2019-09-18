// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

namespace devmand {

#define Notify(_notifier)                 \
  WillOnce(testing::Invoke([&]() -> int { \
    _notifier.notify();                   \
    return 1;                             \
  }))

class Notifier final {
 public:
  Notifier() = default;
  ~Notifier() = default;
  Notifier(const Notifier&) = delete;
  Notifier& operator=(const Notifier&) = delete;
  Notifier(Notifier&&) = delete;
  Notifier& operator=(Notifier&&) = delete;

 public:
  void wait() {
    std::unique_lock<std::mutex> lk(m);
    cv.wait(lk, [this] { return ready; });
    ready = false;
  }

  void notify() {
    {
      std::lock_guard<std::mutex> lk(m);
      ready = true;
    }
    cv.notify_one();
  }

 private:
  std::mutex m;
  std::condition_variable cv;
  bool ready{false};
};

} // namespace devmand
