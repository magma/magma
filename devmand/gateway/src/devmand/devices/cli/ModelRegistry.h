// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

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

class Model {
 public:
  Model() = delete;
  ~Model() = default;

 protected:
  explicit Model(std::string _dir) : dir(_dir) {}

 private:
  const string dir;

 public:
  const string& getDir() const {
    return dir;
  }
  bool operator<(const Model& x) const {
    return dir < x.dir;
  }

  static const Model OPENCONFIG_0_1_6;
  static const Model IETF_0_1_5;
};

class Bundle {
 public:
  explicit Bundle(const Model& model);
  Bundle() = delete;
  ~Bundle() = default;
  Bundle(const Bundle&) = delete;
  Bundle& operator=(const Bundle&) = delete;
  Bundle(Bundle&&) = delete;
  Bundle& operator=(Bundle&&) = delete;

 private:
  mutex lock; // A bundle is expected to be shared, protect it
  Repository repo;
  CodecServiceProvider codecServiceProvider;
  JsonSubtreeCodec jsonSubtreeCodec;

 public:
  std::string encode(Entity& entity);
  shared_ptr<Entity> decode(const string& payload, shared_ptr<Entity> pointer);
};

class ModelRegistry {
 public:
  ModelRegistry();
  ~ModelRegistry();
  ModelRegistry(const ModelRegistry&) = delete;
  ModelRegistry& operator=(const ModelRegistry&) = delete;
  ModelRegistry(ModelRegistry&&) = delete;
  ModelRegistry& operator=(ModelRegistry&&) = delete;

 public:
  Bundle& getBundle(const Model& model);
  ulong cacheSize();

 private:
  mutex lock; // A bundle is expected to be shared, protect it
  std::map<string, Bundle> cache;
};

class SerializationException : public std::exception {
 private:
  string msg;

 public:
  SerializationException(Entity& _entity, string _cause) {
    std::stringstream buffer;
    buffer << "Failed to encode: " << typeid(_entity).name() << " due to "
           << _cause;
    msg = buffer.str();
  };

  SerializationException(shared_ptr<Entity>& _entity, string _cause) {
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
