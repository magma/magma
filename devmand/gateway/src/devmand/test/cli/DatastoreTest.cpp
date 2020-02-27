// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/datastore/BindingAwareDatastore.h>
#include <devmand/channels/cli/datastore/BindingAwareDatastoreTransaction.h>
#include <devmand/channels/cli/datastore/Datastore.h>
#include <devmand/channels/cli/datastore/DatastoreDiff.h>
#include <devmand/channels/cli/datastore/DatastoreTransaction.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/cli/schema/BindingContext.h>
#include <devmand/devices/cli/schema/Model.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/SampleJsons.h>
#include <folly/json.h>
#include <gtest/gtest.h>
#include <ydk/path_api.hpp>
#include <ydk_openconfig/iana_if_type.hpp>
#include <ydk_openconfig/openconfig_interfaces.hpp>
#include <ydk_openconfig/openconfig_vlan_types.hpp>

namespace devmand {
namespace test {
namespace cli {

using devmand::channels::cli::datastore::BindingAwareDatastore;
using devmand::channels::cli::datastore::BindingAwareDatastoreTransaction;
using devmand::channels::cli::datastore::Datastore;
using devmand::channels::cli::datastore::DatastoreDiff;
using devmand::channels::cli::datastore::DatastoreException;
using devmand::channels::cli::datastore::DatastoreTransaction;
using devmand::devices::cli::BindingCodec;
using devmand::devices::cli::SchemaContext;
using devmand::test::utils::cli::interface02state;
using devmand::test::utils::cli::interface02TopPath;
using devmand::test::utils::cli::interfaceCounters;
using devmand::test::utils::cli::newInterface;
using devmand::test::utils::cli::newInterfaceTopPath;
using devmand::test::utils::cli::openconfigInterfacesInterfaces;
using devmand::test::utils::cli::operStatus;
using folly::parseJson;
using folly::toPrettyJson;
using std::to_string;
using std::unique_ptr;
using OpenconfigInterfaces = openconfig::openconfig_interfaces::Interfaces;
using OpenconfigInterface = OpenconfigInterfaces::Interface;
using OpenconfigConfig = OpenconfigInterface::Config;
using VlanType = openconfig::openconfig_vlan_types::VlanModeType;

class DatastoreTest : public ::testing::Test {
 protected:
  unique_ptr<channels::cli::Engine> cliEngine;
  SchemaContext& schemaContext;
  shared_ptr<BindingCodec> bindingCodec;

 public:
  DatastoreTest()
      : cliEngine(std::make_unique<channels::cli::Engine>(dynamic::object())),
        schemaContext(cliEngine->getModelRegistry()->getSchemaContext(
            Model::OPENCONFIG_2_4_3)) {}

 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
    Model model = Model::OPENCONFIG_2_4_3;
    ydk::path::Repository repo(
        model.getDir(), ydk::path::ModelCachingOption::COMMON);
    bindingCodec =
        std::make_shared<BindingCodec>(repo, model.getDir(), schemaContext);
  }
};

static shared_ptr<OpenconfigInterfaces> ydkInterfaces() {
  auto interfaces = make_shared<OpenconfigInterfaces>();
  auto interface = make_shared<OpenconfigInterface>();
  interface->name = "0/2";
  interface->config->name = "0/2";
  interface->config->enabled = true;
  interface->config->mtu = 1500;
  interface->config->description = "this is a config description";
  interface->config->type = openconfig::iana_if_type::EthernetCsmacd();
  interface->state->admin_status = "UP";
  interface->state->description = "dummy state";
  interface->state->enabled = true;
  interface->state->mtu = 1518;
  interface->state->oper_status = "DOWN";
  interface->state->name = "0/2";
  interface->state->type = openconfig::iana_if_type::EthernetCsmacd();
  interfaces->interface.append(interface);
  return interfaces;
}

TEST_F(DatastoreTest, commitWorks) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite("/", parseJson(openconfigInterfacesInterfaces));
  transaction->commit();
  transaction = datastore.newTx();
  dynamic data = transaction->read("/openconfig-interfaces:interfaces");
  // in-broadcast-pkts has values 2767640, 2767641, 2767642 in the given
  // interfaces
  for (int i = 0; i < 3; ++i) {
    EXPECT_EQ(
        data["openconfig-interfaces:interfaces"]["interface"][i]["state"]
            ["counters"]["in-broadcast-pkts"],
        "276764" + to_string(i));
  }
}

TEST_F(DatastoreTest, twoTransactionsAtTheSameTimeNotPermitted) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  EXPECT_THROW(datastore.newTx(), DatastoreException);
}

