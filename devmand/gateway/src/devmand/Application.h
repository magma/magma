// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <chrono>
#include <functional>
#include <list>
#include <map>
#include <memory>
#include <set>

#include <folly/dynamic.h>
#include <folly/io/async/EventBase.h>

#include <devmand/Service.h>
#include <devmand/UnifiedView.h>
#include <devmand/cartography/Cartographer.h>
#include <devmand/channels/Engine.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <devmand/channels/packet/Engine.h>
#include <devmand/channels/ping/Engine.h>
#include <devmand/channels/snmp/Engine.h>
#include <devmand/devices/Device.h>
#include <devmand/devices/Factory.h>
#include <devmand/syslog/Manager.h>

namespace devmand {

using Services = std::list<std::unique_ptr<Service>>;
using ChannelEngines = std::set<std::unique_ptr<channels::Engine>>;
using Devices = std::map<devices::Id, std::shared_ptr<devices::Device>>;
using IPVersion = channels::ping::IPVersion;

class Application : public MetricSink {
 public:
  Application();
  virtual ~Application() = default;
  Application(const Application&) = delete;
  Application& operator=(const Application&) = delete;
  Application(Application&&) = delete;
  Application& operator=(Application&&) = delete;

 public:
  void init();
  void run();
  int status() const;

  void add(const cartography::DeviceConfig& deviceConfig);
  void del(const cartography::DeviceConfig& deviceConfig);

  void addDevice(std::shared_ptr<devices::Device>&& device);
  void addService(std::unique_ptr<Service>&& service);

  void addPlatform(
      const std::string& platform,
      devices::Factory::PlatformBuilder platformBuilder);
  void setDefaultPlatform(devices::Factory::PlatformBuilder platformBuilder);

  void addDeviceDiscoveryMethod(
      const std::shared_ptr<cartography::Method>& method);

  std::string getName() const;
  std::string getVersion() const;

  UnifiedView getUnifiedView();

  folly::EventBase& getEventBase();

  void scheduleEvery(
      std::function<void()> event,
      const std::chrono::seconds& seconds);
  void scheduleIn(
      std::function<void()> event,
      const std::chrono::seconds& seconds);

  channels::snmp::Engine& getSnmpEngine();
  channels::ping::Engine& getPingEngine(IPVersion ipv = IPVersion::v4);
  channels::ping::Engine& getPingEngine(folly::IPAddress ip);
  channels::cli::Engine& getCliEngine();

  syslog::Manager& getSyslogManager();

  virtual void setGauge(
      const std::string& key,
      double value,
      const std::string& labelName,
      const std::string& labelValue);

 private:
  void pollDevices();
  void tryToApplyRunningDatastoreToDevices();
  void doDebug();

  template <class EngineType, class... Args>
  EngineType* addEngine(Args&&... args) {
    return static_cast<EngineType*>(
        channelEngines
            .emplace(std::make_unique<EngineType>(std::forward<Args>(args)...))
            .first->get());
  }

 private:
  /*
   * Ths status code of the applicaton.
   */
  int statusCode{EXIT_SUCCESS};

  /*
   * Event base which handles the event loop
   */
  // TODO spawn one per thread or similar and assign as able?
  folly::EventBase eventBase;

  /*
   * Services provided by this application for access into the unified view.
   */
  Services services;

  /*
   * A unified view of device data to be reported back through the services.
   * Synced as it may be accessed via a few threads.
   */
  SharedUnifiedView unifiedView;

  /*
   * Channel Engines maintain the state and actions required by the various
   * channel libraries. Channels are libraries with various means of
   * communicating with devices such as SNMP, TR069, CLI, NETCONF, HTTP, etc.
   *
   * These are stored in the set but we also store ptrs to specific ones as
   * needed.
   */
  ChannelEngines channelEngines;
  channels::snmp::Engine* snmpEngine;
  channels::ping::Engine* pingEngine;
  channels::ping::Engine* pingEngineIpv6;
  channels::cli::Engine* cliEngine = nullptr;

  /*
   * Devices communicate with the off host devices through any number of
   * channels. The abstract device type defines an interface which all
   * instantiations implement.
   */
  Devices devices;
  devices::Factory deviceFactory;

  /*
   * The cartographer is a class which implements a number of methods by which
   * to discover devices on the network.
   */
  cartography::Cartographer cartographer;

  /*
   * This manages the syslog collection agent.
   */
  syslog::Manager syslogManager;

  static constexpr auto name = "devmand";
  static constexpr auto version = "0.0";
};

} // namespace devmand
