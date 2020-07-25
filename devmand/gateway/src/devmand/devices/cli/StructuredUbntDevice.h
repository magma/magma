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

#include <devmand/Application.h>
#include <devmand/channels/cli/Channel.h>
#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/ReadCachingCli.h>
#include <devmand/channels/cli/TreeCache.h>
#include <devmand/channels/cli/datastore/Datastore.h>
#include <devmand/devices/Device.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;

class StructuredUbntDevice : public Device {
 public:
  StructuredUbntDevice(
      Application& application,
      const Id _id,
      bool readonly_,
      const std::shared_ptr<Channel> _channel,
      const std::shared_ptr<ModelRegistry> mreg,
      std::unique_ptr<ReaderRegistry>&& _rReg,
      std::unique_ptr<WriterRegistry>&& _wReg,
      const std::shared_ptr<CliCache> _cmdCache,
      const std::shared_ptr<TreeCache> _treeCache);
  StructuredUbntDevice() = delete;
  virtual ~StructuredUbntDevice() = default;
  StructuredUbntDevice(const StructuredUbntDevice&) = delete;
  StructuredUbntDevice& operator=(const StructuredUbntDevice&) = delete;
  StructuredUbntDevice(StructuredUbntDevice&&) = delete;
  StructuredUbntDevice& operator=(StructuredUbntDevice&&) = delete;

  static std::unique_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

  // visible for testing
  static std::unique_ptr<devices::Device> createDeviceWithEngine(
      Application& app,
      const cartography::DeviceConfig& deviceConfig,
      Engine& engine);

 public:
  std::shared_ptr<Datastore> getOperationalDatastore() override;

 protected:
  void setIntendedDatastore(const folly::dynamic& config) override;

 private:
  std::shared_ptr<Channel> channel;
  std::shared_ptr<CliCache> cmdCache;
  std::shared_ptr<TreeCache> treeCache;
  std::shared_ptr<ModelRegistry> mreg;
  std::unique_ptr<ReaderRegistry> rReg;
  std::unique_ptr<WriterRegistry> wReg;
  std::unique_ptr<devmand::channels::cli::datastore::Datastore> configCache;
  vector<DiffPath> diffPaths;

  void reconcile(DeviceAccess& access);
};

} // namespace cli
} // namespace devices
} // namespace devmand
