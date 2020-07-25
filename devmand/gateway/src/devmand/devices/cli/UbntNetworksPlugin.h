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

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/translation/BindingReaderRegistry.h>
#include <devmand/devices/cli/translation/BindingWriterRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;

class UbntNetworksPlugin : public Plugin {
 public:
  UbntNetworksPlugin(BindingContext& openconfigContext);

  DeviceType getDeviceType() const override;
  void provideReaders(ReaderRegistryBuilder& registry) const override;
  void provideWriters(WriterRegistryBuilder& registry) const override;

 private:
  BindingContext& openconfigContext;
};

} // namespace cli
} // namespace devices
} // namespace devmand
