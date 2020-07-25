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

#include <devmand/channels/cli/plugin/protocpp/WriterPlugin.grpc.pb.h>
#include <devmand/devices/cli/translation/GrpcCliHandler.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <grpc++/grpc++.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace folly;
using namespace devmand::channels::cli;

class GrpcWriter : public Writer, public GrpcCliHandler {
 private:
  unique_ptr<devmand::channels::cli::plugin::WriterPlugin::Stub> stub_;

 public:
  GrpcWriter(
      shared_ptr<grpc::Channel> channel,
      const string id,
      shared_ptr<Executor> executor);

  Future<Unit> create(const Path& path, dynamic cfg, const DeviceAccess& device)
      const override;

  Future<Unit> update(
      const Path& path,
      dynamic before,
      dynamic after,
      const DeviceAccess& device) const override;

  Future<Unit> remove(
      const Path& path,
      dynamic before,
      const DeviceAccess& device) const override;
};

} // namespace cli
} // namespace devices
} // namespace devmand
