// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/schema/ModelRegistry.h>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;

ModelRegistry::ModelRegistry() {
  // Set plugin directory for YDK libyang according to:
  // https://github.com/CESNET/libyang/blob/c38295963669219b7aad2618b9f1dd31fa667caa/FAQ.md
  // and CMakeLists.ydk
  setenv("LIBYANG_EXTENSIONS_PLUGINS_DIR", LIBYANG_PLUGINS_DIR, false);
}

ModelRegistry::~ModelRegistry() {
  bindingCache.clear();
}

ulong ModelRegistry::bindingCacheSize() {
  lock_guard<std::mutex> lg(lock);

  return bindingCache.size();
}

BindingContext& ModelRegistry::getBindingContext(const Model& model) {
  const SchemaContext& schemaCtx = getSchemaContext(model);
  lock_guard<std::mutex> lg(lock);

  auto it = bindingCache.find(model.getDir());
  if (it != bindingCache.end()) {
    return it->second;
  } else {
    auto pair = bindingCache.emplace(
        piecewise_construct,
        forward_as_tuple(model.getDir()),
        forward_as_tuple(model, schemaCtx));
    return pair.first->second;
  }
}

SchemaContext& ModelRegistry::getSchemaContext(const Model& model) {
  lock_guard<std::mutex> lg(lock);

  auto it = schemaCache.find(model.getDir());
  if (it != schemaCache.end()) {
    return it->second;
  } else {
    auto pair = schemaCache.emplace(model.getDir(), model);
    return pair.first->second;
  }
}

} // namespace cli
} // namespace devices
} // namespace devmand