TEST_F(DatastoreTest, abortDisablesRunningTransaction) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->abort();
  EXPECT_THROW(transaction->read("/whatever"), DatastoreException);
  EXPECT_THROW(
      transaction->overwrite("/", dynamic::object()), DatastoreException);
  EXPECT_THROW(transaction->merge("/", dynamic::object()), DatastoreException);
  EXPECT_THROW(transaction->abort(), DatastoreException);
  EXPECT_THROW(transaction->delete_("/whatever"), DatastoreException);
  EXPECT_THROW(transaction->commit(), DatastoreException);
  EXPECT_THROW(transaction->isValid(), DatastoreException);
  EXPECT_THROW(transaction->diff(), DatastoreException);
}

TEST_F(DatastoreTest, commitDisablesRunningTransaction) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite("/", parseJson(openconfigInterfacesInterfaces));

  transaction->commit();
  EXPECT_THROW(transaction->read("/whatever"), DatastoreException);
  EXPECT_THROW(
      transaction->overwrite("/", dynamic::object()), DatastoreException);
  EXPECT_THROW(transaction->merge("/", dynamic::object()), DatastoreException);
  EXPECT_THROW(transaction->abort(), DatastoreException);
  EXPECT_THROW(transaction->delete_("/whatever"), DatastoreException);
  EXPECT_THROW(transaction->commit(), DatastoreException);
  EXPECT_THROW(transaction->isValid(), DatastoreException);
  EXPECT_THROW(transaction->diff(), DatastoreException);
}

TEST_F(DatastoreTest, commitEndsRunningTransaction) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite("/", parseJson(openconfigInterfacesInterfaces));

  transaction->commit();
  EXPECT_NO_THROW(datastore.newTx());
}

TEST_F(DatastoreTest, dontAllowEmptyCommit) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  EXPECT_THROW(transaction->commit(), DatastoreException);
}

TEST_F(DatastoreTest, abortEndsRunningTransaction) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->abort();
  EXPECT_NO_THROW(datastore.newTx());
}

TEST_F(DatastoreTest, deleteSubtrees) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite("/", parseJson(openconfigInterfacesInterfaces));
  const char* interface03 =
      "/openconfig-interfaces:interfaces/interface[name='0/3']";
  EXPECT_TRUE(transaction->read(interface03) != nullptr);
  transaction->delete_(interface03);
  EXPECT_TRUE(transaction->read(interface03) == nullptr);

  transaction->abort();
}

TEST_F(DatastoreTest, writeNewInterface) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite("/", parseJson(openconfigInterfacesInterfaces));
  const char* interface85 =
      "/openconfig-interfaces:interfaces/interface[name='0/85']";

  transaction->overwrite(interface85, parseJson(newInterface));

  dynamic data = transaction->read(
      "/openconfig-interfaces:interfaces/interface[name='0/85']");
  transaction->abort();
  EXPECT_EQ(
      "0/85", data["openconfig-interfaces:interface"][0]["name"].getString());
}

// this test is disabled due to lowered threshold for validation and it would
// not detect a missing config section in interface. The validation threshold is
// lowered because each device supports a different subset of YANG models (one
// device has BGP activated and another does not, the device where BGP is not
// activated would have "missing" mandatory BGP leafs and the datastore would
// report its data as invalid)
TEST_F(DatastoreTest, DISABLED_identifyInvalidTree) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path(""), parseJson(openconfigInterfacesInterfaces));
  transaction->delete_(
      "/openconfig-interfaces:interfaces/interface[name='0/2']/config");

  EXPECT_FALSE(transaction->isValid());
  transaction->abort();
}

