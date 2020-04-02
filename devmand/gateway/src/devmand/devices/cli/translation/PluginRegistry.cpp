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

void DeviceContext::provideWriters(WriterRegistryBuilder& registry) const {
  for (const auto& plugin : plugins) {
    plugin->provideWriters(registry);
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

void PluginRegistry::registerFlavours(
    map<DeviceType, shared_ptr<CliFlavourParameters>> newFlavours) {
  for (auto entry : newFlavours) {
    flavours.insert_or_assign(entry.first, CliFlavour::create(entry.second));
  }
}

shared_ptr<CliFlavour> PluginRegistry::getCliFlavour(
    const DeviceType& deviceType) {
  auto result = flavours[deviceType];
  if (not result) {
    MLOG(MDEBUG) << "Flavour not found, using default for " << deviceType;
    result = CliFlavour::getDefaultInstance();
  }
  return result;
}

} // namespace cli
} // namespace devices
} // namespace devmand
