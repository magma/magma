// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <folly/IPAddress.h>
#include <folly/Subprocess.h>

#include <devmand/MetricSink.h>
#include <devmand/channels/snmp/Channel.h>
#include <devmand/channels/snmp/Engine.h>
#include <devmand/channels/snmp/IfMib.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/Device.h>
#include <devmand/models/interface/Model.h>
#include <devmand/test/EventBaseTest.h>
#include <devmand/test/TestUtils.h>

namespace devmand {
namespace test {

class SnmpChannelTest : public EventBaseTest, public MetricSink {
 public:
  SnmpChannelTest() = default;
  ~SnmpChannelTest() override = default;
  SnmpChannelTest(const SnmpChannelTest&) = delete;
  SnmpChannelTest& operator=(const SnmpChannelTest&) = delete;
  SnmpChannelTest(SnmpChannelTest&&) = delete;
  SnmpChannelTest& operator=(SnmpChannelTest&&) = delete;

 protected:
  // TODO these should be gmocked out.
  virtual void
  setGauge(const std::string&, double, const std::string&, const std::string&) {
  }
  virtual void setGauge(const std::string&, double) {}

 protected:
  std::string local{"127.0.0.1"};
  std::string community{"public"};
  std::string version{"v1"};
};

TEST_F(SnmpChannelTest, checkSnmpTimeout) {
  channels::snmp::Engine engine(eventBase, "checkSnmpTimeout");
  channels::snmp::Oid oid{".1.3.6.1.2.1.1.4.0"};
  auto channel = std::make_shared<channels::snmp::Channel>(
      engine, local, community, version);
  EXPECT_THROW(channel->asyncGet(oid).get(), std::runtime_error);
  stop();
}

TEST_F(SnmpChannelTest, checkSnmp) {
  folly::Subprocess snmpd(std::vector<std::string>{"/usr/sbin/snmpd", "-f"});
  channels::snmp::Engine engine(eventBase, "checkSnmpTimeout");
  channels::snmp::Oid oid{".1.3.6.1.2.1.1.4.0"};
  auto channel = std::make_shared<channels::snmp::Channel>(
      engine, local, community, version);
  EXPECT_EQ(
      devmand::channels::snmp::Response(oid, folly::dynamic{""}),
      channel->asyncGet(oid).get());
  stop();
  snmpd.kill();
  snmpd.wait();
}

TEST_F(SnmpChannelTest, checkSnmpWithState) {
  folly::Subprocess snmpd(std::vector<std::string>{"/usr/sbin/snmpd", "-f"});
  channels::snmp::Engine engine(eventBase, "checkSnmpTimeout");
  channels::snmp::Oid oid{".1.3.6.1.2.1.1.4.0"};
  auto channel = std::make_shared<channels::snmp::Channel>(
      engine, local, community, version);

  auto state = devices::Datastore::make(*this, "testdid");
  state->setStatus(false);
  state->update([](auto& lockedState) {
    devmand::models::interface::Model::init(lockedState);
  });

  for (int i = 0; i < 10; ++i) {
    state->addRequest(
        channel->walk(channels::snmp::Oid{".1"}).thenValue([](auto) {}));
  }
  state->collect().wait();
  state = nullptr;

  EXPECT_EQ(0, utils::LifetimeTracker<devices::Datastore>::getLivingCount());

  stop();
  snmpd.kill();
  snmpd.wait();
}

TEST_F(SnmpChannelTest, checkSnmpTimeoutWithState) {
  channels::snmp::Engine engine(eventBase, "checkSnmpTimeout");
  channels::snmp::Oid oid{".1.3.6.1.2.1.1.4.0"};
  auto channel = std::make_shared<channels::snmp::Channel>(
      engine, local, community, version);

  auto state = devices::Datastore::make(*this, "testdid");
  state->setStatus(false);
  state->update([](auto& lockedState) {
    devmand::models::interface::Model::init(lockedState);
  });

  for (int i = 0; i < 10; ++i) {
    state->addRequest(
        channel->walk(channels::snmp::Oid{".1"}).thenValue([](auto) {}));
  }
  state->collect().wait();
  state = nullptr;

  EXPECT_EQ(0, utils::LifetimeTracker<devices::Datastore>::getLivingCount());
}

} // namespace test
} // namespace devmand
