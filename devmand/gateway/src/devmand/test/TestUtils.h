// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <chrono>
#include <thread>

#define EXPECT_BECOMES_TRUE(exp)                                     \
  do {                                                               \
    constexpr std::chrono::seconds _maxExpectsWait{10};              \
    auto _expectsWait = _maxExpectsWait;                             \
    while ((not(exp)) and _expectsWait != std::chrono::seconds(0)) { \
      _expectsWait -= std::chrono::seconds(1);                       \
      std::this_thread::sleep_for(std::chrono::seconds(1));          \
    }                                                                \
    EXPECT_TRUE(exp);                                                \
  } while (false);
