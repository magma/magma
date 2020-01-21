// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/translation/PluginRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace folly;
using namespace std;

bool DeviceType::operator==(const DeviceType& rhs) const {
  return device == rhs.device && version == rhs.version;
}

bool DeviceType::operator!=(const DeviceType& rhs) const {
  return !(rhs == *this);
}

bool DeviceType::operator<(const DeviceType& rhs) const {
  if (device < rhs.device)
    return true;
  if (rhs.device < device)
    return false;
  return version < rhs.version;
}

bool DeviceType::operator>(const DeviceType& rhs) const {
  return rhs < *this;
}

bool DeviceType::operator<=(const DeviceType& rhs) const {
  return !(rhs < *this);
}

bool DeviceType::operator>=(const DeviceType& rhs) const {
  return !(*this < rhs);
}

ostream& operator<<(ostream& os, const DeviceType& type) {
  os << "{" << type.device << ": " << type.version << "}";
  return os;
}

string DeviceType::str() const {
  stringstream strStream = stringstream();
  strStream << *this;
  return strStream.str();
}

DeviceContext::DeviceContext(vector<shared_ptr<Plugin>> _plugins)
    : plugins(_plugins) {
  if (_plugins.empty()) {
    throw PluginRegistryException(
        "Cannot create device context, no plugins available");
  }
}

DeviceType DeviceContext::getDeviceType() const {
  return plugins.at(0)->getDeviceType();
}

void DeviceContext::provideReaders(ReaderRegistryBuilder& registry) const {
  for (const auto& plugin : plugins) {
    plugin->provideReaders(registry);
  }
}

void DeviceContext::provideWriters() const {
  for (const auto& plugin : plugins) {
    plugin->provideWriters();
  }
}

shared_ptr<DeviceContext> DeviceContext::addPlugin(
    shared_ptr<Plugin> plugin) const {
  vector<shared_ptr<Plugin>> newPlugins = plugins;
  newPlugins.push_back(plugin);
  return make_shared<DeviceContext>(newPlugins);
}

void PluginRegistry::registerPlugin(shared_ptr<Plugin> plugin) {
  if (containsDeviceType(plugin->getDeviceType())) {
    contexts.insert_or_assign(
        plugin->getDeviceType(),
        contexts[plugin->getDeviceType()]->addPlugin(plugin));
  } else {
    contexts.emplace(
        plugin->getDeviceType(),
        make_shared<DeviceContext>(vector<shared_ptr<Plugin>>{plugin}));
  }
}

shared_ptr<DeviceContext> PluginRegistry::getDeviceContext(
    const DeviceType& deviceType) {
  if (containsDeviceType(deviceType)) {
    return contexts[deviceType];
  }

  throw PluginRegistryException(
      "Device not supported: " + deviceType.str() +
      ". No plugins registered for type");
}

ostream& operator<<(ostream& os, const PluginRegistry& reg) {
  os << "PluginRegistry[";
  for (auto& contexts : reg.contexts) {
    os << contexts.first;
  }
  os << "]";
  return os;
}

bool PluginRegistry::containsDeviceType(const DeviceType& type) {
  return contexts.find(type) != contexts.end();
}

} // namespace cli
} // namespace devices
} // namespace devmand
