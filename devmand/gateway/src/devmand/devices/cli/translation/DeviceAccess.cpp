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

#include <devmand/devices/cli/translation/DeviceAccess.h>

namespace devmand {
namespace devices {
namespace cli {

DeviceAccess::DeviceAccess(
    shared_ptr<Cli> _cliChannel,
    string _deviceId,
    shared_ptr<Executor> _workerExecutor)
    : cliChannel(_cliChannel),
      deviceId(_deviceId),
      workerExecutor(_workerExecutor) {}

shared_ptr<Cli> DeviceAccess::cli() const {
  return cliChannel;
}

shared_ptr<Executor> DeviceAccess::executor() const {
  return workerExecutor;
}

string DeviceAccess::id() const {
  return deviceId;
}

} // namespace cli
} // namespace devices
} // namespace devmand
