// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cambium/Device.h>

#include <iostream>

#include <folly/dynamic.h>
#include <folly/json.h>

#include <devmand/Application.h>
#include <devmand/error/ErrorHandler.h>

namespace devmand {
namespace devices {
namespace cambium {

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  const auto& cambiumKv = channelConfigs.at("cambium").kvPairs;
  return std::make_unique<devices::cambium::Device>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
      folly::IPAddress(deviceConfig.ip),
      cambiumKv.at("client-mac"),
      cambiumKv.at("client-id"),
      cambiumKv.at("client-secret"));
}

Device::Device(
    Application& application,
    const Id& id_,
    bool readonly_,
    const folly::IPAddress& deviceIp_,
    const std::string& clientMac_,
    const std::string& clientId_,
    const std::string& clientSecret_)
    : devices::Device(application, id_, readonly_),
      channel(deviceIp_.str(), clientId_, clientSecret_),
      deviceIp(deviceIp_),
      clientMac(clientMac_) {
  lastUpdate = setupReturnData();
  setupOpenconfig(lastUpdate);
  setupOpenconfigInterfaces(lastUpdate);
  setupIpv4(lastUpdate, 0);
  connect();
}

Device::~Device() {
  // TODO disconnect();
}

void Device::connect() {
  auto retval = channel.setupChannel();
}

folly::dynamic Device::setupReturnData() {
  folly::dynamic data = folly::dynamic::object;
  data["ietf-system:system"] = folly::dynamic::object;

  data["ietf-system:system"]["contact"] = "";
  data["ietf-system:system"]["location"] = "INVALID";
  data["ietf-system:system"]["name"] = "INVALID";

  return data;
}

void Device::setupOpenconfig(folly::dynamic& data) {
  // TODO: This could also be scripted in for loops
  // TODO: Break up all array populations into separate functions.
  data["openconfig-wifi-mac:ssids"] = folly::dynamic::object;

  data["openconfig-wifi-mac:ssids"]["ssid"] = folly::dynamic::array;
  auto& ssidArray = data["openconfig-wifi-mac:ssids"]["ssid"];
  ssidArray.push_back(folly::dynamic::object);
  ssidArray[0]["config"] = folly::dynamic::object;
  ssidArray[0]["state"] = folly::dynamic::object;
  ssidArray[0]["bssids"] = folly::dynamic::object;
  ssidArray[0]["wmm"] = folly::dynamic::object;
  ssidArray[0]["dot11r"] = folly::dynamic::object;
  ssidArray[0]["dot11v"] = folly::dynamic::object;
  ssidArray[0]["clients"] = folly::dynamic::object;
  ssidArray[0]["dot1x-timers"] = folly::dynamic::object;
  ssidArray[0]["band-steering"] = folly::dynamic::object;

  ssidArray[0]["bssids"]["bssid"] = folly::dynamic::array;
  auto& bssidArray = ssidArray[0]["bssids"]["bssid"];
  bssidArray.push_back(folly::dynamic::object);
  bssidArray[0]["state"] = folly::dynamic::object;

  bssidArray[0]["state"]["counters"] = folly::dynamic::object;
  bssidArray[0]["state"]["counters"]["rx-data-dist"] = folly::dynamic::object;
  bssidArray[0]["state"]["counters"]["rx-data-wmm"] = folly::dynamic::object;
  bssidArray[0]["state"]["counters"]["rx-mcs"] = folly::dynamic::object;
  bssidArray[0]["state"]["counters"]["tx-data-dist"] = folly::dynamic::object;
  bssidArray[0]["state"]["counters"]["tx-data-wmm"] = folly::dynamic::object;
  bssidArray[0]["state"]["counters"]["tx-mcs"] = folly::dynamic::object;

  ssidArray[0]["wmm"]["config"] = folly::dynamic::object;
  ssidArray[0]["wmm"]["state"] = folly::dynamic::object;

  ssidArray[0]["dot11r"]["config"] = folly::dynamic::object;
  ssidArray[0]["dot11r"]["state"] = folly::dynamic::object;

  ssidArray[0]["dot11v"]["config"] = folly::dynamic::object;
  ssidArray[0]["dot11v"]["state"] = folly::dynamic::object;

  ssidArray[0]["clients"]["client"] = folly::dynamic::array;
  auto& clientArray = ssidArray[0]["clients"]["client"];
  clientArray.push_back(folly::dynamic::object);
  clientArray[0]["state"] = folly::dynamic::object;
  clientArray[0]["state"]["counters"] = folly::dynamic::object;
  clientArray[0]["client-rf"] = folly::dynamic::object;
  clientArray[0]["client-rf"]["state"] = folly::dynamic::object;
  clientArray[0]["client-capabilities"] = folly::dynamic::object;
  clientArray[0]["client-capabilities"]["state"] = folly::dynamic::object;
  clientArray[0]["dot11k-neighbors"] = folly::dynamic::object;
  clientArray[0]["dot11k-neighbors"]["state"] = folly::dynamic::object;
  clientArray[0]["client-connection"] = folly::dynamic::object;
  clientArray[0]["client-connection"]["state"] = folly::dynamic::object;

  ssidArray[0]["dot1x-timers"]["config"] = folly::dynamic::object;
  ssidArray[0]["dot1x-timers"]["state"] = folly::dynamic::object;

  ssidArray[0]["band-steering"]["config"] = folly::dynamic::object;
  ssidArray[0]["band-steering"]["state"] = folly::dynamic::object;
}

void Device::setupOpenconfigInterfaces(folly::dynamic& data) {
  data["openconfig-interfaces:interfaces"] = folly::dynamic::object;
  data["openconfig-interfaces:interfaces"]["interface"] = folly::dynamic::array;
  addOpenconfigInterface(data);
}

void Device::addOpenconfigInterface(folly::dynamic& data) {
  auto& interfaces = data["openconfig-interfaces:interfaces"]["interface"];
  int index = static_cast<int>(interfaces.size());
  interfaces.push_back(folly::dynamic::object);
  interfaces[index]["name"] = "tmp";
  interfaces[index]["config"] = folly::dynamic::object;
  interfaces[index]["state"] = folly::dynamic::object;
  interfaces[index]["state"]["counters"] = folly::dynamic::object;
  interfaces[index]["hold-time"] = folly::dynamic::object;
  interfaces[index]["subinterfaces"] = folly::dynamic::object;
  interfaces[index]["subinterfaces"]["subinterface"] = folly::dynamic::array;

  // TODO: Figure out how to better represent these identiyrefs and enumerations
  interfaces[index]["config"]["type"] = "iana-if-type:other";

  interfaces[index]["state"]["type"] = "iana-if-type:other";
  interfaces[index]["state"]["admin-status"] = "UP";
  interfaces[index]["state"]["oper-status"] = "UP";

  interfaces[index]["hold-time"]["config"] = folly::dynamic::object;
  interfaces[index]["hold-time"]["state"] = folly::dynamic::object;
}

void Device::setupIpv4(folly::dynamic& data, int interfaceNum) {
  // TODO: pull all lists' populations out into separate functions

  int index =
      static_cast<int>(data["openconfig-interfaces:interfaces"]["interface"]
                           [interfaceNum]["subinterfaces"]["subinterface"]
                               .size());

  data["openconfig-interfaces:interfaces"]["interface"][interfaceNum]
      ["subinterfaces"]["subinterface"]
          .push_back(folly::dynamic::object);
  data["openconfig-interfaces:interfaces"]["interface"][interfaceNum]
      ["subinterfaces"]["subinterface"][index]["openconfig-if-ip:ipv4"] =
          folly::dynamic::object;
  auto& subinterface =
      data["openconfig-interfaces:interfaces"]["interface"][interfaceNum]
          ["subinterfaces"]["subinterface"][index]["openconfig-if-ip:ipv4"];

  subinterface["addresses"] = folly::dynamic::object;
  subinterface["addresses"]["address"] = folly::dynamic::array;
  subinterface["proxy-arp"] = folly::dynamic::object;
  subinterface["neighbors"] = folly::dynamic::object;
  subinterface["neighbors"]["neighbor"] = folly::dynamic::array;
  subinterface["unnumbered"] = folly::dynamic::object;
  subinterface["config"] = folly::dynamic::object;
  subinterface["state"] = folly::dynamic::object;

  subinterface["neighbors"]["neighbor"].push_back(folly::dynamic::object);
  subinterface["neighbors"]["neighbor"][0]["config"] = folly::dynamic::object;
  subinterface["neighbors"]["neighbor"][0]["state"] = folly::dynamic::object;

  subinterface["addresses"]["address"].push_back(folly::dynamic::object);
  auto& address = subinterface["addresses"]["address"][0];
  address["config"] = folly::dynamic::object;
  address["state"] = folly::dynamic::object;
  address["vrrp"] = folly::dynamic::object;
  address["vrrp"]["vrrp-group"] = folly::dynamic::array;
  address["vrrp"]["vrrp-group"].push_back(folly::dynamic::object);
  address["vrrp"]["vrrp-group"][0]["config"] = folly::dynamic::object;
  address["vrrp"]["vrrp-group"][0]["state"] = folly::dynamic::object;
  auto& intTracking = address["vrrp"]["vrrp-group"][0]["interface-tracking"];
  intTracking = folly::dynamic::object;
  intTracking["config"] = folly::dynamic::object;
  intTracking["state"] = folly::dynamic::object;

  subinterface["proxy-arp"]["config"] = folly::dynamic::object;
  subinterface["proxy-arp"]["state"] = folly::dynamic::object;

  subinterface["unnumbered"]["config"] = folly::dynamic::object;
  subinterface["unnumbered"]["state"] = folly::dynamic::object;
  subinterface["unnumbered"]["interface-ref"] = folly::dynamic::object;
  subinterface["unnumbered"]["interface-ref"]["state"] = folly::dynamic::object;
  subinterface["unnumbered"]["interface-ref"]["config"] =
      folly::dynamic::object;

  subinterface["state"]["counters"] = folly::dynamic::object;
}

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  // TODO: Improve error handling
  folly::dynamic data = setupReturnData();
  folly::dynamic returnedData = channel.getDeviceInfo(clientMac);

