// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/translation/ReaderRegistry.h>
#include <devmand/devices/cli/translation/WriterRegistry.h>
#include <folly/dynamic.h>
#include <ostream>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;

struct DeviceType {
  string device;
  string version;

  friend ostream& operator<<(ostream& os, const DeviceType& type);
  string str() const;

  bool operator==(const DeviceType& rhs) const;
  bool operator!=(const DeviceType& rhs) const;

  bool operator<(const DeviceType& rhs) const;
  bool operator>(const DeviceType& rhs) const;
  bool operator<=(const DeviceType& rhs) const;
  bool operator>=(const DeviceType& rhs) const;
};

class Plugin {
 public:
  virtual DeviceType getDeviceType() const = 0;

  virtual void provideReaders(ReaderRegistryBuilder& registry) const = 0;
  virtual void provideWriters(WriterRegistryBuilder& registry) const = 0;
};

class DeviceContext : public Plugin {
 public:
  DeviceContext(vector<shared_ptr<Plugin>> _plugins);
  ~DeviceContext() = default;
  DeviceContext(const DeviceContext&) = delete;
  DeviceContext& operator=(const DeviceContext&) = delete;
  DeviceContext(DeviceContext&&) = delete;
  DeviceContext& operator=(DeviceContext&&) = delete;

  DeviceType getDeviceType() const override;
  void provideReaders(ReaderRegistryBuilder& registry) const override;
  void provideWriters(WriterRegistryBuilder& registry) const override;

  shared_ptr<DeviceContext> addPlugin(shared_ptr<Plugin> plugin) const;

 private:
  vector<shared_ptr<Plugin>> plugins;
};

class PluginRegistryException : public runtime_error {
 public:
  PluginRegistryException(string reason) : runtime_error(reason){};
};

class PluginRegistry {
 public:
  PluginRegistry() = default;
  ~PluginRegistry() = default;
  PluginRegistry(const PluginRegistry&) = delete;
  PluginRegistry& operator=(const PluginRegistry&) = delete;
  PluginRegistry(PluginRegistry&&) = delete;
  PluginRegistry& operator=(PluginRegistry&&) = delete;

  void registerPlugin(shared_ptr<Plugin> plugin);
  shared_ptr<DeviceContext> getDeviceContext(const DeviceType& type);

  friend ostream& operator<<(ostream& os, const PluginRegistry& reg);

 private:
  map<DeviceType, shared_ptr<DeviceContext>> contexts;
  bool containsDeviceType(const DeviceType& type);
};

} // namespace cli
} // namespace devices
} // namespace devmand
