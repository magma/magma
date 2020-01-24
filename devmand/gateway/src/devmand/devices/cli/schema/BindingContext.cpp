// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/devices/cli/schema/BindingContext.h>
#include <folly/json.h>
#include <ydk/codec_provider.hpp>
#include <regex>
#include <sstream>
#include <typeinfo>

namespace devmand {
namespace devices {
namespace cli {

using namespace ydk;
using namespace std;
using namespace ydk::path;

static void (*const noop)() = []() {};

BindingCodec::BindingCodec(
    Repository& repo,
    const string& schemaDir,
    const SchemaContext& _schemaCtx)
    : codecServiceProvider(repo, EncodingFormat::JSON),
      jsonSubtreeCodec(),
      schemaCtx(_schemaCtx) {
  codecServiceProvider.initialize(schemaDir, schemaDir, noop);
}

BindingContext::BindingContext(
    const Model& model,
    const SchemaContext& _schemaCtx)
    : repo(model.getDir(), ModelCachingOption::COMMON),
      bindingCodec(repo, model.getDir(), _schemaCtx) {}

BindingCodec& BindingContext::getCodec() {
  return bindingCodec;
}

static path::RootSchemaNode& rootSchema(
    CodecServiceProvider& codecServiceProvider) {
  return codecServiceProvider.get_root_schema_for_bundle("");
}

static const regex LEAF_SET_WITH_KEY = regex("([^[]+)\\[\\.=\"(.+)\"\\]");

static dynamic leafToDynamic(
    Path leafPath,
    string leafAsString,
    const SchemaContext& schemaCtx) {
  vector<LLLY_DATA_TYPE> type = schemaCtx.leafType(leafPath.unkeyed());

  // Unions
  if (type.size() > 1) {
    for (auto unionType : type) {
      if (LLLY_TYPE_STRING == unionType) {
        return leafAsString;
      } else if (LLLY_TYPE_UNION < unionType) {
        return stoll(leafAsString);
      } else if (LLLY_TYPE_BOOL == unionType) {
        return "true" == leafAsString;
      } else {
        continue;
      }
    }
  }

  // TODO double ?
  if (LLLY_TYPE_STRING == type[0]) {
    return leafAsString;
  } else if (LLLY_TYPE_UNION < type[0]) {
    return stoll(leafAsString);
  } else if (LLLY_TYPE_BOOL == type[0]) {
    return "true" == leafAsString;
  } else {
    // treat all other types as string
    return leafAsString;
  }
}

static dynamic
toDynamicRecurs(Path path, Entity& entity, const SchemaContext& schemaCtx) {
  dynamic asDynamic = dynamic::object();

  // leaves
  for (auto& leaf : entity.get_name_leaf_data()) {
    smatch leafSetMatch;
    // leaf-set entry
    if (regex_match(leaf.first, leafSetMatch, LEAF_SET_WITH_KEY)) {
      string name = leafSetMatch[1];
      Path leafPath = path.getChild(name);
      string value = leafSetMatch[2];
      if (!asDynamic[name].isArray()) {
        asDynamic[name] = dynamic::array();
      }
      asDynamic[name].push_back(leafToDynamic(leafPath, value, schemaCtx));
    } else {
      Path leafPath = path.getChild(leaf.first);
      asDynamic[leaf.first] =
          leafToDynamic(leafPath, leaf.second.value, schemaCtx);
    }
  }

  // complex children
  for (auto& child : entity.get_children()) {
    Path childPath = path.getChild(child.first);

    // list entry
    if (schemaCtx.isList(childPath.unkeyed())) {
      if (!asDynamic[child.second->yang_name].isArray()) {
        asDynamic[child.second->yang_name] = dynamic::array();
      }
      dynamic childAsDynamic =
          toDynamicRecurs(childPath, *child.second, schemaCtx);
      if (!childAsDynamic.empty()) {
        asDynamic[child.second->yang_name].push_back(childAsDynamic);
      }
    } else {
      asDynamic[child.first] =
          toDynamicRecurs(childPath, *child.second, schemaCtx);
    }

    // remove empty nodes
    if (asDynamic[child.first].empty()) {
      asDynamic.erase(child.first);
    }
  }

  return asDynamic;
}

dynamic BindingCodec::toDom(Path path, Entity& entity) {
  if (path.getLastSegment() != entity.get_segment_path()) {
    throw BindingSerializationException(
        entity,
        "Unable to serialize entity, path and entity are not pointing to the same node, path: " +
            path.str() + ", entity: " + entity.get_segment_path());
  }
  try {
    string key = path.prefixAllSegments().unkeyed().getLastSegment();
    if (schemaCtx.isList(path.unkeyed())) {
      return dynamic::object(
          key, dynamic::array(toDynamicRecurs(path, entity, schemaCtx)));
    } else {
      return dynamic::object(key, toDynamicRecurs(path, entity, schemaCtx));
    }
  } catch (runtime_error& e) {
    throw BindingSerializationException(
        entity,
        "Unable to serialize entity on path: " + path.str() +
            " due to: " + e.what());
  }
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

shared_ptr<Entity> BindingCodec::fromDom(
    const dynamic& payload,
    shared_ptr<Entity> pointer) {
  dynamic key = payload.items().begin()->first;
  dynamic value = payload.items().begin()->second;
  if (value.isArray()) {
    if (value.size() > 1) {
      throw BindingSerializationException(
          pointer,
          "Unable to parse entity, DOM contains an array with multiple entries: " +
              payload.asString());
    }
    // for lists, there should be an array in dynamic, however ydk cannot
    // process that so strip the array
    return decode(folly::toJson(dynamic::object(key, value[0])), pointer);
  }

  // Using json step to simplify parsing
  return decode(folly::toJson(payload), pointer);
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