  if (returnedData.isNull()) {
    auto state = Datastore::make(app, getId());
    state->update(
        [&data](auto& lockedState) { lockedState = std::move(data); });
    return state;
  }

  setupOpenconfig(data);
  setupOpenconfigInterfaces(data);
  setupIpv4(data, 0);

  folly::dynamic parsed = folly::parseJson(returnedData.asString());

  data["ietf-system:system"]["contact"] = "";
  data["ietf-system:system"]["location"] = parsed["data"][0]["site"];
  data["ietf-system:system"]["name"] = parsed["data"][0]["name"];

  // TODO: Check to see if the 1st one is "tmp". If it is, then override.
  // Otherwise, parse and find the right one to update. Normal IP
  // As of now, it is assumed that the first subinterface is for the gateway and
  // the second is for the VLAN.
  auto& interface = data["openconfig-interfaces:interfaces"]["interface"];
  interface[0]["name"] = parsed["data"][0]["name"];
  interface[0]["config"]["name"] = parsed["data"][0]["name"];
  interface[0]["subinterfaces"]["subinterface"][0]["openconfig-if-ip:ipv4"]
           ["addresses"]["address"][0]["config"]["ip"] =
               parsed["data"][0]["ip"];
  interface[0]["subinterfaces"]["subinterface"][0]["openconfig-if-ip:ipv4"]
           ["addresses"]["address"][0]["ip"] = parsed["data"][0]["ip"];