TEST_F(DatastoreTest, mergeChangesInterface) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));
  transaction->overwrite(newInterfaceTopPath, parseJson(newInterface));

  dynamic state = transaction->read(newInterfaceTopPath + "/state");
  state["openconfig-interfaces:state"]["mtu"] = 1555;
  state["openconfig-interfaces:state"]["oper-status"] = "UP";
  transaction->merge(newInterfaceTopPath + "/state", state);

  state = transaction->read(newInterfaceTopPath + "/state");
  transaction->abort();
  EXPECT_EQ(state["openconfig-interfaces:state"]["mtu"], 1555);
  EXPECT_EQ(state["openconfig-interfaces:state"]["oper-status"], "UP");
}

TEST_F(DatastoreTest, mergeErasedValueOriginalValueUnchangedInterface) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));

  dynamic state = transaction->read(interface02TopPath + "/state");
  state["openconfig-interfaces:state"].erase("mtu");
  transaction->merge(interface02TopPath + "/state", state);
  state = transaction->read(interface02TopPath + "/state");
  EXPECT_EQ(state["openconfig-interfaces:state"]["mtu"], 1518);
}

TEST_F(DatastoreTest, changeLeaf) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));

  dynamic enabled = transaction->read(interface02TopPath + "/state/enabled");
  MLOG(MINFO) << toPrettyJson(enabled);
  enabled["openconfig-interfaces:enabled"] = false;
  transaction->merge(interface02TopPath + "/state/enabled", enabled);

  enabled = transaction->read(interface02TopPath + "/state/enabled");
  EXPECT_EQ(enabled["openconfig-interfaces:enabled"], false);
}

TEST_F(DatastoreTest, changeLeafDiff) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));

  transaction->commit();
  transaction = datastore.newTx();
  dynamic errors = transaction->read(interface02TopPath + "/state/counters");
  errors["openconfig-interfaces:counters"]["out-errors"] = "777";
  errors["openconfig-interfaces:counters"]["out-discards"] = "17";
  transaction->merge(interface02TopPath + "/state/counters", errors);
  const map<Path, DatastoreDiff>& map = transaction->diff();
  EXPECT_EQ(map.size(), 2);

  Path outDiscards(interfaceCounters + "/openconfig-interfaces:out-discards");
  Path outErrors(interfaceCounters + "/openconfig-interfaces:out-errors");
  EXPECT_EQ(
      map.at(outDiscards).after["openconfig-interfaces:out-discards"], "17");
  EXPECT_EQ(map.at(outErrors).after["openconfig-interfaces:out-errors"], "777");
  EXPECT_EQ(
      map.at(outDiscards).before["openconfig-interfaces:out-discards"], "0");
  EXPECT_EQ(map.at(outErrors).before["openconfig-interfaces:out-errors"], "0");
  EXPECT_EQ(map.at(outDiscards).type, DatastoreDiffType::update);
  EXPECT_EQ(map.at(outErrors).type, DatastoreDiffType::update);
}

TEST_F(DatastoreTest, deleteSubtreeDiff) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));

  transaction->commit();
  transaction = datastore.newTx();
  Path stateToDelete(interface02TopPath + "/state");
  transaction->delete_(stateToDelete);

  const map<Path, DatastoreDiff>& map = transaction->diff();
  EXPECT_EQ(map.size(), 1);
  Path stateKey(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:state");
  EXPECT_EQ(map.begin()->first.str(), stateKey.str());
  EXPECT_EQ(toPrettyJson(map.at(stateKey).before), interface02state);
  EXPECT_EQ(toPrettyJson(map.at(stateKey).after), "{}");
  EXPECT_EQ(map.at(stateKey).type, DatastoreDiffType::deleted);
}

