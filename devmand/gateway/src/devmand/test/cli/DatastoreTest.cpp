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
#include <algorithm>

namespace devmand {
namespace test {
namespace cli {

using devmand::channels::cli::datastore::BindingAwareDatastore;
using devmand::channels::cli::datastore::BindingAwareDatastoreTransaction;
using devmand::channels::cli::datastore::Datastore;
using devmand::channels::cli::datastore::DatastoreDiff;
using devmand::channels::cli::datastore::DatastoreException;
using devmand::channels::cli::datastore::DatastoreTransaction;
using devmand::channels::cli::datastore::DiffPath;
using devmand::devices::cli::BindingCodec;
using devmand::devices::cli::SchemaContext;
using devmand::test::utils::cli::counterPath;
using devmand::test::utils::cli::ifaces02;
using devmand::test::utils::cli::interface02state;
using devmand::test::utils::cli::interface02TopPath;
using devmand::test::utils::cli::interfaceCounters;
using devmand::test::utils::cli::interfaceCountersWithKey;
using devmand::test::utils::cli::networkInstances;
using devmand::test::utils::cli::newInterface;
using devmand::test::utils::cli::newInterfaceTopPath;
using devmand::test::utils::cli::openconfigInterfacesInterfaces;
using devmand::test::utils::cli::operStatus;
using devmand::test::utils::cli::simpleInterfaces;
using devmand::test::utils::cli::simpleReplaceInterface;
using devmand::test::utils::cli::statePath;
using devmand::test::utils::cli::statePathWithKey;
using devmand::test::utils::cli::threeTrees;
using devmand::test::utils::cli::updated011Interface;
using devmand::test::utils::cli::vlans2;
using folly::parseJson;
using folly::toPrettyJson;
using std::to_string;
using std::unique_ptr;
using OpenconfigInterfaces = openconfig::openconfig_interfaces::Interfaces;
using OpenconfigInterface = OpenconfigInterfaces::Interface;
using OpenconfigConfig = OpenconfigInterface::Config;
using VlanType = openconfig::openconfig_vlan_types::VlanModeType;
using channels::cli::datastore::DiffResult;

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
  EXPECT_THROW(transaction->diff({}), DatastoreException);
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
  EXPECT_THROW(transaction->diff({}), DatastoreException);
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
  transaction->print();
  EXPECT_TRUE(toPrettyJson(transaction->read(interface03)) == "{}");

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

  EXPECT_THROW(transaction->isValid(), DatastoreException);
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

  vector<DiffPath> paths;
  Path p1(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface"
      "/openconfig-interfaces:state/openconfig-interfaces:counters");
  paths.emplace_back(p1, false);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:state/openconfig-interfaces:counters");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/2']/openconfig-interfaces:state/openconfig-interfaces:counters");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::update);
  EXPECT_EQ(
      multimap.begin()
          ->second.before["openconfig-interfaces:counters"]["out-errors"]
          .asString(),
      "0");
  EXPECT_EQ(
      multimap.begin()
          ->second.before["openconfig-interfaces:counters"]["out-discards"]
          .asString(),
      "0");
  EXPECT_EQ(
      multimap.begin()
          ->second.after["openconfig-interfaces:counters"]["out-errors"]
          .asString(),
      "777");
  EXPECT_EQ(
      multimap.begin()
          ->second.after["openconfig-interfaces:counters"]["out-discards"]
          .asString(),
      "17");
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

  vector<DiffPath> paths;
  Path p1("/openconfig-interfaces:interfaces/openconfig-interfaces:interface");
  paths.emplace_back(p1, true);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/2']");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::deleted);
  EXPECT_EQ(
      multimap.begin()
          ->second.before["openconfig-interfaces:interface"][0]["state"]["name"]
          .asString(),
      "0/2");
  EXPECT_ANY_THROW(
      multimap.begin()
          ->second.after["openconfig-interfaces:interface"][0]["state"]);
}

