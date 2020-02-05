// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/BindingContext.h>
#include <devmand/devices/cli/schema/Model.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <mutex>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace devmand::devices::cli;

class ModelRegistry {
 public:
  ModelRegistry();
  ~ModelRegistry();
  ModelRegistry(const ModelRegistry&) = delete;
  ModelRegistry& operator=(const ModelRegistry&) = delete;
  ModelRegistry(ModelRegistry&&) = delete;
  ModelRegistry& operator=(ModelRegistry&&) = delete;

 public:
  // Contexts are cached internally, returning just a reference
  BindingContext& getBindingContext(const Model& model);
  SchemaContext& getSchemaContext(const Model& model);

  // Visible for testing
  ulong bindingCacheSize();
  ulong schemaCacheSize();

 private:
  mutex lock; // A bundle is expected to be shared, protect it
  map<string, BindingContext> bindingCache;
  map<string, SchemaContext> schemaCache;
};

} // namespace cli
} // namespace devices
} // namespace devmand
