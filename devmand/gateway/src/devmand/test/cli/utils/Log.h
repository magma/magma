// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG

#include <magma_logging.h>

namespace devmand {
namespace test {
namespace utils {
namespace log {

using namespace std;

extern void initLog(uint32_t verbosity = MDEBUG);

} // namespace log
} // namespace utils
} // namespace test
} // namespace devmand
