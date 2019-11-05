// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/utils/Time.h>

#include <sstream>

namespace devmand {
namespace utils {
namespace Time {

TimePoint now() {
  return Clock::now();
}

TimePoint from(struct timeval tv) {
  return TimePoint{std::chrono::seconds{tv.tv_sec} +
                   std::chrono::microseconds{tv.tv_usec}};
}

std::string toString(const TimePoint& time) {
  std::stringstream stream;
  stream << time;
  return stream.str();
}

} // namespace Time
} // namespace utils
} // namespace devmand

std::ostream& operator<<(
    std::ostream& stream,
    const devmand::utils::TimePoint& time) {
  std::time_t epoch =
      std::chrono::duration_cast<std::chrono::seconds>(time.time_since_epoch())
          .count();
  return stream << std::put_time(std::localtime(&epoch), "%c");
}
