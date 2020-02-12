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
#include <devmand/devices/cli/schema/Path.h>

namespace devmand::channels::cli::datastore {
using LeafVector = std::vector<pair<string, string>>;
using devmand::devices::cli::BindingCodec;
using devmand::devices::cli::Path;
using std::unique_ptr;

class BindingAwareDatastoreTransaction {
 private:
  unique_ptr<DatastoreTransaction> datastoreTransaction;
  shared_ptr<BindingCodec> codec;

 public:
  BindingAwareDatastoreTransaction(
      unique_ptr<DatastoreTransaction>&& _transaction,
      shared_ptr<BindingCodec> _codec);

 public:
  template <typename T>
  shared_ptr<T> read(Path path) {
    const shared_ptr<T>& ydkData = make_shared<T>();
    const dynamic& data = datastoreTransaction->read(path);
    return std::static_pointer_cast<T>(codec->decode(toJson(data), ydkData));
  }
  map<Path, DatastoreDiff> diff();
  void delete_(Path path);
  void overwrite(Path path, shared_ptr<Entity> entity);
  void merge(Path path, shared_ptr<Entity> entity);
  bool isValid();
  void print();
  void commit();
  void abort();
};

} // namespace devmand::channels::cli::datastore
