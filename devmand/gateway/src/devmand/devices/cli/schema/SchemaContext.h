// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/Model.h>
#include <devmand/devices/cli/schema/Path.h>
#include <libyang/libyang.h>

namespace devmand::devices::cli {

using devmand::devices::cli::Path;
using std::string;
using std::vector;

class SchemaContext {
 private:
  llly_ctx* ctx;
  string model = "NO_MODELS";
  llly_set* getNodes(Path path) const;

 public:
  // Special empty schema context, indicating no models are present
  static const SchemaContext NO_MODELS;

  SchemaContext(const Model& model);
  ~SchemaContext();
  SchemaContext(const SchemaContext&) = delete;
  SchemaContext& operator=(const SchemaContext&) = delete;
  SchemaContext(SchemaContext&&) = delete;
  SchemaContext& operator=(SchemaContext&&) = delete;

  llly_ctx* getLyContext() const;
  bool isPathValid(Path path) const;
  bool isList(Path p) const;
  bool isConfig(Path p) const;
  vector<string> getKeys(Path path) const;
  // return resolved type of a leaf
  vector<LLLY_DATA_TYPE> leafType(Path path) const;

  bool operator==(const SchemaContext& rhs) const;
  bool operator!=(const SchemaContext& rhs) const;

 private:
  SchemaContext(llly_ctx* _ctx) : ctx(_ctx){};
  static void loadModules(llly_ctx* context, const Model& _model);
};

class SchemaContextException : public runtime_error {
 public:
  SchemaContextException(const string& _msg) : runtime_error(_msg){};
};
} // namespace devmand::devices::cli
