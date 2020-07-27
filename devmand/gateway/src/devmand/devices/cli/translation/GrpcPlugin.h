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

#include <devmand/channels/cli/plugin/protocpp/PluginRegistration.pb.h>
#include <devmand/channels/cli/plugin/protocpp/ReaderPlugin.grpc.pb.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <grpc++/grpc++.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace devmand::channels::cli;

class GrpcPlugin : public Plugin {
 private:
  shared_ptr<grpc::Channel> channel;
  const string id;
  shared_ptr<Executor> executor;
  devmand::channels::cli::plugin::CapabilitiesResponse capabilities;

  GrpcPlugin(
      shared_ptr<grpc::Channel> channel,
      const string id,
      shared_ptr<Executor> executor,
      devmand::channels::cli::plugin::CapabilitiesResponse capabilities);

 public:
  static shared_ptr<GrpcPlugin> create(
      shared_ptr<grpc::Channel> channel,
      const string id,
      shared_ptr<Executor> executor);

  DeviceType getDeviceType() const override;

  void provideReaders(ReaderRegistryBuilder& registry) const override;

  void provideWriters(WriterRegistryBuilder& registry) const override;

  Optional<CliFlavourParameters> getCliFlavourParameters() const;
};

} // namespace cli
} // namespace devices
} // namespace devmand