TEST_F(DatastoreTest, deleteSubtreeDiffNotifyChildren) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleInterfaces));

  transaction->commit();
  transaction = datastore.newTx();
  transaction->delete_("/openconfig-interfaces:interfaces");

  vector<DiffPath> paths;
  string configPath(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/config");
  Path p1(configPath);
  paths.emplace_back(p1, false);

  DiffResult diffResult = transaction->diff(paths);

  const std::multimap<Path, DatastoreDiff>& diffs = diffResult.diffs;

  auto it = diffs.equal_range(configPath.c_str());

  string handledConfigPath01 =
      "/openconfig-interfaces:interfaces/interface[name='0/1']/config";
  string handledConfigPath02 =
      "/openconfig-interfaces:interfaces/interface[name='0/2']/config";

  // check counters
  EXPECT_EQ(diffs.size(), 2);
  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(configPath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::deleted, itr->second.type);
    if (handledConfigPath01 == itr->second.keyedPath.str()) {
      EXPECT_EQ(handledConfigPath01, itr->second.keyedPath.str());
    } else {
      EXPECT_EQ(handledConfigPath02, itr->second.keyedPath.str());
    }
  }

  vector<string> expectedUnhandled;
  expectedUnhandled.emplace_back("/interfaces");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/1']");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/1']/state");
  expectedUnhandled.emplace_back(
      "/interfaces/interface[name='0/1']/state/counters");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/2']");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/2']/state");
  expectedUnhandled.emplace_back(
      "/interfaces/interface[name='0/2']/state/counters");

  vector<string> actuallyUnhandled;
  for (const auto& path : diffResult.unhandledDiffs) {
    actuallyUnhandled.emplace_back(path.unprefixAllSegments().str());
  }
  std::sort(expectedUnhandled.begin(), expectedUnhandled.end());
  std::sort(actuallyUnhandled.begin(), actuallyUnhandled.end());
  EXPECT_EQ(expectedUnhandled, actuallyUnhandled);
}

TEST_F(DatastoreTest, deleteCreateAndUpdateScenario) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleInterfaces));

  transaction->commit();
  transaction = datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleReplaceInterface));

  string ifaceCounterPath =
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/state/counters";
  vector<DiffPath> paths;
  Path p1(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/config");
  Path p2(ifaceCounterPath);
  paths.emplace_back(p1, false);
  paths.emplace_back(p2, false);

  DiffResult diffResult = transaction->diff(paths);
  const std::multimap<Path, DatastoreDiff>& diffs = diffResult.diffs;

  auto it = diffs.equal_range(ifaceCounterPath.c_str());

  for (auto itr = it.first; itr != it.second; ++itr) {
    if (itr->first.str() == ifaceCounterPath) {
      if (itr->second.type == DatastoreDiffType::update) {
        EXPECT_EQ(ifaceCounterPath, itr->first.str());
        EXPECT_EQ(
            "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/2']/openconfig-interfaces:state/openconfig-interfaces:counters",
            itr->second.keyedPath.str());

      } else {
        EXPECT_EQ(ifaceCounterPath, itr->first.str());
        EXPECT_EQ(
            "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/1']/state/counters",
            itr->second.keyedPath.str());
        EXPECT_EQ(DatastoreDiffType::deleted, itr->second.type);
      }
    }
  }

  it = diffs.equal_range(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/config");

  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(
        "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/config",
        itr->first.str());
    EXPECT_EQ(
        "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/1']/config",
        itr->second.keyedPath.str());
    EXPECT_EQ(DatastoreDiffType::deleted, itr->second.type);
  }

  vector<string> expectedUnhandled;
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/1']");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/1']/state");
  expectedUnhandled.emplace_back(
      "/interfaces/interface[name='0/2']/subinterfaces");
  expectedUnhandled.emplace_back(
      "/interfaces/interface[name='0/2']/subinterfaces/subinterface[index='0']");
  expectedUnhandled.emplace_back(
      "/interfaces/interface[name='0/2']/subinterfaces/subinterface[index='0']/config");

  vector<string> actuallyUnhandled;
  for (const auto& path : diffResult.unhandledDiffs) {
    actuallyUnhandled.emplace_back(path.unprefixAllSegments().str());
  }
  std::sort(expectedUnhandled.begin(), expectedUnhandled.end());
  std::sort(actuallyUnhandled.begin(), actuallyUnhandled.end());
  EXPECT_EQ(expectedUnhandled, actuallyUnhandled);
}

