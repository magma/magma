// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/channels/cli/datastore/BindingAwareDatastore.h>
#include <devmand/channels/cli/datastore/Datastore.h>

namespace devmand::channels::cli::datastore {

using devmand::channels::cli::datastore::Datastore;
BindingAwareDatastore::BindingAwareDatastore(
    const shared_ptr<Datastore> _datastore,
    const shared_ptr<BindingCodec> _codec)
    : datastore(_datastore), codec(_codec) {}

unique_ptr<BindingAwareDatastoreTransaction>
BindingAwareDatastore::newBindingTx() {
  std::unique_lock<std::mutex> lock(_mutex);
  return std::make_unique<BindingAwareDatastoreTransaction>(
      datastore->newTx(), codec);
}

} // namespace devmand::channels::cli::datastore
