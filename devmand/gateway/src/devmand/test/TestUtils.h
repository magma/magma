// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
