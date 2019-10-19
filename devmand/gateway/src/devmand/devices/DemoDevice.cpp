// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/DemoDevice.h>
#include <devmand/models/wifi/Model.h>

namespace devmand {
namespace devices {

std::unique_ptr<devices::Device> DemoDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return std::make_unique<devices::DemoDevice>(app, deviceConfig.id);
}

DemoDevice::DemoDevice(Application& application, const Id& id_)
    : Device(application, id_) {}

std::shared_ptr<State> DemoDevice::getState() {
  auto state = State::make(*reinterpret_cast<MetricSink*>(&app), getId());
  state->update([](auto& lockedState) {
    lockedState = std::move(DemoDevice::getDemoState());
  });
  return state;
}

folly::dynamic DemoDevice::getDemoState() {
  folly::dynamic data = folly::dynamic::object;
  devmand::models::wifi::Model::init(data);

  // ##########################################################################
  auto& papRoot = data["openconfig-ap-manager:provision-aps"];
  auto& paps = papRoot["provision-ap"];
  folly::dynamic& pap = paps[0];
  pap["mac"] = "00:11:22:33:44:55";
  auto& papc = pap["config"];
  papc["mac"] = "00:11:22:33:44:55";
  papc["hostname"] = "faceboook.com";
  papc["country-code"] = "US";
  auto& papt = pap["state"];
  papt["mac"] = "00:11:22:33:44:55";
  papt["hostname"] = "faceboook.com";
  papt["country-code"] = "US";

  auto& japRoot = data["openconfig-ap-manager:joined-aps"];
  auto& japs = japRoot["joined-ap"];
  folly::dynamic& jap = japs[0];
  jap["hostname"] = "facebook.com";
  auto& japt = jap["state"];
  japt["mac"] = "00:11:22:33:44:55";
  japt["hostname"] = "facebook.com";
  japt["opstate"] = "openconfig-wifi-types:UP";
  japt["uptime"] = "0";
  japt["enabled"] = true;
  japt["serial"] = "cherrios";
  japt["model"] = "model";
  japt["software-version"] = "idk";
  japt["ipv4"] = "1.1.1.1";
  japt["ipv6"] = "11::11";
  japt["power-source"] = "PLUG";

  // ##########################################################################
  auto& system = data["ietf-system:system"] = folly::dynamic::object;
  system["hostname"] = "demo";
  system["contact"] = "fb@fb.com";
  system["location"] = "Boston Mass.";

  // ##########################################################################
  auto& interfaces = data["openconfig-interfaces:interfaces"] =
      folly::dynamic::object;
  interfaces["interface"] = folly::dynamic::array;

  folly::dynamic int0 = folly::dynamic::object;
  auto& stateInt0 = int0["state"] = folly::dynamic::object;
  stateInt0["ifindex"] = 0;
  stateInt0["name"] = "eth0";
  stateInt0["oper-status"] = "UP";
  interfaces["interface"].push_back(int0);

  folly::dynamic int1 = folly::dynamic::object;
  auto& stateInt1 = int1["state"] = folly::dynamic::object;
  stateInt1["ifindex"] = 1;
  stateInt1["name"] = "eth0";
  stateInt1["oper-status"] = "DOWN";
  interfaces["interface"].push_back(int1);

  return std::move(data);
}

} // namespace devices
} // namespace devmand
