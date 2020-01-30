// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/schema/BindingContext.h>
#include <ydk/codec_provider.hpp>
#include <sstream>
#include <typeinfo>

namespace devmand {
namespace devices {
namespace cli {

using namespace ydk;
using namespace std;
using namespace ydk::path;

static void (*const noop)() = []() {};

BindingCodec::BindingCodec(Repository& repo, const string& schemaDir)
    : codecServiceProvider(repo, EncodingFormat::JSON), jsonSubtreeCodec() {
  codecServiceProvider.initialize(schemaDir, schemaDir, noop);
}

BindingContext::BindingContext(const Model& model)
    : repo(model.getDir(), ModelCachingOption::COMMON),
      bindingCodec(repo, model.getDir()) {}

BindingCodec& BindingContext::getCodec() {
  return bindingCodec;
}

static path::RootSchemaNode& rootSchema(
    CodecServiceProvider& codecServiceProvider) {
  return codecServiceProvider.get_root_schema_for_bundle("");
}

string BindingCodec::encode(Entity& entity) {
  lock_guard<std::mutex> lg(lock);

  try {
    RootSchemaNode& node = rootSchema(codecServiceProvider);
    return jsonSubtreeCodec.encode(entity, node, true);
  } catch (std::exception& e) {
    throw BindingSerializationException(entity, e.what());
  }
}

shared_ptr<Entity> BindingCodec::decode(
    const string& payload,
    shared_ptr<Entity> pointer) {
  lock_guard<std::mutex> lg(lock);

  try {
    return jsonSubtreeCodec.decode(payload, pointer);
  } catch (std::exception& e) {
    throw BindingSerializationException(pointer, e.what());
  }
}

} // namespace cli
} // namespace devices
} // namespace devmand
