// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/Model.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <ydk/codec_provider.hpp>
#include <ydk/codec_service.hpp>
#include <ydk/json_subtree_codec.hpp>
#include <ydk/path_api.hpp>
#include <mutex>
#include <sstream>

namespace devmand {
namespace devices {
namespace cli {

using namespace std;
using namespace ydk;
using namespace ydk::path;
using namespace devmand::devices::cli;

class BindingCodec {
 public:
  explicit BindingCodec(Repository& repo, const string& schemaDir);
  BindingCodec() = delete;
  ~BindingCodec() = default;
  BindingCodec(const BindingCodec&) = delete;
  BindingCodec& operator=(const BindingCodec&) = delete;
  BindingCodec(BindingCodec&&) = delete;
  BindingCodec& operator=(BindingCodec&&) = delete;

 private:
  mutex lock; // A codec is expected to be shared, protect it
  CodecServiceProvider codecServiceProvider;
  JsonSubtreeCodec jsonSubtreeCodec;

 public:
  string encode(Entity& entity);
  shared_ptr<Entity> decode(const string& payload, shared_ptr<Entity> pointer);
};

class BindingContext {
 public:
  explicit BindingContext(const Model& model);
  BindingContext() = delete;
  ~BindingContext() = default;
  BindingContext(const BindingContext&) = delete;
  BindingContext& operator=(const BindingContext&) = delete;
  BindingContext(BindingContext&&) = delete;
  BindingContext& operator=(BindingContext&&) = delete;

 private:
  Repository repo;
  BindingCodec bindingCodec;

 public:
  BindingCodec& getCodec();
};

class BindingSerializationException : public exception {
 private:
  string msg;

 public:
  BindingSerializationException(Entity& _entity, string _cause) {
    std::stringstream buffer;
    buffer << "Failed to encode: " << typeid(_entity).name() << " due to "
           << _cause;
    msg = buffer.str();
  };

  BindingSerializationException(shared_ptr<Entity>& _entity, string _cause) {
    std::stringstream buffer;
    buffer << "Failed to decode: " << typeid(*_entity).name() << " due to "
           << _cause;
    msg = buffer.str();
  };

 public:
  const char* what() const throw() override {
    return msg.c_str();
  }
};

} // namespace cli
} // namespace devices
} // namespace devmand
