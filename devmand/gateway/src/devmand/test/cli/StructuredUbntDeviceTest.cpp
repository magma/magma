// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <boost/algorithm/string/trim.hpp>
#include <devmand/Application.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/cli/StructuredUbntDevice.h>
#include <devmand/test/cli/utils/Json.h>
#include <devmand/test/cli/utils/Log.h>
#include <folly/futures/Future.h>
#include <folly/json.h>
#include <gtest/gtest.h>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::devices;
using namespace devmand::devices::cli;
using namespace devmand::test::utils::json;
using namespace std;

class StructuredUbntDeviceTest : public testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

class UbntFakeCli : public Cli {
 public:
  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override {
    (void)cmd;
    if (cmd.raw() == "show interfaces description") {
      return "\n"
             "Interface  Admin      Link    Description\n"
             "---------  ---------  ------  ----------------------------------------------------------------\n"
             "0/1        Enable     Up\n"
             "0/2        Enable     Down    Some descr\n"
             "0/3        Enable     Down\n"
             "0/4        Enable     Down\n"
             "0/5        Enable     Down\n"
             "0/6        Enable     Down\n"
             "0/7        Enable     Down\r\n"
             "0/8        Enable     Down\n"
             "3/6        Enable     Down\n";
    } else if (cmd.raw().find("show running-config interface ") == 0) {
      string ifcId =
          cmd.raw().substr(string("show running-config interface ").size() - 1);

      return "\n"
             "!Current Configuration:\n"
             "!\n"
             "interface  " +
          ifcId +
          "\n"
          "description 'This is ifc " +
          ifcId +
          "'\n"
          "mtu 1500\n"
          "exit\n"
          "";
    } else if (cmd.raw().find("show interface ethernet") == 0) {
      return "\n"
             "Total Packets Received (Octets)................ 427066814515\n"
             "\n"
             "Total Packets Received Without Errors.......... 8830072\n"
             "Unicast Packets Received....................... 293117\n"
             "Multicast Packets Received..................... 5769311\n"
             "Broadcast Packets Received..................... 2767644\n"
             "                  \n"
             "Receive Packets Discarded...................... 0\n"
             "\n"
             "Total Packets Received with MAC Errors......... 0\n"
             "\n"
             "Total Received Packets Not Forwarded........... 0\n"
             "802.3x Pause Frames Received................... 0\n"
             "Unacceptable Frame Type........................ 0\n"
             "\n"
             "Total Packets Transmitted (Octets)............. 72669652\n"
             "Max Frame Size................................. 1518\n"
             "\n"
             "Total Packets Transmitted Successfully......... 457849\n"
             "Unicast Packets Transmitted.................... 125182\n"
             "Multicast Packets Transmitted.................. 325039\n"
             "Broadcast Packets Transmitted.................. 7628\n"
             "\n"
             "Transmit Packets Discarded..................... 0\n"
             "\n"
             "Total Transmit Errors.......................... 0\n"
             "\n"
             "Total Transmit Packets Discarded............... 0\n"
             "Single Collision Frames........................ 0\n"
             "Multiple Collision Frames...................... 0\n"
             "Excessive Collision Frames..................... 0\n"
             "\n"
             "Load Interval.................................. 5\n"
             "Bits Per Second Received....................... 38272\n"
             "Bits Per Second Transmitted.................... 4552\n"
             "Packets Per Second Received.................... 14\n"
             "Packets Per Second Transmitted................. 1\n"
             "\n"
             "Time Since Counters Last Cleared............... 7 day 1 hr 24 min 53 sec";
    } else if (cmd.raw() == "show interfaces description") {
      return "\n"
             "Interface  Admin      Link    Description\n"
             "---------  ---------  ------  ----------------------------------------------------------------\n"
             "0/1        Enable     Up\n"
             "0/2        Enable     Down\n"
             "0/3        Enable     Down\n"
             "0/4        Enable     Down\n"
             "0/5        Disable    Down    testing\n"
             "0/6        Enable     Down\n"
             "0/7        Enable     Down\n"
             "0/8        Enable     Down\n"
             "3/6        Enable     Down\n"
             "";
    }

    return "";
  }

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override {
    (void)cmd;
    return folly::Future<string>("");
  }
};

