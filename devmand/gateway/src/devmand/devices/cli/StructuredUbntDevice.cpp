// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/IoConfigurationBuilder.h>
#include <devmand/devices/Datastore.h>
#include <devmand/devices/cli/DeviceType.h>
#include <devmand/devices/cli/StructuredUbntDevice.h>
#include <devmand/devices/cli/UbntStpPlugin.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/translation/PluginRegistry.h>
#include <devmand/devices/cli/translation/ReaderRegistry.h>
#include <devmand/devices/cli/translation/WriterRegistry.h>
#include <folly/json.h>
#include <memory>

namespace devmand {
namespace devices {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;
using namespace folly;

std::unique_ptr<devices::Device> StructuredUbntDevice::createDevice(
    Application& app,
    const cartography::DeviceConfig& deviceConfig) {
  return createDeviceWithEngine(app, deviceConfig, app.getCliEngine());
}

unique_ptr<devices::Device> StructuredUbntDevice::createDeviceWithEngine(
    Application& app,
    const cartography::DeviceConfig& deviceConfig,
    Engine& engine) {
  DeviceType deviceType = DeviceType::create(deviceConfig);
  shared_ptr<CliFlavour> cliFlavour = engine.getCliFlavour(deviceType);
  IoConfigurationBuilder ioConfigurationBuilder(
      deviceConfig, engine, cliFlavour);
  auto cmdCache = ReadCachingCli::createCache();
  auto treeCache = make_shared<TreeCache>(
      ioConfigurationBuilder.getConnectionParameters()->flavour);
  const std::shared_ptr<Channel>& channel = std::make_shared<Channel>(
      deviceConfig.id, ioConfigurationBuilder.createAll(cmdCache, treeCache));

  shared_ptr<DeviceContext> deviceCtx = engine.getDeviceContext(deviceType);
  return std::make_unique<StructuredUbntDevice>(
      app,
      deviceConfig.id,
      deviceConfig.readonly,
      channel,
      engine.getModelRegistry(),
      engine.getReaderRegistry(deviceCtx),
      engine.getWriterRegistry(deviceCtx),
      cmdCache,
      treeCache);
}

StructuredUbntDevice::StructuredUbntDevice(
    Application& application,
    const Id id_,
    bool readonly_,
    const shared_ptr<Channel> _channel,
    const std::shared_ptr<ModelRegistry> _mreg,
    std::unique_ptr<ReaderRegistry>&& _rReg,
    std::unique_ptr<WriterRegistry>&& _wReg,
    const shared_ptr<CliCache> _cmdCache,
    const shared_ptr<TreeCache> _treeCache)
    : Device(application, id_, readonly_),
      channel(_channel),
      cmdCache(_cmdCache),
      treeCache(_treeCache),
      mreg(_mreg),
      rReg(forward<unique_ptr<ReaderRegistry>>(_rReg)),
      wReg(forward<unique_ptr<WriterRegistry>>(_wReg)) {}
//      ,configCache(make_unique<devmand::channels::cli::datastore::Datastore>(
//          DatastoreType::config,
//          app.getCliEngine().getModelRegistry()->getSchemaContext(
//              Model::OPENCONFIG_2_4_3))),
//      diffPaths(vector<DiffPath>()) {
//  for (auto& regPath : wReg->getWriterPaths()) {
//    // TODO support subtree writers
//    diffPaths.push_back(DiffPath(regPath, false));
//  }
//}

// void StructuredUbntDevice::reconcile(DeviceAccess& access) {
//  MLOG(MINFO) << "[" << id << "] "
//              << "Reconciling";
//  auto reconcileTx = configCache->newTx();
//  try {
//    dynamic reconciledData = rReg->readConfiguration("/", access).get();
//    MLOG(MINFO) << "[" << id << "] "
//                << "Reconciled with: " << reconciledData;
//    reconcileTx->overwrite("/", reconciledData);
//  } catch (DatastoreException& e) {
//    reconcileTx->abort();
//    throw runtime_error(
//        "[" + id + "] Invalid configuration reconciled due to: " + e.what());
//  } catch (runtime_error& e) {
//    reconcileTx->abort();
//    throw runtime_error("[" + id + "] Unable to reconcile due to: " +
//    e.what());
//  }
//  if (!reconcileTx->isValid()) {
//    reconcileTx->abort();
//    throw runtime_error(
//        "[" + id +
//        "] Unable to reconcile due to: Reconciled configuration is not
//        valid");
//  }
//  reconcileTx->commit();
//}

void StructuredUbntDevice::setIntendedDatastore(const dynamic& config) {
  (void)config;
  //  MLOG(MINFO) << "[" << id << "] "
  //              << "Writing config";
  //
  //  // Reset cache
  //  cmdCache->wlock()->clear();
  //  treeCache->clear(); // FIXME this is not threadsafe
  //
  //  DeviceAccess access = DeviceAccess(channel, id, getCPUExecutor());
  //  reconcile(access);
  //
  //  // Apply new config
  //  auto tx = configCache->newTx();
  //  try {
  //    tx->merge("/", config);
  //  } catch (DatastoreException& e) {
  //    tx->abort();
  //    throw runtime_error("Invalid configuration for device: " + id);
  //  }
  //
  //  if (!tx->isValid()) {
  //    tx->abort();
  //    throw runtime_error("Invalid configuration for device: " + id);
  //  }
  //
  //  MLOG(MINFO) << "[" << id << "] "
  //              << "Calculating diff";
  //  auto diffResult = tx->diff(diffPaths);
  //  auto diff = diffResult.diffs;
  //
  //  if (!diffResult.unhandledDiffs.empty()) {
  //    stringstream s;
  //    for (auto& unhandledPath : diffResult.unhandledDiffs) {
  //      s << unhandledPath;
  //      s << ",";
  //    }
  //    MLOG(MWARNING) << "[" << id << "] "
  //                   << "Unhandled paths detected from diff: "
  //                   << s.str();
  //  }
  //
  //  if (diff.size() == 0) {
  //    MLOG(MINFO) << "[" << id << "] "
  //                << "No updates detected";
  //    tx->abort();
  //    return;
  //  }
  //
  //  MLOG(MINFO) << "[" << id << "] "
  //              << "Submitting to device";
  //  wReg->write(diff, access);
  //  tx->commit();
  //
  //  MLOG(MINFO) << "[" << id << "] "
  //              << "Config written successfully";
}

shared_ptr<Datastore> StructuredUbntDevice::getOperationalDatastore() {
  MLOG(MINFO) << "[" << id << "] "
              << "Retrieving state";

  // Reset cache
  cmdCache->wlock()->clear();
  treeCache->clear(); // FIXME this is not threadsafe

  auto state = Datastore::make(*reinterpret_cast<MetricSink*>(&app), getId());
  state->setStatus(true);
  //  return state;
  DeviceAccess access = DeviceAccess(channel, id, getCPUExecutor());

  state->addRequest(
      rReg->readState(Path::ROOT, access)
          .thenValue([state](auto v) {
            state->update(
                [&v](auto& lockedState) { lockedState.merge_patch(v); });
          })
          .thenError(
              // TODO unify with ReconnectingCli
              // (DisconnectedException+CommandExecutionException)
              folly::tag_t<DisconnectedException>{},
              [state,
               id = this->id](DisconnectedException const& e) -> Future<Unit> {
                state->setStatus(false);
                throw e;
              })
          .thenError(
              folly::tag_t<std::exception>{},
              [state, id = this->id](std::exception const& e) {
                MLOG(MWARNING) << "[" << id << "] "
                               << "Retrieving state failed: " << e.what();
                state->addError(e.what());
              }));

  return state;
}
} // namespace cli
} // namespace devices
} // namespace devmand
