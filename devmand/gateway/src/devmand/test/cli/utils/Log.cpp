/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
    ::magma::set_verbosity(verbosity);
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