  // VLAN
  if (data["openconfig-interfaces:interfaces"]["interface"].size() < 2) {
    addOpenconfigInterface(data);
    setupIpv4(data, 1);
  }

  interface[1]["name"] = "VLAN_IP";
  interface[1]["config"]["name"] = "VLAN_IP";
  interface[1]["subinterfaces"]["subinterface"][0]["openconfig-if-ip:ipv4"]
           ["addresses"]["address"][0]["config"]["ip"] =
               parsed["data"][0]["config"]["variables"]["VLAN_1_IP"];
  interface[0]["subinterfaces"]["subinterface"][0]["openconfig-if-ip:ipv4"]
           ["addresses"]["address"][0]["ip"] =
               parsed["data"][0]["config"]["variables"]["VLAN_1_IP"];

  auto state = Datastore::make(app, getId());
  state->update([&data](auto& lockedState) { lockedState = std::move(data); });
  return state;
}

void Device::updateYang(
    const folly::dynamic& config,
    std::vector<std::string>& path,
    long unsigned int index,
    const folly::dynamic& original,
    folly::dynamic& updateJson) {
  if (index >= (path.size())) {
    // It shouldn't end on an array, and that is the only thing this catches.
    return;
  }

  switch (config.type()) {
    case folly::dynamic::ARRAY: {
      for (int x = 0; x < static_cast<int>(config.size()); ++x) {
        path[index] = folly::to<std::string>(x);
        updateYang(config[x], path, index + 1, original, updateJson);
      }
      break;
    }
    case folly::dynamic::OBJECT: {
      if (config.find(path[index]) != config.items().end()) {
        if (index + 1 == path.size()) {
          updateDevice(original, path, updateJson);
          return;
        } else {
          updateYang(
              config[path[index]], path, index + 1, original, updateJson);
        }
      }
      break;
    }
    case folly::dynamic::NULLT:
    case folly::dynamic::BOOL:
    case folly::dynamic::DOUBLE:
    case folly::dynamic::INT64:
    case folly::dynamic::STRING:
    default: { return; }
  }
}

