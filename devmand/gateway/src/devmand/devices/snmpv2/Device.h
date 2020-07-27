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

#include <devmand/channels/snmp/Channel.h>
#include <devmand/channels/snmp/IfMib.h>
#include <devmand/devices/ping/Device.h>

namespace devmand {
namespace devices {
namespace snmpv2 {

class Device : public ping::Device {
 public:
  Device(
      Application& application,
      const Id& id,
      bool readonly_,
      const channels::snmp::Peer& peer,
      const channels::snmp::Community& community,
      const channels::snmp::Version& version,
      const std::string& passphrase = "",
      const std::string& securityName = "",
      const channels::snmp::SecurityLevel& securityLevel = "",
      oid proto[] = {});
  Device() = delete;
  virtual ~Device() = default;
  Device(const Device&) = delete;
  Device& operator=(const Device&) = delete;
  Device(Device&&) = delete;
  Device& operator=(Device&&) = delete;

  static std::shared_ptr<devices::Device> createDevice(
      Application& app,
      const cartography::DeviceConfig& deviceConfig);

 public:
  std::shared_ptr<Datastore> getOperationalDatastore() override;

 private:
  folly::Future<folly::Unit> addToStateWithInterfaceIndices(
      std::shared_ptr<Datastore> state,
      const devmand::channels::snmp::InterfaceIndicies& interfaceIndices);

 protected:
  void setIntendedDatastore(const folly::dynamic& config) override {
    (void)config;
    LOG(ERROR) << "set config on unconfigurable device";
  }

 protected:
  channels::snmp::Channel snmpChannel;
};

} // namespace snmpv2
} // namespace devices
} // namespace devmand
