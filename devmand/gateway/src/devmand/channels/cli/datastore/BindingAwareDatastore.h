// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once
#include <devmand/channels/cli/datastore/BindingAwareDatastoreTransaction.h>
#include <devmand/channels/cli/datastore/Datastore.h>
#include <devmand/devices/cli/schema/BindingContext.h>

namespace devmand::channels::cli::datastore {

using devmand::channels::cli::datastore::BindingAwareDatastoreTransaction;
using devmand::channels::cli::datastore::Datastore;
using devmand::channels::cli::datastore::DatastoreTransaction;
using devmand::devices::cli::BindingCodec;

class BindingAwareDatastore {
 private:
  shared_ptr<Datastore> datastore;
  shared_ptr<BindingCodec> codec;
  std::mutex _mutex;

 public:
  unique_ptr<BindingAwareDatastoreTransaction>
  newBindingTx(); // operations on transaction are NOT thread-safe
  BindingAwareDatastore(
      const shared_ptr<Datastore>,
      const shared_ptr<BindingCodec> codec);
};

} // namespace devmand::channels::cli::datastore