void Device::setIntendedDatastore(const folly::dynamic& config) {
  // TODO: Break out successfully so we don't waste lots of for looping
  // TODO: Figure why we couldn't declare the vector in one line.
  folly::dynamic updateJson = folly::dynamic::object;
  updateJson["overrides"] = folly::dynamic::object;
  updateJson["overrides"]["vlans"] = folly::dynamic::array;
  std::vector<std::string> ssidsPath = {
      "openconfig-wifi-mac:ssids", "ssid", "0", "name"};
  std::vector<std::string> interfacesPath = {"openconfig-interfaces:interfaces",
                                             "interface",
                                             "0",
                                             "subinterfaces",
                                             "subinterface",
                                             "0",
                                             "openconfig-if-ip:ipv4",
                                             "addresses",
                                             "address",
                                             "0",
                                             "ip"};

  updateYang(config, ssidsPath, 0, config, updateJson);
  updateYang(config, interfacesPath, 0, config, updateJson);
  // TODO: when this is converted to a response class we may want to do
  // something with the response.
  channel.updateDevice(updateJson, clientMac);
  // LOG(ERROR) << "Output of setRunningDatastore is: " << output;
}

void Device::updateDevice(
    const folly::dynamic& yangIn,
    std::vector<std::string>& path,
    folly::dynamic& updateJson) {
  // This loops through the main keys in order to find one that we identify and
  // associate with a specific YANG model.
  // TODO: Save the last update and use that to know whether we should update
  // the AP or not.

  if (path[0] == std::string("openconfig-wifi-mac:ssids")) {
    updateJson["name"] =
        yangIn[path[0]][path[1]][folly::to<int>(path[2])]["name"];
  } else if (path[0] == std::string("openconfig-interfaces:interfaces")) {
    auto full_path = yangIn[path[0]][path[1]][folly::to<int>(path[2])][path[3]]
                           [path[4]][folly::to<int>(path[5])][path[6]][path[7]]
                           [path[8]][folly::to<int>(path[9])][path[10]];
    if (yangIn[path[0]][path[1]][folly::to<int>(path[2])]["name"] ==
        std::string("VLAN_IP")) {
      auto index = updateJson["overrides"]["vlans"].size();
      // updateJson["overrides"]["auto_set"]["network"] = false;
      updateJson["overrides"]["vlans"].push_back(folly::dynamic::object);
      updateJson["overrides"]["vlans"][index]["id"] = 1;
      updateJson["overrides"]["vlans"][index]["ip"] = full_path;
      updateJson["overrides"]["vlans"][index]["mask"] = "255.255.0.0";
      updateJson["overrides"]["vlans"][index]["mode"] = "static";
    } else {
      return;
      // TODO: Establish rules for updating gateway IP address (right now
      // updating it will cause the cambium device to be unable to sync with
      // cnmaestro)
      updateJson["overrides"]["auto_set"] = folly::dynamic::object;
      updateJson["overrides"]["auto_set"]["network"] = false;
      updateJson["overrides"]["default_gw"] = full_path;
    }
  }
}

} // namespace cambium
} // namespace devices
} // namespace devmand
