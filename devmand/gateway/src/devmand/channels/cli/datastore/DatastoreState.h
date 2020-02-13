// Copyright (c) 2020-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <libyang/libyang.h>
#include <atomic>

namespace devmand::channels::cli::datastore {
using std::atomic_bool;

enum DatastoreType { config, operational };

struct DatastoreState {
  atomic_bool transactionUnderway = ATOMIC_VAR_INIT(false);
  llly_ctx* ctx = nullptr;
  lllyd_node* root = nullptr;
  DatastoreType type;

  virtual ~DatastoreState() {
    if (root != nullptr) {
      lllyd_free(root);
    }
  }

 public:
  DatastoreState(llly_ctx* _ctx, DatastoreType _type)
      : ctx(_ctx), type(_type) {}

  bool isEmpty() {
    return root == nullptr || ctx == nullptr;
  }
};

typedef struct DatastoreState DatastoreState;
} // namespace devmand::channels::cli::datastore
