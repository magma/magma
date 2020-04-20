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

#include <functional>
#include <map>
#include <memory>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/devices/Device.h>

namespace devmand {

class Application;

namespace devices {

class Factory final {
 public:
  Factory(Application& application);
  Factory() = delete;
  ~Factory() = default;
  Factory(const Factory&) = delete;
  Factory& operator=(const Factory&) = delete;
  Factory(Factory&&) = delete;
  Factory& operator=(Factory&&) = delete;

 public:
  using PlatformBuilder = std::function<std::shared_ptr<devices::Device>(
      Application& application,
      const cartography::DeviceConfig& deviceConfig)>;

  std::shared_ptr<devices::Device> createDevice(
      const cartography::DeviceConfig& deviceConfig);

  void addPlatform(
      const std::string& platform,
      PlatformBuilder platformBuilder);

  void setDefaultPlatform(PlatformBuilder defaultPlatformBuilder_);

 private:
  Application& app;
  std::map<std::string, PlatformBuilder> platformBuilders;
  PlatformBuilder defaultPlatformBuilder{nullptr};
};

} // namespace devices
} // namespace devmand