TEST_F(DatastoreTest, deleteSubtreeDiffDontNotifyParent) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleInterfaces));
  transaction->commit();
  transaction = datastore.newTx();
  transaction->delete_(Path(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/1']"
      "/openconfig-interfaces:state/openconfig-interfaces:counters"));

  vector<DiffPath> paths;
  Path p1(statePath);
  paths.emplace_back(p1, false);

  DiffResult diffResult = transaction->diff(paths);

  const std::multimap<Path, DatastoreDiff>& diffs = diffResult.diffs;

  EXPECT_EQ(diffs.size(), 0);
  EXPECT_EQ(
      "/interfaces/interface[name='0/1']/state/counters",
      transaction->diff(paths)
          .unhandledDiffs.front()
          .unprefixAllSegments()
          .str());
}

TEST_F(DatastoreTest, deleteSubtreeNotifyParentOnAsterix) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleInterfaces));
  transaction->commit();
  transaction = datastore.newTx();
  transaction->delete_(Path(
      "/openconfig-interfaces:interfaces/interface[name='0/1']/state/counters"));

  vector<DiffPath> paths;
  Path p1("/openconfig-interfaces:interfaces/interface/state");
  paths.emplace_back(p1, true);

  DiffResult diffResult = transaction->diff(paths);
  const std::multimap<Path, DatastoreDiff>& diffs = diffResult.diffs;

  EXPECT_EQ(
      diffs.begin()->first.str(),
      "/openconfig-interfaces:interfaces/interface/state");
  EXPECT_EQ(
      diffs.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/1']/openconfig-interfaces:state");
  EXPECT_EQ(diffs.begin()->second.type, DatastoreDiffType::deleted);
  EXPECT_EQ(diffResult.unhandledDiffs.size(), 0);
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

  vector<DiffPath> paths;
  Path p1("/openconfig-interfaces:interfaces/openconfig-interfaces:interface");
  paths.emplace_back(p1, false);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/85']");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::create);
  EXPECT_EQ(
      multimap.begin()
          ->second.after["openconfig-interfaces:interface"][0]["name"]
          .asString(),
      "0/85");
  EXPECT_EQ(toPrettyJson(multimap.begin()->second.before), "{}");
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

  vector<DiffPath> paths;
  Path p1(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:state");
  paths.emplace_back(p1, false);

  const std::multimap<Path, DatastoreDiff>& diffs =
      transaction->diff(paths).diffs;

  for (const auto& multi : diffs) {
    MLOG(MINFO) << "key: " << multi.first.str()
                << " handles:  " << multi.second.keyedPath.str()
                << " type: " << multi.second.type;
  }

  EXPECT_EQ(
      diffs.begin()->first.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:state");
  EXPECT_EQ(
      diffs.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/85']/openconfig-interfaces:state");
  EXPECT_EQ(diffs.begin()->second.type, DatastoreDiffType::update);
  EXPECT_EQ(
      diffs.begin()
          ->second.after["openconfig-interfaces:state"]["oper-status"]
          .asString(),
      "UP");
  EXPECT_EQ(
      diffs.begin()
          ->second.before["openconfig-interfaces:state"]["oper-status"]
          .asString(),
      "DOWN");

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

TEST_F(DatastoreTest, diffMultipleOperations) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));
  transaction->delete_(interface02TopPath + "/state");
  transaction->commit();
  transaction = datastore.newTx();
  dynamic interface02 = transaction->read(interface02TopPath);
  interface02["openconfig-interfaces:interface"][0]["state"] =
      folly::dynamic::object();
  interface02["openconfig-interfaces:interface"][0]["state"]["counters"] =
      folly::dynamic::object();
  interface02["openconfig-interfaces:interface"][0]["state"]["counters"]
             ["in-errors"] = 7;
  interface02["openconfig-interfaces:interface"][0]["state"]["admin-status"] =
      "DOWN";
  interface02["openconfig-interfaces:interface"][0]["config"] =
      folly::dynamic::object();
  interface02["openconfig-interfaces:interface"][0]["config"]["mtu"] = 1400;
  interface02["openconfig-interfaces:interface"][0]["config"]["enabled"] =
      false;
  transaction->merge(interface02TopPath, interface02);
  vector<DiffPath> paths;
  Path p1(statePath);
  Path p2(counterPath);
  paths.emplace_back(p1, false);
  paths.emplace_back(p2, false);

  const std::multimap<Path, DatastoreDiff>& diffs =
      transaction->diff(paths).diffs;

  auto it = diffs.equal_range(statePath.c_str());

  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(statePath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::create, itr->second.type);
    EXPECT_EQ(
        "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/2']/openconfig-interfaces:state",
        itr->second.keyedPath.str());
  }

  it = diffs.equal_range(counterPath.c_str());

  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(counterPath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::create, itr->second.type);
    EXPECT_EQ(
        "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/2']/openconfig-interfaces:state/counters",
        itr->second.keyedPath.str());
  }
}

