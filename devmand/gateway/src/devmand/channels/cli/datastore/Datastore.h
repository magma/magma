// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/datastore/DatastoreState.h>
#include <devmand/channels/cli/datastore/DatastoreTransaction.h>
#include <devmand/devices/cli/schema/BindingContext.h>
#include <devmand/devices/cli/schema/ModelRegistry.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <libyang/libyang.h>

namespace devmand::channels::cli::datastore {
using devmand::channels::cli::datastore::DatastoreState;
using devmand::channels::cli::datastore::DatastoreTransaction;
using devmand::channels::cli::datastore::DatastoreType;
using devmand::devices::cli::BindingCodec;
using devmand::devices::cli::ModelRegistry;
using devmand::devices::cli::SchemaContext;
using std::shared_ptr;
using std::unique_ptr;

class Datastore {
 private:
  shared_ptr<DatastoreState> datastoreState;
  SchemaContext& schemaContext;
  std::mutex _mutex;

  void checkIfTransactionRunning();
  void setTransactionRunning();

 public:
  static DatastoreType operational();
  static DatastoreType config();

  Datastore(DatastoreType _type, SchemaContext& _schemaContext);

  unique_ptr<DatastoreTransaction>
  newTx(); // operations on transaction are NOT thread-safe
};

class DatastoreException : public runtime_error {
 public:
  DatastoreException(const string& _msg) : runtime_error(_msg){};
};
} // namespace devmand::channels::cli::datastore
