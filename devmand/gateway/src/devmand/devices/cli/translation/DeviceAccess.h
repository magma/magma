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

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/Channel.h>
#include <folly/executors/CPUThreadPoolExecutor.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace devmand::channels::cli;

class DeviceAccess {
 public:
  DeviceAccess(
      shared_ptr<Cli> _cliChannel,
      string _deviceId,
      shared_ptr<Executor> _workerExecutor);
  ~DeviceAccess() = default;

 public:
  shared_ptr<Cli> cli() const;
  string id() const;
  shared_ptr<Executor> executor() const;

 private:
  shared_ptr<Cli> cliChannel;
  string deviceId;
  shared_ptr<Executor> workerExecutor;
};

} // namespace cli
} // namespace devices
} // namespace devmand
