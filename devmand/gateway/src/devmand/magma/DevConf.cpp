// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <folly/json.h>
#include <yaml-cpp/yaml.h>

#include <devmand/Application.h>
#include <devmand/Config.h>
#include <devmand/magma/DevConf.h>
#include <devmand/utils/FileUtils.h>
#include <devmand/utils/StringUtils.h>

namespace devmand {
namespace magma {

static constexpr const char* ymlExt = ".yml";
static constexpr const char* yamlExt = ".yaml";
static constexpr const char* mconfigExt = ".mconfig";

DevConf::DevConf(
    folly::EventBase& eventBase,
    const std::string& _deviceConfigurationFile)
    : watcher(eventBase),
      deviceConfigurationFile(_deviceConfigurationFile),
      mode(getConfigFileMode(deviceConfigurationFile.native())) {}

void DevConf::enable() {
  // this could go out of scope so handle with a weak ptr.
  std::weak_ptr<DevConf> weak(
      std::dynamic_pointer_cast<DevConf>(this->shared_from_this()));
  LOG(INFO) << "Enabling device configuration @ "
            << deviceConfigurationFile.parent_path().native();
  watcher.addWatch(
      deviceConfigurationFile.parent_path().native(),
      [weak](FileWatchEvent event) {
        if (auto shared = weak.lock()) {
          shared->handleFileWatchEvent(event);
        } else {
          LOG(ERROR) << "weak lock failed";
        }
      },
      true);
}

ConfigFileMode DevConf::getConfigFileMode(
    const std::string& deviceConfigurationFile) {
  if (endsWith(deviceConfigurationFile, ymlExt) or
      endsWith(deviceConfigurationFile, yamlExt)) {
    return ConfigFileMode::Yaml;
  } else if (endsWith(deviceConfigurationFile, mconfigExt)) {
    return ConfigFileMode::Mconfig;
  }
  throw std::runtime_error("unknown device configuration file type");
}

bool DevConf::isDeviceConfDirModifyEvent(FileWatchEvent watchEvent) const {
  return watchEvent.filename.empty() and watchEvent.event == FileEvent::Attrib;
}

bool DevConf::isDeviceConfFileModifyEvent(FileWatchEvent watchEvent) const {
  return watchEvent.filename == deviceConfigurationFile.filename() and
      (watchEvent.event == FileEvent::Modify or
       watchEvent.event == FileEvent::Attrib or
       watchEvent.event == FileEvent::MoveTo);
}

void DevConf::handleFileWatchEvent(FileWatchEvent watchEvent) {
  // TODO make this debug level
  // LOG(INFO) << "Handling file watch event "
  //          << static_cast<int>(watchEvent.event) << " on '"
  //          << watchEvent.filename << "'";

  if (isDeviceConfDirModifyEvent(watchEvent) or
      isDeviceConfFileModifyEvent(watchEvent)) {
    LOG(INFO) << "DevConf modified";
    cartography::DeviceConfigs newDeviceConfigs;
    switch (mode) {
      case ConfigFileMode::Yaml:
        newDeviceConfigs =
            parseYamlDeviceConfigs(deviceConfigurationFile.native());
        break;
      case ConfigFileMode::Mconfig:
        newDeviceConfigs =
            parseMconfigDeviceConfigs(deviceConfigurationFile.native());
        break;
    }

    DiffEventHandler<cartography::DeviceConfig> deh =
        [this](DiffEvent de, const cartography::DeviceConfig& newDeviceConfig) {
          handleDeviceDiff(de, newDeviceConfig);
        };
    diff(oldDeviceConfigs, newDeviceConfigs, deh);
    oldDeviceConfigs = newDeviceConfigs;
  }
}

void DevConf::handleDeviceDiff(
    DiffEvent de,
    const cartography::DeviceConfig& deviceConfig) {
  switch (de) {
    case DiffEvent::Add:
      add(deviceConfig);
      break;
    case DiffEvent::Modify:
      // TODO handle more gracefully
      del(deviceConfig);
      add(deviceConfig);
      break;
    case DiffEvent::Delete:
      del(deviceConfig);
      break;
  }
}

cartography::DeviceConfigs DevConf::parseYamlDeviceConfigs(
    const std::string& deviceConfigurationFile) {
  cartography::DeviceConfigs newDeviceConfigs;
  try {
    YAML::Node devicesFile = YAML::LoadFile(deviceConfigurationFile);
    for (const auto& device : devicesFile["devices"]) {
      cartography::DeviceConfig deviceConfig;
      deviceConfig.id = device["id"].as<std::string>();
      deviceConfig.platform = device["platform"].as<std::string>();
      deviceConfig.ip = device["ip"].as<std::string>();
      deviceConfig.readonly = device["readonly"].as<bool>();

      if (device["yangConfig"]) {
        deviceConfig.yangConfig =
            FileUtils::readContents(device["yangConfig"].as<std::string>());
      }
      for (const auto& channel : device["channels"]) {
        cartography::ChannelConfig channelConfig;
        const auto& channelName = channel.first.as<std::string>();
        for (const auto& kv : channel.second) {
          channelConfig.kvPairs.emplace(
              kv.first.as<std::string>(), kv.second.as<std::string>());
        }
        std::string cName = channelName.substr(0, channelName.length() - 7);
        deviceConfig.channelConfigs.emplace(cName, channelConfig);
      }
      newDeviceConfigs.emplace(deviceConfig);
    }
  } catch (const YAML::Exception& e) {
    LOG(ERROR) << "Bad devices file " << deviceConfigurationFile << " "
               << e.what();
  }
  return newDeviceConfigs;
}

template <typename... Args>
static cartography::ChannelConfig populateChannelConfigKVs(
    const folly::dynamic& channel,
    Args... keys) {
  cartography::ChannelConfig channelConfig;
  for (auto&& key : {keys...}) {
    if (channel.isObject()) {
      auto* value = channel.get_ptr(key);
      if (value != nullptr) {
        channelConfig.kvPairs.emplace(key, value->asString());
      }
    }
  }
  return channelConfig;
}

template <typename... Args>
static void populateChannelConfig(
    const std::string& channelName,
    cartography::DeviceConfig& deviceConfig,
    const folly::dynamic& device,
    Args... keys) {
  auto cName = folly::sformat("{}Channel", channelName);
  auto* channel = device.get_ptr(cName);
  if (channel != nullptr) {
    auto channelConfig = populateChannelConfigKVs(*channel, keys...);
    deviceConfig.channelConfigs.emplace(channelName, channelConfig);
  }
}

static void populateOtherChannelConfig(
    cartography::DeviceConfig& deviceConfig,
    const folly::dynamic& device) {
  auto* channel = device.get_ptr("otherChannel");
  bool isCli = false;
  if (channel != nullptr) {
    if (channel->isObject()) {
      cartography::ChannelConfig channelConfig;
      for (auto&& kv : (*channel)["channelProps"].items()) {
        channelConfig.kvPairs.emplace(
            kv.first.asString(), kv.second.asString());
        if (kv.first.asString() == "cname") {
          isCli = kv.second.asString() == "cli";
        }
      }
      if (isCli) {
        deviceConfig.channelConfigs.emplace("cli", channelConfig);
      } else {
        deviceConfig.channelConfigs.emplace("other", channelConfig);
      }
    }
  }
}

cartography::DeviceConfigs DevConf::parseMconfigDeviceConfigs(
    const std::string& deviceConfigurationFile) {
  cartography::DeviceConfigs newDeviceConfigs;
  try {
    auto contents = FileUtils::readContents(deviceConfigurationFile);
    if (not contents.empty()) {
      auto cfgFile = folly::parseJson(contents);
      auto& devmandCfg = cfgFile["configsByKey"]["devmand"];
      for (const auto& device : devmandCfg["managedDevices"].items()) {
        cartography::DeviceConfig deviceConfig;
        deviceConfig.id = device.first.asString();
        deviceConfig.platform = device.second["platform"].asString();
        deviceConfig.ip = device.second["host"].asString();
        if (device.second.get_ptr("readonly") != nullptr) {
          // TODO make this term configurable
          deviceConfig.readonly = device.second["readonly"].asBool();
        }
        if (device.second.get_ptr("deviceConfig") != nullptr) {
          deviceConfig.yangConfig = device.second["deviceConfig"].asString();
        }
        auto channels = device.second["channels"];
        populateChannelConfig(
            "frinx",
            deviceConfig,
            channels,
            "authorization",
            "deviceType",
            "deviceVersion",
            "frinxPort",
            "host",
            "password",
            "port",
            "transportType",
            "username");
        populateChannelConfig(
            "snmp", deviceConfig, channels, "community", "version");
        populateChannelConfig(
            "cambium",
            deviceConfig,
            channels,
            "clientId",
            "clientSecret",
            "clientMac",
            "clientIp");
        populateOtherChannelConfig(deviceConfig, channels);
        newDeviceConfigs.emplace(deviceConfig);
      }
    }
  } catch (const std::exception& e) {
    LOG(ERROR) << "Bad devices file " << deviceConfigurationFile << " "
               << e.what();
  }
  return newDeviceConfigs;
}

} // namespace magma
} // namespace devmand
