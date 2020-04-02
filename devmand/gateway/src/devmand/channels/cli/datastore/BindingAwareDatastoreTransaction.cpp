// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/algorithm/string.hpp>
#include <devmand/channels/cli/datastore/BindingAwareDatastoreTransaction.h>

namespace devmand::channels::cli::datastore {

DiffResult BindingAwareDatastoreTransaction::diff(
    vector<DiffPath> registeredPaths) {
  return datastoreTransaction->diff(registeredPaths);
}

void BindingAwareDatastoreTransaction::delete_(Path path) {
  datastoreTransaction->delete_(path);
}

void BindingAwareDatastoreTransaction::overwrite(
    Path path,
    shared_ptr<Entity> entity) {
  datastoreTransaction->overwrite(path, codec->toDom(path, *entity));
}

void BindingAwareDatastoreTransaction::merge(
    Path path,
    shared_ptr<Entity> entity) {
  datastoreTransaction->merge(path, codec->toDom(path, *entity));
}

void BindingAwareDatastoreTransaction::commit() {
  datastoreTransaction->commit();
}

BindingAwareDatastoreTransaction::BindingAwareDatastoreTransaction(
    unique_ptr<DatastoreTransaction>&& _transaction,
    shared_ptr<BindingCodec> _codec)
    : datastoreTransaction(
          std::forward<unique_ptr<DatastoreTransaction>>(_transaction)),
      codec(_codec) {}

void BindingAwareDatastoreTransaction::isValid() {
  return datastoreTransaction->isValid();
}

void BindingAwareDatastoreTransaction::abort() {
  datastoreTransaction->abort();
}

void BindingAwareDatastoreTransaction::print() {
  datastoreTransaction->print();
}

} // namespace devmand::channels::cli::datastore