TEST_F(DatastoreTest, createParentListenChildDiff) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleInterfaces));

  vector<DiffPath> paths;
  string configPath = "/openconfig-interfaces:interfaces/interface/config";
  string countersPath =
      "/openconfig-interfaces:interfaces/interface/openconfig-interfaces:state/counters";
  Path p1(configPath);
  Path p2(countersPath);
  paths.emplace_back(p1, false);
  paths.emplace_back(p2, false);

  DiffResult diffResult = transaction->diff(paths);
  const std::multimap<Path, DatastoreDiff>& diffs = diffResult.diffs;

  auto it = diffs.equal_range(countersPath.c_str());

  string handledCounterPath01 =
      "/openconfig-interfaces:interfaces/interface[name='0/1']/state/counters";
  string handledCounterPath02 =
      "/openconfig-interfaces:interfaces/interface[name='0/2']/state/counters";

  // check counters
  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(countersPath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::create, itr->second.type);
    if (handledCounterPath01 == itr->second.keyedPath.str()) {
      EXPECT_EQ(handledCounterPath01, itr->second.keyedPath.str());
    } else {
      EXPECT_EQ(handledCounterPath02, itr->second.keyedPath.str());
    }
  }

  // check config
  string handledConfigPath01 =
      "/openconfig-interfaces:interfaces/interface[name='0/1']/config";
  string handledConfigPath02 =
      "/openconfig-interfaces:interfaces/interface[name='0/2']/config";

  it = diffs.equal_range(configPath.c_str());

  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(configPath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::create, itr->second.type);
    if (handledConfigPath01 == itr->second.keyedPath.str()) {
      EXPECT_EQ(handledConfigPath01, itr->second.keyedPath.str());
    } else {
      EXPECT_EQ(handledConfigPath02, itr->second.keyedPath.str());
    }
  }

  // check unhandled
  vector<string> expectedUnhandled;
  expectedUnhandled.emplace_back("/interfaces");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/1']");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/1']/state");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/2']");
  expectedUnhandled.emplace_back("/interfaces/interface[name='0/2']/state");

  vector<string> actuallyUnhandled;
  for (const auto& path : diffResult.unhandledDiffs) {
    actuallyUnhandled.emplace_back(path.unprefixAllSegments().str());
  }
  std::sort(expectedUnhandled.begin(), expectedUnhandled.end());
  std::sort(actuallyUnhandled.begin(), actuallyUnhandled.end());
  EXPECT_EQ(expectedUnhandled, actuallyUnhandled);
}

TEST_F(DatastoreTest, simpleCreateDiff) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(simpleInterfaces));

  vector<DiffPath> paths;
  Path p1(statePath);
  Path p2(counterPath);
  paths.emplace_back(p1, false);
  paths.emplace_back(p2, false);

  const std::multimap<Path, DatastoreDiff>& diffs =
      transaction->diff(paths).diffs;

  auto it = diffs.equal_range(counterPath.c_str());

  string handledCounterPath01 =
      "/openconfig-interfaces:interfaces/interface[name='0/1']/state/counters";
  string handledCounterPath02 =
      "/openconfig-interfaces:interfaces/interface[name='0/2']/state/counters";
  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(counterPath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::create, itr->second.type);
    if (handledCounterPath01 == itr->second.keyedPath.str()) {
      EXPECT_EQ(handledCounterPath01, itr->second.keyedPath.str());
    } else {
      EXPECT_EQ(handledCounterPath02, itr->second.keyedPath.str());
    }
  }

  string handledStatePath01 =
      "/openconfig-interfaces:interfaces/interface[name='0/1']/state";
  string handledStatePath02 =
      "/openconfig-interfaces:interfaces/interface[name='0/2']/state";

  it = diffs.equal_range(statePath.c_str());

  for (auto itr = it.first; itr != it.second; ++itr) {
    EXPECT_EQ(statePath, itr->first.str());
    EXPECT_EQ(DatastoreDiffType::create, itr->second.type);
    if (handledStatePath01 == itr->second.keyedPath.str()) {
      EXPECT_EQ(handledStatePath01, itr->second.keyedPath.str());
    } else {
      EXPECT_EQ(handledStatePath02, itr->second.keyedPath.str());
    }
  }
}

