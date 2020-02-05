// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/test/cli/utils/Json.h>
#include <devmand/test/cli/utils/Log.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
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

class BindingContextTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }

 protected:
  ModelRegistry mreg;
};

TEST_F(BindingContextTest, caching) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);
  BindingContext& bundleOpenconfig2 =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);
  ASSERT_EQ(&bundleOpenconfig, &bundleOpenconfig2);
  ASSERT_EQ(1, mreg.bindingCacheSize());
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
    "  \"openconfig-interfaces:interface\": [\n"
    "    {\n"
    "      \"state\": {\n"
    "        \"ifindex\": 1\n"
    "      },\n"
    "      \"openconfig-if-ethernet:ethernet\": {\n"
    "        \"openconfig-vlan:switched-vlan\": {\n"
    "          \"config\": {\n"
    "            \"trunk-vlans\": [\n"
    "              1,\n"
    "              100\n"
    "            ],\n"
    "            \"access-vlan\": 77,\n"
    "            \"interface-mode\": \"TRUNK\"\n"
    "          }\n"
    "        }\n"
    "      },\n"
    "      \"config\": {\n"
    "        \"name\": \"loopback1\",\n"
    "        \"mtu\": 1500,\n"
    "        \"type\": \"iana-if-type:softwareLoopback\"\n"
    "      },\n"
    "      \"name\": \"loopback1\"\n"
    "    }\n"
    "  ]\n"
    "}";

TEST_F(BindingContextTest, jsonSerializationTopLevel) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);

  shared_ptr<OpenconfigInterfaces> originalIfc = interfacesCpp();

  const string& interfaceEncoded =
      bundleOpenconfig.getCodec().encode(*originalIfc);
  ASSERT_EQ(sortJson(interfaceJson), sortJson(interfaceEncoded));

  const shared_ptr<Entity> decodedIfcEntity =
      bundleOpenconfig.getCodec().decode(
          interfaceEncoded, make_shared<OpenconfigInterfaces>());
  ASSERT_TRUE(*decodedIfcEntity == *originalIfc);
}

TEST_F(BindingContextTest, domSerializationTopLevel) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);

  shared_ptr<OpenconfigInterfaces> originalIfc = interfacesCpp();

  const dynamic& interfaceEncoded = bundleOpenconfig.getCodec().toDom(
      "/openconfig-interfaces:interfaces", *originalIfc);
  ASSERT_EQ(sortJson(interfaceJson), sortJson(toPrettyJson(interfaceEncoded)));

  const shared_ptr<Entity> decodedIfcEntity =
      bundleOpenconfig.getCodec().fromDom(
          interfaceEncoded, make_shared<OpenconfigInterfaces>());
  ASSERT_TRUE(*decodedIfcEntity == *originalIfc);
}

TEST_F(BindingContextTest, domSerializationNested) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);

  shared_ptr<OpenconfigInterface> originalIfc = interfaceCpp();

  const dynamic& interfaceEncoded = bundleOpenconfig.getCodec().toDom(
      "/openconfig-interfaces:interfaces/interface[name='loopback1']",
      *originalIfc);
  ASSERT_EQ(
      sortJson(singleInterfaceJson), sortJson(toPrettyJson(interfaceEncoded)));

  const shared_ptr<Entity> decodedIfcEntity =
      bundleOpenconfig.getCodec().fromDom(
          interfaceEncoded, make_shared<OpenconfigInterface>());
  ASSERT_TRUE(*decodedIfcEntity == *originalIfc);
}

TEST_F(BindingContextTest, toDomFail) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);

  shared_ptr<OpenconfigInterfaces> originalIfc = interfacesCpp();

  ASSERT_THROW(
      bundleOpenconfig.getCodec().toDom(
          "/openconfig-interfaces:interfaces/interface[name='0/1']",
          *originalIfc),
      BindingSerializationException);
}

TEST_F(BindingContextTest, jsonSerializationFail) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);

  class FakeEntity : public OpenconfigInterfaces {
   public:
    string get_segment_path() const override {
      return "230-4932-4-=23";
    }
  };

  const shared_ptr<FakeEntity> ptr = make_shared<FakeEntity>();
  ASSERT_THROW(
      const string& interfaceEncoded = bundleOpenconfig.getCodec().encode(*ptr),
      BindingSerializationException);
}

TEST_F(BindingContextTest, jsonDeserializationFail) {
  BindingContext& bundleOpenconfig =
      mreg.getBindingContext(Model::OPENCONFIG_0_1_6);

  ASSERT_THROW(
      const shared_ptr<Entity> decodedIfcEntity =
          bundleOpenconfig.getCodec().decode(
              "not a json", make_shared<OpenconfigInterface>()),
      BindingSerializationException);
}

} // namespace cli
} // namespace test
} // namespace devmand
