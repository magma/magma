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
#include <devmand/devices/CambiumDevice.h>
#include <devmand/devices/DcsgDevice.h>
#include <devmand/devices/DemoDevice.h>
#include <devmand/devices/EchoDevice.h>
#include <devmand/devices/FrinxDevice.h>
#include <devmand/devices/PingDevice.h>
#include <devmand/devices/Snmpv2Device.h>
#include <devmand/devices/cli/PlaintextCliDevice.h>
#include <devmand/devices/cli/StructuredUbntDevice.h>
#include <devmand/devices/mikrotik/Device.h>
#include <devmand/magma/DevConf.h>
#include <devmand/magma/Service.h>

int main(int argc, char* argv[]) {
  // TODO Work around for magma issue. Get rid of this...
  // rename devmand to magma@devmand.service and remove from startup
  std::this_thread::sleep_for(std::chrono::seconds(10));

  folly::init(&argc, &argv);

  devmand::Application app;
  app.init();

  // Add services which export the shared view
  // FIXME uncomment and find out how to prevent this from blocking !!!
  app.add(std::make_unique<devmand::magma::Service>(app));

  // Add Demo Platforms
  {
    app.addPlatform("Cambium", devmand::devices::CambiumDevice::createDevice);
    app.addPlatform("Dcsg", devmand::devices::DcsgDevice::createDevice);
    app.addPlatform(
        "Cisco Catalyst 3750", devmand::devices::FrinxDevice::createDevice);
    app.addPlatform(
        "Unifi Switch 16", devmand::devices::FrinxDevice::createDevice);
  }

  // Add Production Ready Platforms
  {
    app.addPlatform("Demo", devmand::devices::DemoDevice::createDevice);
    app.addPlatform("Echo", devmand::devices::EchoDevice::createDevice);
    app.addPlatform(
        "MikroTik", devmand::devices::mikrotik::Device::createDevice);
    app.addPlatform("Ping", devmand::devices::PingDevice::createDevice);
    app.addPlatform("Snmp", devmand::devices::Snmpv2Device::createDevice);
  }
  // CLI
  {
    app.addPlatform(
        "PlaintextCli",
        devmand::devices::cli::PlaintextCliDevice::createDevice);
    app.addPlatform(
        "StructuredUbntCli",
        devmand::devices::cli::StructuredUbntDevice::createDevice);
  }

  app.setDefaultPlatform(devmand::devices::Snmpv2Device::createDevice);

  // Add ways to discover devices.
  app.addDeviceDiscoveryMethod(std::make_shared<devmand::magma::DevConf>(
      app.getEventBase(), devmand::FLAGS_device_configuration_file));

  app.run();
  return app.status();
}
