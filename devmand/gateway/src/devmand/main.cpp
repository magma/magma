// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <memory>

#include <chrono>
#include <thread>

#include <folly/init/Init.h>

#include <devmand/Application.h>
#include <devmand/Config.h>
#include <devmand/devices/cambium/Device.h>
#include <devmand/devices/cli/PlaintextCliDevice.h>
#include <devmand/devices/cli/StructuredUbntDevice.h>
#include <devmand/devices/demo/Device.h>
#include <devmand/devices/echo/Device.h>
#include <devmand/devices/frinx/Device.h>
#include <devmand/devices/mikrotik/Device.h>
#include <devmand/devices/ping/Device.h>
#include <devmand/devices/snmpv2/Device.h>
#include <devmand/fscache/Service.h>
#include <devmand/magma/DevConf.h>
#include <devmand/magma/Service.h>

int main(int argc, char* argv[]) {
  folly::init(&argc, &argv);

  devmand::Application app;
  const shared_ptr<devmand::magma::DevConf>& devConf =
      std::make_shared<devmand::magma::DevConf>(
          app.getEventBase(), devmand::FLAGS_device_configuration_file);

  app.init(devConf);

  // Add services which export the unified view
  app.addService(std::make_unique<devmand::magma::Service>(app));
  // app.addService(std::make_unique<devmand::fscache::Service>(app));

  using namespace devmand::devices;

  // Add Demo Platforms
  {
    app.addPlatform("Cambium", cambium::Device::createDevice);
    app.addPlatform("Cisco Catalyst 3750", frinx::Device::createDevice);
  }

  // Add Production Ready Platforms
  {
    app.addPlatform("Demo", demo::Device::createDevice);
    app.addPlatform("Echo", echo::Device::createDevice);
    app.addPlatform("MikroTik", mikrotik::Device::createDevice);
    app.addPlatform("Ping", ping::Device::createDevice);
    app.addPlatform("Snmp", snmpv2::Device::createDevice);
    app.addPlatform("PlaintextCli", cli::PlaintextCliDevice::createDevice);
    app.addPlatform(
        "StructuredUbntCli", cli::StructuredUbntDevice::createDevice);
  }

  app.setDefaultPlatform(snmpv2::Device::createDevice);

  // Add ways to discover devices.
  app.addDeviceDiscoveryMethod(devConf);

  app.run();
  return app.status();
}
