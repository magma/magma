// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <list>

#include <devmand/cartography/Cartographer.h>
#include <devmand/magma/DevConf.h>
#include <devmand/test/EventBaseTest.h>
#include <devmand/test/Notifier.h>
#include <devmand/utils/FileWatcher.h>

namespace devmand {
namespace test {

static constexpr const char* yamlDeviceConfig = "/tmp/deviceConfig.yml";
static constexpr const char* mconfigDeviceConfig = "/tmp/deviceConfig.mconfig";

const char* yamlTemplate = R"template(devices:
{})template";

const char* yamlDeviceTemplate = R"template(    - id: {}
      ip: {}
      type: Device
      platform: {}
      readonly: true
      poll:
          seconds: 10
      channels:
          snmpChannel:
              version: v1
              community: public)template";

const char* mconfigTemplate = R"template(
{{
 "configsByKey": {{
  "control_proxy": {{
   "@type": "type.googleapis.com/magma.mconfig.ControlProxy",
   "logLevel": "INFO"
  }},
  "devmand": {{
   "@type": "type.googleapis.com/magma.mconfig.DevmandGatewayConfig",
   "managedDevices": {{
        {}
   }}
  }},
  "magmad": {{
   "@type": "type.googleapis.com/magma.mconfig.MagmaD",
   "logLevel": "INFO",
   "checkinInterval": 12,
   "checkinTimeout": 10,
   "autoupgradeEnabled": true,
   "autoupgradePollInterval": 300,
   "packageVersion": "0.0.0-0",
   "images": [
   ],
   "tierId": "default",
   "featureFlags": {{
    "newfeature1": true,
    "newfeature2": false
   }},
   "dynamicServices": [
   ]
  }},
  "metricsd": {{
   "@type": "type.googleapis.com/magma.mconfig.MetricsD",
   "logLevel": "INFO"
  }}
 }},
 "metadata": {{
  "createdAt": "1563307102"
 }}
}}
)template";

const char* mconfigDeviceTemplate = R"template(
        "{}": {{
            "deviceConfig": "foobar",
            "deviceType": [
                "wifi",
                "access_point"
            ],
            "host": "{}",
            "platform": "{}",
            "readonly": true,
            "channels" : {{
              "frinxChannel": {{
                  "authorization": "Basic",
                  "deviceType": "ubnt es",
                  "deviceVersion": "1.8.2",
                  "frinxPort": 8181,
                  "host": "203.0.113.2",
                  "password": "string",
                  "port": 22,
                  "transportType": "ssh",
                  "username": "string"
              }},
              "snmpChannel": {{
                  "community": "public",
                  "version": "v1"
              }}
            }}
        }}
)template";

class DevConfTest : public EventBaseTest {
 public:
  DevConfTest() = default;

  ~DevConfTest() override = default;
  DevConfTest(const DevConfTest&) = delete;
  DevConfTest& operator=(const DevConfTest&) = delete;
  DevConfTest(DevConfTest&&) = delete;
  DevConfTest& operator=(DevConfTest&&) = delete;

 protected:
  void expectDeviceConfig(
      std::list<cartography::DeviceConfig>& where,
      cartography::DeviceConfig&& deviceConfig) {
    EXPECT_EQ(deviceConfig, where.front());
    where.pop_front();
  }

 protected:
  cartography::AddHandler add =
      [this](const cartography::DeviceConfig& deviceConfig) {
        EXPECT_LE(1, expectedNumAdd);
        adds.push_back(deviceConfig);
        if (--expectedNumAdd == 0) {
          addNotifier.notify();
        }
      };
  cartography::DeleteHandler del =
      [this](const cartography::DeviceConfig& deviceConfig) {
        EXPECT_LE(1, expectedNumDel);
        dels.push_back(deviceConfig);
        if (--expectedNumDel == 0) {
          delNotifier.notify();
        }
      };
  Notifier addNotifier;
  Notifier delNotifier;
  unsigned int expectedNumAdd = 0;
  unsigned int expectedNumDel = 0;
  std::list<cartography::DeviceConfig> adds;
  std::list<cartography::DeviceConfig> dels;

  template <typename... Args>
  static std::string formatYaml(Args... args) {
    return folly::sformat(
        yamlTemplate, folly::sformat(yamlDeviceTemplate, args...));
  }

  template <typename... Args>
  static std::string formatMConfig(Args... args) {
    return folly::sformat(
        mconfigTemplate, folly::sformat(mconfigDeviceTemplate, args...));
  }
};

