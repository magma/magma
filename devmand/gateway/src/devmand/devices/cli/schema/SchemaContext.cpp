// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/filesystem.hpp>
#include <devmand/devices/cli/schema/SchemaContext.h>

namespace devmand::devices::cli {

const SchemaContext SchemaContext::NO_MODELS(nullptr);

vector<string> SchemaContext::getKeys(Path path) const {
  if (not isPathValid(path)) {
    throw InvalidPathException(
        path.str(), "Not valid according to model from " + model);
  }
  if (not isList(path)) {
    throw InvalidPathException(
        path.str(), "Not a list according to model from " + model);
  }
  llly_set* pSet = getNodes(path);
  lllys_node* node = pSet->set.s[0];
  vector<string> result;

  if (node->nodetype == LLLYS_LIST) {
    auto* list = (lllys_node_list*)node;
    for (uint8_t i = 0; i < list->keys_size; i++) {
      result.emplace_back(string(list->keys[i]->name));
    }
  }
  llly_set_free(pSet);
  return result;
}

bool SchemaContext::isList(Path path) const {
  if (path == Path::ROOT) {
    return false;
  }

  llly_set* pSet = getNodes(path);
  auto result = pSet->set.s[0]->nodetype == LLLYS_LIST;
  llly_set_free(pSet);
  return result;
}

static vector<LLLY_DATA_TYPE> resolveType(lllys_type type) {
  if (-1 == type.base) {
    MLOG(MDEBUG) << "We have a problem";
  }
  if (LLLY_TYPE_DER == type.base) {
    MLOG(MDEBUG) << "We have a problem";
  }
  if (LLLY_TYPE_UNION == type.base) {
    vector<LLLY_DATA_TYPE> unionTypes;
    for (lllys_type* subtype = type.info.uni.types;
         subtype < type.info.uni.types + type.info.uni.count;
         ++subtype) {
      unionTypes.push_back(subtype->base);
    }
    return unionTypes;
  }

  // follow leafref to actual type
  if (LLLY_TYPE_LEAFREF == type.base) {
    while (type.info.lref.target) {
      type = type.info.lref.target->type;
    }
    return vector<LLLY_DATA_TYPE>{type.base};
  }
  return vector<LLLY_DATA_TYPE>{type.base};
}

vector<LLLY_DATA_TYPE> SchemaContext::leafType(Path path) const {
  llly_set* pSet = getNodes(path);
  lllys_node* node = pSet->set.s[0];
  vector<string> result;

  vector<LLLY_DATA_TYPE> resolvedTypes;
  if (node->nodetype == LLLYS_LEAF) {
    resolvedTypes = resolveType(((lllys_node_leaf*)node)->type);
  } else if (node->nodetype == LLLYS_LEAFLIST) {
    resolvedTypes = resolveType(((lllys_node_leaflist*)node)->type);
  } else {
    llly_set_free(pSet);
    throw InvalidPathException(
        path.str(),
        "Unable to lookup leaf type, path is not pointing to a leaf");
  }

  llly_set_free(pSet);
  return resolvedTypes;
}

bool SchemaContext::isConfig(Path path) const {
  if (path == Path::ROOT) {
    return true;
  }

  llly_set* pSet = getNodes(path);
  auto result = pSet->set.s[0]->flags & LLLYS_CONFIG_W;
  llly_set_free(pSet);
  return result;
}

bool SchemaContext::isPathValid(Path path) const {
  if (path == Path::ROOT) {
    return true;
  }

  void* schema =
      llly_path_data2schema(ctx, const_cast<char*>(path.str().c_str()));
  auto result = schema != nullptr;
  free(schema);
  return result;
}

SchemaContext::~SchemaContext() {
  llly_ctx_destroy(ctx, nullptr);
}

SchemaContext::SchemaContext(const Model& _model) : model(_model.getDir()) {
  // set extensions and user_types for non-YDK libyang
  setenv(
      "LLLIBYANG_EXTENSIONS_PLUGINS_DIR",
      LLLIBYANG_EXTENSIONS_PLUGINS_DIR,
      false);
  setenv(
      "LIBYANG_USER_TYPES_PLUGINS_DIR", LIBYANG_USER_TYPES_PLUGINS_DIR, false);
  ctx = llly_ctx_new(_model.getDir().c_str(), LLLY_CTX_ALLIMPLEMENTED);

  int modelCount = 0;
  int failedModelCount = 0;

  try {
    for (boost::filesystem::directory_entry& p :
         boost::filesystem::directory_iterator(
             boost::filesystem::path(_model.getDir()))) {
      if (!boost::filesystem::is_regular_file(p)) {
        continue;
      }
      if (!p.path().has_extension() or p.path().extension() != ".yang") {
        continue;
      }
      auto modelName = p.path().filename().replace_extension();
      modelCount++;
      if (!llly_ctx_load_module(ctx, modelName.string().c_str(), NULL)) {
        failedModelCount++;
        MLOG(MWARNING) << "Unable to parse model: " << p.path() << ". Ignoring";
      };
    }

    if (failedModelCount == modelCount) {
      throw SchemaContextException(
          "Unable to parse schema context from " + _model.getDir() +
          " due to: Failed to parse all models: " +
          to_string(failedModelCount));
    }
  } catch (boost::filesystem::filesystem_error& e) {
    throw SchemaContextException(
        "Unable to parse schema context from " + _model.getDir() +
        " due to: " + e.what());
  }
}

llly_set* SchemaContext::getNodes(Path path) const {
  if (not isPathValid(path)) {
    throw InvalidPathException(
        path.str(), "Not valid according to model from " + model);
  }
  char* schemaPath =
      llly_path_data2schema(ctx, const_cast<char*>(path.str().c_str()));
  Path schemaPathPrefixed(schemaPath);
  if (schemaPath != nullptr) {
    free(schemaPath);
  }
  llly_set* pSet = llly_ctx_find_path(
      ctx,
      const_cast<char*>(schemaPathPrefixed.prefixAllSegments().str().c_str()));

  if (pSet == nullptr || pSet->number != 1) {
    // shouldn't happen
    throw InvalidPathException(
        path.str(),
        "Unable to find schema node, but path seems valid according to model from " +
            model);
  }

  return pSet;
}

bool SchemaContext::operator==(const SchemaContext& rhs) const {
  return ctx == rhs.ctx;
}

bool SchemaContext::operator!=(const SchemaContext& rhs) const {
  return !(rhs == *this);
}
} // namespace devmand::devices::cli
