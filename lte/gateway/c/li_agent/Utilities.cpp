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

#include <chrono>

#include "Utilities.h"
#include "includes/MagmaService.h"
#include "includes/MConfigLoader.h"
#include "magma_logging_init.h"

namespace magma {
namespace lte {

magma::mconfig::LIAgentD get_default_mconfig() {
  magma::mconfig::LIAgentD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  return mconfig;
}

magma::mconfig::LIAgentD load_mconfig() {
  magma::mconfig::LIAgentD mconfig;
  magma::MConfigLoader loader;
  if (!loader.load_service_mconfig(LIAGENTD, &mconfig)) {
    MLOG(MERROR) << "Unable to load mconfig for liagentd, using default";
    return get_default_mconfig();
  }
  return mconfig;
}

uint64_t get_time_in_sec_since_epoch() {
  auto now = std::chrono::system_clock::now();
  return std::chrono::duration_cast<std::chrono::seconds>(
             now.time_since_epoch())
      .count();
}

uint64_t time_difference_from_now(const uint64_t timestamp) {
  const auto now = get_time_in_sec_since_epoch();
  return (now - timestamp);
}

}  // namespace lte
}  // namespace magma
