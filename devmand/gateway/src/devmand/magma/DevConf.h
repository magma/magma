// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#pragma once

#include <experimental/filesystem>

#include <devmand/Diff.h>
#include <devmand/FileWatcher.h>
#include <devmand/cartography/DeviceConfig.h>
#include <devmand/cartography/Method.h>

namespace devmand {
namespace magma {

enum class ConfigFileMode { Yaml, Mconfig };

/*
 * This class implements a simple device discovery method from a file.
 */
class DevConf : public cartography::Method {
 public:
  DevConf(
      folly::EventBase& eventBase,
      const std::string& deviceConfigurationFile);
  DevConf() = delete;
  virtual ~DevConf() = default;
  DevConf(const DevConf&) = delete;
  DevConf& operator=(const DevConf&) = delete;
  DevConf(DevConf&&) = delete;
  DevConf& operator=(DevConf&&) = delete;

 public:
  void enable() override;

 private:
  void handleFileWatchEvent(FileWatchEvent event);
  void handleDeviceDiff(
      DiffEvent event,
      const cartography::DeviceConfig& deviceConfig);

  bool isDeviceConfDirModifyEvent(FileWatchEvent watchEvent) const;
  bool isDeviceConfFileModifyEvent(FileWatchEvent watchEvent) const;

  static cartography::DeviceConfigs parseYamlDeviceConfigs(
      const std::string& deviceConfigurationFile);
  static cartography::DeviceConfigs parseMconfigDeviceConfigs(
      const std::string& deviceConfigurationFile);

  static ConfigFileMode getConfigFileMode(
      const std::string& deviceConfigurationFile);

 private:
  FileWatcher watcher;
  const std::experimental::filesystem::path deviceConfigurationFile;
  ConfigFileMode mode;

  cartography::DeviceConfigs oldDeviceConfigs;
};

} // namespace magma
} // namespace devmand
