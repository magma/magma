/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include <gflags/gflags.h>
#include <glog/logging.h>
#include <stdint.h>
#include <ostream>

#include "orc8r/gateway/c/common/logging/magma_logging.hpp"

namespace magma {

// set_verbosity sets the global logging verbosity. The higher the verbosity,
// the more is logged
static void set_verbosity(uint32_t verbosity) {
  VLOG(0) << "Setting verbosity to " << verbosity;
  FLAGS_v = verbosity;
}

// get_verbosity gets the the global logging verbosity
static google::int32 get_verbosity() {
  (void)get_verbosity;  // casting to void to suppress unused reference warning
  return FLAGS_v;
}

// init_logging initializes glog, sets logging to use std::err, and sets the
// initial verbosity
static void init_logging(const char* service_name) {
  google::InitGoogleLogging(service_name);
  // log to stderr to automatically log to syslog with systemd
  FLAGS_logtostderr = 1;
}
}  // namespace magma