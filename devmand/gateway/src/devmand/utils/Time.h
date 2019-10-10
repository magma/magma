// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <chrono>
#include <ctime>
#include <iomanip>
#include <ostream>

namespace devmand {

namespace utils {

using Clock = std::chrono::steady_clock;
using SteadyTimePoint = std::chrono::steady_clock::time_point;
using SystemTimePoint = std::chrono::system_clock::time_point;
using TimePoint = SteadyTimePoint;

namespace Time {

extern TimePoint now();

extern TimePoint from(struct timeval tv);

extern std::string toString(const TimePoint& time);

} // namespace Time

} // namespace utils

} // namespace devmand

extern std::ostream& operator<<(
    std::ostream& stream,
    const devmand::utils::TimePoint& time);
