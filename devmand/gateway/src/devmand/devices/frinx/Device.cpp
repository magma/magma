// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/frinx/Device.h>

#include <iostream>

#include <folly/dynamic.h>
#include <folly/json.h>

#include <devmand/Application.h>
#include <devmand/error/ErrorHandler.h>

namespace devmand {
namespace devices {
namespace frinx {

const char* connectTemplate = R"template({{
  "network-topology:node": {{
    "network-topology:node-id": "{}",
    "cli-topology:host": "{}",
    "cli-topology:port": "{}",
    "cli-topology:transport-type": "{}",
    "cli-topology:device-type": "{}",
    "cli-topology:device-version": "{}",
    "cli-topology:username": "{}",
    "cli-topology:password": "{}",
    "cli-topology:journal-size": 150
  }}
}})template";

constexpr const char* getOperationalDatastoreEpTemplate =
    "/restconf/config/"
    "network-topology:network-topology/topology/"
    "unified/node/{}/yang-ext:mount";
constexpr const char* setRunningDatastoreEpTemplate =
    "/restconf/config/"
    "network-topology:network-topology/topology/"
    "unified/node/{}/yang-ext:mount";
constexpr const char* connectEpTemplate =
    "/restconf/config/network-topology:network-topology/topology/cli/node/{}";
constexpr const char* checkConnectEpTemplate =
    "/restconf/operational/network-topology:network-topology/topology/cli/node/{}";
constexpr const char* errorTemplate = "Error on endpoint {} ({})";
constexpr const char* contentTypeJson = "application/json";

std::shared_ptr<devices::Device> Device::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  const auto& channelConfigs = deviceConfig.channelConfigs;
  auto& frinxKv = channelConfigs.at("frinx").kvPairs;
  return std::make_unique<devices::frinx::Device>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
      frinxKv.at("host"),
      folly::to<int>(frinxKv.at("frinxPort")),
      folly::IPAddress(deviceConfig.ip),
      folly::to<int>(frinxKv.at("port")),
      frinxKv.at("authorization"),
      deviceConfig.id,
      frinxKv.at("transportType"),
      frinxKv.at("deviceType"),
      frinxKv.at("deviceVersion"),
      frinxKv.at("username"),
      frinxKv.at("password"));
}

// TODO dont pass id twice
Device::Device(
    Application& application,
    const Id& id_,
    bool readonly_,
    const std::string& controllerHost,
    const int controllerPort,
    const folly::IPAddress& deviceIp_,
    const int devicePort_,
    const std::string& authorization_,
    const std::string& deviceId_,
    const std::string& transportType_,
    const std::string& deviceType_,
    const std::string& deviceVersion_,
    const std::string& deviceUsername_,
    const std::string& devicePassword_)
    : devices::Device(application, id_, readonly_),
      channel(controllerHost, controllerPort),
      headers({{"Authorization", authorization_},
               {"Accept", contentTypeJson},
               {"Content-Type", contentTypeJson}}),
      deviceIp(deviceIp_),
      devicePort(devicePort_),
      deviceId(deviceId_),
      transportType(transportType_),
      deviceType(deviceType_),
      deviceVersion(deviceVersion_),
      deviceUsername(deviceUsername_),
      devicePassword(devicePassword_) {
  connect();
}

Device::~Device() {
  // TODO disconnect();
}

void Device::connect() {
  auto ep = folly::sformat(connectEpTemplate, deviceId);
  auto body = folly::sformat(
      connectTemplate,
      deviceId,
      deviceIp.str(),
      devicePort,
      transportType,
      deviceType,
      deviceVersion,
      deviceUsername,
      devicePassword);

  ErrorHandler::thenError(
      channel.asyncPut(headers, ep, body, contentTypeJson)
          .thenValue([this, ep](auto response) {
            static const std::chrono::seconds retry{10}; // TODO make config
            if (response.isError()) {
              std::cerr << "FRINX: connect error so reconnect" << std::endl;
              app.scheduleIn([this]() { this->connect(); }, retry);
            } else {
              std::cerr << "FRINX: connect success so check" << std::endl;
              app.scheduleIn([this]() { this->checkConnection(); }, retry);
            }
          }));
}

void Device::checkConnection() {
  auto ep = folly::sformat(checkConnectEpTemplate, deviceId);
  std::cerr << "FRINX: check ep " << ep << std::endl;

  ErrorHandler::thenError(
      channel.asyncGet(headers, ep).thenValue([this, ep](auto response) {
        static const std::chrono::seconds retry{10}; // TODO make config
        if (response.isError()) {
          std::cerr << "FRINX: check error so reconnect" << std::endl;
          app.scheduleIn([this]() { this->connect(); }, retry);
        } else {
          ErrorHandler::executeWithCatch([&response, this]() {
            auto res = folly::parseJson(response.get());
            std::cerr << "FRINX: response" << response.get() << std::endl;
            auto& connStatus = res["node"][0]["cli-topology:connection-status"];
            this->connected = connStatus == "connected";
            std::cerr << "FRINX: check success so connected status "
                      << connStatus << std::endl;
          });
          if (not connected) {
            std::cerr << "FRINX: not connected so check" << std::endl;
            app.scheduleIn([this]() { this->checkConnection(); }, retry);
          }
        }
      }));
}

void Device::setIntendedDatastore(const folly::dynamic& config) {
  auto ep = folly::sformat(setRunningDatastoreEpTemplate, deviceId);
  folly::dynamic yang{folly::dynamic::object};
  const folly::dynamic* ints{nullptr};
  if (config != nullptr and config.isObject() and
      (ints = config.get_ptr("openconfig-interfaces:interfaces")) != nullptr) {
    yang["frinx-openconfig-interfaces:interfaces"] = *ints;
  }
  channel.asyncPut(headers, ep, folly::toJson(yang), contentTypeJson);
}

std::shared_ptr<Datastore> Device::getOperationalDatastore() {
  auto state = Datastore::make(app, getId());
  if (not connected) {
    return state;
  }

  auto ep = folly::sformat(getOperationalDatastoreEpTemplate, deviceId);
  state->addRequest(
      channel.asyncGet(headers, ep)
          .thenValue([state, ep](auto response) -> folly::dynamic {
            if (response.isError()) {
              state->addError(std::move(response.get()));
              return folly::dynamic::object;
            } else {
              try {
                std::cerr << "response body gb " << ep << " [" << response.get()
                          << "]" << std::endl;
                return folly::parseJson(response.get());
              } catch (...) {
                state->addError(folly::sformat(errorTemplate, ep, "parse"));
                return folly::dynamic::object;
              }
            }
          })
          .thenValue([state](auto v) { // TODO Wish i didnt have to capture this
            if (v != nullptr and v.isObject() and
                v.get_ptr("frinx-openconfig-interfaces:interfaces") !=
                    nullptr) {
              state->setStatus(true);
              state->update([&v](auto& lockedState) {
                lockedState["openconfig-interfaces:interfaces"] =
                    v["frinx-openconfig-interfaces:interfaces"];
              });
            } else {
              // this->connected = false;
              // this->connect();
              state->setStatus(false);
            }
          }));
  return state;
}

} // namespace frinx
} // namespace devices
} // namespace devmand