TEST_F(DatastoreTest, diffDeleteOperation) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));
  transaction->commit();
  transaction = datastore.newTx();
  transaction->delete_(interface02TopPath);

  vector<DiffPath> paths;
  Path p1(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/state");
  paths.emplace_back(p1, false);

  const std::multimap<Path, DatastoreDiff>& diffs =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      diffs.begin()->first.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/state");
  EXPECT_EQ(
      diffs.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/2']/state");
  EXPECT_EQ(DatastoreDiffType::deleted, diffs.begin()->second.type);
}

TEST_F(DatastoreTest, twoIdenpendentTreesDiffUpdateTest) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(networkInstances));
  transaction->commit();
  transaction = datastore.newTx();
  transaction->merge(
      Path("/openconfig-interfaces:interfaces/interface[name='0/11']"),
      parseJson(updated011Interface));

  vector<DiffPath> paths;
  Path p1(
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:config");
  paths.emplace_back(p1, true);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface/openconfig-interfaces:config");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-interfaces:interfaces/openconfig-interfaces:interface[name='0/11']/openconfig-interfaces:config");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::update);
}

TEST_F(
    DatastoreTest,
    threeIndenpendentTreesDiffDeleteAndNotifyParentBecauseAsterixTest) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(threeTrees));
  transaction->commit();
  transaction = datastore.newTx();
  transaction->delete_(Path(
      "/openconfig-network-instance:network-instances/network-instance['default']/config"));

  vector<DiffPath> paths;
  Path p1(
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance");
  paths.emplace_back(p1, true);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance[name='default']");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::deleted);
}

TEST_F(DatastoreTest, diff2changes) {
  shared_ptr<Datastore> datastore =
      std::make_shared<Datastore>(Datastore::operational(), schemaContext);
  const unique_ptr<DatastoreTransaction>& transaction = datastore->newTx();

  transaction->overwrite(Path("/"), parseJson(openconfigInterfacesInterfaces));

  transaction->commit();
  const unique_ptr<DatastoreTransaction>& transaction2 = datastore->newTx();
  dynamic counters = transaction2->read(interface02TopPath + "/state/counters");
  counters["openconfig-interfaces:counters"]["out-errors"] = "777";
  counters["openconfig-interfaces:counters"]["out-discards"] = "17";
  transaction2->merge(interface02TopPath + "/state/counters", counters);
  vector<DiffPath> paths;
  paths.emplace_back(Path(interfaceCounters), false);
  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction2->diff(paths).diffs;
  for (const auto& multi : multimap) {
    EXPECT_EQ(multi.first.str(), interfaceCounters);
    if (interfaceCountersWithKey + "/openconfig-interfaces:out-discards" ==
        multi.second.keyedPath.str()) {
      EXPECT_EQ(
          multi.second.keyedPath.str(),
          interfaceCountersWithKey + "/openconfig-interfaces:out-discards");
    } else {
      EXPECT_EQ(multi.second.keyedPath.str(), interfaceCountersWithKey);
    }

    MLOG(MINFO) << "key: " << multi.first
                << " handles: " << multi.second.keyedPath;
  }
}