TEST_F(DatastoreTest, diffAfterWrite) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite("/", parseJson(openconfigInterfacesInterfaces));
  transaction->commit();
  transaction = datastore.newTx();
  const char* interface85 =
      "/openconfig-interfaces:interfaces/interface[name='0/85']";

  transaction->overwrite(interface85, parseJson(newInterface));
  const map<Path, DatastoreDiff>& map = transaction->diff();
  Path diffKey(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface");
  EXPECT_EQ(map.size(), 1);
  EXPECT_EQ(
      map.at(diffKey)
          .after["openconfig-interfaces:interface"][0]["state"]["name"],
      "0/85");
  EXPECT_EQ(
      map.at(diffKey)
          .before["openconfig-interfaces:interfaces"]["interface"]
          .size(),
      3);
  EXPECT_EQ(map.at(diffKey).type, DatastoreDiffType::create);
  transaction->abort();
}

TEST_F(DatastoreTest, diffAfterMerge) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));
  transaction->overwrite(newInterfaceTopPath, parseJson(newInterface));

  dynamic state = transaction->read(newInterfaceTopPath + "/state");
  state["openconfig-interfaces:state"]["oper-status"] = "UP";
  transaction->commit();
  transaction = datastore.newTx();
  transaction->merge(newInterfaceTopPath + "/state", state);

  const map<Path, DatastoreDiff>& map = transaction->diff();
  Path diffKey(operStatus);
  MLOG(MINFO) << "FIRST: " << map.at(diffKey).before;
  EXPECT_EQ(
      map.at(diffKey).before["openconfig-interfaces:oper-status"], "DOWN");
  EXPECT_EQ(map.at(diffKey).after["openconfig-interfaces:oper-status"], "UP");
  EXPECT_EQ(map.size(), 1);
  EXPECT_EQ(map.at(diffKey).type, DatastoreDiffType::update);

  transaction->abort();
}

TEST_F(DatastoreTest, writeAndReadYdk) {
  shared_ptr<OpenconfigInterfaces> openconfigInterfaces = ydkInterfaces();
  shared_ptr<Datastore> datastore =
      std::make_shared<Datastore>(Datastore::operational(), schemaContext);
  BindingAwareDatastore bindingAwareDatastore(datastore, bindingCodec);
  const unique_ptr<BindingAwareDatastoreTransaction>& transaction =
      bindingAwareDatastore.newBindingTx();
  Path interfaces("/openconfig-interfaces:interfaces");

  transaction->overwrite(interfaces, openconfigInterfaces);

  const shared_ptr<OpenconfigInterfaces>& readIfaces =
      transaction->read<OpenconfigInterfaces>(interfaces);
  transaction->commit();

  EXPECT_EQ(readIfaces->interface.keys().size(), 1);
  shared_ptr<OpenconfigInterface> iface =
      std::static_pointer_cast<OpenconfigInterface>(readIfaces->interface[0]);
  string name = iface->name;
  EXPECT_EQ(name, "0/2");
  string description = iface->state->description;
  EXPECT_EQ(description, "dummy state");
}

TEST_F(DatastoreTest, readSubElementYdk) {
  shared_ptr<OpenconfigInterfaces> openconfigInterfaces = ydkInterfaces();
  shared_ptr<Datastore> datastore =
      std::make_shared<Datastore>(Datastore::operational(), schemaContext);
  BindingAwareDatastore bindingAwareDatastore(datastore, bindingCodec);

  const unique_ptr<BindingAwareDatastoreTransaction>& transaction =
      bindingAwareDatastore.newBindingTx();
  Path interfaces("/openconfig-interfaces:interfaces");

  transaction->overwrite(interfaces, openconfigInterfaces);

  const shared_ptr<OpenconfigConfig>& config =
      transaction->read<OpenconfigConfig>(
          "/openconfig-interfaces:interfaces/interface[name='0/2']/config");
  transaction->commit();

  string configDescription = config->description;
  EXPECT_EQ(configDescription, "this is a config description");
}

TEST_F(DatastoreTest, twoTransactionsAtTheSameTimeNotPermited) {
  shared_ptr<Datastore> datastore =
      std::make_shared<Datastore>(Datastore::operational(), schemaContext);
  BindingAwareDatastore bindingAwareDatastore(datastore, bindingCodec);

  const unique_ptr<DatastoreTransaction>& trans1 = datastore->newTx();
  EXPECT_THROW(bindingAwareDatastore.newBindingTx(), DatastoreException);
}

} // namespace cli
} // namespace test
} // namespace devmand
