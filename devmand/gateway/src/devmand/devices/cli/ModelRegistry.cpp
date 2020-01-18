// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include "ModelRegistry.h"
#include <ydk/codec_provider.hpp>
#include <sstream>
#include <typeinfo>

namespace devmand {
namespace devices {
namespace cli {

using namespace ydk;
using namespace std;
using namespace ydk::path;

const Model Model::OPENCONFIG_0_1_6 = Model("/usr/share/openconfig@0.1.6");
const Model Model::IETF_0_1_5 = Model("/usr/share/ietf@0.1.5");

static void (*const noop)() = []() {};

Bundle::Bundle(const Model& model)
    : repo(Repository(model.getDir(), ModelCachingOption::COMMON)),
      codecServiceProvider(CodecServiceProvider(repo, EncodingFormat::JSON)),
      jsonSubtreeCodec(JsonSubtreeCodec()) {
  codecServiceProvider.initialize(model.getDir(), model.getDir(), noop);
}

static path::RootSchemaNode& rootSchema(
    CodecServiceProvider& codecServiceProvider) {
  return codecServiceProvider.get_root_schema_for_bundle("");
}

string Bundle::encode(Entity& entity) {
  lock_guard<std::mutex> lg(lock);

  try {
    RootSchemaNode& node = rootSchema(codecServiceProvider);
    return jsonSubtreeCodec.encode(entity, node, true);
  } catch (std::exception& e) {
    throw SerializationException(entity, e.what());
  }
}

shared_ptr<Entity> Bundle::decode(
    const string& payload,
    shared_ptr<Entity> pointer) {
  lock_guard<std::mutex> lg(lock);

  try {
    return jsonSubtreeCodec.decode(payload, pointer);
  } catch (std::exception& e) {
    throw SerializationException(pointer, e.what());
  }
}

ModelRegistry::ModelRegistry() {
  // Set plugin directory for libyang according to:
  // https://github.com/CESNET/libyang/blob/c38295963669219b7aad2618b9f1dd31fa667caa/FAQ.md
  // and CMakeLists.ydk
  setenv("LIBYANG_EXTENSIONS_PLUGINS_DIR", LIBYANG_PLUGINS_DIR, false);
}

ModelRegistry::~ModelRegistry() {
  cache.clear();
}

ulong ModelRegistry::cacheSize() {
  lock_guard<std::mutex> lg(lock);

  return cache.size();
}

Bundle& ModelRegistry::getBundle(const Model& model) {
  lock_guard<std::mutex> lg(lock);

  auto it = cache.find(model.getDir());
  if (it != cache.end()) {
    return it->second;
  } else {
    auto pair = cache.emplace(model.getDir(), model);
    return pair.first->second;
  }
}

} // namespace cli
} // namespace devices
} // namespace devmand