TEST_F(DatastoreTest, threeIndenpendentTreesUpdateReadTest) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(threeTrees));
  transaction->commit();
  transaction = datastore.newTx();
  const Path& vlanPath = Path(
      "/openconfig-network-instance:network-instances/network-instance[name='default']/vlans");
  dynamic vlans = transaction->read(vlanPath);

  vlans["openconfig-network-instance:vlans"]["vlan"][0]["state"]["status"] =
      "SUSPENDED";
  vlans["openconfig-network-instance:vlans"]["vlan"][0]["config"]["status"] =
      "SUSPENDED";
  transaction->overwrite(vlanPath, vlans);
  vector<DiffPath> paths;
  Path p1(
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance"
      "/openconfig-network-instance:vlans/openconfig-network-instance:vlan/openconfig-network-instance:state");
  paths.emplace_back(p1, false);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance"
      "/openconfig-network-instance:vlans/openconfig-network-instance:vlan/openconfig-network-instance:state");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance[name='default']/openconfig-network-instance:vlans"
      "/openconfig-network-instance:vlan[vlan-id='1']/openconfig-network-instance:state");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::update);
}

TEST_F(DatastoreTest, threeIndenpendentTreesCreateReadTest) {
  Datastore datastore(Datastore::operational(), schemaContext);
  unique_ptr<channels::cli::datastore::DatastoreTransaction> transaction =
      datastore.newTx();
  transaction->overwrite(Path("/"), parseJson(threeTrees));
  transaction->commit();
  transaction = datastore.newTx();
  const Path& vlanPath = Path(
      "/openconfig-network-instance:network-instances/network-instance[name='default']/vlans");
  dynamic vlans = transaction->read(vlanPath);

  vlans["openconfig-network-instance:vlans"]["vlan"].push_back(
      dynamic::object()); // [4] =;
  vlans["openconfig-network-instance:vlans"]["vlan"][4]["vlan-id"] = 666;

  vlans["openconfig-network-instance:vlans"]["vlan"][4]["state"] =
      dynamic::object();
  vlans["openconfig-network-instance:vlans"]["vlan"][4]["state"]["status"] =
      "SUSPENDED";
  vlans["openconfig-network-instance:vlans"]["vlan"][4]["state"]["vlan-id"] =
      666;

  vlans["openconfig-network-instance:vlans"]["vlan"][4]["config"] =
      dynamic::object();
  vlans["openconfig-network-instance:vlans"]["vlan"][4]["config"]["status"] =
      "SUSPENDED";
  vlans["openconfig-network-instance:vlans"]["vlan"][4]["config"]["vlan-id"] =
      666;

  transaction->overwrite(vlanPath, vlans);
  vector<DiffPath> paths;
  Path p1(
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance/openconfig-network-instance:vlans"
      "/openconfig-network-instance:vlan");
  paths.emplace_back(p1, false);

  const std::multimap<Path, DatastoreDiff>& multimap =
      transaction->diff(paths).diffs;

  EXPECT_EQ(
      multimap.begin()->first.str(),
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance"
      "/openconfig-network-instance:vlans/openconfig-network-instance:vlan");
  EXPECT_EQ(
      multimap.begin()->second.keyedPath.str(),
      "/openconfig-network-instance:network-instances/openconfig-network-instance:network-instance[name='default']"
      "/openconfig-network-instance:vlans/openconfig-network-instance:vlan[vlan-id='666']");
  EXPECT_EQ(multimap.begin()->second.type, DatastoreDiffType::create);
  EXPECT_EQ(
      multimap.begin()->second.after["openconfig-network-instance:vlan"][0]
                                    ["config"]["vlan-id"],
      666);
  EXPECT_EQ(
      multimap.begin()
          ->second
          .after["openconfig-network-instance:vlan"][0]["config"]["status"]
          .asString(),
      "SUSPENDED");
  EXPECT_EQ(
      multimap.begin()->second.after["openconfig-network-instance:vlan"][0]
                                    ["state"]["vlan-id"],
      666);
  EXPECT_EQ(
      multimap.begin()
          ->second
          .after["openconfig-network-instance:vlan"][0]["state"]["status"]
          .asString(),
      "SUSPENDED");
  EXPECT_EQ(
      multimap.begin()
          ->second.after["openconfig-network-instance:vlan"][0]["vlan-id"],
      666);
  EXPECT_EQ(toPrettyJson(multimap.begin()->second.before), "{}");
}

} // namespace cli
} // namespace test
} // namespace devmand
