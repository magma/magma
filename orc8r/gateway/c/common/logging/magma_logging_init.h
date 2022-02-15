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

#include <stdint.h>

#include "orc8r/gateway/c/common/logging/magma_logging.h"

#include <glog/logging.h>

namespace magma {

// set_verbosity sets the global logging verbosity. The higher the verbosity,
// the more is logged
void set_verbosity(uint32_t verbosity);

// get_verbosity gets the the global logging verbosity
google::int32 get_verbosity();

// init_logging initializes glog, sets logging to use std::err, and sets the
// initial verbosity
void init_logging(const char* service_name);
}  // namespace magma
