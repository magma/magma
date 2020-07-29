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

#include <devmand/devices/echo/Device.h>

namespace devmand {
namespace devices {
namespace echo {

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return std::make_unique<devices::echo::Device>(
      app, deviceConfig.id, deviceConfig.readonly);
}

Device::Device(Application& application, const Id& id_, bool readonly_)
    : devices::Device(application, id_, readonly_) {}

void Device::setIntendedDatastore(const folly::dynamic& config) {
  state = config;
}

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  auto stateCopy =
      Datastore::make(*reinterpret_cast<MetricSink*>(&app), getId());
  stateCopy->update([this](auto& lockedDatastore) { lockedDatastore = state; });
  return stateCopy;
}

} // namespace echo
} // namespace devices
} // namespace devmand
