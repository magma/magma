// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/ModelRegistry.h>
#include <devmand/test/cli/utils/Json.h>
#include <devmand/test/cli/utils/Log.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <gtest/gtest.h>
#include <ydk_ietf/iana_if_type.hpp>
#include <ydk_openconfig/openconfig_interfaces.hpp>
#include <ydk_openconfig/openconfig_vlan_types.hpp>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::devices::cli;
using namespace devmand::test::utils::json;

using OpenconfigInterfaces = openconfig::openconfig_interfaces::Interfaces;
using OpenconfigInterface = OpenconfigInterfaces::Interface;
using VlanType = openconfig::openconfig_vlan_types::VlanModeType;

class ModelRegistryTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }

 protected:
  ModelRegistry mreg;
};

TEST_F(ModelRegistryTest, caching) {
  Bundle& bundleOpenconfig = mreg.getBundle(Model::OPENCONFIG_0_1_6);
  Bundle& bundleOpenconfig2 = mreg.getBundle(Model::OPENCONFIG_0_1_6);
  ASSERT_EQ(&bundleOpenconfig, &bundleOpenconfig2);
  ASSERT_EQ(1, mreg.cacheSize());
}

static shared_ptr<OpenconfigInterface> interfaceCpp() {
  auto interface = make_shared<OpenconfigInterface>();
  interface->name = "loopback1";
  interface->config->name = "loopback1";
  interface->config->type = ietf::iana_if_type::SoftwareLoopback();
  interface->config->mtu = 1500;
  interface->state->ifindex = 1;
  interface->ethernet->switched_vlan->config->access_vlan = 77;
  interface->ethernet->switched_vlan->config->interface_mode = VlanType::TRUNK;
  interface->ethernet->switched_vlan->config->trunk_vlans.append(1);
  interface->ethernet->switched_vlan->config->trunk_vlans.append(100);
  return interface;
}

static shared_ptr<OpenconfigInterfaces> interfacesCpp() {
  auto interfaces = make_shared<OpenconfigInterfaces>();
  interfaces->interface.append(interfaceCpp());
  return interfaces;
}

static const string interfaceJson =
    "{\n"
    "  \"openconfig-interfaces:interfaces\": {\n"
    "    \"interface\": [\n"
    "      {\n"
    "        \"config\": {\n"
    "          \"mtu\": 1500,\n"
    "          \"name\": \"loopback1\",\n"
    "          \"type\": \"iana-if-type:softwareLoopback\"\n"
    "        },\n"
    "        \"name\": \"loopback1\",\n"
    "        \"openconfig-if-ethernet:ethernet\": {\n"
    "          \"openconfig-vlan:switched-vlan\": {\n"
    "            \"config\": {\n"
    "              \"access-vlan\": 77,\n"
    "              \"interface-mode\": \"TRUNK\",\n"
    "              \"trunk-vlans\": [\n"
    "                1,\n"
    "                100\n"
    "              ]\n"
    "            }\n"
    "          }\n"
    "        },\n"
    "        \"state\": {\n"
    "          \"ifindex\": 1\n"
    "        }\n"
    "      }\n"
    "    ]\n"
    "  }\n"
    "}";

static const string singleInterfaceJson =
    "{\n"
    "  \"openconfig-interfaces:interface\": {\n"
    "    \"config\": {\n"
    "      \"mtu\": 1500,\n"
    "      \"name\": \"loopback1\",\n"
    "      \"type\": \"iana-if-type:softwareLoopback\"\n"
    "    },\n"
    "    \"name\": \"loopback1\",\n"
    "    \"openconfig-if-ethernet:ethernet\": {\n"
    "      \"openconfig-vlan:switched-vlan\": {\n"
    "        \"config\": {\n"
    "          \"access-vlan\": 77,\n"
    "          \"interface-mode\": \"TRUNK\",\n"
    "          \"trunk-vlans\": [\n"
    "            1,\n"
    "            100\n"
    "          ]\n"
    "        }\n"
    "      }\n"
    "    },\n"
    "    \"state\": {\n"
    "      \"ifindex\": 1\n"
    "    }\n"
    "  }\n"
    "}";

TEST_F(ModelRegistryTest, jsonSerializationTopLevel) {
  Bundle& bundleOpenconfig = mreg.getBundle(Model::OPENCONFIG_0_1_6);

  shared_ptr<OpenconfigInterfaces> originalIfc = interfacesCpp();

  const string& interfaceEncoded = bundleOpenconfig.encode(*originalIfc);
  ASSERT_EQ(sortJson(interfaceJson), sortJson(interfaceEncoded));

  const shared_ptr<Entity> decodedIfcEntity = bundleOpenconfig.decode(
      interfaceEncoded, make_shared<OpenconfigInterfaces>());
  ASSERT_TRUE(*decodedIfcEntity == *originalIfc);
}

TEST_F(ModelRegistryTest, jsonSerializationFail) {
  Bundle& bundleOpenconfig = mreg.getBundle(Model::OPENCONFIG_0_1_6);

  class FakeEntity : public OpenconfigInterfaces {
   public:
    string get_segment_path() const override {
      return "230-4932-4-=23";
    }
  };

  const shared_ptr<FakeEntity> ptr = make_shared<FakeEntity>();
  ASSERT_THROW(
      const string& interfaceEncoded = bundleOpenconfig.encode(*ptr),
      SerializationException);
}

TEST_F(ModelRegistryTest, jsonDeserializationFail) {
  Bundle& bundleOpenconfig = mreg.getBundle(Model::OPENCONFIG_0_1_6);

  ASSERT_THROW(
      const shared_ptr<Entity> decodedIfcEntity = bundleOpenconfig.decode(
          "not a json", make_shared<OpenconfigInterface>()),
      SerializationException);
}

TEST_F(ModelRegistryTest, jsonSerializationNestedMultiThread) {
  folly::CPUThreadPoolExecutor executor(8);

  for (int i = 0; i < 100; i++) {
    folly::Future<folly::Unit> f = folly::via(&executor, [&, i]() {
      Bundle& bundleOpenconfig = mreg.getBundle(Model::OPENCONFIG_0_1_6);

      //      DLOG(INFO) << "Executing: " << i << " on thread "
      //                 << std::this_thread::get_id()
      //                 << " with bundle: " << &bundleOpenconfig << endl;

      shared_ptr<OpenconfigInterface> originalIfc = interfaceCpp();

      const string& interfaceEncoded = bundleOpenconfig.encode(*originalIfc);
      ASSERT_EQ(sortJson(singleInterfaceJson), sortJson(interfaceEncoded));

      const shared_ptr<Entity> decodedIfcEntity = bundleOpenconfig.decode(
          interfaceEncoded, make_shared<OpenconfigInterface>());
      ASSERT_TRUE(*decodedIfcEntity == *originalIfc);
    });
  }
  executor.join();
  ASSERT_EQ(1, mreg.cacheSize());
}

} // namespace cli
} // namespace test
} // namespace devmand
