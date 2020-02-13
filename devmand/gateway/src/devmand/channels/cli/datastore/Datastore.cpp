// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/datastore/Datastore.h>
#include <devmand/channels/cli/engine/Engine.h>
#include <devmand/devices/cli/schema/Model.h>

namespace devmand::channels::cli::datastore {
using devmand::channels::cli::Engine;
using devmand::channels::cli::datastore::DatastoreException;
using devmand::devices::cli::SchemaContext;
using std::make_unique;

Datastore::Datastore(DatastoreType type, SchemaContext& _schemaContext)
    : schemaContext(_schemaContext) {
  llly_ctx* pLyCtx = schemaContext.getLyContext();
  datastoreState = make_shared<DatastoreState>(pLyCtx, type);
}

unique_ptr<DatastoreTransaction> Datastore::newTx() {
  unique_lock<mutex> lock(_mutex);
  checkIfTransactionRunning();
  setTransactionRunning();
  return make_unique<DatastoreTransaction>(datastoreState);
}

void Datastore::checkIfTransactionRunning() {
  if (datastoreState->transactionUnderway) {
    DatastoreException ex(
        "Transaction in datastore already running, only 1 at a time permitted");
    MLOG(MWARNING) << ex.what();
    throw ex;
  }
}

void Datastore::setTransactionRunning() {
  datastoreState->transactionUnderway.store(true);
}

DatastoreType Datastore::operational() {
  return DatastoreType::operational;
}

DatastoreType Datastore::config() {
  return DatastoreType::config;
}

} // namespace devmand::channels::cli::datastore
