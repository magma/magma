// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/engine/Engine.h>
#include "devmand/test/cli/utils/Log.h"

namespace devmand {
namespace test {
namespace utils {
namespace log {

using namespace std;
using namespace devmand::channels::cli;

atomic_bool loggingInitialized(false);

void initLog() {
  if (loggingInitialized.load()) {
    return;
  }
  Engine::initLogging(MDEBUG, true);
  loggingInitialized.store(true);
  MLOG(MDEBUG) << "Logging for test initialized";
}

} // namespace log
} // namespace utils
} // namespace test
} // namespace devmand
