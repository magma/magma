// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include "devmand/test/cli/utils/Log.h"
#include <devmand/channels/cli/engine/Engine.h>

namespace devmand {
namespace test {
namespace utils {
namespace log {

using namespace std;
using namespace devmand::channels::cli;

atomic_bool loggingInitialized(false);

void initLog(uint32_t verbosity) {
  if (loggingInitialized.load()) {
    magma::set_verbosity(verbosity);
    return;
  }
  Engine::initLogging(verbosity, true);
  loggingInitialized.store(true);
  MLOG(MDEBUG) << "Logging for test initialized";
}

} // namespace log
} // namespace utils
} // namespace test
} // namespace devmand
