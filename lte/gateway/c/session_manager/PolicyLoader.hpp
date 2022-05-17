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
#include <lte/protos/apn.pb.h>
#include <lte/protos/policydb.pb.h>
#include <stdint.h>
#include <atomic>
#include <functional>
#include <vector>

namespace magma {
namespace lte {
class PolicyRule;
}
}  // namespace magma

namespace magma {
using namespace lte;
/**
 * PolicyLoader is used to sync policies with Redis every so often
 */
class PolicyLoader {
 public:
  /**
   * start_loop is the main function to call to initiate a load loop. Based on
   * the given loop interval length, this function will load the policies from
   * redis, and call the processor callback.
   */
  void start_loop(std::function<void(std::vector<PolicyRule>)> processor,
                  uint32_t loop_interval_seconds);

  /**
   * Stop the config loop on the next loop
   */
  void stop();

 private:
  std::atomic<bool> is_running_;
};
}  // namespace magma
