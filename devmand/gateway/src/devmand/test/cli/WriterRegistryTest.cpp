// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/translation/BindingWriterRegistry.h>
#include <devmand/devices/cli/translation/WriterRegistry.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <gtest/gtest.h>
#include <ydk_openconfig/openconfig_interfaces.hpp>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::devices::cli;
using namespace devmand::test::utils::cli;
using namespace folly;
using Ifc = openconfig::openconfig_interfaces::Interfaces::Interface;

class WriterRegistryTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

class NoopWriter : public Writer {
 public:
  virtual Future<Unit> create(
      const Path& path,
      dynamic cfg,
      const DeviceAccess& device) const override {
    MLOG(MDEBUG) << path << cfg << &device;
    return Future<Unit>(unit);
  };
  virtual Future<Unit> update(
      const Path& path,
      dynamic before,
      dynamic after,
      const DeviceAccess& device) const override {
    MLOG(MDEBUG) << path << before << after << &device;
    return Future<Unit>(unit);
  }
  virtual Future<Unit> remove(
      const Path& path,
      dynamic before,
      const DeviceAccess& device) const override {
    MLOG(MDEBUG) << path << before << &device;
    return Future<Unit>(unit);
  };
};

class IfcConfigWriter : public BindingWriter<Ifc::Config> {
 public:
  Future<Unit> create(
      const Path& path,
      shared_ptr<Ifc::Config> cfg,
      const DeviceAccess& device) const override {
    (void)path;
    return device.cli()
        ->executeWrite(WriteCommand::create("Writing IFC " + cfg->name.value))
        .via(device.executor().get())
        .thenValue([](auto output) {
          (void)output;
          return unit;
        });
  }

  Future<Unit> update(
      const Path& path,
      shared_ptr<Ifc::Config> before,
      shared_ptr<Ifc::Config> after,
      const DeviceAccess& device) const override {
    (void)path;
    (void)before;
    return device.cli()
        ->executeWrite(
            WriteCommand::create("Updating IFC " + after->name.value))
        .via(device.executor().get())
        .thenValue([](auto output) {
          (void)output;
          return unit;
        });
    ;
  }

  Future<Unit> remove(
      const Path& path,
      shared_ptr<Ifc::Config> before,
      const DeviceAccess& device) const override {
    (void)path;
    return device.cli()
        ->executeWrite(
            WriteCommand::create("Deleting IFC " + before->name.value))
        .via(device.executor().get())
        .thenValue([](auto output) {
          (void)output;
          return unit;
        });
  }
};

TEST_F(WriterRegistryTest, api) {
  //  ModelRegistry models;
  //  auto executor = make_shared<CPUThreadPoolExecutor>(2);
  //  DeviceAccess mockDevice{make_shared<EchoCli>(), "rest", executor};
  //  WriterRegistryBuilder reg;
  //
  //  BindingContext& bindingCtx =
  //      models.getBindingContext(Model::OPENCONFIG_2_4_3);
  //  BINDING_W(reg, bindingCtx)
  //      .add(
  //          "/openconfig-interfaces:interfaces/interface/config",
  //          make_shared<IfcConfigWriter>(),
  //          {"/openconfig-network-instance:network-instances"});
  //  reg.add(
  //      "/openconfig-network-instance:network-instances",
  //      make_shared<NoopWriter>());
  //  auto r = reg.build();
  //
  //  MLOG(MDEBUG) << *r;
  //
  //  const shared_ptr<Ifc::Config>& before = make_shared<Ifc::Config>();
  //  before->name = "eth 0/1";
  //  const shared_ptr<Ifc::Config>& after = make_shared<Ifc::Config>();
  //  after->name = "eth 0/1";
  //  after->description = "descr";
  //
  //  std::multimap<Path, DatastoreDiff> diff = {
  //      {"/openconfig-interfaces:interfaces/interface/config",
  //       DatastoreDiff(
  //           bindingCtx.getCodec().toDom(
  //               "/openconfig-interfaces:interfaces/interface[name='eth
  //               0/1']/config", *before),
  //           bindingCtx.getCodec().toDom(
  //               "/openconfig-interfaces:interfaces/interface[name='eth
  //               0/1']/config", *after),
  //           DatastoreDiffType::update,
  //           "/openconfig-interfaces:interfaces/interface[name='eth
  //           0/1']/config")},
  //  };
  //  r->write(diff, mockDevice);
  //
  //  // Let the executor finish
  //  via(executor.get(), []() {}).get();
  //  executor->join();
}

TEST_F(WriterRegistryTest, writerDependencyLoop) {
  WriterRegistryBuilder reg;

  reg.add(
      "/openconfig-interfaces:interfaces/interface/config",
      make_shared<NoopWriter>(),
      {"/openconfig-network-instance:network-instances"});
  reg.add(
      "/openconfig-network-instance:network-instances",
      make_shared<NoopWriter>(),
      {"/openconfig-interfaces:interfaces/interface/config"});

  ASSERT_THROW(reg.build(), WriterRegistryException);
}

TEST_F(WriterRegistryTest, wrongPath) {
  llly_verb(LLLY_LOG_LEVEL::LLLY_LLDBG);
  ModelRegistry models;
  WriterRegistryBuilder reg{models.getSchemaContext(Model::OPENCONFIG_2_4_3)};

  ASSERT_THROW(
      reg.add("/NOTEXISTING", make_shared<NoopWriter>()),
      WriterRegistryException);

  ASSERT_THROW(
      reg.add(
          "/openconfig-interfaces:interfaces/interface/config",
          make_shared<NoopWriter>(),
          {"/NOTEXISTING:ELEMENT"}),
      WriterRegistryException);
}

} // namespace cli
} // namespace test
} // namespace devmand