static const string EXPECTED_OUTPUT =
    "{\n"
    "  \"openconfig-interfaces:interfaces\": {\n"
    "    \"interface\": [\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/2\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"Some descr\",\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/2\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/2\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/2\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/3\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/3\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/3\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/3\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/4\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/4\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/4\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/4\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/5\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/5\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/5\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/5\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/6\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/6\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/6\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/6\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/7\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/7\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/7\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/7\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"0/8\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"0/8\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"0/8\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  0/8\"\n"
    "        }\n"
    "      },\n"
    "      {\n"
    "        \"state\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"oper-status\": \"DOWN\",\n"
    "          \"name\": \"3/6\",\n"
    "          \"mtu\": 1518,\n"
    "          \"enabled\": true,\n"
    "          \"counters\": {\n"
    "            \"out-unicast-pkts\": 125182,\n"
    "            \"out-octets\": 72669652,\n"
    "            \"out-multicast-pkts\": 325039,\n"
    "            \"out-errors\": 0,\n"
    "            \"out-discards\": 0,\n"
    "            \"out-broadcast-pkts\": 7628,\n"
    "            \"last-clear\": \" 7 day 1 hr 24 min 53 sec\",\n"
    "            \"in-unicast-pkts\": 293117,\n"
    "            \"in-octets\": 427066814515,\n"
    "            \"in-multicast-pkts\": 5769311,\n"
    "            \"in-errors\": 0,\n"
    "            \"in-discards\": 0,\n"
    "            \"in-broadcast-pkts\": 2767644\n"
    "          },\n"
    "          \"admin-status\": \"UP\"\n"
    "        },\n"
    "        \"name\": \"3/6\",\n"
    "        \"config\": {\n"
    "          \"type\": \"iana-if-type:ethernetCsmacd\",\n"
    "          \"name\": \"3/6\",\n"
    "          \"mtu\": 1500,\n"
    "          \"enabled\": true,\n"
    "          \"description\": \"This is ifc  3/6\"\n"
    "        }\n"
    "      }\n"
    "    ]\n"
    "  },\n"
    "  \"fbc-symphony-device:system\": {\n"
    "    \"status\": \"UP\"\n"
    "  },\n"
    "  \"openconfig-network-instance:network-instances\": {\n"
    "    \"network-instance\": [\n"
    "      {\n"
    "        \"name\": \"default\"\n"
    "      }\n"
    "    ]\n"
    "  }"
    "}";

TEST_F(StructuredUbntDeviceTest, DISABLED_getOperationalDatastore) {
  devmand::Application app;
  // app.init(); <-- without this app is not properly initialized

  cartography::DeviceConfig deviceConfig;
  devmand::cartography::ChannelConfig chnlCfg;
  std::map<std::string, std::string> kvPairs;
  chnlCfg.kvPairs = kvPairs;
  deviceConfig.channelConfigs.insert(std::make_pair("cli", chnlCfg));

  auto cli = make_shared<UbntFakeCli>();
  auto channel = make_shared<Channel>("ubntTest", cli);
  std::unique_ptr<devices::Device> dev = std::make_unique<StructuredUbntDevice>(
      app, deviceConfig.id, false, channel, make_shared<ModelRegistry>());
  std::shared_ptr<Datastore> state = dev->getOperationalDatastore();
  const folly::dynamic& stateResult = state->collect().get();

  EXPECT_EQ(folly::parseJson(EXPECTED_OUTPUT), stateResult);
}

} // namespace cli
} // namespace test
} // namespace devmand
