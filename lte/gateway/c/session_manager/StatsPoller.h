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
#include <stdint.h>
#include <atomic>
#include <memory>

#include "LocalEnforcer.h"

namespace magma {
class LocalEnforcer;

class StatsPoller {
 public:
  /**
   * start_loop is the main function to call to initiate a loop. Based
   * on the given loop interval length, this function will poll stats from
   * Pipelined every loop_interval_seconds
   */
  void start_loop(std::shared_ptr<LocalEnforcer> local_enforcer,
                  uint32_t loop_interval_seconds);
};
}  // namespace magma