TEST_F(DevConfTest, initialReadYaml) {
  FileUtils::write(yamlDeviceConfig, formatYaml("foo", "203.0.113.1", "Foo"));
  expectedNumAdd = 1;
  cartography::Cartographer cartographer{add, del};
  cartographer.addDeviceDiscoveryMethod(
      std::make_shared<magma::DevConf>(eventBase, yamlDeviceConfig));
  addNotifier.wait();
  EXPECT_EQ(1, adds.size());
  auto yamlDevice = cartography::DeviceConfig{
      "foo",
      "Foo",
      "203.0.113.1",
      "",
      true,
      std::map<std::string, cartography::ChannelConfig>{
          {"snmp",
           {std::map<std::string, std::string>{{"version", "v1"},
                                               {"community", "public"}}}}}};
  expectDeviceConfig(adds, std::move(yamlDevice));
  stop();
}

TEST_F(DevConfTest, badYaml) {
  FileUtils::write(yamlDeviceConfig, "THIS IS NOT YAML");
  cartography::Cartographer cartographer{add, del};
  cartographer.addDeviceDiscoveryMethod(
      std::make_shared<magma::DevConf>(eventBase, yamlDeviceConfig));
  stop();
}

TEST_F(DevConfTest, badMconfig) {
  FileUtils::write(mconfigDeviceConfig, "THIS IS NOT JSON");
  cartography::Cartographer cartographer{add, del};
  cartographer.addDeviceDiscoveryMethod(
      std::make_shared<magma::DevConf>(eventBase, mconfigDeviceConfig));
  stop();
}

TEST_F(DevConfTest, badExtension) {
  FileUtils::write(
      "/tmp/deviceConfig.txt", formatYaml("foo", "203.0.113.1", "Foo"));
  cartography::Cartographer cartographer{add, del};
  EXPECT_THROW(
      cartographer.addDeviceDiscoveryMethod(
          std::make_shared<magma::DevConf>(eventBase, "/tmp/deviceConfig.txt")),
      std::runtime_error);
  stop();
}

TEST_F(DevConfTest, modifyIpAddress) {
  FileUtils::write(yamlDeviceConfig, formatYaml("foo", "203.0.113.1", "Foo"));
  expectedNumAdd = 1;
  cartography::Cartographer cartographer{add, del};
  cartographer.addDeviceDiscoveryMethod(
      std::make_shared<magma::DevConf>(eventBase, yamlDeviceConfig));
  addNotifier.wait();
  EXPECT_EQ(1, adds.size());

  expectedNumDel = 1;
  expectedNumAdd = 1;

  FileUtils::write(yamlDeviceConfig, formatYaml("foo", "203.0.113.2", "Foo"));
  delNotifier.wait();
  addNotifier.wait();
  stop();
}

TEST_F(DevConfTest, initialReadMConfig) {
  FileUtils::write(
      mconfigDeviceConfig, formatMConfig("foo", "203.0.113.1", "Foo"));
  expectedNumAdd = 1;
  cartography::Cartographer cartographer{add, del};
  cartographer.addDeviceDiscoveryMethod(
      std::make_shared<magma::DevConf>(eventBase, mconfigDeviceConfig));
  addNotifier.wait();
  EXPECT_EQ(1, adds.size());
  auto mconfigDevice = cartography::DeviceConfig{
      "foo",
      "Foo",
      "203.0.113.1",
      "foobar",
      true,
      std::map<std::string, cartography::ChannelConfig>{
          {"snmp",
           {std::map<std::string, std::string>{{"version", "v1"},
                                               {"community", "public"}}}},
          {"frinx",
           {std::map<std::string, std::string>{{"authorization", "Basic"},
                                               {"deviceType", "ubnt es"},
                                               {"deviceVersion", "1.8.2"},
                                               {"frinxPort", "8181"},
                                               {"host", "203.0.113.2"},
                                               {"password", "string"},
                                               {"port", "22"},
                                               {"transportType", "ssh"},
                                               {"username", "string"}}}}}};

  expectDeviceConfig(adds, std::move(mconfigDevice));
  stop();
}

TEST_F(DevConfTest, initialReadEmptyYaml) {
  FileUtils::write(yamlDeviceConfig, "devices: []");
  expectedNumAdd = 0;
  cartography::Cartographer cartographer{add, del};
  cartographer.addDeviceDiscoveryMethod(
      std::make_shared<magma::DevConf>(eventBase, yamlDeviceConfig));
  EXPECT_EQ(0, adds.size());
  stop();
}

} // namespace test
} // namespace devmand
